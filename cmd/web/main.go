package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// This parses the command-line flag.
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize a new instance of application containing the dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

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
	mux.HandleFunc("/", app.home)

	// fixed paths
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippets", app.showSnippets)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	/*
		Go's file server is pretty cool. It has a few really nice features:
		* sanitizes all request paths
		* range requests are fully supported. You can specify the range of bytes you want.
		* Last-Modified and If-Modified-Since headers are transparently supported.
		* Content-Type is automatically set from the file extension
	*/

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	infoLog.Printf("Starting server on %s", *addr)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	err := srv.ListenAndServe()
	// log.Fatal() function will also call os.Exit(1) after writing the message
	errorLog.Fatal(err)

}
