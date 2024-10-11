package utils

func ConvertSingleOrArrayEitherToArray[A any](either Either[A, []A]) []A {
	var array []A
	if either.IsA() {
		array = []A{either.GetA()}
	} else if either.IsB() {
		array = either.GetB()
	} else {
		array = []A{}
	}

	return array
}
