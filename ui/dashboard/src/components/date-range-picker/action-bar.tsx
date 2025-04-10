import { Range } from 'react-date-range';
import { useTranslation } from 'i18n';
import Button from 'components/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import { StaticRangeOption } from '.';

interface Props {
  staticRanges: StaticRangeOption[];
  staticRangeSelected?: StaticRangeOption;
  setRange: (range: Range) => void;
  onApply: () => void;
  onCancel: () => void;
}

const ActionBar = ({
  staticRanges,
  staticRangeSelected,
  setRange,
  onApply,
  onCancel
}: Props) => {
  const { t } = useTranslation(['common', 'form']);

  return (
    <div className="sticky bottom-0 left-0 right-0 flex items-center justify-between w-full p-5 border-t border-gray-200 bg-white">
      <DropdownMenu>
        <DropdownMenuTrigger
          label={staticRangeSelected?.label}
          placeholder={t('form:select-range')}
        />
        <DropdownMenuContent>
          {staticRanges.map((item, index) => (
            <DropdownMenuItem
              key={index}
              label={item.label}
              value={item.label}
              onSelectOption={value => {
                const rangeSelection = staticRanges.find(
                  item => item.label === value
                );
                if (rangeSelection) setRange(rangeSelection.range());
              }}
            />
          ))}
        </DropdownMenuContent>
      </DropdownMenu>
      <div className="flex items-center gap-x-4">
        <Button variant="secondary" onClick={onCancel}>
          {t('cancel')}
        </Button>
        <Button onClick={onApply}>{t('apply')}</Button>
      </div>
    </div>
  );
};

export default ActionBar;
