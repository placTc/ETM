package utils

func PanicOnError(e error) {
	if e != nil {
		panic(e)
	}
}

func PanicOnErrorWithCustomMessage(e error, message string) {
	if e != nil {
		panic(message + "\n" + e.Error())
	}
}
