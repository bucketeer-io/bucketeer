import { useCallback, useEffect, useMemo, useState } from 'react';
import { useQueryEnvironments } from '@queries/environments';
import { getCurrentEnvironment, useAuth } from 'auth';
import useOptions, { FilterOption, FilterTypes } from 'hooks/use-options';
import { useTranslation } from 'i18n';
import isNil from 'lodash/isNil';
import { isEmpty } from 'utils/data-type';
import { checkEnvironmentEmptyId, onFormatEnvironments } from 'utils/function';
import { cn } from 'utils/style';
import { IconPlus, IconTrash } from '@icons';
import { APIKeysFilters } from 'pages/api-keys/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Dropdown, { DropdownOption, DropdownValue } from 'components/dropdown';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';

export type FilterProps = {
  onSubmit: (v: Partial<APIKeysFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<APIKeysFilters>;
};

export interface Option {
  value: string;
  label: string;
}

const FilterAPIKeyModal = ({
  onSubmit,
  isOpen,
  onClose,
  onClearFilters,
  filters
}: FilterProps) => {
  const { t } = useTranslation(['common']);
  const { environmentEnabledFilterOptions, booleanOptions } = useOptions();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
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

  const { data: environmentCollection, isLoading: isLoadingEnvironments } =
    useQueryEnvironments({
      params: {
        cursor: '0',
        organizationId: currentEnvironment.organizationId
      },
      enabled:
        !!currentEnvironment.organizationId &&
        !!selectedFilters.find(
          item => item.value === FilterTypes.ENVIRONMENT_IDs
        )
    });

  const environments = useMemo(
    () => environmentCollection?.environments || [],
    [environmentCollection]
  );

  const { emptyEnvironmentId, formattedEnvironments } =
    onFormatEnvironments(environments);

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

      if (isEnvironmentFilter) {
        return environmentOptions;
      }

      return booleanOptions;
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
      let newFilterValue: string | number | string[] | number[] = value;
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

  const handleChangeOption = (val: DropdownValue, filterIndex: number) => {
    const selectionOption = remainingFilterOptions.find(
      item => item.value === val
    );
    if (!selectionOption) return;
    const filterValue =
      selectionOption.value === FilterTypes.ENVIRONMENT_IDs ? [] : '';
    selectedFilters[filterIndex] = { ...selectionOption, filterValue };
    setSelectedFilters([...selectedFilters]);
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
          const isEnvironmentFilter =
            filterType === FilterTypes.ENVIRONMENT_IDs;
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
                className="w-full truncate"
                options={remainingFilterOptions.map(item => ({
                  ...item,
                  value: item.value || '',
                  label: item.label
                }))}
                value={filterType}
                onChange={val =>
                  handleChangeOption(val as DropdownValue, filterIndex)
                }
                contentClassName="w-[270px]"
              />

              <p className="typo-para-medium text-gray-600">is</p>

              <Dropdown
                isSearchable={isEnvironmentFilter}
                disabled={
                  (isEnvironmentFilter && isLoadingEnvironments) || !filterType
                }
                loading={isEnvironmentFilter && isLoadingEnvironments}
                placeholder={t(`select-value`)}
                labelCustom={handleGetLabelFilterValue(filterOption)}
                className="w-full truncate"
                options={valueOptions as DropdownOption[]}
                multiselect={isEnvironmentFilter}
                value={
                  isEnvironmentFilter
                    ? (filterOption.filterValue as string[])
                    : (filterOption.filterValue as string)
                }
                onChange={val => {
                  handleChangeFilterValue(val as string | number, filterIndex);
                }}
                contentClassName={cn('w-[235px]', {
                  'pt-0 w-[300px]': isEnvironmentFilter,
                  'hidden-scroll': valueOptions?.length > 15
                })}
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

export default FilterAPIKeyModal;
