import { Trans } from 'react-i18next';
import {
  IconRemoveOutlined,
  IconUpdateOutlined
} from 'react-icons-material-design';
import { useTranslation } from 'i18n';
import { IconPlus } from '@icons';
import { DiscardChangesStateData } from 'pages/feature-flag-details/targeting';
import { DiscardChangesType } from 'pages/feature-flag-details/targeting/types';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';

interface Props {
  isOpen: boolean;
  type: DiscardChangesType | undefined;
  data: DiscardChangesStateData[];
  onClose: () => void;
  onSubmit: (type: DiscardChangesType) => void;
}

const ActionIcon = ({
  labelType
}: {
  labelType: DiscardChangesStateData['labelType'];
}) => (
  <Icon
    icon={
      labelType === 'ADD'
        ? IconPlus
        : labelType === 'REMOVE'
          ? IconRemoveOutlined
          : IconUpdateOutlined
    }
    size={'sm'}
    color="gray-600"
  />
);

const PrerequisiteDiscardItem = ({
  labelType,
  label,
  variationIndex,
  variation
}: DiscardChangesStateData) => {
  const { t } = useTranslation(['common', 'form']);
  return (
    <div className="flex flex-col w-full gap-1">
      <div className="flex w-full gap-x-2">
        <div className="mt-[3px]">
          <ActionIcon labelType={labelType} />
        </div>
        <div className="typo-para-medium text-gray-700">
          <Trans
            i18nKey={'form:prerequisite-discard-desc'}
            values={{
              action: t(labelType.toLowerCase()),
              flagName: label
            }}
            components={{
              b: <strong />
            }}
          />
        </div>
      </div>
      <div className="flex items-center gap-x-2 pl-7">
        <div className="flex-center size-fit">
          <FlagVariationPolygon index={variationIndex} />
        </div>
        <p className="typo-para-medium text-gray-700">
          {variation?.name || variation?.value}
        </p>
      </div>
    </div>
  );
};

const IndividualDiscardItem = ({
  labelType,
  label,
  variationIndex,
  variation
}: DiscardChangesStateData) => {
  const { t } = useTranslation(['common', 'form']);
  return (
    <div className="flex flex-col w-full gap-1">
      <div className="flex w-full gap-x-2">
        <div className="mt-[3px]">
          <ActionIcon labelType={labelType} />
        </div>
        <div className="typo-para-medium text-gray-700">
          <Trans
            i18nKey={'form:individual-discard-desc'}
            values={{
              action: t(labelType.toLowerCase()),
              flagName: label
            }}
            components={{
              b: <strong />,
              variantElement: (
                <div className="flex items-center gap-x-2">
                  <div className="flex-center size-fit">
                    <FlagVariationPolygon index={variationIndex} />
                  </div>
                  <p className="typo-para-medium text-gray-700">
                    {variation?.name || variation?.value}
                  </p>
                </div>
              )
            }}
          />
        </div>
      </div>
    </div>
  );
};

const DiscardChangeModal = ({
  isOpen,
  type,
  data,
  onClose,
  onSubmit
}: Props) => {
  const { t } = useTranslation(['common', 'form']);
  return (
    <DialogModal
      className="w-[500px]"
      title={t('form:discard-unsaved-changes')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full gap-y-4 p-5 max-h-[500px] overflow-y-auto small-scroll">
        {data.map((item, index) => {
          if (type === DiscardChangesType.PREREQUISITE)
            return <PrerequisiteDiscardItem key={index} {...item} />;
          if (type === DiscardChangesType.INDIVIDUAL)
            return <IndividualDiscardItem key={index} {...item} />;
        })}
      </div>
      <ButtonBar
        primaryButton={
          <Button type="button" variant="secondary" onClick={onClose}>
            {t(`cancel`)}
          </Button>
        }
        secondaryButton={
          <Button
            type="button"
            variant="negative"
            onClick={() => onSubmit(type!)}
          >
            {t(`discard`)}
          </Button>
        }
      />
    </DialogModal>
  );
};

export default DiscardChangeModal;
