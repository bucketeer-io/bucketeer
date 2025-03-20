import { forwardRef, Ref, useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { AnyObject } from 'yup';
import { Feature } from '@types';
import { cn } from 'utils/style';
import { IconTrash } from '@icons';
import Button from 'components/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import { PrerequisiteRuleType } from '../types';

interface Props {
  features: Feature[];
  type: 'if' | 'and';
  condition: PrerequisiteRuleType;
  isDisabledDelete: boolean;
  prerequisiteIndex: number;
  ruleIndex: number;
  onDeleteCondition: () => void;
  onChangeFormField: (field: string, value: string | number | boolean) => void;
}

const ConditionForm = forwardRef(
  (
    {
      features,
      type,
      condition,
      isDisabledDelete,
      prerequisiteIndex,
      ruleIndex,
      onDeleteCondition,
      onChangeFormField
    }: Props,
    ref: Ref<HTMLDivElement>
  ) => {
    const { t } = useTranslation(['form', 'common', 'table']);

    const methods = useFormContext();
    const { control, watch } = methods;

    const featuresSelected = useMemo(
      () =>
        watch(`prerequisitesRules.${prerequisiteIndex}.rules`)?.map(
          (item: AnyObject) => item?.featureFlag
        ),
      [prerequisiteIndex]
    );

    const commonName = useMemo(
      () => `prerequisitesRules.${prerequisiteIndex}.rules.${ruleIndex}.`,
      [prerequisiteIndex, ruleIndex]
    );

    const featureFlag = watch(`${commonName}featureFlag`);

    const currentFeature = useMemo(
      () => features.find(item => item.id === featureFlag),
      [featureFlag, features]
    );

    const flagOptions = useMemo(
      () =>
        features
          .filter(f => !featuresSelected.includes(f.id))
          .map(item => ({
            label: item.name,
            value: item.id
          })),
      [features, featuresSelected]
    );

    const variationOptions = useMemo(
      () =>
        currentFeature?.variations?.map(item => ({
          label: item.name || item.value,
          value: item.id
        })),
      [currentFeature]
    );

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
          <div className="flex  w-full gap-x-4">
            <Form.Field
              control={control}
              name={`${commonName}featureFlag`}
              render={({ field }) => (
                <Form.Item className="flex flex-col flex-1 self-stretch py-0 min-w-[170px]">
                  <Form.Label required>{t('feature-flags.flag')}</Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        label={
                          flagOptions?.find(item =>
                            [field.value, condition.featureFlag].includes(
                              item.value
                            )
                          )?.label
                        }
                        placeholder={t('experiments.select-flag')}
                        className="w-full"
                        disabled={!flagOptions?.length}
                      />
                      <DropdownMenuContent align="start" {...field}>
                        {flagOptions?.map((item, index) => (
                          <DropdownMenuItem
                            key={index}
                            label={item.label}
                            value={item.value}
                            onSelectOption={value => {
                              field.onChange(value);
                              onChangeFormField('featureFlag', value);
                            }}
                          />
                        ))}
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <div className="flex-center size-fit min-w-fit px-3 py-3.5 mt-6 rounded bg-gray-100 typo-para-medium leading-5 text-gray-700 whitespace-nowrap">
              <Trans
                i18nKey={'form:feature-flags.receiving-state'}
                values={{
                  state: 'ON'
                }}
              />
            </div>
            <Form.Field
              control={control}
              name={`${commonName}variation`}
              render={({ field }) => (
                <Form.Item className="flex flex-col flex-1 self-stretch py-0 min-w-[170px]">
                  <Form.Label required>
                    {t('table:feature-flags.variation')}
                  </Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        label={
                          variationOptions?.find(item =>
                            [field.value, condition.variation].includes(
                              item.value
                            )
                          )?.label
                        }
                        placeholder={t('experiments.select-variation')}
                        className="w-full"
                        disabled={!variationOptions?.length}
                      />
                      <DropdownMenuContent align="start" {...field}>
                        {variationOptions?.map((item, index) => (
                          <DropdownMenuItem
                            key={index}
                            label={item.label}
                            value={item.value}
                            onSelectOption={value => {
                              field.onChange(value);
                              onChangeFormField('variation', value);
                            }}
                          />
                        ))}
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />

            <div className="flex items-center self-stretch order-5">
              <Button
                type="button"
                disabled={isDisabledDelete}
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
