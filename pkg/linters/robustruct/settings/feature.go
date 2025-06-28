package settings

type Feature string

const (
	FeatureFieldsRequire         Feature = "fields_require"
	FeatureFieldsAlign           Feature = "fields_align"
	FeatureConstGroupSwitchCover Feature = "const_group_switch_cover"
)

func (f Feature) String() string { return string(f) }

type Features []Feature

func (fs Features) Contains(f string) bool {
	for _, itr := range fs {
		if itr.String() == f {
			return true
		}
	}
	return false
}
