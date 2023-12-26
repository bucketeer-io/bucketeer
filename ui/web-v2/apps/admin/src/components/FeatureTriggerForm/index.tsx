import { messages } from '@/lang/messages';
import { AppState } from '@/modules';
import {
  createFlagTrigger,
  selectAll,
  deleteFlagTrigger,
  listFlagTriggers,
  updateFlagTrigger,
  resetFlagTrigger,
  disableFlagTrigger,
  enableFlagTrigger,
} from '@/modules/flagTriggers';
import { FlagTrigger } from '@/proto/feature/flag_trigger_pb';
import {
  CreateFlagTriggerResponse,
  ListFlagTriggersResponse,
} from '@/proto/feature/service_pb';
import { AppDispatch } from '@/store';
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
import { FC, memo, useCallback, useEffect, useState } from 'react';
import { useFormContext, Controller } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';

import { ReactComponent as WebhookSvg } from '../../assets/svg/webhook.svg';
import { useCurrentEnvironment, useIsEditable } from '../../modules/me';
import { DetailSkeleton } from '../DetailSkeleton';
import { HoverPopover } from '../HoverPopover';
import { RelativeDateText } from '../RelativeDateText';
import { Select } from '../Select';

const triggerTypeOptions = [
  // {
  //   value: FlagTrigger.Type.TYPE_UNKNOWN.toString(),
  //   label: 'Unknown',
  // },
  {
    value: FlagTrigger.Type.TYPE_WEBHOOK.toString(),
    label: 'Webhook',
  },
];

const actionOptions = [
  {
    value: FlagTrigger.Action.ACTION_OFF.toString(),
    label: 'Action Off',
  },
  {
    value: FlagTrigger.Action.ACTION_ON.toString(),
    label: 'Action On',
  },
  // {
  //   value: FlagTrigger.Action.ACTION_UNKNOWN.toString(),
  //   label: 'Action Unknown',
  // },
];

interface FeatureTriggerFormProps {
  featureId: string;
}

