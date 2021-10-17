package main

import (
	"log"
	"net/http"
)

func main() {

	/*
		http.handleFunc() functions allow us to register routes without declaring
		a servemux. These functions, behind the scene, register their routes with
		the DefaultServeMux which is initialized by default as a global variable.
		in net/http global. The downside of using the DefaultServeMux is that since
		it is a global variable, any package can access it and register a route.
		This means any third party package that our app imports can register routes
		and expose them to our users.
	*/

	// For the sake of this exercise, let's keep using mux.
	mux := http.NewServeMux()

	// subtree path, patterns are matched
	mux.HandleFunc("/", home)

	// fixed paths
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippets", showSnippets)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Println("Starting server on :4000")
	log.Fatal(http.ListenAndServe(":4000", mux))

}
