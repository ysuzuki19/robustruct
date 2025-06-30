package external

import "external/sub"

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

func HardCoded(xyz sub.XYZ) string {
	switch xyz { // want "robustruct/linters/switch_case_cover: case value requires type related const value"
	case sub.X:
		return "X"
	case 1: // hard-coded value
		return "1"
	case sub.Z:
		return "Z"
	default:
		return "unknown"
	}
}
