import { useCallback, useEffect, useMemo, useState } from 'react';
import { useQueryTags } from '@queries/tags';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useAccountsLoader } from 'hooks/use-accounts-loading-more';
import useOptions, { FilterOption, FilterTypes } from 'hooks/use-options';
import { useTranslation } from 'i18n';
import { isEmpty } from 'utils/data-type';
import { cn } from 'utils/style';
import { IconPlus, IconTrash } from '@icons';
import { FlagFilters } from 'pages/feature-flags/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Dropdown, { DropdownOption, DropdownValue } from 'components/dropdown';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';
import DropdownMenuWithSearch from 'elements/dropdown-with-search';

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

  const [selectedFilters, setSelectedFilters] = useState<FilterOption[]>([
    flagFilterOptions[0]
  ]);

  const hasMaintainerFilter = !!selectedFilters.find(
    item => item.value === FilterTypes.MAINTAINER
  );

  const selectedMaintainer = useMemo(() => {
    const filter = selectedFilters.find(
      item => item.value === FilterTypes.MAINTAINER
    );
    return typeof filter?.filterValue === 'string' ? filter.filterValue : '';
  }, [selectedFilters]);

  const {
    emailOptions,
    isInitialLoading: isLoadingAccounts,
    isLoadingMore,
    isSearching,
    hasMore,
    loadMore,
    onSearchChange: onAccountSearchChange,
    getAccountLabel
  } = useAccountsLoader({
    organizationId: currentEnvironment.organizationId,
    environmentId: currentEnvironment.id,
    enabled: hasMaintainerFilter,
    preloadEmails: selectedMaintainer ? [selectedMaintainer] : []
  });

  const { data: tagCollection, isLoading: isLoadingTags } = useQueryTags({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment.id,
      entityType: 'FEATURE_FLAG'
    },
    enabled: !!selectedFilters?.find(item => item.value === FilterTypes.TAGS)
  });

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
      const isTagFilter = filterOption.value === FilterTypes.TAGS;
      const isStatusFilter = filterOption.value === FilterTypes.STATUS;

      if (isTagFilter) {
        return tags.map(item => ({ label: item.name, value: item.name }));
      }
      if (isStatusFilter) return flagStatusOptions;
      return booleanOptions;
    },
    [tags, flagStatusOptions, booleanOptions]
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
          ? getAccountLabel(filterValue as string)
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
    [tags, flagStatusOptions, booleanOptions, getAccountLabel]
  );

  const handleChangeFilterValue = useCallback(
    (value: DropdownValue, filterIndex: number) => {
      const filterOption = selectedFilters[filterIndex];
      const { value: filterType, filterValue } = filterOption;
      const isTagOption = filterType === FilterTypes.TAGS;
      if (isTagOption) {
        const values = filterValue as string[];
        if (Array.isArray(value) && isEmpty(value)) {
          return setSelectedFilters(prev => {
            const next = [...prev];
            next[filterIndex] = { ...next[filterIndex], filterValue: value };
            return next;
          });
        }
        const isExisted = values.find(item => item === value);
        const newValue: string[] = isExisted
          ? values.filter(item => item !== value)
          : [...values, value as string];
        return setSelectedFilters(prev => {
          const next = [...prev];
          next[filterIndex] = { ...next[filterIndex], filterValue: newValue };
          return next;
        });
      }
      setSelectedFilters(prev => {
        const next = [...prev];
        next[filterIndex] = { ...next[filterIndex], filterValue: value };
        return next;
      });
    },
    [selectedFilters]
  );

  const handleChangeOption = (value: string, filterIndex: number) => {
    const selectedOption = flagFilterOptions.find(item => item.value === value);
    if (selectedOption) {
      const filterValue = selectedOption.value === FilterTypes.TAGS ? [] : '';
      setSelectedFilters(prev => {
        const next = [...prev];
        next[filterIndex] = { ...selectedOption, filterValue };
        return next;
      });
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
  }, [selectedFilters, onSubmit]);

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
              <Dropdown
                placeholder={t(`select-filter`)}
                labelCustom={label}
                options={remainingFilterOptions as DropdownOption[]}
                value={filterType}
                onChange={value => {
                  handleChangeOption(value as string, filterIndex);
                }}
                className="w-full truncate"
                contentClassName="w-[270px]"
              />
              <p className="typo-para-medium text-gray-600">is</p>
              {isHaveSearchingDropdown ? (
                <DropdownMenuWithSearch
                  disabled={
                    (isTagFilter && isLoadingTags) ||
                    (isMaintainerFilter && isLoadingAccounts) ||
                    !filterType
                  }
                  isLoading={
                    (isTagFilter && isLoadingTags) ||
                    (isMaintainerFilter && isLoadingAccounts)
                  }
                  isMultiselect={isTagFilter}
                  placeholder={t(`select-value`)}
                  itemSelected={
                    isMaintainerFilter
                      ? (filterOption?.filterValue as string)
                      : undefined
                  }
                  selectedOptions={
                    isTagFilter
                      ? (filterOption?.filterValue as string[])
                      : undefined
                  }
                  label={handleGetLabelFilterValue(filterOption) as string}
                  options={
                    isMaintainerFilter
                      ? emailOptions
                      : (valueOptions as DropdownOption[])
                  }
                  isHasMore={
                    isMaintainerFilter && hasMaintainerFilter && hasMore
                  }
                  isLoadingMore={isMaintainerFilter && isLoadingMore}
                  onHasMoreOptions={isMaintainerFilter ? loadMore : undefined}
                  isSearching={isMaintainerFilter && isSearching}
                  onSearchChange={
                    isMaintainerFilter ? onAccountSearchChange : undefined
                  }
                  onSelectOption={value =>
                    handleChangeFilterValue(value as DropdownValue, filterIndex)
                  }
                  triggerClassName="w-full truncate"
                  contentClassName={cn('w-[300px]', {
                    'hidden-scroll': valueOptions?.length > 15
                  })}
                />
              ) : (
                <Dropdown
                  disabled={!filterType}
                  placeholder={t(`select-value`)}
                  value={filterOption?.filterValue as DropdownValue}
                  labelCustom={handleGetLabelFilterValue(filterOption)}
                  options={valueOptions as DropdownOption[]}
                  onChange={value =>
                    handleChangeFilterValue(value as DropdownValue, filterIndex)
                  }
                  className="w-full truncate"
                  contentClassName="w-[235px]"
                />
              )}

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
