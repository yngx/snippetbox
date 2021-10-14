package main

import (
	"log"
	"net/http"
)

// The http.ResponseWriter parameter provides methods for assembling a HTTP response and sending it to the user.
// The *http.Request parameter is a pointer to a struct which holds information about the current request
func home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If it doesn't, use
	// the http.NotFound() function to send a 404 response to the client.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Wazzup"))
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display"))
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create"))
}

func main() {

	// http.handleFunc() functions allow us to register routes without declaring
	// a servemux. These functions, behind the scene, register their routes with
	// the DefaultServeMux which is initialized by default as a global variable.
	// in net/http global. The downside of using the DefaultServeMux is that since
	// it is a global variable, any package can access it and register a route.
	// This means any third party package that our app imports can register routes
	// and expose them to our users.

	// For the sake of this exercise, let's keep using mux.
	mux := http.NewServeMux()

	// subtree path, patterns are matched
	mux.HandleFunc("/", home)

	// fixed paths
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Println("Starting server on :4000")
	log.Fatal(http.ListenAndServe(":4000", mux))

}
