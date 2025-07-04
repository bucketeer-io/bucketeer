import { useCallback, useMemo } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { apiKeyCreator, APIKeyResponse, apiKeyUpdater } from '@api/api-key';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateAPIKeys } from '@queries/api-keys';
import { useQueryClient } from '@tanstack/react-query';
import { useAuthAccess } from 'auth';
import { useToast } from 'hooks';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { APIKey, APIKeyRole, Environment } from '@types';
import { checkEnvironmentEmptyId, onFormatEnvironments } from 'utils/function';
import { cn } from 'utils/style';
import { apiKeyOptions } from 'pages/api-keys/constants';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
  DropdownOption
} from 'components/dropdown';
import Form from 'components/form';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import { RadioGroup, RadioGroupItem } from 'components/radio';
import TextArea from 'components/textarea';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import DropdownList from 'elements/dropdown-list';
import FormLoading from 'elements/form-loading';

interface APIKeyCreateUpdateModalProps {
  isOpen: boolean;
  isLoadingApiKey: boolean;
  apiKeyEnvironmentId: string;
  environments: Environment[];
  isLoadingEnvs: boolean;
  apiKey?: APIKey;
  onClose: () => void;
}

export interface APIKeyCreateUpdateForm {
  name: string;
  environmentId: string;
  description?: string;
  role: APIKeyRole;
}

export const formSchema = ({ requiredMessage }: FormSchemaProps) =>
  yup.object().shape({
    name: yup.string().required(requiredMessage),
    environmentId: yup.string().required(requiredMessage),
    description: yup.string(),
    role: yup.mixed<APIKeyRole>().required(requiredMessage)
  });

