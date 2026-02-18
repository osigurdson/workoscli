package workoscli

import (
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

type WorkOSToken struct {
	RefreshToken string    `json:"refresh_token"`
	AccessToken  string    `json:"access_token"`
	Expiry       time.Time `json:"expiry"`
}

type WorkOSTokenMgr struct {
	saveTokenFn func(WorkOSToken) error
	loadTokenFn func() (WorkOSToken, error)
}

func NewWorkOSTokenMgr(
	saveTokenFn func(WorkOSToken) error,
	loadTokenFn func() (WorkOSToken, error),
) (WorkOSTokenMgr, error) {
	if saveTokenFn == nil {
		return WorkOSTokenMgr{}, fmt.Errorf("saveTokenFn required")
	}

	if loadTokenFn == nil {
		return WorkOSTokenMgr{}, fmt.Errorf("loadTokenFn required")
	}

	return WorkOSTokenMgr{
		saveTokenFn: saveTokenFn,
		loadTokenFn: loadTokenFn,
	}, nil
}

func (c WorkOSToken) toOauthToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  c.AccessToken,
		RefreshToken: c.RefreshToken,
		Expiry:       c.Expiry,
	}
}
