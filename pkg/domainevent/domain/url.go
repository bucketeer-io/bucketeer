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
	"errors"
	"fmt"

	proto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

const (
	urlTemplateCodeRef      = "%s/%s/coderefs/%s"
	urlTemplateFeature      = "%s/%s/features/%s"
	urlTemplateGoal         = "%s/%s/goals/%s"
	urlTemplateExperiment   = "%s/%s/experiments/%s"
	urlTemplateAccount      = "%s/%s/accounts/%s"
	urlTemplateAPIKey       = "%s/%s/apikeys/%s"
	urlTemplateSegment      = "%s/%s/segments/%s"
	urlTemplateAutoOpsRule  = "%s/%s/features/%s/autoops"
	urlTemplatePush         = "%s/%s/pushes/%s?environmentId=%s"
	urlTemplateSubscription = "%s/%s/notifications/%s?environmentId=%s"
	urlTemplateTag          = "%s/%s/tags/%s"
	urlTemplateTeam         = "%s/%s/teams/%s"

	urlTemplateAdminSubscription = "%s/%s/notifications/%s"
	urlTemplateEnvironment       = "%s/%s/environments/%s"
	urlTemplateAdminAccount      = "%s/%s/accounts/%s"
	urlTemplateProject           = "%s/%s/projects/%s"
	urlTemplateOrganization      = "%s/organizations/%s/settings"
)

var (
	ErrUnknownEntityType = errors.New("domain: unknown entity type")
)

func URL(entityType proto.Event_EntityType, url, envURLCode, id string) (string, error) {
	switch entityType {
	case proto.Event_FEATURE:
		return fmt.Sprintf(urlTemplateFeature, url, envURLCode, id), nil
	case proto.Event_GOAL:
		return fmt.Sprintf(urlTemplateGoal, url, envURLCode, id), nil
	case proto.Event_EXPERIMENT:
		return fmt.Sprintf(urlTemplateExperiment, url, envURLCode, id), nil
	case proto.Event_ACCOUNT:
		return fmt.Sprintf(urlTemplateAccount, url, envURLCode, id), nil
	case proto.Event_APIKEY:
		return fmt.Sprintf(urlTemplateAPIKey, url, envURLCode, id), nil
	case proto.Event_SEGMENT:
		return fmt.Sprintf(urlTemplateSegment, url, envURLCode, id), nil
	case proto.Event_AUTOOPS_RULE, proto.Event_PROGRESSIVE_ROLLOUT:
		return fmt.Sprintf(urlTemplateAutoOpsRule, url, envURLCode, id), nil
	case proto.Event_PUSH:
		return fmt.Sprintf(urlTemplatePush, url, envURLCode, id, envURLCode), nil
	case proto.Event_SUBSCRIPTION:
		return fmt.Sprintf(urlTemplateSubscription, url, envURLCode, id, envURLCode), nil
	case proto.Event_ADMIN_SUBSCRIPTION:
		return fmt.Sprintf(urlTemplateAdminSubscription, url, envURLCode, id), nil
	case proto.Event_ENVIRONMENT:
		return fmt.Sprintf(urlTemplateEnvironment, url, envURLCode, id), nil
	case proto.Event_ADMIN_ACCOUNT:
		return fmt.Sprintf(urlTemplateAdminAccount, url, envURLCode, id), nil
	case proto.Event_PROJECT:
		return fmt.Sprintf(urlTemplateProject, url, envURLCode, id), nil
	case proto.Event_ORGANIZATION:
		return fmt.Sprintf(urlTemplateOrganization, url, id), nil
	case proto.Event_FLAG_TRIGGER:
		return fmt.Sprintf(urlTemplateFeature, url, envURLCode, id), nil
	case proto.Event_TAG:
		return fmt.Sprintf(urlTemplateTag, url, envURLCode, id), nil
	case proto.Event_CODEREF:
		return fmt.Sprintf(urlTemplateCodeRef, url, envURLCode, id), nil
	case proto.Event_TEAM:
		return fmt.Sprintf(urlTemplateTeam, url, envURLCode, id), nil
	}
	return "", ErrUnknownEntityType
}
