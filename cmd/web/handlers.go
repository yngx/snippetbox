package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

/*
 The http.ResponseWriter parameter provides methods for assembling a HTTP response and sending it to the user.
 The *http.Request parameter is a pointer to a struct which holds information about the current request
*/
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	/*
		Check if the current request URL path exactly matches "/". If it doesn't, use
		the http.NotFound() function to send a 404 response to the client.
	*/
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// We initialize a slice containing the paths to the two template files. Note that the
	// home.page.tmpl file must be the *first* file in the slice.
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// We use the template.ParseFiles() function to read the template file into a
	// template set. We can pass the slice of file paths as a variadic parameter.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// We then use the Execute() method on the template set to write the template
	// content as the response body. The last parameter to Execute() represents any
	// dynamic data that we want to pass in, which for now we'll leave as nil.
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id parameter from the query string and try to
	// convert it to an integer using the strconv.Atoi() function.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
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

func (app *application) showSnippets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Display"))
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		/*
			Any changes to the header after w.WriteHeader will not
			have an effect on the headers!

			w.WriteHeader(405)
			w.Write([]byte("Method Not Allowed"))
		*/

		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Create some variables holding dummy data. We'll remove these later on
	// during the build.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)

	// /*
	// 	By not explicitly calling w.WriteHeader(), the first call to
	// 	w.Write() will send a 200 OK status code to the user.
	// 	In this case it is fine. But if we want to sen a non-200 code
	// 	we must call w.WriteHeader() before the next line.
	// */
	// w.Write([]byte("Created a snippet..."))
}
