package inaturalist

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

var i = NewClient()

type apiTest struct {
	fn          interface{}
	payload     interface{}
	method      string
	path        string
	contentType string
	bodyPattern string
}
var 	 testOP = ObservationParameters{}

var apiTests = []apiTest{
	{ i.GetTaxonDetails, 2222,
		"GET", "/taxa/2222", "application/json",
		"null"},
	{ i.Observations, ObservationParameters{"taxon_id":"1", "place_id":"18","quality_grade":"research"},
		"GET", "/observations?taxon_id=1&place_id=18&quality_grade=research", "application/json",
		"null"},
}

// TestAPIMethods tests that http requests are constructed correctly.
func TestApiMethods(t *testing.T){




	var body []byte
	var req *http.Request

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		req = r
		body, _ = ioutil.ReadAll(req.Body)
	}))

	defer ts.Close()
	i.Host = ts.URL


	for _, tt := range apiTests{

		f := reflect.ValueOf(tt.fn)
		var pl []reflect.Value

		if tt.payload != nil{
			pl = append(pl, reflect.ValueOf(tt.payload))
		}
		f.Call(pl)

		// Test request was correctly formed.
		assertEqual(t, req.Header.Get("Accept"), "application/json")
		assertMatch(t, req.Header.Get("Content-Type"), tt.contentType)
		assertEqual(t, req.Method, tt.method)
		if req.URL.RawQuery != ""{
		assertEqual(t, fmt.Sprintf("%s?%s",req.URL.Path,req.URL.RawQuery), tt.path)
		} else {
			assertEqual(t, req.URL.Path, tt.path)
		}
		assertMatch(t, string(body), tt.bodyPattern)
	}

}


// TestAPITimeout verifies that HTTP timeouts work
func TestAPITimeout(t *testing.T) {
	i.Timeout = 1 * time.Millisecond
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(i.Timeout * 2) // sleep longer than timeout
		fmt.Fprintln(w, "null")
	}))
	defer ts.Close()
	i.Host = ts.URL

	_, err := i.GetTaxonDetails(1)
	if err == nil {
		t.Errorf("Expected HTTP timeout error, but call succeeded")
	} else if !strings.Contains(err.Error(), "Client.Timeout exceeded") {
		// go1.4 doesn't set the timeout error but closes the connection.
		if !strings.Contains(err.Error(), "use of closed network connection") {
			t.Errorf("Expected HTTP timeout error, got %s", err)
		}
	}
}

func assertEqual(t *testing.T, actual, expected interface{}) {
	if actual != expected {
		t.Errorf("Expected %#v, got %#v", expected, actual)
	}
}

func assertMatch(t *testing.T, actual, pattern string) {
	re := regexp.MustCompile(pattern)
	if !re.MatchString(actual) {
		t.Errorf("Expected to match %#v, got %#v", pattern, actual)
	}
}
