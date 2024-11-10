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

func Sample3() {
	// ignore:robustruct
	s := Sample{
		B: "hello",
		A: 1,
	}
	_ = s
}

func Sample4() {
	s := Sample{
		A: 1,
		B: "hello",
	} // ignore:fields_require
	_ = s
}
