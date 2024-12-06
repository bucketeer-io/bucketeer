import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { apiKeyCreator } from '@api/api-key';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateAPIKeys } from '@queries/api-keys';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { APIKeyRole } from '@types';
import { IconInfo } from '@icons';
import { useFetchEnvironments } from 'pages/project-details/environments/collection-loader/use-fetch-environments';
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

interface AddNotificationModalProps {
  isOpen: boolean;
  onClose: () => void;
}

type APIKeyOption = {
  id: string;
  label: string;
  description: string;
  value: APIKeyRole;
};

export interface AddNotificationForm {
  name: string;
  description?: string;
  environmentId: string;
  role: APIKeyRole;
}

export const formSchema = yup.object().shape({
  name: yup.string().required(),
  environmentId: yup.string().required(),
  description: yup.string(),
  role: yup.mixed<APIKeyRole>().required()
});

const AddNotificationModal = ({ isOpen, onClose }: AddNotificationModalProps) => {
  const { consoleAccount } = useAuth();
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form']);
  const { notify } = useToast();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { data: collection, isLoading: isLoadingEnvs } = useFetchEnvironments({
    organizationId: currentEnvironment.organizationId
  });
  const environments = (collection?.environments || []).filter(item => item.id);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      name: '',
      environmentId: '',
      description: '',
      role: 'SDK_CLIENT'
    }
  });

  const options: APIKeyOption[] = [
    {
      id: 'client-sdk',
      label: t('form:api-key.client-sdk'),
      description: t('form:api-key.client-sdk-desc'),
      value: 'SDK_CLIENT'
    },
    {
      id: 'server-sdk',
      label: t('form:api-key.server-sdk'),
      description: t('form:api-key.server-sdk-desc'),
      value: 'SDK_SERVER'
    },
    {
      id: 'public-api-read-only',
      label: t('form:api-key.public-api-read-only'),
      description: t('form:api-key.public-api-read-only-desc'),
      value: 'PUBLIC_API_READ_ONLY'
    },
    {
      id: 'public-api-write',
      label: t('form:api-key.public-api-write'),
      description: t('form:api-key.public-api-write-desc'),
      value: 'PUBLIC_API_WRITE'
    },
    {
      id: 'public-api-admin',
      label: t('form:api-key.public-api-admin'),
      description: t('form:api-key.public-api-admin-desc'),
      value: 'PUBLIC_API_ADMIN'
    }
  ];

  const {
    getValues,
    formState: { isValid, isSubmitting }
  } = form;

  const onSubmit: SubmitHandler<AddNotificationForm> = values => {
    return apiKeyCreator({
      environmentId: values.environmentId,
      command: {
        name: values.name,
        role: values.role
      }
    }).then(() => {
      notify({
        toastType: 'toast',
        messageType: 'success',
        message: (
          <span>
            <b>{values.name}</b> {` has been successfully created!`}{' '}
          </span>
        )
      });
      invalidateAPIKeys(queryClient);
      onClose();
    });
  };

  return (
    <SlideModal title={t('new-notification')} isOpen={isOpen} onClose={onClose}>
      <div className="w-full p-5 pb-28">
        <div className="typo-para-small text-gray-600 mb-1">
          {t('new-notification-subtitle')} 
        </div>
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
                        label={
                          environments.find(
                            item => item.id === getValues('environmentId')
                          )?.name
                        }
                        disabled={isLoadingEnvs}
                        variant="secondary"
                        className="w-full"
                      />
                      <DropdownMenuContent
                        className="w-[502px]"
                        align="start"
                        {...field}
                      >
                        {environments.map((item, index) => (
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

            <RadioGroup defaultValue={getValues('role')}>
              {options.map(({ id, label, description, value }) => (
                <div
                  key={id}
                  className="flex items-center last:border-b-0 border-b py-4 gap-x-5"
                >
                  <label htmlFor={id} className="flex-1 cursor-pointer">
                    <p className="typo-para-medium text-gr ay-700">{label}</p>
                    <p className="typo-para-small text-gray-600">
                      {description}
                    </p>
                  </label>
                  <Form.Field
                    control={form.control}
                    name="role"
                    render={({ field }) => (
                      <Form.Item>
                        <Form.Control>
                          <RadioGroupItem
                            onChange={() => field.onChange(value)}
                            value={value}
                            id={id}
                          />
                        </Form.Control>
                      </Form.Item>
                    )}
                  />
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
                  <Button
                    type="submit"
                    disabled={!isValid}
                    loading={isSubmitting}
                  >
                    {t(`create-api-key`)}
                  </Button>
                }
              />
            </div>
          </Form>
        </FormProvider>
      </div>
    </SlideModal>
  );
};

export default AddNotificationModal;
