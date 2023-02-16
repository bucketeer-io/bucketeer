import { getSelectedLanguage, LanguageTypes } from '@/lang/getSelectedLanguage';
import { Listbox, Transition } from '@headlessui/react';
import { Translate, ArrowDropDown } from '@material-ui/icons';
import React, { FC, Fragment, memo, useState } from 'react';

import { classNames } from '../../utils/css';

interface ILanguageItem {
  readonly label: string;
  readonly value: LanguageTypes;
}

const languageList: ILanguageItem[] = [
  { label: '日本語', value: LanguageTypes.JAPANESE },
  { label: 'English', value: LanguageTypes.ENGLISH },
];

export const LanguageSelect: FC = memo(() => {
  const [isOptionsOpen, setIsOptionsOpen] = useState(false);
  const [selectedLanguage, setSelectedLanguage] = useState<LanguageTypes>(
    getSelectedLanguage()
  );

  const handleSwitchLanguage = async (value) => {
    if (value !== selectedLanguage) {
      setSelectedLanguage(value);
      window.localStorage.setItem('language', value);
      window.location.reload();
    }
  };

  const language = languageList.find((lang) => lang.value === selectedLanguage);
  return (
    <Listbox value={selectedLanguage} onChange={handleSwitchLanguage}>
      <div className="relative mt-2">
        <Listbox.Button
          className={classNames(
            'h-10 py-2 pl-2 flex items-center',
            'text-left bg-white',
            'cursor-pointer',
            'rounded-md',
            'hover:bg-gray-100',
            'text-gray-700',
            'focus:outline-none focus-visible:ring-2 focus-visible:ring-opacity-75 focus-visible:ring-white focus-visible:ring-offset-orange-300 focus-visible:ring-offset-2 focus-visible:border-indigo-500',
            'text-sm'
          )}
          onMouseEnter={() => setIsOptionsOpen(true)}
          onMouseLeave={() => setIsOptionsOpen(false)}
        >
          <Translate fontSize="small" />
          <span className="ml-1">{language.label}</span>
          <ArrowDropDown />
        </Listbox.Button>
        <Transition
          as={Fragment}
          leave="transition ease-in duration-100"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
          show={isOptionsOpen}
        >
          <div
            className="h-full"
            onMouseEnter={() => setIsOptionsOpen(true)}
            onMouseLeave={() => setIsOptionsOpen(false)}
          >
            <Listbox.Options
              className={classNames(
                'absolute',
                'w-full right-0 mt-1',
                'py-2 overflow-auto',
                'bg-white text-sm z-50',
                'rounded-md',
                'shadow-lg max-h-60 ring-1 ring-black',
                'ring-opacity-5 focus:outline-none',
                'text-gray-700'
              )}
            >
              {languageList.map((item) => (
                <Listbox.Option
                  key={item.value}
                  value={item.value}
                  className={classNames(
                    'cursor-default select-none px-2',
                    'text-sm',
                    'text-gray-700'
                  )}
                >
                  {({ selected }) => (
                    <p
                      className={`p-2 ${
                        selected
                          ? 'text-primary font-bold'
                          : 'hover:bg-gray-100 cursor-pointer'
                      }`}
                    >
                      {item.label}
                    </p>
                  )}
                </Listbox.Option>
              ))}
            </Listbox.Options>
          </div>
        </Transition>
      </div>
    </Listbox>
  );
});
