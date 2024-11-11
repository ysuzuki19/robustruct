package settings

type Settings struct {
	Features Features `json:"features"`
	Flags    Flags    `json:"flags"`
}
