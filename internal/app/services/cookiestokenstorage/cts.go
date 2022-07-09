package cookiestokenstorage

import "net/http"

type CookiesTokenStorage struct {
	accessToken  string
	refreshToken string
	r            *http.Request
}

func (s *CookiesTokenStorage) SetParams(params map[string]interface{}) {
	r, ok := params["request"]
	if ok {
		s.r = r.(*http.Request)
	}
}

func (s *CookiesTokenStorage) GetTokens() (string, string, error) {
	if s.accessToken != "" && s.refreshToken != "" {
		return s.accessToken, s.refreshToken, nil
	}
	atCookie, err := s.r.Cookie("access_token")
	if err != nil {
		return "", "", err
	}
	rtCookie, err := s.r.Cookie("refresh_token")
	if err != nil {
		return "", "", err
	}
	return atCookie.Value, rtCookie.Value, nil
}

func (s *CookiesTokenStorage) SetTokens(accessToken, refreshToken string) error {
	s.accessToken = accessToken
	s.refreshToken = refreshToken
	return nil
}
