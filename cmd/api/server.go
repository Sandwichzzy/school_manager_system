package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	mw "github.com/Sandwichzzy/REST_API_GO/internal/api/middlewares"
)

type user struct {
	Name string `json:"name"`
	Age  string `json:"age"`
	City string `json:"city"`
}

//parse form data (necessary for x-www-form-urlencoded)
//err := r.ParseForm()
//if err != nil {
//	http.Error(w, "Error parsing form", http.StatusBadRequest)
//	return
//}
//fmt.Println("Form:", r.Form)
////prepare response data
//response := make(map[string]interface{})
//for key, value := range r.Form {
//	response[key] = value[0]
//}
//fmt.Println("Processed Response Map", response)
////RAW body
//body, err := io.ReadAll(r.Body)
//if err != nil {
//	return
//}
//defer r.Body.Close()
//fmt.Println("Raw body:", string(body))
//
//// if you expect json data, then unmarshal it
//var userInstance user
//err = json.Unmarshal(body, &userInstance)
//if err != nil {
//	return
//}
//fmt.Println("Unmarshaled JSON into an instance of user struct", userInstance)
//fmt.Println("Received user name as:", userInstance.Name)
//
////prepare response data
//response1 := make(map[string]interface{})
//for key, value := range r.Form {
//	response1[key] = value[0]
//}
//err = json.Unmarshal(body, &response1)
//if err != nil {
//	return
//}
//fmt.Println("Unmarshaled JSON into a map", response1)

//Access the request details
//fmt.Println("Body:", r.Body)
//fmt.Println("Form:", r.Form)
//fmt.Println("Header:", r.Header)
//fmt.Println("Context:", r.Context())
//fmt.Println("ContentLength:", r.ContentLength)
//fmt.Println("Host:", r.Host)
//fmt.Println("Method:", r.Method)
//fmt.Println("Proto:", r.Proto)
//fmt.Println("RemoteAddr:", r.RemoteAddr)
//fmt.Println("RequestURI:", r.RequestURI)
//fmt.Println("TLS:", r.TLS)
//fmt.Println("Trailer:", r.Trailer)
//fmt.Println("TransferEncoding:", r.TransferEncoding)
//fmt.Println("URL:", r.URL)
//fmt.Println("User Agent:", r.UserAgent())
//fmt.Println("Port:", r.URL.Port())

//teachers/{id}
//teachers?key=value&query=value2&sortby=email&sortorder=ASC
// 	fmt.Println(r.URL.Path)
// path := strings.TrimPrefix(r.URL.Path, "/teachers/")
// userID := strings.TrimSuffix(path, "/")

// fmt.Println("The ID is ", userID)

// fmt.Println("Query Params", r.URL.Query())
// queryParams := r.URL.Query()
// sortby := queryParams.Get("sortby")
// sortorder := queryParams.Get("sortorder")
// key := queryParams.Get("key")

// if sortorder == "" {
// 	sortorder = "DESC"
// }
// fmt.Printf("Sortby:%v , SortOrder: %v, Key: %v", sortby, sortorder, key)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Root Route"))
	fmt.Println("Hello Root Route")
}

func teachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//teachers/{id}
	//teachers?key=value&query=value2&sortby=email&sortorder=ASC
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

	//w.Write([]byte("hello Teachers route"))
	//fmt.Println("Hello Teachers Route")
}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("hello GET METHOD on students route"))
		fmt.Println("hello GET METHOD on students route")
	case http.MethodPost:
		w.Write([]byte("hello POST METHOD on students route"))
		fmt.Println("hello POST METHOD on students route")
	case http.MethodPut:
		w.Write([]byte("hello PUT METHOD on students route"))
		fmt.Println("hello PUT METHOD on students route")
	case http.MethodPatch:
		w.Write([]byte("hello PATCH METHOD on students route"))
		fmt.Println("hello PATCH METHOD on students route")
	case http.MethodDelete:
		w.Write([]byte("hello DEL METHOD on students route"))
		fmt.Println("hello DEL METHOD on students route")
	}

	w.Write([]byte("hello students route"))
	fmt.Println("Hello students Route")
}

func execsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("hello GET METHOD on Execs route"))
		fmt.Println("hello GET METHOD on Execs route")
	case http.MethodPost:
		w.Write([]byte("hello POST METHOD on Execs route"))
		fmt.Println("hello POST METHOD on Execs route")
	case http.MethodPut:
		w.Write([]byte("hello PUT METHOD on Execs route"))
		fmt.Println("hello PUT METHOD on Execs route")
	case http.MethodPatch:
		w.Write([]byte("hello PATCH METHOD on Execs route"))
		fmt.Println("hello PATCH METHOD on Execs route")
	case http.MethodDelete:
		w.Write([]byte("hello DEL METHOD on Execs route"))
		fmt.Println("hello DEL METHOD on Execs route")
	}

	w.Write([]byte("hello execs route"))
	fmt.Println("Hello execs Route")
}

func main() {
	port := ":3000"

	cert := "cert.pem"
	key := "key.pem"

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)

	mux.HandleFunc("/teachers/", teachersHandler)

	mux.HandleFunc("/students/", studentsHandler)

	mux.HandleFunc("/execs/", execsHandler)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	r1 := mw.NewRateLimiter(5, time.Minute)

	//create custom server
	server := &http.Server{
		Addr: port,
		// Handler:		middlewares.SecurityHeaders(mux),
		Handler:   r1.Middleware(mw.Comporession(mw.ResponseTimeMiddleware(mw.Cors(mw.SecurityHeaders(mux))))),
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is runnning on port " + port)
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error starting server:", err)
	}
}
