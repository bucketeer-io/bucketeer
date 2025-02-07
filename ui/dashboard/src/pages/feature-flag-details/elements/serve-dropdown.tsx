import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';

const serveOptions = [
  { label: 'False', value: 0 },
  { label: 'True', value: 1 }
];

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
  const { t } = useTranslation(['table']);
  const currentServe = serveOptions.find(item => item.value === serveValue);

  return (
    <div className={cn('flex items-end gap-x-6', { 'w-full': isExpand })}>
      <p className="typo-para-small text-gray-600 py-[14px] uppercase">
        {t('feature-flags.serve')}
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
                <FlagVariationPolygon
                  color={currentServe?.value === 0 ? 'pink' : 'blue'}
                />
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
              icon={() => (
                <FlagVariationPolygon color={index === 0 ? 'pink' : 'blue'} />
              )}
              onSelectOption={value => onChangeServe(Number(value))}
            />
          ))}
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
};

export default ServeDropdown;
