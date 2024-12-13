package account

/*
import "fmt"

type account int

const (
	zero  account = 0 // receive default value
	admin account = iota
	user  account = iota
	guest account = iota
)

type Account struct {
	tag account
}

var _ fmt.Stringer = Account{} // ignore:fields_require

func Zero() Account {
	return Account{} //ignore:fields_require
}

func Admin() Account {
	return Account{tag: admin}
}

func User() Account {
	return Account{tag: user}
}

func Guest() Account {
	return Account{tag: guest}
}

func (a Account) IsAdmin() bool {
	return a.tag == admin
}

func (a Account) IsUser() bool {
	return a.tag == user
}

func (a Account) IsGuest() bool {
	return a.tag == guest
}

func FromString(s string) (Account, error) {
	switch s {
	case "admin":
		return Admin(), nil
	case "user":
		return User(), nil
	case "guest":
		return Guest(), nil
	default:
		return Zero(), fmt.Errorf("invalid account: %s", s)
	}
}

func (a Account) String() string {
	switch a.tag {
	case admin:
		return "admin"
	case user:
		return "user"
	case guest:
		return "guest"
	case zero:
		return "guest"
	default:
		panic("unreachable")
	}
}

type Switcher struct {
	Admin func()
	User  func()
	Guest func()
}

func (a Account) Switch(s Switcher) {
	switch a.tag {
	case admin:
		s.Admin()
	case user:
		s.User()
	case guest:
		s.Guest()
	case zero:
		s.Guest()
	}
}
*/
