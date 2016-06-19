package emailHunter

import (
	"time"
)

const (
	PersonalType string = "personal"
	GenericType  string = "generic"
)

type DomainSearchOptions struct {
	Domain  string
	Company string
	Offset  int64
	Type    string
}

type DomainSearchResults struct {
	Status  string  `json:"status"`
	Domain  string  `json:"domain"`
	Results int64   `json:"results"`
	Webmail bool    `json:"webmail"`
	Pattern string  `json:"pattern"`
	Offset  int     `json:"offset"`
	Emails  []Email `json:"emails"`
}

type Email struct {
	Value      string   `json:"value"`
	Type       string   `json:"type"`
	Confidence float64  `json:"confidence"`
	Sources    []Source `json:"sources"`
}

type Source struct {
	Domain       string `json:"domain"`
	URI          string `json:"uri"`
	Extracted_on Date   `json:"extracted_on"`
}

type EmailFinderOptions struct {
	Domain    string
	Company   string
	FirstName string
	LastName  string
}

type EmailFinderResults struct {
	Status  string   `json:"status"`
	Email   string   `json:"email"`
	Score   float64  `json:"score"`
	Sources []Source `json:"sources"`
}

type EmailVerificationResults struct {
	Status     string   `json:"status"`
	Email      string   `json:"email"`
	Score      float64  `json:"score"`
	Result     string   `json:"result"`
	Regexp     bool     `json:"regexp"`
	Gibberish  bool     `json:"gibberish"`
	Disposable bool     `json:"disposable"`
	Webmail    bool     `json:"webmail"`
	MXRecords  bool     `json:"mx_records"`
	SMTPServer bool     `json:"smtp_server"`
	SMTPCheck  bool     `json:"smtp_check"`
	AcceptAll  bool     `json:"accept_all"`
	Sources    []Source `json:"sources"`
}

type EmailCountResults struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type AccountInformationResults struct {
	Status    string `json:"status"`
	Email     string `json:"email"`
	PlanName  string `json:"plan_name"`
	PlanLevel int    `json:"plan_level"`
	ResetDate Date   `json:"reset_date"`
	Calls     struct {
		Used      int64 `json:"used"`
		Available int64 `json:"available"`
	}
}

type Date struct {
	time.Time
}

const TimestampFormat = "2006-01-02"

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
