import { useCallback, useEffect, useMemo, useState } from 'react';
import { getEditorEnvironments, useAuth } from 'auth';
import useOptions, { FilterOption, FilterTypes } from 'hooks/use-options';
import isNil from 'lodash/isNil';
import { isEmpty } from 'utils/data-type';
import { checkEnvironmentEmptyId, onFormatEnvironments } from 'utils/function';
import { NotificationFilters } from 'pages/notifications/types';
import { DropdownValue } from 'components/dropdown';

export interface Option {
  value: string;
  label: string;
}

const useFilterPushLogic = (
  onSubmit: (v: Partial<NotificationFilters>) => void,
  filters?: Partial<NotificationFilters>
) => {
  const { environmentEnabledFilterOptions, booleanOptions } = useOptions();
  const { consoleAccount } = useAuth();
  const { editorEnvironments } = getEditorEnvironments(consoleAccount!);

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

  const { emptyEnvironmentId, formattedEnvironments } =
    onFormatEnvironments(editorEnvironments);

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
      return isEnvironmentFilter ? environmentOptions : booleanOptions;
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
      let newFilterValue: string | number | string[] = value;
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

  const handleChangeOption = (value: DropdownValue, filterIndex: number) => {
    const selected = environmentEnabledFilterOptions.find(
      item => item.value === value
    );
    const filterValue = value === FilterTypes.ENVIRONMENT_IDs ? [] : '';
    selectedFilters[filterIndex] = {
      ...selected!,
      filterValue
    };
    setSelectedFilters([...selectedFilters]);
  };

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

  useEffect(() => {
    handleSetFilterOnInit();
  }, [filters]);

  return {
    remainingFilterOptions,
    isDisabledAddButton,
    isDisabledSubmitButton,
    environmentOptions,

    environmentEnabledFilterOptions,
    booleanOptions,
    selectedFilters,
    handleChangeOption,
    setSelectedFilters,
    getValueOptions,
    handleGetLabelFilterValue,
    handleChangeFilterValue,
    onConfirmHandler,
    handleSetFilterOnInit
  };
};

export default useFilterPushLogic;