const APIKeyCreateUpdateModal = ({
  isOpen,
  isLoadingApiKey,
  isLoadingEnvs,
  apiKeyEnvironmentId,
  apiKey,
  environments,
  onClose
}: APIKeyCreateUpdateModalProps) => {
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form', 'message']);
  const { notify } = useToast();

  const { envEditable, isOrganizationAdmin } = useAuthAccess();

  const { emptyEnvironmentId, formattedEnvironments } =
    onFormatEnvironments(environments);

  const isEditApiKey = useMemo(
    () => !!apiKey || !!isLoadingApiKey || !!apiKeyEnvironmentId,
    [apiKey, isLoadingApiKey, apiKeyEnvironmentId]
  );

  const disabled = useMemo(
    () => !envEditable || !isOrganizationAdmin,
    [envEditable, isOrganizationAdmin]
  );

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    values: {
      name: apiKey?.name || '',
      environmentId: apiKey ? apiKeyEnvironmentId || emptyEnvironmentId : '',
      description: apiKey?.description,
      role: apiKey?.role || 'SDK_CLIENT'
    }
  });

  const {
    watch,
    formState: { isValid, isSubmitting, isDirty }
  } = form;

  const environmentIdWatch = watch('environmentId');

  const currentEnv = useMemo(
    () => formattedEnvironments.find(item => item.id === environmentIdWatch),
    [formattedEnvironments, environmentIdWatch]
  );

  const environmentOptions = useMemo(
    () =>
      formattedEnvironments.map(item => ({
        label: item.name,
        value: item.id
      })),
    [formattedEnvironments]
  );

  const onSubmit: SubmitHandler<APIKeyCreateUpdateForm> = useCallback(
    async values => {
      let resp: APIKeyResponse | null = null;
      const { environmentId, name, description, role } = values;
      const envId = checkEnvironmentEmptyId(environmentId);
      if (isEditApiKey) {
        resp = await apiKeyUpdater({
          id: apiKey!.id,
          environmentId: envId,
          description,
          name
        });
      } else {
        resp = await apiKeyCreator({
          environmentId: envId,
          name,
          role,
          description
        });
      }
      if (resp) {
        notify({
          message: t('message:collection-action-success', {
            collection: t('source-type.api-key'),
            action: t(apiKey ? 'updated' : 'created')
          })
        });
        invalidateAPIKeys(queryClient);
        onClose();
      }
    },
    [apiKey, isEditApiKey]
  );

  return (
    <SlideModal
      title={t(isEditApiKey ? 'update-api-key' : 'new-api-key')}
      isOpen={isOpen}
      onClose={onClose}
    >
      {isLoadingApiKey ? (
        <FormLoading />
      ) : (
        <div className="w-full p-5 pb-28">
          <p className="text-gray-800 typo-head-bold-small">
            {t('form:general-info')}
          </p>
          <FormProvider {...form}>
            <Form onSubmit={form.handleSubmit(onSubmit)}>
              <Form.Field
                control={form.control}
                name="name"
                render={({ field }) => (
                  <Form.Item>
                    <Form.Label required>{t('name')}</Form.Label>
                    <Form.Control>
                      <Input
                        placeholder={`${t('form:placeholder-name')}`}
                        disabled={disabled}
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
                  <Form.Item>
                    <Form.Label optional>{t('form:description')}</Form.Label>
                    <Form.Control>
                      <TextArea
                        placeholder={t('form:placeholder-desc')}
                        rows={4}
                        disabled={disabled}
                        {...field}
                      />
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />

              <p className="text-gray-800 typo-head-bold-small">
                {t('environment')}
              </p>
              <Form.Field
                control={form.control}
                name={`environmentId`}
                render={({ field }) => (
                  <Form.Item className="py-2">
                    <Form.Label required>{t('environment')}</Form.Label>
                    <Form.Control>
                      <DropdownMenu>
                        <DropdownMenuTrigger
                          placeholder={t(`form:select-environment`)}
                          label={currentEnv?.name}
                          disabled={!!apiKey || disabled}
                          variant="secondary"
                          loading={isLoadingEnvs}
                          className="w-full"
                        />
                        <DropdownMenuContent
                          className={cn('w-[502px]', {
                            'hidden-scroll': environmentOptions?.length > 15
                          })}
                          align="start"
                          {...field}
                        >
                          <DropdownList
                            options={environmentOptions as DropdownOption[]}
                            onSelectOption={field.onChange}
                          />
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />

              <div className="flex items-center gap-2 mt-4">
                <p className="text-gray-800 typo-head-bold-small">
                  {t('key-role')}
                </p>
              </div>
              <Form.Field
                control={form.control}
                name="role"
                render={({ field }) => (
                  <Form.Item>
                    <Form.Control>
                      <RadioGroup
                        defaultValue={field.value}
                        disabled={!!apiKey || disabled}
                        onValueChange={field.onChange}
                      >
                        {apiKeyOptions.map(
                          ({ id, label, description, value }) => (
                            <div
                              key={id}
                              className="flex items-center last:border-b-0 border-b py-4 gap-x-5"
                            >
                              <label
                                htmlFor={id}
                                className={cn('flex-1', {
                                  'opacity-50': !!apiKey || disabled
                                })}
                              >
                                <p className="typo-para-medium text-gray-700">
                                  {label}
                                </p>
                                <p className="typo-para-small text-gray-600">
                                  {description}
                                </p>
                              </label>
                              <RadioGroupItem
                                disabled={!!apiKey || disabled}
                                value={value}
                                id={id}
                              />
                            </div>
                          )
                        )}
                      </RadioGroup>
                    </Form.Control>
                  </Form.Item>
                )}
              />
              <div className="absolute left-0 bottom-0 bg-gray-50 w-full rounded-b-lg">
                <ButtonBar
                  primaryButton={
                    <Button variant="secondary" onClick={onClose}>
                      {t(`cancel`)}
                    </Button>
                  }
                  secondaryButton={
                    <DisabledButtonTooltip
                      type={!isOrganizationAdmin ? 'admin' : 'editor'}
                      hidden={envEditable && isOrganizationAdmin}
                      trigger={
                        <Button
                          type="submit"
                          disabled={!isValid || !isDirty || disabled}
                          loading={isSubmitting}
                        >
                          {t(`submit`)}
                        </Button>
                      }
                    />
                  }
                />
              </div>
            </Form>
          </FormProvider>
        </div>
      )}
    </SlideModal>
  );
};

export default APIKeyCreateUpdateModal;
