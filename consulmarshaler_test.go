package consulmarshaler_test

import (
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/nii236/consulmarshaler"
)

type KeyValue struct {
	First  string
	Second int
	Third  bool
}

type Root struct {
	Test map[string]string
	Root Nested
}

type Nested struct {
	KeyField string
	IntField int
}

var KVCase = &KeyValue{
	First:  "Hello",
	Second: 1,
	Third:  true,
}

var NestedCase = &Root{
	Test: map[string]string{"Key": "Value"},
	Root: Nested{
		KeyField: "Value",
		IntField: 5,
	},
}

func TestMarshal(t *testing.T) {
	m, err := consulmarshaler.New(api.DefaultConfig())
	if err != nil {
		t.Error(err)
	}

	err = m.Marshal("testmarshal", KVCase)
	if err != nil {
		t.Error(err)
	}

	err = m.Marshal("testmarshal2", NestedCase)
	if err != nil {
		t.Error(err)
	}
}

func TestUnmarshal(t *testing.T) {
	m, err := consulmarshaler.New(api.DefaultConfig())
	if err != nil {
		t.Error(err)
	}
	result := &KeyValue{}
	m.Unmarshal("testmarshal", result)
	if result.First != "Hello" {
		t.Errorf("Wrong value, expected %v, got %v", "Hello", result.First)
	}
	if result.Second != 1 {
		t.Errorf("Wrong value, expected %v, got %v", 1, result.Second)
	}
	if result.Third != true {
		t.Errorf("Wrong value, expected %v, got %v", true, result.Third)
	}

	nestedResult := &Root{}
	m.Unmarshal("testmarshal", nestedResult)

	if nestedResult.Root.IntField != 5 {
		t.Errorf("Wrong value, expected %v, got %v", 1, nestedResult.Root.IntField)
	}

	if nestedResult.Root.KeyField != "Value" {
		t.Errorf("Wrong value, expected %v, got %v", "Value", nestedResult.Root.KeyField)
	}
}
