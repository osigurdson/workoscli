package workoscli

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type server struct {
	srv    *http.Server
	errCh  chan error
	codeCh chan string
}

func startServer(state string, conf WorkOSConf) *server {
	errCh := make(chan error, 1)
	codeCh := make(chan string, 1)
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(w, "state mismatch", http.StatusBadRequest)
			errCh <- fmt.Errorf("state mismatch")
			return
		}
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "missing code", http.StatusBadRequest)
			errCh <- fmt.Errorf("missing code")
			return
		}
		w.Write([]byte("Login complete. You can close this tab"))
		codeCh <- code
	})

	srv := &http.Server{Addr: conf.ListenAddr(), Handler: mux}
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			if err != http.ErrServerClosed {
				errCh <- err
			}
		}
	}()

	return &server{
		srv:    srv,
		errCh:  errCh,
		codeCh: codeCh,
	}
}

func (s *server) waitForCode(ctx context.Context) (string, error) {
	select {
	case code := <-s.codeCh:
		return code, nil
	case err := <-s.errCh:
		return "", err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (s *server) shutdown(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	_ = s.srv.Shutdown(shutdownCtx)
	defer cancel()
}
