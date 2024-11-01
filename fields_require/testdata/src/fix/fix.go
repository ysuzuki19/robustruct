package fix

type SampleStruct struct {
	Field1 int
	Field2 string
	Field3 bool
}

func Defined123() SampleStruct {
	return SampleStruct{
		Field1: 1,
		Field2: "hello",
		Field3: true,
	}
}

func Defined12() SampleStruct {
	return SampleStruct{ // want "fields 'Field3' are not initialized"
		Field1: 1,
		Field2: "hello",
	}
}

func Defined23() SampleStruct {
	return SampleStruct{ // want "fields 'Field1' are not initialized"
		Field2: "hello",
		Field3: true,
	}
}

func Defined31() SampleStruct {
	return SampleStruct{ // want "fields 'Field2' are not initialized"
		Field3: true,
		Field1: 1,
	}
}

func Defined1() SampleStruct {
	return SampleStruct{ // want "fields 'Field2, Field3' are not initialized"
		Field1: 1,
	}
}

func Defined2() SampleStruct {
	return SampleStruct{ // want "fields 'Field1, Field3' are not initialized"
		Field2: "hello",
	}
}

func Defined3() SampleStruct {
	return SampleStruct{ // want "fields 'Field1, Field2' are not initialized"
		Field3: true,
	}
}

func EmptyInit() SampleStruct {
	return SampleStruct{} // want "fields 'Field1, Field2, Field3' are not initialized"
}

func DefinedUnnamed123() SampleStruct {
	return SampleStruct{0, "", false}
}

type Sample2 struct {
	Num    int
	Sample SampleStruct
}

func UndefinedStructField() Sample2 {
	return Sample2{ // want "fields 'Sample' are not initialized"
		Num: 1,
	}
}

func UndefinedNestedStructField() Sample2 {
	return Sample2{
		Num: 1,
		Sample: SampleStruct{ // want "fields 'Field2, Field3' are not initialized"
			Field1: 1,
		},
	}
}