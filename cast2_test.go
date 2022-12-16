package cast2

import (
	"testing"
)

type Demo struct {
	Key1 string `json:"key_1"`
	Key2 int    `json:"key_2"`
}

var demo1 Demo = Demo{
	Key1: "Key1 Value",
	Key2: 2,
}

func TestToMap(t *testing.T) {
	tests := ToMap(demo1)
	if tests["Key1"] != "Key1 Value" {
		t.Error(`tests["Key1"] != "Key1 Value"`)
	}
}

func TestToMapByTag(t *testing.T) {
	tests := ToMapByTag(demo1, "json")
	if tests["key_1"] != "Key1 Value" {
		t.Error(`tests["key_1"] != "Key1 Value"`)
	}
}
