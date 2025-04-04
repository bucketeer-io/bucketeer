import { useCallback } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { featureCreator } from '@api/features/feature-creator';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateFeatures } from '@queries/features';
import { useQueryTags } from '@queries/tags';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { AxiosError } from 'axios';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { cloneDeep } from 'lodash';
import { v4 as uuid } from 'uuid';
import { Feature, FeatureVariation, FeatureVariationType } from '@types';
import { onGenerateSlug } from 'utils/converts';
import { cn } from 'utils/style';
import {
  IconFlagJSON,
  IconFlagNumber,
  IconFlagString,
  IconFlagSwitch,
  IconInfo
} from '@icons';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import { CreatableSelect } from 'components/creatable-select';
import Divider from 'components/divider';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import TextArea from 'components/textarea';
import { Tooltip } from 'components/tooltip';
import { formSchema } from './formSchema';
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

export const flagTypeOptions = [
  {
    label: 'Boolean',
    value: 'BOOLEAN',
    icon: IconFlagSwitch
  },
  {
    label: 'String',
    value: 'STRING',
    icon: IconFlagString
  },
  {
    label: 'Number',
    value: 'NUMBER',
    icon: IconFlagNumber
  },
  {
    label: 'JSON',
    value: 'JSON',
    icon: IconFlagJSON
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

  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form']);
  const { notify } = useToast();

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
  const { watch } = form;

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
      const newVariations =
        value === 'BOOLEAN'
          ? cloneVariations
          : cloneVariations.map(item => ({
              ...item,
              value: ''
            }));
      form.setValue('variations', newVariations);
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
          message: 'Feature flag created successfully.'
        });
        invalidateFeatures(queryClient);
        onCompleted?.(resp.feature);
        onClose();
      }
    } catch (error) {
      const _error = error as AxiosError;
      const { status, message } = _error || {};
      notify({
        messageType: 'error',
        message:
          status === 409
            ? 'The same data already exists'
            : message || 'Something went wrong.'
      });
    }
  };
  return (
    <div className={cn('w-full p-5 pb-28', className)}>
      <p className="text-gray-700 typo-head-bold-small mb-2">
        {t('form:general-info')}
      </p>
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(onSubmit)}>
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
                  <Icon
                    icon={IconInfo}
                    size="xs"
                    color="gray-500"
                    className="absolute -right-6"
                  />
                </Form.Label>
                <Form.Control>
                  <Input
                    placeholder={`${t('form:feature-flags.placeholder-flag')}`}
                    {...field}
                  />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Form.Field
            control={form.control}
            name="description"
            render={({ field }) => (
              <Form.Item className="py-2.5">
                <Form.Label optional className="!mb-2">
                  {t('form:description')}
                </Form.Label>
                <Form.Control>
                  <TextArea
                    placeholder={t('form:placeholder-desc')}
                    rows={4}
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
                  <CreatableSelect
                    disabled={isLoadingTags}
                    loading={isLoadingTags}
                    placeholder={t(`form:placeholder-tags`)}
                    options={tags?.map(tag => ({
                      label: tag.name,
                      value: tag.id
                    }))}
                    onChange={value =>
                      field.onChange(value.map(tag => tag.label))
                    }
                  />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
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
                    trigger={
                      <div className="flex-center absolute top-0 -right-6">
                        <Icon icon={IconInfo} size={'sm'} color="gray-500" />
                      </div>
                    }
                    content={
                      <Trans
                        i18nKey={'table:feature-flags.variation-type'}
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
                    className="!z-[100] max-w-[400px]"
                  />
                </Form.Label>
                <Form.Control>
                  <DropdownMenu>
                    <DropdownMenuTrigger
                      placeholder={t(`form:feature-flags.flag-type`)}
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
                      variant="secondary"
                      className="w-full"
                    />
                    <DropdownMenuContent
                      className="w-[502px]"
                      align="start"
                      {...field}
                    >
                      {flagTypeOptions.map((item, index) => (
                        <DropdownMenuItem
                          {...field}
                          key={index}
                          icon={item.icon}
                          value={item.value}
                          label={item.label}
                          onSelectOption={value =>
                            handleOnChangeVariationType(
                              value as FeatureVariationType,
                              field.onChange
                            )
                          }
                        />
                      ))}
                    </DropdownMenuContent>
                  </DropdownMenu>
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Form.Field
            control={form.control}
            name="variations"
            render={({ field }) => (
              <Form.Item>
                <Form.Control>
                  <Variations
                    variationType={watch('variationType')}
                    variations={currentVariations}
                    onChangeVariations={field.onChange}
                  />
                </Form.Control>
              </Form.Item>
            )}
          />
          <div className="flex items-center w-full gap-x-4">
            <Form.Field
              control={form.control}
              name={`defaultOnVariation`}
              render={({ field }) => (
                <Form.Item className="py-2.5 flex-1">
                  <Form.Label className="!mb-2">
                    <Trans
                      i18nKey={'form:feature-flags.serve-targeting'}
                      values={{
                        state: 'ON'
                      }}
                    />
                  </Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        placeholder={t(`form:placeholder-tags`)}
                        trigger={
                          <div className="flex items-center gap-x-2">
                            <FlagVariationPolygon
                              index={currentVariations?.findIndex(
                                item => item.id === field.value
                              )}
                            />
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
                        variant="secondary"
                        className="w-full"
                      />
                      <DropdownMenuContent align="start" {...field}>
                        {currentVariations?.map((item, index) => (
                          <DropdownMenuItem
                            {...field}
                            key={index}
                            value={item.id}
                            label={`Variation ${index + 1}`}
                            onSelectOption={() => {
                              field.onChange(item.id);
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
            <Form.Field
              control={form.control}
              name={`defaultOffVariation`}
              render={({ field }) => (
                <Form.Item className="py-2.5 flex-1">
                  <Form.Label className="!mb-2">
                    <Trans
                      i18nKey={'form:feature-flags.serve-targeting'}
                      values={{
                        state: 'OFF'
                      }}
                    />
                  </Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        placeholder={t(`form:placeholder-tags`)}
                        trigger={
                          <div className="flex items-center gap-x-2">
                            <FlagVariationPolygon
                              index={currentVariations?.findIndex(
                                item => item.id === field.value
                              )}
                            />
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
                        variant="secondary"
                        className="w-full"
                      />
                      <DropdownMenuContent align="start" {...field}>
                        {currentVariations?.map((item, index) => (
                          <DropdownMenuItem
                            {...field}
                            key={index}
                            value={item.id}
                            label={`Variation ${index + 1}`}
                            onSelectOption={() => {
                              field.onChange(item.id);
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
          </div>
          <div className="absolute left-0 bottom-0 bg-gray-50 w-full rounded-b-lg z-[999]">
            <ButtonBar
              primaryButton={
                <Button variant="secondary" onClick={onClose}>
                  {t(`cancel`)}
                </Button>
              }
              secondaryButton={
                <Button
                  type="submit"
                  disabled={!form.formState.isDirty}
                  loading={form.formState.isSubmitting}
                >
                  {t(`create-flag`)}
                </Button>
              }
            />
          </div>
        </Form>
      </FormProvider>
    </div>
  );
};

export default CreateFlagForm;
