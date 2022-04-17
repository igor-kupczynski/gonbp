package nbpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
	"time"
)

type httpClient interface {
	Get(url string) (resp *http.Response, err error)
}

// Client is a low-level client over the NBP rates API
type Client struct {
	http httpClient
}

// Init returns *Rates instance with net/http.DefaultClient
func Init(client httpClient) *Client {
	return &Client{
		http: client,
	}
}

// Rates represents the return value of the NBP rates API
type Rates struct {
	Table    string      `json:"table"`
	Currency string      `json:"currency"`
	Code     string      `json:"code"`
	Rates    []DailyRate `json:"rates"`
}

// DailyRate represent a rate for a single day
type DailyRate struct {
	No            string          `json:"no"`
	EffectiveDate string          `json:"effectiveDate"`
	Mid           decimal.Decimal `json:"mid"`
}

// ErrNoExchangeRateForGivenDay represents a failure where there are no published rates for a given day
var ErrNoExchangeRateForGivenDay = errors.New("no exchange rate for given date")

// ErrNoRatesForCurrency represents a failure where NBP doesn't publish exchange rates for the given currency
var ErrNoRatesForCurrency = errors.New("currency without published rates")

// ErrApiCallUnsuccessful represents a generic unsuccessful API call
type ErrApiCallUnsuccessful struct {
	Code int
	Body string
}

func (e ErrApiCallUnsuccessful) Error() string {
	return fmt.Sprintf("unsuccssesful API call, code %d, body %s", e.Code, e.Body)
}

const (
	apiBase = "https://api.nbp.pl/api/exchangerates/rates/A"
)

// Get returns the currency exchange rate for a given date from NBP table A
func (c *Client) Get(curr string, day time.Time) (*Rates, error) {
	resp, err := c.http.Get(fmt.Sprintf("%s/%s/%s", apiBase, curr, day.Format("2006-01-02")))
	if err != nil {
		return nil, fmt.Errorf("can't connect to NBP api: %w", err)
	}

	if resp.StatusCode == 404 {
		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("can't read response: %w", err)
		}
		if bytes.Contains(buf, []byte("Brak danych")) {
			return nil, ErrNoExchangeRateForGivenDay
		}
		return nil, ErrNoRatesForCurrency
	}

	if resp.StatusCode != 200 {
		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("can't read response: %w", err)
		}
		return nil, ErrApiCallUnsuccessful{Code: resp.StatusCode, Body: string(buf)}
	}

	var rates Rates
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		return nil, fmt.Errorf("can't decode response: %w", err)
	}
	return &rates, nil
}
