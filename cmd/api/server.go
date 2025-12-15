package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	mw "github.com/Sandwichzzy/REST_API_GO/internal/api/middlewares"
	"github.com/Sandwichzzy/REST_API_GO/internal/api/router"
	"github.com/Sandwichzzy/REST_API_GO/internal/repository/sqlconnect"
	"github.com/joho/godotenv"
)

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

func main() {

	err := godotenv.Load()
	if err != nil {
		return
	}

	_, err = sqlconnect.ConnectDb()
	if err != nil {
		log.Fatalln("Database connection error:", err)
		return
	}

	port := os.Getenv("API_PORT")

	cert := "cert.pem"
	key := "key.pem"

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// r1 := mw.NewRateLimiter(5, time.Minute)

	// hppOptions := mw.HPPOptions{
	// 	CheckQuery:                  true,
	// 	CheckBody:                   true,
	// 	CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
	// 	Whitelist:                   []string{"sortBy", "sortOrder", "page", "limit", "name", "age", "class"},
	// }

	// secureMux := mw.Cors(r1.Middleware(mw.ResponseTimeMiddleware(mw.SecurityHeaders(mw.Compression(mw.Hpp(hppOptions)(mux))))))
	// secureMux := utils.ApplyMiddlewares(mux, mw.Hpp(hppOptions), mw.Compression, mw.SecurityHeaders, mw.ResponseTimeMiddleware, r1.Middleware, mw.Cors)
	router := router.MainRouter()
	secureMux := mw.SecurityHeaders(router)
	//create custom server
	server := &http.Server{
		Addr: port,
		// Handler:		middlewares.SecurityHeaders(mux),
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is runnning on port " + port)
	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error starting server:", err)
	}
}
