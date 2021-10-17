package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

/*
 The http.ResponseWriter parameter provides methods for assembling a HTTP response and sending it to the user.
 The *http.Request parameter is a pointer to a struct which holds information about the current request
*/
func home(w http.ResponseWriter, r *http.Request) {
	/*
		Check if the current request URL path exactly matches "/". If it doesn't, use
		the http.NotFound() function to send a 404 response to the client.
	*/
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Use the template.ParseFiles() function to read the template file into a
	// template set. If there's an error, we log the detailed error message and use
	// the http.Error() function to send a generic 500 Internal Server Error
	// response to the user.
	ts, err := template.ParseFiles("./ui/html/home.page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// We then use the Execute() method on the template set to write the template
	// content as the response body. The last parameter to Execute() represents any
	// dynamic data that we want to pass in, which for now we'll leave as nil.
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id parameter from the query string and try to
	// convert it to an integer using the strconv.Atoi() function.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	/*
		We pass our http.ResponseWriter to fmt.Fprintf because our http.ResponseWriter
		satisfies the io.Writer interface that Fprintf requires.
	*/
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)

	/*
		Go will automatically set the following headers for us:
		Date, Content-Length, and Content-Type.

		The Content-Type header is set by content sniffing the response
		body with http.DetectContentType() function.

		A caveat here though is that Go can't distinguish JSON from plain
		text. So by default, JSON responses will be sent with
		`Content-Type: text/plain; charset=utf-8` unless we explicitlky
		set the header.
	*/
	//w.Write([]byte("Display"))
}

func showSnippets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	w.Write([]byte("Display"))
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		/*
			Any changes to the header after w.WriteHeader will not
			have an effect on the headers!

			w.WriteHeader(405)
			w.Write([]byte("Method Not Allowed"))
		*/

		http.Error(w, "Method Not Allowed", 405)
		return
	}

	/*
		By not explicitly calling w.WriteHeader(), the first call to
		w.Write() will send a 200 OK status code to the user.
		In this case it is fine. But if we want to sen a non-200 code
		we must call w.WriteHeader() before the next line.
	*/
	w.Write([]byte("Created a snippet..."))
}
