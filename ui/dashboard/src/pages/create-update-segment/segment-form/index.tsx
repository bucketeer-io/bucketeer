import { useCallback, useMemo, useState } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { IconCloudDownloadOutlined } from 'react-icons-material-design';
import { useNavigate } from 'react-router';
import {
  userSegmentBulkDownload,
  userSegmentBulkUpload,
  userSegmentCreator
} from '@api/user-segment';
import { userSegmentUpdater } from '@api/user-segment/user-segment-updater';
import { yupResolver } from '@hookform/resolvers/yup';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_USER_SEGMENTS } from 'constants/routing';
import { useToast, useToggleOpen } from 'hooks';
import useFormSchema from 'hooks/use-form-schema';
import { useConfirm, useUnsavedLeavePage } from 'hooks/use-unsaved-leave-page';
import { useTranslation } from 'i18n';
import { UserSegment } from '@types';
import { covertFileToUint8ToBase64 } from 'utils/converts';
import { cn } from 'utils/style';
import { IconInfoFilled } from '@icons';
import { formSchema } from 'pages/user-segments/form-schema';
import { UserSegmentForm } from 'pages/user-segments/types';
import SegmentWarning from 'pages/user-segments/user-segment-modal/segment-warning';
import {
  hasSegmentRulesChanged,
  toSegmentFormRules,
  toSegmentRulesPayload
} from 'pages/user-segments/utils';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import { RadioGroup, RadioGroupItem } from 'components/radio';
import TextArea from 'components/textarea';
import Upload from 'components/upload-files';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import ConfirmRulesImpactModal from './confirm-rules-impact-modal';
import SegmentRules from './segment-rules';

interface SegmentFormProps {
  isUpdate: boolean;
  userSegment?: UserSegment;
  isDisabled: boolean;
}

