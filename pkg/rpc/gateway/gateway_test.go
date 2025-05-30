package gateway

import (
	"encoding/json"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestCustomJSONPb_Unmarshal(t *testing.T) {
	c := &customJSONPb{
		JSONPb: runtime.JSONPb{
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		},
	}

	tests := []struct {
		name    string
		input   []byte
		want    bool
		wantErr bool
	}{
		{
			name:    "success: boolean true",
			input:   []byte("true"),
			want:    true,
			wantErr: false,
		},
		{
			name:    "success: string 'true'",
			input:   []byte(`"true"`),
			want:    true,
			wantErr: false,
		},
		{
			name:    "success: string 'True'",
			input:   []byte(`"True"`),
			want:    true,
			wantErr: false,
		},
		{
			name:    "success: string 'TRUE'",
			input:   []byte(`"TRUE"`),
			want:    true,
			wantErr: false,
		},
		{
			name:    "success: string '1'",
			input:   []byte(`"1"`),
			want:    true,
			wantErr: false,
		},
		{
			name:    "success: boolean false",
			input:   []byte("false"),
			want:    false,
			wantErr: false,
		},
		{
			name:    "success: string 'false'",
			input:   []byte(`"false"`),
			want:    false,
			wantErr: false,
		},
		{
			name:    "success: string 'False'",
			input:   []byte(`"False"`),
			want:    false,
			wantErr: false,
		},
		{
			name:    "success: string 'FALSE'",
			input:   []byte(`"FALSE"`),
			want:    false,
			wantErr: false,
		},
		{
			name:    "success: string '0'",
			input:   []byte(`"0"`),
			want:    false,
			wantErr: false,
		},
		{
			name:    "error: invalid boolean string",
			input:   []byte(`"invalid"`),
			want:    false,
			wantErr: true,
		},
		{
			name:    "error: invalid JSON",
			input:   []byte(`{invalid json}`),
			want:    false,
			wantErr: true,
		},
		{
			name:    "error: non-boolean value",
			input:   []byte(`123`),
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got wrapperspb.BoolValue
			err := c.Unmarshal(tt.input, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Value != tt.want {
				t.Errorf("Unmarshal() got = %v, want %v", got.Value, tt.want)
			}
		})
	}
}

func TestCustomJSONPb_Unmarshal_NestedStructures(t *testing.T) {
	// Define a test message type for nested structures
	type TestMessage struct {
		BoolField1 bool                   `json:"boolField1"`
		BoolField2 bool                   `json:"boolField2"`
		Nested     map[string]bool        `json:"nested"`
		Array      []bool                 `json:"array"`
		Mixed      map[string]interface{} `json:"mixed"`
	}

	tests := []struct {
		name    string
		input   string
		check   func(t *testing.T, data []byte)
		wantErr bool
	}{
		{
			name: "nested object with boolean strings",
			input: `{
				"boolField1": "true",
				"boolField2": "false",
				"nested": {
					"key1": "TRUE",
					"key2": "0"
				}
			}`,
			check: func(t *testing.T, data []byte) {
				var result map[string]interface{}
				json.Unmarshal(data, &result)

				// Check that boolean strings were converted
				if result["boolField1"] != true {
					t.Errorf("Expected boolField1 to be true, got %v", result["boolField1"])
				}
				if result["boolField2"] != false {
					t.Errorf("Expected boolField2 to be false, got %v", result["boolField2"])
				}

				nested := result["nested"].(map[string]interface{})
				if nested["key1"] != true {
					t.Errorf("Expected nested.key1 to be true, got %v", nested["key1"])
				}
				if nested["key2"] != false {
					t.Errorf("Expected nested.key2 to be false, got %v", nested["key2"])
				}
			},
			wantErr: false,
		},
		{
			name: "array with boolean strings",
			input: `{
				"array": ["true", "false", "1", "0", "TRUE", "FALSE"]
			}`,
			check: func(t *testing.T, data []byte) {
				var result map[string]interface{}
				json.Unmarshal(data, &result)

				array := result["array"].([]interface{})
				expected := []bool{true, false, true, false, true, false}

				for i, v := range array {
					if v != expected[i] {
						t.Errorf("Expected array[%d] to be %v, got %v", i, expected[i], v)
					}
				}
			},
			wantErr: false,
		},
		{
			name: "mixed types with boolean strings",
			input: `{
				"mixed": {
					"bool1": "true",
					"bool2": true,
					"string": "not a boolean",
					"number": 42,
					"nested": {
						"deepBool": "false"
					}
				}
			}`,
			check: func(t *testing.T, data []byte) {
				var result map[string]interface{}
				json.Unmarshal(data, &result)

				mixed := result["mixed"].(map[string]interface{})
				if mixed["bool1"] != true {
					t.Errorf("Expected mixed.bool1 to be true, got %v", mixed["bool1"])
				}
				if mixed["bool2"] != true {
					t.Errorf("Expected mixed.bool2 to be true, got %v", mixed["bool2"])
				}
				if mixed["string"] != "not a boolean" {
					t.Errorf("Expected mixed.string to remain unchanged, got %v", mixed["string"])
				}
				if mixed["number"] != float64(42) {
					t.Errorf("Expected mixed.number to be 42, got %v", mixed["number"])
				}

				nested := mixed["nested"].(map[string]interface{})
				if nested["deepBool"] != false {
					t.Errorf("Expected mixed.nested.deepBool to be false, got %v", nested["deepBool"])
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// First, let's see what the conversion produces
			var rawData interface{}
			json.Unmarshal([]byte(tt.input), &rawData)
			convertBooleanStrings(rawData)
			modifiedData, _ := json.Marshal(rawData)

			// Run the check function on the modified data
			if tt.check != nil {
				tt.check(t, modifiedData)
			}
		})
	}
}
