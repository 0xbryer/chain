package keeper

import (
	"github.com/bandprotocol/chain/v2/x/feeds/types"
)

// calculateIntervalAndDeviation calculates feed interval and deviation from power
func CalculateIntervalAndDeviation(power int64, param types.Params) (interval int64, deviation int64) {
	if power < param.PowerThreshold {
		return 0, 0
	}

	// divide power by power threshold to create steps
	powerFactor := power / param.PowerThreshold

	interval = param.MaxInterval / powerFactor
	if interval < param.MinInterval {
		interval = param.MinInterval
	}

	deviation = param.MaxDeviationInThousandth / powerFactor
	if deviation < param.MinDeviationInThousandth {
		deviation = param.MinDeviationInThousandth
	}

	return
}

// sumPower sums power from a list of signals
func sumPower(signals []types.Signal) (sum int64) {
	for _, signal := range signals {
		sum += signal.Power
	}
	return
}
