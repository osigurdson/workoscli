package workoscli

import (
	"net/http"
)

type AuthTransport struct {
	Base http.RoundTripper
	Auth *WorkOSCli
}

func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}

	// Get token
	tok, err := t.Auth.getValidToken(req.Context(), false)
	if err != nil {
		return nil, err
	}

	// Clone request
	req2 := req.Clone(req.Context())
	req2.Header.Set("Authorization", "Bearer "+tok.AccessToken)
	resp, err := base.RoundTrip(req2)
	if err != nil {
		return nil, err
	}

	// If not 401, we are done
	if resp.StatusCode != http.StatusUnauthorized {
		return resp, nil
	}

	// Try explicit refresh with the refresh token
	newTok, err := t.Auth.getValidToken(req.Context(), true)
	if err != nil {
		return nil, err
	}
	retry := req.Clone(req.Context())
	retry.Header.Set("Authorization", "Bearer "+newTok.AccessToken)

	return base.RoundTrip(retry)
}
