import { useEffect, useMemo, useState } from 'react';
import { i18n, useTranslation } from 'i18n';
import { isNotEmpty } from 'utils/data-type';
import { UserSegmentsFilters } from 'pages/user-segments/types';
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
  onSubmit: (v: Partial<UserSegmentsFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<UserSegmentsFilters>;
};

export interface Option {
  value: string;
  label: string;
}

export enum FilterTypes {
  STATUS = 'status'
}

export enum FilterValue {
  IN_USE = 'in-use',
  NOT_IN_USE = 'not-in-use'
}

const translation = i18n.t;

export const filterOptions: Option[] = [
  {
    value: FilterTypes.STATUS,
    label: translation('common:status')
  }
];

export const statusOptions: Option[] = [
  {
    value: FilterValue.IN_USE,
    label: translation('common:in-use')
  },
  {
    value: FilterValue.NOT_IN_USE,
    label: translation('common:not-in-use')
  }
];

const FilterUserSegmentModal = ({
  onSubmit,
  isOpen,
  onClose,
  onClearFilters,
  filters
}: FilterProps) => {
  const { t } = useTranslation(['common']);
  const [selectedFilterType, setSelectedFilterType] = useState<Option>();
  const [valueOption, setValueOption] = useState<Option>();

  const isDisabledSubmitBtn = useMemo(
    () => !selectedFilterType || !valueOption,
    [selectedFilterType, valueOption]
  );

  const onConfirmHandler = () => {
    switch (selectedFilterType?.value) {
      case FilterTypes.STATUS:
        if (valueOption?.value) {
          onSubmit({
            isInUseStatus: valueOption?.value === FilterValue.IN_USE
          });
        }
        return;
    }
  };

  useEffect(() => {
    if (isNotEmpty(filters?.isInUseStatus)) {
      setSelectedFilterType(filterOptions[0]);
      setValueOption(statusOptions[filters?.isInUseStatus ? 0 : 1]);
    } else {
      setSelectedFilterType(undefined);
      setValueOption(undefined);
    }
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
              {statusOptions.map((item, index) => (
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
          <Button disabled={isDisabledSubmitBtn} onClick={onConfirmHandler}>
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

export default FilterUserSegmentModal;
