import { SetStateAction } from 'react';
import { FilterOption } from 'hooks/use-options';
import { useTranslation } from 'i18n';
import { IconPlus } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import { DropdownOption, DropdownValue } from 'components/dropdown';
import Icon from 'components/icon';
import SlideModal from 'components/modal/slide';
import FilterItem from './filter-item';

interface FilterProps {
  isOpen: boolean;
  remainingFilterOptions: DropdownOption[];
  isLoadingTags: boolean;
  isLoading: boolean;
  isDisabledAddButton?: boolean;
  isDisabledSubmitButton?: boolean;
  selectedFilters: FilterOption[];
  onClose: () => void;
  getValueOptions: (filterOption: FilterOption) => FilterOption[];
  handleChangeOption: (value: string, filterIndex: number) => void;
  onConfirmHandler: () => void;
  onClearFilters: () => void;
  handleGetLabelFilterValue: (
    filterOption?: FilterOption
  ) => string | number | true | string[];
  setSelectedFilters: (value: SetStateAction<FilterOption[]>) => void;
  handleChangeFilterValue: (value: DropdownValue, filterIndex: number) => void;
}

const FilterSlide = ({
  isOpen,
  selectedFilters,
  remainingFilterOptions,
  isLoadingTags,
  isLoading,
  isDisabledAddButton,
  isDisabledSubmitButton,
  onClose,
  getValueOptions,
  handleChangeOption,
  handleGetLabelFilterValue,
  setSelectedFilters,
  handleChangeFilterValue,
  onConfirmHandler,
  onClearFilters
}: FilterProps) => {
  const { t } = useTranslation(['common', 'form']);
  return (
    <SlideModal title={t('filters')} isOpen={isOpen} onClose={onClose}>
      <div className="w-full h-full flex flex-col justify-between relative">
        <div className="flex flex-col w-full items-start p-5 gap-y-4">
          {selectedFilters.map((filterOption, filterIndex) => (
            <FilterItem
              key={filterIndex}
              filterIndex={filterIndex}
              filterOption={filterOption}
              filterOptions={remainingFilterOptions}
              isLoadingTags={isLoadingTags}
              isLoading={isLoading}
              selectedFilters={selectedFilters}
              getValueOptions={getValueOptions}
              handleChangeOption={handleChangeOption}
              handleGetLabelFilterValue={handleGetLabelFilterValue}
              setSelectedFilters={setSelectedFilters}
              handleChangeFilterValue={handleChangeFilterValue}
            />
          ))}

          <Button
            disabled={isDisabledAddButton}
            variant={'text'}
            className="h-6 px-0"
            onClick={() => {
              setSelectedFilters([
                ...selectedFilters,
                {
                  label: '',
                  value: undefined,
                  filterValue: ''
                }
              ]);
            }}
          >
            <Icon icon={IconPlus} />
            {t('add-filter')}
          </Button>
        </div>

        <ButtonBar
          className="sticky bottom-0 left-0 bg-white"
          secondaryButton={
            <Button
              disabled={isDisabledSubmitButton}
              onClick={onConfirmHandler}
            >
              {t('confirm')}
            </Button>
          }
          primaryButton={
            <Button onClick={onClearFilters} variant="secondary">
              {t('clear')}
            </Button>
          }
        />
      </div>
    </SlideModal>
  );
};

export default FilterSlide;
