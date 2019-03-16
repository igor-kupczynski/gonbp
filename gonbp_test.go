package gonbp

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestCurrencyRateJsonDecoding(t *testing.T) {
	const rateBody = `{"no":"194/A/NBP/2016","effectiveDate":"2016-10-06","mid":4.2974}`
	expected := CurrencyRate{Number: "194/A/NBP/2016", EffectiveDate: "2016-10-06", Mid: "4.2974"}
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
		Table:    "A",
		Currency: "euro",
		Code:     "EUR",
		Rates:    []CurrencyRate{{Number: "194/A/NBP/2016", EffectiveDate: "2016-10-06", Mid: "4.2974"}},
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
		Table:    "A",
		Currency: "euro",
		Code:     "EUR",
		Rates:    []CurrencyRate{{Number: "194/A/NBP/2016", EffectiveDate: "2016-10-06", Mid: "4.2974"}},
	}
	rates, err := DefaultNbpClient.Day("A", "EUR", "2016-10-06")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, *rates) {
		t.Error(fmt.Sprintf("Expected %s, got %s", expected, rates))
	}
}

func TestApiException(t *testing.T) {
	_, err := DefaultNbpClient.Day("A", "EUR", "2100-10-06")
	if err == nil {
		t.Error("NbpAPIError expected")
	}
	if err.(NbpAPIError).StatusCode != 400 {
		t.Errorf("StatusCode == 400 expected, but got %d", err.(NbpAPIError).StatusCode)
	}
	if !strings.Contains(err.Error(), "Invalid date range") {
		t.Errorf("'Invalid date range' error  expected, but got %s", err.Error())
	}
}
