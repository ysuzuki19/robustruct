package settings

type Flag string

const (
	FlagDisableTest Flag = "disable_test"
)

func (f Flag) String() string { return string(f) }

type Flags []Flag

func (fs Flags) Contains(f string) bool {
	for _, itr := range fs {
		if itr.String() == f {
			return true
		}
	}
	return false
}

func (fs Flags) DisableTest() bool {
	return fs.Contains(FlagDisableTest.String())
}
