package priceservice

import (
	"encoding/json"
	"net/url"
	"path"
	"strings"
	"time"

	bothanproto "github.com/bandprotocol/bothan-api/go-proxy/proto"
	"github.com/levigross/grequests"
)

type RestService struct {
	url     string
	timeout time.Duration
}

func NewRestService(url string, timeout time.Duration) *RestService {
	return &RestService{url: url, timeout: timeout}
}

type PriceData struct {
	Prices map[string]float64 `json:"prices"`
}

type QueryStruct struct {
	SignalIds []string `json:"signal_ids"`
}

func (rs *RestService) Query(signalIds []string) ([]*bothanproto.PriceData, error) {
	concatSignalIds := strings.Join(signalIds, ",")
	u, err := url.Parse(rs.url)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, concatSignalIds)

	resp, err := grequests.Get(
		u.String(),
		&grequests.RequestOptions{
			RequestTimeout: rs.timeout,
		},
	)
	if err != nil {
		return nil, err
	}

	var priceResp bothanproto.QueryPricesResponse
	err = json.Unmarshal(resp.Bytes(), &priceResp)
	if err != nil {
		return nil, err
	}

	return priceResp.Prices, nil
}
