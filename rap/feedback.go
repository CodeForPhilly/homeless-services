package rap

import (
	"appengine"
	"appengine/datastore"
	"errors"
	"log"
	"net/http"
	//"appengine/datastore"
	"time"
)

func feedback(w http.ResponseWriter, r *http.Request) *appError {
	c := appengine.NewContext(r)
	c.Debugf("method: ", r.Method)

	if r.Method != "POST" {
		return &appError{
			errors.New("Unsupported method call to feedback"),
			"Feedback most be POSTed",
			http.StatusMethodNotAllowed,
		}
	}

	r.ParseForm()

	f := &feedbackform{
		Topic:           r.Form.Get("topic"),
		Challenges:      r.Form.Get("challenges"),
		Recommendations: r.Form.Get("recommendations"),
		Email:           r.Form.Get("email"),
		SubmittedTime:   time.Now().UTC(),
		SubmitterIP:     r.RemoteAddr,
		IsActive:        true,
	}

	key := datastore.NewIncompleteKey(c, "Feedback", nil)

	_, err := datastore.Put(c, key, f)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return &appError{err, "Error submitting feedback", http.StatusInternalServerError}
	}

	http.Redirect(w, r, "/index.html", http.StatusFound)
	return nil
}

type feedbackform struct {
	Topic, Challenges, Recommendations, SubmitterIP, Email string
	SubmittedTime                                          time.Time `datastore:",noindex"`
	IsActive                                               bool
}
