package utils

func ErrHandle(err error) {
	if err != nil {
		panic(err)
	}
}
