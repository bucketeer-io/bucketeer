import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { pushCreator } from '@api/push';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidatePushes } from '@queries/pushes';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { IconInfo } from '@icons';
import { useFetchEnvironments } from 'pages/project-details/environments/collection-loader/use-fetch-environments';
import {
  renderTag,
  tagOptions
} from 'pages/push/collection-layout/data-collection';
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
import TextArea from 'components/textarea';

interface AddPushModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export interface AddPushForm {
  name: string;
  fcmServiceAccount: string;
  tags: string[];
  environmentId: string;
}

export const formSchema = yup.object().shape({
  name: yup.string().required(),
  fcmServiceAccount: yup.string().required(),
  tags: yup.array().required(),
  environmentId: yup.string().required()
});

const AddPushModal = ({ isOpen, onClose }: AddPushModalProps) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form']);
  const { notify } = useToast();

  const { data: collection, isLoading: isLoadingEnvs } = useFetchEnvironments({
    organizationId: currentEnvironment.organizationId
  });
  const environments = (collection?.environments || []).filter(item => item.id);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      name: '',
      fcmServiceAccount: '',
      tags: [],
      environmentId: ''
    }
  });

  const {
    getValues,
    formState: { isValid, isSubmitting }
  } = form;

  const onSubmit: SubmitHandler<AddPushForm> = async values => {
    try {
      const resp = await pushCreator(values);
      if (resp) {
        notify({
          toastType: 'toast',
          messageType: 'success',
          message: (
            <span>
              <b>{values.name}</b> {` has been successfully created!`}
            </span>
          )
        });
        invalidatePushes(queryClient);
        onClose();
      }
    } catch (error) {
      console.log(error);
    }
  };

  return (
    <SlideModal title={t('new-push')} isOpen={isOpen} onClose={onClose}>
      <div className="w-full p-5 pb-28">
        <div className="typo-para-small text-gray-600 mb-3">
          {t('new-push-subtitle')}
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
              name="fcmServiceAccount"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required className="relative w-fit">
                    <>
                      {t('fcm-api-key')}
                      <Icon
                        icon={IconInfo}
                        className="absolute -right-8"
                        size={'sm'}
                      />
                    </>
                  </Form.Label>
                  <Form.Control>
                    <TextArea
                      placeholder={`${t('form:placeholder-firebase')}`}
                      rows={3}
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
                <Form.Item className="py-2">
                  <Form.Label required>{t('tags')}</Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        placeholder={t(`form:placeholder-tags`)}
                        variant="secondary"
                        className="w-full"
                        trigger={
                          field.value?.length > 0 && renderTag(field.value)
                        }
                      />
                      <DropdownMenuContent
                        className="w-[502px]"
                        align="start"
                        {...field}
                      >
                        {tagOptions.map((item, index) => (
                          <DropdownMenuItem
                            {...field}
                            key={index}
                            value={item.value}
                            label={item.label}
                            isMultiselect={true}
                            isSelected={field.value.includes(item.value)}
                            onSelectOption={value => {
                              const _tags = field.value.includes(value)
                                ? field.value.filter(tag => tag !== value)
                                : [...field.value, value];
                              field.onChange(_tags);
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
                    {t(`submit`)}
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

export default AddPushModal;
