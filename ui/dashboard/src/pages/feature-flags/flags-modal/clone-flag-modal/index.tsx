import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { useQueryEnvironments } from '@queries/environments';
// import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { LIST_PAGE_SIZE } from 'constants/app';
// import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';

interface CloneFlagModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export interface CloneFlagForm {
  name: string;
  originEnvironmentId: string;
  destinationEnvironmentId: string;
}

const formSchema = yup.object().shape({
  name: yup.string().required(),
  originEnvironmentId: yup.string().required(),
  destinationEnvironmentId: yup.string().required()
});

const CloneFlagModal = ({ isOpen, onClose }: CloneFlagModalProps) => {
  const { consoleAccount } = useAuth();
  //   const queryClient = useQueryClient();
  const { t } = useTranslation(['common', 'form']);
  //   const { notify } = useToast();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { data: collection, isLoading: isLoadingEnvs } = useQueryEnvironments({
    params: {
      organizationId: currentEnvironment.organizationId,
      pageSize: LIST_PAGE_SIZE,
      cursor: String(0)
    }
  });

  const environments = collection?.environments || [];

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      name: '',
      originEnvironmentId: '',
      destinationEnvironmentId: ''
    }
  });

  const onSubmit: SubmitHandler<CloneFlagForm> = values => {
    console.log(values);
  };

  return (
    <SlideModal
      title={t('form:feature-flags.clone-title')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="w-full p-5">
        <p className="text-gray-600 typo-para-small">
          {t('form:feature-flags.clone-desc')}
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
              name={`originEnvironmentId`}
              render={({ field }) => (
                <Form.Item className="py-2">
                  <Form.Label required>{t('form:origin-env')}</Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        placeholder={t(`form:select-environment`)}
                        label={
                          environments.find(item => item.id === field.value)
                            ?.name
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
              name={`destinationEnvironmentId`}
              render={({ field }) => (
                <Form.Item className="py-2">
                  <Form.Label required>{t('form:destination-env')}</Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        placeholder={t(`form:select-environment`)}
                        label={
                          environments.find(item => item.id === field.value)
                            ?.name
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
                    disabled={!form.formState.isDirty}
                    loading={form.formState.isSubmitting}
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

export default CloneFlagModal;
