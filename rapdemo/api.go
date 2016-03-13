package rapdemo

import (
	"appengine"
	"appengine/datastore"
	//"fmt"
	//"github.com/VividCortex/siesta"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

/* we need to format our resources like so
from the geojson spec:
http://geojson.org/geojson-spec.html#geometrycollection

I think we'll want to implement OData querying on resources with an extra flag for whether we return the data as json or geojson.
We need the geojson for easy population of maps, but the regular json would be more useful for anything else (like webpages that update the resource).

But... that's all future stuff. For now the only thing I am doing is geojson and caching.

{ "type": "GeometryCollection",
    "geometries": [
      { "type": "Point",
        "coordinates": [100.0, 0.0]
        },
      { "type": "LineString",
        "coordinates": [ [101.0, 0.0], [102.0, 1.0] ]
        }
    ]
}
*/

func resources(w http.ResponseWriter, r *http.Request) *appError {
	log.Println(r.Method)
	c := appengine.NewContext(r)
	switch r.Method {
	case "GET":
		res := make([]*resource, 0)
		q := datastore.NewQuery("Resource").
			Filter("IsActive =", true)

		keys, err := q.GetAll(c, &res)

		if err != nil {
			return &appError{err, "Error querying database", http.StatusInternalServerError}
		}

		for i, k := range keys {
			res[i].ID = k.IntID()
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		//err := json.NewEncoder(w).Encode(res)
		tmp, err := json.Marshal(res)
		if err != nil {
			return &appError{err, "Error creating response", http.StatusInternalServerError}
		}
		w.Write(tmp)
		return nil
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
	default:
		return &appError{
			errors.New("Attempted unsupport HTTP method " + r.Method),
			"Unsupported HTTP method",
			http.StatusMethodNotAllowed,
		}
	}
}
