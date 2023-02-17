package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// Defining structure to 3rd party service request
type ProxyRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

// Defining structure to server response
type ProxyResponse struct {
	ID      string            `json:"id"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Length  int64             `json:"length"`
}

// Map to store requests from client
var requests map[string]ProxyRequest

// Map to stove responses from 3rd party service
var responses map[string]ProxyResponse

// running application
func main() {
	requests = make(map[string]ProxyRequest)
	responses = make(map[string]ProxyResponse)

	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(":8080", nil)
}

// Handler function to handle requests from client
func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// Just to test that everything works
	case "DELETE":
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	// Clients all other HTTP Methods handling
	default:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var proxyRequest ProxyRequest
		err = json.Unmarshal(body, &proxyRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id := uuid.NewV4().String()
		requests[id] = proxyRequest

		client := &http.Client{}
		req, err := http.NewRequest(proxyRequest.Method, proxyRequest.URL, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for key, value := range proxyRequest.Headers {
			req.Header.Set(key, value)
		}

		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		proxyResponse := ProxyResponse{
			ID:      id,
			Status:  resp.StatusCode,
			Headers: make(map[string]string),
			Length:  int64(len(body)),
		}

		for key, value := range proxyRequest.Headers {
			if _, ok := proxyResponse.Headers[key]; !ok {
				proxyResponse.Headers[key] = value
			}
		}

		responses[id] = proxyResponse

		jsonResponse, err := json.Marshal(proxyResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", jsonResponse)
	}
}
