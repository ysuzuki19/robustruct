package sample

type Sample struct {
	A  int
	B  string
	AA bool
}

func Sample1() {
	s := Sample{ // want "fields 'AA' are not initialized"
		A: 1,
		B: "hello",
	}
	_ = s
}

func Sample2() {
	s := Sample{ // want "all fields of the struct must be sorted by defined order"
		B:  "hello",
		A:  1,
		AA: false,
	}
	_ = s
}
