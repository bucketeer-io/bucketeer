import { useCallback, useEffect, useRef, useState } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { featureCreator } from '@api/features/feature-creator';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateFeatures } from '@queries/features';
import { useQueryTags } from '@queries/tags';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, hasEditable, useAuth } from 'auth';
import { getDefaultYamlValue } from 'constants/feature-flag';
import { useToast } from 'hooks';
import useFormSchema from 'hooks/use-form-schema';
import useOptions from 'hooks/use-options';
import { getLanguage, Language, useTranslation } from 'i18n';
import cloneDeep from 'lodash/cloneDeep';
import { v4 as uuid } from 'uuid';
import { Feature, FeatureVariation, FeatureVariationType } from '@types';
import { onGenerateSlug } from 'utils/converts';
import { cn } from 'utils/style';
import { IconInfo } from '@icons';
import { createFlagFormSchema } from 'pages/create-flag/form-schema';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Dropdown, { DropdownOption } from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import { Tooltip } from 'components/tooltip';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import SelectMenu from 'elements/select-menu';
import Variations from './variations';

export interface AddFlagForm {
  name: string;
  flagId: string;
  tags: string[];
  variationType: FeatureVariationType;
  variations: FeatureVariation[];
  defaultOnVariation: string;
  defaultOffVariation: string;
  description?: string;
}

const defaultVariations: FeatureVariation[] = [
  {
    id: uuid(),
    value: 'true',
    name: '',
    description: ''
  },
  {
    id: uuid(),
    value: 'false',
    name: '',
    description: ''
  }
];

