package for_testdocgen_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	lib "github.com/ysuzuki19/robustruct/example/for_testdocgen"
)

func TestNewUser(t *testing.T) {
	require := require.New(t)

	// testdoc begin NewUser
	u := lib.NewUser("Alice", 30)
	require.Equal(lib.User{
		Name: "Alice",
		Age:  30,
	}, *u)
	// testdoc end
}

func TestGetName(t *testing.T) {
	require := require.New(t)

	// testdoc begin GetName
	u := lib.NewUser("Alice", 30)
	require.Equal("Alice", u.GetName())
	// testdoc end
}

func TestGetAge(t *testing.T) {
	require := require.New(t)

	// testdoc begin GetAge
	u := lib.NewUser("Alice", 30)
	require.Equal(30, u.GetAge())
	// testdoc end
}

func TestSetName(t *testing.T) {
	require := require.New(t)

	// testdoc begin SetName
	u := lib.NewUser("Alice", 30)
	u.SetName("Bob")
	require.Equal("Bob", u.GetName())
	// testdoc end
}

func TestSetAge(t *testing.T) {
	require := require.New(t)

	// testdoc begin SetAge
	u := lib.NewUser("Alice", 30)
	u.SetAge(40)
	require.Equal(40, u.GetAge())
	// testdoc end
}
