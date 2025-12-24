package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Sandwichzzy/REST_API_GO/pkg/utils"
	"github.com/microcosm-cc/bluemonday"
)

func XSSMiddleware(next http.Handler) http.Handler {
	fmt.Println("******Init XSS Middleware******")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("XSS Middleware Ran")
		//sanitize the URL path
		//URL: http://example.com/users/123/profile
		//r.URL.Path = "/users/123/profile"
		sanitizedPath, err := clean(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("Original Path:", r.URL.Path)
		fmt.Println("Sanitized Path:", sanitizedPath)

		//sanitize query parameters
		// URL: http://example.com/search?q=golang&page=2&sort=desc
		// query := r.URL.Query()
		// q := query.Get("q")        // "golang"
		// page := query.Get("page")  // "2"
		// sort := query.Get("sort")  // "desc
		params := r.URL.Query()
		sanitizedQuery := make(map[string][]string)
		for key, values := range params {
			sanitizedKey, err := clean(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var sanitizedValues []string
			for _, value := range values {
				cleanValue, err := clean(value)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				sanitizedValues = append(sanitizedValues, cleanValue.(string))
			}
			sanitizedQuery[sanitizedKey.(string)] = sanitizedValues
			fmt.Printf("Original Query %s:%s\n", key, strings.Join(values, ","))
			fmt.Printf("Sanitized Query %s:%s\n", sanitizedKey, strings.Join(sanitizedValues, ","))
		}
		r.URL.Path = sanitizedPath.(string)
		//replace original query with sanitized query
		r.URL.RawQuery = url.Values(sanitizedQuery).Encode()
		fmt.Println("Updated URL:", r.URL.String())

		//Sanitize request Body
		if r.Header.Get("Content-Type") == "application/json" {
			if r.Body != nil {
				bodyBytes, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, utils.ErrorHandler(err, "Error reading request body").Error(), http.StatusInternalServerError)
					return
				}
				bodyString := strings.TrimSpace(string(bodyBytes))
				fmt.Println("Original Body:", bodyString)

				//reset the request body
				r.Body = io.NopCloser(bytes.NewReader([]byte(bodyString)))
				if len(bodyString) > 0 {
					var inputData interface{}
					err := json.NewDecoder(bytes.NewReader([]byte(bodyString))).Decode(&inputData)
					if err != nil {
						http.Error(w, utils.ErrorHandler(err, "Error decoding JSON").Error(), http.StatusBadRequest)
						return
					}
					fmt.Println("Original JSON data:", inputData)
					//sanitize the json data
					sanitizedData, err := clean(inputData)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					fmt.Println("Sanitized JSON data:", sanitizedData)

					//Marshal sanitized data back to json
					sanitizedBody, err := json.Marshal(sanitizedData)
					if err != nil {
						http.Error(w, utils.ErrorHandler(err, "Error sanitizing body").Error(), http.StatusBadRequest)
						return
					}
					r.Body = io.NopCloser(bytes.NewReader(sanitizedBody))
					fmt.Println("Sanitized Body:", string(sanitizedBody))

				} else {
					log.Println("Request body is empty")
				}
			} else {
				log.Println("Request body is nil")
			}
		} else if r.Header.Get("Content-Type") != "" {
			log.Printf("received request with unsupported content-type:%s. Expected application/json.\n", r.Header.Get("Content-Type"))
			http.Error(w, "Unsupported Content-Type. Expected application/json", http.StatusUnsupportedMediaType)
			return
		}

		next.ServeHTTP(w, r)
		fmt.Println("Sending response from XSS Middleware Ran")
	})
}

// clean sanitizes input data to prevent XSS attacks
func clean(data interface{}) (interface{}, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			v[key] = sanitizeValue(value)
		}
		return v, nil
	case []interface{}:
		for i, value := range v {
			v[i] = sanitizeValue(value)
		}
		return v, nil
	case string:
		return sanitizeString(v), nil
	default:
		return nil, utils.ErrorHandler(fmt.Errorf("unsupported type:%T", data), fmt.Sprintf("unsupported type:%T", data))
	}
}

func sanitizeValue(data interface{}) interface{} {
	switch v := data.(type) {
	case string:
		return sanitizeString(v)
	case map[string]interface{}:
		for key, value := range v {
			v[key] = sanitizeValue(value)
		}
		return v
	case []interface{}:
		for i, value := range v {
			v[i] = sanitizeValue(value)
		}
		return v
	default:
		//return v as it is unsported type
		return v
	}
}

func sanitizeString(value string) string {
	return bluemonday.UGCPolicy().Sanitize(value)
}
