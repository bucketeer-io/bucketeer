package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOrganization(t *testing.T) {
	t.Parallel()
	organization, err := NewOrganization(
		"organization-name",
		"organization-code",
		"organization desc",
		false, false,
	)
	assert.NoError(t, err)
	assert.IsType(t, &Organization{}, organization)
	assert.NotEqual(t, "organization-name", organization.Id)
	assert.Equal(t, "organization-name", organization.Name)
	assert.Equal(t, "organization-code", organization.UrlCode)
	assert.Equal(t, false, organization.Trial)
}

func TestNewTrialOrganization(t *testing.T) {
	t.Parallel()
	organization, err := NewOrganization(
		"organization-name",
		"organization-code",
		"organization desc",
		true, false,
	)
	assert.NoError(t, err)
	assert.IsType(t, &Organization{}, organization)
	assert.NotEqual(t, "organization-name", organization.Id)
	assert.Equal(t, "organization-name", organization.Name)
	assert.Equal(t, "organization-code", organization.UrlCode)
	assert.Equal(t, true, organization.Trial)
}

func TestChangeDescriptionOrganization(t *testing.T) {
	t.Parallel()
	organization, err := NewOrganization(
		"organization-name",
		"organization-code",
		"organization desc",
		false, false,
	)
	assert.NoError(t, err)
	newDesc := "new org desc"
	organization.ChangeDescription(newDesc)
	assert.Equal(t, newDesc, organization.Description)
}

func TestChangeNameOrganization(t *testing.T) {
	t.Parallel()
	organization, err := NewOrganization(
		"organization-name",
		"organization-code",
		"organization desc",
		false, false,
	)
	assert.NoError(t, err)
	newName := "new-organization-name"
	organization.ChangeName(newName)
	assert.Equal(t, newName, organization.Name)
}

func TestEnableOrganization(t *testing.T) {
	t.Parallel()
	organization, err := NewOrganization(
		"organization-name",
		"organization-code",
		"organization desc",
		false, false,
	)
	assert.NoError(t, err)
	organization.Disabled = true
	organization.Enable()
	assert.False(t, organization.Disabled)
}

func TestDisableOrganization(t *testing.T) {
	t.Parallel()
	organization, err := NewOrganization(
		"organization-name",
		"organization-code",
		"organization desc",
		false, false,
	)
	assert.NoError(t, err)
	err = organization.Disable()
	assert.NoError(t, err)
	assert.True(t, organization.Disabled)
}

func TestCannotDisableSystemAdminOrganization(t *testing.T) {
	t.Parallel()
	organization, err := NewOrganization(
		"organization-name",
		"organization-code",
		"organization desc",
		false, true,
	)
	assert.NoError(t, err)
	err = organization.Disable()
	assert.Equal(t, ErrCannotDisableSystemAdmin, err)
}

func TestArchiveOrganization(t *testing.T) {
	t.Parallel()
	organization, err := NewOrganization(
		"organization-name",
		"organization-code",
		"organization desc",
		false, false,
	)
	assert.NoError(t, err)
	err = organization.Archive()
	assert.NoError(t, err)
	assert.True(t, organization.Archived)
}

func TestCannotArchiveSystemAdminOrganization(t *testing.T) {
	t.Parallel()
	organization, err := NewOrganization(
		"organization-name",
		"organization-code",
		"organization desc",
		false, true,
	)
	assert.NoError(t, err)
	err = organization.Archive()
	assert.Equal(t, ErrCannotArchiveSystemAdmin, err)
}

func TestUnarchiveOrganization(t *testing.T) {
	t.Parallel()
	organization, err := NewOrganization(
		"organization-name",
		"organization-code",
		"organization desc",
		false, false,
	)
	assert.NoError(t, err)
	organization.Archived = true
	organization.Unarchive()
	assert.False(t, organization.Archived)
}

func TestConvertTrialOrganization(t *testing.T) {
	t.Parallel()
	organization, err := NewOrganization(
		"organization-name",
		"organization-code",
		"organization desc",
		false, false,
	)
	assert.NoError(t, err)
	organization.ConvertTrial()
	assert.False(t, organization.Trial)
}
