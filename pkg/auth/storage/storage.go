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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package storage

import (
	"context"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auth/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
)

// CredentialsStorage defines the interface for managing account credentials
type CredentialsStorage interface {
	CreateCredentials(ctx context.Context, email, passwordHash string) error
	GetCredentials(ctx context.Context, email string) (*domain.AccountCredentials, error)
	UpdatePassword(ctx context.Context, email, passwordHash string) error
	DeleteCredentials(ctx context.Context, email string) error

	SetPasswordResetToken(ctx context.Context, email, token string, expiresAt int64) error
	GetPasswordResetToken(ctx context.Context, token string) (*domain.PasswordResetToken, error)
	DeletePasswordResetToken(ctx context.Context, token string) error
}

// DomainPolicyStorage defines the interface for managing domain authentication policies
type DomainPolicyStorage interface {
	CreateDomainPolicy(ctx context.Context, policy *domain.DomainAuthPolicy) error
	GetDomainPolicy(ctx context.Context, domainName string) (*domain.DomainAuthPolicy, error)
	UpdateDomainPolicy(ctx context.Context, policy *domain.DomainAuthPolicy) error
	DeleteDomainPolicy(ctx context.Context, domainName string) error
	ListDomainPolicies(ctx context.Context, options *mysql.ListOptions) ([]*authproto.DomainAuthPolicy, int, int64, error)
}
