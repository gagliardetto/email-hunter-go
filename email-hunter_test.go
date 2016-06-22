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
	defer func() {
		client.domain = "https://api.emailhunter.co/"
	}()
	response, _, err := client.fetchAndReturnPage(path, method, headers, queryParameters, bodyPayload)
	if err != nil {
		t.Error("reading reponse body: %v, want %q", err, testBody)
	}

	if gotBody, wantBody := string(response), testBody+"\n"; gotBody != wantBody {
		t.Errorf("request body mismatch: got %q, want %q", gotBody, wantBody)
	}

}

func TestDomainSearch(t *testing.T) {
	domainSearchOptions := DomainSearchOptions{
		Domain: "stripe.com",
	}
	_, err := client.DomainSearch(domainSearchOptions)
	if err != nil {
		t.Error("calling DomainSearch: %q", err)
	}
}

func TestEmailFinder(t *testing.T) {
	emailFinderOptions := EmailFinderOptions{
		Domain:    "asana.com",
		FirstName: "Dustin",
		LastName:  "Moskovitz",
	}
	_, err := client.EmailFinder(emailFinderOptions)
	if err != nil {
		t.Error("calling EmailFinder: %q", err)
	}
}

func TestEmailVerification(t *testing.T) {
	_, err := client.EmailVerification("steli@close.io")
	if err != nil {
		t.Error("calling EmailVerification: %q", err)
	}
}

func TestEmailCount(t *testing.T) {
	_, err := client.EmailCount("stripe.com")
	if err != nil {
		t.Error("calling EmailCount: %q", err)
	}
}

func TestAccountInformation(t *testing.T) {
	_, err := client.AccountInformation()
	if err != nil {
		t.Error("calling AccountInformation: %q", err)
	}
}
