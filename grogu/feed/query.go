package feed

import (
	"math"
	"strconv"
	"strings"
	"time"

	bothanproto "github.com/bandprotocol/bothan-api/go-proxy/proto"
	"golang.org/x/exp/maps"

	grogucontext "github.com/bandprotocol/chain/v2/grogu/context"
	"github.com/bandprotocol/chain/v2/x/feeds/types"
)

func StartQuerySignalIDs(c *grogucontext.Context, l *grogucontext.Logger) {
	for {
		QuerySignalIDs(c, l)
	}
}

func QuerySignalIDs(c *grogucontext.Context, l *grogucontext.Logger) {
	signalIDsWithTimeLimit := <-c.PendingSignalIDs

GetAllSignalIDs:
	for {
		select {
		case nextSignalIDs := <-c.PendingSignalIDs:
			maps.Copy(signalIDsWithTimeLimit, nextSignalIDs)
		default:
			break GetAllSignalIDs
		}
	}

	signalIDs := maps.Keys(signalIDsWithTimeLimit)

	l.Info("Try to get prices for signal ids: %+v", signalIDs)
	prices, err := c.PriceService.Query(signalIDs)
	if err != nil {
		l.Error(":exploding_head: Failed to get prices from price-service with error: %s", c, err.Error())
	}

	maxSafePrice := math.MaxUint64 / uint64(math.Pow10(9))
	now := time.Now()
	submitPrices := []types.SubmitPrice{}
	for _, priceData := range prices {
		switch priceData.PriceOption {
		case bothanproto.PriceOption_PRICE_OPTION_UNSUPPORTED:
			submitPrices = append(submitPrices, types.SubmitPrice{
				PriceOption: types.PriceOptionUnsupported,
				SignalID:    priceData.SignalId,
				Price:       0,
			})
			continue

		case bothanproto.PriceOption_PRICE_OPTION_AVAILABLE:
			price, err := strconv.ParseFloat(strings.TrimSpace(priceData.Price), 64)
			if err != nil || price > float64(maxSafePrice) || price < 0 {
				l.Error(":exploding_head: Failed to parse price from singal id:", c, priceData.SignalId, err)
				priceData.PriceOption = bothanproto.PriceOption_PRICE_OPTION_UNAVAILABLE
				priceData.Price = ""
			} else {
				submitPrices = append(submitPrices, types.SubmitPrice{
					PriceOption: types.PriceOptionAvailable,
					SignalID:    priceData.SignalId,
					Price:       uint64(price * math.Pow10(9)),
				})
				continue
			}
		}

		if signalIDsWithTimeLimit[priceData.SignalId].Before(now) {
			submitPrices = append(submitPrices, types.SubmitPrice{
				PriceOption: types.PriceOptionUnavailable,
				SignalID:    priceData.SignalId,
				Price:       0,
			})
		}
	}

	// delete signal id from in progress map if its price is not found
	signalIDPriceMap := convertToSignalIDPriceMap(submitPrices)
	for _, signalID := range signalIDs {
		if _, found := signalIDPriceMap[signalID]; !found {
			c.InProgressSignalIDs.Delete(signalID)
		}
	}

	if len(submitPrices) == 0 {
		l.Debug(":exploding_head: query signal got no prices with signal ids: %+v", signalIDs)
		return
	}
	l.Info("got prices for signal ids: %+v", maps.Keys(signalIDPriceMap))
	c.PendingPrices <- submitPrices
}

// convertToSignalIDPriceMap converts an array of SubmitPrice to a map of signal id to price.
func convertToSignalIDPriceMap(data []types.SubmitPrice) map[string]uint64 {
	signalIDPriceMap := make(map[string]uint64)

	for _, entry := range data {
		signalIDPriceMap[entry.SignalID] = entry.Price
	}

	return signalIDPriceMap
}
