//The RAP site is composed of three parts:
//-static pages - for the public
//-API - for RESTful CRUD ops on Datastore
//-logged in pages - as a web based way access the CRUD ops - not in the initial release
//-an import page for the bulk updates

//google app engine handles auth fairly well on its own

/* sometime performance, json, and geocode info
talks.golang.or/2015/json.slide#1
https://github.com/nf/geocode
https://github.com/nf/geocode
talks.golang.org/2013/highperf.slide#1
github.com/mjibson/appstats
*/

package rap

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
	"os"
	"time"
)

const basePath = "rap"

var recaptchaServerKey, recaptchaBrowserKey, mapsBrowserKey, geocodingServerKey string

func init() {
	//load the api keys
	recaptchaBrowserKey = os.Getenv("RECAPTCHA_SITE_KEY")
	recaptchaServerKey = os.Getenv("RECAPTCHA_SECRET_KEY")
	mapsBrowserKey = os.Getenv("MAPS_BROWSER_KEY")
	geocodingServerKey = os.Getenv("GECODING_SERVER_KEY")

	//basePath = "rapdemo"
	fs := http.FileServer(http.Dir(basePath + "/static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.Handle("/auth", appHandler(authdemo))

	//bulk import from csv
	http.Handle("/csvimport", appHandler(csvimport))

	//api
	http.Handle("/resources", appHandler(resources))

	//feedback page submit
	http.Handle("/feedback", appHandler(feedback))

	//handles the templated but otherwise mostly static html pages
	http.Handle("/", appHandler(serveTemplate))
}

//The resource type is what most of the application will focus on.
type resource struct {
	//ID int64 //db id

	//display fields
	Category         string             `datastore:"category"` //Medical, Food, et cetera
	OrganizationName string             `datastore:"organizationname"`
	Address          string             `datastore:"address"`
	ZipCode          string             `datastore:"zipcode"`
	Days             string             `datastore:"days"`
	TimeOpen         string             `datastore:"timeopen"`
	TimeClose        string             `datastore:"timeclose"`
	PeopleServed     string             `datastore:"peopleserved"`
	Description      string             `datastore:"description"`
	PhoneNumber      string             `datastore:"phonenumber"`
	Location         appengine.GeoPoint `datastore:"location"` //lng lat

	//audit fields
	LastUpdatedTime time.Time `datastore:"lastupdatedtime,noindex"`
	LastUpdatedBy   string    `datastore:"lastupdatedby,noindex"`
	IsActive        bool      `datastore:"isactive"`
}

//following the error pattern suggested in the Go Blog
//http://blog.golang.org/error-handling-and-go

type appError struct {
	Error   error
	Message string
	Code    int
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		c := appengine.NewContext(r)
		log.Errorf(c, "%v", e.Error)
		http.Error(w, e.Message, e.Code)
	}
}
