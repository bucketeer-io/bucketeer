import { ReactNode, SetStateAction } from 'react';
import { FilterOption } from 'hooks/use-options';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { IconTrash } from '@icons';
import Button from 'components/button';
import Divider from 'components/divider';
import Dropdown, { DropdownOption, DropdownValue } from 'components/dropdown';
import Icon from 'components/icon';

interface FilterProps {
  valueOption: DropdownValue;
  valueLabel?: string | number | true | string[];
  optionFilter: DropdownOption[];
  isLoadingTags: boolean;
  isLoading: boolean;
  isListItem?: boolean;
  isMultiselect?: boolean;
  isSearchable?: boolean;
  selectedFilters: FilterOption[];
  filterIndex: number;
  filterType?: string | number;
  filterOption: FilterOption;
  label: ReactNode;
  disable?: boolean;
  valueOptions: DropdownOption[];
  onClose: () => void;
  getValueOptions: (filterOption: FilterOption) => FilterOption[];
  handleChangeOption: (value: string, filterIndex: number) => void;
  handleGetLabelFilterValue: (
    filterOption?: FilterOption
  ) => string | number | true | string[];
  setSelectedFilters: (value: SetStateAction<FilterOption[]>) => void;
  handleChangeFilterValue: (value: DropdownValue, filterIndex: number) => void;
}

const FilterSlideItem = ({
  filterIndex,
  filterType,
  isListItem,
  isMultiselect,
  label,
  valueOption,
  optionFilter,
  valueOptions,
  disable,
  filterOption,
  isSearchable,
  selectedFilters,
  isLoading,

  handleChangeOption,
  handleGetLabelFilterValue,
  setSelectedFilters,
  handleChangeFilterValue
}: FilterProps) => {
  const { t } = useTranslation(['common', 'form']);
  return (
    <div
      className="flex items-start w-full h-[100px] gap-x-3"
      key={filterIndex}
    >
      <div className="h-full flex flex-col gap-y-4 items-center justify-center">
        <div
          className={cn(
            'mt-2 typo-para-small text-center py-[3px] w-[42px] min-w-[42px] rounded text-accent-pink-500 bg-accent-pink-50',
            {
              'bg-gray-200 text-gray-600': filterIndex !== 0
            }
          )}
        >
          {t(filterIndex === 0 ? `if` : 'and')}
        </div>
        <Divider vertical={true} className="border-primary-500" />
      </div>
      <div className="flex flex-col w-full">
        <Dropdown
          placeholder={t(`select-filter`)}
          labelCustom={label}
          options={optionFilter as DropdownOption[]}
          value={filterType}
          onChange={value => {
            handleChangeOption(value as string, filterIndex);
          }}
          className="w-full truncate py-2"
          contentClassName="w-[270px]"
        />
        <div className="flex items-center gap-3 mt-3 pl-3">
          <p className="typo-para-medium text-gray-600">is</p>
          <Dropdown
            disabled={disable}
            loading={isLoading}
            isListItem={isListItem}
            multiselect={isMultiselect}
            placeholder={t(`select-value`)}
            value={valueOption}
            labelCustom={handleGetLabelFilterValue(filterOption)}
            isSearchable={isSearchable}
            options={valueOptions as DropdownOption[]}
            onChange={value =>
              handleChangeFilterValue(value as DropdownValue, filterIndex)
            }
            className="w-full truncate py-2"
            contentClassName={cn('w-[235px]', {
              'pt-0 w-[300px]': isListItem,
              'hidden-scroll': valueOptions?.length > 15
            })}
          />

          <Button
            variant={'grey'}
            className="px-0 w-fit"
            disabled={selectedFilters.length <= 1}
            onClick={() =>
              setSelectedFilters(
                selectedFilters.filter((_, index) => filterIndex !== index)
              )
            }
          >
            <Icon icon={IconTrash} size={'sm'} />
          </Button>
        </div>
      </div>
    </div>
  );
};

export default FilterSlideItem;
