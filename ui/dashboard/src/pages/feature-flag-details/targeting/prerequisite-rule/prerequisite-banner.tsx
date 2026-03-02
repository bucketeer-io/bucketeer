import { useState } from 'react';
import { Trans } from 'react-i18next';
import { Link } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import {
  PAGE_PATH_FEATURE_TARGETING,
  PAGE_PATH_FEATURES
} from 'constants/routing';
import { useTranslation } from 'i18n';
import { Feature } from '@types';
import { cn } from 'utils/style';
import { IconChevronDown, IconInfoFilled } from '@icons';
import Icon from 'components/icon';

const PrerequisiteBanner = ({
  hasPrerequisiteFlags
}: {
  hasPrerequisiteFlags: Feature[];
}) => {
  const { t } = useTranslation(['form', 'common', 'table']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const [isExpanded, setIsExpanded] = useState(false);

  return (
    <div className="flex flex-col w-full rounded border-l-4 p-4 border-accent-blue-500 bg-accent-blue-50">
      <button
        type="button"
        onClick={() => setIsExpanded(prev => !prev)}
        className="flex items-center justify-between w-full gap-x-4 cursor-pointer"
      >
        <div className="flex items-center gap-x-2 min-w-0">
          <Icon icon={IconInfoFilled} size="xxs" color="accent-blue-500" />
          <p className="typo-para-small leading-[14px] text-accent-blue-500">
            <Trans
              i18nKey={'form:targeting.prerequisite-flags'}
              values={{
                quantity: hasPrerequisiteFlags.length,
                flag: t(
                  hasPrerequisiteFlags.length > 1
                    ? 'table:flags'
                    : 'common:flag'
                )?.toLowerCase()
              }}
            />
          </p>
        </div>
        <Icon
          icon={IconChevronDown}
          size="xxs"
          color="accent-blue-500"
          className={cn(
            'transition-transform duration-200 flex-shrink-0',
            isExpanded && 'rotate-180'
          )}
        />
      </button>

      {isExpanded && (
        <div className="flex flex-col gap-y-1 pl-6 pt-3">
          <p className="typo-para-small text-gray-700">
            <Trans
              i18nKey={'form:targeting.prerequisite-flags-desc'}
              values={{
                text: t(
                  hasPrerequisiteFlags.length > 1
                    ? 'table:flags'
                    : 'common:flag'
                )?.toLowerCase()
              }}
            />
          </p>
          <ul className="flex flex-col gap-y-1 pt-1">
            {hasPrerequisiteFlags.map((item, index) => (
              <li
                key={index}
                className="typo-para-small text-primary-500 underline w-fit max-w-full truncate"
              >
                <Link
                  to={`/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${item.id}${PAGE_PATH_FEATURE_TARGETING}`}
                >
                  {item.name}
                </Link>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};

export default PrerequisiteBanner;
