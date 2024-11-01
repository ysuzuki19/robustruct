package e2e

type SampleStruct struct {
	Field1 int
	Field2 string
	Field3 bool
}

// succeed all
func DefinedAlignedAll() SampleStruct {
	return SampleStruct{
		Field1: 1,
		Field2: "hello",
		Field3: true,
	}
}

// failed fields_require
func DefinedLacked() SampleStruct {
	return SampleStruct{ // want "fields 'Field2, Field3' are not initialized"
		Field1: 1,
	}
}

// failed fields_align
func DefinedUnalignedAll() SampleStruct {
	return SampleStruct{ // want "all fields of the struct must be sorted by defined order"
		Field2: "hello",
		Field1: 1,
		Field3: true,
	}
}

// failed fields_require, fields_align
func DefinedLackedUnaligned() SampleStruct {
	return SampleStruct{ // want "fields 'Field2' are not initialized" "all fields of the struct must be sorted by defined order"
		Field3: true,
		Field1: 1,
	}
}

// ignore by packages
func IgnorePackages() SampleStruct {
	// ignore:fields_require
	// ignore:fields_align
	return SampleStruct{
		Field3: true,
		Field1: 1,
	}
}

// ignore by mod
func IgnoreMod() SampleStruct {
	// ignore:robustruct
	return SampleStruct{
		Field3: true,
		Field1: 1,
	}
}
