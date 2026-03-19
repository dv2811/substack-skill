package substack

import (
	"bytes"
	"errors"
	// "fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"entext-applications/internal/models"
	"github.com/goccy/go-json"
)

var (
	ErrSubstackHTTPFailed = errors.New(`Substack connection failure`)
	ErrNotLoggedIn        = errors.New("not logged in")
	ErrSessionExpired     = errors.New("session expired, login required")
	EmptyReadingList      = errors.New("empty reading")
)

// substack login session
type Session struct {
	// User email will be used to trigger password less authentication flow
	Email      string                  `json:"user_email"`
	Headers    http.Header             `json:"headers"`
	Cookies    map[string]*http.Cookie `json:"ckv"`
	LastUpdate time.Time               `json:"last_update"`
	Expiry     time.Time               `json:"expiry"`
	sync.RWMutex
}

// save data to []byte format
func (s *Session) Load(data []byte) error {
	newSession := new(Session)
	err := json.Unmarshal(data, newSession)
	if err != nil {
		return err
	}

	// Copy values while respecting mutex
	s.Lock()
	defer s.Unlock()
	newSession.Lock()
	defer newSession.Unlock()

	s.Email = newSession.Email
	s.Headers = newSession.Headers
	s.Cookies = newSession.Cookies
	s.LastUpdate = newSession.LastUpdate
	s.Expiry = newSession.Expiry

	return nil
}

// save data to []byte format
func (s *Session) Save() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Session) UpdatedAt() time.Time {
	return s.LastUpdate
}

// Copy copies the source session to the current session
func (s *Session) Copy(source any) error {
	src, ok := source.(*Session)
	if !ok {
		return models.ErrSessionTypeNotMatch
	}

	s.Lock()
	defer s.Unlock()
	src.RLock()
	defer src.RUnlock()

	s.Email = src.Email
	s.Headers = src.Headers
	s.LastUpdate = src.LastUpdate
	s.Expiry = src.Expiry

	// Copy cookies map
	s.Cookies = make(map[string]*http.Cookie, len(src.Cookies))
	for k, v := range src.Cookies {
		s.Cookies[k] = v
	}

	return nil
}

// AuthenticatedURLFlow generates Substack session from authenticated login link
func (s *Session) LoadFromFile(content io.Reader) error {
	s.Lock()
	defer s.Unlock()

	return json.NewDecoder(content).Decode(s)
}

// LoadFromResponse populates session value using http.Response data
func (s *Session) LoadFromResponse(resp *http.Response) error {
	s.Lock()
	defer s.Unlock()

	// prevent nil dereference
	if s.Cookies == nil {
		s.Cookies = map[string]*http.Cookie{}
	}

	now := time.Now()
	for _, c := range resp.Cookies() {
		// removing expired cookies
		if c.MaxAge < 0 || c.Expires.Before(now) {
			delete(s.Cookies, c.Name)
			continue
		}
		s.Cookies[c.Name] = c
		// important: session ID for substack
		if c.Name == "substack.sid" {
			s.Expiry = c.Expires
		}
	}
	s.LastUpdate = now
	return nil
}

// AuthorizeHeaders populates header value for authorization and checks session expiry
func (s *Session) AuthorizeHeaders(headers http.Header) error {
	s.RLock()
	defer s.RUnlock()

	// Check if session has expired
	now := time.Now()
	if !s.Expiry.IsZero() && s.Expiry.Before(now) {
		return ErrSessionExpired
	}

	// add s.Headers to authorization headers - overwrite over addition
	for k, v := range s.Headers {
		headers[k] = v
	}

	// set cookies
	buf, _ := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	for _, c := range s.Cookies {
		if c.MaxAge < 0 || c.Expires.Before(now) {
			continue
		}
		buf.WriteString(c.Name)
		buf.WriteByte('=')
		buf.WriteString(c.Value)
		buf.WriteString("; ")
	}

	// write cookie string
	if buf.Len() > 0 {
		headers.Set("Cookie", buf.String()[:buf.Len()-2])
	}

	return nil
}

// Refresh checks if the session needs renewal and updates it using a fresh response
// This method checks if current time passes the mid-point between Expiry and LastUpdate
// This method assumes that the response contains updated session information
func (s *Session) Refresh(resp *http.Response) error {
	// Use Rlock() / RUnlock() for checking refresh need
	s.RLock()
	now := time.Now()

	// Check if renewal is needed based on the midpoint condition
	var needsRenewal bool
	if s.Expiry.Before(now) {
		needsRenewal = true
	} else {
		duration := s.Expiry.Sub(s.LastUpdate)
		midpoint := s.LastUpdate.Add(duration / 2)
		needsRenewal = now.After(midpoint)
	}

	s.RUnlock()

	// Return early if refresh not needed
	if !needsRenewal {
		return nil
	}

	// Call LoadFromResponse method which handles its own locking
	return s.LoadFromResponse(resp)
}

// NewSessionFromFile creates a new session by loading from a file path
func NewSessionFromFile(path string) (*Session, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	session := &Session{}
	if err := session.LoadFromFile(file); err != nil {
		return nil, err
	}

	return session, nil
}
