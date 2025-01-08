import { useState } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Form from 'components/form';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import { RadioGroup, RadioGroupItem } from 'components/radio';
import TextArea from 'components/textarea';
import Upload from 'components/upload-files';

interface AddUserSegmentModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export interface AddUserSegmentForm {
  id?: string;
  name: string;
  description?: string;
  userIds?: string;
}

const formSchema = yup.object().shape({
  name: yup.string().required(),
  description: yup.string(),
  id: yup.string(),
  userIds: yup.string()
});

const AddUserSegmentModal = ({ isOpen, onClose }: AddUserSegmentModalProps) => {
  const { t } = useTranslation(['common', 'form']);
  const [userIdsType, setUserIdsType] = useState('upload');
  const [files, setFiles] = useState<File[]>([]);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      id: '',
      name: '',
      description: '',
      userIds: ''
    }
  });

  const {
    setValue
    // formState: { isValid, isSubmitting }
  } = form;

  const onSubmit: SubmitHandler<AddUserSegmentForm> = values => {
    console.log({ values });
  };

  return (
    <SlideModal title={t('new-user-segment')} isOpen={isOpen} onClose={onClose}>
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
              defaultValue="upload"
              onValueChange={value => setUserIdsType(value)}
              className="flex flex-col w-full gap-y-4"
            >
              <div className="flex flex-col w-full h-fit gap-y-3">
                <div className="flex items-center h-fit gap-x-2">
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
                  <div className="flex w-full max-w-full h-fit pl-8">
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
                    {t(`create-user-segment`)}
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

export default AddUserSegmentModal;
