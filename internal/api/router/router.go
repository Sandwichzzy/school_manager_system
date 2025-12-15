package router

import (
	"net/http"

	"github.com/Sandwichzzy/REST_API_GO/internal/api/handlers"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.RootHandler)

	mux.HandleFunc("GET /teachers", handlers.GetTeachersHandler)
	mux.HandleFunc("POST /teachers", handlers.AddTeachersHandler)
	mux.HandleFunc("PATCH /teachers", handlers.PatchTeachersHandler)
	mux.HandleFunc("DELETE /teachers", handlers.DeleteTeachersHandler)

	mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeacherHandler)
	mux.HandleFunc("GET /teachers/{id}", handlers.GetOneTeacherHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteOneTeacherHandler)
	mux.HandleFunc("PATCH /teachers/{id}", handlers.PatchOneTeacherHandler)

	mux.HandleFunc("GET /teachers/{id}/students", handlers.GetStudentsByTeacherIdHandler)
	mux.HandleFunc("GET /teachers/{id}/studentcount", handlers.GetStudentsCountByTeacherIdHandler)

	mux.HandleFunc("GET /students", handlers.GetStudentsHandler)
	mux.HandleFunc("POST /students", handlers.AddStudentsHandler)
	mux.HandleFunc("PATCH /students", handlers.PatchStudentsHandler)
	mux.HandleFunc("DELETE /students", handlers.DeleteStudentsHandler)

	mux.HandleFunc("PUT /students/{id}", handlers.UpdateStudentHandler)
	mux.HandleFunc("GET /students/{id}", handlers.GetOneStudentHandler)
	mux.HandleFunc("DELETE /students/{id}", handlers.DeleteOneStudentHandler)
	mux.HandleFunc("PATCH /students/{id}", handlers.PatchOneStudentHandler)

	mux.HandleFunc("/execs/", handlers.ExecsHandler)

	return mux
}
