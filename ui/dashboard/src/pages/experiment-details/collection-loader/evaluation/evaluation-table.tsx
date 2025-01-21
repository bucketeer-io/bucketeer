import { useTranslation } from 'react-i18next';
import { cn } from 'utils/style';
import { IconInfo } from '@icons';
import Icon from 'components/icon';

const headerList = [
  {
    name: 'variation',
    tooltip: ''
  },
  {
    name: 'evaluation-user',
    tooltip: ''
  },
  {
    name: 'goal-total',
    tooltip: ''
  },
  {
    name: 'goal-user',
    tooltip: ''
  },
  {
    name: 'conversion-rate',
    tooltip: ''
  },
  {
    name: 'value-total',
    tooltip: ''
  },
  {
    name: 'value-user',
    tooltip: ''
  }
];

const HeaderItem = ({
  text,
  isShowIcon = true,
  className
}: {
  text: string;
  isShowIcon?: boolean;
  className?: string;
}) => {
  const formatText = text.replace(' ', '<br />');
  return (
    <div
      className={cn(
        'flex items-center size-fit p-4 pt-0 gap-x-3 text-[13px] leading-[13px] text-gray-500 uppercase',
        { 'min-w-[255px]': !isShowIcon },
        className
      )}
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
          />
        ))}
      </div>
    </div>
  );
};

export default EvaluationTable;
