package spkg_tst5

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}
