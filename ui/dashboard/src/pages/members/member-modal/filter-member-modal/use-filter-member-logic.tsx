import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useQueryTeams } from '@queries/teams';
import { getCurrentEnvironment, useAuth } from 'auth';
import useOptions, { FilterOption, FilterTypes } from 'hooks/use-options';
import { isEmpty } from 'utils/data-type';
import { MembersFilters } from 'pages/members/types';

const useFilterMemberLogic = (
  onSubmit: (v: Partial<MembersFilters>) => void,
  filters?: Partial<MembersFilters>
) => {
  const { booleanOptions, memberFilterOptions, roleOptions } = useOptions();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const [selectedFilters, setSelectedFilters] = useState<FilterOption[]>([
    memberFilterOptions[0]
  ]);
  const initialized = useRef(false);

  const { data: teamCollection, isLoading: isLoadingTeams } = useQueryTeams({
    params: {
      cursor: String(0),
      organizationId: currentEnvironment.organizationId
    },
    enabled: !!selectedFilters.find(item => item.value === FilterTypes.TEAMS)
  });
  const teams = useMemo(() => teamCollection?.teams || [], [teamCollection]);
  const teamOptions = useMemo(
    () =>
      teams?.map(item => ({
        label: item.name,
        value: item.name
      })) || [],
    [teams]
  );

  const remainingFilterOptions = useMemo(
    () =>
      memberFilterOptions.filter(
        option => !selectedFilters.find(item => item.value === option.value)
      ),
    [selectedFilters, memberFilterOptions]
  );

  const isDisabledAddButton = useMemo(
    () =>
      !memberFilterOptions.length ||
      selectedFilters.length >= memberFilterOptions.length,

    [memberFilterOptions, selectedFilters]
  );

  const isDisabledSubmitButton = useMemo(() => {
    return !!selectedFilters.find(item => isEmpty(item.filterValue));
  }, [selectedFilters]);

  const handleGetLabelFilterValue = useCallback(
    (filterOption?: FilterOption) => {
      if (filterOption) {
        const { value: filterType, filterValue } = filterOption;
        const isTeamsFilter = filterType === FilterTypes.TEAMS;
        if (isTeamsFilter) {
          return (
            (Array.isArray(filterValue) &&
              teams.length &&
              filterValue
                .map(item => teams.find(team => team.name === item)?.name)
                ?.join(', ')) ||
            ''
          );
        }

        const currentOption = (
          filterType === FilterTypes.ENABLED ? booleanOptions : roleOptions
        ).find(item => item.value === filterValue);
        if (currentOption) return currentOption.label;
      }
      return '';
    },
    [teams, booleanOptions, roleOptions]
  );

  const getValueOptions = useCallback(
    (filterOption: FilterOption) => {
      if (!filterOption.value) return [];
      const isRoleFilter = filterOption.value === FilterTypes.ROLE;
      const isTeamsFilter = filterOption.value === FilterTypes.TEAMS;

      if (isTeamsFilter) return teamOptions;
      if (isRoleFilter) return roleOptions;
      return booleanOptions;
    },
    [teamOptions, roleOptions, booleanOptions]
  );

  const handleChangeFilterValue = useCallback(
    (value: string | number, filterIndex: number) => {
      const filterOption = selectedFilters[filterIndex];
      const { value: filterType, filterValue } = filterOption;
      const isTeamsFilter = filterType === FilterTypes.TEAMS;

      setSelectedFilters(prev =>
        prev.map((item, i) => {
          if (i !== filterIndex) return item;
          if (isTeamsFilter) {
            if (Array.isArray(value) && value.length === 0) {
              return { ...item, filterValue: [] };
            }
            const values = filterValue as string[];
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

  const handleSetFilterOnInit = useCallback(() => {
    if (filters) {
      const { disabled, organizationRole, teams } = filters || {};
      const filterTypeArr: FilterOption[] = [];
      const addFilterOption = (
        index: number,
        value: FilterOption['filterValue']
      ) => {
        if (!isEmpty(value)) {
          const option = memberFilterOptions[index];

          filterTypeArr.push({
            ...option,
            filterValue:
              option.value! === FilterTypes.TEAMS
                ? value
                : option.value! === FilterTypes.ROLE
                  ? value?.toString()
                  : value
                    ? 0
                    : 1
          });
        }
      };
      addFilterOption(0, disabled);
      addFilterOption(1, organizationRole);
      addFilterOption(2, teams);
      setSelectedFilters(
        filterTypeArr.length ? filterTypeArr : [memberFilterOptions[0]]
      );
    }
  }, [filters]);

  const onConfirmHandler = () => {
    const defaultFilters = {
      disabled: undefined,
      organizationRole: undefined,
      teams: undefined
    };

    const newFilters = {};

    selectedFilters.forEach(filter => {
      const isEnabledFilter = filter.value === FilterTypes.ENABLED;
      const isTeamsFilter = filter.value === FilterTypes.TEAMS;

      Object.assign(newFilters, {
        [isEnabledFilter ? 'disabled' : filter.value!]: isEnabledFilter
          ? !filter.filterValue
          : isTeamsFilter
            ? filter.filterValue
            : Number(filter?.filterValue)
      });
    });

    onSubmit({
      ...defaultFilters,
      ...newFilters
    });
  };

  useEffect(() => {
    if (initialized.current) return;
    initialized.current = true;
    handleSetFilterOnInit();
  }, [filters]);

  return {
    isLoadingTeams,
    booleanOptions,
    memberFilterOptions,
    roleOptions,
    selectedFilters,
    setSelectedFilters,
    remainingFilterOptions,
    isDisabledAddButton,
    isDisabledSubmitButton,
    handleGetLabelFilterValue,
    getValueOptions,
    handleChangeFilterValue,
    handleSetFilterOnInit,
    onConfirmHandler
  };
};

export default useFilterMemberLogic;
