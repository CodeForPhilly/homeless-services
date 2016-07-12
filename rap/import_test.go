package rap

import (
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
)

func TestImport(t *testing.T) {
	c, done, err := aetest.NewContext()
	//ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	testTimeParse(t, c)
	//testDaysParse(t, c)
}

var (
	simpleTime     = []string{"Mon Tue Wed 8:00AM - 6:00PM", "", ""}
	separatedTime  = []string{"Mon, Tue, & Wed 8:00AM - 6:00PM", "", ""}
	spanOfTime     = []string{"Thu - Sat 8:00AM - 6:00PM", "", ""}
	differentTimes = []string{"Mon - Wed 8:00AM - 6:00PM Sat 1:00PM - 3:00PM", "", ""}
	complexTimes   = []string{"Fri - Sun 8:00AM - 6:00PM & Tue through Wed 9:00AM to 11:00am", "", ""}
	noSpan         = []string{"Mon, Tue, & Wed", "6:00PM", "10:00PM"}
	noSpan24       = []string{"Mon, Tue, & Wed", "24 hours", ""}
	badInput       = []string{"Fri - 8:00AM - 30:00PM & Tue through Wed 9:00PM to 11:00am", "WED", "9:00PM"}
)

//This should be checking the specific times and days instead of just the length
//This might be a good place to try some table driven tests
func testTimeParse(t *testing.T, c context.Context) {
	r := GetTimes(simpleTime, c)
	if len(r) != 3 {
		t.Errorf("Failed simple time parse: %s", r)
	}

	r = GetTimes(separatedTime, c)
	if len(r) != 3 {
		t.Errorf("Failed separated time parse: %s", r)
	}

	r = GetTimes(spanOfTime, c)
	if len(r) != 3 {
		t.Errorf("Failed span of time parse: %s", r)
	}

	r = GetTimes(differentTimes, c)
	if len(r) != 4 {
		t.Errorf("Failed different times parse: %s", r)
	}

	r = GetTimes(complexTimes, c)
	if len(r) != 5 {
		t.Errorf("Failed complex times parse: %s", r)
	}

	r = GetTimes(noSpan, c)
	if len(r) != 3 {
		t.Errorf("Failed no span parse: %s", r)
	}

	r = GetTimes(noSpan24, c)
	if len(r) != 3 {
		t.Errorf("Failed no span 24 parse: %s", r)
	}

	r = GetTimes(complexTimes, c)
	if len(r) != 0 {
		t.Errorf("Failed bad input parse: %s", r)
	}
}

/*
func hasDays(current []time.Weekday, potential ...time.Weekday) bool {
	var found int

	for _, p := range potential {
		if HasDay(p, current) {
			found++
		}
	}

	//if the potential slice has been emptied, then true
	return len(potential) == found
}
*/
