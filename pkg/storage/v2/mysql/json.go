// Copyright 2024 The Bucketeer Authors.
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

package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONObject struct {
	Val interface{}
}

func (o JSONObject) Value() (driver.Value, error) {
	return json.Marshal(o.Val)
}

func (o *JSONObject) Scan(src interface{}) error {
	var _src []byte
	switch s := src.(type) {
	case []byte:
		_src = s
	case nil:
		return nil
	default:
		return errors.New("incompatible type for JSONObject")
	}
	return json.Unmarshal(_src, &o.Val)
}
