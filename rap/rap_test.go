//Integration tests that call the public site and API to simulate actual use.
package rap_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"google.golang.org/appengine/aetest"
)

func TestPages(t *testing.T) {
	_, done, err := aetest.NewContext()
	//ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	//test that the webpages pages come back
	basePageTester(t, "index", "map", "calendar")

	//this could be switched over to move of a table driven test if we test a known import file first to get a consistent set of data to test.
	apiTester(
		t,
		"",
		"filter=category%20eq%20Legal",
		"filter=category%20eq%20Medical&top=5&skip=2",
		"filter=category%20eq%20Legal&orderBy=OrganizationName&skip=1",
	)
}

func basePageTester(t *testing.T, pns ...string) {
	for _, pn := range pns {
		getTester(t, "http://localhost:8080/"+pn+".html")
	}
}

func apiTester(t *testing.T, queries ...string) {
	for _, q := range queries {
		getTester(t, "http://localhost:8080/resources?"+q)
	}
}

func getTester(t *testing.T, url string) {
	res, err := http.Get(url)
	if err != nil {
		t.Errorf("Error Getting response: %v", err)
	}
	if res.StatusCode != 200 && res.StatusCode != 302 {
		t.Errorf("Error Getting %s: %s", url, res.Status)
		message, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		t.Errorf("Body of response: %s", message)
	}
	res.Body.Close()
}
