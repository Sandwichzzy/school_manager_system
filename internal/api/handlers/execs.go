package handlers

import (
	"net/http"
)

func ExecsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("hello GET METHOD on Execs route"))
	case http.MethodPost:
		w.Write([]byte("hello POST METHOD on Execs route"))
	case http.MethodPut:
		w.Write([]byte("hello PUT METHOD on Execs route"))
	case http.MethodPatch:
		w.Write([]byte("hello PATCH METHOD on Execs route"))
	case http.MethodDelete:
		w.Write([]byte("hello DEL METHOD on Execs route"))
	}
	w.Write([]byte("hello execs route"))
}

// ExecsHandler handles requests to the /execs/ endpoint
// func ExecsHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "ExecsHandler is working!")
// }
