package rap

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"

	"testing"
	"time"
)

func TestImport(t *testing.T) {
	c, done, err := aetest.NewContext()
	//ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	testTimeParse(t, c)
	testDaysParse(t, c)
}

var (
	simpleTime     = "Mon Tue Wed 8:00AM - 6:00PM"
	separatedTime  = "Mon, Tue, & Wed 8:00AM - 6:00PM"
	spanOfTime     = "Thu - Sat 8:00AM - 6:00PM"
	differentTimes = "Mon - Wed 8:00AM - 6:00PM Sat 1:00PM - 3:00PM"
	complexTimes   = "Fri - Sun 8:00AM - 6:00PM Tue through Wed 9:00AM to 11:00am"
)

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
}

var (
	simpleDays     = "Mon Tue Wed"
	separatedDays  = "Mon, Tue, & Wed"
	spanOfDays     = "Thu - Sat"
	longSpanOfDays = "Tue - Sat"
	complexDays    = "Fri - Sun Wed"
)

func testDaysParse(t *testing.T, c context.Context) {
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

	r = GetDays(longSpanOfDays, c)
	if len(r) != 5 || !hasDays(r, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday) {
		t.Errorf("Failed long span of days parse: %s", r)
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
