package emailHunter

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func NewClient(APIKey string) (*Client, error) {
	if len(APIKey) == 0 {
		return &Client{}, errors.New("APIKey not valid")
	}

	return &Client{
		httpClient: http.DefaultClient,
		APIKey:     APIKey,
		domain:     "https://api.emailhunter.co/",
	}, nil
}

func (client *Client) fetchAndReturnPage(path string, method string, headers http.Header, queryParameters url.Values, bodyPayload interface{}) ([]byte, http.Header, error) {

	if client.APIKey == "" {
		return []byte(""), http.Header{}, fmt.Errorf("%s", "APIKey not provided")
	}
	queryParameters.Add("api_key", client.APIKey)

	requestURL, err := url.Parse(client.domain)
	if err != nil {
		return []byte(""), http.Header{}, err
	}
	requestURL.Path = path
	requestURL.RawQuery = queryParameters.Encode()

	if method != "GET" && method != "POST" && method != "PUT" && method != "PATCH" && method != "DELETE" {
		return []byte(""), http.Header{}, fmt.Errorf("Method not supported: %v", method)
	}

	encodedBody, err := json.Marshal(bodyPayload)
	if err != nil {
		return []byte(""), http.Header{}, err
	}

	//fmt.Println(requestURL.String())
	request, err := http.NewRequest(method, requestURL.String(), bytes.NewBuffer(encodedBody))
	if err != nil {
		return []byte(""), http.Header{}, fmt.Errorf("Failed to get the URL %s: %s", requestURL, err)
	}
	request.Header = headers
	request.Header.Add("Content-Length", strconv.Itoa(len(encodedBody)))

	request.Header.Add("Connection", "Keep-Alive")
	request.Header.Add("Accept-Encoding", "gzip, deflate")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("User-Agent", "github.com/gagliardetto/email-hunter-go")

	response, err := client.httpClient.Do(request)
	if err != nil {
		return []byte(""), http.Header{}, fmt.Errorf("Failed to get the URL %s: %s", requestURL, err)
	}
	defer response.Body.Close()

	var responseReader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		decompressedBodyReader, err := gzip.NewReader(response.Body)
		if err != nil {
			return []byte(""), http.Header{}, err
		}
		responseReader = decompressedBodyReader
		defer responseReader.Close()
	default:
		responseReader = response.Body
	}

	responseBody, err := ioutil.ReadAll(responseReader)
	if err != nil {
		return []byte(""), http.Header{}, err
	}

	if response.StatusCode > 299 || response.StatusCode < 199 {
		var errorDescription string

		switch response.StatusCode {
		case 200:
			errorDescription = "The request was successful"
		case 401:
			errorDescription = "No valid API key provided"
		case 429:
			errorDescription = "You have reached your usage limit. Upgrade your plan if necessary"
		default:
			errorDescription = "Unknown error"
		}

		if response.StatusCode >= 500 && response.StatusCode <= 599 {
			errorDescription = "Something went wrong on Email Hunter's end"
		}

		return []byte(""), http.Header{}, fmt.Errorf("HTTPStatus %s: %s", strconv.Itoa(response.StatusCode), errorDescription)
	}

	return responseBody, response.Header, nil
}

func (client *Client) DomainSearch(options DomainSearchOptions) (DomainSearchResults, error) {

	path := fmt.Sprintf("/v1/search")
	method := "GET"

	headers := http.Header{}
	queryParameters := url.Values{}

	if len(options.Domain) > 0 {
		queryParameters.Add("domain", options.Domain)
	}

	if len(options.Company) > 0 {
		queryParameters.Add("company", options.Company)
	}

	if len(options.Domain) == 0 && len(options.Company) == 0 {
		return DomainSearchResults{}, errors.New("You must provide at least Domain or Company (or both)")
	}

	if options.Offset > 0 {
		queryParameters.Add("offset", strconv.FormatInt(options.Offset, 10))
	}

	if options.Type == GenericType || options.Type == PersonalType {
		queryParameters.Add("type", options.Type)
	}

	var bodyPayload interface{}

	response, _, err := client.fetchAndReturnPage(path, method, headers, queryParameters, bodyPayload)
	if err != nil {
		return DomainSearchResults{}, err
	}

	var domainSearchResults DomainSearchResults
	err = json.Unmarshal(response, &domainSearchResults)
	if err != nil {
		return DomainSearchResults{}, err
	}
	return domainSearchResults, nil
}

