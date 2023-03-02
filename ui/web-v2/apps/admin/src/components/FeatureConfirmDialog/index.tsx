import { PAGE_PATH_FEATURES, PAGE_PATH_ROOT } from '@/constants/routing';
import { AppState } from '@/modules';
import {
  listFeatures,
  selectAll as selectAllFeatures,
} from '@/modules/features';
import { useCurrentEnvironment } from '@/modules/me';
import { Feature } from '@/proto/feature/feature_pb';
import { ListFeaturesRequest } from '@/proto/feature/service_pb';
import { Dialog } from '@headlessui/react';
import { ExclamationCircleIcon, CheckCircleIcon } from '@heroicons/react/solid';
import { FC, useEffect, useState } from 'react';
import { Controller, useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useSelector } from 'react-redux';
import { Link } from 'react-router-dom';

import {
  FEATURE_LIST_PAGE_SIZE,
  FEATURE_UPDATE_COMMENT_MAX_LENGTH,
} from '../../constants/feature';
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
  isArchive?: boolean;
  featureId?: string;
  feature?: Feature.AsObject;
}

export const FeatureConfirmDialog: FC<FeatureConfirmDialogProps> = ({
  open,
  handleSubmit,
  onClose,
  title,
  description,
  displayResetSampling,
  isArchive,
  featureId,
  feature,
}) => {
  const { formatMessage: f } = useIntl();
  const methods = useFormContext();
  const currentEnvironment = useCurrentEnvironment();
  const [flagList, setFlagList] = useState([]);
  const [isNotUsedFor7Days, setIsNotUsedFor7Days] = useState(false);

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
        environmentNamespace: currentEnvironment.namespace,
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
      console.log('render ======');

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
      getFlagStatus(feature, new Date()) === FlagStatus.INACTIVE
    ) {
      console.log('render ======');
      setIsNotUsedFor7Days(true);
    }
  }, [isArchive, open, feature]);

  return (
    <Modal open={open} onClose={onClose}>
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
        <div className="mt-4">
          <div className="px-4 py-2 bg-red-100 border border-red-600 space-x-3 flex text-sm items-center">
            <ExclamationCircleIcon className="text-red-600 w-8 self-start" />
            <div>
              <span>{f(messages.feature.confirm.flagUsedAsPrerequisite)}</span>
              <ul className="list-disc pl-3 mt-1">
                {flagList.map((flag) => (
                  <li key={flag.id}>
                    <Link
                      className="link text-left"
                      to={`${PAGE_PATH_ROOT}${currentEnvironment.id}${PAGE_PATH_FEATURES}/${flag.id}`}
                    >
                      <p className="truncate w-60">{flag.name}</p>
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </div>
      )}
      {isNotUsedFor7Days && (
        <div className="mt-4">
          <div className="px-4 py-2 bg-green-50 border border-green-500 space-x-3 flex text-sm items-center">
            <CheckCircleIcon className="text-green-500 w-6" />
            <span>{f(messages.feature.confirm.flagNotRequested)}</span>
          </div>
        </div>
      )}
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
            disabled={!isDirty || !isValid || isSubmitting}
            onClick={handleSubmit}
          >
            {f(messages.button.submit)}
          </button>
        </div>
      </div>
    </Modal>
  );
};
