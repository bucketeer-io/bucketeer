import { AppState } from '../../modules';
import { Tag } from '../../proto/tag/tag_pb';
import { Dialog } from '@headlessui/react';
import { ChangeEvent, FC, memo, useState } from 'react';
import { Controller, useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';
import FileUploadIcon from '@material-ui/icons/CloudUpload';
import FilePresentIcon from '@material-ui/icons/FileCopyOutlined';
import { messages } from '../../lang/messages';
import { selectAll as selectAllTags } from '../../modules/tags';
import { Option } from '../CreatableSelect';
import { classNames } from '../../utils/css';
import { Select } from '../Select';
import { Push } from '../../proto/push/push_pb';
import { selectAll as selectAllPushes } from '../../modules/pushes';

export interface PushAddFormProps {
  onSubmit: () => void;
  onCancel: () => void;
}

export const PushAddForm: FC<PushAddFormProps> = memo(
  ({ onSubmit, onCancel }) => {
    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const {
      register,
      control,
      formState: { errors, isValid, isSubmitted },
      trigger
    } = methods;
    const [selectedFile, setSelectedFile] = useState<File | null>(null);

    const pushList = useSelector<AppState, Push.AsObject[]>(
      (state) => selectAllPushes(state.push),
      shallowEqual
    );

    const tagsList = useSelector<AppState, Tag.AsObject[]>(
      (state) => selectAllTags(state.tags),
      shallowEqual
    );
    const featureFlagTagsList = tagsList.filter(
      (tag) => tag.entityType === Tag.EntityType.FEATURE_FLAG
    );

    // Filter out tags that are already used in existing pushes
    const filteredTagsList = featureFlagTagsList.filter(
      (tag) => !pushList.some((push) => push.tagsList.includes(tag.name))
    );

    const onFileInput = (event: ChangeEvent<HTMLInputElement>): void => {
      const file = event.target.files?.[0];
      if (file) {
        setSelectedFile(file);
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
                  {f(messages.push.add.header.title)}
                </Dialog.Title>
              </div>
              <div className="mt-1">
                <p className="text-sm text-indigo-300">
                  {f(messages.push.add.header.description)}
                </p>
              </div>
            </div>
            <div className="flex-1 flex flex-col justify-between">
              <div className="space-y-6 px-5 pt-6 pb-5 flex flex-col">
                <div className="">
                  <label htmlFor="name">
                    <span className="input-label">{f(messages.name)}</span>
                  </label>
                  <div className="mt-1">
                    <input
                      {...register('name')}
                      type="text"
                      name="name"
                      id="name"
                      className="input-text w-full"
                      disabled={isSubmitted}
                    />
                    <p className="input-error">
                      {errors.name && (
                        <span role="alert">{errors.name.message}</span>
                      )}
                    </p>
                  </div>
                </div>
                <div className="">
                  <label htmlFor="tags">
                    <span className="input-label">
                      {f(messages.tags.title)}
                    </span>
                  </label>
                  <Controller
                    name="tags"
                    control={control}
                    render={({ field }) => {
                      return (
                        <Select
                          isMulti
                          onChange={(options: Option[]) => {
                            field.onChange(options.map((o) => o.value));
                          }}
                          disabled={isSubmitted}
                          options={filteredTagsList.map((tag) => ({
                            label: tag.name,
                            value: tag.name
                          }))}
                          closeMenuOnSelect={false}
                          placeholder={f(messages.tags.tagsPlaceholder)}
                        />
                      );
                    }}
                  />
                  <p className="input-error">
                    {errors.tags && (
                      <span role="alert">{errors.tags.message}</span>
                    )}
                  </p>
                </div>
                <div>
                  <label htmlFor="file">
                    <span className="input-label whitespace-nowrap">
                      {f(messages.push.input.fcmServiceAccount)}
                    </span>
                  </label>
                  <div className="mb-2 mt-1">
                    <div
                      className={classNames(
                        'relative h-[90px] rounded-lg border-dashed',
                        'border-2 border-gray-300 bg-gray-100',
                        'flex justify-center items-center'
                      )}
                    >
                      <div className="absolute">
                        <div className="flex flex-col items-center ">
                          <div className="text-gray-500">
                            <FileUploadIcon />
                          </div>
                          <span className="block text-gray-500">
                            {f(messages.fileUpload.browseFiles)}
                          </span>
                        </div>
                      </div>
                      <input
                        {...register('file')}
                        id="file"
                        name="file"
                        type="file"
                        className="input-file"
                        accept=".json"
                        disabled={isSubmitted}
                        onInput={onFileInput}
                      />
                    </div>
                    <div className="flex text-gray-400 my-2">
                      {f(messages.fileUpload.fileFormatJson)}
                    </div>
                  </div>
                  {selectedFile && (
                    <div
                      className={classNames(
                        'flex h-14 rounded-lg border-dashed border-2',
                        errors.file ? 'border-red-400' : 'border-gray-300'
                      )}
                    >
                      <div className="flex items-center ml-3">
                        <div
                          className={classNames(
                            errors.file ? 'text-red-400' : 'text-gray-400'
                          )}
                        >
                          <FilePresentIcon />
                        </div>
                        <div className="ml-3">
                          <p
                            className={classNames(
                              'text-base w-96 truncate ...',
                              errors.file ? 'text-red-600' : 'text-gray-400'
                            )}
                          >
                            {selectedFile.name}
                          </p>
                          <p
                            className={classNames(
                              'text-xs',
                              errors.file ? 'text-red-500' : 'text-gray-500'
                            )}
                          >
                            {f(messages.fileUpload.fileSize, {
                              fileSize: selectedFile.size.toLocaleString()
                            })}
                          </p>
                        </div>
                      </div>
                    </div>
                  )}
                  <p className="input-error mt-1">
                    {errors.file && (
                      <span role="alert">{errors.file.message}</span>
                    )}
                  </p>
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
            <button
              type="button"
              className="btn-submit"
              disabled={!isValid || isSubmitted}
              onClick={onSubmit}
            >
              {f(messages.button.submit)}
            </button>
          </div>
        </form>
      </div>
    );
  }
);
