import { useCallback, useEffect, useMemo, useState } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { featureUpdater } from '@api/features';
import { yupResolver } from '@hookform/resolvers/yup';
import { useQueryAccounts } from '@queries/accounts';
import { invalidateFeature } from '@queries/feature-details';
import { invalidateFeatures } from '@queries/features';
import { invalidateTags, useQueryTags } from '@queries/tags';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import { Feature, TagChange } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { cn } from 'utils/style';
import { IconInfo, IconWatch } from '@icons';
import Button from 'components/button';
import { DropdownOption } from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import TextArea from 'components/textarea';
import { Tooltip } from 'components/tooltip';
import Card from 'elements/card';
import DateTooltip from 'elements/date-tooltip';
import DropdownMenuWithSearch from 'elements/dropdown-with-search';
import { generalInfoFormSchema, GeneralInfoFormType } from './form-schema';
import SaveWithCommentModal from './modals/save-with-comment';

const GeneralInfoForm = ({ feature }: { feature: Feature }) => {
  const { t } = useTranslation(['form', 'common', 'table', 'message']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { notify, errorNotify } = useToast();
  const queryClient = useQueryClient();
  const formatDateTime = useFormatDateTime();

  const [isOpenSaveModal, onOpenSaveModal, onCloseSaveModal] =
    useToggleOpen(false);
  const [tagOptions, setTagOptions] = useState<DropdownOption[]>([]);

  const { data: tagCollection, isLoading: isLoadingTags } = useQueryTags({
    params: {
      cursor: String(0),
      entityType: 'FEATURE_FLAG',
      environmentId: currentEnvironment.id
    }
  });

  const { data: accountCollection, isLoading: isLoadingAccounts } =
    useQueryAccounts({
      params: {
        cursor: String(0),
        environmentId: currentEnvironment.id,
        organizationId: currentEnvironment.organizationId,
        environmentRole: 2,
        organizationRole: 2
      }
    });

  const form = useForm<GeneralInfoFormType>({
    resolver: yupResolver(generalInfoFormSchema),
    defaultValues: {
      maintainer: feature.maintainer,
      name: feature.name,
      flagId: feature.id,
      description: feature.description,
      tags: feature.tags,
      comment: ''
    },
    mode: 'onChange'
  });
  const {
    formState: { isValid, isDirty },
    watch,
    getValues,
    setError,
    resetField
  } = form;
  const maintainer = watch('maintainer');
  const tags = tagCollection?.tags || [];
  const accounts = accountCollection?.accounts || [];

  const accountOptions = accounts.map(item => ({
    label: item.email,
    value: item.email
  }));

  const maintainerLabel = useMemo(() => {
    const currentAccount = accounts.find(item => item.email === maintainer);
    if (currentAccount?.firstName && currentAccount?.lastName)
      return `${currentAccount.firstName} ${currentAccount.lastName}`;
    return currentAccount?.email || maintainer;
  }, [accounts, maintainer]);

  const handleCheckTags = useCallback(
    (tagValues: string[]) => {
      const tagChanges: TagChange[] = [];
      const { tags } = feature;
      tags?.forEach(item => {
        if (!tagValues.find(tag => tag === item)) {
          tagChanges.push({
            changeType: 'DELETE',
            tag: item
          });
        }
      });
      tagValues.forEach(item => {
        const currentTag = tags.find(tag => tag === item);
        if (!currentTag) {
          tagChanges.push({
            changeType: 'CREATE',
            tag: item
          });
        }
      });

      return {
        tagChanges
      };
    },
    [feature]
  );

  const onSubmit = useCallback(async () => {
    try {
      const values = getValues();
      const { flagId, comment, tags, ...rest } = values;
      if (currentEnvironment.requireComment && !comment)
        return setError('comment', {
          message: t('message:required-field')
        });

      const resp = await featureUpdater({
        id: flagId,
        environmentId: currentEnvironment.id,
        comment,
        ...handleCheckTags(tags),
        ...rest
      });

      if (resp) {
        notify({
          message: t('message:flag-updated')
        });
        form.reset({
          ...values,
          comment: ''
        });
        invalidateFeature(queryClient);
        invalidateFeatures(queryClient);
        invalidateTags(queryClient);
        onCloseSaveModal();
      }
    } catch (error) {
      errorNotify(error);
    }
  }, [currentEnvironment, feature]);

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
    <FormProvider {...form}>
      <Form onSubmit={form.handleSubmit(onSubmit)}>
        <Card>
          <div className="flex items-center gap-x-2">
            <p className="typo-head-bold-small text-gray-800">
              {t('general-info')}
            </p>
            <DateTooltip
              trigger={
                <div className="flex items-center gap-x-2 text-gray-700 typo-para-small whitespace-nowrap -mb-1">
                  <Icon icon={IconWatch} size={'xxs'} />
                  {Number(feature.createdAt) === 0 ? (
                    t('never')
                  ) : (
                    <Trans
                      i18nKey={'common:time-created'}
                      values={{
                        time: formatDateTime(feature.createdAt)
                      }}
                    />
                  )}
                </div>
              }
              date={Number(feature.createdAt) === 0 ? null : feature.createdAt}
            />
          </div>

          <Form.Field
            name="maintainer"
            render={({ field }) => (
              <Form.Item className="w-full py-0">
                <Form.Label required>{t('common:maintainer')}</Form.Label>
                <Form.Control>
                  <DropdownMenuWithSearch
                    isLoading={isLoadingAccounts}
                    placeholder={t('placeholder-maintainer')}
                    label={maintainerLabel}
                    options={accountOptions}
                    selectedOptions={[field.value]}
                    onSelectOption={field.onChange}
                  />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Form.Field
            name="name"
            render={({ field }) => (
              <Form.Item className="w-full py-0">
                <Form.Label required>{t('common:name')}</Form.Label>
                <Form.Control>
                  <Input {...field} placeholder={t('placeholder-name')} />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Form.Field
            name="flagId"
            render={({ field }) => (
              <Form.Item className="w-full py-0">
                <Form.Label required>{t('feature-flags.flag-id')}</Form.Label>
                <Form.Control>
                  <Input
                    {...field}
                    disabled
                    placeholder={t('feature-flags.placeholder-flag')}
                  />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Form.Field
            name="description"
            render={({ field }) => (
              <Form.Item className="w-full py-0">
                <Form.Label>{t('description')}</Form.Label>
                <Form.Control>
                  <TextArea
                    {...field}
                    placeholder={t('placeholder-desc')}
                    rows={2}
                    style={{
                      resize: 'vertical',
                      maxHeight: 98
                    }}
                  />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Form.Field
            name="tags"
            render={({ field }) => (
              <Form.Item className="w-full py-0">
                <Form.Label required className="relative w-fit">
                  {t('common:tags')}
                  <Tooltip
                    align="start"
                    alignOffset={-46}
                    content={t('tags-tooltip')}
                    trigger={
                      <div className="flex-center size-fit absolute top-0 -right-6">
                        <Icon icon={IconInfo} size="xs" color="gray-500" />
                      </div>
                    }
                    className="max-w-[400px]"
                  />
                </Form.Label>
                <Form.Control>
                  <DropdownMenuWithSearch
                    label={field.value?.join(', ') || ''}
                    isExpand
                    isMultiselect
                    disabled={isLoadingTags}
                    placeholder={t('select-or-create-tags')}
                    options={tagOptions}
                    selectedOptions={field.value}
                    onKeyDown={({
                      event,
                      searchValue,
                      matchOptions,
                      onClearSearchValue
                    }) => {
                      const value: string = matchOptions?.length
                        ? (matchOptions[0].value as string)
                        : searchValue;
                      if (
                        event.key === 'Enter' &&
                        !field.value?.includes(value)
                      ) {
                        if (!matchOptions?.length)
                          setTagOptions([
                            ...tagOptions,
                            {
                              label: value,
                              value
                            }
                          ]);
                        field.onChange([...field.value, value]);
                        onClearSearchValue();
                      }
                    }}
                    onSelectOption={value => {
                      const isExisted = field.value?.find(
                        (item: string) => item === value
                      );
                      field.onChange(
                        isExisted
                          ? field.value?.filter(
                              (item: string) => item !== value
                            )
                          : [...field.value, value]
                      );
                    }}
                    notFoundOption={(searchValue, onChangeValue) => {
                      const isExisted = field.value?.find(
                        (item: string) => item === searchValue
                      );
                      return (
                        searchValue && (
                          <div
                            className={cn(
                              'flex items-center py-2 px-4 my-1 rounded pointer-events-none',
                              {
                                'hover:bg-gray-100 cursor-pointer pointer-events-auto':
                                  !isExisted
                              }
                            )}
                            onClick={() => {
                              field.onChange([...field.value, searchValue]);
                              tagOptions.push({
                                label: searchValue,
                                value: searchValue
                              });
                              onChangeValue('');
                            }}
                          >
                            <p className="text-gray-700">
                              {t('create-tag-name', {
                                name: searchValue
                              })}
                            </p>
                          </div>
                        )
                      );
                    }}
                  />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Button
            type="button"
            variant={'secondary'}
            disabled={!isValid || !isDirty}
            className="w-fit"
            onClick={onOpenSaveModal}
          >
            {t('common:save-with-comment')}
          </Button>
        </Card>
        {isOpenSaveModal && (
          <SaveWithCommentModal
            isOpen={isOpenSaveModal}
            isRequired={currentEnvironment.requireComment}
            onClose={() => {
              onCloseSaveModal();
              resetField('comment');
            }}
            onSubmit={onSubmit}
          />
        )}
      </Form>
    </FormProvider>
  );
};

export default GeneralInfoForm;
