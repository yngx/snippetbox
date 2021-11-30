package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/yngx/snippetbox/pkg/models"
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

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id parameter from the query string and try to
	// convert it to an integer using the strconv.Atoi() function.
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	/*
		We pass our http.ResponseWriter to fmt.Fprintf because our http.ResponseWriter
		satisfies the io.Writer interface that Fprintf requires.
	*/
	//fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)

	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Use the new render helper.
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

func (app *application) showSnippets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	snippets, err := app.snippets.GetAll()
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	for i := range snippets {
		fmt.Fprintf(w, "%v", snippets[i])
	}
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new snippet..."))
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
