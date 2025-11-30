package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := ":3000"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Root Route"))
		fmt.Println("Hello Root Route")
	})

	http.HandleFunc("/teachers", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method)
		switch r.Method {
		case http.MethodGet:
			w.Write([]byte("hello GET METHOD on Teachers route"))
			fmt.Println("hello GET METHOD on Teachers route")
		case http.MethodPost:
			w.Write([]byte("hello POST METHOD on Teachers route"))
			fmt.Println("hello POST METHOD on Teachers route")
		case http.MethodPut:
			w.Write([]byte("hello PUT METHOD on Teachers route"))
			fmt.Println("hello PUT METHOD on Teachers route")
		case http.MethodPatch:
			w.Write([]byte("hello PATCH METHOD on Teachers route"))
			fmt.Println("hello PATCH METHOD on Teachers route")
		case http.MethodDelete:
			w.Write([]byte("hello DEL METHOD on Teachers route"))
			fmt.Println("hello DEL METHOD on Teachers route")
		}

		w.Write([]byte("hello Teachers route"))
		fmt.Println("Hello Teachers Route")
	})

	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello students route"))
		fmt.Println("Hello students Route")
	})

	http.HandleFunc("/execs", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello execs route"))
		fmt.Println("Hello execs Route")
	})

	fmt.Println("Server is runnning on port " + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("Error starting server:", err)
	}
}
