//The RAP site is composed of three parts:
//-static pages - for the public
//-API - for RESTful CRUD ops on the data
//-logged in pages - as a web based way access the CRUD ops
//-an import page for the bulk updates

//The static pages are based on this tutorial
//http://www.alexedwards.net/blog/serving-static-sites-with-go

//I was thinking of trying siesta (https://github.com/VividCortex/siesta) for the api

//google app engine handles auth fairly well on its own

package rapdemo

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
)

const basePath = "rapdemo"

func init() {
	//basePath = "rapdemo"
	fs := http.FileServer(http.Dir(basePath + "/static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.HandleFunc("/auth", authdemo)

	//datastore testing
	http.HandleFunc("/dsdemo", dsdemo)

	//service := siesta.NewService("/api/")
	http.HandleFunc("/", serveTemplate)
}

//it would be good for this function to pass a token to the page in case the page has a form (a lot of them will)
func serveTemplate(w http.ResponseWriter, r *http.Request) {
	log.Println("serving a non-static request")
	lp := path.Join(basePath+"/templates", "layout.html")
	fp := path.Join(basePath+"/templates", r.URL.Path)

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("couldn't find the template")
			log.Println(err.Error())
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		log.Println("found a directory insteead of a file")
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		// Log the detailed error
		log.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
