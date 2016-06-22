package emailHunter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

var (
	APIKey string
	client *Client
	err    error
)

func init() {
	APIKey = GetTestKey()

	client, err = NewClient(APIKey)
	if err != nil {
		fmt.Println("client setup error: ", err)
		return
	}
}

func GetTestKey() string {
	key := os.Getenv("EMAIL_HUNTER_KEY")

	if len(key) == 0 {
		panic("EMAIL_HUNTER_KEY environment variable is not set, but is needed to run tests!\n")
	}

	return key
}

func TestFetchAndReturnPage(t *testing.T) {
	// server serves json
	// TestFetchAndReturnPage asks for that json
	// if the two are equal, success

	testBody := `{some:"json"}`

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, testBody)
	}))
	defer testServer.Close()

	path := fmt.Sprintf(testServer.URL)
	method := "GET"
	headers := http.Header{}
	queryParameters := url.Values{}
	var bodyPayload interface{}
	client.domain = ""
	response, _, err := client.fetchAndReturnPage(path, method, headers, queryParameters, bodyPayload)
	if err != nil {
		t.Error("reading reponse body: %v, want %q", err, testBody)
	}

	if gotBody, wantBody := string(response), testBody+"\n"; gotBody != wantBody {
		t.Errorf("request body mismatch: got %q, want %q", gotBody, wantBody)
	}

}

func main() {

	domainSearchOptions := DomainSearchOptions{
		Domain: "stripe.com",
	}
	domainSearchResults, err := client.DomainSearch(domainSearchOptions)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	fmt.Printf("%#v\n\n", domainSearchResults)

	emailFinderOptions := EmailFinderOptions{
		Domain:    "asana.com",
		FirstName: "Dustin",
		LastName:  "Moskovitz",
	}
	emailFinderResults, err := client.EmailFinder(emailFinderOptions)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	fmt.Printf("%#v\n\n", emailFinderResults)

	emailVerificationResults, err := client.EmailVerification("steli@close.io")
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	fmt.Printf("%#v\n\n", emailVerificationResults)

	emailCountResults, err := client.EmailCount("stripe.com")
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	fmt.Printf("%#v\n\n", emailCountResults)

	accountInformationResults, err := client.AccountInformation()
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	fmt.Printf("%#v\n\n", accountInformationResults)
}
