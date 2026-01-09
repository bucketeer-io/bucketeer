import { useCallback, useEffect, useMemo, useState } from 'react';
import { useQueryAccounts } from '@queries/accounts';
import { useQueryTags } from '@queries/tags';
import { getCurrentEnvironment, useAuth } from 'auth';
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
import SlideModal from 'components/modal/slide';

export type FilterProps = {
  isOpen: boolean;
  filters?: Partial<FlagFilters>;
  onSubmit: (v: Partial<FlagFilters>) => void;
  onClose: () => void;
  onClearFilters: () => void;
};

const FilterFlagSlideModal = ({
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
    [accounts, tags]
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
    (value: DropdownValue, filterIndex: number) => {
      const filterOption = selectedFilters[filterIndex];
      const { value: filterType, filterValue } = filterOption;
      const isTagOption = filterType === FilterTypes.TAGS;
      if (isTagOption) {
        const values = filterValue as string[];
        if (Array.isArray(value) && isEmpty(value)) {
          selectedFilters[filterIndex] = {
            ...selectedFilters[filterIndex],
            filterValue: value
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

  const handleChangeOption = (value: string, filterIndex: number) => {
    const selectedOption = flagFilterOptions.find(item => item.value === value);
    if (selectedOption) {
      const filterValue = selectedOption.value === FilterTypes.TAGS ? [] : '';

      selectedFilters[filterIndex] = {
        ...selectedOption,
        filterValue
      };
      setSelectedFilters([...selectedFilters]);
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
  }, [selectedFilters]);

  useEffect(() => {
    handleSetFilterOnInit();
  }, [filters]);

  return (
    <SlideModal title={t('filters')} isOpen={isOpen} onClose={onClose}>
      <div className="flex flex-col w-full items-start p-5 gap-y-4">
        {selectedFilters.map((filterOption, filterIndex) => {
          const { label, value: filterType } = filterOption;
          const isMaintainerFilter = filterType === FilterTypes.MAINTAINER;
          const isTagFilter = filterType === FilterTypes.TAGS;
          const isHaveSearchingDropdown = isMaintainerFilter || isTagFilter;
          const valueOptions = getValueOptions(filterOption);
          return (
            <div
              className="flex items-start w-full h-[100px] gap-x-3"
              key={filterIndex}
            >
              <div className="h-full flex flex-col gap-y-4 items-center justify-center">
                <div
                  className={cn(
                    'mt-2 typo-para-small text-center py-[3px] w-[42px] min-w-[42px] rounded text-accent-pink-500 bg-accent-pink-50',
                    {
                      'bg-gray-200 text-gray-600': filterIndex !== 0
                    }
                  )}
                >
                  {t(filterIndex === 0 ? `if` : 'and')}
                </div>
                <Divider vertical={true} className="border-primary-500" />
              </div>
              <div className="flex flex-col w-full">
                <Dropdown
                  placeholder={t(`select-filter`)}
                  labelCustom={label}
                  options={remainingFilterOptions as DropdownOption[]}
                  value={filterType}
                  onChange={value => {
                    handleChangeOption(value as string, filterIndex);
                  }}
                  className="w-full truncate py-2"
                  contentClassName="w-[270px]"
                />
                <div className="flex items-center gap-3 mt-3 pl-3">
                  <p className="typo-para-medium text-gray-600">is</p>
                  <Dropdown
                    disabled={
                      (isTagFilter && isLoadingTags) ||
                      (isMaintainerFilter && isLoading) ||
                      !filterType
                    }
                    loading={
                      (isTagFilter && isLoadingTags) ||
                      (isMaintainerFilter && isLoading)
                    }
                    isListItem={isHaveSearchingDropdown}
                    multiselect={isTagFilter}
                    placeholder={t(`select-value`)}
                    value={filterOption?.filterValue as DropdownValue}
                    labelCustom={handleGetLabelFilterValue(filterOption)}
                    isSearchable={isHaveSearchingDropdown}
                    options={valueOptions as DropdownOption[]}
                    onChange={value =>
                      handleChangeFilterValue(
                        value as DropdownValue,
                        filterIndex
                      )
                    }
                    className="w-full truncate py-2"
                    contentClassName={cn('w-[235px]', {
                      'pt-0 w-[300px]': isHaveSearchingDropdown,
                      'hidden-scroll': valueOptions?.length > 15
                    })}
                  />

                  <Button
                    variant={'grey'}
                    className="px-0 w-fit"
                    disabled={selectedFilters.length <= 1}
                    onClick={() =>
                      setSelectedFilters(
                        selectedFilters.filter(
                          (_, index) => filterIndex !== index
                        )
                      )
                    }
                  >
                    <Icon icon={IconTrash} size={'sm'} />
                  </Button>
                </div>
              </div>
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
    </SlideModal>
  );
};

export default FilterFlagSlideModal;
