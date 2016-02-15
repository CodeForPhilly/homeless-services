package rapdemo_test

import (
	"google.golang.org/appengine/aetest"
	"net/http"

	"io/ioutil"
	"log"
	"testing"
)

func TestRapDemo(t *testing.T) {
	_, done, err := aetest.NewContext()
	//ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	//test that the base pages come back
	basePageTester(t, "index", "map", "calendar")
}

func basePageTester(t *testing.T, pns ...string) {
	for _, pn := range pns {
		res, err := http.Get("http://localhost:8080/" + pn + ".html")
		if err != nil {
			t.Errorf("Error Getting Index: %v", err)
		}
		if res.StatusCode != 200 && res.StatusCode != 302 {
			t.Errorf("Error Getting %s: %s", pn, res.Status)
			message, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			t.Errorf("Body of response: %s", message)
		}
		res.Body.Close()
	}
}
