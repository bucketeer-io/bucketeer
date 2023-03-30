import { isLanguageJapanese } from '@/lang/getSelectedLanguage';
import { FC, memo } from 'react';
import { FormattedDate } from 'react-intl';
import { register, format } from 'timeago.js';
import ja from 'timeago.js/lib/lang/ja';

import { HoverPopover } from '../HoverPopover';

export interface RelativeDateTextProps {
  date: Date;
}

if (isLanguageJapanese) {
  register('ja', ja);
}

export const RelativeDateText: FC<RelativeDateTextProps> = memo(({ date }) => {
  return (
    <HoverPopover
      render={() => {
        return (
          <div className="whitespace-nowrap bg-gray-900 text-white p-2 text-xs rounded">
            <FormattedDate
              value={date}
              year="numeric"
              month="long"
              day="numeric"
              weekday="short"
              hour="2-digit"
              minute="2-digit"
            />
          </div>
        );
      }}
    >
      <span>{`${format(date, isLanguageJapanese ? 'ja' : null)}`}</span>
    </HoverPopover>
  );
});
