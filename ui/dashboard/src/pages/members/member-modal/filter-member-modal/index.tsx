import { useCallback, useEffect, useMemo, useState } from 'react';
import { useQueryTags } from '@queries/tags';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { isNotEmpty } from 'utils/data-type';
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
import DialogModal from 'components/modal/dialog';

export type FilterProps = {
  onSubmit: (v: Partial<MembersFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<MembersFilters>;
};

export interface Option {
  value: string;
  label: string;
}

export enum FilterTypes {
  ENABLED = 'enabled',
  ROLE = 'role',
  TAGS = 'tags'
}

export const filterOptions: Option[] = [
  {
    value: FilterTypes.ENABLED,
    label: 'Enabled'
  },
  {
    value: FilterTypes.ROLE,
    label: 'Role'
  },
  {
    value: FilterTypes.TAGS,
    label: 'Tags'
  }
];

export const enabledOptions: Option[] = [
  {
    value: 'yes',
    label: 'Yes'
  },
  {
    value: 'no',
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
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const [selectedFilterType, setSelectedFilterType] = useState<Option>();
  const [filterValue, setFilterValue] = useState<string | string[]>();

  const isEnabledFilter = useMemo(
    () => selectedFilterType?.value === FilterTypes.ENABLED,
    [selectedFilterType]
  );
  const isTagFilter = useMemo(
    () => selectedFilterType?.value === FilterTypes.TAGS,
    [selectedFilterType]
  );

  const { data: tagCollection, isLoading: isLoadingTags } = useQueryTags({
    params: {
      cursor: String(0),
      organizationId: currentEnvironment?.organizationId,
      entityType: 'ACCOUNT'
    },
    enabled: isTagFilter
  });

  const tags = tagCollection?.tags || [];

  const dropdownOptions = useMemo(() => {
    switch (selectedFilterType?.value) {
      case FilterTypes.ENABLED:
        return enabledOptions;
      case FilterTypes.ROLE:
        return roleOptions;
      case FilterTypes.TAGS:
        return tags?.map(item => ({
          label: item.name,
          value: item.id
        }));
      default:
        return [];
    }
  }, [selectedFilterType, tags]);

  const handleSetFilterOnInit = useCallback(() => {
    if (filters) {
      const { organizationRole, disabled, tags } = filters || {};
      const isNotEmptyRole = isNotEmpty(organizationRole);
      const isNotTag = isNotEmpty(tags);
      const isNotEmptyDisabled = isNotEmpty(disabled);
      if (isNotEmptyDisabled) {
        setSelectedFilterType(filterOptions[0]);
        return setFilterValue(enabledOptions[disabled ? 1 : 0].value);
      }
      if (isNotEmptyRole) {
        setSelectedFilterType(filterOptions[1]);
        return setFilterValue(
          roleOptions.find(
            item => item.value === String(filters?.organizationRole)
          )?.value
        );
      }

      if (isNotTag && tags) {
        setFilterValue(Array.isArray(tags) ? tags : [tags]);
        return setSelectedFilterType(filterOptions[2]);
      }
    }

    setSelectedFilterType(undefined);
    setFilterValue(undefined);
  }, [filters]);

  const labelFilterValue = useMemo(() => {
    if (!isTagFilter)
      return (
        (isEnabledFilter ? enabledOptions : roleOptions).find(
          item => item.value === filterValue
        )?.label || ''
      );

    return (
      (Array.isArray(filterValue) &&
        tags.length &&
        filterValue
          .map(item => tags.find(tag => tag.id === item)?.name)
          ?.join(', ')) ||
      ''
    );
  }, [selectedFilterType, filterValue, isTagFilter, tags]);

  const handleChangeFilterValue = useCallback(
    (value: string) => {
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

  const onConfirmHandler = useCallback(() => {
    switch (selectedFilterType?.value) {
      case FilterTypes.ENABLED: {
        if (filterValue) {
          onSubmit({
            disabled: filterValue === 'no'
          });
        }
        return;
      }
      case FilterTypes.ROLE: {
        if (filterValue) {
          onSubmit({
            organizationRole: +filterValue
          });
        }
        return;
      }
      case FilterTypes.TAGS: {
        return onSubmit({
          tags: Array.isArray(filterValue) ? filterValue : []
        });
      }
    }
  }, [selectedFilterType, filterValue]);

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
              label={selectedFilterType?.label}
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
                    setSelectedFilterType(item);
                    setFilterValue(item.value === FilterTypes.TAGS ? [] : '');
                  }}
                />
              ))}
            </DropdownMenuContent>
          </DropdownMenu>
          <p className="typo-para-medium text-gray-600">{`is`}</p>
          <DropdownMenu>
            <DropdownMenuTrigger
              placeholder={t(`select-value`)}
              label={labelFilterValue}
              disabled={!selectedFilterType || isLoadingTags}
              variant="secondary"
              className="w-full max-w-[235px] truncate"
            />
            <DropdownMenuContent
              className={isTagFilter ? 'w-[300px]' : 'w-[235px]'}
              align="start"
            >
              {dropdownOptions.map((item, index) => (
                <DropdownMenuItem
                  isSelected={
                    isTagFilter &&
                    Array.isArray(filterValue) &&
                    filterValue.includes(item.value as string)
                  }
                  isMultiselect={isTagFilter}
                  key={index}
                  value={item.value}
                  label={item.label}
                  onSelectOption={() => handleChangeFilterValue(item.value)}
                />
              ))}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      <ButtonBar
        secondaryButton={
          <Button onClick={onConfirmHandler}>{t(`confirm`)}</Button>
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
