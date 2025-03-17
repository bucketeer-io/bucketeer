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

package v2

// func TestGetAdminAccountV2(t *testing.T) {
// 	t.Parallel()
// 	mockController := gomock.NewController(t)
// 	defer mockController.Finish()
// 	patterns := []struct {
// 		desc        string
// 		setup       func(*accountStorage)
// 		email       string
// 		expectedErr error
// 	}{
// 		{
// 			desc: "ErrAdminAccountNotFound",
// 			setup: func(s *accountStorage) {
// 				qe := mock.NewMockQueryExecer(mockController)
// 				s.client.(*mock.MockClient).EXPECT().Qe(
// 					gomock.Any(),
// 				).Return(qe)
// 				row := mock.NewMockRow(mockController)
// 				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
// 				qe.EXPECT().QueryRowContext(
// 					gomock.Any(), gomock.Any(), gomock.Any(),
// 				).Return(row)
// 			},
// 			email:       "bucketeer@example.com",
// 			expectedErr: ErrSystemAdminAccountNotFound,
// 		},
// 		{
// 			desc: "Error",
// 			setup: func(s *accountStorage) {
// 				qe := mock.NewMockQueryExecer(mockController)
// 				s.client.(*mock.MockClient).EXPECT().Qe(
// 					gomock.Any(),
// 				).Return(qe)
// 				row := mock.NewMockRow(mockController)
// 				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
// 				qe.EXPECT().QueryRowContext(
// 					gomock.Any(), gomock.Any(), gomock.Any(),
// 				).Return(row)

// 			},
// 			email:       "bucketeer@example.com",
// 			expectedErr: errors.New("error"),
// 		},
// 		{
// 			desc: "Success",
// 			setup: func(s *accountStorage) {
// 				qe := mock.NewMockQueryExecer(mockController)
// 				s.client.(*mock.MockClient).EXPECT().Qe(
// 					gomock.Any(),
// 				).Return(qe)
// 				row := mock.NewMockRow(mockController)
// 				row.EXPECT().Scan(gomock.Any()).Return(nil)
// 				qe.EXPECT().QueryRowContext(
// 					gomock.Any(), gomock.Any(), gomock.Any(),
// 				).Return(row)
// 			},
// 			email:       "bucketeer@example.com",
// 			expectedErr: nil,
// 		},
// 	}
// 	for _, p := range patterns {
// 		t.Run(p.desc, func(t *testing.T) {
// 			storage := newAccountStorageWithMock(t, mockController)
// 			if p.setup != nil {
// 				p.setup(storage)
// 			}
// 			_, err := storage.GetSystemAdminAccountV2(context.Background(), p.email)
// 			assert.Equal(t, p.expectedErr, err)
// 		})
// 	}
// }
