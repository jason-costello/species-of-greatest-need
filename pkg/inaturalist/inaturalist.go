// Package inaturalist interfaces with the inaturalist api and the structure is largley taken from
// from https://github.com/Medium/medium-sdk-go/blob/master/mediui.go as a way to learn different
// way to write a service
package inaturalist

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const(
	// host is the default host of inaturalist api
	host = "https://api.inaturalist.org/v1"

	// defaultTimeout is  the default duration of http requests before timeout occurs.
	defaultTimeout = 15 * time.Second

	// defaultCode is the default error code for failures
	defaultCode = -1
)

// marshalling formats for requests
const(
	formatJSON = "json"
	formatCSV = "csv"
)


// Error defines an error received when making a request to the api.
type Error struct{
	Message string `json:"message"`
	Code int `json:"code""`
}

// Error returns a string representing the error, implementing the error interface
func (e Error) Error() string{
	return fmt.Sprintf("inat: %s (%d)", e.Message, e.Code)
}


// INaturalist defines the iNaturalist Client
type INaturalist struct{

	Host string
	Timeout time.Duration
	Transport http.RoundTripper

}


// TaxonDetail defines the returned data from GetTaxonDetails
type TaxonDetail struct{

}


// NewClient returns a new INaturalist client which can be used to make requests
func NewClient() *INaturalist{
	return &INaturalist{
		Host:      host,
		Timeout:   defaultTimeout,
		Transport: http.DefaultTransport,
	}
}


func(i *INaturalist) GetTaxonDetails(taxonID ...int) (*TaxaResponse, error){
	var r clientRequest
	if taxonID == nil {
		return nil, Error{ fmt.Sprintf("No Taxon ID provided: %s", errors.New("nil taxon")), defaultCode}
	}

	var tIDS string
	for i, x := range taxonID{
		if i != 0 {
			tIDS = fmt.Sprintf("%s,%d", tIDS, x)

		} else {
			tIDS = fmt.Sprintf("%d", x)
		}
	}
		r = clientRequest{
			method: "GET",
			path:   fmt.Sprintf("/taxa/%s", tIDS),
		}
	t := &TaxaResponse{}
	err := i.request(r, t)


	spew.Dump(t)
	return t, err

}


func buildQueryString(path string, op ObservationParameters) string{

	var qs string
	var cnt int


	for k,v := range op{
		if cnt == 0 {
			qs += fmt.Sprintf("%s=%s", k, v)
		} else {
			qs += fmt.Sprintf("&%s=%s", k, v)

		}
		cnt++

	}



	if qs != "" {
		return fmt.Sprintf("%s?%s", path, qs)
	}
	return path

}
func(i *INaturalist) Observations(observationParameters ObservationParameters) (*ObservationResponse, error) {
//var d1 string // observed on or after this date
//var d2 string // observed on or before this date
//var created_d1 string // created on or after this datetime
//var created_d2 string // created on or before this datetime
//var acc_below string // positional accuracy below this value (meters)
//var acc_above string // positional accuracy above this value (meters)


path := buildQueryString("/observations", observationParameters)




	var r clientRequest




	r = clientRequest{
		method: "GET",
		path:   path,
	}

	t := &ObservationResponse{}
	err := i.request(r, t)

	tjson, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(string(tjson))
	return t, err


}

// generateJSONRequestData returns the body and content type for a JSON request.
func (i *INaturalist) generateJSONRequestData(cr clientRequest) ([]byte, string, error) {
	body, err := json.Marshal(cr.data)
	if err != nil {
		return nil, "", Error{fmt.Sprintf("Could not marshal JSON: %s", err), defaultCode}
	}
	return body, "application/json", nil
}


// generateJSONRequestData returns the body and content type for a JSON request.
func (i *INaturalist) generateCSVRequestData(cr clientRequest) ([]byte, string, error) {
	body, err := json.Marshal(cr.data)
	if err != nil {
		return nil, "", Error{fmt.Sprintf("Could not marshal CSV: %d", err), defaultCode}
	}
	return body, "text/csv", nil
}

