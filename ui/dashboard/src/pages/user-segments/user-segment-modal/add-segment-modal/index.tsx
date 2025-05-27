import { useState } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { userSegmentBulkUpload, userSegmentCreator } from '@api/user-segment';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateUserSegments } from '@queries/user-segments';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { covertFileToUint8ToBase64 } from 'utils/converts';
import { UserSegmentForm } from 'pages/user-segments/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Form from 'components/form';
import Input from 'components/input';
import SlideModal from 'components/modal/slide';
import { RadioGroup, RadioGroupItem } from 'components/radio';
import TextArea from 'components/textarea';
import Upload from 'components/upload-files';
import { formSchema } from '../../form-schema';

interface AddUserSegmentModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const AddUserSegmentModal = ({ isOpen, onClose }: AddUserSegmentModalProps) => {
  const { t } = useTranslation(['common', 'form', 'message']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { notify, errorNotify } = useToast();
  const queryClient = useQueryClient();

  const [userIdsType, setUserIdsType] = useState('upload');
  const [files, setFiles] = useState<File[]>([]);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      id: '',
      name: '',
      description: '',
      userIds: '',
      file: null
    }
  });

  const {
    formState: { isValid, isDirty, isSubmitting },
    trigger
  } = form;

  const addSuccess = (isUpload = false) => {
    if (!isUpload) {
      notify({
        message: t('message:collection-action-success', {
          collection: t('source-type.segment'),
          action: t('created').toLowerCase()
        })
      });
      onClose();
    }
    invalidateUserSegments(queryClient);
  };

  const onAddSuccess = (isUpload = false) => {
    let timerId: NodeJS.Timeout | null = null;
    if (timerId) clearTimeout(timerId);
    if (isUpload)
      return (timerId = setTimeout(() => addSuccess(isUpload), 10000));
    return addSuccess(isUpload);
  };

  const onSubmit: SubmitHandler<UserSegmentForm> = async values => {
    try {
      const resp = await userSegmentCreator({
        environmentId: currentEnvironment.id,
        name: values.name,
        description: values.description
      });

      if (resp) {
        let file: File | null = null;
        if (values.file || files.length) {
          file = (values.file as File) || files[0];
        } else if (values.userIds?.length) {
          file = new File([values.userIds], 'filename.txt', {
            type: 'text/plain'
          });
        }
        if (file) {
          covertFileToUint8ToBase64(file, async base64String => {
            const uploadResp = await userSegmentBulkUpload({
              segmentId: resp.segment.id,
              environmentId: currentEnvironment.id,
              state: 'INCLUDED',
              data: base64String
            });
            onAddSuccess();
            if (uploadResp) onAddSuccess(true);
          });
        } else onAddSuccess(false);
      }
    } catch (error) {
      errorNotify(error);
    }
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
              defaultValue={userIdsType}
              onValueChange={setUserIdsType}
              className="flex flex-col w-full gap-y-4"
            >
              <Form.Field
                control={form.control}
                name="file"
                render={({ field }) => (
                  <Form.Item className="py-0">
                    <Form.Control>
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
                              onChange={files => {
                                setFiles(files);
                                field.onChange(files?.length ? files[0] : null);
                                trigger('file');
                              }}
                            />
                          </div>
                        )}
                      </div>
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
              <Form.Field
                control={form.control}
                name="userIds"
                render={({ field }) => (
                  <Form.Item className="py-0">
                    <Form.Control>
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
                              onChange={e => field.onChange(e.target.value)}
                            />
                          </div>
                        )}
                      </div>
                    </Form.Control>
                    <Form.Message />
                  </Form.Item>
                )}
              />
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
                    disabled={!isDirty || !isValid || isSubmitting}
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

export default AddUserSegmentModal;
