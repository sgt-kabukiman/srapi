// Copyright (c) 2015, Sgt. Kabukiman | MIT licensed

package srapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ErrorBadJSON represents an invalid response from the API server, usually due to
// server-side downtimes or bugs.
const ErrorBadJSON = 900

// ErrorNetwork represents connection timeouts and other network issues.
const ErrorNetwork = 901

// ErrorBadURL represents the unlikely case of trying to fetch an invalid URL.
// Except for bugs in this package, this should never occur.
const ErrorBadURL = 902

// ErrorBadLogic represents a programmer mistake, like trying to get a leaderboard
// without specifying the game and category.
const ErrorBadLogic = 903

// ErrorNoSuchLink represents the case when the package wants to follow a link
// in the resource which is suddenly not present. As the code relies on links to
// move around, this is bad.
const ErrorNoSuchLink = 904

// BaseURL is the base URL for all API calls.
const BaseURL = "http://www.speedrun.com/api/v1"

// our http client, initialized by init
var httpClient *apiClient

// initialize the httpClient
func init() {
	httpClient = &apiClient{
		baseURL: BaseURL,
		client:  &http.Client{},
	}
}

// request represents all options relevant for making an actual HTTP request.
// These fields are mapped to a http.Request when performing a request.
type request struct {
	// HTTP method, like "GET"
	method string

	// the URL, relative to BaseURL
	url string

	// optional filter (will be applied to the query string)
	filter filter

	// optional sorting (will be applied to the query string)
	sorting *Sorting

	// optional cursor (will be applied to the query string)
	cursor *Cursor
}

// apiClient is our helper to not pollute the package-wide variables
type apiClient struct {
	// the effective base url
	baseURL string

	// the underlying, concurrency-safe HTTP client
	client *http.Client
}

// do performs a HTTP request by transforming the request and applying the
// filters. The data is parsed as JSON and unmarshaled into dst. An error is
// returned when the request failed or when invalid JSON was received.
func (ac *apiClient) do(request request, dst interface{}) *Error {
	// prepare the actual net.http.Request
	u, err := url.Parse(ac.baseURL + request.url)
	if err != nil {
		return failedRequest(request, nil, err, ErrorBadURL)
	}

	request.filter.applyToURL(u)

	if request.cursor != nil {
		request.cursor.applyToURL(u)
	}

	if request.sorting != nil {
		request.sorting.applyToURL(u)
	}

	req := http.Request{
		Method: request.method,
		URL:    u,
	}

	// hit the network
	response, err := ac.client.Do(&req)
	if err != nil {
		return failedRequest(request, nil, err, ErrorNetwork)
	}

	// decode a successful response
	if response.StatusCode == 200 || response.StatusCode == 201 {
		defer response.Body.Close()

		err = json.NewDecoder(response.Body).Decode(dst)
		if err != nil {
			return failedRequest(request, nil, err, ErrorBadJSON)
		}

		// everything went fine
		return nil
	}

	// something went wrong
	return failedRequest(request, response, nil, 0)
}

// Error is an error that occured in this package. It contains basic information
// about the failed request (if any, some errors are independent of requests)
// and about what failed.
type Error struct {
	// the HTTP method of the request that failed, empty if no request involved
	Method string

	// the URL that failed, empty if no request involved
	URL string

	// the HTTP status code, set to one of the Error* constants of this package for
	// internal errors
	Status int

	// a description of what failed
	Message string
}

// Error returns a string including all details of the Error struct.
func (e *Error) Error() string {
	return fmt.Sprintf("[%d] %s (%s %s)", e.Status, e.Message, e.Method, e.URL)
}

// failedRequest is a helper to assemble a Error struct when a request failed.
func failedRequest(request request, response *http.Response, previous error, errorCode int) *Error {
	// build an incomplete error struct
	result := &Error{
		Method: request.method,
		URL:    request.url,
	}

	// decode the body into an Error
	if previous != nil {
		result.Status = errorCode
		result.Message = previous.Error()
	} else {
		defer response.Body.Close()
		err := json.NewDecoder(response.Body).Decode(result)
		if err != nil {
			result.Status = ErrorBadJSON
			result.Message = "Could not decode response body as JSON. Site is probably having issues."
		}
	}

	return result
}
