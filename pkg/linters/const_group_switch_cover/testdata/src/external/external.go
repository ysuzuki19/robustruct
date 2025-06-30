package external

import "external/sub"

type Dummy int

const (
	A Dummy = iota
	B
	C
)

func EnoughSwitch(xyz sub.XYZ) string {
	switch xyz {
	case sub.X:
		return "X"
	case sub.Y:
		return "Y"
	case sub.Z:
		return "Z"
	default:
		return "unknown"
	}
}

func LackedSwitch(xyz sub.XYZ) string {
	switch xyz { // want "robustruct/linters/switch_case_cover: case body uncovered grouped const value"
	case sub.X:
		return "X"
	case sub.Y:
		return "Y"
	default:
		return "unknown"
	}
}
