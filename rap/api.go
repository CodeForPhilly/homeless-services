package rap

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"math"
	"net/http"
	"strconv"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

//ErrUnknownOperator represents an error for unknown logical operators.
var ErrUnknownOperator = errors.New("Unknown or unsupported operator")

type featurecollection struct {
	GeoType  string     `json:"type"`
	Features []*feature `json:"features"`
}

type feature struct {
	Geotype    string    `json:"type"`
	Geometry   *geometry `json:"geometry"`
	Properties *Resource `json:"properties"`
}

type geometry struct {
	Geotype     string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
	ID          string    `json:"id"`
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
	//log.Infof(r.Method)

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
	//qselect :=u.Get("select") //dunno if we'll ever support this
	filter := strings.Trim(u.Get("filter"), " ")
	//count := strings.Trim(u.Get("count"), " ") //not yet implemented, maybe once there is a web interface for editting data
	//everything is text right now so this might not work as expected
	orderby := strings.ToLower(strings.Trim(u.Get("orderby"), " "))
	skip := strings.Trim(u.Get("skip"), " ")

	//This text search will hit the description field... or we would except that datastore doesn't support that kinda thing
	//you can get around that restriction by creating a field that is a unique list of all the strings and searching that. Of course we're not doing that right now.
	//https://cloud.google.com/appengine/articles/indexselection
	//search := strings.Trim(u.Get("search"), " ")

	//eventually allow for straight json as opposed to the geojson we normally serve - not supported yet :P
	//format:=u.Get("format")

	//here we should check the cache to see if it holds an identical query so we can return that
	h := fnv.New64a()
	h.Write([]byte("rap_query" + top + filter + orderby + skip))
	cacheKey := h.Sum64()

	cv, err := memcache.Get(c, fmt.Sprint(cacheKey))
	if err == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(cv.Value)
		return nil
	}

	//build query against the db
	q := datastore.NewQuery("Resource").Filter("isactive =", true)

	//loop through the filters
	log.Debugf(c, "filter %s", filter)
	if len(filter) > 0 {
		if cn, eo, fv, err := filterParser(filter); err == nil {
			log.Debugf(c, "cn %s, eo %s, fv %s", cn, eo, fv)
			q = q.Filter(cn+" "+eo, fv)
		} else {
			log.Debugf(c, "filter parsing error %s", err)
		}
	}

	//handle orderby - works for a single order
	if len(orderby) > 0 { //need to actually parse
		//check to see if the string ends with desc
		if strings.LastIndex(orderby, "desc") > len(orderby)-len(" desc") {
			orderby = strings.Trim(orderby[:len(orderby)-len(" desc")], " ")
			orderby = "-" + orderby
		}
		q = q.Order(orderby)
	} else {
		q = q.Order("organizationname")
	}

	var res []*Resource

	//based on https://cloud.google.com/appengine/docs/go/datastore/queries#Go_Sort_orders
	//use the cursor to handle the top and skip
	//In this case we don't care about the error because the int will then be zero. We'll proceed with the query with what worked.
	//Alternatively we could check the errors and return a bad request accordingly
	s, _ := strconv.Atoi(skip)
	t, _ := strconv.Atoi(top)
	if s > 0 || t > 0 {
		if t == 0 {
			t = math.MaxInt32
		}

		// Iterate over the results.
		iter := q.Run(c)
		var tmp Resource
		var err error
		//There's a bug here that causes this to return the same resource top times
		for i := 0; s+t > i; i++ {
			if s > i {
				_, err = iter.Next(nil) //don't pull anything till we're done skipping
			} else {
				_, err = iter.Next(&tmp)
			}

			if err == datastore.Done {
				log.Debugf(c, "i %d, s %d, t %d", i, s, t)
				break
			}
			if err != nil {
				return &appError{
					err,
					"Error in skip/top querying",
					http.StatusInternalServerError,
				}
			}

			//don't keep anything till we're done skipping
			if i >= s {
				res = append(res, &tmp)
			}
		}
	} else {
		//execute query against the db
		_, err := q.GetAll(c, &res)

		if err != nil {
			return &appError{err, "Error querying database", http.StatusInternalServerError}
		}
	}

	//make geojson from the db results
	var f []*feature

	log.Debugf(c, "res contains %v items", len(res))
	for _, v := range res {
		f = append(f, &feature{
			Geotype: "Feature",
			Geometry: &geometry{
				Geotype: "Point",
				Coordinates: []float64{
					v.Location.Lng,
					v.Location.Lat,
					0,
				},
				//Id: strconv.FormatInt(v.ID, 10), //not needed
			},
			Properties: v,
		})
	}

	gc := featurecollection{
		"FeatureCollection",
		f,
	}

	tmp, err := json.Marshal(gc)
	if err != nil {
		return &appError{err, "Error creating response", http.StatusInternalServerError}
	}

	//cache this result with it's query as the key
	memcache.Set(c, &memcache.Item{
		Key:   fmt.Sprint(cacheKey),
		Value: tmp,
	})

	//return the json version
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(tmp)
	return nil
}

//I'd like to refactor this function and logicalOperatorConverter so that the logical operators are a single map used in both
func filterParser(f string) (cn, eo, fv string, err error) {
	odataOperators := []string{"eq", "gt", "ge", "lt", "le"}

	//find the operator in this filter and use it to split the filter
	for _, lo := range odataOperators {
		if strings.Index(f, lo) == -1 {
			continue
		}

		parts := strings.Split(f, lo)
		if len(parts) != 2 {
			err = errors.New("Incorrect number of filter parameters")
			break
		}

		if eo, err = logicalOperatorConverter(lo); err == nil {
			cn = strings.Trim(parts[0], " ")
			fv = strings.Trim(parts[1], " ")
			return
		}

		break
	}

	if err == nil {
		err = errors.New("No supported logical operator was found")
	}

	return "", "", "", err
}

func logicalOperatorConverter(eo string) (string, error) {
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
		return "", ErrUnknownOperator
	}
}
