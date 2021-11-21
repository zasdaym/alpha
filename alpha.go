package alpha

import (
	_ "embed"
	"html/template"
	"net/http"
	"sync"
)

// Server provides aggregate from Clients.
type Server struct {
	mu    sync.Mutex
	count map[string]int
}

// NewServer creates a new Server.
func NewServer() *Server {
	return &Server{
		count: make(map[string]int),
	}
}

// Increment increments count for given clientID.
func (s *Server) Increment(clientID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.count[clientID]++
}

func (s *Server) HandleIncrement() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		clientID := r.FormValue("client-id")
		s.Increment(clientID)
	}
}

//go:embed "views/index.html"
var indexHTML string

func (s *Server) HandleIndex() (http.HandlerFunc, error) {
	tmpl, err := template.New("index").Parse(indexHTML)
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, s.count); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}, nil
}
