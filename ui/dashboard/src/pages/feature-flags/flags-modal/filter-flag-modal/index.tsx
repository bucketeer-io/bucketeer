import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useQueryAccounts } from '@queries/accounts';
import { useQueryTags } from '@queries/tags';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { debounce } from 'lodash';
import { isEmpty, isNotEmpty } from 'utils/data-type';
import { cn } from 'utils/style';
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
import DialogModal from 'components/modal/dialog';

export type FilterProps = {
  isOpen: boolean;
  filters?: Partial<FlagFilters>;
  onSubmit: (v: Partial<FlagFilters>) => void;
  onClose: () => void;
  onClearFilters: () => void;
};

export interface Option {
  value: string | number;
  label: string;
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
    label: 'Has Experiment'
  },
  {
    value: FilterTypes.HAS_PREREQUISITES,
    label: 'Has Prerequisites'
  },
  {
    value: FilterTypes.MAINTAINER,
    label: 'Maintainer'
  },
  {
    value: FilterTypes.ENABLED,
    label: 'Enabled'
  },
  {
    value: FilterTypes.TAGS,
    label: 'Tags'
  }
];

export const booleanOptions: Option[] = [
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
  const menuContentRef = useRef<HTMLDivElement>(null);

  const [selectedFilter, setSelectedFilter] = useState<Option>(
    filterOptions[0]
  );
  const [filterValue, setFilterValue] = useState<string | number | string[]>(
    ''
  );
  const [searchValue, setSearchValue] = useState('');
  const [debounceValue, setDebounceValue] = useState('');

  const isMaintainerFilter = useMemo(
    () => selectedFilter.value === FilterTypes.MAINTAINER,
    [selectedFilter]
  );

  const isTagFilter = useMemo(
    () => selectedFilter.value === FilterTypes.TAGS,
    [selectedFilter]
  );

  const isHaveSearchingDropdown = useMemo(
    () => isMaintainerFilter || isTagFilter,
    [isMaintainerFilter, isTagFilter]
  );

  const { data: collection, isLoading } = useQueryAccounts({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment?.id,
      organizationId: currentEnvironment?.organizationId
    },
    enabled: isMaintainerFilter
  });

  const { data: tagCollection, isLoading: isLoadingTags } = useQueryTags({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment?.id,
      entityType: 'FEATURE_FLAG'
    },
    enabled: isTagFilter
  });

  const accounts = collection?.accounts || [];
  const tags = tagCollection?.tags || [];

  const valueOptions = useMemo(() => {
    if (isHaveSearchingDropdown) {
      const options = isMaintainerFilter
        ? accounts.map(item => ({ label: item.email, value: item.email }))
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
  }, [
    isMaintainerFilter,
    accounts,
    isTagFilter,
    tags,
    searchValue,
    isHaveSearchingDropdown
  ]);

  const handleFocusSearchInput = useCallback(() => {
    let timerId: NodeJS.Timeout | null = null;
    if (timerId) clearTimeout(timerId);
    timerId = setTimeout(() => inputSearchRef?.current?.focus(), 50);
  }, []);

  const onConfirmHandler = useCallback(() => {
    const defaultFilters = {
      hasExperiment: undefined,
      hasPrerequisites: undefined,
      maintainer: undefined,
      enabled: undefined,
      tags: undefined
    };

    onSubmit({
      ...defaultFilters,
      [selectedFilter.value]: isHaveSearchingDropdown
        ? filterValue
        : !!filterValue
    });
  }, [isMaintainerFilter, isTagFilter, filterValue, isHaveSearchingDropdown]);

  const handleSetFilterOnInit = useCallback(() => {
    if (filters) {
      const { maintainer, hasExperiment, hasPrerequisites, enabled, tags } =
        filters || {};
      const isNotEmptyMaintainer = isNotEmpty(maintainer);
      const isNotTagMaintainer = isNotEmpty(tags);
      const isNotEmptyExperiment = isNotEmpty(hasExperiment);
      const isNotEmptyPrerequisites = isNotEmpty(hasPrerequisites);
      const isNotEmptyEnabled = isNotEmpty(enabled);

      if (isNotEmptyMaintainer) {
        setFilterValue(maintainer!);
        return setSelectedFilter(filterOptions[2]);
      }
      if (isNotTagMaintainer) {
        setFilterValue(tags!);
        return setSelectedFilter(filterOptions[4]);
      }
      if (
        isNotEmptyExperiment ||
        isNotEmptyPrerequisites ||
        isNotEmptyEnabled
      ) {
        setFilterValue(hasExperiment || hasPrerequisites || enabled ? 1 : 0);
        return setSelectedFilter(
          filterOptions[
            isNotEmptyExperiment ? 0 : isNotEmptyPrerequisites ? 1 : 3
          ]
        );
      }
      setSelectedFilter(filterOptions[0]);
    }
  }, [filters]);

  const handleGetLabelFilterValue = useCallback(() => {
    return isMaintainerFilter
      ? String(filterValue)
      : isTagFilter
        ? (Array.isArray(filterValue) &&
            tags.length &&
            filterValue
              .map(item => tags.find(tag => tag.name === item)?.name)
              ?.join(', ')) ||
          ''
        : booleanOptions.find(item => item.value === filterValue)?.label || '';
  }, [filterValue, isMaintainerFilter, isTagFilter, tags]);

  const handleChangeFilterValue = useCallback(
    (value: string | number) => {
      if (!isTagFilter) return setFilterValue(value);
      if (Array.isArray(filterValue)) {
        const isExistedTag = filterValue.includes(value as string);
        setFilterValue(
          isExistedTag
            ? filterValue.filter(item => item !== value)
            : [...filterValue, value as string]
        );
      }
    },
    [isTagFilter, filterValue]
  );

  const debouncedSearch = useCallback(
    debounce(value => {
      menuContentRef.current?.scrollTo({ top: 0, behavior: 'smooth' });
      setSearchValue(value);
    }, 500),
    []
  );

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
        <div className="flex items-center w-full h-12 gap-x-4">
          <div className="typo-para-small text-center py-[3px] px-4 rounded text-accent-pink-500 bg-accent-pink-50">
            {t(`if`)}
          </div>
          <Divider vertical={true} className="border-primary-500" />
          <DropdownMenu>
            <DropdownMenuTrigger
              placeholder={t(`select-filter`)}
              label={selectedFilter.label}
              variant="secondary"
              className="w-full"
            />
            <DropdownMenuContent className="w-[235px]" align="start">
              {filterOptions.map((item, index) => (
                <DropdownMenuItem
                  key={index}
                  value={item.value}
                  label={item.label}
                  onSelectOption={() => {
                    setSelectedFilter(item);
                    setFilterValue(item.value === FilterTypes.TAGS ? [] : '');
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
              disabled={isLoading || isLoadingTags}
              placeholder={t(`select-value`)}
              label={handleGetLabelFilterValue()}
              variant="secondary"
              className="w-full max-w-[235px] truncate"
            />
            <DropdownMenuContent
              ref={menuContentRef}
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
                      Array.isArray(filterValue) &&
                      filterValue.includes(item.value as string)
                    }
                    isMultiselect={isTagFilter}
                    value={item.value}
                    label={item.label}
                    className="flex items-center max-w-full truncate"
                    onSelectOption={value => handleChangeFilterValue(value)}
                  />
                ))
              ) : (
                <div className="flex-center py-2.5 typo-para-medium text-gray-600">
                  {t('no-options-found')}
                </div>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      <ButtonBar
        secondaryButton={
          <Button disabled={isEmpty(filterValue)} onClick={onConfirmHandler}>
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
