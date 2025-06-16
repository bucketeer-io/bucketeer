import { useCallback, useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'i18n';
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
  DropdownMenuTrigger
} from 'components/dropdown';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';

export type FilterProps = {
  onSubmit: (v: Partial<MembersFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<MembersFilters>;
};

export interface Option {
  value: string | number;
  label: string;
  filterValue?: string | number | boolean;
}

export enum FilterTypes {
  ENABLED = 'enabled',
  ROLE = 'organizationRole'
}

export const filterOptions: Option[] = [
  {
    value: FilterTypes.ENABLED,
    label: 'Enabled',
    filterValue: ''
  },
  {
    value: FilterTypes.ROLE,
    label: 'Role',
    filterValue: ''
  }
];

export const enabledOptions: Option[] = [
  {
    value: 1,
    label: 'Yes'
  },
  {
    value: 0,
    label: 'No'
  }
];

export const roleOptions: Option[] = [
  {
    value: '1',
    label: 'Member'
  },
  {
    value: '2',
    label: 'Admin'
  },
  {
    value: '3',
    label: 'Owner'
  }
];

const FilterMemberModal = ({
  onSubmit,
  isOpen,
  onClose,
  onClearFilters,
  filters
}: FilterProps) => {
  const { t } = useTranslation(['common']);
  const [selectedFilters, setSelectedFilters] = useState<Option[]>([
    filterOptions[0]
  ]);

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

  const isDisabledSubmitButton = useMemo(() => {
    return !!selectedFilters.find(item => isEmpty(item.filterValue));
  }, [[...selectedFilters]]);

  const handleGetLabelFilterValue = useCallback((filterOption?: Option) => {
    if (filterOption) {
      const { value: filterType, filterValue } = filterOption;
      const currentOption = (
        filterType === FilterTypes.ENABLED ? enabledOptions : roleOptions
      ).find(item => item.value === filterValue);
      if (currentOption) return currentOption.label;
    }
    return '';
  }, []);

  const handleSetFilterOnInit = useCallback(() => {
    if (filters) {
      const { disabled, organizationRole } = filters || {};
      const filterTypeArr: Option[] = [];
      const addFilterOption = (index: number, value: Option['filterValue']) => {
        if (!isEmpty(value)) {
          const option = filterOptions[index];

          filterTypeArr.push({
            ...option,
            filterValue:
              option.value! === FilterTypes.ROLE
                ? value?.toString()
                : value
                  ? 0
                  : 1
          });
        }
      };
      addFilterOption(0, disabled);
      addFilterOption(1, organizationRole);
      setSelectedFilters(
        filterTypeArr.length ? filterTypeArr : [filterOptions[0]]
      );
    }
  }, [filters]);

  const onConfirmHandler = () => {
    const defaultFilters = {
      disabled: undefined,
      organizationRole: undefined
    };

    const newFilters = {};

    selectedFilters.forEach(filter => {
      Object.assign(newFilters, {
        [filter.value! === FilterTypes.ENABLED ? 'disabled' : filter.value!]:
          filter.value === FilterTypes.ENABLED
            ? !filter.filterValue
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
      className="w-[665px]"
      title={t('filters')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full items-start p-5 gap-y-4">
        {selectedFilters.map((filterOption, filterIndex) => {
          const { label, value: filterType } = filterOption;

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
                      value={item.value}
                      label={item.label}
                      onSelectOption={() => {
                        selectedFilters[filterIndex] = item;
                        setSelectedFilters([...selectedFilters]);
                      }}
                    />
                  ))}
                </DropdownMenuContent>
              </DropdownMenu>
              <p className="typo-para-medium text-gray-600">{`is`}</p>
              <DropdownMenu>
                <DropdownMenuTrigger
                  placeholder={t(`select-value`)}
                  label={handleGetLabelFilterValue(filterOption)}
                  disabled={!filterType}
                  variant="secondary"
                  className="w-full"
                />
                <DropdownMenuContent className="w-[235px]" align="start">
                  {(filterType === FilterTypes.ENABLED
                    ? enabledOptions
                    : roleOptions
                  ).map((item, index) => (
                    <DropdownMenuItem
                      key={index}
                      value={item.value}
                      label={item.label}
                      onSelectOption={() => {
                        selectedFilters[filterIndex] = {
                          ...selectedFilters[filterIndex],
                          filterValue: item.value
                        };
                        setSelectedFilters([...selectedFilters]);
                      }}
                    />
                  ))}
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
