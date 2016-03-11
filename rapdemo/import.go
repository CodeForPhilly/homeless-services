//File upload based on https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/04.5.html
//datastore usage based on https://cloud.google.com/appengine/docs/go/gettingstarted/usingdatastore

package rapdemo

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"time"
)

type resource struct {
	//ID int

	//display fields
	Category         string `json:"Category"`
	OrganizationName string `json:"Organization Name"`
	Address          string `json:"Address"`
	ZipCode          string `json:"Zip Code"`
	Days             string `json:"Days"`
	TimeOpen         string `json:"Time: Open"`
	TimeClose        string `json:"Time: Close"`
	PeopleServed     string `json:"People Served"`
	Description      string `json:"Description"`
	PhoneNumber      string `json:"Phone Number"`

	//audit fields
	LastUpdatedTime time.Time `datastore:",noindex"`
	LastUpdatedBy   string    `datastore:",noindex"`
	IsActive        bool
}

//so... now we need to get the csv file into memory
//we'll do that by a form submit from a static
//then parse the form data from the request into resources
//then into datastore

//I need to implement a token on this form
func dsdemo(w http.ResponseWriter, r *http.Request) {
	log.Println("method:", r.Method)

	if r.Method != "POST" {
		http.Redirect(w, r, "/index.html?wrongmethod=true", http.StatusFound)
		return
	}

	c := appengine.NewContext(r)

	//this block for check the user's credentials should eventually be broken out into a filter
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}

	//some crappy security so that only a certain person can upload things
	//we should probably have a users entity in datastore that we manage manually for this kinda thing
	if u.Email != "test@example.com" {
		http.Redirect(w, r, "/index.html?wronguser=true", http.StatusFound)
		return
	}

	//r.ParseMultipartForm(1 << 10)

	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		log.Println("error: ", err)
		return
	}
	defer file.Close()

	log.Println(handler.Filename)

	cr := csv.NewReader(file)

	//at the moment we always insert a new item, this should be an insert or update based on OrganizationName
	//also need to switch to batch operations with GetMulti and PutMulti
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		//we may want IDs in there eventually
		//_, err = strconv.ParseInt(rec[0], 2, 64)
		res := resource{
			Category:         rec[1],
			OrganizationName: rec[2],
			Address:          rec[3],
			ZipCode:          rec[4],
			Days:             rec[5],
			TimeOpen:         rec[6],
			TimeClose:        rec[7],
			PeopleServed:     rec[8],
			Description:      rec[9],
			PhoneNumber:      rec[10],
			LastUpdatedBy:    u.Email,
			LastUpdatedTime:  time.Now().UTC(),
			IsActive:         true,
		}

		dk := datastore.NewKey(c, "RAP", "default_rap", 0, nil)

		key := datastore.NewIncompleteKey(c, "Resource", dk)

		_, err = datastore.Put(c, key, &res)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/index.html", http.StatusFound)
}
