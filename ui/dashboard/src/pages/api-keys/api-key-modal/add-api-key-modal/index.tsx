import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { apiKeyCreator } from '@api/api-key';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateAPIKeys } from '@queries/api-keys';
import { useQueryClient } from '@tanstack/react-query';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { APIKeyRole, Environment } from '@types';
import { checkEnvironmentEmptyId } from 'utils/function';
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

interface AddAPIKeyModalProps {
  isOpen: boolean;
  environments: Environment[];
  isLoadingEnvs: boolean;
  onClose: () => void;
}

export interface AddAPIKeyForm {
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

const AddAPIKeyModal = ({
  isOpen,
  isLoadingEnvs,
  environments,
  onClose
}: AddAPIKeyModalProps) => {
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form']);
  const { notify } = useToast();

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      name: '',
      environmentId: '',
      description: '',
      role: 'SDK_CLIENT'
    }
  });

  const {
    getValues,
    formState: { isValid, isSubmitting }
  } = form;

  const onSubmit: SubmitHandler<AddAPIKeyForm> = values => {
    return apiKeyCreator({
      environmentId: checkEnvironmentEmptyId(values.environmentId),
      name: values.name,
      role: values.role,
      description: values.description
    }).then(() => {
      notify({
        toastType: 'toast',
        messageType: 'success',
        message: (
          <span>
            <b>{values.name}</b> {` has been successfully created!`}
          </span>
        )
      });
      invalidateAPIKeys(queryClient);
      onClose();
    });
  };

  return (
    <SlideModal title={t('new-api-key')} isOpen={isOpen} onClose={onClose}>
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
                          )?.name || ''
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
            <Form.Field
              control={form.control}
              name="role"
              render={({ field }) => (
                <Form.Item>
                  <Form.Control>
                    <RadioGroup
                      defaultValue={field.value}
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
                              className="flex-1 cursor-pointer"
                            >
                              <p className="typo-para-medium text-gray-700">
                                {label}
                              </p>
                              <p className="typo-para-small text-gray-600">
                                {description}
                              </p>
                            </label>

                            <RadioGroupItem value={value} id={id} />
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

export default AddAPIKeyModal;
