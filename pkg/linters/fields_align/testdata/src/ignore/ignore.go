package ignore

type SampleStruct struct {
	Field1 int
	Field2 string
	Field3 bool
}

func DefinedUnaligned() SampleStruct {
	return SampleStruct{ // want "all fields of the struct must be sorted by defined order"
		Field2: "hello",
		Field1: 1,
		Field3: true,
	}
}

func Ignored() SampleStruct {
	// ignore:fields_align
	return SampleStruct{
		Field2: "hello",
		Field1: 1,
		Field3: true,
	}
}

func IgnoredEndOfStruct() SampleStruct {
	return SampleStruct{
		Field2: "hello",
		Field1: 1,
		Field3: true,
	} // ignore:fields_align
}
