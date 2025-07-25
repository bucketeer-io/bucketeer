import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useQueryAccounts } from '@queries/accounts';
import { useQueryTags } from '@queries/tags';
import { getCurrentEnvironment, useAuth } from 'auth';
import useOptions, { FilterOption, FilterTypes } from 'hooks/use-options';
import { useTranslation } from 'i18n';
import debounce from 'lodash/debounce';
import { isEmpty } from 'utils/data-type';
import { cn } from 'utils/style';
import { IconPlus, IconTrash } from '@icons';
import { FlagFilters } from 'pages/feature-flags/types';
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
  isOpen: boolean;
  filters?: Partial<FlagFilters>;
  onSubmit: (v: Partial<FlagFilters>) => void;
  onClose: () => void;
  onClearFilters: () => void;
};

const FilterFlagModal = ({
  isOpen,
  filters,
  onSubmit,
  onClose,
  onClearFilters
}: FilterProps) => {
  const { t } = useTranslation(['common']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { booleanOptions, flagFilterOptions, flagStatusOptions } = useOptions();
  const inputSearchRef = useRef<HTMLInputElement>(null);

  const [selectedFilters, setSelectedFilters] = useState<FilterOption[]>([
    flagFilterOptions[0]
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

        return options?.filter(item =>
          searchValue
            ? item.value.toLowerCase().includes(searchValue.toLowerCase())
            : item
        );
      }

      return booleanOptions;
    },
    [accounts, tags, searchValue]
  );

  const handleFocusSearchInput = useCallback(() => {
    let timerId: NodeJS.Timeout | null = null;
    if (timerId) clearTimeout(timerId);
    timerId = setTimeout(() => inputSearchRef?.current?.focus(), 50);
  }, []);

  const handleSetFilterOnInit = useCallback(() => {
    if (filters) {
      const {
        maintainer,
        hasExperiment,
        hasPrerequisites,
        enabled,
        tags,
        status,
        hasFeatureFlagAsRule
      } = filters || {};
      const filterTypeArr: FilterOption[] = [];
      const addFilterOption = (
        index: number,
        value: FilterOption['filterValue']
      ) => {
        if (!isEmpty(value)) {
          const option = flagFilterOptions[index];

          filterTypeArr.push({
            ...option,
            filterValue: [
              FilterTypes.TAGS,
              FilterTypes.MAINTAINER,
              FilterTypes.STATUS
            ].includes(option.value as FilterTypes)
              ? value
              : value
                ? 1
                : 0
          });
        }
      };
      addFilterOption(0, hasPrerequisites);
      addFilterOption(1, hasFeatureFlagAsRule);
      addFilterOption(2, hasExperiment);
      addFilterOption(3, enabled);
      addFilterOption(4, tags);
      addFilterOption(5, status);
      addFilterOption(6, maintainer);

      setSelectedFilters(
        filterTypeArr.length ? filterTypeArr : [flagFilterOptions[0]]
      );
    }
  }, [filters]);

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
    [tags]
  );

  const handleChangeFilterValue = useCallback(
    (value: string | number, filterIndex: number) => {
      const filterOption = selectedFilters[filterIndex];
      const { value: filterType, filterValue } = filterOption;
      const isTagOption = filterType === FilterTypes.TAGS;
      if (isTagOption) {
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

  const onConfirmHandler = useCallback(() => {
    const defaultFilters = {
      hasExperiment: undefined,
      hasPrerequisites: undefined,
      maintainer: undefined,
      enabled: undefined,
      tags: undefined,
      status: undefined,
      hasFeatureFlagAsRule: undefined
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

    onSubmit({
      ...defaultFilters,
      ...newFilters
    });
  }, [selectedFilters]);

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
          const isMaintainerFilter = filterType === FilterTypes.MAINTAINER;
          const isTagFilter = filterType === FilterTypes.TAGS;
          const isHaveSearchingDropdown = isMaintainerFilter || isTagFilter;
          const valueOptions = getValueOptions(filterOption);
          return (
            <div
              className="flex items-center w-full h-12 gap-x-4"
              key={filterIndex}
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
                  className="w-full truncate"
                />
                <DropdownMenuContent className="w-[270px]" align="start">
                  {remainingFilterOptions.map((item, index) => (
                    <DropdownMenuItem
                      key={index}
                      value={item.value || ''}
                      label={item.label}
                      onSelectOption={() => {
                        const filterValue =
                          item.value === FilterTypes.TAGS ? [] : '';
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
                  disabled={
                    (isTagFilter && isLoadingTags) ||
                    (isMaintainerFilter && isLoading) ||
                    !filterType
                  }
                  loading={
                    (isTagFilter && isLoadingTags) ||
                    (isMaintainerFilter && isLoading)
                  }
                  placeholder={t(`select-value`)}
                  label={handleGetLabelFilterValue(filterOption)}
                  variant="secondary"
                  className="w-full truncate"
                />
                <DropdownMenuContent
                  className={cn('w-[235px]', {
                    'pt-0 w-[300px]': isHaveSearchingDropdown,
                    'hidden-scroll': valueOptions?.length > 15
                  })}
                  align="start"
                >
                  {isHaveSearchingDropdown && (
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
                    <DropdownList
                      isMultiselect={isTagFilter}
                      selectedOptions={
                        isTagFilter && Array.isArray(filterOption?.filterValue)
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

export default FilterFlagModal;
