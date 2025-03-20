import { useMemo } from 'react';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';

interface Props {
  serveValue: number;
  label?: string;
  isExpand?: boolean;
  onChangeServe: (value: number) => void;
}

const ServeDropdown = ({
  serveValue,
  label,
  isExpand,
  onChangeServe
}: Props) => {
  const { t } = useTranslation(['common', 'table']);

  const serveOptions = useMemo(
    () => [
      { label: t('false'), value: 0 },
      { label: t('true'), value: 1 }
    ],
    []
  );

  const currentServe = useMemo(
    () => serveOptions.find(item => item.value === serveValue),
    [serveValue]
  );

  return (
    <div className={cn('flex items-end gap-x-6', { 'w-full': isExpand })}>
      <p className="typo-para-small text-gray-600 py-[14px] uppercase">
        {t('table:feature-flags.serve')}
      </p>
      <DropdownMenu>
        <div className={cn('flex flex-col gap-y-2', { 'w-full': isExpand })}>
          {label && (
            <p className="typo-para-small leading-[14px] text-gray-600">
              {label}
            </p>
          )}
          <DropdownMenuTrigger
            label="test"
            trigger={
              <div className={cn('flex items-center gap-x-2')}>
                <FlagVariationPolygon index={serveValue} />
                <p className="typo-para-medium leading-5 text-gray-700">
                  {currentServe?.label}
                </p>
              </div>
            }
            className={isExpand ? 'w-full' : ''}
          />
        </div>
        <DropdownMenuContent align="start" isExpand={isExpand}>
          {serveOptions.map((item, index) => (
            <DropdownMenuItem
              key={index}
              label={item.label}
              value={item.value}
              icon={() => <FlagVariationPolygon index={index} />}
              onSelectOption={value => onChangeServe(Number(value))}
            />
          ))}
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
};

export default ServeDropdown;
