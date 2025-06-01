//go:generate go run github.com/ysuzuki19/robustruct/cmd/gen/testdocgen -file=$GOFILE
package lib

type User struct {
	Name string
	Age  int
}

// Example:
//
//	u := lib.NewUser("Alice", 30)
//	require.Equal(lib.User{
//		Name: "Alice",
//		Age:  30,
//	}, *u)
func NewUser(name string, age int) *User {
	return &User{
		Name: name,
		Age:  age,
	}
}

// Example:
//
//	u := lib.NewUser("Alice", 30)
//	require.Equal("Alice", u.GetName())
func (u User) GetName() string {
	return u.Name
}

// Example:
//
//	u := lib.NewUser("Alice", 30)
//	require.Equal(30, u.GetAge())
func (u User) GetAge() int {
	return u.Age
}

// Example:
//
//	u := lib.NewUser("Alice", 30)
//	u.SetName("Bob")
//	require.Equal("Bob", u.GetName())
func (u *User) SetName(name string) {
	u.Name = name
}

// Example:
//
//	u := lib.NewUser("Alice", 30)
//	u.SetAge(40)
//	require.Equal(40, u.GetAge())
func (u *User) SetAge(age int) {
	u.Age = age
}
