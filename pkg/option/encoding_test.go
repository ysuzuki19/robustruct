package option_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ysuzuki19/robustruct/pkg/option"
)

type User struct {
	Name string             `json:"name"`
	Age  option.Option[int] `json:"age"`
}

type Case struct {
	User User
	Json string
}

func cases() []Case {
	return []Case{
		{
			User: User{
				Name: "Alice",
				Age:  option.NewSome(25),
			},
			Json: `{"name":"Alice","age":25}`,
		},
		{
			User: User{
				Name: "Bob",
				Age:  option.None[int](),
			},
			Json: `{"name":"Bob","age":null}`,
		},
	}
}

func TestMarshalJson(t *testing.T) {
	require := require.New(t)

	for _, c := range cases() {
		b, err := json.Marshal(c.User)
		require.NoError(err)
		require.Equal(c.Json, string(b))
	}
}

func TestUnmarshalJSON(t *testing.T) {
	require := require.New(t)

	for _, c := range cases() {
		var user User
		err := json.Unmarshal([]byte(c.Json), &user)
		require.NoError(err)
		require.Equal(c.User, user)
	}
}
