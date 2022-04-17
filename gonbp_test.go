package gonbp

import (
	"github.com/shopspring/decimal"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func day(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

type mockResponse struct {
	code int
	body string
	err  error
}

type mockClient struct {
	urls map[string]mockResponse
}

func (m *mockClient) Get(url string) (*http.Response, error) {
	var resp mockResponse
	var ok bool
	if resp, ok = m.urls[url]; !ok {
		panic("response not set up for url: " + url)
	}
	if resp.err != nil {
		return nil, resp.err
	}
	return &http.Response{
		StatusCode: resp.code,
		Body:       io.NopCloser(strings.NewReader(resp.body)),
	}, nil
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
				"https://api.nbp.pl/api/exchangerates/rates/A/EUR/2022-04-15": {
					code: 200,
					body: `{"table":"A","currency":"euro","code":"EUR","rates":[{"no":"074/A/NBP/2022","effectiveDate":"2022-04-15","mid":4.6378}]}`,
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
				"https://api.nbp.pl/api/exchangerates/rates/A/CHF/2021-04-15": {
					code: 200,
					body: `{"table":"A","currency":"frank szwajcarski","code":"CHF","rates":[{"no":"072/A/NBP/2021","effectiveDate":"2021-04-15","mid":4.1198}]}`,
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
				"https://api.nbp.pl/api/exchangerates/rates/A/EUR/2022-04-16": {
					code: 404,
					body: `404 NotFound - Not Found - Brak danych`,
				},
			},
			curr:    EUR,
			day:     day(2022, 4, 16),
			wantErr: true,
		},
		{
			name: "Non existing currency",
			urls: map[string]mockResponse{
				"https://api.nbp.pl/api/exchangerates/rates/A/DOGE/2022-04-15": {
					code: 404,
					body: `404 NotFound`,
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
				httpClient: &mockClient{urls: tt.urls},
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
				"https://api.nbp.pl/api/exchangerates/rates/A/EUR/2022-04-15": {
					code: 200,
					body: `{"table":"A","currency":"euro","code":"EUR","rates":[{"no":"074/A/NBP/2022","effectiveDate":"2022-04-15","mid":4.6378}]}`,
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
				"https://api.nbp.pl/api/exchangerates/rates/A/EUR/2022-04-17": {
					code: 404,
					body: `404 NotFound - Not Found - Brak danych`,
				},
				"https://api.nbp.pl/api/exchangerates/rates/A/EUR/2022-04-16": {
					code: 404,
					body: `404 NotFound - Not Found - Brak danych`,
				},
				"https://api.nbp.pl/api/exchangerates/rates/A/EUR/2022-04-15": {
					code: 200,
					body: `{"table":"A","currency":"euro","code":"EUR","rates":[{"no":"074/A/NBP/2022","effectiveDate":"2022-04-15","mid":4.6378}]}`,
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
				"https://api.nbp.pl/api/exchangerates/rates/A/DOGE/2022-04-15": {
					code: 404,
					body: `404 NotFound`,
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
				httpClient: &mockClient{urls: tt.urls},
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
