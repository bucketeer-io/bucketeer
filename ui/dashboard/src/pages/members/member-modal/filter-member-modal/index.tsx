import { useCallback, useEffect, useMemo, useState } from 'react';
import { useQueryTeams } from '@queries/teams';
import { getCurrentEnvironment, useAuth } from 'auth';
import useOptions, { FilterOption, FilterTypes } from 'hooks/use-options';
import { useTranslation } from 'i18n';
import { isEmpty } from 'utils/data-type';
import { cn } from 'utils/style';
import { IconPlus, IconTrash } from '@icons';
import { MembersFilters } from 'pages/members/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Dropdown, { DropdownOption } from 'components/dropdown';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';

export type FilterProps = {
  onSubmit: (v: Partial<MembersFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<MembersFilters>;
};

const FilterMemberModal = ({
  isOpen,
  filters,
  onSubmit,
  onClose,
  onClearFilters
}: FilterProps) => {
  const { t } = useTranslation(['common']);
  const { booleanOptions, memberFilterOptions, roleOptions } = useOptions();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const [selectedFilters, setSelectedFilters] = useState<FilterOption[]>([
    memberFilterOptions[0]
  ]);

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

    [memberFilterOptions, selectedFilters, memberFilterOptions]
  );

  const isDisabledSubmitButton = useMemo(() => {
    return !!selectedFilters.find(item => isEmpty(item.filterValue));
  }, [[...selectedFilters]]);

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
    [teams]
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
    [teamOptions]
  );

  const handleChangeFilterValue = useCallback(
    (value: string | number, filterIndex: number) => {
      const filterOption = selectedFilters[filterIndex];
      const { value: filterType, filterValue } = filterOption;
      const isTeamsFilter = filterType === FilterTypes.TEAMS;
      if (isTeamsFilter) {
        const values = filterValue as string[];
        if (Array.isArray(value) && value.length === 0) {
          selectedFilters[filterIndex] = {
            ...selectedFilters[filterIndex],
            filterValue: []
          };
          return setSelectedFilters([...selectedFilters]);
        }
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
    handleSetFilterOnInit();
  }, [filters]);

  return (
    <DialogModal
      className="w-[750px]"
      title={t('filters')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full items-start p-5 gap-y-4">
        {selectedFilters.map((filterOption, filterIndex) => {
          const { label, value: filterType } = filterOption;
          const isTeamsFilter = filterType === FilterTypes.TEAMS;
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
              <Dropdown
                placeholder={t(`select-filter`)}
                labelCustom={label}
                options={remainingFilterOptions as DropdownOption[]}
                value={filterType}
                onChange={value => {
                  const selected = remainingFilterOptions.find(
                    item => item.value === value
                  );
                  if (selected) {
                    selectedFilters[filterIndex] = {
                      ...selected,
                      filterValue: value === FilterTypes.TEAMS ? [] : ''
                    };
                    setSelectedFilters([...selectedFilters]);
                  }
                }}
              />

              <p className="typo-para-medium text-gray-600">is</p>
              <Dropdown
                placeholder={t(`select-value`)}
                labelCustom={handleGetLabelFilterValue(filterOption)}
                disabled={(isTeamsFilter && isLoadingTeams) || !filterType}
                loading={isTeamsFilter && isLoadingTeams}
                multiselect={isTeamsFilter}
                isSearchable={isTeamsFilter}
                value={
                  isTeamsFilter && Array.isArray(filterOption?.filterValue)
                    ? (filterOption.filterValue as string[])
                    : (filterOption.filterValue as string)
                }
                options={valueOptions as DropdownOption[]}
                onChange={value =>
                  handleChangeFilterValue(value as string, filterIndex)
                }
                className="w-full truncate"
                contentClassName={cn('w-[235px]', {
                  'pt-0 w-[300px]': isTeamsFilter,
                  'hidden-scroll': valueOptions?.length > 15
                })}
                menuContentSide="bottom"
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
                value: '',
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

export default FilterMemberModal;
