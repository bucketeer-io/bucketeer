import { PAGE_PATH_FEATURES, PAGE_PATH_ROOT } from '@/constants/routing';
import { Dialog } from '@headlessui/react';
import { ExclamationCircleIcon } from '@heroicons/react/solid';
import FileUploadIcon from '@material-ui/icons/CloudUpload';
import FilePresentIcon from '@material-ui/icons/FileCopyOutlined';
import { FC, memo, ChangeEvent, useState, useEffect } from 'react';
import { useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { Link } from 'react-router-dom';

import { messages } from '../../lang/messages';
import { useCurrentEnvironment, useIsEditable } from '../../modules/me';
import { Segment } from '../../proto/feature/segment_pb';
import { classNames } from '../../utils/css';
import { userIdListTypes, UserIdListTypes } from '../SegmentAddForm';

export interface SegmentUpdateFormProps {
  onSubmit: () => void;
  onCancel: () => void;
}

export const SegmentUpdateForm: FC<SegmentUpdateFormProps> = memo(
  ({ onSubmit, onCancel }) => {
    const currentEnvironment = useCurrentEnvironment();

    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const [selectedFile, setSelectedFile] = useState(null);
    const [isDirty, setIsDirty] = useState(false);
    const editable = useIsEditable();
    const {
      register,
      trigger,
      formState,
      getValues,
      formState: { errors, isSubmitted, isValid, dirtyFields },
    } = methods;

    const { isInUseStatus, featureList } = getValues();
    const [selectedUserIdListType, setSelectedUserIdListType] = useState(
      UserIdListTypes.BROWSE
    );

    useEffect(() => {
      setIsDirty(
        dirtyFields.name ||
          dirtyFields.description ||
          selectedFile ||
          dirtyFields.userIds
      );
    }, [formState]);

    const onFileInput = (e: ChangeEvent<HTMLInputElement>): void => {
      if (e.target.files.length > 0) {
        setSelectedFile(e.target.files[0]);
      } else {
        setSelectedFile(null);
      }
      trigger('file');
    };

    return (
      <div className="w-[500px]">
        <form className="flex flex-col">
          <div className="flex-1 h-0">
            <div className="py-6 px-4 bg-primary">
              <div className="flex items-center justify-between">
                <Dialog.Title className="text-lg font-medium text-white">
                  {f(messages.segment.update.header.title)}
                </Dialog.Title>
              </div>
              <div className="mt-1">
                <p className="text-sm text-indigo-300">
                  {f(messages.segment.update.header.description)}
                </p>
              </div>
            </div>
            <div
              className="
                flex-1
                flex flex-col
                justify-between
              "
            >
              <div
                className="
                  space-y-6 px-5 pt-6 pb-5
                  flex flex-col
                "
              >
                <div className="">
                  <label htmlFor="name">
                    <span className="input-label">{f({ id: 'name' })}</span>
                  </label>
                  <div className="mt-1">
                    <input
                      {...register('name')}
                      type="text"
                      name="name"
                      id="name"
                      className="input-text w-full"
                      disabled={!editable || isSubmitted}
                    />
                    <p className="input-error">
                      {errors.name && (
                        <span role="alert">{errors.name.message}</span>
                      )}
                    </p>
                  </div>
                </div>
                <div className="">
                  <label htmlFor="description" className="block">
                    <span className="input-label">
                      {f(messages.description)}
                    </span>
                    &nbsp;
                    <span className="input-label-optional">
                      {f(messages.input.optional)}
                    </span>
                  </label>
                  <div className="mt-1">
                    <textarea
                      {...register('description')}
                      id="description"
                      name="description"
                      rows={4}
                      className="input-text w-full"
                      disabled={!editable || isSubmitted}
                    />
                    <p className="input-error">
                      {errors.description && (
                        <span role="alert">{errors.description.message}</span>
                      )}
                    </p>
                  </div>
                </div>
                {isInUseStatus && (
                  <div className="rounded-md bg-yellow-50 p-4 mt-2">
                    <div className="flex">
                      <div className="flex-shrink-0">
                        <ExclamationCircleIcon
                          className="h-5 w-5 text-yellow-400"
                          aria-hidden="true"
                        />
                      </div>
                      <div className="ml-3 flex-1">
                        <p className="text-sm text-yellow-700">
                          {f(messages.segment.update.userId, {
                            length: featureList.length,
                          })}
                        </p>
                        <div className="mt-2 text-sm text-yellow-700">
                          <ul className="list-disc space-y-1 pl-5">
                            {featureList.map((feature) => (
                              <li key={feature.id}>
                                <Link
                                  className="link text-left"
                                  to={`${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${feature.id}`}
                                >
                                  <p className="truncate w-60">
                                    {feature.name}
                                  </p>
                                </Link>
                              </li>
                            ))}
                          </ul>
                        </div>
                      </div>
                    </div>
                  </div>
                )}
                <div className="">
                  <label htmlFor="file" className="block">
                    <span className="input-label">
                      {f(messages.segment.fileUpload.userList)}
                    </span>
                    &nbsp;
                    <span className="input-label-optional">
                      {f(messages.input.optional)}
                    </span>
                  </label>
                  <div className="mt-1 flex items-center space-x-4">
                    {userIdListTypes.map(({ id, title }) => (
                      <div key={id} className="flex items-center">
                        <input
                          id={id}
                          name="notification-method"
                          type="radio"
                          className="h-4 w-4 border-gray-300 text-primary focus:ring-primary disabled:text-gray-300"
                          defaultChecked={id === selectedUserIdListType}
                          onClick={() => setSelectedUserIdListType(id)}
                          disabled={isInUseStatus}
                        />
                        <label
                          htmlFor={id}
                          className={classNames(
                            'ml-3 block text-sm leading-6',
                            isInUseStatus && 'text-gray-400'
                          )}
                        >
                          {title}
                        </label>
                      </div>
                    ))}
                  </div>
                  <div className="mt-2">
                    {selectedUserIdListType === UserIdListTypes.USER_IDS ? (
                      <div className="mt-1">
                        <textarea
                          {...register('userIds')}
                          id="userIds"
                          name="userIds"
                          rows={4}
                          className="input-text w-full"
                          placeholder={f(
                            messages.segment.enterUserIdsPlaceholder
                          )}
                          disabled={isInUseStatus}
                        />
                      </div>
                    ) : (
                      <div>
                        <div className="mb-2">
                          <div
                            className={classNames(
                              'relative h-[90px] rounded-lg border-dashed',
                              'border-2 border-gray-300 bg-gray-100',
                              'flex justify-center items-center',
                              isInUseStatus && 'opacity-60'
                            )}
                          >
                            <div className="absolute">
                              <div
                                className={classNames(
                                  'flex flex-col items-center'
                                )}
                              >
                                <div className="text-gray-500">
                                  <FileUploadIcon />
                                </div>
                                <span className="block text-gray-500">
                                  {f(messages.segment.fileUpload.browseFiles)}
                                </span>
                              </div>
                            </div>
                            <input
                              {...register('file')}
                              id="file"
                              name="file"
                              type="file"
                              className="input-file"
                              onInput={onFileInput}
                              accept=".csv,.txt"
                              disabled={
                                !editable ||
                                isInUseStatus ||
                                getValues('status') ==
                                  Segment.Status.UPLOADING ||
                                isSubmitted
                              }
                            />
                          </div>
                          <div
                            className={classNames(
                              'flex my-2',
                              isInUseStatus ? 'text-gray-300' : 'text-gray-400'
                            )}
                          >
                            {f(messages.segment.fileUpload.fileFormat)}
                          </div>
                        </div>
                        {selectedFile && (
                          <div
                            className={classNames(
                              'h-14 rounded-lg border-dashed',
                              'border-2 border-gray-300',
                              'flex'
                            )}
                          >
                            <div className="flex items-center ml-3 truncate ...">
                              <div className="text-gray-500">
                                <FilePresentIcon />
                              </div>
                              <div className="ml-3">
                                <p className="text-base text-sm text-gray-700 w-96 truncate ...">
                                  {selectedFile.name}
                                </p>
                                <p className="text-xs text-gray-500">
                                  {f(messages.segment.fileUpload.fileSize, {
                                    fileSize:
                                      selectedFile.size.toLocaleString(),
                                  })}
                                </p>
                              </div>
                            </div>
                          </div>
                        )}
                        <p className="input-error">
                          {errors.file && (
                            <span role="alert">{errors.file.message}</span>
                          )}
                        </p>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div className="flex-shrink-0 px-4 py-4 flex justify-end">
            <div className="mr-3">
              <button
                type="button"
                className="btn-cancel"
                disabled={isSubmitted}
                onClick={onCancel}
              >
                {f(messages.button.cancel)}
              </button>
            </div>
            {editable && (
              <button
                type="button"
                className="btn-submit"
                disabled={!editable || isSubmitted || !isDirty || !isValid}
                onClick={onSubmit}
              >
                {f(messages.button.submit)}
              </button>
            )}
          </div>
        </form>
      </div>
    );
  }
);
