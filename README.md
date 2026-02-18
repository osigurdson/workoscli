# workoscli

A small Go library for CLI-style OAuth login with WorkOS AuthKit, including:

- PKCE + state generation
- Local callback server (`/callback`) to capture the authorization code
- Token exchange and refresh handling
- `http.Client` transport that injects `Authorization: Bearer <token>`

## Requirements

- Go toolchain compatible with `go.mod` (`go 1.25.7`)
- A WorkOS AuthKit app/client configured with a localhost redirect URL

## Install

```bash
go get github.com/osigurdson/workoscli
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "os/exec"

    "github.com/osigurdson/workoscli"
)

func main() {
    ctx := context.Background()

    conf := workoscli.WorkOSConf{
        AuthKitURI:   "https://your-org.authkit.app",
        RedirectIP:   "127.0.0.1",
        RedirectPort: 21222,
        ClientID:     "client_xxx",
        Scopes:       []string{"openid", "profile", "email", "offline_access"},
    }

    var token *workoscli.WorkOSToken

    mgr, err := workoscli.NewWorkOSTokenMgr(
        func(newToken workoscli.WorkOSToken) error {
            token = &newToken
            return nil
        },
        func() (workoscli.WorkOSToken, error) {
            if token == nil {
                return workoscli.WorkOSToken{}, fmt.Errorf("token not found")
            }
            return *token, nil
        },
    )
    if err != nil {
        panic(err)
    }

    browser := func(url string) {
        _ = exec.Command("xdg-open", url).Start()
    }

    cli := workoscli.NewWorkOSCli(conf, browser, mgr)

    if err := cli.Login(ctx); err != nil {
        panic(err)
    }

    httpClient := cli.NewHttpClient(ctx)
    _ = httpClient
}
```

## How It Works

1. `Login` creates PKCE verifier/challenge + random state.
2. Browser is opened to WorkOS `/oauth2/authorize`.
3. A local HTTP server listens on `http://<RedirectIP>:<RedirectPort>/callback`.
4. On callback, code is exchanged at `/oauth2/token`.
5. Access/refresh token is saved via your `WorkOSTokenMgr` callbacks.
6. `NewHttpClient` returns a client with automatic bearer auth and refresh.

## Public API

- `type WorkOSConf`
- `func NewWorkOSTokenMgr(saveFn, loadFn) (WorkOSTokenMgr, error)`
- `func NewWorkOSCli(conf, browserFn, tokenMgr) *WorkOSCli`
- `func (c *WorkOSCli) Login(ctx context.Context) error`
- `func (c *WorkOSCli) NewHttpClient(ctx context.Context) *http.Client`

## Demo

Run:

```bash
go run ./cmd/demo
```

The demo currently contains hardcoded staging values in `cmd/demo/main.go`; replace them with your own WorkOS/AuthKit configuration.

## License

MIT (see `LICENSE`).
