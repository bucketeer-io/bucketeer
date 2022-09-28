import { QuestionMarkCircleIcon } from '@heroicons/react/solid';
import { FC, memo } from 'react';
import { useIntl } from 'react-intl';
import { Link } from 'react-router-dom';

import { messages } from '../../lang/messages';
import { classNames } from '../../utils/css';
import { HoverPopover } from '../HoverPopover';

export interface HelpTextTooltipProps {
  helpText: string;
  link?: string;
}

export const HelpTextTooltip: FC<HelpTextTooltipProps> = memo(
  ({ helpText, link }) => {
    const { formatMessage: f } = useIntl();
    return (
      <HoverPopover
        render={() => {
          return (
            <div
              className={classNames(
                'border shadow-sm bg-white text-gray-500 p-1',
                'w-64 text-xs rounded whitespace-normal break-words'
              )}
            >
              {helpText}
              {link && (
                <div>
                  <span>
                    {link && (
                      <Link to={link} className="link">
                        {f(messages.readMore)}
                      </Link>
                    )}
                  </span>
                </div>
              )}
            </div>
          );
        }}
      >
        <div className={classNames('hover:text-gray-400')}>
          <QuestionMarkCircleIcon
            className="w-5 h-5 text-gray-300"
            aria-hidden="true"
          />
        </div>
      </HoverPopover>
    );
  }
);
