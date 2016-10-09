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


// NBP API Error
type NbpApiError struct {
	StatusCode int
	Message    string
}

func (c NbpApiError) Error() string {
	return c.Message
}


// Average currency rate published a day
type CurrencyRate struct {
	Number        string      `json:"no"`
	EffectiveDate string      `json:"effectiveDate"`
	Mid           json.Number `json:"mid,Number"`
}


// List of rates spanning for multiple days
type CurrencyRateList struct {
	Table    string         `json:"table"`
	Currency string         `json:"currency"`
	Code     string         `json:"code"`
	Rates    []CurrencyRate `json:"rates"`
}


// NbpClient is client to connect to the NBP exchange rate api
type NbpClient struct {
	// Http client to use to fetch the rate
	Client   *http.Client
}

// NbpClient with default paramters
var DefaultNbpClient = &NbpClient{Client:http.DefaultClient}


// Current exchange rate for a currency
func (c *NbpClient) Current(table string, currCode string) (*CurrencyRateList, error) {
	url := fmt.Sprintf("%s/%s/%s", rates, table, currCode)
	return c.fetchRates(url)
}

// Current exchange rate for a for given day
func (c *NbpClient) Day(table string, currCode string, day string) (*CurrencyRateList, error) {
	url := fmt.Sprintf("%s/%s/%s/%s", rates, table, currCode, day)
	return c.fetchRates(url)
}

func (c *NbpClient) fetchRates(url string)  (*CurrencyRateList, error) {
	response, err := c.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 && response.StatusCode < 500 {
		return nil, NbpApiError{response.StatusCode, response.Status}
	}

	rates := &CurrencyRateList{}
	err = json.NewDecoder(response.Body).Decode(rates)
	if err != nil {
		return nil, err
	}

	return rates, nil
}
