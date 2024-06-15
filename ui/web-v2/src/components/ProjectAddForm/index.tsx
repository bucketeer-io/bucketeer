import { Dialog } from '@headlessui/react';
import { FC, memo } from 'react';
import { useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';

import { messages } from '../../lang/messages';

export interface ProjectAddFormProps {
  onSubmit: () => void;
  onCancel: () => void;
}

export const ProjectAddForm: FC<ProjectAddFormProps> = memo(
  ({ onSubmit, onCancel }) => {
    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const {
      register,
      formState: { errors, isValid, isSubmitted },
    } = methods;

    return (
      <div className="w-[500px]">
        <form className="flex flex-col">
          <div className="flex-1 h-0">
            <div className="py-6 px-4 bg-primary">
              <div className="flex items-center justify-between">
                <Dialog.Title className="text-lg font-medium text-white">
                  {f(messages.adminProject.add.header.title)}
                </Dialog.Title>
              </div>
              <div className="mt-1">
                <p className="text-sm text-indigo-300">
                  {f(messages.adminProject.add.header.description)}
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
                  <label htmlFor="urlCode">
                    <span className="input-label">{f(messages.urlCode)}</span>
                  </label>
                  <div className="mt-1">
                    <input
                      {...register('urlCode')}
                      type="text"
                      name="urlCode"
                      id="urlCode"
                      className="input-text w-full"
                      disabled={isSubmitted}
                    />
                    <p className="input-error">
                      {errors.urlCode && (
                        <span role="alert">{errors.urlCode.message}</span>
                      )}
                    </p>
                  </div>
                </div>
                <div className="">
                  <label htmlFor="description" className="block">
                    <span className="input-label">
                      {f(messages.description)}
                    </span>
                  </label>
                  <div className="mt-1">
                    <textarea
                      {...register('description')}
                      name="description"
                      id="description"
                      className="input-text w-full h-48"
                      disabled={isSubmitted}
                    />
                    <p className="input-error">
                      {errors.description && (
                        <span role="alert">{errors.description.message}</span>
                      )}
                    </p>
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
