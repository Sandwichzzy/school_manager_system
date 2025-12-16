package router

import (
	"net/http"
)

func MainRouter() *http.ServeMux {

	eRouter := execsRouter()
	tRouter := teachersRouter()
	sRouter := studentsRouter()

	tRouter.Handle("/", sRouter)
	sRouter.Handle("/", eRouter)

	return tRouter
}
