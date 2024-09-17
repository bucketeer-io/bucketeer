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
import FlagBoolean from '../../assets/svg/flag-boolean.svg';
import FlagNumber from '../../assets/svg/flag-number.svg';
import FlagString from '../../assets/svg/flag-string.svg';
import FlagJson from '../../assets/svg/flag-json.svg';
import { FlagStatus, getFlagStatus } from '../FeatureList';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { Reason } from '../../proto/feature/reason_pb';
import { HoverPopover } from '../HoverPopover';

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
                          <FlagIcon variationType={variationType} />
                        </div>
                        <div>
                          <div className="flex items-center space-x-3">
                            <Link
                              to={`${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${featureId}${PAGE_PATH_FEATURE_TARGETING}`}
                              className="link text-left text-base font-normal w-full max-w-[60ch] truncate"
                            >
                              {featureDetails.name}
                            </Link>
                            <FlagStatusBadge
                              flagStatus={getFlagStatus(
                                featureDetails,
                                new Date()
                              )}
                            />
                          </div>
                          <p className="text-sm text-gray-400 font-light w-full max-w-[60ch] truncate">
                            {featureId}
                          </p>
                        </div>
                      </div>
                    </td>
                    <td>
                      <div className="whitespace-nowrap px-3 py-4 text-sm text-gray-500 flex space-x-3 items-center">
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

interface FlagIconProps {
  variationType: Feature.AsObject['variationType'];
}
const FlagIcon = ({ variationType }: FlagIconProps) => {
  let icon;
  let msg;

  if (variationType === Feature.VariationType.BOOLEAN) {
    icon = <FlagBoolean width={18} />;
    msg = intl.formatMessage(messages.feature.type.boolean);
  } else if (variationType === Feature.VariationType.STRING) {
    icon = <FlagString width={18} />;
    msg = intl.formatMessage(messages.feature.type.string);
  } else if (variationType === Feature.VariationType.NUMBER) {
    icon = <FlagNumber width={18} />;
    msg = intl.formatMessage(messages.feature.type.number);
  } else if (variationType === Feature.VariationType.JSON) {
    icon = <FlagJson width={18} />;
    msg = intl.formatMessage(messages.feature.type.json);
  } else {
    return null;
  }

  return (
    <HoverPopover
      render={() => {
        return (
          <div className="bg-gray-900 text-white p-2 text-xs rounded whitespace-pre">
            {msg}
          </div>
        );
      }}
    >
      <div className="w-9 h-9 flex justify-center items-center">{icon}</div>
    </HoverPopover>
  );
};
