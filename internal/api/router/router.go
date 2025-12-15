package router

import (
	"net/http"
)

func MainRouter() *http.ServeMux {

	tRouter := teachersRouter()
	sRouter := studentsRouter()

	tRouter.Handle("/", sRouter)
	return tRouter

	// mux := http.NewServeMux()
	// mux.HandleFunc("/", handlers.RootHandler)
	// mux.HandleFunc("/execs/", handlers.ExecsHandler)

	// return mux
}