export const FeatureTriggerForm: FC<FeatureTriggerFormProps> = memo(
  ({ featureId }) => {
    const dispatch = useDispatch<AppDispatch>();
    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const {
      control,
      formState: { errors, isDirty },
      watch,
      handleSubmit,
      register,
      reset,
      setValue,
    } = methods;
    const currentEnvironment = useCurrentEnvironment();

    const isLoading = useSelector<AppState, boolean>(
      (state) => state.flagTriggers.loading
    );
    const flagTriggers = useSelector<
      AppState,
      ListFlagTriggersResponse.FlagTriggerWithUrl.AsObject[]
    >((state) => selectAll(state.flagTriggers), shallowEqual);

    const [isAddTriggerOpen, setIsAddTriggerOpen] = useState(false);
    const [selectedFlagTriggerWithUrl, setSelectedFlagTriggerWithUrl] =
      useState<CreateFlagTriggerResponse.AsObject>(null);

    const fetchFlagTriggers = useCallback(() => {
      dispatch(
        listFlagTriggers({
          environmentNamespace: currentEnvironment.id,
          featureId,
        })
      );
    }, []);

    useEffect(() => {
      fetchFlagTriggers();
    }, []);

    const handleDelete = useCallback((flagTriggerId) => {
      dispatch(
        deleteFlagTrigger({
          id: flagTriggerId,
          environmentNamespace: currentEnvironment.id,
        })
      ).then(() => fetchFlagTriggers());
    }, []);

    const handleReset = useCallback((flagTriggerId) => {
      dispatch(
        resetFlagTrigger({
          id: flagTriggerId,
          environmentNamespace: currentEnvironment.id,
        })
      ).then(() => fetchFlagTriggers());
    }, []);

    const handleEnable = useCallback((flagTriggerId) => {
      dispatch(
        enableFlagTrigger({
          id: flagTriggerId,
          environmentNamespace: currentEnvironment.id,
        })
      ).then(() => fetchFlagTriggers());
    }, []);

    const handleDisable = useCallback((flagTriggerId) => {
      dispatch(
        disableFlagTrigger({
          id: flagTriggerId,
          environmentNamespace: currentEnvironment.id,
        })
      ).then(() => fetchFlagTriggers());
    }, []);

    if (isLoading) {
      return (
        <div className="p-9 bg-gray-100">
          <DetailSkeleton />
        </div>
      );
    }

    return (
      <div className="px-10 py-6 bg-white">
        <div className="shadow-md space-y-4 p-5 rounded-sm">
          <p className="text-[#334155]">{f(messages.feature.tab.trigger)}</p>
          <p className="text-sm text-[#728BA3]">
            {f(messages.trigger.description)}
          </p>
          {flagTriggers.map((flagTriggerWithUrl) =>
            flagTriggerWithUrl.flagTrigger.id ===
            selectedFlagTriggerWithUrl?.flagTrigger?.id ? (
              <AddUpdateTrigger
                key={flagTriggerWithUrl.flagTrigger.id}
                close={() => {
                  reset();
                  setSelectedFlagTriggerWithUrl(null);
                }}
                featureId={featureId}
                fetchFlagTriggers={fetchFlagTriggers}
                flagTriggerWithUrl={flagTriggerWithUrl}
              />
            ) : (
              <div
                key={flagTriggerWithUrl.flagTrigger.id}
                className="p-5 border border-[#CBD5E1] rounded-lg flex space-x-3"
              >
                <WebhookSvg className="mt-1" />
                <div className="space-y-3 flex-1">
                  <div className="flex justify-between">
                    <p className="text-[#475569]">
                      {
                        triggerTypeOptions.find(
                          (d) =>
                            d.value ===
                            flagTriggerWithUrl.flagTrigger.type.toString()
                        )?.label
                      }
                    </p>
                    <Popover className="relative flex">
                      <Popover.Button>
                        <div className="flex items-center cursor-pointer text-gray-500">
                          <DotsHorizontalIcon width={20} />
                        </div>
                      </Popover.Button>
                      <Popover.Panel className="absolute z-10 bg-white text-gray-500 right-0 top-7 rounded-lg p-1 whitespace-nowrap shadow-md">
                        <button
                          onClick={() => {
                            setIsAddTriggerOpen(false);
                            reset();
                            setSelectedFlagTriggerWithUrl(flagTriggerWithUrl);
                          }}
                          className="flex w-full space-x-3 px-2 py-1.5 items-center hover:bg-gray-100"
                        >
                          <PencilAltIcon width={20} />
                          <span className="text-sm">
                            {f(messages.trigger.editDescription)}
                          </span>
                        </button>
                        {flagTriggerWithUrl.flagTrigger.disabled ? (
                          <button
                            onClick={() =>
                              handleEnable(flagTriggerWithUrl.flagTrigger.id)
                            }
                            className="flex w-full space-x-3 px-2 py-1.5 items-center hover:bg-gray-100"
                          >
                            <CheckCircleIcon width={20} />
                            <span className="text-sm">
                              {f(messages.trigger.enableTrigger)}
                            </span>
                          </button>
                        ) : (
                          <button
                            onClick={() =>
                              handleDisable(flagTriggerWithUrl.flagTrigger.id)
                            }
                            className="flex w-full space-x-3 px-2 py-1.5 items-center hover:bg-gray-100"
                          >
                            <BanIcon width={20} />
                            <span className="text-sm">
                              {f(messages.trigger.disableTrigger)}
                            </span>
                          </button>
                        )}
                        <button
                          onClick={() =>
                            handleReset(flagTriggerWithUrl.flagTrigger.id)
                          }
                          className="flex w-full space-x-3 px-2 py-1.5 items-center hover:bg-gray-100"
                        >
                          <RefreshIcon width={20} />
                          <span className="text-sm">
                            {f(messages.trigger.resetURL)}
                          </span>
                        </button>
                        <button
                          onClick={() =>
                            handleDelete(flagTriggerWithUrl.flagTrigger.id)
                          }
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
                    {flagTriggerWithUrl.flagTrigger.description}
                  </p>
                  <div className="flex pt-3 border-t border-gray-200 justify-between">
                    <div>
                      <p className="text-gray-400 uppercase text-sm">
                        {f(messages.trigger.flagTarget)}
                      </p>
                      <p className="text-gray-700 mt-1">
                        {FlagTrigger.Action.ACTION_OFF ===
                          flagTriggerWithUrl.flagTrigger.action && 'OFF -> ON'}
                        {FlagTrigger.Action.ACTION_ON ===
                          flagTriggerWithUrl.flagTrigger.action && 'ON -> OFF'}
                        {/* {FlagTrigger.Action.ACTION_UNKNOWN ===
                        flagTriggerWithUrl.flagTrigger.action && 'UNKNOWN'} */}
                      </p>
                    </div>
                    <div>
                      <p className="text-gray-400 uppercase text-sm">
                        {f(messages.trigger.triggerURL)}
                      </p>
                      <a
                        className="text-primary mt-1"
                        href={flagTriggerWithUrl.url}
                        target="_blank"
                        rel="noreferrer"
                      >
                        <span className="underline">
                          {flagTriggerWithUrl.url}
                        </span>
                      </a>
                    </div>
                    <div>
                      <p className="text-gray-400 uppercase text-sm">
                        {f(messages.trigger.triggeredTimes)}
                      </p>
                      <p className="text-gray-700 mt-1">
                        {flagTriggerWithUrl.flagTrigger.triggerCount}
                      </p>
                    </div>
                    <div>
                      <p className="text-gray-400 uppercase text-sm">
                        {f(messages.trigger.lastTriggered)}
                      </p>
                      <p className="text-gray-700 mt-1">
                        {flagTriggerWithUrl.flagTrigger.lastTriggeredAt ? (
                          <RelativeDateText
                            date={
                              new Date(
                                flagTriggerWithUrl.flagTrigger.lastTriggeredAt *
                                  1000
                              )
                            }
                          />
                        ) : (
                          '-'
                        )}
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            )
          )}
          {isAddTriggerOpen && (
            <AddUpdateTrigger
              close={() => {
                reset();
                setIsAddTriggerOpen(false);
              }}
              featureId={featureId}
              fetchFlagTriggers={fetchFlagTriggers}
            />
          )}
          {(!isAddTriggerOpen || selectedFlagTriggerWithUrl) && (
            <button
              onClick={() => {
                reset();
                setSelectedFlagTriggerWithUrl(null);
                setIsAddTriggerOpen(true);
              }}
              className="text-primary flex items-center space-x-2 py-1"
            >
              <PlusIcon width={20} />
              <span>{f(messages.trigger.addTrigger)}</span>
            </button>
          )}
        </div>
      </div>
    );
  }
);

