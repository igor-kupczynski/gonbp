/*
Copyright 2016 Igor Kupczynski

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package gonbp

import (
	"net/http"
	"fmt"
	"encoding/json"
)

const apiRoot = "http://api.nbp.pl/api"
const rates = apiRoot + "/exchangerates/rates"

// Currencies

// Currency represents a foreign currency together with NBP table that stores it
type Currency struct {
	// Currency exchange table to use
	// "A" - table of average foreign currency exchange rates;
	// "B" - table of average exchange rates of inconvertible currencies;
	// "C" - table of purchase and sale exchange rates;
	Table string

	// Currency code
	Code  string
}


// Predefined currencies
var (
	Eur = &Currency{"A", "EUR"}
)

// Average currency rate published a day
type CurrencyRate struct {
	// NBP exchange rate publication number
	Number        string      `json:"no"`

	// Rate publication date YYYY-MM-DD
	EffectiveDate string      `json:"effectiveDate"`

	// Average exchange rate stored as string to retain original precision
	Mid           json.Number `json:"mid,Number"`
}


// List of rates spanning for multiple days
type CurrencyRateList struct {
	//NBP currency table
	Table        string         `json:"table"`

	// Currency name
	CurrencyName string         `json:"currency"`

	// Currency code
	CurrencyCode string         `json:"code"`

	// List of rates
	Rates        []CurrencyRate `json:"rates"`
}


// CurrencyRateClient is the NBP single currency exchange rate API client
type CurrencyRateClient struct {
	// Currency for client
	Currency *Currency

	// Http client to use to fetch the rate
	client   *http.Client
}

// Create new CurrencyRateClient with default http client
func NewCurrencyRateClient(curr *Currency) *CurrencyRateClient {
	return &CurrencyRateClient{Currency: curr, client: http.DefaultClient}
}


// Current exchange rate for a currency
func (cc *CurrencyRateClient) Current() (*CurrencyRateList, error) {
	url := fmt.Sprintf("%s/%s/%s", rates, cc.Currency.Table, cc.Currency.Code)
	response, err := cc.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	rates := &CurrencyRateList{}
	err = json.NewDecoder(response.Body).Decode(rates)
	if err != nil {
		return nil, err
	}

	return rates, nil
}
