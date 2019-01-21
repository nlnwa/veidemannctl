package connection

import (
	"fmt"
	"golang.org/x/net/context"
	"net/http"
)

type LoginResponseServer struct {
	server  *http.Server
	channel chan string
}

func (s *LoginResponseServer) handleLoginResponse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		fmt.Fprintf(w, "<html><body><h1>%s</h1></body></html>", "Window can safely be closed")
		go s.setResponse(query.Get("code"), query.Get("state"))
	}
}

func (s *LoginResponseServer) setResponse(code string, state string) {
	s.channel <- code
	s.channel <- state
}

func (a *auth) listen() (string, string) {
	s := &LoginResponseServer{
		channel: make(chan string),
		server: &http.Server{
			Addr: ":9876",
		},
	}

	http.HandleFunc("/", s.handleLoginResponse())
	go s.server.ListenAndServe()
	//go log.Fatal(s.server.ListenAndServe())
	code, state := <- s.channel, <- s.channel
	go s.server.Shutdown(context.Background())
	return code, state
}
