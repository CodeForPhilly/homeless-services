//File upload based on https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/04.5.html
//datastore usage based on https://cloud.google.com/appengine/docs/go/gettingstarted/usingdatastore

package rap

import (
	"encoding/csv"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/user"
)

//so... now we need to get the csv file into memory
//we'll do that by a form submit from a static
//then parse the form data from the request into resources
//then into datastore

//I need to implement a token on this form
func csvimport(w http.ResponseWriter, r *http.Request) *appError {
	c := appengine.NewContext(r)
	log.Infof(c, "method: ", r.Method)

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

	log.Infof(c, "New import file: %s ", handler.Filename)

	cr := csv.NewReader(file)
	var res []*Resource
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

		//if the first row has column headers then skip to the next one
		if strings.ToLower(strings.Trim(rec[1], " ")) == "category" {
			continue
		}

		//Search for this Resource by OrganizationName
		q := datastore.NewQuery("Resource").Filter("organizationname =", rec[2]).KeysOnly().Limit(2)
		tmpKey := datastore.NewIncompleteKey(c, "Resource", nil)
		if tmpKeys, err := q.GetAll(c, nil); len(tmpKeys) == 1 && err == nil {
			tmpKey = tmpKeys[0]
		}

		//we may want IDs in there eventually
		//_, err = strconv.ParseInt(rec[0], 2, 64)
		tmp := &Resource{
			Categories:       getSliceFromString(rec[1]),
			OrganizationName: rec[2],
			Address:          rec[3],
			ZipCode:          rec[4],
			Days:             GetDays(rec[5], c),
			TimeOpenClose:    GetTimes(rec[6], c),
			PeopleServed:     getSliceFromString(rec[8]),
			Description:      rec[9],
			PhoneNumber:      rec[10],
			LastUpdatedBy:    u.Email,
			LastUpdatedTime:  time.Now().UTC(),
			IsActive:         true,
			Location:         appengine.GeoPoint{},
		}

		//log.Infof(c, "len slice check: %x, len rec LatLng check: %x, check for comma: %x", len(rec) > 11, len(rec[11]) > 0, strings.Index(rec[11], ",") != -1)

		if len(rec) > 11 && len(rec[11]) > 0 && strings.Index(rec[11], ",") != -1 {
			tmp.Location.Lng, _ = strconv.ParseFloat(strings.Split(rec[11], ",")[0], 64)
			tmp.Location.Lat, _ = strconv.ParseFloat(strings.Split(rec[11], ",")[1], 64)
			//log.Println(tmp.Location)
		}

		res = append(res, tmp)

		keys = append(keys, tmpKey)
	}

	_, err = datastore.PutMulti(c, keys, res)
	if err != nil {
		log.Debugf(c, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return &appError{err, "Error updating database", http.StatusInternalServerError}
	}

	// clear the cache
	memcache.Flush(c)

	http.Redirect(w, r, "/index.html", http.StatusFound)
	return nil
}

func getSliceFromString(c string) []string {
	return strings.Split(c, ",")
}

func HasDay(d time.Weekday, days []time.Weekday) bool {
	for _, t := range days {
		if d == t {
			return true
		}
	}
	return false
}

func GetDays(d string, c context.Context) []time.Weekday {
	//c := appengine.NewContext(r)

	var span bool
	var days []time.Weekday

	//log.Infof(c, "d: "+d)
	for _, v := range strings.Split(d, " ") {
		//log.Infof(c, "v: "+v)

		//ignore these separators
		if v == "&" || v == "," || len(v) == 0 {
			continue
		}

		//this signals the start of a span so we'll need to loop
		if v == "-" || v == "through" || v == "to" {
			span = true
			continue
		}

		//if there is a trailing ":" or "," drop it
		lastChar := v[len(v)-1:]
		if lastChar == "," {
			v = v[:len(v)-1]
		}

		if d, f := dayTranslations[v]; f == true {
			//log.Infof(c, "do: %s", d)
			if span {
				//add a span of days, start the span with last day in the slice
				startDay := days[len(days)-1]
				//log.Infof(c, "d: %s", d)
				//log.Infof(c, "startDay: %s", startDay)

				//there is a bug here where we could add days we already have
				wds := getWeekdaySpan(startDay, d)

				//log.Infof(c, "wds: %s", wds)

				days = append(days, wds...)

				span = false
			} else if !HasDay(d, days) {
				days = append(days, d)
			}
		}
	}

	//log.Infof(c, "days: %s", days)
	return days
}

