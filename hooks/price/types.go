package price

import (
	"github.com/bandprotocol/chain/v3/pkg/obi"
)

const DefaultMultiplier = uint64(1000000000)

type CommonOutput struct {
	Symbols    []string
	Rates      []uint64
	Multiplier uint64
}

type LegacyInput struct {
	Symbols    []string `json:"symbols"`
	Multiplier uint64   `json:"multiplier"`
}

type LegacyOutput struct {
	Rates []uint64 `json:"rates"`
}

type Input struct {
	Symbols            []string
	MinimumSourceCount uint8
}

type Output struct {
	Responses []Response
}

type Response struct {
	Symbol       string
	ResponseCode uint8
	Rate         uint64
}

func MustDecodeResult(calldata, result []byte) CommonOutput {
	commonOutput, err := DecodeResult(result)
	if err == nil {
		return commonOutput
	}

	commonOutput, err = DecodeLegacyResult(calldata, result)
	if err != nil {
		panic(err)
	}

	return commonOutput
}

func DecodeLegacyResult(calldata, result []byte) (CommonOutput, error) {
	var legacyInput LegacyInput
	var legacyOutput LegacyOutput

	err := obi.Decode(calldata, &legacyInput)
	if err != nil {
		return CommonOutput{}, err
	}

	err = obi.Decode(result, &legacyOutput)
	if err != nil {
		return CommonOutput{}, err
	}

	return CommonOutput{
		Symbols:    legacyInput.Symbols,
		Rates:      legacyOutput.Rates,
		Multiplier: DefaultMultiplier,
	}, nil
}

func DecodeResult(result []byte) (CommonOutput, error) {
	var out Output
	var symbols []string
	var rates []uint64

	err := obi.Decode(result, &out)
	if err != nil {
		return CommonOutput{}, err
	}

	for _, r := range out.Responses {
		if r.ResponseCode != 0 {
			continue
		}

		symbols = append(symbols, r.Symbol)
		rates = append(rates, r.Rate)
	}

	return CommonOutput{
		Symbols:    symbols,
		Rates:      rates,
		Multiplier: DefaultMultiplier,
	}, nil
}
