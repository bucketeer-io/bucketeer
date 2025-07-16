import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { getEditorEnvironments, useAuth } from 'auth';
import useOptions, { FilterOption, FilterTypes } from 'hooks/use-options';
import { useTranslation } from 'i18n';
import debounce from 'lodash/debounce';
import isNil from 'lodash/isNil';
import { isEmpty } from 'utils/data-type';
import { checkEnvironmentEmptyId, onFormatEnvironments } from 'utils/function';
import { cn } from 'utils/style';
import { IconPlus, IconTrash } from '@icons';
import { NotificationFilters } from 'pages/notifications/types';
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
  onSubmit: (v: Partial<NotificationFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<NotificationFilters>;
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
  const { editorEnvironments } = getEditorEnvironments(consoleAccount!);

  const [selectedFilters, setSelectedFilters] = useState<FilterOption[]>([
    environmentEnabledFilterOptions[0]
  ]);
  const [searchValue, setSearchValue] = useState('');
  const [debounceValue, setDebounceValue] = useState('');
  const inputSearchRef = useRef<HTMLInputElement>(null);

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

  const debouncedSearch = useCallback(
    debounce(value => {
      setSearchValue(value);
    }, 500),
    []
  );

  const { emptyEnvironmentId, formattedEnvironments } =
    onFormatEnvironments(editorEnvironments);

  const environmentOptions = useMemo(
    () =>
      formattedEnvironments.map(item => ({
        label: item.name,
        value: item.id
      })),
    [formattedEnvironments]
  );

  const handleFocusSearchInput = useCallback(() => {
    let timerId: NodeJS.Timeout | null = null;
    if (timerId) clearTimeout(timerId);
    timerId = setTimeout(() => inputSearchRef?.current?.focus(), 50);
  }, []);

  const getValueOptions = useCallback(
    (filterOption: FilterOption) => {
      if (!filterOption.value) return [];
      const isEnvironmentFilter =
        filterOption.value === FilterTypes.ENVIRONMENT_IDs;

      if (isEnvironmentFilter) {
        return environmentOptions?.filter(item =>
          searchValue
            ? item.value.toLowerCase().includes(searchValue.toLowerCase())
            : item
        );
      }

      return booleanOptions;
    },
    [booleanOptions, searchValue, environmentOptions]
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
      let newFilterValue: string | number | string[] = value;
      if (isEnvironmentFilter) {
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
                          item.value === FilterTypes.ENVIRONMENT_IDs ? [] : '';
                        selectedFilters[filterIndex] = {
                          ...item,
                          filterValue
                        };
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
                  disabled={!filterType}
                  placeholder={t(`select-value`)}
                  label={handleGetLabelFilterValue(filterOption)}
                  variant="secondary"
                  className="w-full truncate"
                />
                <DropdownMenuContent
                  className={cn('w-[235px]', {
                    'pt-0 w-[300px]': isEnvironmentFilter,
                    'hidden-scroll': valueOptions?.length > 15
                  })}
                  align="start"
                >
                  {isEnvironmentFilter && (
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
                      isMultiselect={isEnvironmentFilter}
                      selectedOptions={
                        isEnvironmentFilter &&
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
