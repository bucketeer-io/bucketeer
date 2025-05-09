import { useCallback, useMemo } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { featureUpdater } from '@api/features';
import { yupResolver } from '@hookform/resolvers/yup';
import { useQueryAccounts } from '@queries/accounts';
import { invalidateFeature } from '@queries/feature-details';
import { useQueryTags } from '@queries/tags';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast, useToggleOpen } from 'hooks';
import { useTranslation } from 'i18n';
import { uniqBy } from 'lodash';
import { Feature } from '@types';
import Button from 'components/button';
import { CreatableSelect } from 'components/creatable-select';
import Form from 'components/form';
import Input from 'components/input';
import TextArea from 'components/textarea';
import Card from 'elements/card';
import DropdownMenuWithSearch from 'elements/dropdown-with-search';
import { generalInfoFormSchema, GeneralInfoFormType } from './form-schema';
import SaveWithCommentModal from './modals/save-with-comment';

const GeneralInfoForm = ({ feature }: { feature: Feature }) => {
  const { t } = useTranslation(['form', 'common', 'table', 'message']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { notify, errorNotify } = useToast();
  const queryClient = useQueryClient();

  const [isOpenSaveModal, onOpenSaveModal, onCloseSaveModal] =
    useToggleOpen(false);

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

  const tagOptions = uniqBy(
    tags.map(item => ({
      label: item.name,
      value: item.name
    })),
    'label'
  );
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
        tags: {
          values: tags
        },
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
        onCloseSaveModal();
      }
    } catch (error) {
      errorNotify(error);
    }
  }, [currentEnvironment]);

  return (
    <FormProvider {...form}>
      <Form onSubmit={form.handleSubmit(onSubmit)}>
        <Card>
          <p className="typo-head-bold-small text-gray-800">
            {t('general-info')}
          </p>
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
                <Form.Label required>{t('common:tags')}</Form.Label>
                <Form.Control>
                  <CreatableSelect
                    disabled={isLoadingTags || !tagOptions.length}
                    value={tagOptions.filter(tag =>
                      field.value?.includes(tag.value)
                    )}
                    loading={isLoadingTags}
                    placeholder={t(
                      !tagOptions.length && !isLoadingTags
                        ? `form:no-tags-found`
                        : `form:placeholder-tags`
                    )}
                    options={tagOptions}
                    onChange={value =>
                      field.onChange(value.map(tag => tag.value))
                    }
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
