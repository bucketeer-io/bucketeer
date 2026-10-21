import { useEffect, useMemo, useState } from 'react';
import { FilterOption } from 'hooks/use-options';
import { useTranslation } from 'i18n';
import { isEmpty } from 'utils/data-type';
import { cn } from 'utils/style';
import { IconPlus, IconTrash } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Dropdown, { DropdownValue } from 'components/dropdown';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';
import DropdownMenuWithSearch from 'elements/dropdown-with-search';
import { FilterFieldDef, FilterModalConfig } from './types';

export type FilterModalProps<F> = {
  config: FilterModalConfig<F>;
  isOpen: boolean;
  filters?: Partial<F>;
  onSubmit: (v: Partial<F>) => void;
  onClose: () => void;
  onClearFilters: () => void;
};

const isMulti = (kind: FilterFieldDef<unknown>['valueKind']) =>
  kind === 'multiselect' || kind === 'searchable';

function FilterModal<F>({
  config,
  isOpen,
  filters,
  onSubmit,
  onClose,
  onClearFilters
}: FilterModalProps<F>) {
  const { t } = useTranslation(['common']);
  const { fields, mode } = config;

  const filterTypeOptions: FilterOption[] = useMemo(
    () =>
      fields.map(field => ({
        value: field.type,
        label: t(field.labelKey),
        filterValue: field.emptyValue
      })),
    [fields, t]
  );

  const [selectedFilters, setSelectedFilters] = useState<FilterOption[]>([
    filterTypeOptions[0]
  ]);

  const findField = (type?: FilterOption['value']) =>
    fields.find(field => field.type === type);

  const remainingFilterOptions = useMemo(
    () =>
      filterTypeOptions.filter(
        option => !selectedFilters.find(item => item.value === option.value)
      ),
    [filterTypeOptions, selectedFilters]
  );

  const isDisabledAddButton = useMemo(
    () =>
      !remainingFilterOptions.length ||
      selectedFilters.length >= filterTypeOptions.length,
    [filterTypeOptions, selectedFilters, remainingFilterOptions]
  );

  const isDisabledSubmitButton = useMemo(
    () => !!selectedFilters.find(item => isEmpty(item.filterValue)),
    [selectedFilters]
  );

  const handleSetFilterOnInit = () => {
    if (!filters) return;
    if (config.shouldHydrate && !config.shouldHydrate(filters)) return;
    const next: FilterOption[] = [];
    fields.forEach(field => {
      const value = field.fromFilter(filters);
      if (!isEmpty(value)) {
        const option = filterTypeOptions.find(
          item => item.value === field.type
        );
        if (option) next.push({ ...option, filterValue: value });
      }
    });
    setSelectedFilters(next.length ? next : [filterTypeOptions[0]]);
  };

  useEffect(() => {
    handleSetFilterOnInit();
  }, [filters]);

  const handleChangeOption = (value: DropdownValue, filterIndex: number) => {
    if (selectedFilters[filterIndex].value === value) return;
    const option = filterTypeOptions.find(item => item.value === value);
    const field = findField(value);
    if (!option || !field) return;
    setSelectedFilters(prev => {
      const nextFilters = [...prev];
      nextFilters[filterIndex] = { ...option, filterValue: field.emptyValue };
      return nextFilters;
    });
  };

  const handleChangeFilterValue = (
    value: DropdownValue,
    filterIndex: number
  ) => {
    const current = selectedFilters[filterIndex];
    const field = findField(current.value);
    setSelectedFilters(prev => {
      const nextFilters = [...prev];
      if (field && isMulti(field.valueKind)) {
        const values = (current.filterValue as string[]) || [];
        if (Array.isArray(value) && isEmpty(value)) {
          nextFilters[filterIndex] = { ...current, filterValue: value };
          return nextFilters;
        }
        const exists = values.includes(value as string);
        const newValue = exists
          ? values.filter(item => item !== value)
          : [...values, value as string];
        nextFilters[filterIndex] = { ...current, filterValue: newValue };
        return nextFilters;
      }
      nextFilters[filterIndex] = { ...current, filterValue: value };
      return nextFilters;
    });
  };

  const onConfirmHandler = () => {
    const defaultFilters =
      config.defaultFilters ??
      fields.reduce<Partial<F>>(
        (acc, field) => ({ ...acc, ...emptyKeys(field.toFilter) }),
        {} as Partial<F>
      );
    const newFilters = {} as Partial<F>;
    selectedFilters.forEach(filter => {
      const field = findField(filter.value);
      if (field) {
        const mapped = field.toFilter(filter.filterValue);
        Object.assign(newFilters, mapped);
      }
    });
    const submitted = {
      ...defaultFilters,
      ...newFilters,
      ...config.submitExtra
    };

    onSubmit(submitted);
  };

  return (
    <DialogModal
      className="w-[750px]"
      title={t('filters')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full items-start p-5 gap-y-4">
        {selectedFilters.map((filterOption, filterIndex) => {
          const field = findField(filterOption.value);

          return (
            <div
              key={filterIndex}
              className="flex items-center w-full h-12 gap-x-4"
            >
              <div
                className={cn(
                  'typo-para-small text-center py-[3px] rounded text-accent-pink-500 bg-accent-pink-50',
                  mode === 'single' ? 'px-4' : 'w-[42px] min-w-[42px]',
                  {
                    'bg-gray-200 text-gray-600': filterIndex !== 0
                  }
                )}
              >
                {t(filterIndex === 0 ? 'if' : 'and')}
              </div>
              <Divider vertical={true} className="border-primary-500" />
              <Dropdown
                placeholder={t('select-filter')}
                labelCustom={filterOption.label}
                className="w-full truncate"
                contentClassName="w-[270px]"
                options={(mode === 'single'
                  ? filterTypeOptions
                  : remainingFilterOptions
                ).map(item => ({
                  ...item,
                  value: item.value || '',
                  label: item.label
                }))}
                value={filterOption.value}
                onChange={value =>
                  handleChangeOption(value as DropdownValue, filterIndex)
                }
              />

              <p className="typo-para-medium text-gray-600">is</p>

              <FilterValueInput
                key={field?.type ?? 'empty'}
                field={field}
                value={filterOption.filterValue}
                onChange={value => handleChangeFilterValue(value, filterIndex)}
              />

              {mode === 'multi' && (
                <Button
                  variant="grey"
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
                  <Icon icon={IconTrash} size="sm" />
                </Button>
              )}
            </div>
          );
        })}

        {mode === 'multi' && (
          <Button
            disabled={isDisabledAddButton}
            variant="text"
            className="h-6 px-0"
            onClick={() =>
              setSelectedFilters([
                ...selectedFilters,
                { label: '', value: undefined, filterValue: '' }
              ])
            }
          >
            <Icon icon={IconPlus} />
            {t('add-filter')}
          </Button>
        )}
      </div>

      <ButtonBar
        secondaryButton={
          <Button disabled={isDisabledSubmitButton} onClick={onConfirmHandler}>
            {t('confirm')}
          </Button>
        }
        primaryButton={
          <Button onClick={onClearFilters} variant="secondary">
            {t('clear')}
          </Button>
        }
      />
    </DialogModal>
  );
}

// Keep query data and its callbacks in the component that renders the value.
// A field-type key remounts this component when its hook implementation changes.
function FilterValueInput<F>({
  field,
  value,
  onChange
}: {
  field?: FilterFieldDef<F>;
  value: FilterOption['filterValue'];
  onChange: (value: DropdownValue) => void;
}) {
  const { t } = useTranslation(['common']);
  const data = field?.useData
    ? field.useData({ enabled: true, value })
    : { options: [] };
  const kind = field?.valueKind;
  const multiselect = isMulti(kind ?? 'enum');
  const isSearchable = kind === 'searchable';
  const label = data.getLabel
    ? data.getLabel(value)
    : multiselect
      ? ((value as string[]) || [])
          .map(v => data.options.find(option => option.value === v)?.label)
          .filter(Boolean)
          .join(', ')
      : data.options.find(option => option.value === value)?.label || '';

  if (kind === 'searchable-paginated') {
    return (
      <DropdownMenuWithSearch
        disabled={data.isLoading || !field}
        isLoading={data.isLoading}
        placeholder={t('select-value')}
        itemSelected={value as string}
        label={label}
        options={data.options}
        isHasMore={data.hasMore}
        isLoadingMore={data.isLoadingMore}
        isSearching={data.isSearching}
        onHasMoreOptions={data.loadMore}
        onSearchChange={data.onSearchChange}
        onSelectOption={value => onChange(value as DropdownValue)}
        triggerClassName="w-full truncate"
        contentClassName="w-[300px]"
      />
    );
  }

  return (
    <Dropdown
      disabled={data.isLoading || !field}
      loading={data.isLoading}
      isSearchable={isSearchable}
      multiselect={multiselect}
      placeholder={t('select-value')}
      labelCustom={label}
      className="w-full truncate"
      contentClassName={cn('w-[235px]', {
        'pt-0 w-[300px]': isSearchable,
        'hidden-scroll': data.options.length > 15
      })}
      value={multiselect ? (value as string[]) : (value as DropdownValue)}
      options={data.options}
      onChange={value => onChange(value as DropdownValue)}
    />
  );
}

/**
 * Derive the keys a field writes by invoking its `toFilter` with `undefined`,
 * so de-selected filters reset to `undefined` on submit.
 */
function emptyKeys<F>(
  toFilter: FilterFieldDef<F>['toFilter']
): Record<string, undefined> {
  const result = toFilter(undefined) as Record<string, unknown>;
  return Object.keys(result).reduce(
    (acc, key) => ({ ...acc, [key]: undefined }),
    {}
  );
}

export default FilterModal;
