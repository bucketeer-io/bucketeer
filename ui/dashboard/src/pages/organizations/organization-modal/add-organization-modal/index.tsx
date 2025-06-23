import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { organizationCreator } from '@api/organization';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateOrganizations } from '@queries/organizations';
import { useQueryClient } from '@tanstack/react-query';
import { requiredMessage, translation } from 'constants/message';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { onGenerateSlug } from 'utils/converts';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Checkbox from 'components/checkbox';
import Form from 'components/form';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import TextArea from 'components/textarea';

interface AddOrganizationModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export interface AddOrganizationForm {
  name: string;
  urlCode: string;
  ownerEmail: string;
  isTrial?: boolean;
  description?: string;
}

const formSchema = yup.object().shape({
  name: yup.string().required(requiredMessage),
  urlCode: yup
    .string()
    .required(requiredMessage)
    .matches(
      /^[a-zA-Z0-9][a-zA-Z0-9-]*$/,
      translation('message:validation.id-rule', {
        name: translation('common:url-code')
      })
    ),
  description: yup.string(),
  ownerEmail: yup.string().email().required(requiredMessage),
  isTrial: yup.bool()
});

const AddOrganizationModal = ({
  isOpen,
  onClose
}: AddOrganizationModalProps) => {
  const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form', 'message']);
  const { notify, errorNotify } = useToast();

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      name: '',
      urlCode: '',
      description: '',
      ownerEmail: '',
      isTrial: true
    }
  });

  const onSubmit: SubmitHandler<AddOrganizationForm> = async values => {
    return organizationCreator({
      ...values,
      isSystemAdmin: false
    })
      .then(() => {
        notify({
          message: t('message:collection-action-success', {
            collection: t('organization'),
            action: t('created')
          })
        });
        invalidateOrganizations(queryClient);
        onClose();
      })
      .catch(error => errorNotify(error));
  };

  return (
    <SlideModal title={t('new-org')} isOpen={isOpen} onClose={onClose}>
      <div className="w-full p-5">
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
                      onChange={value => {
                        const isUrlCodeDirty =
                          form.getFieldState('urlCode').isDirty;
                        const urlCode = form.getValues('urlCode');
                        field.onChange(value);
                        form.setValue(
                          'urlCode',
                          isUrlCodeDirty ? urlCode : onGenerateSlug(value)
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
              name="urlCode"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required>{t('form:url-code')}</Form.Label>
                  <Form.Control>
                    <Input
                      placeholder={`${t('form:placeholder-code')}`}
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
            <Form.Field
              control={form.control}
              name="ownerEmail"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required>{t('form:owner-email')}</Form.Label>
                  <Form.Control className="w-full">
                    <Input
                      placeholder={`${t('form:placeholder-email')}`}
                      {...field}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name="isTrial"
              render={({ field }) => (
                <Form.Item>
                  <Form.Control>
                    <Checkbox
                      onCheckedChange={checked => field.onChange(checked)}
                      checked={field.value}
                      title={`${t(`form:trial`)}`}
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
                    {t(`common:cancel`)}
                  </Button>
                }
                secondaryButton={
                  <Button
                    type="submit"
                    disabled={!form.formState.isDirty}
                    loading={form.formState.isSubmitting}
                  >
                    {t(`create-org`)}
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

export default AddOrganizationModal;
