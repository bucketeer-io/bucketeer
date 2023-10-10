import { PAGE_PATH_FEATURES, PAGE_PATH_ROOT } from '@/constants/routing';
import { AppState } from '@/modules';
import { createAutoOpsRule } from '@/modules/autoOpsRules';
import {
  listFeatures,
  selectAll as selectAllFeatures,
} from '@/modules/features';
import { useCurrentEnvironment } from '@/modules/me';
import { addToast } from '@/modules/toasts';
import { OpsType } from '@/proto/autoops/auto_ops_rule_pb';
import { DatetimeClause } from '@/proto/autoops/clause_pb';
import { CreateAutoOpsRuleCommand } from '@/proto/autoops/command_pb';
import { Feature } from '@/proto/feature/feature_pb';
import { ListFeaturesRequest } from '@/proto/feature/service_pb';
import { AppDispatch } from '@/store';
import { Dialog } from '@headlessui/react';
import {
  XCircleIcon,
  ExclamationIcon,
  InformationCircleIcon,
} from '@heroicons/react/solid';
import { FC, useEffect, useState } from 'react';
import ReactDatePicker from 'react-datepicker';
import { Controller, useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import { Link } from 'react-router-dom';

import { FEATURE_UPDATE_COMMENT_MAX_LENGTH } from '../../constants/feature';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { classNames } from '../../utils/css';
import { CheckBox } from '../CheckBox';
import { getFlagStatus, FlagStatus } from '../FeatureList';
import { Modal } from '../Modal';

interface FeatureConfirmDialogProps {
  open: boolean;
  handleSubmit: () => void;
  onClose: () => void;
  title: string;
  description: string;
  displayResetSampling?: boolean;
  isSwitchEnabledConfirm?: boolean;
  isEnabled?: boolean;
  isArchive?: boolean;
  featureId?: string;
  feature?: Feature.AsObject;
}

const SwitchEnabledType = {
  ENABLE_NOW: intl.formatMessage(messages.feature.confirm.enableNow),
  DISABLE_NOW: intl.formatMessage(messages.feature.confirm.disableNow),
  SCHEDULE: intl.formatMessage(messages.feature.confirm.schedule),
};

export const FeatureConfirmDialog: FC<FeatureConfirmDialogProps> = ({
  open,
  handleSubmit,
  onClose,
  title,
  description,
  displayResetSampling,
  isSwitchEnabledConfirm,
  isEnabled,
  isArchive,
  featureId,
  feature,
}) => {
  const dispatch = useDispatch<AppDispatch>();
  const { formatMessage: f } = useIntl();
  const methods = useFormContext();
  const currentEnvironment = useCurrentEnvironment();
  const [flagList, setFlagList] = useState([]);
  const [isFlagActive, setIsFlagActive] = useState(false);
  const [selectedSwitchEnabledType, setSelectedSwitchEnabledType] = useState(
    isEnabled ? SwitchEnabledType.DISABLE_NOW : SwitchEnabledType.ENABLE_NOW
  );
  const [scheduleErrorMessage, setScheduleErrorMessage] = useState('');

  const date = new Date();
  date.setDate(date.getDate() + 1);

  const [datetime, setDatetime] = useState(date);

  const {
    register,
    control,
    formState: { errors, isSubmitting, isDirty, isValid },
  } = methods;

  const features = useSelector<AppState, Feature.AsObject[]>(
    (state) => selectAllFeatures(state.features),
    shallowEqual
  );

  useEffect(() => {
    if (isArchive && open) {
      listFeatures({
        environmentNamespace: currentEnvironment.id,
        pageSize: 0,
        cursor: '',
        tags: [],
        searchKeyword: null,
        maintainerId: null,
        orderBy: ListFeaturesRequest.OrderBy.DEFAULT,
        orderDirection: ListFeaturesRequest.OrderDirection.ASC,
      });
    }
  }, [isArchive, open]);

  useEffect(() => {
    if (isArchive && open && features.length > 0) {
      setFlagList(
        features.reduce((acc, feature) => {
          if (
            feature.prerequisitesList.find(
              (prerequisite) => prerequisite.featureId === featureId
            )
          ) {
            return [
              ...acc,
              {
                id: feature.id,
                name: feature.name,
              },
            ];
          }
          return acc;
        }, [])
      );
    }
  }, [isArchive, open, features, featureId]);

  useEffect(() => {
    if (
      isArchive &&
      open &&
      feature &&
      getFlagStatus(feature, new Date()) === FlagStatus.RECEIVING_REQUESTS
    ) {
      setIsFlagActive(true);
    }
  }, [isArchive, open, feature]);

  const getSubmitBtnLabel = () => {
    if (isSwitchEnabledConfirm) {
      return selectedSwitchEnabledType === SwitchEnabledType.ENABLE_NOW
        ? f(messages.button.enable)
        : selectedSwitchEnabledType === SwitchEnabledType.DISABLE_NOW
        ? f(messages.button.disable)
        : f(messages.button.schedule);
    }
    return f(messages.button.submit);
  };

  const handleScheduleSubmit = () => {
    const command = new CreateAutoOpsRuleCommand();
    command.setFeatureId(featureId);
    if (isEnabled) {
      command.setOpsType(OpsType.DISABLE_FEATURE);
    } else {
      command.setOpsType(OpsType.ENABLE_FEATURE);
    }
    const clause = new DatetimeClause();
    clause.setTime(Math.round(datetime.getTime() / 1000));
    command.setDatetimeClausesList([clause]);
    dispatch(
      createAutoOpsRule({
        environmentNamespace: currentEnvironment.id,
        command: command,
      })
    ).then(() => {
      dispatch(
        addToast({
          message: f(messages.feature.successMessages.schedule),
          severity: 'success',
        })
      );
      onClose();
    });
  };

  const checkSubmitBtnDisabled = () => {
    if (
      selectedSwitchEnabledType === SwitchEnabledType.SCHEDULE &&
      !scheduleErrorMessage
    ) {
      return false;
    }
    return !isDirty || !isValid || isSubmitting;
  };

  return (
    <Modal
      open={open}
      onClose={onClose}
      overflowVisible={isSwitchEnabledConfirm}
    >
      <Dialog.Title
        as="h3"
        className="text-lg font-medium leading-6 text-gray-900"
      >
        {title}
      </Dialog.Title>
      <div className="mt-2">
        <p className="text-sm text-gray-500">{description}</p>
      </div>
      {flagList.length > 0 && (
        <div className="rounded-md bg-red-50 p-4 mt-4">
          <div className="flex">
            <div className="flex-shrink-0">
              <XCircleIcon
                className="h-5 w-5 text-red-400"
                aria-hidden="true"
              />
            </div>
            <div className="ml-3">
              <h3 className="text-sm font-medium text-red-800">
                {f(messages.feature.confirm.flagUsedAsPrerequisite)}
              </h3>
              <div className="mt-2 text-sm text-red-700">
                <ul className="list-disc space-y-1 pl-5">
                  {flagList.map((flag) => (
                    <li key={flag.id}>
                      <Link
                        className="link text-left"
                        to={`${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${flag.id}`}
                      >
                        <p className="truncate w-60">{flag.name}</p>
                      </Link>
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </div>
        </div>
      )}
      {isFlagActive && (
        <div className="bg-yellow-50 p-4 mt-4">
          <div className="flex">
            <div className="flex-shrink-0">
              <ExclamationIcon
                className="h-5 w-5 text-yellow-400"
                aria-hidden="true"
              />
            </div>
            <div className="ml-3">
              <p className="text-sm text-yellow-700">
                {f(messages.feature.confirm.flagIsActive)}
              </p>
            </div>
          </div>
        </div>
      )}
      {selectedSwitchEnabledType !== SwitchEnabledType.SCHEDULE && (
        <div className="mt-5">
          <label
            htmlFor="about"
            className="block text-sm font-medium text-gray-700"
          >
            {f(messages.feature.updateComment)}
          </label>
          <div className="mt-1">
            <textarea
              {...register('comment', {
                required: true,
                maxLength: FEATURE_UPDATE_COMMENT_MAX_LENGTH,
              })}
              id="description"
              rows={3}
              className="input-text w-full"
              disabled={flagList.length > 0}
            />
            <p className="input-error">
              {errors.comment && (
                <span role="alert">{errors.comment.message}</span>
              )}
            </p>
          </div>
          {displayResetSampling && (
            <div className="mt-3 flex space-x-2 items-center">
              <Controller
                name="resetSampling"
                control={control}
                render={({ field }) => {
                  return (
                    <CheckBox
                      id="resample"
                      value={'resample'}
                      onChange={(_: string, checked: boolean): void =>
                        field.onChange(checked)
                      }
                      defaultChecked={false}
                    />
                  );
                }}
              />
              <label htmlFor="resample" className={classNames('input-label')}>
                {f(messages.feature.resetRandomSampling)}
              </label>
            </div>
          )}
        </div>
      )}
      {isSwitchEnabledConfirm && (
        <div className="mt-4 space-y-2">
          <div className="flex items-center space-x-2">
            <input
              id="enable-disable-now"
              type="radio"
              checked={
                selectedSwitchEnabledType === SwitchEnabledType.DISABLE_NOW ||
                selectedSwitchEnabledType === SwitchEnabledType.ENABLE_NOW
              }
              className="h-4 w-4 text-primary focus:ring-primary border-gray-300"
              onChange={() => {
                setSelectedSwitchEnabledType(
                  isEnabled
                    ? SwitchEnabledType.DISABLE_NOW
                    : SwitchEnabledType.ENABLE_NOW
                );
              }}
            />
            <label htmlFor="enable-disable-now">
              {isEnabled
                ? SwitchEnabledType.DISABLE_NOW
                : SwitchEnabledType.ENABLE_NOW}
            </label>
          </div>
          <div>
            <div className="flex items-center space-x-2">
              <input
                id="schedule"
                type="radio"
                checked={
                  selectedSwitchEnabledType === SwitchEnabledType.SCHEDULE
                }
                className="h-4 w-4 text-primary focus:ring-primary border-gray-300"
                onChange={() => {
                  setSelectedSwitchEnabledType(SwitchEnabledType.SCHEDULE);
                }}
              />
              <label htmlFor="schedule">
                {SwitchEnabledType.SCHEDULE}
                <div className="rounded-sm bg-[#F3F9FD] text-[#399CE4] px-2 py-[6px] text-sm ml-3 inline-block">
                  New
                </div>
              </label>
            </div>
            {selectedSwitchEnabledType === SwitchEnabledType.SCHEDULE && (
              <div className="my-3">
                <div className="bg-blue-50 p-4 border-l-4 border-blue-400">
                  <div className="flex">
                    <div className="flex-shrink-0">
                      <InformationCircleIcon
                        className="h-5 w-5 text-blue-400"
                        aria-hidden="true"
                      />
                    </div>
                    <div className="ml-3 flex-1">
                      <p className="text-sm text-blue-700">
                        {f(messages.feature.confirm.scheduleInfo)}
                      </p>
                    </div>
                  </div>
                </div>
                <div className="mt-2">
                  <span className="input-label">
                    {f(messages.autoOps.startDate)}
                  </span>
                  <ReactDatePicker
                    dateFormat="yyyy-MM-dd HH:mm"
                    showTimeSelect
                    timeIntervals={60}
                    placeholderText=""
                    className={classNames('input-text w-full')}
                    wrapperClassName="w-full"
                    selected={datetime}
                    onChange={(d) => {
                      setDatetime(d);
                      if (d.getTime() < new Date().getTime()) {
                        setScheduleErrorMessage(
                          f(messages.input.error.notLaterThanCurrentTime)
                        );
                      } else {
                        setScheduleErrorMessage('');
                      }
                    }}
                  />
                  <p className="input-error">
                    {scheduleErrorMessage && (
                      <span role="alert">{scheduleErrorMessage}</span>
                    )}
                  </p>
                </div>
              </div>
            )}
          </div>
        </div>
      )}
      <div className="pt-5">
        <div className="flex justify-end">
          <button
            type="button"
            className="btn-cancel mr-3"
            disabled={false}
            onClick={onClose}
          >
            {f(messages.button.cancel)}
          </button>
          <button
            type="button"
            className="btn-submit"
            disabled={checkSubmitBtnDisabled()}
            onClick={() => {
              if (
                isSwitchEnabledConfirm &&
                selectedSwitchEnabledType === SwitchEnabledType.SCHEDULE
              ) {
                handleScheduleSubmit();
              } else {
                handleSubmit();
              }
            }}
          >
            {getSubmitBtnLabel()}
          </button>
        </div>
      </div>
    </Modal>
  );
};
