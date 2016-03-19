//File upload based on https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/04.5.html
//datastore usage based on https://cloud.google.com/appengine/docs/go/gettingstarted/usingdatastore

package rap

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"appengine/user"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

//so... now we need to get the csv file into memory
//we'll do that by a form submit from a static
//then parse the form data from the request into resources
//then into datastore

//I need to implement a token on this form
func csvimport(w http.ResponseWriter, r *http.Request) *appError {
	c := appengine.NewContext(r)
	c.Debugf("method: ", r.Method)

	if r.Method != "POST" {
		return &appError{
			errors.New("Unsupported method call to import"),
			"Imports most be POSTed",
			http.StatusMethodNotAllowed,
		}
	}

	//this block for check the user's credentials should eventually be broken out into a filter
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			return &appError{err, "Could not determine LoginURL", http.StatusInternalServerError}
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return nil
	}

	//some crappy security so that only a certain person can upload things
	//we should probably have a users entity in datastore that we manage manually for this kinda thing
	if u.Email != "test@example.com" {
		return &appError{
			errors.New("Illegal import attempted by " + u.Email),
			"Your user is not allowed to import",
			http.StatusForbidden,
		}
	}

	//r.ParseMultipartForm(1 << 10)

	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		return &appError{err, "Error uploading file", http.StatusInternalServerError}
	}
	defer file.Close()

	log.Println(handler.Filename)

	cr := csv.NewReader(file)
	var res []resource
	var keys []*datastore.Key

	//at the moment we always insert a new item, this should be an insert or update based on OrganizationName
	//if we get a large enough data set we'll need to implement two loops so that we only batch a certain number of records at a time
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &appError{err, "Error reading file", http.StatusInternalServerError}
		}

		//we may want IDs in there eventually
		//_, err = strconv.ParseInt(rec[0], 2, 64)
		res = append(res, resource{
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
			Location: appengine.GeoPoint{
				Lat: 39.9522 - rand.Float64(),
				Lng: -75.1635 - rand.Float64(),
			},
		})

		keys = append(keys, datastore.NewIncompleteKey(c, "Resource", nil))
	}

	_, err = datastore.PutMulti(c, keys, res)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return &appError{err, "Error updating database", http.StatusInternalServerError}
	}

	// clear the cache
	memcache.Flush(c)

	http.Redirect(w, r, "/index.html", http.StatusFound)
	return nil
}
