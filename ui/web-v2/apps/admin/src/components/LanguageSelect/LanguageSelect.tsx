import { getSelectedLanguage, LanguageTypes } from '@/lang/getSelectedLanguage';
import { Listbox, Transition } from '@headlessui/react';
import { CheckIcon, SelectorIcon } from '@heroicons/react/solid';
import React, { FC, Fragment, memo, useState } from 'react';

import englishIcon from '../../assets/english.png';
import japanIcon from '../../assets/japan.png';

import { classNames } from '../../utils/css';

interface ILanguageItem {
  readonly label: string;
  readonly value: LanguageTypes;
  readonly icon: string;
}

const languageList: ILanguageItem[] = [
  { label: '日本国', value: LanguageTypes.JAPAN, icon: japanIcon },
  { label: 'English', value: LanguageTypes.ENGLISH, icon: englishIcon },
];

export const LanguageSelect: FC = memo(() => {
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
      <div className="relative">
        <Listbox.Button
          className={classNames(
            'h-10 py-2 px-2 space-x-3 border flex items-center',
            'text-left bg-white',
            'cursor-default',
            'rounded-md',
            'hover:bg-gray-100',
            'text-gray-700',
            'focus:outline-none focus-visible:ring-2 focus-visible:ring-opacity-75 focus-visible:ring-white focus-visible:ring-offset-orange-300 focus-visible:ring-offset-2 focus-visible:border-indigo-500',
            'text-sm'
          )}
        >
          <img alt={language.label} src={language.icon} />
          <span>{language.label}</span>
          {/* <SelectorIcon className="w-5 h-5 text-gray-700" aria-hidden="true" /> */}
        </Listbox.Button>
        <Transition
          as={Fragment}
          leave="transition ease-in duration-100"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <Listbox.Options
            className={classNames(
              'absolute',
              'w-32 right-0',
              'py-1 mt-1 overflow-auto',
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
                  'cursor-default select-none py-2 flex px- space-x-2 justify-center whitespace-nowrap',
                  'text-sm hover:bg-gray-100',
                  'text-gray-700'
                )}
              >
                {({ selected }) => (
                  <>
                    <img alt={item.label} src={item.icon} />
                    <span>{item.label}</span>
                    {selected ? (
                      <CheckIcon className="w-5 h-5" aria-hidden="true" />
                    ) : (
                      <div className="w-6" />
                    )}
                  </>
                )}
              </Listbox.Option>
            ))}
          </Listbox.Options>
        </Transition>
      </div>
    </Listbox>
  );
});
