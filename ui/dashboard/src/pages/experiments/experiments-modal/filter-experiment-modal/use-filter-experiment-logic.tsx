import { useCallback, useEffect, useMemo, useState } from 'react';
import { useQueryAccounts } from '@queries/accounts';
import { getCurrentEnvironment, useAuth } from 'auth';
import useOptions, { FilterOption, FilterTypes } from 'hooks/use-options';
import { isEmpty } from 'utils/data-type';
import { ExperimentFilters } from 'pages/experiments/types';
import { DropdownValue } from 'components/dropdown';

const useExperimentFilterLogic = (
  filters?: Partial<ExperimentFilters>,
  onSubmit?: (v: Partial<ExperimentFilters>) => void
) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { experimentStatusOptions, experimentFilterOptions } = useOptions();
  const [selectedFilters, setSelectedFilters] = useState<FilterOption[]>([
    experimentFilterOptions[0]
  ]);
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
  const getValueOptions = useCallback(
    (filterOption: FilterOption) => {
      if (!filterOption.value) return [];
      const isMaintainerFilter = filterOption.value === FilterTypes.MAINTAINER;
      if (isMaintainerFilter) {
        const options = accounts.map(item => ({
          label: item.email,
          value: item.email
        }));
        return options;
      }
      return experimentStatusOptions;
    },
    [accounts, experimentStatusOptions]
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
          filterTypeArr.push({ ...option, filterValue: value });
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
        if (Array.isArray(value) && isEmpty(value)) {
          selectedFilters[filterIndex] = {
            ...selectedFilters[filterIndex],
            filterValue: []
          };
          return setSelectedFilters([...selectedFilters]);
        }
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
  const handleChangeOption = (value: DropdownValue, filterIndex: number) => {
    {
      const selectedOption = experimentFilterOptions.find(
        item => item.value === value
      );
      if (selectedOption) {
        const filterValue =
          selectedOption.value === FilterTypes.STATUSES ? [] : '';
        selectedFilters[filterIndex] = { ...selectedOption, filterValue };
        setSelectedFilters([...selectedFilters]);
      }
    }
  };
  const onConfirmHandler = () => {
    const defaultFilters = { statuses: undefined, maintainer: undefined };
    const newFilters = {};
    selectedFilters.forEach(filter => {
      Object.assign(newFilters, { [filter.value!]: filter.filterValue });
    });
    onSubmit?.({ ...defaultFilters, ...newFilters, isFilter: true });
  };
  useEffect(() => {
    if (filters?.isFilter || filters?.filterBySummary) {
      handleSetFilterOnInit();
    }
  }, [filters]);
  return {
    selectedFilters,
    isDisabledAddButton,
    isDisabledSubmitButton,
    isLoading,
    remainingFilterOptions,
    setSelectedFilters,
    onConfirmHandler,
    getValueOptions,
    handleGetLabelFilterValue,
    handleChangeFilterValue,
    handleChangeOption
  };
};
export default useExperimentFilterLogic;
