package sub

type SubSample struct {
	A int
	b int
	C int
}

type SubPubOnly struct {
	A int
	B int
	C int
}

func AvoidLint() {
	ss := SubSample{}
	print(ss.A, ss.b, ss.C)
	pub := SubPubOnly{}
	print(pub.A, pub.B, pub.C)
}
