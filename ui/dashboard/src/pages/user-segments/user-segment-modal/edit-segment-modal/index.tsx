import { useState } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { yupResolver } from '@hookform/resolvers/yup';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import { UserSegment } from '@types';
import { IconToastWarning } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import { RadioGroup, RadioGroupItem } from 'components/radio';
import TextArea from 'components/textarea';
import Upload from 'components/upload-files';

interface EditUserSegmentModalProps {
  userSegment: UserSegment;
  isOpen: boolean;
  onClose: () => void;
}

export interface EditUserSegmentForm {
  id?: string;
  name: string;
  description?: string;
  userIds?: string;
  files?: File[];
}

const formSchema = yup.object().shape({
  name: yup.string().required(),
  description: yup.string(),
  id: yup.string(),
  userIds: yup.string(),
  files: yup.array()
});

const EditUserSegmentModal = ({
  userSegment,
  isOpen,
  onClose
}: EditUserSegmentModalProps) => {
  const { t } = useTranslation(['common', 'form']);
  const [userIdsType, setUserIdsType] = useState('typing');
  const [files, setFiles] = useState<File[]>([]);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      id: userSegment?.id || '',
      name: userSegment?.name || '',
      description: userSegment?.description || '',
      userIds: '',
      files: []
    }
  });

  const {
    setValue
    // formState: { isValid, isSubmitting }
  } = form;

  const onSubmit: SubmitHandler<EditUserSegmentForm> = values => {
    console.log({ values });
  };

  return (
    <SlideModal
      title={t('update-user-segment')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="w-full p-5 pb-28">
        <p className="text-gray-600 typo-para-medium mb-4">
          {t('form:update-user-segment')}
        </p>
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
            <Divider className="mt-1 mb-5" />
            <p className="text-gray-900 typo-head-bold-small mb-5">{`${t('form:list-of-users-ids')} (${t('form:optional')})`}</p>
            <RadioGroup
              defaultValue="typing"
              onValueChange={value => setUserIdsType(value)}
              className="flex flex-col w-full gap-y-4"
            >
              <div className="flex flex-col w-full gap-y-3">
                <div className="flex items-center gap-x-2">
                  <RadioGroupItem
                    value={'upload'}
                    checked={userIdsType === 'upload'}
                    id={'upload'}
                  />
                  <label
                    htmlFor={'upload'}
                    className="cursor-pointer typo-para-small text-gray-700"
                  >
                    {t('form:browse-files')}
                  </label>
                </div>
                {userIdsType === 'upload' && (
                  <div className="flex w-full max-w-full h-fit gap-x-4 pl-8">
                    <Upload
                      files={files}
                      className="border-l border-primary-500 pl-4"
                      uploadClassName="min-h-[200px] h-[200px]"
                      onChange={files => setFiles(files)}
                    />
                  </div>
                )}
              </div>
              <div className="flex flex-col w-full gap-y-3">
                <div className="flex items-center gap-x-2">
                  <RadioGroupItem
                    id={'typing'}
                    value={'typing'}
                    checked={userIdsType === 'typing'}
                  />
                  <label
                    htmlFor={'typing'}
                    className="cursor-pointer typo-para-small text-gray-700"
                  >
                    {t('form:enter-user-ids')}
                  </label>
                </div>
                {userIdsType === 'typing' && (
                  <div className="flex w-full max-w-full h-fit gap-x-4 pl-8">
                    <Divider
                      vertical
                      width={1}
                      className="border-primary-500 !h-[120px]"
                    />
                    <TextArea
                      placeholder={t('form:placeholder-enter-user-ids')}
                      rows={4}
                      onChange={e => setValue('userIds', e.target.value)}
                    />
                  </div>
                )}
              </div>
            </RadioGroup>

            <div className="flex flex-col w-full px-4 py-3 bg-accent-yellow-50 border-l-4 border-accent-yellow-500 rounded mt-5">
              <div className="flex gap-x-2 w-full pr-3">
                <Icon
                  icon={IconToastWarning}
                  size={'xxs'}
                  color="accent-yellow-500"
                  className="mt-1"
                />
                <Trans
                  i18nKey="form:update-user-segment-warning"
                  values={{ count: 1 }}
                  components={{
                    p: <p className="typo-para-medium text-accent-yellow-500" />
                  }}
                />
              </div>
              <div className="flex gap-x-2 w-full pl-6 typo-para-medium text-primary-500">
                <p>1.</p>
                <p className="hover:underline">{userSegment.name}</p>
              </div>
            </div>

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

export default EditUserSegmentModal;
