import { useEffect, useState } from 'react';
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
  ROLE = 'role'
}

export const filterOptions: Option[] = [
  {
    value: FilterTypes.ENABLED,
    label: 'Enabled'
  },
  {
    value: FilterTypes.ROLE,
    label: 'Role'
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
  const [selectedFilterType, setSelectedFilterType] = useState<Option>();
  const [valueOption, setValueOption] = useState<Option>();

  const onConfirmHandler = () => {
    switch (selectedFilterType?.value) {
      case FilterTypes.ENABLED:
        if (valueOption?.value) {
          onSubmit({
            disabled: valueOption?.value === 'no'
          });
        }
        return;
      case FilterTypes.ROLE:
        if (valueOption?.value) {
          onSubmit({
            organizationRole: +valueOption?.value
          });
        }
        return;
    }
  };

  useEffect(() => {
    if (isNotEmpty(filters?.disabled)) {
      setSelectedFilterType(filterOptions[0]);
      return setValueOption(enabledOptions[filters?.disabled ? 1 : 0]);
    }
    if (isNotEmpty(filters?.organizationRole)) {
      setSelectedFilterType(filterOptions[1]);
      return setValueOption(
        roleOptions.find(
          item => item.value === String(filters?.organizationRole)
        )
      );
    }
    setSelectedFilterType(undefined);
    setValueOption(undefined);
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
                  onSelectOption={() => setSelectedFilterType(item)}
                />
              ))}
            </DropdownMenuContent>
          </DropdownMenu>
          <p className="typo-para-medium text-gray-600">{`is`}</p>
          <DropdownMenu>
            <DropdownMenuTrigger
              placeholder={t(`select-value`)}
              label={valueOption?.label}
              disabled={!selectedFilterType}
              variant="secondary"
              className="w-full"
            />
            <DropdownMenuContent className="w-[235px]" align="start">
              {(selectedFilterType?.label === FilterTypes.ENABLED
                ? enabledOptions
                : roleOptions
              ).map((item, index) => (
                <DropdownMenuItem
                  key={index}
                  value={item.value}
                  label={item.label}
                  onSelectOption={() => setValueOption(item)}
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
