package resolver

import (
	"errors"
	"strings"
)

// A MassUnit is a measure of a physical body.
type MassUnit float64

type MassUnitArgs struct {
	Unit string
}

const (
	Gram               = Kilogram * 1e-3
	Kilogram  MassUnit = 1e0
	MetricTon          = Gram * 1e6
	Pound              = Kilogram * 0.45359237
)

// ToMassUnit converts a string value to a MassUnit enumeration value.
func ToMassUnit(s string) (MassUnit, error) {
	u, ok := _strToMassUnit[strings.ToUpper(s)]
	if !ok {
		return MassUnit(0), errors.New("unknown MassUnit: " + s)
	}
	return u, nil
}

// String returns the string representation of the MassUnit enumeration.
func (u MassUnit) String() string {
	return _MassUnitToStr[u]
}

var _MassUnitToStr = map[MassUnit]string{
	Gram:      "GRAM",
	Kilogram:  "KILOGRAM",
	MetricTon: "METRIC_TON",
	Pound:     "POUND",
}

var _strToMassUnit = map[string]MassUnit{
	"GRAM":       Gram,
	"KILOGRAM":   Kilogram,
	"METRIC_TON": MetricTon,
	"POUND":      Pound,
}

func ConvertMass(v float64, from, to MassUnit) float64 {
	return (v * float64(from)) / float64(to)
}