func (client *Client) EmailFinder(options EmailFinderOptions) (EmailFinderResults, error) {

	path := fmt.Sprintf("/v1/generate")
	method := "GET"

	headers := http.Header{}
	queryParameters := url.Values{}

	if len(options.Domain) > 0 {
		queryParameters.Add("domain", options.Domain)
	}

	if len(options.Company) > 0 {
		queryParameters.Add("company", options.Company)
	}

	if len(options.Domain) == 0 && len(options.Company) == 0 {
		return EmailFinderResults{}, errors.New("You must provide at least Domain or Company (or both)")
	}

	if len(options.FirstName) > 0 {
		queryParameters.Add("first_name", options.FirstName)
	} else {
		return EmailFinderResults{}, errors.New("You need to provide a first name")
	}

	if len(options.LastName) > 0 {
		queryParameters.Add("last_name", options.LastName)
	} else {
		return EmailFinderResults{}, errors.New("You need to provide a last name")
	}

	var bodyPayload interface{}

	response, _, err := client.fetchAndReturnPage(path, method, headers, queryParameters, bodyPayload)
	if err != nil {
		return EmailFinderResults{}, err
	}

	var emailFinderResults EmailFinderResults
	err = json.Unmarshal(response, &emailFinderResults)
	if err != nil {
		return EmailFinderResults{}, err
	}
	return emailFinderResults, nil
}

func (client *Client) EmailVerification(email string) (EmailVerificationResults, error) {

	path := fmt.Sprintf("/v1/verify")
	method := "GET"

	headers := http.Header{}
	queryParameters := url.Values{}

	if len(email) > 0 {
		queryParameters.Add("email", email)
	} else {
		return EmailVerificationResults{}, errors.New("You need to provide an email")
	}

	var bodyPayload interface{}

	response, _, err := client.fetchAndReturnPage(path, method, headers, queryParameters, bodyPayload)
	if err != nil {
		return EmailVerificationResults{}, err
	}

	var emailVerificationResults EmailVerificationResults
	err = json.Unmarshal(response, &emailVerificationResults)
	if err != nil {
		return EmailVerificationResults{}, err
	}
	return emailVerificationResults, nil
}

func (client *Client) EmailCount(domain string) (EmailCountResults, error) {

	path := fmt.Sprintf("/v1/email-count")
	method := "GET"

	headers := http.Header{}
	queryParameters := url.Values{}

	if len(domain) > 0 {
		queryParameters.Add("domain", domain)
	} else {
		return EmailCountResults{}, errors.New("You need to provide a domain")
	}

	var bodyPayload interface{}

	response, _, err := client.fetchAndReturnPage(path, method, headers, queryParameters, bodyPayload)
	if err != nil {
		return EmailCountResults{}, err
	}

	var emailCountResults EmailCountResults
	err = json.Unmarshal(response, &emailCountResults)
	if err != nil {
		return EmailCountResults{}, err
	}
	return emailCountResults, nil
}

func (client *Client) AccountInformation() (AccountInformationResults, error) {

	path := fmt.Sprintf("/v1/account")
	method := "GET"

	headers := http.Header{}
	queryParameters := url.Values{}

	var bodyPayload interface{}

	response, _, err := client.fetchAndReturnPage(path, method, headers, queryParameters, bodyPayload)
	if err != nil {
		return AccountInformationResults{}, err
	}

	var accountInformationResults AccountInformationResults
	err = json.Unmarshal(response, &accountInformationResults)
	if err != nil {
		return AccountInformationResults{}, err
	}
	return accountInformationResults, nil
}

func (ct *Date) UnmarshalJSON(b []byte) error {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	var err error
	ct.Time, err = time.Parse(TimestampFormat, string(b))
	return err
}

func (ct *Date) MarshalJSON() ([]byte, error) {
	return []byte(ct.Time.Format(TimestampFormat)), nil
}