const SegmentForm = ({
  isUpdate,
  userSegment,
  isDisabled
}: SegmentFormProps) => {
  const { t } = useTranslation(['common', 'form', 'message', 'table']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const navigate = useNavigate();
  const { notify, errorNotify } = useToast();
  const { allowNavigation } = useConfirm();

  const isDisabledUserIds = useMemo(
    () =>
      userSegment &&
      (userSegment?.isInUseStatus || userSegment?.features?.length > 0),
    [userSegment]
  );

  const [userIdsType, setUserIdsType] = useState('upload');
  const [files, setFiles] = useState<File[]>([]);
  const [
    isOpenConfirmImpactModal,
    onOpenConfirmImpactModal,
    onCloseConfirmImpactModal
  ] = useToggleOpen(false);

  const schema = useFormSchema(formSchema);
  const resolver = useMemo(() => yupResolver(schema), [schema]);

  const form = useForm({
    resolver,
    values: {
      id: userSegment?.id || '',
      name: userSegment?.name || '',
      description: userSegment?.description || '',
      userIds: '',
      file: null,
      rules: toSegmentFormRules(userSegment?.rules)
    },
    mode: 'onChange'
  });

  const {
    formState: { isValid, isDirty, isSubmitting },
    getFieldState
  } = form;

  const includedUserCount = Number(userSegment?.includedUserCount) || 0;

  const onDownloadUserList = useCallback(async () => {
    if (!userSegment) return;
    try {
      const resp = await userSegmentBulkDownload({
        segmentId: userSegment.id,
        environmentId: currentEnvironment.id
      });
      if (resp.data) {
        const url = window.URL.createObjectURL(
          new Blob([atob(String(resp.data))])
        );
        const link = window.document.createElement('a');
        link.href = url;
        link.setAttribute(
          'download',
          `${currentEnvironment.name}-${userSegment.name}.csv`
        );
        window.document.body.appendChild(link);
        link.click();
        if (link.parentNode) {
          link.parentNode.removeChild(link);
        }
        window.URL.revokeObjectURL(url);
      }
    } catch (error) {
      errorNotify(error);
    }
  }, [userSegment, currentEnvironment]);

  const goToList = useCallback(
    (uploadingSegment?: UserSegment) =>
      navigate(`/${currentEnvironment.urlCode}${PAGE_PATH_USER_SEGMENTS}`, {
        // Let the list page show the uploading indicator immediately,
        // before the server reports the segment status as UPLOADING.
        state: uploadingSegment ? { segmentUploading: uploadingSegment } : null
      }),
    [currentEnvironment]
  );

  const handleSave: SubmitHandler<UserSegmentForm> = useCallback(
    async values => {
      try {
        const { id, name, description } = values;
        const rulesChanged = userSegment
          ? hasSegmentRulesChanged(values.rules, userSegment.rules || [])
          : false;
        let segmentId = id;
        let newSegment = userSegment;
        if (!userSegment) {
          const resp = await userSegmentCreator({
            environmentId: currentEnvironment.id,
            name: values.name,
            description: values.description,
            ...(values.rules?.length
              ? { rules: toSegmentRulesPayload(values.rules, []) }
              : {})
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
          // Await the upload so a failure reaches the catch below instead of
          // being swallowed by the FileReader callback while the UI shows a
          // success toast and navigates away.
          const base64String = await new Promise<string>(resolve =>
            covertFileToUint8ToBase64(file, resolve)
          );
          await userSegmentBulkUpload({
            segmentId: segmentId as string,
            environmentId: currentEnvironment.id,
            state: 'INCLUDED',
            data: base64String
          });
        }

        if (userSegment) {
          if (
            getFieldState('name').isDirty ||
            getFieldState('description').isDirty ||
            rulesChanged
          ) {
            await userSegmentUpdater({
              id: segmentId as string,
              name,
              description,
              environmentId: currentEnvironment.id,
              // The rules field is a full replacement; leave it absent when
              // the user did not modify rules so they stay unchanged.
              ...(rulesChanged
                ? {
                    rules: {
                      values: toSegmentRulesPayload(
                        values.rules,
                        userSegment.rules || []
                      )
                    }
                  }
                : {})
            });
          }
        }
        onCloseConfirmImpactModal();
        notify({
          message: t('message:collection-action-success', {
            collection: t('source-type.segment'),
            action: t(userSegment ? 'updated' : 'created')
          })
        });
        // The form is still flagged dirty at this point; skip the
        // unsaved-changes prompt since everything was just saved.
        allowNavigation(() => goToList(file ? newSegment : undefined));
      } catch (error) {
        errorNotify(error);
      }
    },
    [files, currentEnvironment, userSegment, goToList]
  );

  const onSubmit: SubmitHandler<UserSegmentForm> = useCallback(
    async values => {
      const rulesChanged = userSegment
        ? hasSegmentRulesChanged(values.rules, userSegment.rules || [])
        : false;
      // Changing the rules changes segment membership, which affects the
      // targeting of every flag referencing this segment — confirm first.
      if (rulesChanged && userSegment && userSegment.features?.length > 0) {
        return onOpenConfirmImpactModal();
      }
      return handleSave(values);
    },
    [handleSave, userSegment]
  );

  useUnsavedLeavePage({
    isShow: isDirty && !isSubmitting
  });

  return (
    <FormProvider {...form}>
      <Form
        onSubmit={form.handleSubmit(onSubmit)}
        className="flex flex-col size-full gap-y-6"
      >
        <div className="flex flex-col w-full gap-y-4">
          <div className="flex flex-col w-full p-5 gap-y-5 bg-white rounded-lg shadow-card">
            <p className="typo-para-medium text-gray-700">
              {t('form:general-info')}
            </p>
            <Form.Field
              control={form.control}
              name="name"
              render={({ field }) => (
                <Form.Item className="py-0">
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
                <Form.Item className="py-0">
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
          </div>

          <div className="flex flex-col w-full p-5 gap-y-5 bg-white rounded-lg shadow-card">
            <div>
              <p className="typo-para-medium text-gray-700 mb-2">
                {t('form:list-of-users-ids')}
                <span className="typo-para-small text-gray-500 ml-2">
                  ({t('form:optional')})
                </span>
              </p>
              <p className="typo-para-small text-gray-600">
                {t('form:list-of-users-ids-description')}
              </p>
            </div>

            {isUpdate && userSegment && (
              <div className="flex items-center justify-between gap-x-4 w-full rounded border-l-4 p-4 border-accent-blue-500 bg-accent-blue-50">
                <div className="flex items-center gap-x-2">
                  <Icon
                    icon={IconInfoFilled}
                    size="xxs"
                    color="accent-blue-500"
                  />
                  <p className="typo-para-small leading-[14px] text-accent-blue-500">
                    {t('form:current-user-ids-count', {
                      total: includedUserCount.toLocaleString()
                    })}
                  </p>
                </div>
                {includedUserCount > 0 && (
                  <Button
                    type="button"
                    variant="text"
                    className="w-fit h-auto gap-x-1.5 !p-0 whitespace-nowrap"
                    onClick={onDownloadUserList}
                  >
                    <Icon
                      icon={IconCloudDownloadOutlined}
                      size="sm"
                      className="flex-center"
                    />
                    {t('table:popover.download-segment')}
                  </Button>
                )}
              </div>
            )}

            {isDisabledUserIds && userSegment && (
              <SegmentWarning
                className="!mt-0"
                features={userSegment.features}
              />
            )}

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
                                field.onChange(files?.length ? files[0] : null);
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
          </div>

          {/* The included-user list and the rules are combined with OR. */}
          <div className="flex-center self-center w-[42px] h-[26px] rounded-[3px] typo-para-small leading-[14px] bg-gray-200 text-gray-600">
            {t('common:or')}
          </div>

          <div className="flex flex-col w-full p-5 bg-white rounded-lg shadow-card">
            <p className="typo-para-medium text-gray-700 mb-2">
              {t('table:feature-flags.rules')}
              <span className="typo-para-small text-gray-500 ml-2">
                ({t('form:optional')})
              </span>
            </p>
            <p className="typo-para-small text-gray-600 mb-4">
              {t('form:segment-rules.description')}
            </p>

            <SegmentRules disabled={isDisabled} />
          </div>
        </div>

        <ButtonBar
          className="!border-0 p-0"
          primaryButton={
            <Button
              type="button"
              variant="secondary"
              onClick={() => goToList()}
            >
              {t(`common:cancel`)}
            </Button>
          }
          secondaryButton={
            <DisabledButtonTooltip
              hidden={!isDisabled}
              trigger={
                <Button
                  type="submit"
                  disabled={!isDirty || !isValid || isSubmitting || isDisabled}
                  loading={isSubmitting}
                >
                  {t(`submit`)}
                </Button>
              }
            />
          }
        />
      </Form>
      {isOpenConfirmImpactModal && userSegment && (
        <ConfirmRulesImpactModal
          features={userSegment.features || []}
          isOpen={isOpenConfirmImpactModal}
          loading={isSubmitting}
          onClose={onCloseConfirmImpactModal}
          onSubmit={() => form.handleSubmit(values => handleSave(values))()}
        />
      )}
    </FormProvider>
  );
};

export default SegmentForm;
