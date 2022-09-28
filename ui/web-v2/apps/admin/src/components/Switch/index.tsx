import { Switch as SwitchReact } from '@headlessui/react';
import { FC } from 'react';

import { classNames } from '../../utils/css';

interface SwitchProps {
  enabled: boolean;
  onChange: () => void;
  size: Size;
  readOnly?: boolean;
}

export type Size = 'small' | 'large';

export const Switch: FC<SwitchProps> = ({
  enabled,
  onChange,
  size,
  readOnly,
}) => {
  return (
    <SwitchReact
      as={readOnly ? 'div' : 'button'}
      checked={enabled}
      onChange={readOnly ? () => {} : onChange}
      className={classNames(
        enabled ? 'bg-primary' : 'bg-gray-200',
        size === 'small' ? 'w-[3.3rem] h-[1.6rem]' : 'w-[4.29rem] h-[2.08rem]',
        `relative inline-flex items-center rounded-full text-gray-700`
      )}
    >
      <span className="sr-only">Switch</span>
      <div
        className={classNames(
          enabled ? 'text-left' : 'text-right',
          size === 'small'
            ? 'text-[0.7rem] p-[0.3rem]'
            : 'text-[0.91rem] p-[0.39rem]',
          'absolute top-0 left-0 w-full h-full',
          'flex flex-col justify-center text-gray-700'
        )}
      >
        {enabled ? (
          <span className="text-white">ON</span>
        ) : (
          <span className="text-gray-700">OFF</span>
        )}
      </div>
      <span
        className={classNames(
          enabled
            ? size === 'small'
              ? 'translate-x-[2.1rem]'
              : 'translate-x-[2.73rem]'
            : size === 'small'
            ? 'translate-x-[0.3rem]'
            : 'translate-x-[0.39rem]',
          'w-[1rem] h-[1rem]',
          size === 'small' ? 'w-[1rem] h-[1rem]' : 'w-[1.3rem] h-[1.3rem]',
          'inline-block transform bg-white rounded-full text-gray-700'
        )}
      />
    </SwitchReact>
  );
};
