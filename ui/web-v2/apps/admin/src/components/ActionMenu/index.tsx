import { Menu, Transition } from '@headlessui/react';
import React, { FC, Fragment, memo } from 'react';

import { classNames } from '../../utils/css';

export enum MenuActions {
  NONE,
  ARCHIVE,
  CLONE,
  CONVERT_PROJECT,
  DELETE,
  DOWNLOAD,
  STOP,
}

export interface MenuItem {
  readonly action: number;
  readonly name: string;
  readonly iconElement: JSX.Element;
  readonly disabled?: boolean;
}

export interface ActionMenuProps {
  onClickAction: (action: number) => void;
  menuItems: Array<MenuItem>;
}

export const ActionMenu: FC<ActionMenuProps> = memo(
  ({ onClickAction, menuItems }) => {
    return (
      <Menu as="div" className="relative text-left">
        <Menu.Button as="div" className="overflow-hidden">
          <button
            className={classNames(
              'disabled:opacity-50 inline-flex justify-center',
              'py-2 px-4 rounded-md',
              'bg-white hover:bg-gray-100',
              'h-10 items-center space-x-1'
            )}
          >
            <div className="w-1 h-1 rounded-full bg-gray-500" />
            <div className="w-1 h-1 rounded-full bg-gray-500" />
            <div className="w-1 h-1 rounded-full bg-gray-500" />
          </button>
        </Menu.Button>
        <Transition
          as={Fragment}
          enter="transition duration-100 ease-out"
          enterFrom="transform scale-95 opacity-0"
          enterTo="transform scale-100 opacity-100"
          leave="transition duration-75 ease-out"
          leaveFrom="transform scale-100 opacity-100"
          leaveTo="transform scale-95 opacity-0"
        >
          <Menu.Items
            className={classNames(
              'absolute right-0 z-10 mt-2 origin-top-right bg-white',
              'divide-gray-100 rounded-md shadow-lg ring-1 ring-black',
              'ring-opacity-5 focus:outline-none'
            )}
          >
            <div className="px-2 py-2">
              {menuItems.map((item) => (
                <Menu.Item key={item.name}>
                  {({ active }) =>
                    item.disabled ? (
                      <button
                        disabled
                        className={classNames(
                          'bg-white text-gray-400 group flex',
                          'rounded-md items-center w-full px-3',
                          'py-2 text-sm disabled:cursor-auto disabled:pointer-events-none'
                        )}
                      >
                        <div className="w-6 h-6 mr-2 text-gray-300">
                          {item.iconElement}
                        </div>
                        <div>{item.name}</div>
                      </button>
                    ) : (
                      <button
                        className={`${
                          active
                            ? 'bg-gray-100 text-gray-700 '
                            : 'bg-white text-gray-700'
                        } group flex rounded-md items-center w-full px-3 py-2 text-sm`}
                        onClick={() => onClickAction(item.action)}
                      >
                        <div className="w-6 h-6 mr-2 text-gray-500">
                          {item.iconElement}
                        </div>
                        <div>{item.name}</div>
                      </button>
                    )
                  }
                </Menu.Item>
              ))}
            </div>
          </Menu.Items>
        </Transition>
      </Menu>
    );
  }
);