// generateFormRequestData returns the body and content type for a form data request.
//func (i *INaturalist) generateFormRequestData(cr clientRequest) ([]byte, string, error) {
//	var body []byte
//	switch d := cr.data.(type) {
//	case string:
//		body = []byte(d)
//	case []byte:
//		body = d
//	default:
//		return nil, "", Error{"Invalid data passed for form request", defaultCode}
//	}
//	return body, "application/x-www-form-urlencoded", nil
//}


// request makes a request to iNaturalist's API
func (i *INaturalist) request(cr clientRequest, result interface{}) error {
	f := cr.format
	if f == "" {
		f = formatJSON
	}

	// Get the body and content type.
	var g requestDataGenerator
	switch f {
	case formatJSON:
		g = i.generateJSONRequestData
	case formatCSV:
		g = i.generateCSVRequestData

	default:
		return Error{fmt.Sprintf("Unknown format: %s", cr.format), defaultCode}
	}
	body, ct, err := g(cr)
	if err != nil {
		return err
	}

	// Construct the request
	req, err := http.NewRequest(cr.method, i.Host+cr.path, bytes.NewReader(body))
	if err != nil {
		return Error{fmt.Sprintf("Could not create request: %s", err), defaultCode}
	}

	req.Header.Add("Content-Type", ct)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Charset", "utf-8")

	// Create the HTTP client
	client := &http.Client{
		Transport: i.Transport,
		Timeout:   i.Timeout,
	}

	fmt.Println(req.URL)
	// Make the request
	res, err := client.Do(req)
	if err != nil {
		return Error{fmt.Sprintf("Failed to make request: %s", err), defaultCode}
	}
	defer res.Body.Close()

	// Parse the response
	c, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Error{fmt.Sprintf("Could not read response: %s", err), defaultCode}
	}

	var env envelope
	if err := json.Unmarshal(c, &env); err != nil {
		return Error{fmt.Sprintf("Could not parse response: %s", err), defaultCode}
	}

	if http.StatusOK <= res.StatusCode && res.StatusCode < http.StatusMultipleChoices {
		if env.Data != nil {
			c, _ = json.Marshal(env.Data)
		}
		return json.Unmarshal(c, &result)
	}
	e := env.Errors[0]
	return Error{e.Message, e.Code}
}




// clientRequest defines information used to make a request to inaturalist
type clientRequest struct{
	method string
	path string
	data interface{}
	format string
}



// payload defines a struct to represent payloads that are returned from inaturalist.
type envelope struct {
	Data   interface{} `json:"data"`
	Errors []Error     `json:"errors,omitempty"`
}

// requestDataGenerator defines a function that can generate request data.
type requestDataGenerator func(cr clientRequest) ([]byte, string, error)


// Borrowed from multipart/writer.go
var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

// escapeQuotes returns the supplied string with quotes escaped.
func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}


