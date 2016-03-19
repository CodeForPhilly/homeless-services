package rap

import (
	"appengine"
	"appengine/datastore"
	//"appengine/memcache"
	//"fmt"
	//"github.com/VividCortex/siesta"
	"encoding/json"
	"errors"
	"log"
	"net/http"
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
						res[i].Location.Lat,
						res[i].Location.Lng,
					},
				},
				Properties: res[i],
			})
		}

		gc := featurecollection{
			"FeatureCollection",
			f,
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		//err := json.NewEncoder(w).Encode(res)
		tmp, err := json.Marshal(gc)
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
