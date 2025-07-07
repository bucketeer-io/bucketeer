import { useCallback } from 'react';
import { Trans } from 'react-i18next';
import { featureUpdater } from '@api/features';
import { invalidateFeature } from '@queries/feature-details';
import { invalidateFeatures } from '@queries/features';
import { invalidateHistories } from '@queries/histories';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import { Feature, FeatureUpdaterParams } from '@types';
import { getFlagStatus } from 'pages/feature-flags/collection-layout/elements/utils';
import ArchiveModal from 'pages/feature-flags/flags-modal/archive-modal';
import { FeatureActivityStatus } from 'pages/feature-flags/types';
import Button from 'components/button';
import Card from 'elements/card';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';

const ArchiveFlag = ({
  feature,
  disabled
}: {
  feature: Feature;
  disabled: boolean;
}) => {
  const { t } = useTranslation(['common', 'form', 'table', 'message']);
  const { notify, errorNotify } = useToast();
  const queryClient = useQueryClient();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const [
    isOpenArchiveFlagModal,
    onOpenArchiveFlagModal,
    onCloseArchiveFlagModal
  ] = useToggleOpen(false);

  const mutation = useMutation({
    mutationFn: async (params: Partial<FeatureUpdaterParams>) => {
      return featureUpdater(params);
    },
    onSuccess: () => {
      onCloseArchiveFlagModal();
      notify({
        message: t('message:collection-action-success', {
          collection: t('common:source-type.feature-flag'),
          action: t('updated')
        })
      });
      invalidateFeature(queryClient);
      invalidateFeatures(queryClient);
      invalidateHistories(queryClient);
      mutation.reset();
    },
    onError: error => errorNotify(error)
  });

  const onSubmit = useCallback(
    async ({ comment }: { comment?: string }) => {
      mutation.mutate({
        id: feature.id,
        environmentId: currentEnvironment.id,
        archived: !feature.archived,
        comment
      });
    },
    [currentEnvironment, feature]
  );

  return (
    <Card>
      <p className="typo-head-bold-small text-gray-800">
        {t(feature.archived ? 'unarchive-flag' : 'archive-flag')}
      </p>
      <p className="typo-para-small text-gray-500">
        {t(
          feature.archived
            ? 'form:unarchive-flag-desc'
            : 'form:archive-flag-desc'
        )}
      </p>
      <DisabledButtonTooltip
        align="start"
        hidden={!disabled}
        trigger={
          <Button
            disabled={disabled}
            className="w-fit"
            variant="secondary"
            onClick={onOpenArchiveFlagModal}
          >
            {t(feature.archived ? 'unarchive-flag' : 'archive-flag')}
          </Button>
        }
      />
      {isOpenArchiveFlagModal && (
        <ArchiveModal
          isOpen={isOpenArchiveFlagModal}
          isArchiving={!feature.archived}
          isLoading={mutation.isPending}
          isShowWarning={
            !feature.archived &&
            getFlagStatus(feature) === FeatureActivityStatus.RECEIVING_TRAFFIC
          }
          title={
            !feature.archived
              ? t(`table:popover.archive-flag`)
              : t(`table:popover.unarchive-flag`)
          }
          description={
            <Trans
              i18nKey={
                !feature.archived
                  ? 'table:feature-flags.confirm-archive-desc'
                  : 'table:feature-flags.confirm-unarchive-desc'
              }
              values={{ name: feature.name }}
              components={{ text: <span /> }}
            />
          }
          onClose={onCloseArchiveFlagModal}
          onSubmit={onSubmit}
        />
      )}
    </Card>
  );
};

export default ArchiveFlag;
