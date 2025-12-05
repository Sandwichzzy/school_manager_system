package handlers

import (
	"net/http"
)

func StudentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("hello GET METHOD on students route"))
	case http.MethodPost:
		w.Write([]byte("hello POST METHOD on students route"))
	case http.MethodPut:
		w.Write([]byte("hello PUT METHOD on students route"))
	case http.MethodPatch:
		w.Write([]byte("hello PATCH METHOD on students route"))
	case http.MethodDelete:
		w.Write([]byte("hello DEL METHOD on students route"))
	}
	w.Write([]byte("hello students route"))
}
