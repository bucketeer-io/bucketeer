import { intl } from '@/lang';
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
  ResetFlagTriggerResponse,
} from '@/proto/feature/service_pb';
import { AppDispatch } from '@/store';
import { classNames } from '@/utils/css';
import { Popover } from '@headlessui/react';
import {
  DotsHorizontalIcon,
  PlusIcon,
  BanIcon,
  RefreshIcon,
  TrashIcon,
  PencilAltIcon,
  CheckCircleIcon,
  XIcon,
  ClockIcon,
  InformationCircleIcon,
} from '@heroicons/react/outline';
import { ExclamationCircleIcon } from '@heroicons/react/solid';
import { FileCopyOutlined } from '@material-ui/icons';
import React, { FC, memo, useCallback, useEffect, useState } from 'react';
import { useFormContext, Controller } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';

import { ReactComponent as OpenInNewSvg } from '../../assets/svg/open-new-tab.svg';
import { ReactComponent as WebhookSvg } from '../../assets/svg/webhook.svg';
import { useCurrentEnvironment, useIsEditable } from '../../modules/me';
import { CopyChip } from '../CopyChip';
import { DetailSkeleton } from '../DetailSkeleton';
import { HoverPopover } from '../HoverPopover';
import { RelativeDateText } from '../RelativeDateText';
import { Select } from '../Select';
import { TriggerDeleteDialog } from '../TriggerDeleteDialog';
import { TriggerResetDialog } from '../TriggerResetDialog';

const triggerTypeOptions = [
  {
    value: FlagTrigger.Type.TYPE_WEBHOOK.toString(),
    label: 'Webhook',
  },
];

const actionOptions = [
  {
    value: FlagTrigger.Action.ACTION_OFF.toString(),
    label: intl.formatMessage(messages.trigger.turnTheFlagOFF),
  },
  {
    value: FlagTrigger.Action.ACTION_ON.toString(),
    label: intl.formatMessage(messages.trigger.turnTheFlagON),
  },
];

interface CopyUrl {
  id: string;
  url: string;
}
interface FeatureTriggerFormProps {
  featureId: string;
}

