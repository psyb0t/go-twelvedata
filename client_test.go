package twelvedata_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/psyb0t/go-twelvedata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAPIKey = "EXAMPLE-DO-NOT-USE"
	testSymbol = "SYMBOL_X"
)

func newTestClient(
	t *testing.T,
	handler http.HandlerFunc,
) *twelvedata.ClientWithResponses {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	auth := twelvedata.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "apikey "+testAPIKey)

			return nil
		},
	)

	client, err := twelvedata.NewClientWithResponses(srv.URL, auth)
	require.NoError(t, err)

	return client
}

func TestGetPrice(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/price", r.URL.Path)
		assert.Equal(t, testSymbol, r.URL.Query().Get("symbol"))
		assert.Empty(t, r.URL.Query().Get("apikey"))
		assert.Equal(t, "apikey "+testAPIKey, r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"price":"123.456789"}`))
		assert.NoError(t, err)
	})

	symbol := testSymbol
	res, err := client.GetPriceWithResponse(
		context.Background(),
		&twelvedata.GetPriceParams{Symbol: &symbol},
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	assert.Equal(t, "123.456789", res.JSON200.Price)
}

func TestGetQuote(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/quote", r.URL.Path)
		assert.Equal(t, testSymbol, r.URL.Query().Get("symbol"))
		assert.Equal(t, "apikey "+testAPIKey, r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{
			"symbol":"SYMBOL_X","close":"123.456789",
			"volume":"100","is_market_open":false
		}`))
		assert.NoError(t, err)
	})

	symbol := testSymbol
	res, err := client.GetQuoteWithResponse(
		context.Background(),
		&twelvedata.GetQuoteParams{Symbol: &symbol},
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	assert.Equal(t, testSymbol, res.JSON200.Symbol)
	assert.Equal(t, "123.456789", res.JSON200.Close)
}

func TestUnauthorized(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Empty(t, r.Header.Get("Authorization"))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte(`{
				"code":401,"message":"missing key","status":"error"
			}`))
			assert.NoError(t, err)
		},
	))
	t.Cleanup(srv.Close)

	client, err := twelvedata.NewClientWithResponses(srv.URL)
	require.NoError(t, err)

	res, err := client.GetPriceWithResponse(
		context.Background(),
		&twelvedata.GetPriceParams{},
	)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode())
	assert.Nil(t, res.JSON200)
	require.NotNil(t, res.JSON401)
	assert.Equal(t, "missing key", res.JSON401.Message)
}

func TestQueryEscaping(t *testing.T) {
	t.Parallel()

	const symbol = "SYMBOL_X&apikey=injected value/#"

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/price", r.URL.Path)
		assert.Equal(t, symbol, r.URL.Query().Get("symbol"))
		assert.Empty(t, r.URL.Query().Get("apikey"))
		assert.Equal(t, "apikey "+testAPIKey, r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"price":"1.25"}`))
		assert.NoError(t, err)
	})

	res, err := client.GetPriceWithResponse(
		context.Background(),
		&twelvedata.GetPriceParams{Symbol: new(symbol)},
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)
	assert.Equal(t, "1.25", res.JSON200.Price)
}

func TestCancelledRequest(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("cancelled request reached the server")
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.GetPriceWithResponse(ctx, &twelvedata.GetPriceParams{})
	require.ErrorIs(t, err, context.Canceled)
}

func TestExchangeRatePrecision(t *testing.T) {
	t.Parallel()

	const (
		pair         = "CURRENCY_X/CURRENCY_Y"
		expectedRate = 123.4567890123
	)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/exchange_rate", r.URL.Path)
		assert.Equal(t, pair, r.URL.Query().Get("symbol"))
		assert.Equal(t, "apikey "+testAPIKey, r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{
			"symbol":"CURRENCY_X/CURRENCY_Y",
			"rate":123.4567890123,"timestamp":1
		}`))
		assert.NoError(t, err)
	})

	res, err := client.GetExchangeRateWithResponse(
		context.Background(),
		&twelvedata.GetExchangeRateParams{Symbol: pair},
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)
	assert.Equal(t, expectedRate, res.JSON200.Rate)
}
