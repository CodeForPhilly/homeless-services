package rap

import (
	"google.golang.org/appengine/aetest"

	"testing"
	"time"
)

/*
func TestImport(t *testing.T) {
	_, done, err := aetest.NewContext()
	//ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	//test that the base pages come back
	basePageTester(t, "index", "map", "calendar")
}
*/

func TestTimeParse(t *testing.T) {
	_, done, err := aetest.NewContext()
	//ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	t.Log("Time Parse Test not yet implemented")
}

var (
	simpleDays    = "Mon Tue Wed"
	separatedDays = "Mon, Tue, & Wed"
	spanOfDays    = "Thu - Sat"
	complexDays   = "Fri - Sun Wed"
)

func TestDaysParse(t *testing.T) {
	c, done, err := aetest.NewContext()
	//ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	r := GetDays(simpleDays, c)
	if len(r) != 3 || !hasDays(r, time.Monday, time.Tuesday, time.Wednesday) {
		t.Errorf("Failed simple days parse: %s", r)
	}

	r = GetDays(separatedDays, c)
	if len(r) != 3 || !hasDays(r, time.Monday, time.Tuesday, time.Wednesday) {
		t.Errorf("Failed separated days parse: %s", r)
	}

	r = GetDays(spanOfDays, c)
	if len(r) != 3 || !hasDays(r, time.Thursday, time.Friday, time.Saturday) {
		t.Errorf("Failed span of days parse: %s", r)
	}

	r = GetDays(complexDays, c)
	if len(r) != 4 || !hasDays(r, time.Friday, time.Saturday, time.Sunday, time.Wednesday) {
		t.Errorf("Failed complex days parse: %s", r)
	}
}

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
