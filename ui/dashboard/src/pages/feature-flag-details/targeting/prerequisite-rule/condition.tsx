import { forwardRef, Ref, useEffect, useMemo } from 'react';
import { useFormContext, useWatch } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { Feature } from '@types';
import { truncateBySide } from 'utils/converts';
import { cn } from 'utils/style';
import { IconTrash } from '@icons';
import Button from 'components/button';
import Dropdown from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import DropdownMenuWithSearch from 'elements/dropdown-with-search';
import FeatureFlagStatus from 'elements/feature-flag-status';
import VariationLabel from 'elements/variation-label';
import { PrerequisiteSchema } from '../types';

interface Props {
  features: Feature[];
  type: 'if' | 'and';
  prerequisiteIndex: number;
  featureId: string;
  onDeleteCondition: () => void;
}

const ConditionForm = forwardRef(
  (
    { features, type, prerequisiteIndex, featureId, onDeleteCondition }: Props,
    ref: Ref<HTMLDivElement>
  ) => {
    const { t } = useTranslation(['form', 'common', 'table']);

    const methods = useFormContext();
    const { control } = methods;

    const prerequisitesWatch = useWatch({ control, name: 'prerequisites' });

    const commonName = useMemo(
      () => `prerequisites.${prerequisiteIndex}`,
      [prerequisiteIndex]
    );
    const currentFeatureId = useWatch({
      control,
      name: `${commonName}.featureId`
    });
    const currentVariationId = useWatch({
      control,
      name: `${commonName}.variationId`
    });
    const currentFeature = useMemo(
      () => features.find(item => item.id === currentFeatureId),
      [currentFeatureId, features]
    );

    const flagOptions = useMemo(() => {
      const featuresSelected = prerequisitesWatch.map(
        (item: PrerequisiteSchema) => item.featureId
      );
      return features
        .filter(f => ![...featuresSelected, featureId].includes(f.id))
        .map(item => ({
          label: item.name,
          value: item.id,
          enabled: item.enabled
        }));
    }, [features, prerequisitesWatch, featureId]);

    const variationOptions = useMemo(
      () =>
        currentFeature?.variations?.map((item, index) => ({
          label: (
            <VariationLabel label={item.name || item.value} index={index} />
          ),
          value: item.id
        })),
      [currentFeature]
    );
    useEffect(() => {
      const firstVariationId = currentFeature?.variations?.[0]?.id;

      const isCurrentValid = currentFeature?.variations?.some(
        v => v.id === currentVariationId
      );

      if ((!currentVariationId || !isCurrentValid) && firstVariationId) {
        methods.setValue(`${commonName}.variationId`, firstVariationId, {
          shouldDirty: true,
          shouldValidate: true
        });
      }
    }, [currentFeatureId, currentFeature, methods.setValue, commonName]);

    return (
      <div ref={ref} className="flex items-center w-full gap-x-4">
        <div
          className={cn(
            'flex-center w-[42px] h-[26px] rounded-[3px] typo-para-small leading-[14px]',
            {
              'bg-accent-pink-50 text-accent-pink-500': type === 'if',
              'bg-gray-200 text-gray-600': type === 'and'
            }
          )}
        >
          {type === 'if' ? t('common:if') : t('common:and')}
        </div>
        <div className="flex items-center w-full flex-1 pl-4 border-l border-primary-500">
          <div className="flex w-full gap-x-4">
            <Form.Field
              control={control}
              name={`${commonName}.featureId`}
              render={({ field }) => (
                <Form.Item className="flex flex-col flex-1 self-stretch py-0 min-w-[170px]">
                  <Form.Label required>{t('feature-flags.flag')}</Form.Label>
                  <Form.Control>
                    <DropdownMenuWithSearch
                      label={truncateBySide(currentFeature?.name || '', 50)}
                      placeholder={t('experiments.select-flag')}
                      isExpand
                      options={flagOptions}
                      selectedOptions={field.value}
                      additionalElement={item => (
                        <FeatureFlagStatus
                          status={t(
                            item.enabled ? 'experiments.on' : 'experiments.off'
                          )}
                          enabled={item.enabled as boolean}
                        />
                      )}
                      onSelectOption={value => {
                        field.onChange(value);
                      }}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <div className="flex-center size-fit min-w-fit px-3 py-3.5 mt-6 rounded bg-gray-100 typo-para-medium leading-5 text-gray-700 whitespace-nowrap">
              <Trans
                i18nKey={'form:feature-flags.receiving-state'}
                values={{
                  state: t('form:experiments.on').toUpperCase()
                }}
              />
            </div>
            <Form.Field
              control={control}
              name={`${commonName}.variationId`}
              render={({ field }) => (
                <Form.Item className="flex flex-col flex-1 self-stretch py-0 min-w-[170px]">
                  <Form.Label required>
                    {t('table:feature-flags.variation')}
                  </Form.Label>
                  <Form.Control>
                    <Dropdown
                      options={variationOptions}
                      value={currentVariationId}
                      labelCustom={
                        variationOptions?.find(
                          item => currentVariationId === item.value
                        )?.label || variationOptions?.[0]?.label
                      }
                      onChange={field.onChange}
                      placeholder={t('experiments.select-variation')}
                      disabled={!variationOptions?.length}
                      className="w-full"
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />

            <div className="flex items-end self-stretch">
              <Button
                type="button"
                variant={'grey'}
                className="flex-center text-gray-500 hover:text-gray-600"
                onClick={onDeleteCondition}
              >
                <Icon icon={IconTrash} size={'sm'} />
              </Button>
            </div>
          </div>
        </div>
      </div>
    );
  }
);

export default ConditionForm;