const CreateFlagForm = ({
  className,
  onClose,
  onCompleted
}: {
  className?: string;
  onClose: () => void;
  onCompleted?: (flag: Feature) => void;
}) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const editable = hasEditable(consoleAccount!);
  const formSchema = useFormSchema(createFlagFormSchema);
  const { flagTypeOptions } = useOptions();
  const refFormModel = useRef<HTMLDivElement>(null);
  const queryClient = useQueryClient();
  const isJapaneseLanguage = getLanguage() === Language.JAPANESE;

  const { t } = useTranslation(['common', 'form', 'message']);
  const { notify, errorNotify } = useToast();

  const [tagOptions, setTagOptions] = useState<DropdownOption[]>([]);

  const { data: collection, isLoading: isLoadingTags } = useQueryTags({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment?.id,
      entityType: 'FEATURE_FLAG'
    }
  });

  const tags = collection?.tags || [];

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      name: '',
      flagId: '',
      description: '',
      variationType: 'BOOLEAN',
      tags: [],
      variations: defaultVariations,
      defaultOnVariation: defaultVariations[0].id,
      defaultOffVariation: defaultVariations[1].id
    },
    mode: 'onChange'
  });

  const {
    formState: { isDirty, isValid },
    watch
  } = form;

  const variationType = watch('variationType');

  const currentFlagOption = flagTypeOptions.find(
    item => item.value === variationType
  );

  const currentVariations = watch('variations') as FeatureVariation[];

  const handleOnChangeVariationType = useCallback(
    (
      value: FeatureVariationType,
      onChange: (value: FeatureVariationType) => void
    ) => {
      const cloneVariations = cloneDeep(defaultVariations);
      const isBoolean = value === 'BOOLEAN';
      const isJSON = value === 'JSON';
      const isYAML = value === 'YAML';

      let newVariations;

      if (isBoolean) {
        newVariations = cloneVariations;
      } else {
        newVariations = cloneVariations.map((item, index) => {
          let defaultValue = '';

          if (isJSON) {
            defaultValue = '{}';
          } else if (isYAML) {
            defaultValue = getDefaultYamlValue(index);
          }

          return {
            ...item,
            value: defaultValue
          };
        });
      }

      form.resetField('variations', { defaultValue: newVariations });

      onChange(value);
      let timerId: NodeJS.Timeout | null = null;
      if (timerId) clearTimeout(timerId);
      timerId = setTimeout(() => form.setFocus('variations.0.value'), 100);
    },
    [form, defaultVariations]
  );

  const onSubmit: SubmitHandler<AddFlagForm> = async values => {
    try {
      const {
        flagId,
        name,
        tags,
        variationType,
        defaultOffVariation,
        defaultOnVariation,
        description,
        variations
      } = values;
      const resp = await featureCreator({
        environmentId: currentEnvironment.id,
        id: flagId,
        name,
        tags,
        defaultOnVariationIndex: variations.findIndex(
          item => item.id === defaultOnVariation
        ),
        defaultOffVariationIndex: variations.findIndex(
          item => item.id === defaultOffVariation
        ),
        variations,
        variationType,
        description
      });
      if (resp) {
        notify({
          message: t('message:collection-action-success', {
            collection: t('source-type.feature-flag'),
            action: t('created')
          })
        });
        invalidateFeatures(queryClient);
        onCompleted?.(resp.feature);
        onClose();
      }
    } catch (error) {
      errorNotify(error);
    }
  };

  useEffect(() => {
    if (tags.length) {
      setTagOptions(
        tags.map(item => ({
          label: item.name,
          value: item.name
        }))
      );
    }
  }, [tags]);

  return (
    <div ref={refFormModel} className={cn('w-full p-5 pb-28', className)}>
      <p className="text-gray-700 typo-head-bold-small mb-2">
        {t('form:general-info')}
      </p>
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(onSubmit)}>
          <div className="w-full flex gap-x-4 [&>div]:flex-1">
            <Form.Field
              control={form.control}
              name="name"
              render={({ field }) => (
                <Form.Item className="py-2.5">
                  <Form.Label required className="!mb-2">
                    {t('name')}
                  </Form.Label>
                  <Form.Control>
                    <Input
                      placeholder={`${t('form:placeholder-name')}`}
                      {...field}
                      onChange={value => {
                        const isFlagIdDirty =
                          form.getFieldState('flagId').isDirty;
                        const flagId = form.getValues('flagId');
                        field.onChange(value);
                        form.setValue(
                          'flagId',
                          isFlagIdDirty ? flagId : onGenerateSlug(value)
                        );
                      }}
                      name="flag-name"
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name="flagId"
              render={({ field }) => (
                <Form.Item className="py-2.5">
                  <Form.Label required className="relative w-fit !mb-2">
                    {t('form:feature-flags.flag-id')}
                    <Tooltip
                      align="start"
                      alignOffset={-56}
                      trigger={
                        <div className="flex-center absolute top-0 -right-6">
                          <Icon icon={IconInfo} size={'sm'} color="gray-500" />
                        </div>
                      }
                      content={t('form:flag-id-tooltip')}
                      className="!z-[100] max-w-[400px]"
                    />
                  </Form.Label>
                  <Form.Control>
                    <Input
                      placeholder={`${t('form:feature-flags.placeholder-flag')}`}
                      {...field}
                      name="flag-id"
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
          </div>
          <div className="w-full flex gap-x-4 [&>div]:flex-1">
            <Form.Field
              control={form.control}
              name="description"
              render={({ field }) => (
                <Form.Item className="py-2.5">
                  <Form.Label optional className="!mb-2">
                    {t('form:description')}
                  </Form.Label>
                  <Form.Control>
                    <Input
                      placeholder={t('form:placeholder-desc')}
                      {...field}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name={`tags`}
              render={({ field }) => (
                <Form.Item className="py-2.5">
                  <Form.Label required className="!mb-2">
                    {t('tags')}
                  </Form.Label>
                  <Form.Control>
                    <SelectMenu
                      options={tagOptions}
                      fieldValues={field.value}
                      onChange={field.onChange}
                      disabled={isLoadingTags}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
          </div>
          <Divider className="mt-2.5 mb-4" />
          <p className="text-gray-700 typo-head-bold-small mb-2">
            {t('form:feature-flags.flag-variations')}
          </p>
          <Form.Field
            control={form.control}
            name={`variationType`}
            render={({ field }) => (
              <Form.Item className="py-2">
                <Form.Label required className="relative w-fit !mb-2">
                  {t('form:feature-flags.flag-type')}
                  <Tooltip
                    align="start"
                    alignOffset={-30}
                    trigger={
                      <div className="flex-center absolute top-0 -right-6">
                        <Icon icon={IconInfo} size={'sm'} color="gray-500" />
                      </div>
                    }
                    content={
                      <Trans
                        i18nKey={'form:flag-type-tooltip'}
                        values={{
                          type:
                            variationType === 'JSON'
                              ? variationType
                              : variationType?.toLowerCase()
                        }}
                        components={{
                          text: <span className="capitalize" />
                        }}
                      />
                    }
                    className="!z-[100] max-w-[300px]"
                  />
                </Form.Label>
                <Form.Control>
                  <Dropdown
                    placeholder={t(`form:feature-flags.flag-type`)}
                    options={flagTypeOptions}
                    trigger={
                      <div className="flex items-center gap-x-2">
                        {currentFlagOption?.icon && (
                          <Icon
                            icon={currentFlagOption?.icon}
                            size={'md'}
                            className="flex-center"
                          />
                        )}
                        <p>{currentFlagOption?.label}</p>
                      </div>
                    }
                    disabled={isLoadingTags}
                    value={field.value}
                    onChange={value =>
                      handleOnChangeVariationType(
                        value as FeatureVariationType,
                        field.onChange
                      )
                    }
                    className="w-full"
                    contentClassName="w-[502px]"
                  />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Form.Field
            control={form.control}
            name="variations"
            render={() => (
              <Form.Item>
                <Form.Control>
                  <Variations
                    refModel={refFormModel}
                    variationType={watch('variationType')}
                  />
                </Form.Control>
              </Form.Item>
            )}
          />
          <div className="flex items-center w-full gap-x-4">
            <Form.Field
              control={form.control}
              name={`defaultOnVariation`}
              render={({ field }) => {
                const variationIndex = currentVariations?.findIndex(
                  item => item.id === field.value
                );
                return (
                  <Form.Item className="py-2.5 flex-1">
                    <Form.Label className="!mb-2">
                      <Trans
                        i18nKey={'form:feature-flags.serve-targeting'}
                        values={{
                          state: isJapaneseLanguage
                            ? t('form:experiments.on')
                            : t('form:experiments.on').toUpperCase()
                        }}
                      />
                    </Form.Label>
                    <Form.Control>
                      <Dropdown
                        placeholder={t(`form:placeholder-tags`)}
                        trigger={
                          <div className="flex items-center gap-x-2">
                            <FlagVariationPolygon index={variationIndex} />
                            <Trans
                              i18nKey={'form:feature-flags.variation'}
                              values={{
                                index:
                                  currentVariations?.findIndex(
                                    item => item.id === field.value
                                  ) + 1
                              }}
                            />
                          </div>
                        }
                        options={currentVariations?.map((item, index) => ({
                          label: `Variation ${index + 1}`,
                          value: item.id
                        }))}
                        value={field.value}
                        onChange={field.onChange}
                        className="w-full"
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                );
              }}
            />
            <Form.Field
              control={form.control}
              name={`defaultOffVariation`}
              render={({ field }) => {
                const variationIndex = currentVariations?.findIndex(
                  item => item.id === field.value
                );
                return (
                  <Form.Item className="py-2.5 flex-1">
                    <Form.Label className="!mb-2">
                      <Trans
                        i18nKey={'form:feature-flags.serve-targeting'}
                        values={{
                          state: isJapaneseLanguage
                            ? t('form:experiments.off')
                            : t('form:experiments.off').toUpperCase()
                        }}
                      />
                    </Form.Label>
                    <Form.Control>
                      <Dropdown
                        placeholder={t(`form:placeholder-tags`)}
                        trigger={
                          <div className="flex items-center gap-x-2">
                            <FlagVariationPolygon index={variationIndex} />
                            <Trans
                              i18nKey={'form:feature-flags.variation'}
                              values={{
                                index:
                                  currentVariations?.findIndex(
                                    item => item.id === field.value
                                  ) + 1
                              }}
                            />
                          </div>
                        }
                        options={currentVariations?.map((item, index) => ({
                          label: `Variation ${index + 1}`,
                          value: item.id
                        }))}
                        value={field.value}
                        onChange={field.onChange}
                        className="w-full"
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                );
              }}
            />
          </div>
          <div className="absolute left-0 bottom-0 bg-gray-50 w-full rounded-b-lg z-[999]">
            <ButtonBar
              primaryButton={
                <Button variant="secondary" onClick={onClose}>
                  {t(`cancel`)}
                </Button>
              }
              secondaryButton={
                <DisabledButtonTooltip
                  hidden={editable}
                  trigger={
                    <Button
                      type="submit"
                      disabled={!isDirty || !isValid || !editable}
                      loading={form.formState.isSubmitting}
                    >
                      {t(`create-flag`)}
                    </Button>
                  }
                />
              }
            />
          </div>
        </Form>
      </FormProvider>
    </div>
  );
};

export default CreateFlagForm;
