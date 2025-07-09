package ignore

type Ignore string

const (
	Yes    Ignore = "yes"
	No     Ignore = "no"
	Extend Ignore = "extend"
)

func (i Ignore) IsYes() bool {
	switch i { // ignore:robustruct
	case Yes:
		return true
	}
	return false
}

func (i Ignore) IsNo() bool {
	// ignore:const_group_switch_cover
	switch i {
	case No:
		return true
	}
	return false
}
