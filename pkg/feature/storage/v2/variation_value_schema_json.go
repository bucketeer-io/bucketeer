// Copyright 2026 The Bucketeer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v2

import (
	"database/sql/driver"
	"errors"

	"google.golang.org/protobuf/encoding/protojson"
	goproto "google.golang.org/protobuf/proto"

	proto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

type VariationValueSchemaJSON struct {
	Val *proto.VariationValueSchema
}

func (o VariationValueSchemaJSON) Value() (driver.Value, error) {
	bytes, err := protojson.Marshal(o.Val)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

type VariationValueSchemaJSONScanner struct {
	Val **proto.VariationValueSchema
}

func (o VariationValueSchemaJSONScanner) Scan(src interface{}) error {
	var bytes []byte
	switch s := src.(type) {
	case []byte:
		bytes = s
	case string:
		bytes = []byte(s)
	case nil:
		*o.Val = nil
		return nil
	default:
		return errors.New("incompatible type for VariationValueSchemaJSONScanner")
	}
	schema := &proto.VariationValueSchema{}
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(bytes, schema); err != nil {
		return err
	}
	*o.Val = schema
	return nil
}

func VariationValueSchemaJSONValue(schema *proto.VariationValueSchema) interface{} {
	if schema == nil || goproto.Equal(schema, &proto.VariationValueSchema{}) {
		return nil
	}
	return VariationValueSchemaJSON{Val: schema}
}
