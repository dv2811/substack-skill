package koyfin

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"entext-applications/internal/models"
	"github.com/goccy/go-json"
)

const (
	KfRefreshInterval int64 = 3600
	KfLoginInterval   int64 = 3600 * 720
	RefreshBuffer     int64 = 600 // 10-min overlapping period
)

var bufferPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

type Session struct {
	UserName     string            `json:"email"`
	Password     string            `json:"password"`
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	CookieMap    map[string]string `json:"ckv"`
	NextRefresh  int64             `json:"next_refresh_after"`
	NextLogin    int64             `json:"must_login_after"`
	LastUpdate   time.Time         `json:"last_update"`
	sync.RWMutex
}

func NewSession(email, password string) *Session {
	s := &Session{
		UserName:  email,
		Password:  password,
		CookieMap: map[string]string{},
	}
	return s
}

func NewSessionFromFile(path string) (*Session, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	s := new(Session)
	err = json.NewDecoder(file).Decode(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Session) SaveToFile(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return json.NewEncoder(file).Encode(s)
}

// save data to []byte format
func (s *Session) Load(data []byte) error {
	newSession := new(Session)
	err := json.Unmarshal(data, newSession)
	if err != nil {
		return err
	}

	// issue with existing field values, create new Session and marshal instead
	_ = s.Copy(newSession)
	return nil
}

// save data to []byte format
func (s *Session) Save() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Session) UpdatedAt() time.Time {
	return s.LastUpdate
}

// Copy perform shallow copy of source object if the type matches
func (s *Session) Copy(source any) error {
	data, valid := source.(*Session)
	if !valid {
		return models.ErrSessionTypeNotMatch
	}
	if data != nil {
		s.UserName = data.UserName
		s.Password = data.Password
		s.AccessToken = data.AccessToken
		s.RefreshToken = data.RefreshToken
		s.CookieMap = data.CookieMap
		s.NextRefresh = data.NextRefresh
		s.NextLogin = data.NextLogin
		s.LastUpdate = data.LastUpdate
	}
	return nil
}

// AuthorizeHeader populates request header before being passed to request
func (s *Session) AuthorizeHeader(header http.Header) {
	header.Set("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))
	var b strings.Builder
	limit := len(s.CookieMap)
	for k, v := range s.CookieMap {
		b.WriteString(fmt.Sprintf("%s=%s", k, v))
		limit--
		if limit > 0 {
			b.WriteString("; ")
		}
	}

	// set cookie
	if b.Len() > 0 {
		header.Set("Cookies", b.String())
	}
}
