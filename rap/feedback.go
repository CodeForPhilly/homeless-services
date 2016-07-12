package rap

import (
	"errors"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func feedback(w http.ResponseWriter, r *http.Request) *appError {
	c := appengine.NewContext(r)
	log.Infof(c, "method: ", r.Method)

	if r.Method != "POST" {
		return &appError{
			errors.New("Unsupported method call to feedback"),
			"Feedback most be POSTed",
			http.StatusMethodNotAllowed,
		}
	}

	r.ParseForm()

	//build the RECAPTCHA request
	rc := &recaptchaRequest{
		Secret:   recaptchaServerKey,
		Response: r.Form.Get("g-recaptcha-response"),
		RemoteIP: r.RemoteAddr,
	}

	log.Infof(c, "rc: ", rc)

	//send to google

	//read response

	//redisplay form with errors if it failed

	//ok, so they passed RECAPTCHA, proceed
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
		log.Debugf(c, err.Error())
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

type recaptchaRequest struct {
	Secret   string `json:"secret"`
	Response string `json:"response"`
	RemoteIP string `json:"remoteip"`
}

type recaptchaResponse struct {
	Success      bool
	Challenge_ts time.Time
	Hostname     string
	ErrorCodes   []string `json:"error-codes"`
}
