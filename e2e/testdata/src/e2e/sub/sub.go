package sub

type SubSample struct {
	A int
	b int
}

type SubPubOnly struct {
	A int
	B int
}

type SubPrivOnly struct {
	a int
	b int
}

func AvoidLint() {
	ss := SubSample{}
	print(ss.A, ss.b)
	pub := SubPubOnly{}
	print(pub.A, pub.B)
	priv := SubPrivOnly{}
	print(priv.a, priv.b)
}
