package main

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/osigurdson/workoscli"
)

func main() {
	ctx := context.Background()
	/*
		original := workoscligo.WorkOSConf{
			AuthKitURI:   "https://your-org.authkit.app",
			RedirectIP:   "127.0.0.1",
			RedirectPort: 21222,
			ClientID:     "client_123",
			Scopes:       []string{"openid", "profile", "email", "offline_access"},
		}
	*/

	conf := workoscli.WorkOSConf{
		AuthKitURI:   "https://surprising-cliff-63-staging.authkit.app",
		RedirectIP:   "127.0.0.1",
		RedirectPort: 21222,
		ClientID:     "client_01KHQ0D7B6Y3J3S7AXH9RQYX85",
		Scopes:       []string{"openid", "profile", "email", "offline_access"},
	}

	var token *workoscli.WorkOSToken

	mgr, err := workoscli.NewWorkOSTokenMgr(
		func(newToken workoscli.WorkOSToken) error {
			token = &newToken
			fmt.Printf("token: %+v\n", token)
			return nil
		},
		func() (workoscli.WorkOSToken, error) {
			if token == nil {
				return workoscli.WorkOSToken{}, fmt.Errorf("WorkOS token not found")
			}
			return *token, nil
		},
	)

	browser := func(url string) {
		exec.Command("xdg-open", url).Start()
	}

	workosCli := workoscli.NewWorkOSCli(
		conf,
		browser,
		mgr,
	)

	err = workosCli.Login(ctx)
	if err != nil {
		panic(err)
	}

	client := workosCli.NewHttpClient(ctx)
	res, err := client.Get("http://localhost:8787/api/me")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	fmt.Println(string(body))
}
