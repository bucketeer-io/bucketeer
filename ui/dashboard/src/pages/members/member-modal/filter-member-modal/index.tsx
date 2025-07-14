import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useQueryTeams } from '@queries/teams';
import { getCurrentEnvironment, useAuth } from 'auth';
import useOptions, { FilterOption, FilterTypes } from 'hooks/use-options';
import { useTranslation } from 'i18n';
import debounce from 'lodash/debounce';
import { isEmpty } from 'utils/data-type';
import { cn } from 'utils/style';
import { IconPlus, IconTrash } from '@icons';
import { MembersFilters } from 'pages/members/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSearch,
  DropdownMenuTrigger,
  DropdownOption
} from 'components/dropdown';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';
import DropdownList from 'elements/dropdown-list';

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
  const [searchValue, setSearchValue] = useState('');
  const [debounceValue, setDebounceValue] = useState('');

  const inputSearchRef = useRef<HTMLInputElement>(null);

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
      !remainingFilterOptions.length ||
      selectedFilters.length >= memberFilterOptions.length,

    [memberFilterOptions, selectedFilters, remainingFilterOptions]
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

  const debouncedSearch = useCallback(
    debounce(value => {
      setSearchValue(value);
    }, 500),
    []
  );

  const getValueOptions = useCallback(
    (filterOption: FilterOption) => {
      if (!filterOption.value) return [];
      const isRoleFilter = filterOption.value === FilterTypes.ROLE;
      const isTeamsFilter = filterOption.value === FilterTypes.TEAMS;

      if (isTeamsFilter)
        return teamOptions.filter(item =>
          searchValue
            ? item.value.toLowerCase().includes(searchValue.toLowerCase())
            : item
        );
      if (isRoleFilter) return roleOptions;
      return booleanOptions;
    },
    [searchValue, teamOptions]
  );

  const handleFocusSearchInput = useCallback(() => {
    let timerId: NodeJS.Timeout | null = null;
    if (timerId) clearTimeout(timerId);
    timerId = setTimeout(() => inputSearchRef?.current?.focus(), 50);
  }, []);

  const handleChangeFilterValue = useCallback(
    (value: string | number, filterIndex: number) => {
      const filterOption = selectedFilters[filterIndex];
      const { value: filterType, filterValue } = filterOption;
      const isTeamsFilter = filterType === FilterTypes.TEAMS;
      if (isTeamsFilter) {
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
                      value={item.value as string}
                      label={item.label}
                      onSelectOption={() => {
                        selectedFilters[filterIndex] = item;
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
                  placeholder={t(`select-value`)}
                  label={handleGetLabelFilterValue(filterOption)}
                  disabled={(isTeamsFilter && isLoadingTeams) || !filterType}
                  loading={isTeamsFilter && isLoadingTeams}
                  variant="secondary"
                  className="w-full truncate"
                />
                <DropdownMenuContent
                  className={cn('w-[235px]', {
                    'pt-0 w-[300px]': isTeamsFilter,
                    'hidden-scroll': valueOptions?.length > 15
                  })}
                  align="start"
                >
                  {isTeamsFilter && (
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
                  {valueOptions.length > 0 ? (
                    <DropdownList
                      isMultiselect={isTeamsFilter}
                      selectedOptions={
                        isTeamsFilter &&
                        Array.isArray(filterOption?.filterValue)
                          ? filterOption.filterValue
                          : undefined
                      }
                      options={valueOptions as DropdownOption[]}
                      onSelectOption={value =>
                        handleChangeFilterValue(value, filterIndex)
                      }
                    />
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
