import { FC, memo } from 'react';
import { UserIcon } from '@heroicons/react/outline';
import { Link } from 'react-router-dom';
import { UserEvaluation } from '../../pages/debugger';
import {
  PAGE_PATH_FEATURE_TARGETING,
  PAGE_PATH_FEATURES,
  PAGE_PATH_ROOT
} from '../../constants/routing';
import { useCurrentEnvironment } from '../../modules/me';
import { Feature } from '../../proto/feature/feature_pb';
import { classNames } from '../../utils/css';
import UserProfileSvg from '../../assets/svg/user-profile.svg';
import FlagBoolean from '../../assets/svg/flag-boolean.svg';
import FlagNumber from '../../assets/svg/flag-number.svg';
import FlagString from '../../assets/svg/flag-string.svg';
import { FlagStatus, getFlagStatus } from '../FeatureList';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { Reason } from '../../proto/feature/reason_pb';

interface FlagStatusIconProps {
  flagStatus: FlagStatus;
}

export const FlagStatusBadge: FC<FlagStatusIconProps> = ({ flagStatus }) => {
  let flagStatusBadge;
  switch (flagStatus) {
    case FlagStatus.NEW:
      flagStatusBadge = (
        <div className="px-2 py-1.5 rounded bg-blue-50 text-blue-500">
          {intl.formatMessage(messages.feature.flagStatus.new)}
        </div>
      );
      break;
    case FlagStatus.RECEIVING_REQUESTS:
      flagStatusBadge = (
        <div className="px-2 py-1.5 rounded bg-green-50 text-green-500">
          {intl.formatMessage(messages.feature.flagStatus.receivingRequests)}
        </div>
      );
      break;
    case FlagStatus.INACTIVE:
      flagStatusBadge = (
        <div className="px-2 py-1.5 rounded bg-yellow-50 text-yellow-500">
          {intl.formatMessage(messages.feature.flagStatus.inactive)}
        </div>
      );
      break;
  }

  if (flagStatusBadge) {
    return flagStatusBadge;
  }
  return null;
};

export interface DebuggerResultProps {
  userId: string;
  userEvaluations: UserEvaluation[];
  editFields: () => void;
  clearAllFields: () => void;
}

export const DebuggerResult: FC<DebuggerResultProps> = memo(
  ({ userId, userEvaluations, editFields, clearAllFields }) => {
    const currentEnvironment = useCurrentEnvironment();

    return (
      <div>
        <div className="flex justify-between">
          <span className="font-medium">Debugger Results</span>
          <div className="flex space-x-4">
            <button
              className="btn btn-submit !bg-white !border-primary !text-primary !shadow-sm"
              onClick={editFields}
            >
              Edit Fields
            </button>
            <button className="btn btn-submit" onClick={clearAllFields}>
              Clear All Fields
            </button>
          </div>
        </div>
        <div className="mt-12 divide-y divide-gray-100">
          <div className="flex space-x-4 items-center">
            <div className="bg-purple-50 rounded p-2 text-primary">
              <UserIcon width={18} />
            </div>
            <span className="text-primary underline">{userId}</span>
          </div>
          <table className="min-w-full mt-6">
            <thead>
              <tr>
                <th
                  scope="col"
                  className="py-3.5 pl-4 pr-3 text-left text-sm font-normal text-gray-400 sm:pl-0"
                >
                  NAME
                </th>
                <th
                  scope="col"
                  className="px-3 py-3.5 text-left text-sm font-normal text-gray-400"
                >
                  VARIATION
                </th>
                <th
                  scope="col"
                  className="px-3 py-3.5 text-left text-sm font-normal text-gray-400"
                >
                  REASON
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {userEvaluations.map((userEvaluation) => {
                const {
                  variationId,
                  featureDetails,
                  variationName,
                  featureId,
                  reason
                } = userEvaluation;

                const { variationType } = featureDetails;

                const getReasonType = (
                  type: Reason.TypeMap[keyof Reason.TypeMap]
                ) => {
                  return Object.keys(Reason.Type).find(
                    (key) => Reason.Type[key] === type
                  );
                };
                return (
                  <tr key={variationId}>
                    <td>
                      <div className="flex items-center space-x-6 whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-0">
                        <div className="bg-purple-50 rounded flex justify-center items-center text-primary w-9 h-9">
                          {variationType === Feature.VariationType.BOOLEAN ? (
                            <FlagBoolean width={18} />
                          ) : variationType === Feature.VariationType.STRING ? (
                            <FlagString width={18} />
                          ) : variationType === Feature.VariationType.NUMBER ? (
                            <FlagNumber width={18} />
                          ) : null}
                        </div>
                        <div>
                          <div className="flex items-center space-x-3">
                            <Link
                              to={`${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${featureId}${PAGE_PATH_FEATURE_TARGETING}`}
                              className="link text-left text-base font-normal"
                            >
                              {featureDetails.name}
                            </Link>
                            <div className="bg-purple-50 w-8 h-8 flex justify-center items-center rounded">
                              <UserProfileSvg />
                            </div>
                            <FlagStatusBadge
                              flagStatus={getFlagStatus(
                                featureDetails,
                                new Date()
                              )}
                            />
                          </div>
                          <p className="text-sm text-gray-400 font-light">
                            {featureId}
                          </p>
                        </div>
                      </div>
                    </td>
                    <td>
                      <div className="whitespace-nowrap px-3 py-4 text-sm text-gray-500 flex space-x-3 items-center">
                        <div
                          className={classNames(
                            'w-3 h-3 rotate-45 rounded-sm',
                            variationType === Feature.VariationType.BOOLEAN
                              ? 'bg-[#0ea5e9]'
                              : variationType === Feature.VariationType.STRING
                                ? 'bg-pink-600'
                                : variationType === Feature.VariationType.NUMBER
                                  ? 'bg-green-500'
                                  : ''
                          )}
                        />
                        <span>{variationName}</span>
                      </div>
                    </td>
                    <td>
                      <div className="whitespace-nowrap px-3 py-4 text-sm text-gray-500 uppercase">
                        {getReasonType(reason.type)}
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>
    );
  }
);