interface AddUpdateTriggerProps {
  close: () => void;
  flagTriggerWithUrl?: CreateFlagTriggerResponse.AsObject;
  featureId: string;
  fetchFlagTriggers: () => void;
}
const AddUpdateTrigger: FC<AddUpdateTriggerProps> = memo(
  ({ close, flagTriggerWithUrl, featureId, fetchFlagTriggers }) => {
    const dispatch = useDispatch<AppDispatch>();
    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const {
      control,
      formState: { errors, isDirty, isValid },
      watch,
      handleSubmit,
      register,
      reset,
      setValue,
    } = methods;
    const editable = useIsEditable();
    const currentEnvironment = useCurrentEnvironment();

    useEffect(() => {
      if (flagTriggerWithUrl) {
        setValue('triggerType', flagTriggerWithUrl.flagTrigger.type.toString());
        setValue('action', flagTriggerWithUrl.flagTrigger.action.toString());
        setValue('description', flagTriggerWithUrl.flagTrigger.description);
      }
    }, [flagTriggerWithUrl]);

    const handleOnSubmit = useCallback((data) => {
      dispatch(
        createFlagTrigger({
          environmentNamespace: currentEnvironment.id,
          featureId,
          triggerType: data.triggerType.value,
          action: data.action.value,
          description: data.description,
        })
      ).then((response) => {
        const payload = response.payload as CreateFlagTriggerResponse.AsObject;
        console.log('url', payload.url);
        fetchFlagTriggers();
        reset();
        close();
      });
    }, []);

    const handleOnSaveSubmit = useCallback(
      (data) => {
        dispatch(
          updateFlagTrigger({
            environmentNamespace: currentEnvironment.id,
            id: flagTriggerWithUrl.flagTrigger.id,
            description: data.description,
          })
        ).then(() => {
          fetchFlagTriggers();
          reset();
          close();
        });
      },
      [flagTriggerWithUrl]
    );

    return (
      <div className="space-y-4 mt-4">
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
                onChange={(o) => field.onChange(o.value)}
                options={triggerTypeOptions}
                value={triggerTypeOptions.find((o) => o.value === field.value)}
                isSearchable={false}
                disabled={!!flagTriggerWithUrl}
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
            name="action"
            control={control}
            render={({ field }) => (
              <Select
                onChange={(o) => field.onChange(o.value)}
                value={actionOptions.find((o) => o.value === field.value)}
                options={actionOptions}
                isSearchable={false}
                disabled={!!flagTriggerWithUrl}
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
            {...register('description')}
            rows={4}
            className="input-textarea w-full"
            disabled={!editable}
          />
          <p className="input-error">
            {errors.description?.message && (
              <span role="alert">{errors.description?.message}</span>
            )}
          </p>
        </div>
        <div className="flex space-x-4">
          <button onClick={close} className="btn-cancel">
            <span>Cancel</span>
          </button>
          {flagTriggerWithUrl ? (
            <button
              onClick={handleSubmit(handleOnSaveSubmit)}
              className="btn-submit"
              disabled={!isValid}
            >
              <span>{f(messages.trigger.save)}</span>
            </button>
          ) : (
            <button
              onClick={handleSubmit(handleOnSubmit)}
              className="btn-submit"
              disabled={!isValid}
            >
              <span>Submit</span>
              {/* <span>{f(messages.trigger.save)}</span> */}
            </button>
          )}
        </div>
      </div>
    );
  }
);
