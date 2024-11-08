package fix

import "fix/sub"

type SampleStruct struct {
	Field1 int
	Field2 string
	Field3 bool
}

func DefinedAligned() SampleStruct {
	return SampleStruct{
		Field1: 1,
		Field2: "hello",
		Field3: true,
	}
}

func DefinedUnaligned() SampleStruct {
	return SampleStruct{ // want "all fields of the struct must be sorted by defined order"
		Field2: "hello",
		Field1: 1,
		Field3: true,
	}
}

func NotEnoughNotAligned() SampleStruct {
	return SampleStruct{ // want "all fields of the struct must be sorted by defined order"
		Field3: true,
		Field2: "hello",
	}
}

type NoisyStruct struct {
	Field2 int
	Field4 int
	Field0 int
}

func NoisyAligned() NoisyStruct {
	return NoisyStruct{
		Field2: 0,
		Field4: 3,
		Field0: 2,
	}
}

func NoisyUnaligned() NoisyStruct {
	return NoisyStruct{ // want "all fields of the struct must be sorted by defined order"
		Field2: 0,
		Field0: 2,
		Field4: 3,
	}
}

type AnythingTypes struct {
	Field1  int
	Field2  string
	Field3  bool
	Field4  float64
	Field5  complex128
	Field6  []int
	Field7  map[string]int
	Field8  chan int
	Field9  func()
	Field10 interface{}
	Field11 struct{}
	Field12 *int
	Field13 **int
	Field14 ***int
	Field15 [3]int
}

func AnythingTypesAligned() AnythingTypes {
	return AnythingTypes{ // want "all fields of the struct must be sorted by defined order"
		Field2:  "hello",
		Field1:  1,
		Field3:  true,
		Field4:  1.0,
		Field5:  1i,
		Field6:  []int{1},
		Field7:  map[string]int{"hello": 1},
		Field8:  make(chan int),
		Field9:  func() {},
		Field10: 1,
		Field11: struct{}{},
		Field12: new(int),
		Field13: new(*int),
		Field14: new(**int),
		Field15: [3]int{1, 2, 3},
	}
}

func SubAligned() sub.SubSample {
	return sub.SubSample{
		A: 1,
		C: 3,
	}
}

func SubUnaligned() sub.SubSample {
	return sub.SubSample{ // want "all fields of the struct must be sorted by defined order"
		C: 3,
		A: 1,
	}
}

func SubNotEnoughAligned() sub.SubSample {
	return sub.SubSample{
		A: 1,
		C: 3,
	}
}

func SubNotEnoughUnaligned() sub.SubSample {
	return sub.SubSample{ // want "all fields of the struct must be sorted by defined order"
		C: 3,
		A: 1,
	}
}
