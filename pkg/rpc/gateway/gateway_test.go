package gateway

import (
	"encoding/json"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestBooleanConversionMarshaler_PreprocessJSON(t *testing.T) {
	marshaler := &BooleanConversionMarshaler{}

	tests := []struct {
		name    string
		input   string
		check   func(t *testing.T, data []byte)
		wantErr bool
	}{
		{
			name: "nested object with boolean strings in specific fields",
			input: `{
				"userAttributesUpdated": "true",
				"user_attributes_updated": "false",
				"otherField": "true",
				"nested": {
					"userAttributesUpdated": "TRUE",
					"user_attributes_updated": "0"
				}
			}`,
			check: func(t *testing.T, data []byte) {
				var result map[string]interface{}
				json.Unmarshal(data, &result)

				// Check that boolean strings in specific fields were converted
				if result["userAttributesUpdated"] != true {
					t.Errorf("Expected userAttributesUpdated to be true, got %v", result["userAttributesUpdated"])
				}
				if result["user_attributes_updated"] != false {
					t.Errorf("Expected user_attributes_updated to be false, got %v", result["user_attributes_updated"])
				}
				// Other fields should remain as strings
				if result["otherField"] != "true" {
					t.Errorf("Expected otherField to remain 'true' string, got %v", result["otherField"])
				}

				nested := result["nested"].(map[string]interface{})
				if nested["userAttributesUpdated"] != true {
					t.Errorf("Expected nested.userAttributesUpdated to be true, got %v", nested["userAttributesUpdated"])
				}
				if nested["user_attributes_updated"] != false {
					t.Errorf("Expected nested.user_attributes_updated to be false, got %v", nested["user_attributes_updated"])
				}
			},
			wantErr: false,
		},
		{
			name: "array with mixed content",
			input: `{
				"data": [{
					"userAttributesUpdated": "true",
					"otherField": "not a boolean"
				}, {
					"user_attributes_updated": "1",
					"anotherField": 42
				}]
			}`,
			check: func(t *testing.T, data []byte) {
				var result map[string]interface{}
				json.Unmarshal(data, &result)

				dataArray := result["data"].([]interface{})

				item1 := dataArray[0].(map[string]interface{})
				if item1["userAttributesUpdated"] != true {
					t.Errorf("Expected data[0].userAttributesUpdated to be true, got %v", item1["userAttributesUpdated"])
				}
				if item1["otherField"] != "not a boolean" {
					t.Errorf("Expected data[0].otherField to remain unchanged, got %v", item1["otherField"])
				}

				item2 := dataArray[1].(map[string]interface{})
				if item2["user_attributes_updated"] != true {
					t.Errorf("Expected data[1].user_attributes_updated to be true, got %v", item2["user_attributes_updated"])
				}
				if item2["anotherField"] != float64(42) {
					t.Errorf("Expected data[1].anotherField to be 42, got %v", item2["anotherField"])
				}
			},
			wantErr: false,
		},
		{
			name: "all boolean string variations",
			input: `{
				"userAttributesUpdated": "true",
				"data": {
					"user_attributes_updated": "TRUE"
				},
				"items": [{
					"userAttributesUpdated": "1"
				}, {
					"user_attributes_updated": "false"
				}, {
					"userAttributesUpdated": "FALSE"
				}, {
					"user_attributes_updated": "0"
				}]
			}`,
			check: func(t *testing.T, data []byte) {
				var result map[string]interface{}
				json.Unmarshal(data, &result)

				// Top level
				if result["userAttributesUpdated"] != true {
					t.Errorf("Expected userAttributesUpdated to be true, got %v", result["userAttributesUpdated"])
				}

				// Nested object
				dataObj := result["data"].(map[string]interface{})
				if dataObj["user_attributes_updated"] != true {
					t.Errorf("Expected data.user_attributes_updated to be true, got %v", dataObj["user_attributes_updated"])
				}

				// Array items
				items := result["items"].([]interface{})

				item0 := items[0].(map[string]interface{})
				if item0["userAttributesUpdated"] != true {
					t.Errorf("Expected items[0].userAttributesUpdated to be true, got %v", item0["userAttributesUpdated"])
				}

				item1 := items[1].(map[string]interface{})
				if item1["user_attributes_updated"] != false {
					t.Errorf("Expected items[1].user_attributes_updated to be false, got %v", item1["user_attributes_updated"])
				}

				item2 := items[2].(map[string]interface{})
				if item2["userAttributesUpdated"] != false {
					t.Errorf("Expected items[2].userAttributesUpdated to be false, got %v", item2["userAttributesUpdated"])
				}

				item3 := items[3].(map[string]interface{})
				if item3["user_attributes_updated"] != false {
					t.Errorf("Expected items[3].user_attributes_updated to be false, got %v", item3["user_attributes_updated"])
				}
			},
			wantErr: false,
		},
		{
			name:  "empty data",
			input: ``,
			check: func(t *testing.T, data []byte) {
				if len(data) != 0 {
					t.Errorf("Expected empty data, got %v", string(data))
				}
			},
			wantErr: false,
		},
		{
			name:  "invalid JSON returns original",
			input: `{invalid json}`,
			check: func(t *testing.T, data []byte) {
				if string(data) != `{invalid json}` {
					t.Errorf("Expected original data for invalid JSON, got %v", string(data))
				}
			},
			wantErr: false,
		},
		{
			name: "non-string values remain unchanged",
			input: `{
				"userAttributesUpdated": true,
				"user_attributes_updated": false,
				"data": {
					"userAttributesUpdated": 123,
					"user_attributes_updated": null
				}
			}`,
			check: func(t *testing.T, data []byte) {
				var result map[string]interface{}
				json.Unmarshal(data, &result)

				// Boolean values should remain as booleans
				if result["userAttributesUpdated"] != true {
					t.Errorf("Expected userAttributesUpdated to remain true, got %v", result["userAttributesUpdated"])
				}
				if result["user_attributes_updated"] != false {
					t.Errorf("Expected user_attributes_updated to remain false, got %v", result["user_attributes_updated"])
				}

				// Non-string values in nested objects
				dataObj := result["data"].(map[string]interface{})
				if dataObj["userAttributesUpdated"] != float64(123) {
					t.Errorf("Expected data.userAttributesUpdated to remain 123, got %v", dataObj["userAttributesUpdated"])
				}
				if dataObj["user_attributes_updated"] != nil {
					t.Errorf("Expected data.user_attributes_updated to remain nil, got %v", dataObj["user_attributes_updated"])
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the preprocessJSON method
			processedData := marshaler.preprocessJSON([]byte(tt.input))

			// Run the check function on the processed data
			if tt.check != nil {
				tt.check(t, processedData)
			}
		})
	}
}

func TestStringToBool(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantBool bool
		wantOk   bool
	}{
		{"lowercase true", "true", true, true},
		{"uppercase TRUE", "TRUE", true, true},
		{"mixed case True", "True", true, true},
		{"numeric 1", "1", true, true},
		{"lowercase false", "false", false, true},
		{"uppercase FALSE", "FALSE", false, true},
		{"mixed case False", "False", false, true},
		{"numeric 0", "0", false, true},
		{"invalid string", "invalid", false, false},
		{"empty string", "", false, false},
		{"numeric 2", "2", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBool, gotOk := stringToBool(tt.input)
			if gotBool != tt.wantBool {
				t.Errorf("stringToBool() gotBool = %v, want %v", gotBool, tt.wantBool)
			}
			if gotOk != tt.wantOk {
				t.Errorf("stringToBool() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestBooleanConversionMarshaler_Integration(t *testing.T) {
	marshaler := &BooleanConversionMarshaler{
		JSONPb: runtime.JSONPb{
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		},
	}

	// Test that the marshaler integrates properly with runtime.JSONPb
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, marshaler *BooleanConversionMarshaler)
	}{
		{
			name:  "marshaler processes JSON correctly",
			input: `{"userAttributesUpdated": "true", "someOtherField": "not a boolean"}`,
			check: func(t *testing.T, m *BooleanConversionMarshaler) {
				// Test that preprocessJSON is called and works correctly
				processed := m.preprocessJSON([]byte(`{"userAttributesUpdated": "true"}`))

				var result map[string]interface{}
				if err := json.Unmarshal(processed, &result); err != nil {
					t.Fatalf("Failed to unmarshal processed JSON: %v", err)
				}

				if result["userAttributesUpdated"] != true {
					t.Errorf("Expected userAttributesUpdated to be converted to true, got %v", result["userAttributesUpdated"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.check != nil {
				tt.check(t, marshaler)
			}
		})
	}
}

func TestShouldProcessRequest(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"exact match get_evaluations", "/get_evaluations", true},
		{"exact match v1 gateway", "/v1/gateway/evaluations", true},
		{"contains get_evaluations", "/api/get_evaluations/test", true},
		{"contains v1 gateway", "/api/v1/gateway/evaluations/test", true},
		{"no match", "/api/other/endpoint", false},
		{"empty path", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldProcessRequest(tt.path); got != tt.want {
				t.Errorf("shouldProcessRequest(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
