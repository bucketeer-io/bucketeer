import { useMemo } from 'react';
import { useFormContext, useWatch } from 'react-hook-form';
import { useFeatureFlagsLoader } from 'hooks/use-feature-loading-more';
import { useTranslation } from 'i18n';
import { Experiment, Feature } from '@types';
import Dropdown from 'components/dropdown';
import Form from 'components/form';
import CreateNewOptionButton from 'elements/create-new-option-button';
import DropdownMenuWithSearch from 'elements/dropdown-with-search';
import FeatureFlagStatus from 'elements/feature-flag-status';
import VariationLabel from 'elements/variation-label';
import { ExperimentCreateUpdateForm } from '.';

interface FeatureFlagFieldProps {
  disabled: boolean;
  isEdit: boolean | string | undefined;
  environmentId: string;
  experiment?: Experiment;
  experimentFeature?: Feature;
  isOpenCreateFlagModal: boolean;
  onOpenCreateFlagModal: () => void;
}

const FeatureFlagField = ({
  disabled,
  isEdit,
  environmentId,
  experiment,
  experimentFeature,
  isOpenCreateFlagModal,
  onOpenCreateFlagModal
}: FeatureFlagFieldProps) => {
  const { t } = useTranslation(['form', 'common']);
  const { control } = useFormContext<ExperimentCreateUpdateForm>();
  const featureId = useWatch({ control, name: 'featureId' });

  const {
    allAvailableFlags,
    remainingFlagOptions,
    isLoadingMore,
    isSearching: isSearchingFeature,
    isInitialLoading: isLoadingFeature,
    hasMore,
    onSearchChange,
    loadMore
  } = useFeatureFlagsLoader({
    environmentId,
    selectedFlagIds: featureId ? [featureId] : [],
    filterSelected: !isEdit
  });

  const featureFlagOptions = useMemo(
    () =>
      allAvailableFlags.map(feature => {
        return {
          value: feature.id,
          label: feature.name,
          enabled: feature.enabled,
          disabled: featureId === feature.id
        };
      }),
    [allAvailableFlags, featureId]
  );

  const variationOptions = useMemo(() => {
    // In edit mode, use variations from fetched experiment feature or experiment data
    const variations =
      isEdit && (experimentFeature || experiment)
        ? experimentFeature?.variations || experiment?.variations
        : allAvailableFlags?.find(item => item.id === featureId)?.variations;

    return (
      variations?.map((item, index) => ({
        label: <VariationLabel label={item.name || item.value} index={index} />,
        value: item.id
      })) || []
    );
  }, [isEdit, experimentFeature, experiment, allAvailableFlags, featureId]);

  return (
    <>
      <Form.Field
        control={control}
        name={`featureId`}
        render={({ field }) => (
          <Form.Item className="flex flex-col w-full">
            <Form.Label required>{t('common:flag')}</Form.Label>
            <Form.Control>
              <DropdownMenuWithSearch
                disabled={!!isEdit || disabled}
                hidden={isOpenCreateFlagModal}
                isLoading={isLoadingFeature}
                isLoadingMore={isLoadingMore}
                isSearching={isSearchingFeature}
                isHasMore={hasMore || isLoadingMore}
                onSearchChange={onSearchChange}
                onHasMoreOptions={loadMore}
                placeholder={t(`experiments.select-flag`)}
                label={
                  (isEdit && experimentFeature
                    ? experimentFeature.name
                    : featureFlagOptions.find(
                        item => item.value === field.value
                      )?.label) || ''
                }
                options={remainingFlagOptions}
                selectedOptions={[field.value]}
                additionalElement={item => (
                  <FeatureFlagStatus
                    status={t(
                      item.enabled ? 'experiments.on' : 'experiments.off'
                    )}
                    enabled={item.enabled as boolean}
                  />
                )}
                createNewOption={
                  disabled ? undefined : (
                    <CreateNewOptionButton
                      text={t('common:create-a-new-flag')}
                      onClick={onOpenCreateFlagModal}
                    />
                  )
                }
                onSelectOption={field.onChange}
              />
            </Form.Control>
            <Form.Message />
          </Form.Item>
        )}
      />
      {featureId && (
        <Form.Field
          control={control}
          name={`baseVariationId`}
          render={({ field }) => (
            <Form.Item className="flex flex-col w-full overflow-hidden">
              <Form.Label required>
                {t('experiments.base-variation')}
              </Form.Label>
              <Form.Control>
                <Dropdown
                  disabled={!!isEdit || disabled}
                  placeholder={t(`experiments.select-variation`)}
                  className="w-full [&>div>p]:truncate [&>div]:max-w-[calc(100%-36px)]"
                  contentClassName="min-w-[502px]"
                  options={variationOptions}
                  value={field.value}
                  onChange={field.onChange}
                />
              </Form.Control>
              <Form.Message />
            </Form.Item>
          )}
        />
      )}
    </>
  );
};

export default FeatureFlagField;
