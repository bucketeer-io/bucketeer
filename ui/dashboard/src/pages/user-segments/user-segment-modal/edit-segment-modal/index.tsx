import { useMemo, useState } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { userSegmentBulkUpload } from '@api/user-segment';
import { userSegmentUpdater } from '@api/user-segment/user-segment-updater';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateUserSegments } from '@queries/user-segments';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { useTranslation } from 'i18n';
import { UserSegment } from '@types';
import { covertFileToUint8ToBase64 } from 'utils/converts';
import { cn } from 'utils/style';
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
import SegmentWarning from './segment-warning';

interface EditUserSegmentModalProps {
  userSegment: UserSegment;
  isOpen: boolean;
  onClose: () => void;
}

const EditUserSegmentModal = ({
  userSegment,
  isOpen,
  onClose
}: EditUserSegmentModalProps) => {
  const { t } = useTranslation(['common', 'form']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const queryClient = useQueryClient();
  const { notify } = useToast();

  const isDisabledUserIds = useMemo(
    () => userSegment.isInUseStatus || userSegment.features?.length > 0,
    [userSegment]
  );

  const [userIdsType, setUserIdsType] = useState(
    isDisabledUserIds ? '' : 'upload'
  );
  const [files, setFiles] = useState<File[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      id: userSegment?.id || '',
      name: userSegment?.name || '',
      description: userSegment?.description || '',
      userIds: '',
      file: null
    }
  });

  const {
    formState: { isValid, isDirty, isSubmitting },
    trigger
  } = form;

  const updateSuccess = (name: string) => {
    notify({
      toastType: 'toast',
      messageType: 'success',
      message: (
        <span>
          <b>{name}</b> {` has been successfully updated!`}
        </span>
      )
    });
    invalidateUserSegments(queryClient);
    setIsLoading(false);
    onClose();
  };

  const onUpdateSuccess = (name: string, isUpload = true) => {
    let timerId: NodeJS.Timeout | null = null;
    if (timerId) clearTimeout(timerId);
    if (isUpload)
      return (timerId = setTimeout(() => updateSuccess(name), 10000));
    return updateSuccess(name);
  };

  const onSubmit: SubmitHandler<UserSegmentForm> = async values => {
    try {
      setIsLoading(true);
      const { id, name, description } = values;
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
          await userSegmentBulkUpload({
            segmentId: id as string,
            environmentId: currentEnvironment.id,
            state: 'INCLUDED',
            data: base64String
          });
        });
      }
      const resp = await userSegmentUpdater({
        id: id as string,
        name,
        description,
        environmentId: currentEnvironment.id
      });
      if (resp?.segment) onUpdateSuccess(name, !!file);
    } catch (error) {
      notify({
        toastType: 'toast',
        messageType: 'error',
        message: (error as Error)?.message || 'Something went wrong.'
      });
    }
  };

  return (
    <SlideModal
      title={t('update-user-segment')}
      isOpen={isOpen}
      shouldCloseOnOverlayClick={!isLoading}
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
              defaultValue={userIdsType}
              onValueChange={setUserIdsType}
              disabled={isDisabledUserIds}
              className="flex flex-col w-full gap-y-4"
            >
              <Form.Field
                control={form.control}
                name="file"
                render={({ field }) => (
                  <Form.Item
                    className={cn('py-0', { 'opacity-50': isDisabledUserIds })}
                  >
                    <Form.Control>
                      <div className="flex flex-col w-full gap-y-3">
                        <div className="flex items-center gap-x-2">
                          <RadioGroupItem
                            value={'upload'}
                            checked={userIdsType === 'upload'}
                            id={'upload'}
                          />
                          <label
                            htmlFor={'upload'}
                            className={cn(
                              'cursor-pointer typo-para-small text-gray-700',
                              { 'cursor-not-allowed': isDisabledUserIds }
                            )}
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
                  <Form.Item
                    className={cn('py-0', { 'opacity-50': isDisabledUserIds })}
                  >
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
                            className={cn(
                              'cursor-pointer typo-para-small text-gray-700',
                              { 'cursor-not-allowed': isDisabledUserIds }
                            )}
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

            {isDisabledUserIds && (
              <SegmentWarning features={userSegment.features} />
            )}

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
                    disabled={!isDirty || !isValid || isLoading || isSubmitting}
                    loading={isSubmitting || isLoading}
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
