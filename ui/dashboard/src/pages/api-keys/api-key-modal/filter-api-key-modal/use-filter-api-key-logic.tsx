import { useCallback, useEffect, useMemo, useState } from 'react';
import { useQueryEnvironments } from '@queries/environments';
import { getCurrentEnvironment, useAuth } from 'auth';
import useOptions, { FilterOption, FilterTypes } from 'hooks/use-options';
import isNil from 'lodash/isNil';
import { isEmpty } from 'utils/data-type';
import { checkEnvironmentEmptyId, onFormatEnvironments } from 'utils/function';
import { APIKeysFilters } from 'pages/api-keys/types';
import { DropdownValue } from 'components/dropdown';

export type FilterProps = {
  onSubmit: (v: Partial<APIKeysFilters>) => void;
  filters?: Partial<APIKeysFilters>;
};

export interface Option {
  value: string;
  label: string;
}

const useFilterAPIKeyLogic = (
  onSubmit: (v: Partial<APIKeysFilters>) => void,
  filters?: Partial<APIKeysFilters>
) => {
  const { environmentEnabledFilterOptions, booleanOptions } = useOptions();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const [selectedFilters, setSelectedFilters] = useState<FilterOption[]>([
    environmentEnabledFilterOptions[0]
  ]);
  const remainingFilterOptions = useMemo(
    () =>
      environmentEnabledFilterOptions.filter(
        option => !selectedFilters.find(item => item.value === option.value)
      ),
    [selectedFilters, environmentEnabledFilterOptions]
  );

  const isDisabledAddButton = useMemo(
    () =>
      !remainingFilterOptions.length ||
      selectedFilters.length >= environmentEnabledFilterOptions.length,

    [environmentEnabledFilterOptions, selectedFilters, remainingFilterOptions]
  );

  const isDisabledSubmitButton = useMemo(
    () => !!selectedFilters.find(item => isEmpty(item.filterValue)),
    [selectedFilters]
  );

  const { data: environmentCollection, isLoading: isLoadingEnvironments } =
    useQueryEnvironments({
      params: {
        cursor: '0',
        organizationId: currentEnvironment.organizationId
      },
      enabled:
        !!currentEnvironment.organizationId &&
        !!selectedFilters.find(
          item => item.value === FilterTypes.ENVIRONMENT_IDs
        )
    });
  const environments = useMemo(
    () => environmentCollection?.environments || [],
    [environmentCollection]
  );

  const { emptyEnvironmentId, formattedEnvironments } =
    onFormatEnvironments(environments);

  const environmentOptions = useMemo(
    () =>
      formattedEnvironments.map(item => ({
        label: item.name,
        value: item.id
      })),
    [formattedEnvironments]
  );

  const getValueOptions = useCallback(
    (filterOption: FilterOption) => {
      if (!filterOption.value) return [];
      const isEnvironmentFilter =
        filterOption.value === FilterTypes.ENVIRONMENT_IDs;

      if (isEnvironmentFilter) {
        return environmentOptions;
      }

      return booleanOptions;
    },
    [booleanOptions, environmentOptions]
  );

  const handleGetLabelFilterValue = useCallback(
    (filterOption?: FilterOption) => {
      if (filterOption) {
        const { value: filterType, filterValue } = filterOption;
        const isEnvironmentFilter = filterType === FilterTypes.ENVIRONMENT_IDs;

        if (isEnvironmentFilter) {
          return (
            (Array.isArray(filterValue) &&
              filterValue
                .map(
                  item =>
                    environmentOptions.find(env => env.value === item)?.label
                )
                ?.join(', ')) ||
            ''
          );
        }
        return (
          booleanOptions.find(item => item.value === filterValue)?.label || ''
        );
      }
      return '';
    },
    [booleanOptions, environmentOptions]
  );

  const handleChangeFilterValue = useCallback(
    (value: string | number, filterIndex: number) => {
      const filterOption = selectedFilters[filterIndex];
      const { value: filterType, filterValue } = filterOption;
      const isEnvironmentFilter = filterType === FilterTypes.ENVIRONMENT_IDs;
      let newFilterValue: string | number | string[] | number[] = value;
      if (isEnvironmentFilter) {
        if (Array.isArray(newFilterValue) && newFilterValue.length === 0) {
          selectedFilters[filterIndex] = {
            ...selectedFilters[filterIndex],
            filterValue: value
          };
          return setSelectedFilters([...selectedFilters]);
        }

        const values = filterValue as string[];
        const isExisted = values.find(item => item === value);
        const newValue: string[] = isExisted
          ? values.filter(item => item !== value)
          : [...values, value as string];
        newFilterValue = newValue;
      }
      selectedFilters[filterIndex] = {
        ...selectedFilters[filterIndex],
        filterValue: newFilterValue
      };
      setSelectedFilters([...selectedFilters]);
    },
    [selectedFilters]
  );

  const onConfirmHandler = useCallback(() => {
    const defaultFilters = {
      disabled: undefined,
      environmentIds: undefined
    };
    const newFilters = {};

    selectedFilters.forEach(filter => {
      const isEnvironmentFilter = filter.value === FilterTypes.ENVIRONMENT_IDs;
      const isEnabledFilter = filter.value === FilterTypes.ENABLED;
      Object.assign(newFilters, {
        [isEnabledFilter ? 'disabled' : FilterTypes.ENVIRONMENT_IDs]:
          isEnvironmentFilter
            ? Array.isArray(filter.filterValue)
              ? filter.filterValue.map(item => checkEnvironmentEmptyId(item))
              : []
            : isEmpty(filter.filterValue)
              ? undefined
              : !filter.filterValue
      });
    });

    onSubmit({
      ...defaultFilters,
      ...newFilters
    });
  }, [selectedFilters]);

  const handleSetFilterOnInit = useCallback(() => {
    if (filters) {
      const { disabled, environmentIds } = filters || {};
      const filterTypeArr: FilterOption[] = [];

      const addFilterOption = (
        index: number,
        value: FilterOption['filterValue']
      ) => {
        if (!isNil(value)) {
          const option = environmentEnabledFilterOptions[index];
          filterTypeArr.push({
            ...option,
            filterValue:
              option.value === FilterTypes.ENVIRONMENT_IDs
                ? Array.isArray(value)
                  ? value.map(item => item || emptyEnvironmentId)
                  : []
                : value
                  ? 0
                  : 1
          });
        }
      };
      addFilterOption(0, environmentIds);
      addFilterOption(1, disabled);

      setSelectedFilters(
        filterTypeArr.length
          ? filterTypeArr
          : [environmentEnabledFilterOptions[0]]
      );
    }
  }, [filters, emptyEnvironmentId]);

  const handleChangeOption = (val: DropdownValue, filterIndex: number) => {
    const selectionOption = remainingFilterOptions.find(
      item => item.value === val
    );
    if (!selectionOption) return;
    const filterValue =
      selectionOption.value === FilterTypes.ENVIRONMENT_IDs ? [] : '';
    selectedFilters[filterIndex] = { ...selectionOption, filterValue };
    setSelectedFilters([...selectedFilters]);
  };

  useEffect(() => {
    handleSetFilterOnInit();
  }, [filters]);

  return {
    environmentEnabledFilterOptions,
    booleanOptions,
    selectedFilters,
    setSelectedFilters,
    remainingFilterOptions,
    isDisabledAddButton,
    isDisabledSubmitButton,
    environmentCollection,
    isLoadingEnvironments,
    environmentOptions,
    getValueOptions,
    handleChangeFilterValue,
    handleGetLabelFilterValue,
    onConfirmHandler,
    handleChangeOption
  };
};

export default useFilterAPIKeyLogic;
