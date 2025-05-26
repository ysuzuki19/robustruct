package ignore

func IgnoredByTest() SampleStruct {
	return SampleStruct{
		Field2: "hello",
		Field1: 1,
		Field3: true,
	}
}
