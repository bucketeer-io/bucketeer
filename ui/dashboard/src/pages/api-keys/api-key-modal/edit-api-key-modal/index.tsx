import { useMemo } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { apiKeyUpdater } from '@api/api-key';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateAPIKeys } from '@queries/api-keys';
import { useQueryClient } from '@tanstack/react-query';
import { useAuthAccess } from 'auth';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { APIKey, Environment } from '@types';
import { checkEnvironmentEmptyId, onFormatEnvironments } from 'utils/function';
import { IconInfo } from '@icons';
import { apiKeyOptions } from 'pages/api-keys/constants';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import { RadioGroup, RadioGroupItem } from 'components/radio';
import TextArea from 'components/textarea';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import FormLoading from 'elements/form-loading';

interface EditAPIKeyModalProps {
  isOpen: boolean;
  isLoadingApiKey: boolean;
  environments: Environment[];
  apiKey?: APIKey;
  onClose: () => void;
}

export interface EditAPIKeyForm {
  name: string;
  environmentId: string;
  description?: string;
}

export const formSchema = yup.object().shape({
  name: yup.string().required(),
  environmentId: yup.string().required(),
  description: yup.string()
});

const EditAPIKeyModal = ({
  isOpen,
  isLoadingApiKey,
  apiKey,
  environments,
  onClose
}: EditAPIKeyModalProps) => {
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form']);
  const { notify } = useToast();

  const { envEditable, isOrganizationAdmin } = useAuthAccess();

  const { emptyEnvironmentId, formattedEnvironments } =
    onFormatEnvironments(environments);

  const apiEnvironmentId = useMemo(
    () =>
      formattedEnvironments.find(item => item.name === apiKey?.environmentName)
        ?.id,
    [formattedEnvironments, apiKey]
  );

  const disabled = useMemo(
    () => !envEditable || !isOrganizationAdmin,
    [envEditable, isOrganizationAdmin]
  );

  const form = useForm({
    resolver: yupResolver(formSchema),
    values: {
      name: apiKey?.name || '',
      environmentId: apiEnvironmentId || emptyEnvironmentId,
      description: apiKey?.description
    }
  });

  const {
    watch,
    formState: { isValid, isSubmitting, isDirty }
  } = form;

  const environmentId = watch('environmentId');

  const currentEnv = useMemo(
    () => formattedEnvironments.find(item => item.id === environmentId),
    [formattedEnvironments, environmentId]
  );

  const onSubmit: SubmitHandler<EditAPIKeyForm> = values => {
    return apiKeyUpdater({
      id: apiKey!.id,
      environmentId: checkEnvironmentEmptyId(values.environmentId),
      description: values.description,
      name: values.name
    }).then(() => {
      notify({
        toastType: 'toast',
        messageType: 'success',
        message: (
          <span>
            <b>{values.name}</b> {` has been successfully updated!`}
          </span>
        )
      });
      invalidateAPIKeys(queryClient);
      onClose();
    });
  };

  return (
    <SlideModal title={t('update-api-key')} isOpen={isOpen} onClose={onClose}>
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
                          disabled
                          variant="secondary"
                          className="w-full"
                        />
                        <DropdownMenuContent
                          className="w-[502px]"
                          align="start"
                          {...field}
                        >
                          {formattedEnvironments.map((item, index) => (
                            <DropdownMenuItem
                              {...field}
                              key={index}
                              value={item.id}
                              label={item.name}
                              onSelectOption={value => {
                                field.onChange(value);
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

              <div className="flex items-center gap-2 mt-4">
                <p className="text-gray-800 typo-head-bold-small">
                  {t('key-role')}
                </p>
                <Icon
                  icon={IconInfo}
                  size="xs"
                  color="gray-500"
                  className="mt-0.5"
                />
              </div>

              <RadioGroup defaultValue={apiKey?.role} disabled={disabled}>
                {apiKeyOptions.map(({ id, label, description, value }) => (
                  <div
                    key={id}
                    className="flex items-center last:border-b-0 border-b py-4 gap-x-5"
                  >
                    <label htmlFor={id} className="flex-1 opacity-50">
                      <p className="typo-para-medium text-gray-700">{label}</p>
                      <p className="typo-para-small text-gray-600">
                        {description}
                      </p>
                    </label>
                    <RadioGroupItem disabled value={value} id={id} />
                  </div>
                ))}
              </RadioGroup>

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

export default EditAPIKeyModal;
