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
	"testing"
	"encoding/json"
	"strings"
	"fmt"
	"reflect"
	"net/http"
	"time"
)

var nbpClient = &NbpClient{Client: &http.Client{Timeout: time.Second * 10}}


func TestCurrencyRateJsonDecoding(t *testing.T) {
	const rateBody = `{"no":"194/A/NBP/2016","effectiveDate":"2016-10-06","mid":4.2974}`
	expected := CurrencyRate{Number:"194/A/NBP/2016", EffectiveDate:"2016-10-06", Mid: "4.2974"}
	var rate CurrencyRate

	if err := json.NewDecoder(strings.NewReader(rateBody)).Decode(&rate); err != nil {
		t.Fatal(err)
	}
	if expected != rate {
		t.Error(fmt.Sprintf("Expected %s, got %s", expected, rate))
	}
}

func TestCurrencyRateListJsonDecoding(t *testing.T) {
	const body = `{"table":"A","currency":"euro","code":"EUR","rates":[{"no":"194/A/NBP/2016","effectiveDate":"2016-10-06","mid":4.2974}]}`
	expected := CurrencyRateList{
		Table: "A",
		Currency: "euro",
		Code: "EUR",
		Rates: []CurrencyRate{{Number:"194/A/NBP/2016", EffectiveDate:"2016-10-06", Mid: "4.2974"}},
	}
	var rates CurrencyRateList

	if err := json.NewDecoder(strings.NewReader(body)).Decode(&rates); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, rates) {
		t.Error(fmt.Sprintf("Expected %s, got %s", expected, rates))
	}
}

func TestExchangeRateForGivenDay(t *testing.T) {
	expected := CurrencyRateList{
		Table: "A",
		Currency: "euro",
		Code: "EUR",
		Rates: []CurrencyRate{{Number:"194/A/NBP/2016", EffectiveDate:"2016-10-06", Mid: "4.2974"}},
	}
	rates, err := nbpClient.Day("A", "EUR", "2016-10-06")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, *rates) {
		t.Error(fmt.Sprintf("Expected %s, got %s", expected, rates))
	}
}
