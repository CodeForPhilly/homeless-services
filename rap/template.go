//The static pages are based on this tutorial
//http://www.alexedwards.net/blog/serving-static-sites-with-go

package rap

import (
	"appengine"
	"errors"
	"html/template"
	"net/http"
	"os"
	"path"
)

//it would be good for this function to pass a token to the page in case the page has a form (a lot of them will)
//it would also be good to default to a page if none is given in the r.URL.Path
func serveTemplate(w http.ResponseWriter, r *http.Request) *appError {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/index.html", 301)
	}

	c := appengine.NewContext(r)
	c.Debugf("serving a non-static request")

	lp := path.Join(basePath+"/templates", "layout.html")
	fp := path.Join(basePath+"/templates", r.URL.Path)

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			return &appError{
				err,
				"Couldn't find the specified template",
				http.StatusNotFound,
			}
		}
		return &appError{err, "Unknown error finding template", http.StatusInternalServerError}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		return &appError{
			errors.New("Attempted to display directory " + r.URL.Path),
			"Can't display record",
			http.StatusNotFound,
		}
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		return &appError{err, "Error parsing template", http.StatusInternalServerError}
	}

	if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
		return &appError{err, "Error executing template", http.StatusInternalServerError}
	}
	return nil
}
