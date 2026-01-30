import { useCallback, useEffect, useMemo, useState } from 'react';
import useOptions, { FilterOption, FilterTypes } from 'hooks/use-options';
import { useTranslation } from 'i18n';
import { isEmpty } from 'utils/data-type';
import { cn } from 'utils/style';
import { IconPlus } from '@icons';
import { ExperimentFilters } from 'pages/experiments/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Dropdown, { DropdownOption, DropdownValue } from 'components/dropdown';
import Icon from 'components/icon';
import SlideModal from 'components/modal/slide';

export type FilterProps = {
  isOpen: boolean;
  filters?: Partial<ExperimentFilters>;
  onSubmit: (v: Partial<ExperimentFilters>) => void;
  onClose: () => void;
  onClearFilters: () => void;
};

const FilterExperimentSlideModal = ({
  isOpen,
  filters,
  onSubmit,
  onClose,
  onClearFilters
}: FilterProps) => {
  const { t } = useTranslation(['common']);

  const { experimentFilterOptions } = useOptions();

  const [selectedFilters, setSelectedFilters] = useState<FilterOption[]>([
    experimentFilterOptions[0]
  ]);

  const remainingFilterOptions = useMemo(
    () =>
      experimentFilterOptions.filter(
        option => !selectedFilters.find(item => item.value === option.value)
      ),
    [selectedFilters, experimentFilterOptions]
  );

  const isDisabledAddButton = useMemo(
    () =>
      !remainingFilterOptions.length ||
      selectedFilters.length >= experimentFilterOptions.length,

    [experimentFilterOptions, selectedFilters, remainingFilterOptions]
  );

  const isDisabledSubmitButton = useMemo(
    () => !!selectedFilters.find(item => isEmpty(item.filterValue)),
    [selectedFilters]
  );

  const handleSetFilterOnInit = useCallback(() => {
    if (filters) {
      const { maintainer, statuses } = filters || {};
      const filterTypeArr: FilterOption[] = [];

      const addFilterOption = (
        index: number,
        value: FilterOption['filterValue']
      ) => {
        if (!isEmpty(value)) {
          const option = experimentFilterOptions[index];
          filterTypeArr.push({
            ...option,
            filterValue: value
          });
        }
      };
      addFilterOption(0, statuses);
      addFilterOption(1, maintainer);
      setSelectedFilters(
        filterTypeArr.length ? filterTypeArr : [experimentFilterOptions[0]]
      );
    }
  }, [filters]);

  const handleChangeOption = (value: DropdownValue, filterIndex: number) => {
    {
      const selectedOption = experimentFilterOptions.find(
        item => item.value === value
      );
      if (selectedOption) {
        const filterValue =
          selectedOption.value === FilterTypes.STATUSES ? [] : '';
        selectedFilters[filterIndex] = {
          ...selectedOption,
          filterValue
        };
        setSelectedFilters([...selectedFilters]);
      }
    }
  };

  const onConfirmHandler = () => {
    const defaultFilters = {
      statuses: undefined,
      maintainer: undefined
    };

    const newFilters = {};

    selectedFilters.forEach(filter => {
      Object.assign(newFilters, {
        [filter.value!]: filter.filterValue
      });
    });
    onSubmit({
      ...defaultFilters,
      ...newFilters,
      isFilter: true
    });
  };

  useEffect(() => {
    if (filters?.isFilter || filters?.filterBySummary) {
      handleSetFilterOnInit();
    }
  }, [filters]);

  return (
    <SlideModal title={t('filters')} isOpen={isOpen} onClose={onClose}>
      <div className="w-full h-full flex flex-col justify-between relative">
        <div className="flex flex-col w-full items-start p-5 gap-y-4">
          {selectedFilters.map((filterOption, filterIndex) => {
            const { label, value: filterType } = filterOption;

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
                    options={remainingFilterOptions as DropdownOption[]}
                    value={filterType}
                    onChange={value => {
                      handleChangeOption(value as string, filterIndex);
                    }}
                    className="w-full truncate py-2"
                    contentClassName="w-[270px]"
                  />
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

export default FilterExperimentSlideModal;
