package helper

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	Field1 string `json:"field1"`
	Field2 string
}

func TestFindJsonField(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		key   string
		want  bool
	}{
		{"field1 should be valid", TestStruct{}, "field1", true},
		{"field2 should be valid", TestStruct{}, "field2", false},
		{"field3 should be invalid", TestStruct{}, "field3", false},
		{"slice should be invalid", []TestStruct{}, "test", false},
		{"int should be invalid", 0, "test", false},
		{"string should be invalid", "", "test", false},
		{"struct pointer should be valid", &TestStruct{}, "field1", true},
		{"slice pointer should be invalid", &[]TestStruct{}, "test", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := reflect.ValueOf(tt.value)
			if got := FindJsonField(value, tt.key); got.IsValid() != tt.want {
				t.Errorf("FindJsonField(%s) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestMergeDeep(t *testing.T) {
	tests := []struct {
		name   string
		a      map[string]interface{}
		b      map[string]interface{}
		target map[string]interface{}
	}{
		{"should merge empty maps", map[string]interface{}{}, map[string]interface{}{}, map[string]interface{}{}},
		{"should add new fields", map[string]interface{}{
			"field1": "value1",
		}, map[string]interface{}{
			"field2": "value2",
		}, map[string]interface{}{
			"field1": "value1",
			"field2": "value2",
		}},
		{"should add new nested maps", map[string]interface{}{
			"field1": "value1",
		}, map[string]interface{}{
			"field2": map[string]interface{}{
				"nested": true,
			},
		}, map[string]interface{}{
			"field1": "value1",
			"field2": map[string]interface{}{
				"nested": true,
			},
		}},
		{"should add new slices", map[string]interface{}{
			"field1": "value1",
		}, map[string]interface{}{
			"field2": []interface{}{"nested"},
		}, map[string]interface{}{
			"field1": "value1",
			"field2": []interface{}{"nested"},
		}},
		{"should override existing fields", map[string]interface{}{
			"field1": "value1",
		}, map[string]interface{}{
			"field1": "newValue",
		}, map[string]interface{}{
			"field1": "newValue",
		}},
		{"should override nested maps", map[string]interface{}{
			"field1": "value1",
		}, map[string]interface{}{
			"field1": map[string]interface{}{
				"nested": true,
			},
		}, map[string]interface{}{
			"field1": map[string]interface{}{
				"nested": true,
			},
		}},
		{"should override nested slices", map[string]interface{}{
			"field1": "value1",
		}, map[string]interface{}{
			"field1": []interface{}{"nested"},
		}, map[string]interface{}{
			"field1": []interface{}{"nested"},
		}},
		{"should merge slices", map[string]interface{}{
			"field1": []interface{}{"a"},
		}, map[string]interface{}{
			"field1": []interface{}{"b"},
		}, map[string]interface{}{
			"field1": []interface{}{"a", "b"},
		}},
		{"should merge nested maps", map[string]interface{}{
			"field1": map[string]interface{}{
				"nested1": "A",
			},
		}, map[string]interface{}{
			"field1": map[string]interface{}{
				"nested2": "B",
			},
		}, map[string]interface{}{
			"field1": map[string]interface{}{
				"nested1": "A",
				"nested2": "B",
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeDeep(tt.a, tt.b); !reflect.DeepEqual(got, tt.target) {
				t.Errorf("MergeDeep() = %v, want %v", got, tt.target)
			}
		})
	}
}

type NestedStruct struct {
	NestedField string `json:"nestedField"`
}

type MappedStruct struct {
	Field1       string `json:"field1"`
	Field2       string
	NestedMap    map[string]int `json:"nestedMap"`
	NestedSlice  []int          `json:"nestedSlice"`
	NestedStruct NestedStruct   `json:"nestedStruct"`
}

func TestMapToStruct(t *testing.T) {
	tests := []struct {
		name string
		data map[string]interface{}
		want MappedStruct
	}{
		{"should return empty struct", map[string]interface{}{}, MappedStruct{}},
		{"should assign basic values", map[string]interface{}{
			"field1": "value1",
		}, MappedStruct{
			Field1: "value1",
		}},
		{"should skip fields without json tag", map[string]interface{}{
			"Field2": "value2",
		}, MappedStruct{}},
		{"should set nested map", map[string]interface{}{
			"nestedMap": map[string]interface{}{
				"key": 1,
			},
		}, MappedStruct{
			NestedMap: map[string]int{
				"key": 1,
			},
		}},
		{"should set nested slice", map[string]interface{}{
			"nestedSlice": []interface{}{1},
		}, MappedStruct{
			NestedSlice: []int{1},
		}},
		// {"should set nested struct", map[string]interface{}{
		// 	"nestedStruct": map[string]interface{}{
		// 		"nestedField": "nestedValue",
		// 	},
		// }, MappedStruct{
		// 	NestedStruct: NestedStruct{"nestedValue"},
		// }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapToStruct(tt.data, MappedStruct{}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapToStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}
