package twelvedata_test

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-twelvedata"
)

//nolint:testableexamples // requires provider credentials; compile check only
func Example() {
	auth := twelvedata.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set(
				"Authorization",
				"apikey "+os.Getenv("TWELVEDATA_API_KEY"),
			)

			return nil
		},
	)

	client, err := twelvedata.NewClientWithResponses(
		"https://api.twelvedata.com",
		auth,
	)
	if err != nil {
		panic(err)
	}

	const symbol = "SYMBOL_X"

	res, err := client.GetPriceWithResponse(
		context.Background(),
		&twelvedata.GetPriceParams{Symbol: new(symbol)},
	)
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
