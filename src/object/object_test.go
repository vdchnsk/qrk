package object

import "testing"

func TestStringHashKey(t *testing.T) {
	input1 := &String{Value: "Hello"}
	input2 := &String{Value: "Hello"}
	input3 := &String{Value: "Hello World"}

	if input1.HashKey() != input2.HashKey() {
		t.Errorf(
			"hash keys of strings with values %s and %s don't match",
			input1.Value, input2.Value,
		)
	}

	if input1.HashKey() == input3.HashKey() {
		t.Errorf(
			"hash keys of different strings match %s, %s",
			input1.Value, input3.Value,
		)
	}
}

func TestIntegerHashKey(t *testing.T) {
	input1 := &Integer{Value: 1}
	input2 := &Integer{Value: 1}
	input3 := &Integer{Value: 3}

	if input1.HashKey() != input2.HashKey() {
		t.Errorf(
			"hash keys of ints with values %d and %d don't match",
			input1.Value, input2.Value,
		)
	}

	if input1.HashKey() == input3.HashKey() {
		t.Errorf(
			"hash keys of different ints match %d, %d",
			input1.Value, input3.Value,
		)
	}
}

func TestBooleanHashKey(t *testing.T) {
	input1 := &Boolean{Value: true}
	input2 := &Boolean{Value: true}
	input3 := &Boolean{Value: false}

	if input1.HashKey() != input2.HashKey() {
		t.Errorf(
			"hash keys of booleans with values %t and %t don't match",
			input1.Value, input2.Value,
		)
	}

	if input1.HashKey() == input3.HashKey() {
		t.Errorf(
			"hash keys of different booleans match %t, %t",
			input1.Value, input3.Value,
		)
	}
}
