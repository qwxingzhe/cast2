package cast2

import (
	"reflect"
	"testing"
)

type Demo1 struct {
	Key1 string `json:"key_1"`
	Key2 int    `json:"key_2"`
}
type Demo2 struct {
	Key1 string  `json:"key_1"`
	Key2 int     `json:"key_2"`
	Key3 float64 `json:"key_3"`
}

var demo1 = Demo1{
	Key1: "Key1 Value",
	Key2: 2,
}
var demo2 = Demo2{}

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

func TestToMapByTagJson(t *testing.T) {
	tests := ToMapByTagJson(demo1)
	if tests["key_1"] != "Key1 Value" {
		t.Error(`tests["key_1"] != "Key1 Value"`)
	}
}

func TestSetStructValue(t *testing.T) {
	demo1 = SetStructValue(demo1, "Key2", 5)
	if demo1.Key2 != 5 {
		t.Error(`demo1.Key2 != 5`)
	}
}

func TestGetStructKeyKind(t *testing.T) {
	k2Kind := GetStructKeyKind(demo1, "Key2")
	if k2Kind != reflect.Int {
		t.Error(`k2Kind != reflect.Int`)
	}
}

func TestCopyStruct(t *testing.T) {
	demo2 = CopyStruct(demo1, demo2)
	if demo2.Key2 != 2 {
		t.Error(`demo2.Key2!=2`)
	}
}

func TestCopyStructMapping(t *testing.T) {
	demo2 = CopyStructMapping(demo1, demo2, map[string]string{
		"Key2": "Key3",
	})
	if demo2.Key3 != 2 {
		t.Error(`demo2.Key3 != 2`)
	}
}
