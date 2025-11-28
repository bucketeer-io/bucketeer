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

package storage

import (
	"context"
	_ "embed"
	"errors"
	"time"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auth/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
)

var (
	ErrEmailVerificationTokenNotFound = errors.New("email_verification: token not found")
)

var (
	//go:embed sql/email_verification/insert_verification_token.sql
	insertVerificationTokenSQL string
	//go:embed sql/email_verification/select_verification_token.sql
	selectVerificationTokenSQL string
	//go:embed sql/email_verification/mark_verified.sql
	markVerifiedSQL string
	//go:embed sql/email_verification/delete_expired_tokens.sql
	deleteExpiredTokensSQL string
)

type emailVerificationStorage struct {
	qe mysql.QueryExecer
}

// NewEmailVerificationStorage creates a new email verification storage instance
func NewEmailVerificationStorage(qe mysql.QueryExecer) EmailVerificationStorage {
	return &emailVerificationStorage{qe: qe}
}

func (s *emailVerificationStorage) CreateVerificationToken(
	ctx context.Context,
	email, token string,
	expiresAt int64,
	ipAddress, userAgent string,
) error {
	createdAt := time.Now().Unix()
	_, err := s.qe.ExecContext(
		ctx,
		insertVerificationTokenSQL,
		email,
		token,
		createdAt,
		expiresAt,
		ipAddress,
		userAgent,
	)
	if err != nil {
		// ON DUPLICATE KEY UPDATE means no error on duplicate
		return err
	}
	return nil
}

func (s *emailVerificationStorage) GetVerificationToken(
	ctx context.Context,
	token string,
) (*domain.EmailVerificationToken, error) {
	var vToken domain.EmailVerificationToken
	var verifiedAt *int64
	err := s.qe.QueryRowContext(
		ctx,
		selectVerificationTokenSQL,
		token,
	).Scan(
		&vToken.Email,
		&vToken.Token,
		&vToken.CreatedAt,
		&vToken.ExpiresAt,
		&verifiedAt,
		&vToken.IPAddress,
		&vToken.UserAgent,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrEmailVerificationTokenNotFound
		}
		return nil, err
	}
	vToken.VerifiedAt = verifiedAt
	return &vToken, nil
}

func (s *emailVerificationStorage) MarkVerified(
	ctx context.Context,
	token string,
	verifiedAt int64,
) error {
	result, err := s.qe.ExecContext(
		ctx,
		markVerifiedSQL,
		verifiedAt,
		token,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// Token not found or already verified
		return ErrEmailVerificationTokenNotFound
	}
	return nil
}

func (s *emailVerificationStorage) DeleteExpiredTokens(
	ctx context.Context,
	before int64,
) error {
	_, err := s.qe.ExecContext(
		ctx,
		deleteExpiredTokensSQL,
		before,
		before,
	)
	return err
}
