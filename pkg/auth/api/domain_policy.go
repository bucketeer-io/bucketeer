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

package api

import (
	"context"
	"errors"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
	authdomain "github.com/bucketeer-io/bucketeer/v2/pkg/auth/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/auth/storage"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	authproto "github.com/bucketeer-io/bucketeer/v2/proto/auth"
)

func (s *authService) GetAuthOptionsByEmail(
	ctx context.Context,
	request *authproto.GetAuthOptionsByEmailRequest,
) (*authproto.GetAuthOptionsByEmailResponse, error) {
	err := validateGetAuthOptionsByEmailRequest(request)
	if err != nil {
		s.logger.Error("GetAuthOptionsByEmail request validation failed", zap.Error(err))
		return nil, err
	}

	// Normalize email and extract domain
	normalizedEmail, err := authdomain.NormalizeEmail(request.Email)
	if err != nil {
		s.logger.Error("Failed to normalize email",
			zap.Error(err),
			zap.String("email", request.Email),
		)
		return nil, auth.StatusInvalidArguments.Err()
	}

	domain, err := authdomain.ExtractDomain(normalizedEmail)
	if err != nil {
		s.logger.Error("Failed to extract domain",
			zap.Error(err),
			zap.String("email", normalizedEmail),
		)
		return nil, auth.StatusInvalidArguments.Err()
	}

	// Lookup domain policy
	policy, err := s.domainPolicyStorage.GetDomainPolicy(ctx, domain)
	if err != nil && !errors.Is(err, storage.ErrDomainPolicyNotFound) {
		s.logger.Error("Failed to get domain policy",
			zap.Error(err),
			zap.String("domain", domain),
		)
		return nil, auth.StatusInternal.Err()
	}

	// Build auth options response
	options := buildAuthOptions(domain, policy, s.config)

	return &authproto.GetAuthOptionsByEmailResponse{
		Options: options,
	}, nil
}

