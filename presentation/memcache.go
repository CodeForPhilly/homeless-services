package rap
import (
    	"hash/fnv"
        	"google.golang.org/appengine/memcache"
)

func getResources(w http.ResponseWriter, r *http.Request) *appError {
	//omitted code to get top, filter et cetera from the Request

    c := appengine.NewContext(r)
    h := fnv.New64a()
	h.Write([]byte("rap_query" + top + filter + orderby + skip))
	cacheKey := h.Sum64()

	cv, err := memcache.Get(c, fmt.Sprint(cacheKey))
	if err == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(cv.Value)
		return nil
	}

    //omitting code that get resources from DataStore

    //cache this result with it's query as the key
	memcache.Set(c, &memcache.Item{
		Key:   fmt.Sprint(cacheKey),
		Value: result,
	})
    //return the json version
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(result)
	return nil
}

