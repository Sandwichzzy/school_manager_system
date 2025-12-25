package main

import (
	"crypto/tls"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	mw "github.com/Sandwichzzy/REST_API_GO/internal/api/middlewares"
	"github.com/Sandwichzzy/REST_API_GO/internal/api/router"
	"github.com/Sandwichzzy/REST_API_GO/pkg/utils"
	"github.com/joho/godotenv"
)

//go:embed .env
var envFile embed.FS

func loadEnvFromEmbeddedFile() {
	//read the embedded .env file
	content, err := envFile.ReadFile(".env")
	if err != nil {
		log.Fatalf("Error reading embedded .env file:%v", err)
	}
	//Create a temporary file to store the .env content
	tmpFile, err := os.CreateTemp("", ".env")
	if err != nil {
		log.Fatalf("Error creating temp .env file:%v", err)
	}
	defer os.Remove(tmpFile.Name())

	//write content of the embedded .env file to the temp file
	_, err = tmpFile.Write(content)
	if err != nil {
		log.Fatalf("Error writing to temp .env file:%v", err)
	}
	err = tmpFile.Close()
	if err != nil {
		log.Fatalf("Error closing temp .env file:%v", err)
	}

	//load the env variables from the temp file
	err = godotenv.Load(tmpFile.Name())
	if err != nil {
		log.Fatalf("Error loading env variables from embedded .env file:%v", err)
	}
}

func main() {
	// only in production, for running source code
	// err := godotenv.Load()
	// if err != nil {
	// 	return
	// }

	//load env variables from the embedded .env file
	loadEnvFromEmbeddedFile()

	fmt.Println("Environment variables CERT_FILE:", os.Getenv("CERT_FILE"))

	port := os.Getenv("API_PORT")

	// cert := "cert.pem"
	// key := "key.pem"

	cert := os.Getenv("CERT_FILE")
	key := os.Getenv("KEY_FILE")

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS10,
	}

	r1 := mw.NewRateLimiter(5, time.Minute)

	hppOptions := mw.HPPOptions{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		Whitelist:                   []string{"sortBy", "sortOrder", "page", "limit", "name", "age", "class"},
	}

	// secureMux := mw.Cors(r1.Middleware(mw.ResponseTimeMiddleware(mw.SecurityHeaders(mw.Compression(mw.Hpp(hppOptions)(mux))))))
	// secureMux := mw.SecurityHeaders(router)
	// secureMux := mw.XSSMiddleware(router)
	// secureMux := jwtMiddleware(mw.SecurityHeaders(router))
	router := router.MainRouter()
	jwtMiddleware := mw.MiddlewareExcludePaths(mw.JWTMiddleware, "/execs/login", "/execs/forgotpassword", "/execs/resetpassword/reset/")
	secureMux := utils.ApplyMiddlewares(router, mw.SecurityHeaders, mw.Compression, mw.Hpp(hppOptions), mw.XSSMiddleware, jwtMiddleware, mw.ResponseTimeMiddleware, r1.Middleware, mw.Cors)
	//create custom server
	server := &http.Server{
		Addr: port,
		// Handler:		middlewares.SecurityHeaders(mux),
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is runnning on port " + port)
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error starting server:", err)
	}
}
