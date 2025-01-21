import {
  ControllerRenderProps,
  FormProvider,
  SubmitHandler,
  useForm
} from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { IconInfo } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
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
import DefineAudience from './define-audience';

interface AddExperimentModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export interface AddExperimentForm {
  id: string;
  name: string;
  description?: string;
  audience: {
    inExperiment: number;
    notInExperiment: number;
    served: boolean;
    variationReassignment: boolean;
  };
  linkFlag: string;
  linkGoal: string;
}

export type DefineAudienceField = ControllerRenderProps<
  AddExperimentForm,
  'audience'
>;

export const formSchema = yup.object().shape({
  id: yup.string().required(),
  name: yup.string().required(),
  description: yup.string(),
  audience: yup.mixed(),
  linkFlag: yup.string().required(),
  linkGoal: yup.string().required()
});

const AddExperimentModal = ({ isOpen, onClose }: AddExperimentModalProps) => {
  const { t } = useTranslation(['form', 'common']);
  const { notify } = useToast();

  const flagOptions = [
    {
      label: 'Flag 1',
      value: 'flag-1'
    },
    {
      label: 'Flag 2',
      value: 'flag-2'
    }
  ];

  const goalOptions = [
    {
      label: 'Goal 1',
      value: 'goal-1'
    },
    {
      label: 'Goal 2',
      value: 'goal-2'
    }
  ];

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      id: '',
      name: '',
      description: '',
      audience: {
        inExperiment: 5,
        notInExperiment: 95,
        served: true,
        variationReassignment: false
      },
      linkFlag: '',
      linkGoal: ''
    }
  });

  const {
    formState: { isValid, isSubmitting }
  } = form;

  const onSubmit: SubmitHandler<AddExperimentForm> = async values => {
    try {
      console.log(values);
    } catch (error) {
      const errorMessage = (error as Error)?.message;
      notify({
        toastType: 'toast',
        messageType: 'error',
        message: errorMessage || 'Something went wrong.'
      });
    }
  };

  return (
    <SlideModal
      title={t('common:new-experiment')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="w-[542px] p-5 pb-28">
        <p className="text-gray-800 typo-head-bold-small">
          {t('general-info')}
        </p>
        <FormProvider {...form}>
          <Form onSubmit={form.handleSubmit(onSubmit)}>
            <Form.Field
              control={form.control}
              name="name"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required>{t('common:name')}</Form.Label>
                  <Form.Control>
                    <Input
                      placeholder={`${t('placeholder-name')}`}
                      {...field}
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name="id"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required className="relative w-fit">
                    {t('experiments.experiment-id')}
                    <Icon
                      icon={IconInfo}
                      className="absolute -right-8"
                      size={'sm'}
                    />
                  </Form.Label>
                  <Form.Control>
                    <Input
                      placeholder={`${t('experiments.placeholder-id')}`}
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
            <Divider className="mb-5" />
            <div className="mb-3">
              <p className="text-gray-800 typo-head-bold-small mb-3">
                {t('experiments.define-audience.title')}
              </p>
              <Form.Field
                control={form.control}
                name={`audience`}
                render={({ field }) => (
                  <Form.Item className="flex flex-col w-full py-2 gap-y-5">
                    <DefineAudience field={field as DefineAudienceField} />
                  </Form.Item>
                )}
              />
            </div>
            <Divider className="mb-5" />
            <p className="text-gray-800 typo-head-bold-small mb-3">
              {t('link')}
            </p>
            <Form.Field
              control={form.control}
              name={`linkFlag`}
              render={({ field }) => (
                <Form.Item className="py-2">
                  <Form.Label required>{t('experiments.link-flag')}</Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        placeholder={t(`experiments.select-flag`)}
                        label={''}
                        variant="secondary"
                        className="w-full"
                      />
                      <DropdownMenuContent
                        className="w-[502px]"
                        align="start"
                        {...field}
                      >
                        {flagOptions.map((item, index) => (
                          <DropdownMenuItem
                            {...field}
                            key={index}
                            value={item.value}
                            label={item.label}
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
              name={`linkGoal`}
              render={({ field }) => (
                <Form.Item className="py-2">
                  <Form.Label required>{t('experiments.link-goal')}</Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        placeholder={t(`experiments.select-goal`)}
                        label={''}
                        variant="secondary"
                        className="w-full"
                      />
                      <DropdownMenuContent
                        className="w-[502px]"
                        align="start"
                        {...field}
                      >
                        {goalOptions.map((item, index) => (
                          <DropdownMenuItem
                            {...field}
                            key={index}
                            value={item.value}
                            label={item.label}
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
                    {t(`common:cancel`)}
                  </Button>
                }
                secondaryButton={
                  <Button
                    type="submit"
                    disabled={!isValid}
                    loading={isSubmitting}
                  >
                    {t(`common:submit`)}
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

export default AddExperimentModal;
