package deleter

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/batch/jobs"
	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client/mock"
	storagemock "github.com/bucketeer-io/bucketeer/pkg/environment/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	envproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

func TestDemoOrganizationDeleter_Run(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	deleter := newMockDemoOrganizationDeleter(t, mockController)
	err := os.Setenv("DEMO_TRIAL_PERIOD_DAY", "7")
	assert.Nil(t, err)

	patterns := []struct {
		desc     string
		setup    func(deleter *demoOrganizationDeleter)
		expected error
	}{
		{
			desc: "Error internal",
			setup: func(deleter *demoOrganizationDeleter) {
				deleter.organizationStorage.(*storagemock.MockOrganizationStorage).EXPECT().ListOrganizations(
					gomock.Any(), gomock.Any(),
				).Return(nil, 0, int64(0), errors.New("internal error"))
			},
			expected: errors.New("internal error"),
		},
		{
			desc: "Success no outdated organizations",
			setup: func(deleter *demoOrganizationDeleter) {
				deleter.organizationStorage.(*storagemock.MockOrganizationStorage).EXPECT().ListOrganizations(
					gomock.Any(), gomock.Any(),
				).Return([]*envproto.Organization{}, 0, int64(0), nil)
			},
		},
		{
			desc: "Success delete outdated organizations",
			setup: func(deleter *demoOrganizationDeleter) {
				deleter.organizationStorage.(*storagemock.MockOrganizationStorage).EXPECT().ListOrganizations(
					gomock.Any(), gomock.Any(),
				).Return([]*envproto.Organization{
					{
						Id:          "org-id-1",
						CreatedAt:   17000000000,
						SystemAdmin: false,
					},
					{
						Id:          "org-id-2",
						CreatedAt:   17000000000,
						SystemAdmin: false,
					},
				}, 2, int64(2), nil)
				deleter.environmentClient.(*environmentclient.MockClient).EXPECT().DeleteBucketeerData(
					gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			p.setup(deleter)
			err := deleter.Run(ctx)
			if p.expected != nil {
				assert.Equal(t, p.expected, err)
				return
			}
		})
	}
}

func newMockDemoOrganizationDeleter(t *testing.T, c *gomock.Controller) *demoOrganizationDeleter {
	t.Helper()
	logger, err := log.NewLogger()
	assert.Nil(t, err)
	return &demoOrganizationDeleter{
		organizationStorage: storagemock.NewMockOrganizationStorage(c),
		environmentClient:   environmentclient.NewMockClient(c),
		opts: &jobs.Options{
			Timeout: 5 * time.Second,
		},
		logger: logger,
	}
}
