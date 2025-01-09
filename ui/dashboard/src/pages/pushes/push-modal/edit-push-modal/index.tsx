import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { pushUpdater } from '@api/push';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidatePushes } from '@queries/pushes';
import { useQueryTags } from '@queries/tags';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { LIST_PAGE_SIZE } from 'constants/app';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { Push } from '@types';
import { useFetchEnvironments } from 'pages/project-details/environments/collection-loader/use-fetch-environments';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import { CreatableSelect } from 'components/creatable-select';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';

interface EditPushModalProps {
  isOpen: boolean;
  onClose: () => void;
  push: Push;
}

export interface EditPushForm {
  name: string;
  fcmServiceAccount?: string;
  tags: string[];
}

const formSchema = yup.object().shape({
  name: yup.string().required(),
  fcmServiceAccount: yup.string(),
  tags: yup.array().required(),
  environmentId: yup.string().required()
});

const EditPushModal = ({ isOpen, onClose, push }: EditPushModalProps) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form']);
  const { notify } = useToast();

  const { data: collection } = useFetchEnvironments({
    organizationId: currentEnvironment.organizationId
  });

  const { data: tagCollection, isLoading: isLoadingTags } = useQueryTags({
    params: {
      cursor: String(0),
      pageSize: LIST_PAGE_SIZE,
      environmentId: currentEnvironment.id
    }
  });
  const tagOptions = tagCollection?.tags || [];
  const environments = (collection?.environments || []).filter(item => item.id);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      name: push.name,
      fcmServiceAccount: push.fcmServiceAccount,
      tags: push.tags,
      environmentId: push.environmentId
    }
  });

  const {
    getValues,
    formState: { isValid, isSubmitting, isDirty }
  } = form;

  const onSubmit: SubmitHandler<EditPushForm> = async values => {
    await pushUpdater({
      ...values,
      id: push.id
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
      invalidatePushes(queryClient);
      onClose();
    });
  };

  return (
    <SlideModal title={t('edit-push')} isOpen={isOpen} onClose={onClose}>
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
                        disabled
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
            <Form.Field
              control={form.control}
              name={`tags`}
              render={({ field }) => (
                <Form.Item className="py-2">
                  <Form.Label required>{t('tags')}</Form.Label>
                  <Form.Control>
                    <CreatableSelect
                      defaultValues={field.value?.map(tag => ({
                        label: tag,
                        value: tag
                      }))}
                      disabled={isLoadingTags}
                      placeholder={t(`form:placeholder-tags`)}
                      options={tagOptions?.map(tag => ({
                        label: tag.id,
                        value: tag.id
                      }))}
                      onChange={value =>
                        field.onChange(value.map(tag => tag.value))
                      }
                    />
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
                    disabled={!isValid || !isDirty}
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

export default EditPushModal;
