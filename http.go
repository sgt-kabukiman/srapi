package srapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const ErrorBadJSON = 900
const ErrorNetwork = 901
const ErrorBadURL = 902
const ErrorBadLogic = 903

const BaseUrl = "http://www.speedrun.com/api/v1"

var httpClient *apiClient

func init() {
	httpClient = &apiClient{
		baseUrl: BaseUrl,
		client:  &http.Client{},
	}
}

type request struct {
	method  string
	url     string
	filter  filter
	sorting *Sorting
	cursor  *Cursor
}

type apiClient struct {
	baseUrl string
	client  *http.Client
}

func (self *apiClient) do(request request, dst interface{}) *Error {
	// prepare the actual net.http.Request
	u, err := url.Parse(self.baseUrl + request.url)
	if err != nil {
		return failedRequest(request, nil, err, ErrorBadURL)
	}

	if request.filter != nil {
		request.filter.applyToURL(u)
	}

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
	response, err := self.client.Do(&req)
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

type Error struct {
	Method  string
	URL     string
	Status  int
	Message string
}

func (self *Error) Error() string {
	return fmt.Sprintf("[%d] %s (%s %s)", self.Status, self.Message, self.Method, self.URL)
}

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
