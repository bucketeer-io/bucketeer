import { FC, memo, useCallback, useState } from 'react';
import { useIntl } from 'react-intl';

import { messages } from '../../lang/messages';
import { classNames } from '../../utils/css';
import { HoverPopover } from '../HoverPopover';

export interface CopyChipProps {
  text: string;
}

export const CopyChip: FC<CopyChipProps> = memo(({ text, children }) => {
  const { formatMessage: f } = useIntl();
  const [textClicked, setTextClicked] = useState<boolean>(false);

  const handleTextClick = useCallback(
    (text: string) => {
      navigator.clipboard.writeText(text);
      setTextClicked(true);
    },
    [setTextClicked]
  );
  return (
    <HoverPopover
      onClick={() => handleTextClick(text)}
      onMouseLeave={() => setTextClicked(false)}
      render={() => {
        return (
          <div
            className={classNames(
              'bg-gray-900 text-white p-2 text-xs',
              'rounded cursor-pointer whitespace-pre'
            )}
          >
            {textClicked ? (
              <span>{f(messages.copy.copied)}</span>
            ) : (
              <span>{f(messages.copy.copyToClipboard)}</span>
            )}
          </div>
        );
      }}
    >
      {children}
    </HoverPopover>
  );
});
