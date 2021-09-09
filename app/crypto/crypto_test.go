package crypto

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/tradersclub/TCStocksCrypto/mocks"
	"github.com/tradersclub/TCStocksCrypto/model"
)

func TestGetMarkets(t *testing.T) {
	startedAt := time.Now()
	inputFloat := float64(15.5)
	expectedResult := []model.Market{{ID: "bitcoin", Symbol: "btc", Name: "Bitcoin", Image: "https://assets.coingecko.com/coins/images/1/large/bitcoin.png?1547033579", CurrentPrice: 47201, MarketCap: 887534283153, MarketCapRank: 1, FullyDilutedValuation: &inputFloat, TotalVolume: 34612719880, High24H: 48713, Low24H: 46811, PriceChange24H: -1038.574618445593, PriceChangePercentage24H: -2.15294, MarketCapChange24H: -18456933846.516357, MarketCapChangePercentage24H: -2.03721, CirculatingSupply: 18802743, TotalSupply: &inputFloat, MaxSupply: &inputFloat, Ath: 64805, AthChangePercentage: -27.16403, AthDate: time.Time{}, Atl: 67.81, AtlChangePercentage: 69508.97064, AtlDate: time.Now(), Roi: nil, LastUpdated: time.Now()}, {ID: "ethereum", Symbol: "eth", Name: "Ethereum", Image: "https://assets.coingecko.com/coins/images/279/large/ethereum.png?1595348880", CurrentPrice: 3370.68, MarketCap: 395391546668, MarketCapRank: 2, FullyDilutedValuation: nil, TotalVolume: 31353762178, High24H: 3457.63, Low24H: 3203.43, PriceChange24H: 72.08, PriceChangePercentage24H: 2.18524, MarketCapChange24H: 11571839246, MarketCapChangePercentage24H: 3.01492, CirculatingSupply: 117334257.999, TotalSupply: nil, MaxSupply: nil, Ath: 4356.99, AthChangePercentage: -22.65786, AthDate: time.Now(), Atl: 0.432979, AtlChangePercentage: 778180.44378, AtlDate: time.Now(), Roi: nil, LastUpdated: time.Now()}}
	cases := map[string]struct {
		ExpectedErr  error
		ExpectedData []model.Market

		InputVersion  string
		InputDatetime time.Time

		PrepareMock func(mock *mocks.MockApp)
	}{
		"deve retornar sucesso": {
			ExpectedErr:   nil,
			ExpectedData:  expectedResult,
			InputVersion:  "1",
			InputDatetime: startedAt,
			PrepareMock: func(mock *mocks.MockApp) {
				mock.EXPECT().GetCryptoMarkets(gomock.Any()).Return(expectedResult, nil).Times(1)
			}},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			mockCrypto := mocks.NewMockApp(ctrl)
			cs.PrepareMock(mockCrypto)
			result, err := mockCrypto.GetCryptoMarkets(ctx)

			if diff := cmp.Diff(result, cs.ExpectedData); diff != "" {
				t.Error(diff)
			}

			if diff := cmp.Diff(err, cs.ExpectedErr); diff != "" {
				t.Error(diff)
			}
		})
	}
}
