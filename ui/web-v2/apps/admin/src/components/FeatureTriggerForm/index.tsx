import { messages } from '@/lang/messages';
import { Popover } from '@headlessui/react';
import {
  DotsHorizontalIcon,
  InformationCircleIcon,
  PlusIcon,
  BanIcon,
  RefreshIcon,
  TrashIcon,
  PencilAltIcon,
  CheckCircleIcon,
} from '@heroicons/react/outline';
import { FC, memo } from 'react';
import { useFormContext, Controller } from 'react-hook-form';
import { useIntl } from 'react-intl';

import { ReactComponent as WebhookSvg } from '../../assets/svg/webhook.svg';
import { useIsEditable } from '../../modules/me';
import { HoverPopover } from '../HoverPopover';
import { Select } from '../Select';

interface FeatureTriggerFormProps {
  // onOpenConfirmDialog: () => void;
}

export const FeatureTriggerForm: FC<FeatureTriggerFormProps> = memo(
  ({ onOpenConfirmDialog }) => {
    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const {
      control,
      formState: { errors, isDirty },
      watch,
    } = methods;
    const editable = useIsEditable();

    return (
      <div className="px-10 py-6 bg-white">
        <div className="shadow-md space-y-4 p-5 rounded-sm">
          <p className="text-[#334155]">{f(messages.feature.tab.trigger)}</p>
          <p className="text-sm text-[#728BA3]">
            {f(messages.trigger.description)}
          </p>
          <div>
            <div className="flex space-x-2 items-center mb-1">
              <label htmlFor="triggerType" className="text-sm text=[#64748B]">
                {f(messages.trigger.triggerType)}
              </label>
              <HoverPopover
                render={() => {
                  return (
                    <div className="shadow p-2 rounded bg-white text-sm whitespace-nowrap -ml-28 mt-[-60px]">
                      Trigger type popover message
                    </div>
                  );
                }}
              >
                <InformationCircleIcon width={18} className="text-gray-400" />
              </HoverPopover>
            </div>
            <Controller
              name="triggerType"
              control={control}
              render={({ field }) => (
                <Select
                  onChange={field.onChange}
                  options={[
                    {
                      value: '',
                      label: '',
                    },
                  ]}
                  disabled={!editable}
                  value={field.value}
                  isSearchable={false}
                />
              )}
            />
          </div>
          <div>
            <div className="flex space-x-2 items-center mb-1">
              <label htmlFor="triggerType" className="text-sm text=[#64748B]">
                {f(messages.trigger.action)}
              </label>
              <HoverPopover
                render={() => {
                  return (
                    <div className="shadow p-2 rounded bg-white text-sm whitespace-nowrap -ml-28 mt-[-60px]">
                      Action popover message
                    </div>
                  );
                }}
              >
                <InformationCircleIcon width={18} className="text-gray-400" />
              </HoverPopover>
            </div>
            <Controller
              name="triggerType"
              control={control}
              render={({ field }) => (
                <Select
                  onChange={field.onChange}
                  options={[
                    {
                      value: '',
                      label: f(messages.trigger.turnTheFlagON),
                    },
                    {
                      value: '',
                      label: f(messages.trigger.turnTheFlagOFF),
                    },
                  ]}
                  disabled={!editable}
                  value={field.value}
                  isSearchable={false}
                />
              )}
            />
          </div>
          <div>
            <div className="flex space-x-1 items-center mb-1">
              <label htmlFor="triggerType" className="text-sm text=[#64748B]">
                {f(messages.description)}
              </label>
              <label htmlFor="optional" className="text-sm text-gray-400">
                {f(messages.input.optional)}
              </label>
            </div>
            <textarea
              id="description"
              rows={4}
              className="input-textarea w-full"
              disabled={!editable}
            />
          </div>
          <div className="p-5 border border-[#CBD5E1] rounded-lg flex space-x-3">
            <WebhookSvg className="mt-1" />
            <div className="space-y-3 flex-1">
              <div className="flex justify-between">
                <p className="text-[#475569]">Generic Webhook</p>
                <Popover className="relative flex">
                  <Popover.Button>
                    <div className="flex items-center cursor-pointer text-gray-500">
                      <DotsHorizontalIcon width={20} />
                    </div>
                  </Popover.Button>
                  <Popover.Panel className="absolute z-10 bg-white text-gray-500 right-0 top-7 rounded-lg p-1 whitespace-nowrap shadow-md">
                    <button
                      // onClick={() => handleOpenUpdate(rule)}
                      className="flex w-full space-x-3 px-2 py-1.5 items-center hover:bg-gray-100"
                    >
                      <PencilAltIcon width={20} />
                      <span className="text-sm">
                        {f(messages.trigger.editDescription)}
                      </span>
                    </button>
                    <button
                      // onClick={() => handleOpenUpdate(rule)}
                      className="flex w-full space-x-3 px-2 py-1.5 items-center hover:bg-gray-100"
                    >
                      <CheckCircleIcon width={20} />
                      <span className="text-sm">
                        {f(messages.trigger.enableTrigger)}
                      </span>
                    </button>
                    <button
                      // onClick={() => handleOpenUpdate(rule)}
                      className="flex w-full space-x-3 px-2 py-1.5 items-center hover:bg-gray-100"
                    >
                      <BanIcon width={20} />
                      <span className="text-sm">
                        {f(messages.trigger.disableTrigger)}
                      </span>
                    </button>
                    <button
                      // onClick={() => handleOpenUpdate(rule)}
                      className="flex w-full space-x-3 px-2 py-1.5 items-center hover:bg-gray-100"
                    >
                      <RefreshIcon width={20} />
                      <span className="text-sm">
                        {f(messages.trigger.resetURL)}
                      </span>
                    </button>
                    <button
                      // onClick={() => handleDelete(rule.id)}
                      className="flex space-x-3 w-full px-2 py-1.5 items-center text-red-500 hover:bg-red-100"
                    >
                      <TrashIcon width={20} />
                      <span className="text-sm">
                        {f(messages.trigger.deleteTrigger)}
                      </span>
                    </button>
                  </Popover.Panel>
                </Popover>
              </div>
              <p className="text-[#728BA3]">
                That's an example of a great description.
              </p>
              <div className="flex pt-3 border-t border-gray-200 justify-between">
                <div>
                  <p className="text-gray-400 uppercase text-sm">
                    {f(messages.trigger.flagTarget)}
                  </p>
                  <p className="text-gray-700 mt-1">{`On -> Off`}</p>
                </div>
                <div>
                  <p className="text-gray-400 uppercase text-sm">
                    {f(messages.trigger.triggerURL)}
                  </p>
                  <a
                    className="text-primary mt-1"
                    href="https://loremIpsum.com/*****fmrio"
                    target="_blank"
                    rel="noreferrer"
                  >
                    <span className="underline">
                      https://loremIpsum.com/*****fmrio
                    </span>
                  </a>
                </div>
                <div>
                  <p className="text-gray-400 uppercase text-sm">
                    {f(messages.trigger.triggeredTimes)}
                  </p>
                  <p className="text-gray-700 mt-1">3</p>
                </div>
                <div>
                  <p className="text-gray-400 uppercase text-sm">
                    {f(messages.trigger.lastTriggered)}
                  </p>
                  <p className="text-gray-700 mt-1">1 hour ago</p>
                </div>
              </div>
            </div>
          </div>
          <button className="text-primary flex items-center space-x-2 py-1">
            <PlusIcon width={20} />
            <span>{f(messages.trigger.addTrigger)}</span>
          </button>
          <button className="flex items-center text-primary px-6 border border-primary rounded-md h-12">
            <span>{f(messages.trigger.save)}</span>
          </button>
        </div>
      </div>
    );
  }
);
