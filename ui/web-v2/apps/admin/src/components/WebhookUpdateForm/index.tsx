import { Dialog } from '@headlessui/react';
import FileCopyOutlined from '@material-ui/icons/FileCopyOutlined';
import { FC, memo, useEffect } from 'react';
import { useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { useDispatch, useSelector } from 'react-redux';
import { useParams } from 'react-router-dom';

import { CopyChip } from '../../components/CopyChip';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import { useIsEditable, useCurrentEnvironment } from '../../modules/me';
import { getWebhook } from '../../modules/webhooks';
import { AppDispatch } from '../../store';
import { DetailSkeleton } from '../DetailSkeleton';
export interface WebhookUpdateFormProps {
  onSubmit: () => void;
  onCancel: () => void;
}

export const WebhookUpdateForm: FC<WebhookUpdateFormProps> = memo(
  ({ onSubmit, onCancel }) => {
    const { webhookId } = useParams<{ webhookId: string }>();
    const dispatch = useDispatch<AppDispatch>();
    const currentEnvironment = useCurrentEnvironment();
    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const editable = useIsEditable();
    const {
      register,
      formState: { errors, isValid, isDirty, isSubmitted },
    } = methods;

    const webhookLoading = useSelector<AppState, boolean>(
      (state) => state.webhook.webhookLoading
    );
    const webhookUrl = useSelector<AppState, string>(
      (state) => state.webhook.webhookUrl
    );

    useEffect(() => {
      dispatch(
        getWebhook({
          environmentNamespace: currentEnvironment.id,
          id: webhookId,
        })
      );
    }, [webhookId]);

    return webhookLoading ? (
      <div className="p-9 bg-gray-100">
        <DetailSkeleton />
      </div>
    ) : (
      <div className="w-[500px]">
        <form className="flex flex-col">
          <div className="flex-1 h-0">
            <div className="py-6 px-4 bg-primary">
              <div className="flex items-center justify-between">
                <Dialog.Title className="text-lg font-medium text-white">
                  {f(messages.webhook.update.header.title)}
                </Dialog.Title>
              </div>
              <div className="mt-1">
                <p className="text-sm text-indigo-300">
                  {f(messages.webhook.update.header.description)}
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
                    <span className="input-label-optional">
                      {' '}
                      {f(messages.input.optional)}
                    </span>
                  </label>
                  <div className="mt-1">
                    <textarea
                      {...register('description')}
                      name="description"
                      id="description"
                      rows={4}
                      className="input-text w-full h-48 break-all"
                    />
                    <p className="input-error">
                      {errors.description && (
                        <span role="alert">{errors.description}</span>
                      )}
                    </p>
                  </div>
                </div>
                <div className="">
                  <label htmlFor="description" className="block">
                    <span className="input-label">Webhook URL</span>
                  </label>
                  <div className="mt-1 relative">
                    <input
                      type="text"
                      name="name"
                      id="name"
                      className="input-text w-full"
                      disabled={true}
                      value={webhookUrl}
                    />
                    <div className="absolute right-[2px] bottom-[1px] top-[1px] h-[38px] pl-3 pr-2 bg-gray-100 cursor-pointer border-l border-gray-300 flex items-center">
                      <CopyChip key={webhookUrl} text={webhookUrl}>
                        <FileCopyOutlined fontSize="small" />
                      </CopyChip>
                    </div>
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
                disabled={!isValid || !isDirty || isSubmitted}
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
