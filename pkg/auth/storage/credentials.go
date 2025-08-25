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

	"github.com/bucketeer-io/bucketeer/pkg/auth/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
)

var (
	ErrCredentialsAlreadyExists          = errors.New("credentials: already exists")
	ErrCredentialsNotFound               = errors.New("credentials: not found")
	ErrCredentialsUnexpectedAffectedRows = errors.New("credentials: unexpected affected rows")
	ErrPasswordResetTokenNotFound        = errors.New("credentials: password reset token not found")
	ErrPasswordResetTokenAlreadyExists   = errors.New("credentials: password reset token already exists")
)

var (
	//go:embed sql/credentials/insert_credentials.sql
	insertCredentialsSQL string
	//go:embed sql/credentials/select_credentials.sql
	selectCredentialsSQL string
	//go:embed sql/credentials/update_password.sql
	updatePasswordSQL string
	//go:embed sql/credentials/delete_credentials.sql
	deleteCredentialsSQL string
	//go:embed sql/credentials/set_password_reset_token.sql
	setPasswordResetTokenSQL string
	//go:embed sql/credentials/get_password_reset_token.sql
	getPasswordResetTokenSQL string
	//go:embed sql/credentials/delete_password_reset_token.sql
	deletePasswordResetTokenSQL string
)

type credentialsStorage struct {
	qe mysql.QueryExecer
}

// NewCredentialsStorage creates a new credentials storage instance
func NewCredentialsStorage(qe mysql.QueryExecer) CredentialsStorage {
	return &credentialsStorage{qe: qe}
}

func (s *credentialsStorage) CreateCredentials(ctx context.Context, email, passwordHash string) error {
	now := time.Now().Unix()
	_, err := s.qe.ExecContext(
		ctx,
		insertCredentialsSQL,
		email,
		passwordHash,
		nil, // password_reset_token (NULL)
		nil, // password_reset_token_expires_at (NULL)
		now,
		now,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrDuplicateEntry) {
			return ErrCredentialsAlreadyExists
		}
		return err
	}
	return nil
}

func (s *credentialsStorage) GetCredentials(ctx context.Context, email string) (*domain.AccountCredentials, error) {
	var credentials domain.AccountCredentials
	err := s.qe.QueryRowContext(
		ctx,
		selectCredentialsSQL,
		email,
	).Scan(
		&credentials.Email,
		&credentials.PasswordHash,
		&credentials.CreatedAt,
		&credentials.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrCredentialsNotFound
		}
		return nil, err
	}
	return &credentials, nil
}

func (s *credentialsStorage) UpdatePassword(ctx context.Context, email, passwordHash string) error {
	now := time.Now().Unix()
	result, err := s.qe.ExecContext(
		ctx,
		updatePasswordSQL,
		passwordHash,
		now,
		email,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrCredentialsUnexpectedAffectedRows
	}
	return nil
}

func (s *credentialsStorage) DeleteCredentials(ctx context.Context, email string) error {
	result, err := s.qe.ExecContext(
		ctx,
		deleteCredentialsSQL,
		email,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrCredentialsUnexpectedAffectedRows
	}
	return nil
}

func (s *credentialsStorage) SetPasswordResetToken(ctx context.Context, email, token string, expiresAt int64) error {
	now := time.Now().Unix()
	result, err := s.qe.ExecContext(
		ctx,
		setPasswordResetTokenSQL,
		token,
		expiresAt,
		now,
		email,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrCredentialsUnexpectedAffectedRows
	}
	return nil
}

func (s *credentialsStorage) GetPasswordResetToken(
	ctx context.Context, token string,
) (*domain.PasswordResetToken, error) {
	var resetToken domain.PasswordResetToken
	err := s.qe.QueryRowContext(
		ctx,
		getPasswordResetTokenSQL,
		token,
	).Scan(
		&resetToken.Token,
		&resetToken.Email,
		&resetToken.ExpiresAt,
		&resetToken.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRows) {
			return nil, ErrPasswordResetTokenNotFound
		}
		return nil, err
	}
	return &resetToken, nil
}

func (s *credentialsStorage) DeletePasswordResetToken(ctx context.Context, token string) error {
	now := time.Now().Unix()
	result, err := s.qe.ExecContext(
		ctx,
		deletePasswordResetTokenSQL,
		now,
		token,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrCredentialsUnexpectedAffectedRows
	}
	return nil
}
