// Package gonbp is a go wrapper over NBP API
package gonbp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

const (
	apiBase = "https://api.nbp.pl/api/exchangerates/rates/A"
)

type rawRate struct {
	Table    string `json:"table"`
	Currency string `json:"currency"`
	Code     string `json:"code"`
	Rates    []struct {
		No            string          `json:"no"`
		EffectiveDate string          `json:"effectiveDate"`
		Mid           decimal.Decimal `json:"mid"`
	} `json:"rates"`
}

type httpClient interface {
	Get(url string) (resp *http.Response, err error)
}

// NBP is the NBP API client
type NBP struct {
	httpClient httpClient
}

// New returns *NBP instance with net/http.DefaultClient
func New() *NBP {
	return &NBP{
		httpClient: http.DefaultClient,
	}
}

// WithClient returns *NBP instance with given http.Client
func WithClient(client *http.Client) *NBP {
	return &NBP{
		httpClient: client,
	}
}

// Currency enumerates supported currencies
type Currency string

const (
	CHF Currency = "CHF"
	EUR Currency = "EUR"
	USD Currency = "USD"
)

// Rate represents the currency exchange rate for a given date
type Rate struct {
	TableNo string
	Day     time.Time
	Mid     decimal.Decimal
}

//Rate returns the currency exchange rate for a given date from NBP table A
func (n *NBP) Rate(curr Currency, day time.Time) (*Rate, error) {
	resp, err := n.httpClient.Get(fmt.Sprintf("%s/%s/%s", apiBase, curr, day.Format("2006-01-02")))
	if err != nil {
		return nil, fmt.Errorf("can't connect to NBP api: %w", err)
	}
	var rawRate rawRate
	if err := json.NewDecoder(resp.Body).Decode(&rawRate); err != nil {
		return nil, fmt.Errorf("can't decode response: %w", err)
	}
	if len(rawRate.Rates) != 1 {
		return nil, fmt.Errorf("expectation failed: wanted a single rate, instead got %v", rawRate.Rates)
	}
	rate := rawRate.Rates[0]
	effectiveDay, err := time.Parse("2006-01-02", rate.EffectiveDate)
	if err != nil {
		return nil, fmt.Errorf("expectation failed: can't parse date as day %s", rate.EffectiveDate)
	}
	return &Rate{
		TableNo: rate.No,
		Day:     effectiveDay,
		Mid:     rate.Mid,
	}, nil
}
