import { useTranslation } from 'react-i18next';
import { cn } from 'utils/style';
import { IconInfo } from '@icons';
import { Polygon } from 'pages/experiment-details/elements/header-details';
import Icon from 'components/icon';

const headerList = [
  {
    name: 'variation',
    tooltip: '',
    minSize: 255
  },
  {
    name: 'evaluation-user',
    tooltip: '',
    minSize: 143
  },
  {
    name: 'goal-total',
    tooltip: '',
    minSize: 119
  },
  {
    name: 'goal-user',
    tooltip: '',
    minSize: 123
  },
  {
    name: 'conversion-rate',
    tooltip: '',
    minSize: 147
  },
  {
    name: 'value-total',
    tooltip: '',
    minSize: 125
  },
  {
    name: 'value-user',
    tooltip: '',
    minSize: 123
  }
];

const mockData = [
  [
    {
      value: true,
      tooltip: '',
      minSize: 255
    },
    {
      value: 0,
      tooltip: '',
      minSize: 143
    },
    {
      value: 0,
      tooltip: '',
      minSize: 119
    },
    {
      value: 0,
      tooltip: '',
      minSize: 123
    },
    {
      value: 'N/A',
      tooltip: '',
      minSize: 147
    },
    {
      value: 0,
      tooltip: '',
      minSize: 125
    },
    {
      value: 'N/A',
      tooltip: '',
      minSize: 123
    }
  ],
  [
    {
      value: false,
      tooltip: '',
      minSize: 255
    },
    {
      value: 0,
      tooltip: '',
      minSize: 143
    },
    {
      value: 0,
      tooltip: '',
      minSize: 119
    },
    {
      value: 0,
      tooltip: '',
      minSize: 123
    },
    {
      value: 'N/A',
      tooltip: '',
      minSize: 147
    },
    {
      value: 0,
      tooltip: '',
      minSize: 125
    },
    {
      value: 'N/A',
      tooltip: '',
      minSize: 123
    }
  ]
];

const HeaderItem = ({
  text,
  minSize,
  isShowIcon = true,
  className
}: {
  text: string;
  minSize: number;
  isShowIcon?: boolean;
  className?: string;
}) => {
  const formatText = text.replace(' ', '<br />');
  return (
    <div
      className={cn(
        'flex items-center size-fit w-full p-4 pt-0 gap-x-3 text-[13px] leading-[13px] text-gray-500 uppercase',
        className
      )}
      style={{
        minWidth: minSize
      }}
    >
      <p
        dangerouslySetInnerHTML={{
          __html: formatText
        }}
      />
      {isShowIcon && <Icon icon={IconInfo} size={'xxs'} color="gray-500" />}
    </div>
  );
};

const RowItem = ({
  value,
  minSize,
  isFirstItem,
  className
}: {
  value: string | number | boolean;
  minSize: number;
  isFirstItem?: boolean;
  className?: string;
}) => {
  return (
    <div
      className={cn(
        'flex items-center size-fit w-full px-4 py-5 gap-x-2 text-gray-500 capitalize',
        className
      )}
      style={{ minWidth: minSize }}
    >
      {isFirstItem && (
        <Polygon
          className={
            value
              ? 'bg-accent-blue-500 border-none size-3'
              : 'bg-accent-pink-500 border-none size-3'
          }
        />
      )}
      <p className="typo-para-medium leading-4 text-gray-800">
        {String(value)}
      </p>
    </div>
  );
};

const EvaluationTable = () => {
  const { t } = useTranslation(['common', 'table']);
  return (
    <div>
      <div className="flex w-full">
        {headerList.map((item, index) => (
          <HeaderItem
            key={index}
            text={t(`table:metrics.${item.name}`)}
            isShowIcon={index > 0}
            minSize={item.minSize}
          />
        ))}
      </div>
      <div className="divide-y divide-gray-300">
        {...mockData.map((data, index) => (
          <div key={index} className="flex w-full">
            {data.map((item, i) => (
              <RowItem
                isFirstItem={!i}
                key={i}
                value={item.value}
                minSize={item.minSize}
              />
            ))}
          </div>
        ))}
      </div>
    </div>
  );
};

export default EvaluationTable;
