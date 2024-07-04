package main

import (
	"cmp"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	AccessToken               = "AccessToken"
	ValidToken                = "583-asgl-1s4gh-789b"
	datasetPath               = "./dataset.xml"
	ageField                  = "Age"
	idField                   = "Id"
	nameField                 = "Name"
	internalServerErorrMarker = "SIMULATE_INTERNAL_SERVER_ERROR"
	replyInvalidJSON          = "REPLY_INVALID_JOSN"
)

var (
	datasetUsers               Users
	BadRequestError            error  = errors.New("ErrorBadOrderField")
	OrderByInvalidError        error  = errors.New("invalid order_by")
	InternalServerErrorContent []byte = []byte("{\"status\": 500, \"reason\": \"Internal Server Error\"}")
	invalidJsonResponse               = []byte("{\"some': \"invalid\", }")
)

type Users struct {
	Members []UserEntry `xml:"row"`
	ready   bool
}

type UserEntry struct {
	Id        int    `xml:"id"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

func (ue UserEntry) toUser() User {
	return User{Id: ue.Id, Age: ue.Age, About: ue.About, Name: ue.FirstName + " " + ue.LastName, Gender: ue.Gender}
}

// Parse xml file with users info.
// Executed on init of test.
func prepareSearchData() {
	if datasetUsers.ready { // xml file was parsed already
		return
	}
	data, err := os.ReadFile(datasetPath)
	if err != nil { // handle any problem of file read
		panic(fmt.Sprintf("error reading dataset file [%s]: %v", datasetPath, err))
	}
	parseResult := Users{Members: make([]UserEntry, 0, 40), ready: true}
	err = xml.Unmarshal(data, &parseResult)
	if err != nil || len(parseResult.Members) < 1 { // check errors and parse result
		panic(fmt.Sprintf("error parsing xml [%s]: %v", datasetPath, err))
	}
	datasetUsers = parseResult
}

func produceErrorResponse(reason string) SearchErrorResponse {
	return SearchErrorResponse{Error: reason}
}

func (e SearchErrorResponse) Msg() []byte {
	if msg, err := json.Marshal(e); err != nil {
		return InternalServerErrorContent
	} else {
		return msg
	}
}

func authorize(r *http.Request) (reason string, isAuthorized bool) {
	switch r.Header.Get(AccessToken) {
	case internalServerErorrMarker:
		panic("internal server error")
	case ValidToken:
		return "success", true
	default:
		return "invalid token", false
	}
}

// Extract integer value of search param. Otherwise - handle error
func getIntParam(vals url.Values, paramName string) (int, error) {
	err_ := fmt.Errorf("invalid integer param [%s]", paramName)
	if !vals.Has(paramName) {
		return -1, err_
	}
	result, e := strconv.Atoi(vals.Get(paramName))
	if e != nil {
		return -1, err_
	}
	return result, nil
}

// Function validates request search params.
// On empty search param - is valid ( use default).
// Otherwise verify is search param listed in allowed values
func validateAllowedValues(value interface{}, allowedValues ...interface{}) error {
	for i := 0; i < len(allowedValues); i++ {
		if value == allowedValues[i] {
			return nil
		}
	}
	return errors.New("invalid param")
}

func validateSearchParams(r *http.Request) (*SearchRequest, error) {
	q := r.URL.Query()
	limit, err := getIntParam(q, "limit")
	if err != nil {
		return nil, BadRequestError
	}
	offset, err := getIntParam(q, "offset")
	if err != nil {
		return nil, BadRequestError
	}
	orderBy, err := getIntParam(q, "order_by")
	if err != nil {
		return nil, BadRequestError
	}

	if err := validateAllowedValues(orderBy, OrderByAsc, OrderByAsIs, OrderByDesc); err != nil {
		return nil, OrderByInvalidError
	}
	if err := validateAllowedValues(q.Get("order_field"), "", ageField, idField, nameField); err != nil {
		return nil, errors.New("ErrorBadOrderField")
	}
	candidate := SearchRequest{Query: q.Get("query"), Limit: limit, Offset: offset, OrderField: q.Get("order_field"), OrderBy: orderBy}
	return &candidate, nil
}

func handleErrorResponse(w http.ResponseWriter, status int, reason string) {
	w.WriteHeader(status)
	if _, err := w.Write(produceErrorResponse(reason).Msg()); err != nil {
		panic("error to response")
	}
}

// Build sort function based on order_field and order_by from search params.
func resolveSortFunc(slice []User, searchParams *SearchRequest) func(i, j int) bool {
	switch searchParams.OrderField {
	default:
		return func(i, j int) bool { return cmp.Compare(slice[i].Name, slice[j].Name) == searchParams.OrderBy }
	case ageField:
		return func(i, j int) bool { return cmp.Compare(slice[i].Age, slice[j].Age) == searchParams.OrderBy }
	case idField:
		return func(i, j int) bool { return cmp.Compare(slice[i].Id, slice[j].Id) == searchParams.OrderBy }
	}
}

func sortUsersBeforeSearch(searchParams *SearchRequest, users []User) {
	if searchParams.OrderBy == OrderByAsIs {
		return
	}
	sort.Slice(users, resolveSortFunc(users, searchParams))
}

// Search predicate.
func matches(user *User, searchParams *SearchRequest) bool {
	return strings.Contains(user.Name, searchParams.Query) || strings.Contains(user.About, searchParams.Query)
}

func search(searchParams *SearchRequest, w http.ResponseWriter) {
	if searchParams.Query == replyInvalidJSON { // siulate invalid JSON response on search result
		w.WriteHeader(http.StatusOK)
		n, err := w.Write(invalidJsonResponse)
		if n != len(invalidJsonResponse) || err != nil {
			panic("failed to process response")
		}
		return
	}
	searchCopy := make([]User, len(datasetUsers.Members)) // make a copy to keep original sequence between search requests
	for i := 0; i < len(datasetUsers.Members); i++ {
		searchCopy[i] = datasetUsers.Members[i].toUser()
	}
	searchResult := make([]User, 0, len(searchCopy))
	for i := 0; i < len(searchCopy); i++ {
		if matches(&searchCopy[i], searchParams) {
			searchResult = append(searchResult, searchCopy[i])
		}
	}
	sortUsersBeforeSearch(searchParams, searchResult) // sort result if needed accordingly to search params
	response, err := json.Marshal(searchResult)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	n, err := w.Write(response)
	if n != len(response) || err != nil {
		panic("failed to process response")
	}
}

func init() {
	prepareSearchData() // Prepare datastore.
}

func InvalidJsonHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	if _, err := w.Write(invalidJsonResponse); err != nil {
		panic("error to response")
	}
}

func TimeOutHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1100 * time.Millisecond)
	handleErrorResponse(w, http.StatusGatewayTimeout, "gateway timed out")
}

func InvalidResponseHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(invalidJsonResponse); err != nil {
		panic("failed to process response")
	}
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	defer func(writer http.ResponseWriter) {
		if r := recover(); r != nil {
			handleErrorResponse(writer, http.StatusInternalServerError, "internal server error")
			return
		}
	}(w)
	// 1. authorize request.
	if _, authorized := authorize(r); !authorized {
		fmt.Println()
		handleErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	// 2. validate search params.
	searchParams, err := validateSearchParams(r)
	if err != nil {
		handleErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	// 3. search data -> handle errrors -> prodive response result.
	search(searchParams, w)
}

func TestTimeOut(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(TimeOutHandler))
	client := SearchClient{AccessToken: ValidToken, URL: ts.URL}
	_, err := client.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !strings.HasPrefix(err.Error(), "timeout for") {
		t.Errorf("expected time error, got %s", err)
	}
}

func TestJsonUnpackError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(InvalidJsonHandler))
	client := SearchClient{AccessToken: ValidToken, URL: ts.URL}
	_, err := client.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !strings.HasPrefix(err.Error(), "cant unpack error json") {
		t.Errorf("expected unpack error, got %s", err)
	}
}

func TestUnknownNetworkError(t *testing.T) {
	client := SearchClient{AccessToken: ValidToken, URL: "http://127.0.0.1:1234"}
	_, err := client.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !strings.HasPrefix(err.Error(), "unknown error") {
		t.Errorf("expected time error, got %s", err)
	}
}

func TestFailSearchServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	for caseNum, item := range failCases {
		client := item.client
		client.URL = ts.URL
		_, err := client.FindUsers(item.search)
		if err == nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if err.Error() != item.err.Error() { //make sence to compare exact error result with expected
			t.Errorf("[%d] expected error [%s] doesn't match got [%s]", caseNum, item.err, err)
		}
	}
	ts.Close()
}

func TestSucccessSearchServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	for caseNum, item := range successCases {
		client := item.client
		client.URL = ts.URL
		result, err := client.FindUsers(item.search)
		if err != nil {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if !reflect.DeepEqual(item.result, result) {
			t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, item.result, result)
		}
	}
	ts.Close()
}
