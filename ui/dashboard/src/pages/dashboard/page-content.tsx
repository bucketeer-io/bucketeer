import { useCallback, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import ChartWrapper from 'elements/charts';
import { TempTableDataType } from 'elements/charts/data-collection';
import { ChartDataType } from 'elements/charts/line-chart';
import PageLayout from 'elements/page-layout';
import CollectionLoader from './collection-loader';
import { clientRequestDataPoints, clientRequestDataPoints2 } from './fake-data';
import Overview from './overview';
import { findMinMax, formatNumber, getMonthName } from './utils';

export type TempChartData = {
  x: string;
  y: number;
  type: string;
};

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
  const { t } = useTranslation(['common']);

  const saveXAxisRef = useRef('');

  const [timeValue, setTimeValue] = useState('1d');
  const [dropdownValue, setDropdownValue] = useState('all');
  const [tabValue, setTabValue] = useState('mau');

  const labels = clientRequestDataPoints.map(item => {
    const month = new Date(item.x)?.getMonth();
    const currentValue = getMonthName(month + 1);
    return currentValue;
  });

  const tabs = useMemo(
    () => [
      {
        label: t('mau-count'),
        value: 'mau'
      },
      {
        label: t('request-count'),
        value: 'request'
      }
    ],
    []
  );

  const chartData: ChartDataType = {
    labels,
    datasets: [
      {
        data: clientRequestDataPoints.map(p => p.y),
        borderColor: 'rgba(228, 57, 172, 1)',
        fill: true,
        backgroundColor: 'rgba(228, 57, 172, 0.1)',
        borderWidth: 1,
        pointBackgroundColor: 'transparent',
        pointBorderColor: 'transparent',
        pointHoverRadius: 6,
        pointHoverBackgroundColor: 'rgba(228, 57, 172, 1)'
        // tension: 0.1, --> smooth the line
      },
      {
        data: clientRequestDataPoints2.map(p => p.y),
        borderColor: 'rgba(87, 55, 146, 1)',
        fill: true,
        backgroundColor: 'rgba(87, 55, 146, 0.1)',
        borderWidth: 1,
        pointBackgroundColor: 'transparent',
        pointBorderColor: 'transparent',
        pointHoverRadius: 6,
        pointHoverBackgroundColor: 'rgba(87, 55, 146, 1)'
        // tension: 0.1, --> smooth the line
      }
    ]
  };

  const getTableData = useCallback(
    ({
      data,
      options,
      cb
    }: {
      data: TempChartData[][];
      options?: Intl.NumberFormatOptions;
      cb?: (value: string) => string;
    }): TempTableDataType[] => {
      const _data = data.map(item => {
        const { min, max } = findMinMax(item, 'y');
        const _min = formatNumber(min, options);
        const _max = formatNumber(max, options);
        const _current = formatNumber(item.at(-1)?.y || 0, options);
        return {
          name: item[0].type,
          min: cb ? cb(_min) : _min,
          max: cb ? cb(_max) : _max,
          current: cb ? cb(_current) : _current
        };
      });
      return _data;
    },
    []
  );

  const tableData: TempTableDataType[] = useMemo(
    () =>
      getTableData({
        data: [clientRequestDataPoints, clientRequestDataPoints2],
        cb: (value: string) => {
          return value.length >= 5 ? `${value} K` : value;
        }
      }),
    []
  );

  const formatLabel = useCallback(
    (index: number) => {
      const currentValue = labels[index];
      if (currentValue !== saveXAxisRef.current) {
        saveXAxisRef.current = currentValue;
        return saveXAxisRef.current;
      }
      return '';
    },
    [labels]
  );

  const renderName = useCallback((tempData: TempTableDataType) => {
    return (
      <div className="flex items-center w-full gap-x-2">
        <div
          className={cn('size-3 rounded-sm bg-accent-pink-500', {
            'bg-primary-500': tempData.name === 'RegisterEvents'
          })}
        ></div>
        <p className="text-gray-700 typo-para-medium">{tempData.name}</p>
      </div>
    );
  }, []);

  const onSelectTimeOption = useCallback((value: string) => {
    setTimeValue(value);
  }, []);

  const onSelectDropdownOption = useCallback((value: string) => {
    setDropdownValue(value);
  }, []);

  const onChangeTabs = useCallback((value: string) => {
    setTabValue(value);
  }, []);

  return (
    <PageLayout.Content className="w-full gap-y-5 overflow-auto">
      <Overview />
      <ChartWrapper
        title={`${t('client-request-count')}`}
        timeValue={timeValue}
        dropdownValue={dropdownValue}
        timeOptions={timeOptions}
        dropdownOptions={dropdownOptions}
        tableData={tableData}
        chartData={chartData}
        formatLabel={formatLabel}
        renderName={renderName}
        onSelectTimeOption={onSelectTimeOption}
        onSelectDropdownOption={onSelectDropdownOption}
      />
      <ChartWrapper
        title={`${t('response-time')}`}
        timeValue={timeValue}
        dropdownValue={dropdownValue}
        timeOptions={timeOptions}
        dropdownOptions={dropdownOptions}
        tableData={tableData}
        chartData={chartData}
        formatLabel={formatLabel}
        renderName={renderName}
        onSelectTimeOption={onSelectTimeOption}
        onSelectDropdownOption={onSelectDropdownOption}
      />
      <ChartWrapper
        tabs={tabs}
        timeValue={timeValue}
        tabValue={tabValue}
        dropdownValue={dropdownValue}
        timeOptions={timeOptions}
        dropdownOptions={dropdownOptions}
        tableData={tableData}
        chartData={chartData}
        formatLabel={formatLabel}
        renderName={renderName}
        onSelectTimeOption={onSelectTimeOption}
        onSelectDropdownOption={onSelectDropdownOption}
        onChangeTabs={onChangeTabs}
      />
      <CollectionLoader />
    </PageLayout.Content>
  );
};

export default PageContent;
