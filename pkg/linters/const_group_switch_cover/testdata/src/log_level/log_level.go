package loglevel

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

func EnoughCase(ll LogLevel) string {
	switch ll {
	case LogLevelDebug:
		return "debug"
	case LogLevelInfo:
		return "info"
	case LogLevelWarn:
		return "warn"
	case LogLevelError:
		return "error"
	default:
		return "unknown"
	}
}

func MultipleCase(ll LogLevel) string {
	switch ll {
	case LogLevelDebug, LogLevelInfo:
		return "debug"
	case LogLevelWarn:
		return "warn"
	case LogLevelError:
		return "error"
	default:
		return "unknown"
	}
}

func InvalidCase(ll LogLevel) string {
	switch ll { // want "robustruct/linters/switch_case_cover: case value requires type related const value"
	case 0:
		return "debug"
	case 1:
		return "info"
	case 2:
		return "warn"
	case 3:
		return "error"
	default:
		return "unknown"
	}
}

func LackedCase(ll LogLevel) string {
	switch ll { // want "robustruct/linters/switch_case_cover: case body uncovered grouped const value"
	case LogLevelDebug:
		return "debug"
	case LogLevelInfo:
		return "info"
	case LogLevelWarn:
		return "warn"
	default:
		return "unknown"
	}
}

func (ll LogLevel) String() string {
	switch ll {
	case LogLevelDebug:
		return "debug"
	case LogLevelInfo:
		return "info"
	case LogLevelWarn:
		return "warn"
	case LogLevelError:
		return "error"
	default:
		return "unknown"
	}
}

func FromString(s string) LogLevel {
	switch s {
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "warn":
		return LogLevelWarn
	case "error":
		return LogLevelError
	default:
		panic("unknown log level: " + s)
	}
}