func (s *authService) CreateDomainAuthPolicy(
	ctx context.Context,
	request *authproto.CreateDomainAuthPolicyRequest,
) (*authproto.CreateDomainAuthPolicyResponse, error) {
	// Check system admin permission
	if err := s.checkSystemAdminPermission(ctx); err != nil {
		return nil, err
	}

	err := validateCreateDomainAuthPolicyRequest(request)
	if err != nil {
		s.logger.Error("CreateDomainAuthPolicy request validation failed", zap.Error(err))
		return nil, err
	}

	now := time.Now().Unix()
	policy := &authdomain.DomainAuthPolicy{
		Domain:     request.Domain,
		AuthPolicy: request.AuthPolicy,
		Enabled:    true, // New policies are enabled by default
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	err = s.domainPolicyStorage.CreateDomainPolicy(ctx, policy)
	if err != nil {
		if errors.Is(err, storage.ErrDomainPolicyAlreadyExists) {
			s.logger.Error("Domain policy already exists", zap.String("domain", request.Domain))
			return nil, auth.StatusAlreadyExists.Err()
		}
		s.logger.Error("Failed to create domain policy",
			zap.Error(err),
			zap.String("domain", request.Domain),
		)
		return nil, auth.StatusInternal.Err()
	}

	s.logger.Info("Domain policy created successfully", zap.String("domain", request.Domain))
	return &authproto.CreateDomainAuthPolicyResponse{
		Policy: policy.ToProto(),
	}, nil
}

func (s *authService) UpdateDomainAuthPolicy(
	ctx context.Context,
	request *authproto.UpdateDomainAuthPolicyRequest,
) (*authproto.UpdateDomainAuthPolicyResponse, error) {
	// Check system admin permission
	if err := s.checkSystemAdminPermission(ctx); err != nil {
		return nil, err
	}

	err := validateUpdateDomainAuthPolicyRequest(request)
	if err != nil {
		s.logger.Error("UpdateDomainAuthPolicy request validation failed", zap.Error(err))
		return nil, err
	}

	// Get existing policy
	existingPolicy, err := s.domainPolicyStorage.GetDomainPolicy(ctx, request.Domain)
	if err != nil {
		if errors.Is(err, storage.ErrDomainPolicyNotFound) {
			s.logger.Error("Domain policy not found", zap.String("domain", request.Domain))
			return nil, auth.StatusNotFound.Err()
		}
		s.logger.Error("Failed to get domain policy",
			zap.Error(err),
			zap.String("domain", request.Domain),
		)
		return nil, auth.StatusInternal.Err()
	}

	// Update policy
	existingPolicy.AuthPolicy = request.AuthPolicy
	existingPolicy.Enabled = request.Enabled
	existingPolicy.UpdatedAt = time.Now().Unix()

	err = s.domainPolicyStorage.UpdateDomainPolicy(ctx, existingPolicy)
	if err != nil {
		s.logger.Error("Failed to update domain policy",
			zap.Error(err),
			zap.String("domain", request.Domain),
		)
		return nil, auth.StatusInternal.Err()
	}

	s.logger.Info("Domain policy updated successfully", zap.String("domain", request.Domain))
	return &authproto.UpdateDomainAuthPolicyResponse{
		Policy: existingPolicy.ToProto(),
	}, nil
}

func (s *authService) GetDomainAuthPolicy(
	ctx context.Context,
	request *authproto.GetDomainAuthPolicyRequest,
) (*authproto.GetDomainAuthPolicyResponse, error) {
	// Check system admin permission
	if err := s.checkSystemAdminPermission(ctx); err != nil {
		return nil, err
	}

	err := validateGetDomainAuthPolicyRequest(request)
	if err != nil {
		s.logger.Error("GetDomainAuthPolicy request validation failed", zap.Error(err))
		return nil, err
	}

	policy, err := s.domainPolicyStorage.GetDomainPolicy(ctx, request.Domain)
	if err != nil {
		if errors.Is(err, storage.ErrDomainPolicyNotFound) {
			s.logger.Error("Domain policy not found", zap.String("domain", request.Domain))
			return nil, auth.StatusNotFound.Err()
		}
		s.logger.Error("Failed to get domain policy",
			zap.Error(err),
			zap.String("domain", request.Domain),
		)
		return nil, auth.StatusInternal.Err()
	}

	return &authproto.GetDomainAuthPolicyResponse{
		Policy: policy.ToProto(),
	}, nil
}

func (s *authService) DeleteDomainAuthPolicy(
	ctx context.Context,
	request *authproto.DeleteDomainAuthPolicyRequest,
) (*authproto.DeleteDomainAuthPolicyResponse, error) {
	// Check system admin permission
	if err := s.checkSystemAdminPermission(ctx); err != nil {
		return nil, err
	}

	err := validateDeleteDomainAuthPolicyRequest(request)
	if err != nil {
		s.logger.Error("DeleteDomainAuthPolicy request validation failed", zap.Error(err))
		return nil, err
	}

	err = s.domainPolicyStorage.DeleteDomainPolicy(ctx, request.Domain)
	if err != nil {
		if errors.Is(err, storage.ErrDomainPolicyUnexpectedAffectedRows) {
			s.logger.Error("Domain policy not found", zap.String("domain", request.Domain))
			return nil, auth.StatusNotFound.Err()
		}
		s.logger.Error("Failed to delete domain policy",
			zap.Error(err),
			zap.String("domain", request.Domain),
		)
		return nil, auth.StatusInternal.Err()
	}

	s.logger.Info("Domain policy deleted successfully", zap.String("domain", request.Domain))
	return &authproto.DeleteDomainAuthPolicyResponse{}, nil
}

func (s *authService) ListDomainAuthPolicies(
	ctx context.Context,
	request *authproto.ListDomainAuthPoliciesRequest,
) (*authproto.ListDomainAuthPoliciesResponse, error) {
	// Check system admin permission
	if err := s.checkSystemAdminPermission(ctx); err != nil {
		return nil, err
	}

	err := validateListDomainAuthPoliciesRequest(request)
	if err != nil {
		s.logger.Error("ListDomainAuthPolicies request validation failed", zap.Error(err))
		return nil, err
	}

	// Parse cursor as offset
	offset := 0
	if request.Cursor != "" {
		parsedOffset, err := strconv.Atoi(request.Cursor)
		if err != nil {
			s.logger.Error("Failed to parse cursor", zap.Error(err), zap.String("cursor", request.Cursor))
			return nil, auth.StatusInvalidArguments.Err()
		}
		offset = parsedOffset
	}

	// Default page size
	limit := int(request.PageSize)
	if limit == 0 {
		limit = 50
	}

	// Construct orders
	var orders []*mysql.Order
	if request.OrderBy != "" {
		direction := mysql.OrderDirectionAsc
		if request.OrderDirection == "DESC" || request.OrderDirection == "desc" {
			direction = mysql.OrderDirectionDesc
		}
		orders = append(orders, &mysql.Order{
			Column:    request.OrderBy,
			Direction: direction,
		})
	} else {
		// Default order by created_at DESC
		orders = append(orders, &mysql.Order{
			Column:    "created_at",
			Direction: mysql.OrderDirectionDesc,
		})
	}

	options := &mysql.ListOptions{
		Limit:  limit,
		Offset: offset,
		Orders: orders,
	}

	policies, nextOffset, totalCount, err := s.domainPolicyStorage.ListDomainPolicies(ctx, options)
	if err != nil {
		s.logger.Error("Failed to list domain policies", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	return &authproto.ListDomainAuthPoliciesResponse{
		Policies:   policies,
		Cursor:     strconv.Itoa(nextOffset),
		TotalCount: totalCount,
	}, nil
}

// Helper functions

func (s *authService) checkSystemAdminPermission(ctx context.Context) error {
	accessToken, ok := rpc.GetAccessToken(ctx)
	if !ok || accessToken == nil {
		s.logger.Error("No access token in context")
		return auth.StatusUnauthenticated.Err()
	}

	if !accessToken.IsSystemAdmin {
		s.logger.Error("Permission denied: not a system admin",
			zap.String("email", accessToken.Email),
		)
		return auth.StatusPermissionDenied.Err()
	}

	return nil
}

func buildAuthOptions(
	domain string,
	policy *authdomain.DomainAuthPolicy,
	config *auth.OAuthConfig,
) *authproto.DomainAuthOptions {
	// If no policy exists, return global defaults
	if policy == nil || !policy.Enabled {
		return &authproto.DomainAuthOptions{
			Domain:                domain,
			PasswordEnabled:       false,                 // Password not enabled by default
			GoogleOidcEnabled:     true,                  // Google OIDC is global default
			GoogleOidcDisplayName: "Sign in with Google", // Default display name
			OidcRequired:          false,
		}
	}

	// If company OIDC is required, only show company OIDC
	if policy.IsCompanyOidcRequired() {
		return &authproto.DomainAuthOptions{
			Domain:                 domain,
			PasswordEnabled:        false,
			GoogleOidcEnabled:      false,
			CompanyOidcEnabled:     true,
			CompanyOidcDisplayName: policy.AuthPolicy.CompanyOidc.DisplayName,
			OidcRequired:           true,
		}
	}

	// Otherwise, show all enabled options
	options := &authproto.DomainAuthOptions{
		Domain:          domain,
		PasswordEnabled: policy.IsPasswordEnabled(),
		OidcRequired:    false,
	}

	// Add Google OIDC if enabled
	if policy.IsGoogleOidcEnabled() {
		options.GoogleOidcEnabled = true
		if policy.AuthPolicy.GoogleOidc != nil && policy.AuthPolicy.GoogleOidc.DisplayName != "" {
			options.GoogleOidcDisplayName = policy.AuthPolicy.GoogleOidc.DisplayName
		} else {
			options.GoogleOidcDisplayName = "Sign in with Google" // Default display name
		}
	}

	// Add Company OIDC if enabled
	if policy.IsCompanyOidcEnabled() {
		options.CompanyOidcEnabled = true
		options.CompanyOidcDisplayName = policy.AuthPolicy.CompanyOidc.DisplayName
	}

	return options
}