export const FeatureTriggerForm: FC<FeatureTriggerFormProps> = memo(
  ({ featureId }) => {
    const dispatch = useDispatch<AppDispatch>();
    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const { reset } = methods;
    const currentEnvironment = useCurrentEnvironment();
    const [isDeleteConfirmDialogOpen, setIsDeleteConfirmDialogOpen] =
      useState(false);
    const [isResetConfirmDialogOpen, setIsResetConfirmDialogOpen] =
      useState(false);

    const isLoading = useSelector<AppState, boolean>(
      (state) => state.flagTriggers.loading
    );
    const flagTriggers = useSelector<
      AppState,
      ListFlagTriggersResponse.FlagTriggerWithUrl.AsObject[]
    >((state) => selectAll(state.flagTriggers), shallowEqual);

    const [isAddTriggerOpen, setIsAddTriggerOpen] = useState(false);
    const [selectedFlagTrigger, setSelectedFlagTrigger] =
      useState<CreateFlagTriggerResponse.AsObject>(null);
    const [selectedFlagTriggerForCopyUrl, setSelectedFlagTriggerForCopyUrl] =
      useState<CopyUrl>(null);

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

    const handleDelete = () => {
      setIsDeleteConfirmDialogOpen(false);
      setSelectedFlagTrigger(null);
      dispatch(
        deleteFlagTrigger({
          id: selectedFlagTrigger.flagTrigger.id,
          environmentNamespace: currentEnvironment.id,
        })
      ).then(() => fetchFlagTriggers());
    };

    const handleReset = () => {
      setIsResetConfirmDialogOpen(false);
      setSelectedFlagTrigger(null);
      setSelectedFlagTriggerForCopyUrl(null);
      dispatch(
        resetFlagTrigger({
          id: selectedFlagTrigger.flagTrigger.id,
          environmentNamespace: currentEnvironment.id,
        })
      ).then((response) => {
        const payload = response.payload as ResetFlagTriggerResponse.AsObject;
        setSelectedFlagTriggerForCopyUrl({
          id: payload.flagTrigger.id,
          url: payload.url,
        });
        fetchFlagTriggers();
      });
    };

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
      <>
        <div className="px-10 py-6 bg-white">
          <div className="shadow-md space-y-4 p-5 rounded-sm">
            <p className="text-[#334155]">{f(messages.feature.tab.triggers)}</p>
            <div className="flex">
              <p className="text-sm text-[#728BA3] flex items-center">
                {f(messages.trigger.description, {
                  link: (
                    <a
                      href="https://docs.bucketeer.io"
                      target="_blank"
                      rel="noreferrer"
                      className="underline text-primary flex items-center space-x-1 ml-1"
                    >
                      <span>{f(messages.trigger.documentation)}</span>
                      <OpenInNewSvg />
                    </a>
                  ),
                })}
                {f(messages.fullStop)}
              </p>
            </div>
            {flagTriggers.map((flagTriggerWithUrl) =>
              flagTriggerWithUrl.flagTrigger.id ===
                selectedFlagTrigger?.flagTrigger?.id &&
              !isDeleteConfirmDialogOpen &&
              !isResetConfirmDialogOpen ? (
                <AddUpdateTrigger
                  key={flagTriggerWithUrl.flagTrigger.id}
                  close={() => {
                    reset();
                    setSelectedFlagTrigger(null);
                  }}
                  featureId={featureId}
                  fetchFlagTriggers={fetchFlagTriggers}
                  flagTriggerWithUrl={flagTriggerWithUrl}
                  setSelectedFlagTriggerForCopyUrl={
                    setSelectedFlagTriggerForCopyUrl
                  }
                />
              ) : (
                <div
                  key={flagTriggerWithUrl.flagTrigger.id}
                  className="border border-[#CBD5E1] rounded-lg relative"
                >
                  <div className="absolute left-5 top-5">
                    <WebhookSvg className="w-6 h-6 flex-shrink-0" />
                  </div>
                  <div className="divide-y py-5 pr-5 pl-14">
                    <div className="pb-[10px]">
                      <div className="flex justify-between">
                        <div className="flex space-x-4 items-center">
                          <p className="text-[#475569]">
                            {
                              triggerTypeOptions.find(
                                (d) =>
                                  d.value ===
                                  flagTriggerWithUrl.flagTrigger.type.toString()
                              )?.label
                            }
                          </p>
                          {flagTriggerWithUrl.flagTrigger.disabled ? (
                            <div className="text-sm border rounded border-gray-300 px-[8px] py-[2px] text-gray-500">
                              Off
                            </div>
                          ) : (
                            <div className="text-sm rounded bg-primary px-[8px] py-[2px] text-white">
                              On
                            </div>
                          )}
                        </div>
                        <div className="flex space-x-4 items-center">
                          <div className="flex space-x-1 text-gray-400">
                            <ClockIcon width={17} />
                            {flagTriggerWithUrl.flagTrigger.updatedAt && (
                              <div className="flex text-sm">
                                <p className="">
                                  {f(messages.trigger.updated)}&nbsp;
                                </p>
                                <RelativeDateText
                                  date={
                                    new Date(
                                      flagTriggerWithUrl.flagTrigger.updatedAt *
                                        1000
                                    )
                                  }
                                />
                              </div>
                            )}
                          </div>
                          <Popover className="relative flex">
                            <Popover.Button>
                              <div className="flex items-center cursor-pointer text-gray-500">
                                <DotsHorizontalIcon width={20} />
                              </div>
                            </Popover.Button>
                            <Popover.Panel className="absolute z-10 bg-white text-gray-500 right-1 top-5 rounded-lg p-1 whitespace-nowrap shadow-md">
                              <button
                                onClick={() => {
                                  setIsAddTriggerOpen(false);
                                  reset();
                                  setSelectedFlagTrigger(flagTriggerWithUrl);
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
                                    handleEnable(
                                      flagTriggerWithUrl.flagTrigger.id
                                    )
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
                                    handleDisable(
                                      flagTriggerWithUrl.flagTrigger.id
                                    )
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
                                onClick={() => {
                                  setIsResetConfirmDialogOpen(true);
                                  setSelectedFlagTrigger(flagTriggerWithUrl);
                                }}
                                className="flex w-full space-x-3 px-2 py-1.5 items-center hover:bg-gray-100"
                              >
                                <RefreshIcon width={20} />
                                <span className="text-sm">
                                  {f(messages.trigger.resetURL)}
                                </span>
                              </button>
                              <button
                                onClick={() => {
                                  setIsDeleteConfirmDialogOpen(true);
                                  setSelectedFlagTrigger(flagTriggerWithUrl);
                                }}
                                className="flex space-x-3 w-full px-2 py-[5px] items-center text-red-500 hover:bg-red-100"
                              >
                                <TrashIcon width={20} />
                                <span className="text-sm">
                                  {f(messages.trigger.deleteTrigger)}
                                </span>
                              </button>
                            </Popover.Panel>
                          </Popover>
                        </div>
                      </div>
                      <p className="text-[#728BA3] mt-[10px]">
                        {flagTriggerWithUrl.flagTrigger.description}
                      </p>
                    </div>
                    <div className="pt-3">
                      {selectedFlagTriggerForCopyUrl?.id ===
                      flagTriggerWithUrl.flagTrigger.id ? (
                        <div>
                          <div className="space-y-2">
                            <span className="text-gray-500 text-sm">
                              {f(messages.trigger.triggerURL)}
                            </span>
                            <div className="flex space-x-2">
                              <div className="border border-[#CBD5E1] rounded-lg text-gray-700 px-3 h-[44px] pt-2 truncate ...">
                                {selectedFlagTriggerForCopyUrl.url}
                              </div>
                              <CopyChip
                                key={selectedFlagTriggerForCopyUrl?.url}
                                text={selectedFlagTriggerForCopyUrl?.url}
                              >
                                <div className="flex text-gray-400 border cursor-pointer hover:border-gray-400 border-[#CBD5E1] h-[44px] px-[10px] items-center rounded-lg">
                                  <FileCopyOutlined fontSize="small" />
                                </div>
                              </CopyChip>
                            </div>
                          </div>
                          <div className="flex mt-3 items-center space-x-2">
                            <ExclamationCircleIcon
                              className="h-4 w-4 text-yellow-500"
                              aria-hidden="true"
                            />
                            <div className="flex text-yellow-500">
                              <p>{f(messages.trigger.triggerUrlTitle)}</p>
                              <p>
                                &nbsp;
                                {f(messages.trigger.triggerUrlDescription)}
                              </p>
                            </div>
                          </div>
                        </div>
                      ) : (
                        <div className="flex justify-between">
                          <div>
                            <p className="text-gray-400 uppercase text-sm">
                              {f(messages.trigger.action)}
                            </p>
                            <p className="text-gray-700 mt-1">
                              {FlagTrigger.Action.ACTION_OFF ===
                                flagTriggerWithUrl.flagTrigger.action &&
                                f(messages.trigger.turnTheFlagOFF)}
                              {FlagTrigger.Action.ACTION_ON ===
                                flagTriggerWithUrl.flagTrigger.action &&
                                f(messages.trigger.turnTheFlagON)}
                            </p>
                          </div>
                          <div>
                            <p className="text-gray-400 uppercase text-sm">
                              {f(messages.trigger.triggerURL)}
                            </p>
                            <p className="text-gray-700 mt-1">
                              {flagTriggerWithUrl.url}
                            </p>
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
                              {flagTriggerWithUrl.flagTrigger
                                .lastTriggeredAt ? (
                                <RelativeDateText
                                  date={
                                    new Date(
                                      flagTriggerWithUrl.flagTrigger
                                        .lastTriggeredAt * 1000
                                    )
                                  }
                                />
                              ) : (
                                '-'
                              )}
                            </p>
                          </div>
                        </div>
                      )}
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
                setSelectedFlagTriggerForCopyUrl={
                  setSelectedFlagTriggerForCopyUrl
                }
              />
            )}
            {(!isAddTriggerOpen || selectedFlagTrigger) && (
              <button
                onClick={() => {
                  reset();
                  setSelectedFlagTrigger(null);
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
        <TriggerDeleteDialog
          open={isDeleteConfirmDialogOpen}
          onConfirm={handleDelete}
          onClose={() => {
            setIsDeleteConfirmDialogOpen(false);
            setSelectedFlagTrigger(null);
          }}
        />
        <TriggerResetDialog
          open={isResetConfirmDialogOpen}
          onConfirm={handleReset}
          onClose={() => {
            setIsResetConfirmDialogOpen(false);
            setSelectedFlagTrigger(null);
          }}
        />
      </>
    );
  }
);

interface AddUpdateTriggerProps {
  close: () => void;
  flagTriggerWithUrl?: CreateFlagTriggerResponse.AsObject;
  featureId: string;
  fetchFlagTriggers: () => void;
  setSelectedFlagTriggerForCopyUrl: React.Dispatch<
    React.SetStateAction<CopyUrl>
  >;
}
const AddUpdateTrigger: FC<AddUpdateTriggerProps> = memo(
  ({
    close,
    flagTriggerWithUrl,
    featureId,
    fetchFlagTriggers,
    setSelectedFlagTriggerForCopyUrl,
  }) => {
    const dispatch = useDispatch<AppDispatch>();
    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const {
      control,
      formState: { errors, isValid },
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
          triggerType: data.triggerType,
          action: data.action,
          description: data.description,
        })
      ).then((response) => {
        const payload = response.payload as CreateFlagTriggerResponse.AsObject;
        setSelectedFlagTriggerForCopyUrl({
          id: payload.flagTrigger.id,
          url: payload.url,
        });
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
          setSelectedFlagTriggerForCopyUrl(null);
        });
      },
      [flagTriggerWithUrl]
    );

    return (
      <div className="space-y-4 mt-6">
        <div>
          <div className="flex space-x-2 items-center mb-1">
            <label htmlFor="triggerType" className="text-sm text=[#64748B]">
              {f(messages.trigger.triggerType)}
            </label>
            <HoverPopover
              render={() => {
                return (
                  <div
                    className={classNames(
                      'bg-gray-900 text-white p-2 text-xs',
                      'rounded cursor-pointer whitespace-pre'
                    )}
                  >
                    <span>Tooltip</span>
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
                options={triggerTypeOptions}
                onChange={(o) => field.onChange(o.value)}
                value={triggerTypeOptions.find((o) => o.value === field.value)}
                disabled={!!flagTriggerWithUrl}
                isSearchable={false}
                formatOptionLabel={({ label, value }) => {
                  return (
                    <div className="flex space-x-4 items-center">
                      {value === FlagTrigger.Type.TYPE_WEBHOOK.toString() && (
                        <WebhookSvg className="w-6 h-6" />
                      )}
                      <span className="flex-1 truncate">{label}</span>
                    </div>
                  );
                }}
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
                  <div
                    className={classNames(
                      'bg-gray-900 text-white p-2 text-xs',
                      'rounded cursor-pointer whitespace-pre'
                    )}
                  >
                    <span>Tooltip</span>
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
        </div>
        <div className="flex space-x-4">
          <button onClick={close} className="btn-cancel">
            <span>{f(messages.button.cancel)}</span>
          </button>
          {flagTriggerWithUrl ? (
            <button
              onClick={handleSubmit(handleOnSaveSubmit)}
              className="btn-submit"
              disabled={!isValid}
            >
              <span>{f(messages.button.save)}</span>
            </button>
          ) : (
            <button
              onClick={handleSubmit(handleOnSubmit)}
              className="btn-submit"
              disabled={!isValid}
            >
              <span>{f(messages.button.submit)}</span>
            </button>
          )}
        </div>
      </div>
    );
  }
);
