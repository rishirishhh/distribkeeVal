package web

import (
	"fmt"
	"net/http"

	"github.com/rishirishhh/distribkv/db"
)

// Server contains HTTP method handlers to be used for the database
type Server struct {
	db *db.Database
}

// NewServer creates a new instance with HTTP handlers to be used to get and set values
func NewServer(db *db.Database) *Server{
	return &Server{
		db:db,
	}
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	key := r.Form.Get("key")

	value, err := s.db.GetKey(key)
	if err != nil {
		fmt.Fprintf(w, "error= %v", err)

	}
	fmt.Fprintf(w, "Value = %q", value)
}

func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	err := s.db.SetKey(key, []byte(value))
	if err != nil {
		fmt.Fprintf(w, "Error = %v", err)
	}
}