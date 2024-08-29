package types

import (
	"gopkg.in/yaml.v2"
)

const (
	// Default values for Params
	DefaultAllowableBlockTimeDiscrepancy = int64(60)
	DefaultGracePeriod                   = int64(30)
	DefaultMinInterval                   = int64(60)
	DefaultMaxInterval                   = int64(3600)
	DefaultPowerStepThreshold            = int64(1_000_000_000)
	DefaultMaxCurrentFeeds               = uint64(300)
	DefaultCooldownTime                  = int64(30)
	DefaultMinDeviationBasisPoint        = int64(50)
	DefaultMaxDeviationBasisPoint        = int64(3000)
	// estimated from block time of 3 seconds, aims for 1 day update
	DefaultCurrentFeedsUpdateInterval = int64(28800)
	DefaultMaxSignalIDsPerSigning     = uint64(10)
)

// NewParams creates a new Params instance
func NewParams(
	admin string,
	allowableBlockTimeDiscrepancy int64,
	gracePeriod int64,
	minInterval int64,
	maxInterval int64,
	powerStepThreshold int64,
	maxCurrentFeeds uint64,
	cooldownTime int64,
	minDeviationBasisPoint int64,
	maxDeviationBasisPoint int64,
	currentFeedsUpdateInterval int64,
	maxSignalIDsPerSigning uint64,
) Params {
	return Params{
		Admin:                         admin,
		AllowableBlockTimeDiscrepancy: allowableBlockTimeDiscrepancy,
		GracePeriod:                   gracePeriod,
		MinInterval:                   minInterval,
		MaxInterval:                   maxInterval,
		PowerStepThreshold:            powerStepThreshold,
		MaxCurrentFeeds:               maxCurrentFeeds,
		CooldownTime:                  cooldownTime,
		MinDeviationBasisPoint:        minDeviationBasisPoint,
		MaxDeviationBasisPoint:        maxDeviationBasisPoint,
		CurrentFeedsUpdateInterval:    currentFeedsUpdateInterval,
		MaxSignalIDsPerSigning:        maxSignalIDsPerSigning,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		"[NOT_SET]",
		DefaultAllowableBlockTimeDiscrepancy,
		DefaultGracePeriod,
		DefaultMinInterval,
		DefaultMaxInterval,
		DefaultPowerStepThreshold,
		DefaultMaxCurrentFeeds,
		DefaultCooldownTime,
		DefaultMinDeviationBasisPoint,
		DefaultMaxDeviationBasisPoint,
		DefaultCurrentFeedsUpdateInterval,
		DefaultMaxSignalIDsPerSigning,
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	fields := []struct {
		validateFn     func(string, bool, interface{}) error
		name           string
		val            interface{}
		isPositiveOnly bool
	}{
		{validateString, "admin", p.Admin, false},
		{validateInt64, "allowable block time discrepancy", p.AllowableBlockTimeDiscrepancy, true},
		{validateInt64, "grace period", p.GracePeriod, true},
		{validateInt64, "min interval", p.MinInterval, true},
		{validateInt64, "max interval", p.MaxInterval, true},
		{validateInt64, "power threshold", p.PowerStepThreshold, true},
		{validateUint64, "max current feeds", p.MaxCurrentFeeds, false},
		{validateInt64, "cooldown time", p.CooldownTime, true},
		{validateInt64, "min deviation basis point", p.MinDeviationBasisPoint, true},
		{validateInt64, "max deviation basis point", p.MaxDeviationBasisPoint, true},
		{validateInt64, "current feeds update interval", p.CurrentFeedsUpdateInterval, true},
		{validateUint64, "max signalIDs per Signing", p.MaxSignalIDsPerSigning, true},
	}

	for _, f := range fields {
		if err := f.validateFn(f.name, f.isPositiveOnly, f.val); err != nil {
			return err
		}
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
