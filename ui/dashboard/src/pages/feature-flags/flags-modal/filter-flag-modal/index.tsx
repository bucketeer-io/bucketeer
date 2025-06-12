import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useQueryAccounts } from '@queries/accounts';
import { useQueryTags } from '@queries/tags';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { debounce } from 'lodash';
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
  DropdownMenuTrigger
} from 'components/dropdown';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';

export type FilterProps = {
  isOpen: boolean;
  filters?: Partial<FlagFilters>;
  onSubmit: (v: Partial<FlagFilters>) => void;
  onClose: () => void;
  onClearFilters: () => void;
};

export interface Option {
  value: FilterTypes | undefined;
  label: string;
  filterValue?: string | boolean | number | string[];
}

export enum FilterTypes {
  HAS_EXPERIMENT = 'hasExperiment',
  HAS_PREREQUISITES = 'hasPrerequisites',
  MAINTAINER = 'maintainer',
  ENABLED = 'enabled',
  TAGS = 'tags'
}

export const filterOptions: Option[] = [
  {
    value: FilterTypes.HAS_EXPERIMENT,
    label: 'Has Experiment',
    filterValue: ''
  },
  {
    value: FilterTypes.HAS_PREREQUISITES,
    label: 'Has Prerequisites',
    filterValue: ''
  },
  {
    value: FilterTypes.MAINTAINER,
    label: 'Maintainer',
    filterValue: ''
  },
  {
    value: FilterTypes.ENABLED,
    label: 'Enabled',
    filterValue: ''
  },
  {
    value: FilterTypes.TAGS,
    label: 'Tags',
    filterValue: ''
  }
];

export const booleanOptions = [
  {
    value: 1,
    label: 'Yes'
  },
  {
    value: 0,
    label: 'No'
  }
];

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

  const inputSearchRef = useRef<HTMLInputElement>(null);

  const [selectedFilters, setSelectedFilters] = useState<Option[]>([
    filterOptions[0]
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

  const accounts = collection?.accounts || [];
  const tags = tagCollection?.tags || [];

  const remainingFilterOptions = useMemo(
    () =>
      filterOptions.filter(
        option => !selectedFilters.find(item => item.value === option.value)
      ),
    [selectedFilters, filterOptions]
  );

  const isDisabledAddButton = useMemo(
    () =>
      !remainingFilterOptions.length ||
      selectedFilters.length >= filterOptions.length,

    [filterOptions, selectedFilters, remainingFilterOptions]
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
    (filterOption: Option) => {
      if (!filterOption.value) return [];
      const isMaintainerFilter = filterOption.value === FilterTypes.MAINTAINER;
      const isTagFilter = filterOption.value === FilterTypes.TAGS;
      const isHaveSearchingDropdown = isMaintainerFilter || isTagFilter;
      if (isHaveSearchingDropdown) {
        const options = isMaintainerFilter
          ? accounts.map(item => ({
              label: item.email,
              value: item.email
            }))
          : tags.map(item => ({
              label: item.name,
              value: item.name
            }));
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
      const { maintainer, hasExperiment, hasPrerequisites, enabled, tags } =
        filters || {};
      const filterTypeArr: Option[] = [];
      const addFilterOption = (index: number, value: Option['filterValue']) => {
        if (!isEmpty(value)) {
          const option = filterOptions[index];

          filterTypeArr.push({
            ...option,
            filterValue: [FilterTypes.TAGS, FilterTypes.MAINTAINER].includes(
              option.value!
            )
              ? value
              : value
                ? 1
                : 0
          });
        }
      };
      addFilterOption(0, hasExperiment);
      addFilterOption(1, hasPrerequisites);
      addFilterOption(2, maintainer);
      addFilterOption(3, enabled);
      addFilterOption(4, tags);
      setSelectedFilters(
        filterTypeArr.length ? filterTypeArr : [filterOptions[0]]
      );
    }
  }, [filters]);

  const handleGetLabelFilterValue = useCallback(
    (filterOption?: Option) => {
      if (filterOption) {
        const { value: filterType, filterValue } = filterOption;
        const isMaintainerFilter = filterType === FilterTypes.MAINTAINER;
        const isTagFilter = filterType === FilterTypes.TAGS;
        return isMaintainerFilter
          ? filterValue || ''
          : isTagFilter
            ? (Array.isArray(filterValue) &&
                tags.length &&
                filterValue
                  .map(item => tags.find(tag => tag.name === item)?.name)
                  ?.join(', ')) ||
              ''
            : booleanOptions.find(item => item.value === filterValue)?.label ||
              '';
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
      tags: undefined
    };

    const newFilters = {};

    selectedFilters.forEach(filter => {
      const filterByText = [FilterTypes.MAINTAINER, FilterTypes.TAGS].includes(
        filter.value!
      );
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
      className="w-[665px]"
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
                          item.value === FilterTypes.TAGS ? [] : '';
                        selectedFilters[filterIndex] = { ...item, filterValue };
                        setSelectedFilters([...selectedFilters]);
                      }}
                    />
                  ))}
                </DropdownMenuContent>
              </DropdownMenu>
              <p className="typo-para-medium text-gray-600">{`is`}</p>
              <DropdownMenu
                onOpenChange={open => {
                  if (open) return handleFocusSearchInput();
                  setDebounceValue('');
                  setSearchValue('');
                }}
              >
                <DropdownMenuTrigger
                  disabled={isLoading || isLoadingTags || !filterType}
                  placeholder={t(`select-value`)}
                  label={handleGetLabelFilterValue(filterOption)}
                  variant="secondary"
                  className="w-full max-w-[235px] truncate"
                />
                <DropdownMenuContent
                  className={cn('w-[235px]', {
                    'pt-0 w-[300px]': isHaveSearchingDropdown
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
                    valueOptions.map((item, index) => (
                      <DropdownMenuItem
                        key={index}
                        isSelected={
                          isTagFilter &&
                          Array.isArray(filterOption?.filterValue) &&
                          filterOption.filterValue?.includes(
                            item.value as string
                          )
                        }
                        isMultiselect={isTagFilter}
                        value={item.value}
                        label={item.label}
                        className="flex items-center max-w-full truncate"
                        onSelectOption={value =>
                          handleChangeFilterValue(value, filterIndex)
                        }
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

export default FilterFlagModal;
