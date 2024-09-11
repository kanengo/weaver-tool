package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/ServiceWeaver/weaver"

	"html"
	"log"
	"net/http"
	"time"
)

//go:embed index.html
var indexHtml string

func main() {
	if err := weaver.Run(context.Background(), serve); err != nil {
		log.Fatal(err)
	}
}

type server struct {
	weaver.Implements[weaver.Main]
	reverser weaver.Ref[Reverser]
	lis      weaver.Listener `weaver:"reverser"`
}

func serve(ctx context.Context, s *server) error {
	// Setup the HTTP handler.
	var mux http.ServeMux
	mux.Handle("/", weaver.InstrumentHandlerFunc("root", s.handleRoot))
	mux.Handle("/reverse", weaver.InstrumentHandlerFunc("reverse", s.handleReverse))
	mux.HandleFunc(weaver.HealthzURL, weaver.HealthzHandler)
	s.Logger(ctx).Info("Reverser server running", "address", s.lis)
	//连接redis

	go func() {
		for {
			ret, err := s.reverser.Get().Reverse(ctx, "!akeel, ih")
			if err != nil {
				s.Logger(ctx).Error(fmt.Sprintf("self reverse err: %w\n", err))
			} else {
				s.Logger(ctx).Info(fmt.Sprintf("self reverse: %s\n", ret))
			}

			time.Sleep(10 * time.Second)
		}
	}()
	return http.Serve(s.lis, &mux)
}

// Handle requests to the "/" endpoint.
func (s *server) handleRoot(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, indexHtml)
}

// Handle requests to the "/reverse?s=<string>" endpoint.
func (s *server) handleReverse(w http.ResponseWriter, r *http.Request) {
	reversed, err := s.reverser.Get().Reverse(r.Context(), r.URL.Query().Get("s"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, html.EscapeString(reversed))
}
