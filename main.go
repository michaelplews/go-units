package units

import (
	"fmt"
	"strings"
)

// Return all registered Units
func All() (units []Unit) {
	for _, u := range unitMap {
		units = append(units, u)
	}
	return units
}

// MustConvertFloat converts a provided float from one Unit to another, panicking on error
func MustConvertFloat(x float64, from, to Unit) Value {
	val, err := ConvertFloat(x, from, to)
	if err != nil {
		panic(err)
	}
	return val
}

// ConvertFloat converts a provided float from one Unit to another
func ConvertFloat(x float64, from, to Unit) (Value, error) {
	path, err := ResolveConversion(from, to)
	if err != nil {
		return Value{}, err
	}

	for _, c := range path {
		x = c.Fn(x)
	}

	return Value{x, to}, nil
}

// Find Unit matching name or symbol provided
func Find(s string) (Unit, error) {

	// first try case-sensitive match
	for _, u := range unitMap {
		if matchUnit(s, u, true) {
			return u, nil
		}
	}

	// then case-insensitive
	for _, u := range unitMap {
		if matchUnit(s, u, false) {
			return u, nil
		}
	}

	// finally, try stripping plural suffix
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "S") {
		s = strings.TrimSuffix(s, "s")
		s = strings.TrimSuffix(s, "S")
		for _, u := range unitMap {
			if matchUnit(s, u, false) {
				return u, nil
			}
		}
	}

	return Unit{}, fmt.Errorf("unit \"%s\" not found", s)
}

func matchUnit(s string, u Unit, matchCase bool) bool {
	for _, name := range u.Names() {
		if matchCase {
			if name == s {
				return true
			}
		} else {
			if strings.EqualFold(s, name) {
				return true
			}
		}
	}

	return false
}
