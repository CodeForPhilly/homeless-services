//This page implements some authentication based on this tutorial
//https://cloud.google.com/appengine/docs/go/gettingstarted/usingusers
//The intention is that eventually, specific users could log in to update the RAP information
package rapdemo

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

func authdemo(w http.ResponseWriter, r *http.Request) *appError {
	c := appengine.NewContext(r)
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
	fmt.Fprintf(w, "Hello, %v!", u)
	return nil
}