func getWeekdaySpan(s, e time.Weekday, days ...time.Weekday) []time.Weekday {
	//I need to rework this...
	for _, w := range Weekdays {
		if w != s {
			continue
		}
		s++

		if w == e {
			break
		}

		if w+1 == time.Sunday {
			getWeekdaySpan(time.Sunday, e, days...)
		} else {
			days = append(days, w+1)
		}
	}

	return days
}

//getTimes takes a string of days and times and returns a slice of dailyAvailability. It does not return errors. If parsing fails, an error is logged and an empty slice is returned.
/*
getTimes intends to be flexible with formats. "," ":" and "&" are ignored.
"Mon 8:30AM - 9:00PM" -> Day: Monday, Open: 8:30AM, Close: 9PM
"Sat - Sun 2:00pm - 5:00pm" -> [ { Day: Saturday, Open: 2PM, Close: 5PM }, { Day: Monday, Open: 2PM, Close: 9PM} ]
"Tues & Wed: 9-3pm" -> [ { Day: Tuesday, Open: 9AM, Close: 3PM }, { Day: Wednesday, Open: 9AM, Close: 3PM} ]

Technically we could use a lexer/parser here but that's an overkill pipedream :(
*/
func GetTimes(s string, c context.Context) []dailyAvailability {
	//this works by accumulating Weekdays until a time span is encoutered
	//then that time span is applied to those days
	var times []dailyAvailability
	var span bool
	var days []time.Weekday
	var open, close time.Time

	log.Infof(c, "s: "+s)
	for _, dt := range strings.Split(strings.Replace(s, ",", " ", -1), " ") {

		log.Infof(c, "dt: "+dt)

		//this signals the start of a span so we'll need to loop
		if dt == "-" || dt == "through" || dt == "to" {
			span = true
			continue
		}

		//ignore these separators
		if dt == ":" || dt == "&" || len(dt) <= 1 {
			continue
		}

		//if there is a trailing ":" or "," drop it
		lastChar := dt[len(dt)-1:]
		if lastChar == ":" {
			dt = dt[:len(dt)-1]
		}

		if d, f := dayTranslations[dt]; f == true {
			log.Infof(c, "d: %s", d)
			if span {
				//add a span of days, start the span with last day in the slice
				startDay := days[len(days)-1]
				//log.Infof(c, "d: %s", d)
				//log.Infof(c, "startDay: %s", startDay)

				//there is a bug here where we could add days we already have
				wds := getWeekdaySpan(startDay, d)

				//log.Infof(c, "wds: %s", wds)

				days = append(days, wds...)

				span = false
			} else if !HasDay(d, days) {
				days = append(days, d)
			}
		}
		/*
			if d, err := time.Parse("Mon", v); err == nil {
				if span {
					//add a span of days, start the span with last day in the slice
					startDay := days[len(days)-1]
					for _, w := range Weekdays {
						if w != startDay {
							continue
						}
						days = append(days, d.Weekday())
						if w == d.Weekday() {
							break
						}
					}
					continue
				}
				days = append(days, d.Weekday())
			}
		*/

		//I would like to handle multiple formats here... for now just kitchen
		if t, err := time.Parse(time.Kitchen, dt); err == nil {
			log.Infof(c, "t: %s", t)
			log.Infof(c, "span: %s", span)
			if !span {
				open = t
				continue
			}
			close = t

			//if we didn't get a valid time span then move on
			if close.After(open) {
				//add the resulting days to the output slice
				for _, v := range days {
					times = append(
						times,
						dailyAvailability{
							Day:   v,
							Open:  open,
							Close: close,
						},
					)
				}
			}

			//reset for the next set of days and times
			open = time.Time{}
			close = time.Time{}
			days = days[:0]
			span = false
		}
	}

	log.Infof(c, "times: %s", times)

	return times
}

var (
	//Weekdays is a slice of time.Weekday used for import.
	Weekdays = []time.Weekday{
		time.Sunday,
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday,
		time.Saturday,
	}

	//dayTranslations maps strings to time.Weekday for parsing.
	dayTranslations = map[string]time.Weekday{
		"Sunday":    time.Sunday,
		"Monday":    time.Monday,
		"Tuesday":   time.Tuesday,
		"Wednesday": time.Wednesday,
		"Thursday":  time.Thursday,
		"Friday":    time.Friday,
		"Saturday":  time.Saturday,
		"Sun":       time.Sunday,
		"Mon":       time.Monday,
		"Tue":       time.Tuesday,
		"Wed":       time.Wednesday,
		"Thu":       time.Thursday,
		"Fri":       time.Friday,
		"Sat":       time.Saturday,
	}
)
