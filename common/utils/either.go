package utils

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type Either[A, B any] struct {
	value any
}

func (e Either[A, B]) IsA() bool {
	_, ok := e.value.(A)
	return ok
}

func (e Either[A, B]) IsB() bool {
	_, ok := e.value.(B)
	return ok
}

func (e Either[A, B]) GetA() A {
	return e.value.(A)
}

func (e Either[A, B]) GetB() B {
	return e.value.(B)
}

func (e *Either[A, B]) SetA(a A) {
	e.value = a
}

func (e *Either[A, B]) SetB(b B) {
	e.value = b
}

func (e Either[A, B]) IsNil() bool {
	return e.value == nil
}

func (input *Either[A, B]) UnmarshalYAML(value *yaml.Node) (err error) {
	err = nil

	var a A
	err = value.Decode(&a)
	if err != nil {
		var b B
		err = value.Decode(&b)
		input.value = b
		return err
	}

	input.value = a
	return err
}

func (input Either[A, B]) MarshalYAML() (interface{}, error) {
	if input.IsA() {
		return input.value.(A), nil
	} else {
		return input.value.(B), nil
	}
}

func (input Either[A, B]) MarshalJSON() ([]byte, error) {
	if input.IsA() {
		return json.Marshal(input.value.(A))
	} else {
		return json.Marshal(input.value.(B))
	}
}
