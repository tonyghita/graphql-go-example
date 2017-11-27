package resolver

import (
	"errors"
	"strings"
)

// A LengthUnit is a measure of the physical property length.
type LengthUnit float64

// LengthUnitArgs ...
type LengthUnitArgs struct {
	Unit string
}

const (
	Millimeter            = Meter * 1e-2
	Centimeter            = Meter * 1e-1
	Meter      LengthUnit = 1e0
	Kilometer             = Meter * 1e3
	Inch                  = Meter * 0.0254
	Foot                  = Inch * 12
	Yard                  = Foot * 3
	Mile                  = Meter * 1609.344
)

var _strToLengthUnit = map[string]LengthUnit{
	"MILLIMETER": Millimeter,
	"CENTIMETER": Centimeter,
	"METER":      Meter,
	"KILOMETER":  Kilometer,
	"INCH":       Inch,
	"FOOT":       Foot,
	"YARD":       Yard,
	"MILE":       Mile,
}

var _LengthUnitToStr = map[LengthUnit]string{
	Millimeter: "MILLIMETER",
	Centimeter: "CENTIMETER",
	Meter:      "METER",
	Kilometer:  "KILOMETER",
	Inch:       "INCH",
	Foot:       "FOOT",
	Yard:       "YARD",
	Mile:       "MILE",
}

// ToLengthUnit converts a string value to a LengthUnit enumeration value.
func ToLengthUnit(s string) (LengthUnit, error) {
	unit, ok := _strToLengthUnit[strings.ToUpper(s)]
	if !ok {
		return LengthUnit(0), errors.New("unknown LengthUnit: " + s)
	}
	return unit, nil
}

// String returns the string representation of the LengthUnit enumeration.
func (u LengthUnit) String() string {
	return _LengthUnitToStr[u]
}

// Convert a float value to the length unit.
func ConvertLength(v float64, from LengthUnit, to LengthUnit) float64 {
	return (v * float64(from)) / float64(to)
}
