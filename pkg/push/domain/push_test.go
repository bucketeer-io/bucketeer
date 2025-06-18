// Copyright 2025 The Bucketeer Authors.
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

package domain

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pushproto "github.com/bucketeer-io/bucketeer/proto/push"
)

func TestNewPush(t *testing.T) {
	t.Parallel()
	
	patterns := []struct {
		name           string
		pushName       string
		serviceAccount string
		tags           []string
		expected       func(*testing.T, *Push, error)
	}{
		{
			name:           "with tags",
			pushName:       "name-1",
			serviceAccount: "sa",
			tags:           []string{"tag-1", "tag-2"},
			expected: func(t *testing.T, push *Push, err error) {
				assert.NoError(t, err)
				assert.IsType(t, &Push{}, push)
				assert.NotEqual(t, "", push.Id)
				assert.Equal(t, "sa", push.FcmServiceAccount)
				assert.Equal(t, []string{"tag-1", "tag-2"}, push.Tags)
				assert.Equal(t, 2, len(push.Tags))
			},
		},
		{
			name:           "without tags",
			pushName:       "name-1",
			serviceAccount: "sa",
			tags:           []string{},
			expected: func(t *testing.T, push *Push, err error) {
				assert.NoError(t, err)
				assert.IsType(t, &Push{}, push)
				assert.NotEqual(t, "", push.Id)
				assert.Equal(t, "sa", push.FcmServiceAccount)
				assert.Equal(t, []string{}, push.Tags)
				assert.Equal(t, 0, len(push.Tags))
			},
		},
		{
			name:           "with nil tags",
			pushName:       "name-1",
			serviceAccount: "sa",
			tags:           nil,
			expected: func(t *testing.T, push *Push, err error) {
				assert.NoError(t, err)
				assert.IsType(t, &Push{}, push)
				assert.NotEqual(t, "", push.Id)
				assert.Equal(t, "sa", push.FcmServiceAccount)
				assert.Nil(t, push.Tags)
				assert.Equal(t, 0, len(push.Tags))
			},
		},
	}
	
	for _, p := range patterns {
		t.Run(p.name, func(t *testing.T) {
			actual, err := NewPush(p.pushName, p.serviceAccount, p.tags)
			p.expected(t, actual, err)
		})
	}
}

func TestSetDeleted(t *testing.T) {
	t.Parallel()
	name := "name-1"
	serviceAccount := "sa"
	tags := []string{"tag-1", "tag-2"}
	actual, err := NewPush(name, serviceAccount, tags)
	assert.NoError(t, err)
	assert.Equal(t, false, actual.Deleted)
	actual.SetDeleted()
	assert.Equal(t, true, actual.Deleted)
}

func TestAddTags(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		origin      *Push
		input       []string
		expectedErr error
		expected    []string
	}{
		{
			desc:        "success: one",
			origin:      &Push{&pushproto.Push{Tags: []string{"tag-0", "tag-1"}}},
			input:       []string{"tag-2"},
			expectedErr: nil,
			expected:    []string{"tag-0", "tag-1", "tag-2"},
		},
		{
			desc:        "success: two",
			origin:      &Push{&pushproto.Push{Tags: []string{"tag-0", "tag-1"}}},
			input:       []string{"tag-2", "tag-3"},
			expectedErr: nil,
			expected:    []string{"tag-0", "tag-1", "tag-2", "tag-3"},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			err := p.origin.AddTags(p.input)
			assert.Equal(t, p.expectedErr, err)
			sort.Strings(p.expected)
			sort.Strings(p.origin.Tags)
			assert.Equal(t, p.expected, p.origin.Tags)
		})
	}
}

func TestDeleteTags(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		origin      *Push
		input       []string
		expectedErr error
		expected    []string
	}{
		{
			desc:        "success: one",
			origin:      &Push{&pushproto.Push{Tags: []string{"tag-0", "tag-1"}}},
			input:       []string{"tag-1"},
			expectedErr: nil,
			expected:    []string{"tag-0"},
		},
		{
			desc:        "success: two",
			origin:      &Push{&pushproto.Push{Tags: []string{"tag-0", "tag-1"}}},
			input:       []string{"tag-0", "tag-1"},
			expectedErr: nil,
			expected:    []string{},
		},
		{
			desc:        "fail: not found: one",
			origin:      &Push{&pushproto.Push{Tags: []string{"tag-0", "tag-1"}}},
			input:       []string{"tag-2"},
			expectedErr: ErrTagNotFound,
			expected:    []string{"tag-0", "tag-1"},
		},
		{
			desc:        "fail: not found: two",
			origin:      &Push{&pushproto.Push{Tags: []string{"tag-0", "tag-1"}}},
			input:       []string{"tag-0", "tag-2"},
			expectedErr: ErrTagNotFound,
			expected:    []string{"tag-0", "tag-1"},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			err := p.origin.DeleteTags(p.input)
			assert.Equal(t, p.expectedErr, err)
			sort.Strings(p.expected)
			sort.Strings(p.origin.Tags)
			assert.Equal(t, p.expected, p.origin.Tags)
		})
	}
}

func TestExistTag(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc     string
		origin   *Push
		input    string
		expected bool
	}{
		{
			desc:     "true",
			origin:   &Push{&pushproto.Push{Tags: []string{"tag-0", "tag-1"}}},
			input:    "tag-1",
			expected: true,
		},
		{
			desc:     "false: no tags",
			origin:   &Push{&pushproto.Push{}},
			input:    "tag-1",
			expected: false,
		},
		{
			desc:     "false: not found",
			origin:   &Push{&pushproto.Push{Tags: []string{"tag-0", "tag-1"}}},
			input:    "tag-2",
			expected: false,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := p.origin.ExistTag(p.input)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestRename(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		origin      *Push
		input       string
		expectedErr error
		expected    string
	}{
		{
			desc:        "success",
			origin:      &Push{&pushproto.Push{Name: "a"}},
			input:       "b",
			expectedErr: nil,
			expected:    "b",
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			err := p.origin.Rename(p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, p.origin.Name)
		})
	}
}

func TestUpdatePush(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		desc        string
		origin      *Push
		inputName   *wrapperspb.StringValue
		inputTags   []string
		expectedErr error
		expected    *Push
	}{
		{
			desc: "success",
			origin: &Push{&pushproto.Push{
				Name:     "a",
				Tags:     []string{"tag-0"},
				Disabled: false,
			}},
			inputName:   &wrapperspb.StringValue{Value: "b"},
			inputTags:   []string{"tag-0", "tag-1"},
			expectedErr: nil,
			expected: &Push{&pushproto.Push{
				Name:     "b",
				Tags:     []string{"tag-0", "tag-1"},
				Disabled: true,
			}},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			updatedPush, err := p.origin.Update(p.inputName, p.inputTags, wrapperspb.Bool(true))
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected.Name, updatedPush.Name)
			assert.Equal(t, p.expected.Tags, updatedPush.Tags)
			assert.Equal(t, p.expected.Disabled, updatedPush.Disabled)
		})
	}
}
