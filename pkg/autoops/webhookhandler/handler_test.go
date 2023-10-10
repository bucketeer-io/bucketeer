// Copyright 2023 The Bucketeer Authors.
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

package webhookhandler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	authclientmock "github.com/bucketeer-io/bucketeer/pkg/auth/client/mock"
	autoopsdomain "github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	featureclientmock "github.com/bucketeer-io/bucketeer/pkg/feature/client/mock"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

type dummyWebhookCryptoUtil struct{}

func (u *dummyWebhookCryptoUtil) Encrypt(ctx context.Context, data []byte) ([]byte, error) {
	return []byte(data), nil
}

func (u *dummyWebhookCryptoUtil) Decrypt(ctx context.Context, data []byte) ([]byte, error) {
	return []byte(data), nil
}

func TestNewHandler(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	issuer := "test_issuer"
	clientID := "test_client_id"
	testcases := []struct {
		desc       string
		rawIDToken string
		valid      bool
	}{
		{
			desc:       "err: malformed jwt",
			rawIDToken: "",
			valid:      false,
		},
		{
			desc:       "err: invalid jwt",
			rawIDToken: "testdata/invalid-token",
			valid:      false,
		},
		{
			desc:       "success",
			rawIDToken: "testdata/valid-token",
			valid:      true,
		},
	}
	verifier, err := token.NewVerifier("testdata/valid-public.pem", issuer, clientID)
	require.NoError(t, err)
	for _, p := range testcases {
		t.Run(p.desc, func(t *testing.T) {
			h, err := NewHandler(
				mysqlmock.NewMockClient(mockController),
				authclientmock.NewMockClient(mockController),
				featureclientmock.NewMockClient(mockController),
				publishermock.NewMockPublisher(mockController),
				verifier,
				p.rawIDToken,
				&dummyWebhookCryptoUtil{},
			)
			if p.valid {
				assert.NotNil(t, h)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, h)
				assert.Error(t, err)
			}
		})
	}
}

func TestServeHTTP(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	h := &handler{
		mysqlClient:       mysqlmock.NewMockClient(mockController),
		authClient:        authclientmock.NewMockClient(mockController),
		featureClient:     featureclientmock.NewMockClient(mockController),
		webhookCryptoUtil: &dummyWebhookCryptoUtil{},
		logger:            logger,
	}
	patterns := []struct {
		desc     string
		setup    func(*testing.T, *handler)
		input    *http.Request
		expected int
	}{
		{
			desc: "fail: bad params",
			input: httptest.NewRequest("POST",
				"/hook?foo=bar",
				nil),
			expected: http.StatusBadRequest,
		},
		{
			desc: "fail: auth error",
			input: httptest.NewRequest("POST",
				"/hook?auth=secret",
				nil),
			expected: http.StatusInternalServerError,
		},
		{
			desc: "success",
			setup: func(t *testing.T, h *handler) {
				h.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				h.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			// The test secret below will return the following proto
			// &autoopsdomain.webhookSecret {
			// 	WebhookID: "id-0",
			// 	EnvironmentNamespace: "ns0",
			// }
			input: httptest.NewRequest("POST",
				"/hook?auth=eyJ3ZWJob29rX2lkIjoiaWQtMCIsImVudmlyb25tZW50X25hbWVzcGFjZSI6Im5zMCJ9",
				nil),
			expected: http.StatusOK,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			if p.setup != nil {
				p.setup(t, h)
			}
			actual := httptest.NewRecorder()
			h.ServeHTTP(actual, p.input)
			assert.Equal(t, p.expected, actual.Code)
		})
	}
}

