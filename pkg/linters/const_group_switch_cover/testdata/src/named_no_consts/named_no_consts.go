package named_no_consts

type S string

func (s S) String() string {
	return string(s)
}

func LackedSwitch(s S) {
	switch s {
	default:
		// do something
	}
}
