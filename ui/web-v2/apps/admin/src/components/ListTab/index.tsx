import React, { FC, Fragment, useState, memo } from 'react';
import { useIntl } from 'react-intl';

import { messages } from '../../lang/messages';
import { Tab } from '../../types/list';
import { classNames } from '../../utils/css';

export const ListTab: FC = memo(() => {
  const { formatMessage: f } = useIntl();
  const [tab, setTab] = useState<Tab>('active');
  const tabs = [
    {
      id: 'active',
      message: f(messages.feature.list.active),
    },
    {
      id: 'archive',
      message: f(messages.feature.list.archive),
    },
  ];

  return (
    <nav className="-mb-px flex" aria-label="Tabs">
      {tabs.map((tabItem, idx) => (
        <button
          key={idx}
          className={classNames(
            tabItem.id == tab ? 'border-primary text-primary' : 'text-gray-500',
            'px-5 py-5',
            'hover:text-gray-700 hover:border-gray-300',
            'whitespace-nowrap py-4 border-b-2',
            'font-medium text-sm'
          )}
        >
          {tabItem.message}
        </button>
      ))}
    </nav>
  );
});