func TestHandleWebhook(t *testing.T) {
	t.Parallel()
	convert := func(wcs []autoopsproto.WebhookClause) []*autoopsproto.Clause {
		var clauses []*autoopsproto.Clause
		for _, w := range wcs {
			c, err := ptypes.MarshalAny(&w)
			require.NoError(t, err)
			clauses = append(clauses, &autoopsproto.Clause{Clause: c})
		}
		return clauses
	}
	ctx := context.TODO()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	cases := []struct {
		name        string
		webhookId   string
		payload     string
		autoOpsRule *autoopsdomain.AutoOpsRule
		wantErr     bool
		expected    bool
	}{
		{
			name:      "1. execute rule-1",
			webhookId: "webhook-1",
			payload:   `{"body":{"Alert id": 123}}`,
			expected:  true,
			autoOpsRule: &autoopsdomain.AutoOpsRule{
				AutoOpsRule: &autoopsproto.AutoOpsRule{
					Id: "rule-1",
					Clauses: convert([]autoopsproto.WebhookClause{
						{
							WebhookId: "webhook-1",
							Conditions: []*autoopsproto.WebhookClause_Condition{
								{
									Filter:   `.body."Alert id"`,
									Value:    `123`,
									Operator: autoopsproto.WebhookClause_Condition_EQUAL,
								},
							},
						},
					}),
				},
			},
		},
		{
			name:      "2. execute rule-1",
			webhookId: "webhook-1",
			payload:   `{"body":{"Alert id": 123, "Status": "Open", "Impacted plays": 10}}`,
			expected:  true,
			autoOpsRule: &autoopsdomain.AutoOpsRule{
				AutoOpsRule: &autoopsproto.AutoOpsRule{
					Id: "rule-1",
					Clauses: convert([]autoopsproto.WebhookClause{
						{
							WebhookId: "webhook-1",
							Conditions: []*autoopsproto.WebhookClause_Condition{
								{
									Filter:   `.body."Alert id"`,
									Value:    `123`,
									Operator: autoopsproto.WebhookClause_Condition_EQUAL,
								},
								{
									Filter:   `.body.Status`,
									Value:    `"Open"`,
									Operator: autoopsproto.WebhookClause_Condition_EQUAL,
								},
								{
									Filter:   `.body."Impacted plays"`,
									Value:    `5`,
									Operator: autoopsproto.WebhookClause_Condition_MORE_THAN_OR_EQUAL,
								},
							},
						},
					}),
				},
			},
		},
		{
			name:      "execute rule-1 with converting",
			webhookId: "webhook-1",
			payload:   `{"body":{"Alert id": 123, "Status": "Open", "Impacted users": "100"}}`,
			expected:  true,
			autoOpsRule: &autoopsdomain.AutoOpsRule{
				AutoOpsRule: &autoopsproto.AutoOpsRule{
					Id: "rule-1",
					Clauses: convert([]autoopsproto.WebhookClause{
						{
							WebhookId: "webhook-1",
							Conditions: []*autoopsproto.WebhookClause_Condition{
								{
									Filter:   `.body."Alert id"`,
									Value:    `"123"`,
									Operator: autoopsproto.WebhookClause_Condition_EQUAL,
								},
								{
									Filter:   `.body.Status`,
									Value:    `"Open"`,
									Operator: autoopsproto.WebhookClause_Condition_EQUAL,
								},
								{
									Filter:   `.body."Impacted users"`,
									Value:    `50`,
									Operator: autoopsproto.WebhookClause_Condition_MORE_THAN_OR_EQUAL,
								},
							},
						},
					}),
				},
			},
		},
		{
			name:      "execute rule-1 with string comparison for lexical order",
			webhookId: "webhook-1",
			payload:   `{"body":{"foo": "abc", "bar": "ABC"}}`,
			expected:  true,
			autoOpsRule: &autoopsdomain.AutoOpsRule{
				AutoOpsRule: &autoopsproto.AutoOpsRule{
					Id: "rule-1",
					Clauses: convert([]autoopsproto.WebhookClause{
						{
							WebhookId: "webhook-1",
							Conditions: []*autoopsproto.WebhookClause_Condition{
								{
									Filter:   `.body.foo`,
									Value:    `"ab"`,
									Operator: autoopsproto.WebhookClause_Condition_MORE_THAN,
								},
								{
									Filter:   `.body.bar`,
									Value:    `"XYZ"`,
									Operator: autoopsproto.WebhookClause_Condition_LESS_THAN,
								},
							},
						},
					}),
				},
			},
		},
		{
			name:      "execute nothing because of conditions not matched",
			webhookId: "webhook-1",
			payload:   `{"body":{"Alert id": 123, "Status": "Open", "Impacted plays": 3}}`,
			expected:  false,
			autoOpsRule: &autoopsdomain.AutoOpsRule{
				AutoOpsRule: &autoopsproto.AutoOpsRule{
					Id: "rule-1",
					Clauses: convert([]autoopsproto.WebhookClause{
						{
							WebhookId: "webhook-1",
							Conditions: []*autoopsproto.WebhookClause_Condition{
								{
									Filter:   `.body."Alert id"`,
									Value:    `123`,
									Operator: autoopsproto.WebhookClause_Condition_EQUAL,
								},
								{
									Filter:   `.body.Status`,
									Value:    `"Open"`,
									Operator: autoopsproto.WebhookClause_Condition_EQUAL,
								},
								{
									Filter:   `.body."Impacted plays"`,
									Value:    `10`,
									Operator: autoopsproto.WebhookClause_Condition_MORE_THAN,
								},
							},
						},
					}),
				},
			},
		},
		{
			name:      "execute nothing because of conditions not matched",
			webhookId: "webhook-1",
			payload:   `{"body":{"Alert id": 123, "Status": "Open", "Impacted plays": 10}}`,
			expected:  false,
			autoOpsRule: &autoopsdomain.AutoOpsRule{
				AutoOpsRule: &autoopsproto.AutoOpsRule{
					Id: "rule-1",
					Clauses: convert([]autoopsproto.WebhookClause{
						{
							WebhookId: "webhook-1",
							Conditions: []*autoopsproto.WebhookClause_Condition{
								{
									Filter:   `.body."Alert id"`,
									Value:    `123`,
									Operator: autoopsproto.WebhookClause_Condition_EQUAL,
								},
								{
									Filter:   `.body.Status`,
									Value:    `"Close"`,
									Operator: autoopsproto.WebhookClause_Condition_EQUAL,
								},
								{
									Filter:   `.body."Impacted plays"`,
									Value:    `5`,
									Operator: autoopsproto.WebhookClause_Condition_MORE_THAN_OR_EQUAL,
								},
							},
						},
					}),
				},
			},
		},
		{
			name:      "execute nothing because of conditions not matched",
			webhookId: "webhook-1",
			payload:   `{"body":{"Alert id": 123, "Status": "Open", "Impacted plays": 10}}`,
			expected:  false,
			autoOpsRule: &autoopsdomain.AutoOpsRule{
				AutoOpsRule: &autoopsproto.AutoOpsRule{
					Id: "rule-1",
					Clauses: convert([]autoopsproto.WebhookClause{
						{
							WebhookId: "webhook-1",
							Conditions: []*autoopsproto.WebhookClause_Condition{
								{
									Filter:   `.body."Alert id"`,
									Value:    `321`,
									Operator: autoopsproto.WebhookClause_Condition_EQUAL,
								},
								{
									Filter:   `.body.Status`,
									Value:    `"Open"`,
									Operator: autoopsproto.WebhookClause_Condition_EQUAL,
								},
								{
									Filter:   `.body."Impacted plays"`,
									Value:    `5`,
									Operator: autoopsproto.WebhookClause_Condition_MORE_THAN_OR_EQUAL,
								},
							},
						},
					}),
				},
			},
		},
		{
			name:      "execute nothing because of webhook id is not matched",
			webhookId: "webhook-1",
			payload:   `{"body":{"Alert id": 123}}`,
			expected:  false,
			autoOpsRule: &autoopsdomain.AutoOpsRule{
				AutoOpsRule: &autoopsproto.AutoOpsRule{
					Id: "rule-1",
					Clauses: convert([]autoopsproto.WebhookClause{
						{
							WebhookId: "webhook-2",
							Conditions: []*autoopsproto.WebhookClause_Condition{
								{
									Filter:   `.body."Alert id"`,
									Value:    `123`,
									Operator: autoopsproto.WebhookClause_Condition_EQUAL,
								},
							},
						},
					}),
				},
			},
		},
		{
			name:      `execute nothing because of typo "Alert id" with "Alert Id"`,
			webhookId: "webhook-1",
			payload:   `{"body":{"Alert id": 123, "Status": "Open", "Impacted plays": 10}}`,
			expected:  false,
			autoOpsRule: &autoopsdomain.AutoOpsRule{
				AutoOpsRule: &autoopsproto.AutoOpsRule{
					Id: "rule-1",
					Clauses: convert([]autoopsproto.WebhookClause{
						{
							WebhookId: "webhook-1",
							Conditions: []*autoopsproto.WebhookClause_Condition{
								{
									Filter:   `.body."Alert Id"`,
									Value:    `123`,
									Operator: autoopsproto.WebhookClause_Condition_EQUAL,
								},
							},
						},
					}),
				},
			},
		},
		{
			name:      `execute nothing because of invalid filter`,
			webhookId: "webhook-1",
			payload:   `{"body":{"Alert id": 123, "Status": "Open", "Impacted plays": 10}}`,
			expected:  false,
			wantErr:   true,
			autoOpsRule: &autoopsdomain.AutoOpsRule{
				AutoOpsRule: &autoopsproto.AutoOpsRule{
					Id: "rule-1",
					Clauses: convert([]autoopsproto.WebhookClause{
						{
							WebhookId: "webhook-1",
							Conditions: []*autoopsproto.WebhookClause_Condition{
								{
									Filter:   `.body | ..`,
									Value:    `123`,
									Operator: autoopsproto.WebhookClause_Condition_EQUAL,
								},
							},
						},
					}),
				},
			},
		},
		{
			name:      `execute nothing because of invalid data`,
			webhookId: "webhook-1",
			payload:   `{"body":{"Alert id": 123, "Status": "Open", "Impacted plays": 10}}`,
			expected:  false,
			wantErr:   true,
			autoOpsRule: &autoopsdomain.AutoOpsRule{
				AutoOpsRule: &autoopsproto.AutoOpsRule{
					Id: "rule-1",
					Clauses: convert([]autoopsproto.WebhookClause{
						{
							WebhookId:  "webhook-1",
							Conditions: []*autoopsproto.WebhookClause_Condition{},
						},
					}),
				},
			},
		},
		{
			name:      `execute nothing because of converting error`,
			webhookId: "webhook-1",
			payload:   `{"body":{"foo": true}}`,
			expected:  false,
			wantErr:   true,
			autoOpsRule: &autoopsdomain.AutoOpsRule{
				AutoOpsRule: &autoopsproto.AutoOpsRule{
					Id: "rule-1",
					Clauses: convert([]autoopsproto.WebhookClause{
						{
							WebhookId: "webhook-1",
							Conditions: []*autoopsproto.WebhookClause_Condition{
								{
									Filter:   `.body.foo`,
									Value:    `123`,
									Operator: autoopsproto.WebhookClause_Condition_EQUAL,
								},
							},
						},
					}),
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var payload interface{}
			err := json.Unmarshal([]byte(c.payload), &payload)
			require.NoError(t, err)
			h := &handler{
				mysqlClient:   mysqlmock.NewMockClient(mockController),
				authClient:    authclientmock.NewMockClient(mockController),
				featureClient: featureclientmock.NewMockClient(mockController),
				logger:        logger,
			}
			result, err := h.assessAutoOpsRule(
				ctx,
				c.autoOpsRule,
				c.webhookId,
				payload,
			)
			assert.Equal(t, c.expected, result)
			assert.Equal(t, c.wantErr, err != nil)
			if err != nil {
				t.Log(err)
			}
		})
	}
}

func TestAuthWebhook(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	h := &handler{
		mysqlClient:       mysqlmock.NewMockClient(mockController),
		authClient:        authclientmock.NewMockClient(mockController),
		featureClient:     featureclientmock.NewMockClient(mockController),
		webhookCryptoUtil: &dummyWebhookCryptoUtil{},
		logger:            logger,
	}
	ctx := context.TODO()

	testcases := []struct {
		desc                 string
		id                   string
		environmentNamespace string
	}{
		{
			desc:                 "success",
			id:                   "id-1",
			environmentNamespace: "ns-1",
		},
	}
	for _, p := range testcases {
		t.Run(p.desc, func(t *testing.T) {
			ws := autoopsdomain.NewWebhookSecret(p.id, p.environmentNamespace)
			encoded, err := json.Marshal(ws)
			require.NoError(t, err)
			actual, err := h.authWebhook(ctx, base64.RawURLEncoding.EncodeToString(encoded))
			require.NoError(t, err)
			assert.Equal(t, ws, actual)
		})
	}
}
