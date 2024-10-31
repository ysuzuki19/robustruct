package external

import "fix"

func Defined123() fix.SampleStruct {
	return fix.SampleStruct{
		Field1: 1,
		Field2: "hello",
		Field3: true,
	}
}

func Defined12() fix.SampleStruct {
	return fix.SampleStruct{ // want "fields 'Field3' are not initialized"
		Field1: 1,
		Field2: "hello",
	}
}

func EmptyInit() fix.SampleStruct {
	return fix.SampleStruct{} // want "fields 'Field1, Field2, Field3' are not initialized"
}

func UndefinedStructField() fix.Sample2 {
	return fix.Sample2{ // want "fields 'Sample' are not initialized"
		Num: 1,
	}
}
