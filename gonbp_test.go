package gonbp

import (
	"fmt"
	"github.com/shopspring/decimal"
	"gonbp/internal/nbpapi"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func day(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

type mockResponse struct {
	rates *nbpapi.Rates
	err   error
}

type mockClient struct {
	urls map[string]mockResponse
}

func (m *mockClient) Get(curr string, day time.Time) (*nbpapi.Rates, error) {
	url := fmt.Sprintf("%s/%s", curr, day.Format("2006-01-02"))
	var resp mockResponse
	var ok bool
	if resp, ok = m.urls[url]; !ok {
		panic("response not set up for " + url)
	}
	if resp.err != nil {
		return nil, resp.err
	}
	return resp.rates, nil
}

func TestNBP_Rate(t *testing.T) {
	tests := []struct {
		name    string
		urls    map[string]mockResponse
		curr    Currency
		day     time.Time
		want    *Rate
		wantErr bool
	}{
		{
			name: "Positive case EUR",
			urls: map[string]mockResponse{
				"EUR/2022-04-15": {
					rates: &nbpapi.Rates{
						Table:    "A",
						Currency: "euro",
						Code:     "EUR",
						Rates: []nbpapi.DailyRate{
							{
								No:            "074/A/NBP/2022",
								EffectiveDate: "2022-04-15",
								Mid:           decimal.NewFromFloat(4.6378),
							},
						},
					},
				},
			},
			curr: EUR,
			day:  day(2022, 4, 15),
			want: &Rate{
				TableNo: "074/A/NBP/2022",
				Day:     day(2022, 4, 15),
				Mid:     decimal.NewFromFloat(4.6378),
			},
			wantErr: false,
		},
		{
			name: "Positive case CHF",
			urls: map[string]mockResponse{
				"CHF/2021-04-15": {
					rates: &nbpapi.Rates{
						Table:    "A",
						Currency: "frank szwajcarski",
						Code:     "CHF",
						Rates: []nbpapi.DailyRate{
							{
								No:            "072/A/NBP/2021",
								EffectiveDate: "2021-04-15",
								Mid:           decimal.NewFromFloat(4.1198),
							},
						},
					},
				},
			},
			curr: CHF,
			day:  day(2021, 4, 15),
			want: &Rate{
				TableNo: "072/A/NBP/2021",
				Day:     day(2021, 4, 15),
				Mid:     decimal.NewFromFloat(4.1198),
			},
			wantErr: false,
		},
		{
			name: "Not found for a given day",
			urls: map[string]mockResponse{
				"EUR/2022-04-16": {
					err: nbpapi.ErrNoExchangeRateForGivenDay,
				},
			},
			curr:    EUR,
			day:     day(2022, 4, 16),
			wantErr: true,
		},
		{
			name: "Non existing currency",
			urls: map[string]mockResponse{
				"DOGE/2022-04-15": {
					err: nbpapi.ErrNoExchangeRateForGivenDay,
				},
			},
			curr:    "DOGE",
			day:     day(2022, 4, 15),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			n := &NBP{
				api: &mockClient{urls: tt.urls},
			}
			got, err := n.Rate(tt.curr, tt.day)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Rate() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNBP_PreviousRate(t *testing.T) {
	tests := []struct {
		name    string
		urls    map[string]mockResponse
		curr    Currency
		day     time.Time
		want    *Rate
		wantErr bool
	}{
		{
			name: "Go back a day",
			urls: map[string]mockResponse{
				"EUR/2022-04-15": {
					rates: &nbpapi.Rates{
						Table:    "A",
						Currency: "euro",
						Code:     "EUR",
						Rates: []nbpapi.DailyRate{
							{
								No:            "074/A/NBP/2022",
								EffectiveDate: "2022-04-15",
								Mid:           decimal.NewFromFloat(4.6378),
							},
						},
					},
				},
			},
			curr: EUR,
			day:  day(2022, 4, 16),
			want: &Rate{
				TableNo: "074/A/NBP/2022",
				Day:     day(2022, 4, 15),
				Mid:     decimal.NewFromFloat(4.6378),
			},
			wantErr: false,
		},
		{
			name: "Go back over a long weekend",
			urls: map[string]mockResponse{
				"EUR/2022-04-17": {
					err: nbpapi.ErrNoExchangeRateForGivenDay,
				},
				"EUR/2022-04-16": {
					err: nbpapi.ErrNoExchangeRateForGivenDay,
				},
				"EUR/2022-04-15": {
					rates: &nbpapi.Rates{
						Table:    "A",
						Currency: "euro",
						Code:     "EUR",
						Rates: []nbpapi.DailyRate{
							{
								No:            "074/A/NBP/2022",
								EffectiveDate: "2022-04-15",
								Mid:           decimal.NewFromFloat(4.6378),
							},
						},
					},
				},
			},
			curr: EUR,
			day:  day(2022, 4, 18),
			want: &Rate{
				TableNo: "074/A/NBP/2022",
				Day:     day(2022, 4, 15),
				Mid:     decimal.NewFromFloat(4.6378),
			},
			wantErr: false,
		},
		{
			name: "Non-existing currency",
			urls: map[string]mockResponse{
				"DOGE/2022-04-15": {
					err: nbpapi.ErrNoRatesForCurrency,
				},
			},
			curr:    "DOGE",
			day:     day(2022, 4, 16),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			n := &NBP{
				api: &mockClient{urls: tt.urls},
			}
			got, err := n.PreviousRate(tt.curr, tt.day)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Rate() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
