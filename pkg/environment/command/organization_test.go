package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

func TestHandleCreateOrganizationCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	organization, err := domain.NewOrganization("organization-name", "organization-code", "organization desc", false, false)
	assert.NoError(t, err)

	h := newOrganizationCommandHandler(t, publisher, organization)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	cmd := &environmentproto.CreateOrganizationCommand{Name: organization.Name, Description: organization.Description}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
}

func TestHandleChangeNameOrganizationCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	organization, err := domain.NewOrganization("organization-name", "organization-code", "organization desc", false, false)
	assert.NoError(t, err)

	h := newOrganizationCommandHandler(t, publisher, organization)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	newName := "new-organization-name"
	cmd := &environmentproto.ChangeNameOrganizationCommand{Name: newName}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, newName, organization.Name)
}

func TestHandleChangeDescriptionOrganizationCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	organization, err := domain.NewOrganization("organization-name", "organization-code", "organization desc", false, false)
	assert.NoError(t, err)

	h := newOrganizationCommandHandler(t, publisher, organization)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	newDesc := "new organization desc"
	cmd := &environmentproto.ChangeDescriptionOrganizationCommand{Description: newDesc}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, newDesc, organization.Description)
}

func TestHandleEnableOrganizationCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	organization, err := domain.NewOrganization("organization-name", "organization-code", "organization desc", false, false)
	assert.NoError(t, err)
	organization.Disabled = true

	h := newOrganizationCommandHandler(t, publisher, organization)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	cmd := &environmentproto.EnableOrganizationCommand{}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.False(t, organization.Disabled)
}

func TestHandleDisableOrganizationCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	organization, err := domain.NewOrganization("organization-name", "organization-code", "organization desc", false, false)
	assert.NoError(t, err)

	h := newOrganizationCommandHandler(t, publisher, organization)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	cmd := &environmentproto.DisableOrganizationCommand{}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.True(t, organization.Disabled)
}

func TestHandleArchiveOrganizationCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	organization, err := domain.NewOrganization("organization-name", "organization-code", "organization desc", false, false)
	assert.NoError(t, err)

	h := newOrganizationCommandHandler(t, publisher, organization)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	cmd := &environmentproto.ArchiveOrganizationCommand{}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.True(t, organization.Archived)
}

func TestHandleUnarchiveOrganizationCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	organization, err := domain.NewOrganization("organization-name", "organization-code", "organization desc", false, false)
	assert.NoError(t, err)
	organization.Archive()

	h := newOrganizationCommandHandler(t, publisher, organization)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	cmd := &environmentproto.UnarchiveOrganizationCommand{}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.False(t, organization.Archived)
}

func TestHandleConvertTrialOrganizationCommand(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	publisher := publishermock.NewMockPublisher(mockController)
	organization, err := domain.NewOrganization("organization-name", "organization-code", "organization desc", true, false)
	assert.NoError(t, err)

	h := newOrganizationCommandHandler(t, publisher, organization)
	publisher.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
	cmd := &environmentproto.ConvertTrialOrganizationCommand{}
	err = h.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.False(t, organization.Trial)
}

func newOrganizationCommandHandler(t *testing.T, publisher publisher.Publisher, organization *domain.Organization) Handler {
	t.Helper()
	return NewOrganizationCommandHandler(
		&eventproto.Editor{
			Email: "email",
		},
		organization,
		publisher,
	)
}
