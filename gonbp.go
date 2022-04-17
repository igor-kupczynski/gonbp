// Package gonbp is a go wrapper over NBP API
package gonbp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	return WithClient(http.DefaultClient)
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

// ErrNoExchangeRateForGivenDay represents a failure where there are no published rates for a given day
type ErrNoExchangeRateForGivenDay struct {
	Day time.Time
}

func (e ErrNoExchangeRateForGivenDay) Error() string {
	return fmt.Sprintf("no exchange rate for given date %s", e.Day.Format("2006-01-02"))
}

// ErrNoRatesForCurrency represents a failure where NBP doesn't publish exchange rates for the given currency
type ErrNoRatesForCurrency struct {
	Curr Currency
}

func (e ErrNoRatesForCurrency) Error() string {
	return fmt.Sprintf("currency without published rates %s", e.Curr)
}

// ErrApiCallUnsuccessful represents a generic unsuccessful API call
type ErrApiCallUnsuccessful struct {
	Code int
	Body string
}

func (e ErrApiCallUnsuccessful) Error() string {
	return fmt.Sprintf("unsuccssesful API call, code %d, body %s", e.Code, e.Body)
}

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

	if resp.StatusCode == 404 {
		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("can't read response: %w", err)
		}
		if bytes.Contains(buf, []byte("Brak danych")) {
			return nil, ErrNoExchangeRateForGivenDay{Day: day}
		}
		return nil, ErrNoRatesForCurrency{Curr: curr}
	}

	if resp.StatusCode != 200 {
		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("can't read response: %w", err)
		}
		return nil, ErrApiCallUnsuccessful{Code: resp.StatusCode, Body: string(buf)}
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
