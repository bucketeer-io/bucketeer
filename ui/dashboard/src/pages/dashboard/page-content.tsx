import { useCallback, useState } from 'react';
import ChartWrapper from 'elements/charts';
import PageLayout from 'elements/page-layout';
import Overview from './overview';

const timeOptions = [
  {
    label: '1D',
    value: '1d'
  },
  {
    label: '7D',
    value: '7d'
  },
  {
    label: '1M',
    value: '1m'
  }
];

const dropdownOptions = [
  {
    label: 'All',
    value: 'all'
  },
  {
    label: 'Android',
    value: 'android'
  },
  {
    label: 'iOS',
    value: 'ios'
  }
];

const PageContent = () => {
  const [timeValue, setTimeValue] = useState('1d');
  const [dropdownValue, setDropdownValue] = useState('all');

  const onSelectTimeOption = useCallback((value: string) => {
    setTimeValue(value);
  }, []);

  const onSelectDropdownOption = useCallback((value: string) => {
    setDropdownValue(value);
  }, []);

  return (
    <PageLayout.Content className="w-full gap-y-5 overflow-auto">
      <Overview />
      <ChartWrapper
        title="Client Request Count"
        timeValue={timeValue}
        dropdownValue={dropdownValue}
        timeOptions={timeOptions}
        dropdownOptions={dropdownOptions}
        onSelectTimeOption={onSelectTimeOption}
        onSelectDropdownOption={onSelectDropdownOption}
      />
    </PageLayout.Content>
  );
};

export default PageContent;
