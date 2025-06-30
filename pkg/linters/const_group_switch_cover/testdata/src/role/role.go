package role

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
	RoleGuest Role = "guest"
)

func EnoughCase(r Role) {
	switch r {
	case RoleAdmin:
		// Handle admin role
	case RoleUser:
		// Handle user role
	case RoleGuest:
		// Handle guest role
	default:
		// Handle unknown role
		panic("unknown role: " + string(r))
	}
}

func MultipleCase(r Role) {
	switch r {
	case RoleAdmin, RoleUser:
		// Handle admin or user role
	case RoleGuest:
		// Handle guest role
	default:
		// Handle unknown role
		panic("unknown role: " + string(r))
	}
}

func InvalidCase(r Role) {
	switch r { // want "robustruct/linters/switch_case_cover: case value requires type related const value"
	case "admin":
		// Handle admin role
	case "user":
		// Handle user role
	case "guest":
		// Handle guest role
	default:
		// Handle unknown role
		panic("unknown role: " + string(r))
	}
}

func LackedCase(r Role) {
	switch r { // want "robustruct/linters/switch_case_cover: case body uncovered grouped const value"
	case RoleAdmin:
		// Handle admin role
	case RoleUser:
		// Handle user role
	default:
		// Handle unknown role
		panic("unknown role: " + string(r))
	}
}

func FromString(r string) Role {
	switch r {
	case "admin":
		return RoleAdmin
	case "user":
		return RoleUser
	case "guest":
		return RoleGuest
	default:
		panic("unknown role: " + r)
	}
}
