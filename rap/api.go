package rap

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type featurecollection struct {
	GeoType  string     `json:"type"`
	Features []*feature `json:"features"`
}

type feature struct {
	Geotype    string    `json:"type"`
	Geometry   *geometry `json:"geometry"`
	Properties *resource `json:"properties"`
}

type geometry struct {
	Geotype     string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
	Id          string    `json:"id"`
}

/* we need to format our resources like so
from the geojson spec:
http://geojson.org/geojson-spec.html#geometrycollection

I think we'll want to implement OData querying on resources with an extra flag for whether we return the data as json or geojson.
We need the geojson for easy population of maps, but the regular json would be more useful for anything else (like webpages that update the resource).

But... that's all future stuff. For now the only thing I am doing is geojson and caching.

{ "type": "FeatureCollection",
    "features": [
      { "type": "Feature",
        "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
        "properties": {"prop0": "value0"}
        },
*/

/*
I'm thinking we store the data under 4 keys, one for each category
func get_data(){
    data = memcache.Get('key')
    if data != nil{
        return data
   }
    else{
        data = self.Query_for_data()
        memcache.Add('key', data, 60)
        return data
   }
}
*/

func resources(w http.ResponseWriter, r *http.Request) *appError {
	log.Println(r.Method)

	switch r.Method {
	case "GET":
		return getResources(w, r)
	case "POST":
		//auth, validate, add to datastore, invalidate cache
		return &appError{
			errors.New("Post not yet implemented"),
			"RESTful updates not yet implemented",
			http.StatusNotImplemented,
		}
	case "PUT":
		//auth, validate, update in datastore, invalidate cache
		return &appError{
			errors.New("Put not yet implemented"),
			"RESTful updates not yet implemented",
			http.StatusNotImplemented,
		}
	case "DELETE":
		//auth, validate, update in datastore, invalidate cache
		return &appError{
			errors.New("Delete not yet implemented"),
			"RESTful updates not yet implemented",
			http.StatusNotImplemented,
		}
	default:
		return &appError{
			errors.New("Attempted unsupport HTTP method " + r.Method),
			"Unsupported HTTP method",
			http.StatusMethodNotAllowed,
		}
	}
}

func getResources(w http.ResponseWriter, r *http.Request) *appError {
	c := appengine.NewContext(r)
	//get query from request
	u := r.URL.Query()
	top := strings.Trim(u.Get("top"), " ")
	//qselect :=u.Get("select") //dunno if we'll support this... definitely not for demo
	filter := strings.Trim(u.Get("filter"), " ")
	//count := strings.Trim(u.Get("count"), " ")
	//everything is text right now so this might not work as expected
	orderby := strings.Trim(u.Get("orderby"), " ")
	skip := strings.Trim(u.Get("skip"), " ")

	//This text search will hit the description field... or we would except that datastore doesn't support that kinda thing
	//you can get around that restriction by creating a field that is a unique list of all the strings and searching that. Of course we're not doing that right now.
	//search := strings.Trim(u.Get("search"), " ")

	//eventually allow for straight json as opposed to the geojson we normally serve - not supported yet :P
	//format:=u.Get("format")

	//build query against the db
	q := datastore.NewQuery("Resource").Filter("IsActive =", true)

	//loop through the filters - eventually, for now just handle one
	if len(filter) > 0 {
		//= q.Filter(
	}

	//handle orderby - works for a single order, Caps matter
	if len(orderby) > 0 { //need to actually parse
		//check to see if the string ends with desc
		if strings.LastIndex(orderby, "desc") > len(orderby)-len(" desc") {
			orderby = strings.Trim(orderby[:len(orderby)-len(" desc")], " ")
			q = q.Order("-" + orderby)
		}
		q = q.Order(orderby)
	}

	res := make([]*resource, 0)

	keys := make([]*datastore.Key, 0)

	//based on https://cloud.google.com/appengine/docs/go/datastore/queries#Go_Sort_orders
	//use the cursor to handle the top and skip
	//need to put in the skip and take fit into cursors... so tired right now
	if len(skip) > 0 || len(top) > 0 {
		item, err := memcache.Get(c, "rap_cursor")
		if err == nil {
			cursor, err := datastore.DecodeCursor(string(item.Value))
			if err == nil {
				q = q.Start(cursor)
			}
		}

		// Iterate over the results.
		t := q.Run(c)
		for {
			var tmp resource
			_, err := t.Next(&tmp)
			if err == datastore.Done {
				break
			}
			if err != nil {
				return &appError{
					err,
					"Fetching next Resource",
					http.StatusInternalServerError,
				}
			}

			res = append(res, &tmp)
		}

		// Get updated cursor and store it for next time.
		if cursor, err := t.Cursor(); err == nil {
			memcache.Set(c, &memcache.Item{
				Key:   "person_cursor",
				Value: []byte(cursor.String()),
			})
		}
	}

	//execute query against the db
	keys, err := q.GetAll(c, &res)

	if err != nil {
		return &appError{err, "Error querying database", http.StatusInternalServerError}
	}

	//make geojson from the db results
	f := make([]*feature, 0)

	//keys = keys[:1]
	//res = res[:1]

	for i, k := range keys {
		res[i].ID = k.IntID()

		f = append(f, &feature{
			Geotype: "Feature",
			Geometry: &geometry{
				Geotype: "Point",
				Coordinates: []float64{
					res[i].Location.Lng,
					res[i].Location.Lat,
					0,
				},
				Id: strconv.FormatInt(res[i].ID, 10),
			},
			Properties: res[i],
		})
	}

	gc := featurecollection{
		"FeatureCollection",
		f,
	}

	//return the json version
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//err := json.NewEncoder(w).Encode(res)
	tmp, err := json.Marshal(gc)
	if err != nil {
		return &appError{err, "Error creating response", http.StatusInternalServerError}
	}
	w.Write(tmp)
	return nil
}

func odataLogicalOperatorConverter(eo string) (string, error) {
	switch strings.ToLower(strings.Trim(eo, " ")) {
	case "eq":
		return "=", nil
		//case "ne": //not directly supported by datastore
	//return "!="
	case "gt":
		return ">", nil
	case "ge":
		return ">=", nil
	case "lt":
		return "<", nil
	case "le":
		return "<=", nil
	default:
		return "", errors.New("Unknown or unsupported operator")
	}
}
