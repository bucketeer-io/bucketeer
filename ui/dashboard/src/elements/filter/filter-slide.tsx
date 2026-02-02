import { SetStateAction } from 'react';
import { FilterOption, FilterTypes } from 'hooks/use-options';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import { IconPlus, IconTrash } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Dropdown, { DropdownOption, DropdownValue } from 'components/dropdown';
import Icon from 'components/icon';
import SlideModal from 'components/modal/slide';

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
          {selectedFilters.map((filterOption, filterIndex) => {
            const { label, value: filterType } = filterOption;
            const isMaintainerFilter = filterType === FilterTypes.MAINTAINER;
            const isTagFilter = filterType === FilterTypes.TAGS;
            const isHaveSearchingDropdown = isMaintainerFilter || isTagFilter;
            const valueOptions = getValueOptions(filterOption);
            return (
              // <FilterSlideItem
              //  key={filterIndex}
              //   filterIndex={filterIndex}
              //   label={label}
              //   filterType={filterType}
              //   valueOption={filterOption?.filterValue as DropdownValue}
              //   isLoading={
              //     (isTagFilter && isLoadingTags) ||
              //     (isMaintainerFilter && isLoading)
              //   }
              //   disable={
              //     (isTagFilter && isLoadingTags) ||
              //     (isMaintainerFilter && isLoading) ||
              //     !filterType
              //   }
              //   isLoadingTags={isLoadingTags}
              //   valueOptions={valueOptions as DropdownOption[]}
              //   isListItem={isHaveSearchingDropdown}
              //   isMultiselect={isTagFilter}
              //   filterOption={filterOption}
              //   selectedFilters={selectedFilters}
              //   optionFilter={remainingFilterOptions}
              //   getValueOptions={getValueOptions}
              //   handleChangeOption={handleChangeOption}
              //   handleGetLabelFilterValue={handleGetLabelFilterValue}
              //   setSelectedFilters={setSelectedFilters}
              //   handleChangeFilterValue={handleChangeFilterValue}
              //   onClose={onClose}

              // />
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
                    options={remainingFilterOptions as DropdownOption[]}
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
                      disabled={
                        (isTagFilter && isLoadingTags) ||
                        (isMaintainerFilter && isLoading) ||
                        !filterType
                      }
                      loading={
                        (isTagFilter && isLoadingTags) ||
                        (isMaintainerFilter && isLoading)
                      }
                      isListItem={isHaveSearchingDropdown}
                      multiselect={isTagFilter}
                      placeholder={t(`select-value`)}
                      value={filterOption?.filterValue as DropdownValue}
                      labelCustom={handleGetLabelFilterValue(filterOption)}
                      isSearchable={isHaveSearchingDropdown}
                      options={valueOptions as DropdownOption[]}
                      onChange={value =>
                        handleChangeFilterValue(
                          value as DropdownValue,
                          filterIndex
                        )
                      }
                      className="w-full truncate py-2"
                      contentClassName={cn('w-[235px]', {
                        'pt-0 w-[300px]': isHaveSearchingDropdown,
                        'hidden-scroll': valueOptions?.length > 15
                      })}
                    />

                    <Button
                      variant={'grey'}
                      className="px-0 w-fit"
                      disabled={selectedFilters.length <= 1}
                      onClick={() =>
                        setSelectedFilters(
                          selectedFilters.filter(
                            (_, index) => filterIndex !== index
                          )
                        )
                      }
                    >
                      <Icon icon={IconTrash} size={'sm'} />
                    </Button>
                  </div>
                </div>
              </div>
            );
          })}

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
              {t(`confirm`)}
            </Button>
          }
          primaryButton={
            <Button onClick={onClearFilters} variant="secondary">
              {t(`clear`)}
            </Button>
          }
        />
      </div>
    </SlideModal>
  );
};

export default FilterSlide;
