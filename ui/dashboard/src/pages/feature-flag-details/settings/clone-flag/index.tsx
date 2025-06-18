import { useToast, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import { Feature } from '@types';
import CloneFlagModal from 'pages/feature-flags/flags-modal/clone-flag-modal';
import Button from 'components/button';
import Card from 'elements/card';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';

const CloneFlag = ({
  feature,
  disabled
}: {
  feature: Feature;
  disabled: boolean;
}) => {
  const { t } = useTranslation(['common', 'form']);
  const { errorNotify } = useToast();

  const [isOpenCloneFlagModal, onOpenCloneFlagModal, onCloseCloneFlagModal] =
    useToggleOpen(false);
  return (
    <Card>
      <p className="typo-head-bold-small text-gray-800">{t('clone-flag')}</p>
      <p className="typo-para-small text-gray-500">
        {t('form:clone-flag-desc')}
      </p>
      <DisabledButtonTooltip
        align="start"
        hidden={!disabled}
        trigger={
          <Button
            className="w-fit"
            variant="secondary"
            disabled={disabled}
            onClick={onOpenCloneFlagModal}
          >
            {t('clone-flag')}
          </Button>
        }
      />
      {isOpenCloneFlagModal && !disabled && (
        <CloneFlagModal
          isOpen={isOpenCloneFlagModal}
          flagId={feature.id}
          onClose={onCloseCloneFlagModal}
          errorToast={errorNotify}
        />
      )}
    </Card>
  );
};

export default CloneFlag;
