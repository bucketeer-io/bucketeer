import { useCallback, useMemo, useState } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { userSegmentBulkUpload, userSegmentCreator } from '@api/user-segment';
import { userSegmentUpdater } from '@api/user-segment/user-segment-updater';
import { yupResolver } from '@hookform/resolvers/yup';
import { invalidateUserSegments } from '@queries/user-segments';
import { useQueryClient } from '@tanstack/react-query';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import useFormSchema from 'hooks/use-form-schema';
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
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import FormLoading from 'elements/form-loading';
import { formSchema } from '../../form-schema';
import SegmentWarning from './segment-warning';

interface SegmentCreateUpdateModalProps {
  isUpdate: boolean;
  isOpen: boolean;
  isLoadingSegment?: boolean;
  userSegment?: UserSegment;
  isDisabled: boolean;
  onClose: () => void;
  setSegmentUploading: (userSegment: UserSegment | null) => void;
}

const SegmentCreateUpdateModal = ({
  isUpdate,
  isOpen,
  isLoadingSegment,
  userSegment,
  isDisabled,
  onClose,
  setSegmentUploading
}: SegmentCreateUpdateModalProps) => {
  const { t } = useTranslation(['common', 'form', 'message']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const queryClient = useQueryClient();
  const { notify, errorNotify } = useToast();

  const isDisabledUserIds = useMemo(
    () =>
      userSegment &&
      (userSegment?.isInUseStatus || userSegment?.features?.length > 0),
    [userSegment]
  );

  const [userIdsType, setUserIdsType] = useState('upload');
  const [files, setFiles] = useState<File[]>([]);

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    values: {
      id: userSegment?.id || '',
      name: userSegment?.name || '',
      description: userSegment?.description || '',
      userIds: '',
      file: null
    },
    mode: 'onChange'
  });

  const {
    formState: { isValid, isDirty, isSubmitting },
    getFieldState
  } = form;

  const updateSuccess = (isUpload = false) => {
    if (!isUpload) {
      notify({
        message: t('message:collection-action-success', {
          collection: t('source-type.segment'),
          action: t(userSegment ? 'updated' : 'created')
        })
      });
      onClose();
    }
    if (isUpload) setSegmentUploading(null);
    invalidateUserSegments(queryClient);
  };

  const onUpdateSuccess = (isUpload = false) => {
    let timerId: NodeJS.Timeout | null = null;
    if (timerId) clearTimeout(timerId);
    if (isUpload)
      return (timerId = setTimeout(() => updateSuccess(isUpload), 10000));
    return updateSuccess(isUpload);
  };

  const onSubmit: SubmitHandler<UserSegmentForm> = useCallback(
    async values => {
      try {
        const { id, name, description } = values;
        let segmentId = id;
        let newSegment = userSegment;
        if (!userSegment) {
          const resp = await userSegmentCreator({
            environmentId: currentEnvironment.id,
            name: values.name,
            description: values.description
          });
          segmentId = resp.segment.id;
          newSegment = resp.segment;
        }

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
              segmentId: segmentId as string,
              environmentId: currentEnvironment.id,
              state: 'INCLUDED',
              data: base64String
            });
            onUpdateSuccess(true);
          });
        }

        if (userSegment) {
          if (
            getFieldState('name').isDirty ||
            getFieldState('description').isDirty
          ) {
            await userSegmentUpdater({
              id: segmentId as string,
              name,
              description,
              environmentId: currentEnvironment.id
            });
          }
        }
        if (file) setSegmentUploading(newSegment!);
        onUpdateSuccess();
      } catch (error) {
        errorNotify(error);
      }
    },
    [files, currentEnvironment, userSegment]
  );

  return (
    <SlideModal
      title={t(isUpdate ? 'update-user-segment' : 'new-user-segment')}
      isOpen={isOpen}
      onClose={onClose}
    >
      {isLoadingSegment ? (
        <FormLoading />
      ) : (
        <div className="w-full p-5 pb-28">
          {isUpdate && (
            <p className="text-gray-600 typo-para-medium mb-4">
              {t('form:update-user-segment')}
            </p>
          )}
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
                        disabled={isDisabled}
                        {...field}
                        name="user-segment-name"
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
                        disabled={isDisabled}
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
                disabled={isDisabled}
                className="flex flex-col w-full gap-y-4"
              >
                <Form.Field
                  control={form.control}
                  name="file"
                  render={({ field }) => (
                    <Form.Item
                      className={cn('py-0', {
                        'opacity-50': isDisabled
                      })}
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
                                {
                                  'cursor-not-allowed': isDisabled
                                }
                              )}
                            >
                              {t('form:browse-files')}
                            </label>
                          </div>
                          {userIdsType === 'upload' && !isDisabled && (
                            <div className="flex w-full max-w-full h-fit gap-x-4 pl-8">
                              <Upload
                                files={files}
                                className="border-l border-primary-500 pl-4"
                                uploadClassName="min-h-[200px] h-[200px]"
                                onChange={files => {
                                  setFiles(files);
                                  field.onChange(
                                    files?.length ? files[0] : null
                                  );
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
                      className={cn('py-0', {
                        'opacity-50': isDisabled
                      })}
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
                                {
                                  'cursor-not-allowed': isDisabled
                                }
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
                                placeholder={t(
                                  'form:placeholder-enter-user-ids'
                                )}
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

              {isDisabledUserIds && userSegment && (
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
                    <DisabledButtonTooltip
                      hidden={!isDisabled}
                      trigger={
                        <Button
                          type="submit"
                          disabled={
                            !isDirty || !isValid || isSubmitting || isDisabled
                          }
                          loading={isSubmitting}
                        >
                          {t(`submit`)}
                        </Button>
                      }
                    />
                  }
                />
              </div>
            </Form>
          </FormProvider>
        </div>
      )}
    </SlideModal>
  );
};

export default SegmentCreateUpdateModal;
