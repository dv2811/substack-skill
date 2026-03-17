package koyfin

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/goccy/go-json"
)

// AuthResponse record authentication response from Koyfin API
type AuthResponse struct {
	JWT          string `json:"jwt"`
	ID           string `json:"_id"`
	RefreshToken string `json:"refreshToken"`
}

func (c *Client) Login(session *Session, currentTimeStamp int64) error {
	// lock session before modifying
	session.Lock()
	defer session.Unlock()

	// if currentTimeStamp isn't within 10 minutes earlier or after next refresh - then ignore
	if currentTimeStamp < session.NextLogin {
		return nil
	}

	input := struct {
		Email         string `json:"email"`
		Password      string `json:"password"`
		RememberMe    bool   `json:"rememberMe"`
		SessionSource string `json:"sessionSource"`
	}{
		Email:         session.UserName,
		Password:      session.Password,
		SessionSource: "WEB",
		RememberMe:    true,
	}

	// get bytes buffer
	b := bufferPool.Get().(*bytes.Buffer)
	b.Reset()
	defer func() {
		bufferPool.Put(b)
	}()

	// request body buffer
	err := json.NewEncoder(b).Encode(input)
	if err != nil {
		return err
	}

	headers := map[string][]string{
		"User-Agent":      {"Mozilla/5.0 (X11; Linux x86_64; rv:146.0) Gecko/20100101 Firefox/146.0"},
		"Accept":          {"application/json, text/plain, */*"},
		"Accept-Language": {"en-US,en;q=0.5"},
		"Accept-Encoding": {"gzip, deflate, br, zstd"},
		"Content-Type":    {"application/json"},
		"Origin":          {"https://app.koyfin.com"},
		"Sec-GPC":         {"1"},
		"x-tab-id":        {"ur3V7E3hifT4El9GKesve"},
		"Connection":      {"keep-alive"},
		"Referer":         {"https://app.koyfin.com/"},
		"Sec-Fetch-Dest":  {"empty"},
		"Sec-Fetch-Mode":  {"cors"},
		"Sec-Fetch-Site":  {"same-site"},
	}

	// prepare request
	req, err := http.NewRequest("POST", "https://auth.koyfin.com/authentication/login", b)
	if err != nil {
		return err
	}
	req.Header = headers
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	// HTTP error handling
	if resp.StatusCode > 399 {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		return errors.Join(KoyfinAPIErr, errors.New(resp.Status))
	}

	authResponse := new(AuthResponse)
	err = json.NewDecoder(resp.Body).Decode(authResponse)
	if err != nil {
		return err
	}

	// make sure cookie header not empty
	if session.CookieMap == nil {
		session.CookieMap = map[string]string{}
	}

	// cookie header
	for _, ck := range resp.Cookies() {
		// remove expired cookies
		if ck.MaxAge < 0 || ck.Expires.Before(time.Now()) {
			delete(session.CookieMap, ck.Name)
			continue
		}
		session.CookieMap[ck.Name] = ck.Value
	}

	// clean up connection
	resp.Body.Close()
	// populate tokens
	session.AccessToken = authResponse.JWT
	session.RefreshToken = authResponse.RefreshToken

	// set timestamp
	session.LastUpdate = time.Now()
	currentTimeStamp = session.LastUpdate.Unix()
	session.NextRefresh = currentTimeStamp + KfRefreshInterval
	session.NextLogin = currentTimeStamp + KfLoginInterval
	return nil
}

// Refresh renews access token
func (c *Client) Refresh(session *Session, currentTimeStamp int64) error {
	// lock session before modifying
	session.Lock()
	defer session.Unlock()

	// double check
	if currentTimeStamp < session.NextRefresh {
		return nil
	}

	input := struct {
		RefreshToken  string `json:"refreshToken"`
		SessionSource string `json:"sessionSource"`
	}{
		RefreshToken:  session.RefreshToken,
		SessionSource: "WEB",
	}

	// get bytes buffer
	b := bufferPool.Get().(*bytes.Buffer)
	b.Reset()
	defer func() {
		bufferPool.Put(b)
	}()

	// request body buffer
	err := json.NewEncoder(b).Encode(input)
	if err != nil {
		return err
	}

	headers := map[string][]string{
		"User-Agent":      {"Mozilla/5.0 (X11; Linux x86_64; rv:146.0) Gecko/20100101 Firefox/146.0"},
		"Accept":          {"application/json, text/plain, */*"},
		"Accept-Language": {"en-US,en;q=0.5"},
		"Accept-Encoding": {"gzip, deflate, br, zstd"},
		"Content-Type":    {"application/json"},
		"Origin":          {"https://app.koyfin.com"},
		"Sec-GPC":         {"1"},
		"Connection":      {"keep-alive"},
		"Referer":         {"https://app.koyfin.com/"},
		"Sec-Fetch-Dest":  {"empty"},
		"Sec-Fetch-Mode":  {"cors"},
		"Sec-Fetch-Site":  {"same-site"},
	}

	// prepare request
	req, err := http.NewRequest("POST", "https://auth.koyfin.com/authentication/v2/token/refresh", b)
	if err != nil {
		return err
	}
	req.Header = headers
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	// HTTP error handling
	if resp.StatusCode > 399 {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		return errors.Join(KoyfinAPIErr, errors.New(resp.Status))
	}

	authResponse := new(AuthResponse)
	err = json.NewDecoder(resp.Body).Decode(authResponse)
	if err != nil {
		return err
	}

	// make sure cookie header not empty
	if session.CookieMap == nil {
		session.CookieMap = map[string]string{}
	}

	// cookie header
	for _, ck := range resp.Cookies() {
		// remove expired cookies
		if ck.MaxAge < 0 || ck.Expires.Before(time.Now()) {
			delete(session.CookieMap, ck.Name)
			continue
		}
		session.CookieMap[ck.Name] = ck.Value
	}

	// clean up connection
	resp.Body.Close()
	// populate tokens
	session.AccessToken = authResponse.JWT
	session.RefreshToken = authResponse.RefreshToken

	// set timestamp
	session.LastUpdate = time.Now()
	currentTimeStamp = session.LastUpdate.Unix()
	session.NextRefresh = currentTimeStamp + KfRefreshInterval
	session.NextLogin = currentTimeStamp + KfLoginInterval
	return nil
}

// ensureValidToken ensures access token is always valid and performs refresh or login when needed
func (c *Client) ensureValidToken(session *Session) error {
	session.RLock()
	// curent timestamp + overlapping period
	currentTimeStamp := time.Now().Unix() + RefreshBuffer
	nextLogin := session.NextLogin
	nextRefresh := session.NextRefresh
	session.RUnlock()

	// if login is outdated - refresh
	if currentTimeStamp >= nextLogin {
		return c.Login(session, currentTimeStamp)

		// refresh session
	} else if currentTimeStamp >= nextRefresh {
		return c.Refresh(session, currentTimeStamp)
	}
	return nil
}
