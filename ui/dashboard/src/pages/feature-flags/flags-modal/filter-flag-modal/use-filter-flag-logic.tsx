import { useCallback, useEffect, useMemo, useState } from 'react';
import { useQueryAccounts } from '@queries/accounts';
import { useQueryTags } from '@queries/tags';
import { getCurrentEnvironment, useAuth } from 'auth';
import useOptions, { FilterOption, FilterTypes } from 'hooks/use-options';
import { isEmpty } from 'utils/data-type';
import { FlagFilters } from 'pages/feature-flags/types';
import { DropdownValue } from 'components/dropdown';

const useFilterFlagLogic = (
  filters?: Partial<FlagFilters>,
  onSubmit?: (v: Partial<FlagFilters>) => void
) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { booleanOptions, flagFilterOptions, flagStatusOptions } = useOptions();

  const [selectedFilters, setSelectedFilters] = useState<FilterOption[]>([
    flagFilterOptions[0]
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

  const { data: tagCollection, isLoading: isLoadingTags } = useQueryTags({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment?.id,
      entityType: 'FEATURE_FLAG'
    },
    enabled: !!selectedFilters?.find(item => item.value === FilterTypes.TAGS)
  });

  const accounts = useMemo(() => collection?.accounts || [], [collection]);
  const tags = useMemo(() => tagCollection?.tags || [], [tagCollection]);

  const remainingFilterOptions = useMemo(
    () =>
      flagFilterOptions.filter(
        option => !selectedFilters.find(item => item.value === option.value)
      ),
    [selectedFilters, flagFilterOptions]
  );

  const isDisabledAddButton = useMemo(
    () =>
      !remainingFilterOptions.length ||
      selectedFilters.length >= flagFilterOptions.length,

    [flagFilterOptions, selectedFilters, remainingFilterOptions]
  );

  const isDisabledSubmitButton = useMemo(
    () => !!selectedFilters.find(item => isEmpty(item.filterValue)),
    [selectedFilters]
  );

  const getValueOptions = useCallback(
    (filterOption: FilterOption) => {
      if (!filterOption.value) return [];
      const isMaintainerFilter = filterOption.value === FilterTypes.MAINTAINER;
      const isTagFilter = filterOption.value === FilterTypes.TAGS;
      const isStatusFilter = filterOption.value === FilterTypes.STATUS;

      const isHaveSearchingDropdown =
        isMaintainerFilter || isTagFilter || isStatusFilter;
      if (isHaveSearchingDropdown) {
        const options = isMaintainerFilter
          ? accounts.map(item => ({
              label: item.email,
              value: item.email
            }))
          : isTagFilter
            ? tags.map(item => ({
                label: item.name,
                value: item.name
              }))
            : flagStatusOptions;

        return options;
      }

      return booleanOptions;
    },
    [accounts, tags, flagStatusOptions, booleanOptions]
  );

  const handleSetFilterOnInit = useCallback(() => {
    if (filters) {
      const {
        maintainer,
        hasExperiment,
        hasPrerequisites,
        enabled,
        tags,
        status,
        hasFeatureFlagAsRule,
        hasActiveAutoOps,
        hasFinishedAutoOps
      } = filters || {};
      const filterTypeArr: FilterOption[] = [];
      const addFilterOption = (
        filterType: FilterTypes,
        value: FilterOption['filterValue']
      ) => {
        const option = flagFilterOptions.find(
          item => item.value === filterType
        );
        if (option && !isEmpty(value)) {
          filterTypeArr.push({
            ...option,
            filterValue: [
              FilterTypes.TAGS,
              FilterTypes.MAINTAINER,
              FilterTypes.STATUS
            ].includes(filterType)
              ? value
              : value
                ? 1
                : 0
          });
        }
      };
      addFilterOption(FilterTypes.HAS_PREREQUISITES, hasPrerequisites);
      addFilterOption(FilterTypes.HAS_RULE, hasFeatureFlagAsRule);
      addFilterOption(FilterTypes.HAS_ACTIVE_AUTO_OPS, hasActiveAutoOps);
      addFilterOption(FilterTypes.HAS_FINISHED_AUTO_OPS, hasFinishedAutoOps);
      addFilterOption(FilterTypes.HAS_EXPERIMENT, hasExperiment);
      addFilterOption(FilterTypes.ENABLED, enabled);
      addFilterOption(FilterTypes.TAGS, tags);
      addFilterOption(FilterTypes.STATUS, status);
      addFilterOption(FilterTypes.MAINTAINER, maintainer);

      setSelectedFilters(
        filterTypeArr.length ? filterTypeArr : [flagFilterOptions[0]]
      );
    }
  }, [filters, flagFilterOptions]);

  const handleGetLabelFilterValue = useCallback(
    (filterOption?: FilterOption) => {
      if (filterOption) {
        const { value: filterType, filterValue } = filterOption;
        const isMaintainerFilter = filterType === FilterTypes.MAINTAINER;
        const isTagFilter = filterType === FilterTypes.TAGS;
        const isStatusFilter = filterType === FilterTypes.STATUS;

        return isMaintainerFilter
          ? filterValue || ''
          : isTagFilter
            ? (Array.isArray(filterValue) &&
                tags.length &&
                filterValue
                  .map(item => tags.find(tag => tag.name === item)?.name)
                  ?.join(', ')) ||
              ''
            : (isStatusFilter ? flagStatusOptions : booleanOptions).find(
                item => item.value === filterValue
              )?.label || '';
      }
      return '';
    },
    [tags, flagStatusOptions, booleanOptions]
  );

  const handleChangeFilterValue = useCallback(
    (value: DropdownValue, filterIndex: number) => {
      const filterOption = selectedFilters[filterIndex];
      const { value: filterType, filterValue } = filterOption;
      const isTagOption = filterType === FilterTypes.TAGS;

      setSelectedFilters(prev =>
        prev.map((item, i) => {
          if (i !== filterIndex) return item;
          if (isTagOption) {
            const values = filterValue as string[];
            if (Array.isArray(value) && isEmpty(value)) {
              return { ...item, filterValue: value };
            }
            const isExisted = values.find(v => v === value);
            const newValue: string[] = isExisted
              ? values.filter(v => v !== value)
              : [...values, value as string];
            return { ...item, filterValue: newValue };
          }
          return { ...item, filterValue: value };
        })
      );
    },
    [selectedFilters]
  );

  const handleChangeOption = (value: string, filterIndex: number) => {
    const selectedOption = flagFilterOptions.find(item => item.value === value);
    if (selectedOption) {
      const filterValue = selectedOption.value === FilterTypes.TAGS ? [] : '';
      setSelectedFilters(prev =>
        prev.map((item, i) =>
          i === filterIndex ? { ...selectedOption, filterValue } : item
        )
      );
    }
  };

  const onConfirmHandler = useCallback(() => {
    const defaultFilters = {
      hasExperiment: undefined,
      hasPrerequisites: undefined,
      maintainer: undefined,
      enabled: undefined,
      tags: undefined,
      status: undefined,
      hasFeatureFlagAsRule: undefined,
      hasActiveAutoOps: undefined,
      hasFinishedAutoOps: undefined
    };

    const newFilters = {};

    selectedFilters.forEach(filter => {
      const filterByText = [
        FilterTypes.MAINTAINER,
        FilterTypes.TAGS,
        FilterTypes.STATUS
      ].includes(filter.value as FilterTypes);
      Object.assign(newFilters, {
        [filter.value!]: filterByText
          ? filter.filterValue
          : !!filter.filterValue
      });
    });

    onSubmit?.({
      ...defaultFilters,
      ...newFilters
    });
  }, [selectedFilters, onSubmit]);

  useEffect(() => {
    handleSetFilterOnInit();
  }, [filters]);

  return {
    selectedFilters,
    accounts,
    tags,
    isLoading,
    isLoadingTags,
    remainingFilterOptions,
    isDisabledAddButton,
    isDisabledSubmitButton,
    getValueOptions,
    handleGetLabelFilterValue,
    handleChangeFilterValue,
    handleChangeOption,
    setSelectedFilters,
    onConfirmHandler
  };
};

export default useFilterFlagLogic;
