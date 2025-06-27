import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useQueryAccounts } from '@queries/accounts';
import { getCurrentEnvironment, useAuth } from 'auth';
import useOptions, { FilterOption, FilterTypes } from 'hooks/use-options';
import { useTranslation } from 'i18n';
import { debounce } from 'lodash';
import { isEmpty } from 'utils/data-type';
import { cn } from 'utils/style';
import { IconPlus, IconTrash } from '@icons';
import { ExperimentFilters } from 'pages/experiments/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSearch,
  DropdownMenuTrigger
} from 'components/dropdown';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';

export type FilterProps = {
  onSubmit: (v: Partial<ExperimentFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<ExperimentFilters>;
};

const FilterExperimentModal = ({
  onSubmit,
  isOpen,
  onClose,
  onClearFilters,
  filters
}: FilterProps) => {
  const { t } = useTranslation(['common']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { experimentStatusOptions, experimentFilterOptions } = useOptions();

  const inputSearchRef = useRef<HTMLInputElement>(null);

  const [selectedFilters, setSelectedFilters] = useState<FilterOption[]>([
    experimentFilterOptions[0]
  ]);
  const [searchValue, setSearchValue] = useState('');
  const [debounceValue, setDebounceValue] = useState('');

  const { data: collection, isLoading } = useQueryAccounts({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment?.id,
      organizationId: currentEnvironment?.organizationId
    },
    enabled: !!selectedFilters?.find(
      item => item.value === FilterTypes.MAINTAINER
    )
  });

  const accounts = collection?.accounts || [];

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

  const debouncedSearch = useCallback(
    debounce(value => {
      setSearchValue(value);
    }, 500),
    []
  );

  const getValueOptions = useCallback(
    (filterOption: FilterOption) => {
      if (!filterOption.value) return [];
      const isMaintainerFilter = filterOption.value === FilterTypes.MAINTAINER;
      if (isMaintainerFilter) {
        const options = accounts.map(item => ({
          label: item.email,
          value: item.email
        }));

        return options?.filter(item =>
          searchValue
            ? item.value.toLowerCase().includes(searchValue.toLowerCase())
            : item
        );
      }
      return experimentStatusOptions;
    },
    [accounts, searchValue, experimentStatusOptions]
  );

  const handleFocusSearchInput = useCallback(() => {
    let timerId: NodeJS.Timeout | null = null;
    if (timerId) clearTimeout(timerId);
    timerId = setTimeout(() => inputSearchRef?.current?.focus(), 50);
  }, []);

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

  const handleGetLabelFilterValue = useCallback(
    (filterOption?: FilterOption) => {
      if (filterOption) {
        const { value: filterType, filterValue } = filterOption;
        const isMaintainerFilter = filterType === FilterTypes.MAINTAINER;
        return isMaintainerFilter
          ? filterValue || ''
          : Array.isArray(filterValue)
            ? filterValue
                .map(
                  value =>
                    experimentStatusOptions.find(item => item.value === value)
                      ?.label
                )
                ?.join(', ')
                ?.toLowerCase()
            : filterValue;
      }
      return '';
    },
    [experimentStatusOptions]
  );

  const handleChangeFilterValue = useCallback(
    (value: string | number, filterIndex: number) => {
      const filterOption = selectedFilters[filterIndex];
      const { value: filterType, filterValue } = filterOption;
      const isStatusOption = filterType === FilterTypes.STATUSES;
      if (isStatusOption) {
        const values = filterValue as string[];
        const isExisted = values.find(item => item === value);
        const newValue: string[] = isExisted
          ? values.filter(item => item !== value)
          : [...values, value as string];
        selectedFilters[filterIndex] = {
          ...selectedFilters[filterIndex],
          filterValue: newValue
        };
        return setSelectedFilters([...selectedFilters]);
      }
      selectedFilters[filterIndex] = {
        ...selectedFilters[filterIndex],
        filterValue: value
      };
      setSelectedFilters([...selectedFilters]);
    },
    [selectedFilters]
  );

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
    <DialogModal
      className="w-[665px]"
      title={t('filters')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full items-start p-5 gap-y-4">
        {selectedFilters.map((filterOption, filterIndex) => {
          const { label, value: filterType } = filterOption;
          const isMaintainerFilter = filterType === FilterTypes.MAINTAINER;
          const isStatusFilter = filterType === FilterTypes.STATUSES;
          const valueOptions = getValueOptions(filterOption);

          return (
            <div
              key={filterIndex}
              className="flex items-center w-full h-12 gap-x-4"
            >
              <div
                className={cn(
                  'typo-para-small text-center py-[3px] w-[42px] min-w-[42px] rounded text-accent-pink-500 bg-accent-pink-50',
                  {
                    'bg-gray-200 text-gray-600': filterIndex !== 0
                  }
                )}
              >
                {t(filterIndex === 0 ? `if` : 'and')}
              </div>
              <Divider vertical={true} className="border-primary-500" />
              <DropdownMenu>
                <DropdownMenuTrigger
                  placeholder={t(`select-filter`)}
                  label={label}
                  variant="secondary"
                  className="w-full"
                />
                <DropdownMenuContent className="w-[235px]" align="start">
                  {remainingFilterOptions.map((item, index) => (
                    <DropdownMenuItem
                      key={index}
                      value={item.value || ''}
                      label={item.label}
                      onSelectOption={() => {
                        const filterValue =
                          item.value === FilterTypes.STATUSES ? [] : '';
                        selectedFilters[filterIndex] = { ...item, filterValue };
                        setSelectedFilters([...selectedFilters]);
                      }}
                    />
                  ))}
                </DropdownMenuContent>
              </DropdownMenu>
              <p className="typo-para-medium text-gray-600">is</p>
              <DropdownMenu
                onOpenChange={open => {
                  if (open) return handleFocusSearchInput();
                  setDebounceValue('');
                  setSearchValue('');
                }}
              >
                <DropdownMenuTrigger
                  disabled={isLoading || !filterType}
                  placeholder={t(`select-value`)}
                  label={handleGetLabelFilterValue(filterOption)}
                  variant="secondary"
                  className={cn('w-full max-w-[235px] truncate', {
                    capitalize: isStatusFilter
                  })}
                />
                <DropdownMenuContent
                  className={cn('w-[235px]', {
                    'pt-0 w-[300px]': isMaintainerFilter
                  })}
                  align="start"
                >
                  {isMaintainerFilter && (
                    <DropdownMenuSearch
                      ref={inputSearchRef}
                      value={debounceValue}
                      onChange={value => {
                        setDebounceValue(value);
                        debouncedSearch(value);
                        handleFocusSearchInput();
                      }}
                    />
                  )}
                  {valueOptions?.length > 0 ? (
                    valueOptions.map((item, index) => (
                      <DropdownMenuItem
                        isSelected={
                          isStatusFilter &&
                          Array.isArray(filterOption?.filterValue) &&
                          filterOption.filterValue?.includes(
                            item.value as string
                          )
                        }
                        isMultiselect={isStatusFilter}
                        key={index}
                        value={item.value as string}
                        label={item.label}
                        onSelectOption={value => {
                          handleChangeFilterValue(value, filterIndex);
                        }}
                      />
                    ))
                  ) : (
                    <div className="flex-center py-2.5 typo-para-medium text-gray-600">
                      {t('no-options-found')}
                    </div>
                  )}
                </DropdownMenuContent>
              </DropdownMenu>
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
        secondaryButton={
          <Button disabled={isDisabledSubmitButton} onClick={onConfirmHandler}>
            {t(`confirm`)}
          </Button>
        }
        primaryButton={
          <Button onClick={onClearFilters} variant="secondary">
            {t(`clear`)}
          </Button>
        }
      />
    </DialogModal>
  );
};

export default FilterExperimentModal;