type TaxaResponse struct {
	TotalResults int `json:"total_results"`
	Page         int `json:"page"`
	PerPage      int `json:"per_page"`
	Results      []struct {
		ObservationsCount int    `json:"observations_count"`
		TaxonSchemesCount int    `json:"taxon_schemes_count"`
		Ancestry          string `json:"ancestry"`
		IsActive          bool   `json:"is_active"`
		FlagCounts        struct {
			Unresolved int `json:"unresolved"`
			Resolved   int `json:"resolved"`
		} `json:"flag_counts"`
		WikipediaURL              string      `json:"wikipedia_url"`
		CurrentSynonymousTaxonIds interface{} `json:"current_synonymous_taxon_ids"`
		IconicTaxonID             int         `json:"iconic_taxon_id"`
		RankLevel                 int         `json:"rank_level"`
		TaxonChangesCount         int         `json:"taxon_changes_count"`
		AtlasID                   interface{} `json:"atlas_id"`
		CompleteSpeciesCount      interface{} `json:"complete_species_count"`
		ParentID                  int         `json:"parent_id"`
		Name                      string      `json:"name"`
		Rank                      string      `json:"rank"`
		Extinct                   bool        `json:"extinct"`
		ID                        int         `json:"id"`
		DefaultPhoto              struct {
			SquareURL          string        `json:"square_url"`
			Attribution        string        `json:"attribution"`
			Flags              []interface{} `json:"flags"`
			MediumURL          string        `json:"medium_url"`
			ID                 int           `json:"id"`
			LicenseCode        string        `json:"license_code"`
			OriginalDimensions interface{}   `json:"original_dimensions"`
			URL                string        `json:"url"`
		} `json:"default_photo"`
		AncestorIds         []int  `json:"ancestor_ids"`
		MatchedTerm         string `json:"matched_term"`
		IconicTaxonName     string `json:"iconic_taxon_name"`
		PreferredCommonName string `json:"preferred_common_name"`
	} `json:"results"`
}


type ObservationParameters map[string]string


type ObservationResponse struct {
	TotalResults int `json:"total_results"`
	Page         int `json:"page"`
	PerPage      int `json:"per_page"`
	Results      []struct {
	ID               int  `json:"id"`
	CachedVotesTotal int  `json:"cached_votes_total"`
	Captive          bool `json:"captive"`
	Comments         []struct {
	ID               int       `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedAtDetails struct {
	Date  string `json:"date"`
	Day   int    `json:"day"`
	Hour  int    `json:"hour"`
	Month int    `json:"month"`
	Week  int    `json:"week"`
	Year  int    `json:"year"`
	} `json:"created_at_details"`
	User struct {
	ID              int    `json:"id"`
	IconContentType string `json:"icon_content_type"`
	IconFileName    string `json:"icon_file_name"`
	IconURL         string `json:"icon_url"`
	Login           string `json:"login"`
	Name            string `json:"name"`
	} `json:"user"`
	} `json:"comments"`
	CommentsCount    int       `json:"comments_count"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedAtDetails struct {
	Date  string `json:"date"`
	Day   int    `json:"day"`
	Hour  int    `json:"hour"`
	Month int    `json:"month"`
	Week  int    `json:"week"`
	Year  int    `json:"year"`
	} `json:"created_at_details"`
	CreatedTimeZone string `json:"created_time_zone"`
	Description     string `json:"description"`
	FavesCount      int    `json:"faves_count"`
	Geojson         struct {
	Type        string   `json:"type"`
	Coordinates []interface{} `json:"coordinates"`
	} `json:"geojson"`
	Geoprivacy                  string `json:"geoprivacy"`
	TaxonGeoprivacy             string `json:"taxon_geoprivacy"`
	IDPlease                    bool   `json:"id_please"`
	IdentificationsCount        int    `json:"identifications_count"`
	IdentificationsMostAgree    bool    `json:"identifications_most_agree"`
	IdentificationsMostDisagree bool    `json:"identifications_most_disagree"`
	IdentificationsSomeAgree    bool    `json:"identifications_some_agree"`
	LicenseCode                 string `json:"license_code"`
	Location                    string `json:"location"`
	Mappable                    bool   `json:"mappable"`
	NonOwnerIds                 []struct {
	ID               int       `json:"id"`
	Body             string    `json:"body"`
	CreatedAt        time.Time `json:"created_at"`
	CreatedAtDetails struct {
	Date  string `json:"date"`
	Day   int    `json:"day"`
	Hour  int    `json:"hour"`
	Month int    `json:"month"`
	Week  int    `json:"week"`
	Year  int    `json:"year"`
	} `json:"created_at_details"`
	User struct {
	ID              int    `json:"id"`
	IconContentType string `json:"icon_content_type"`
	IconFileName    string `json:"icon_file_name"`
	IconURL         string `json:"icon_url"`
	Login           string `json:"login"`
	Name            string `json:"name"`
	} `json:"user"`
	} `json:"non_owner_ids"`
	NumIdentificationAgreements    int       `json:"num_identification_agreements"`
	NumIdentificationDisagreements int       `json:"num_identification_disagreements"`
	Obscured                       bool      `json:"obscured"`
	ObservedOn                     string `json:"observed_on"`
	ObservedOnDetails              struct {
	Date  string `json:"date"`
	Day   int    `json:"day"`
	Hour  int    `json:"hour"`
	Month int    `json:"month"`
	Week  int    `json:"week"`
	Year  int    `json:"year"`
	} `json:"observed_on_details"`
	ObservedOnString string `json:"observed_on_string"`
	ObservedTimeZone string `json:"observed_time_zone"`
	Ofvs             []struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	} `json:"ofvs"`
	OutOfRange bool `json:"out_of_range"`
	Photos     []struct {
	ID          int    `json:"id"`
	Attribution string `json:"attribution"`
	LicenseCode string `json:"license_code"`
	URL         string `json:"url"`
	} `json:"photos"`
	PlaceGuess                 string   `json:"place_guess"`
	PlaceIds                   []int `json:"place_ids"`
	ProjectIds                 []int `json:"project_ids"`
	ProjectIdsWithCuratorID    []int `json:"project_ids_with_curator_id"`
	ProjectIdsWithoutCuratorID []int `json:"project_ids_without_curator_id"`
	QualityGrade               string   `json:"quality_grade"`
	ReviewedBy                 []int `json:"reviewed_by"`
	SiteID                     int   `json:"site_id"`
	Sounds                     []struct {
	ID          int    `json:"id"`
	Attribution string `json:"attribution"`
	LicenseCode string `json:"license_code"`
	} `json:"sounds"`
	SpeciesGuess string   `json:"species_guess"`
	Tags         []string `json:"tags"`
	Taxon        struct {
	ID                  int    `json:"id"`
	IconicTaxonID       int    `json:"iconic_taxon_id"`
	IconicTaxonName     string `json:"iconic_taxon_name"`
	IsActive            bool   `json:"is_active"`
	Name                string `json:"name"`
	PreferredCommonName string `json:"preferred_common_name"`
	Rank                string `json:"rank"`
	RankLevel           int    `json:"rank_level"`
	AncestorIds         []int  `json:"ancestor_ids"`
	Ancestry            string `json:"ancestry"`
	ConservationStatus  struct {
	SourceID   int    `json:"source_id"`
	Authority  string `json:"authority"`
	Status     string `json:"status"`
	StatusName string `json:"status_name"`
	Iucn       int    `json:"iucn"`
	Geoprivacy string `json:"geoprivacy"`
	} `json:"conservation_status"`
	Endemic            bool `json:"endemic"`
	EstablishmentMeans struct {
	EstablishmentMeans string `json:"establishment_means"`
	Place              struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	} `json:"place"`
	} `json:"establishment_means"`
	Introduced bool `json:"introduced"`
	Native     bool `json:"native"`
	Threatened bool `json:"threatened"`
	} `json:"taxon"`
	TimeObservedAt string `json:"time_observed_at"`
	TimeZoneOffset string    `json:"time_zone_offset"`
	UpdatedAt      string `json:"updated_at"`
	URI            string    `json:"uri"`
	User           struct {
	ID              int    `json:"id"`
	IconContentType string `json:"icon_content_type"`
	IconFileName    string `json:"icon_file_name"`
	IconURL         string `json:"icon_url"`
	Login           string `json:"login"`
	Name            string `json:"name"`
	} `json:"user"`
	Verifiable bool `json:"verifiable"`
	} `json:"results"`
	}
