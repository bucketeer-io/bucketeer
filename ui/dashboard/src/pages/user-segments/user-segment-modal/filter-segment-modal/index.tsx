import { useEffect, useMemo, useState } from 'react';
import useOptions, { FilterOption, FilterTypes } from 'hooks/use-options';
import { useTranslation } from 'i18n';
import { isNotEmpty } from 'utils/data-type';
import { UserSegmentsFilters } from 'pages/user-segments/types';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Dropdown from 'components/dropdown';
import DialogModal from 'components/modal/dialog';

export type FilterProps = {
  onSubmit: (v: Partial<UserSegmentsFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<UserSegmentsFilters>;
};

const FilterUserSegmentModal = ({
  onSubmit,
  isOpen,
  onClose,
  onClearFilters,
  filters
}: FilterProps) => {
  const { t } = useTranslation(['common']);
  const { filterStatusOptions, segmentStatusOptions } = useOptions();
  const [selectedFilterType, setSelectedFilterType] = useState<FilterOption>();
  const [valueOption, setValueOption] = useState<FilterOption>();

  const isDisabledSubmitBtn = useMemo(
    () => !selectedFilterType || !valueOption,
    [selectedFilterType, valueOption]
  );

  const onConfirmHandler = () => {
    switch (selectedFilterType?.value) {
      case FilterTypes.STATUS:
        if (valueOption?.value) {
          onSubmit({
            isInUseStatus: valueOption?.value === FilterTypes.IN_USE
          });
        }
        return;
    }
  };

  useEffect(() => {
    if (isNotEmpty(filters?.isInUseStatus)) {
      setSelectedFilterType(filterStatusOptions[0]);
      setValueOption(segmentStatusOptions[filters?.isInUseStatus ? 0 : 1]);
    } else {
      setSelectedFilterType(undefined);
      setValueOption(undefined);
    }
  }, [filters]);

  return (
    <DialogModal
      className="w-[750px]"
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
          <Dropdown
            value={selectedFilterType?.value}
            onChange={value => {
              const selected = filterStatusOptions.find(
                item => item.value === value
              );
              setSelectedFilterType(selected);
            }}
            placeholder={t(`select-filter`)}
            options={filterStatusOptions.map(item => ({
              ...item,
              label: item.label,
              value: item.value || ''
            }))}
            className="w-full"
            contentClassName="w-[235px]"
          />

          <p className="typo-para-medium text-gray-600">is</p>
          <Dropdown
            placeholder={t(`select-value`)}
            disabled={!selectedFilterType}
            options={segmentStatusOptions.map(item => ({
              ...item,
              label: item.label,
              value: item.value || ''
            }))}
            value={valueOption?.value}
            onChange={value => {
              const selected = segmentStatusOptions.find(
                item => item.value === value
              );
              setValueOption(selected);
            }}
            className="w-full"
            contentClassName="w-[235px]"
          />
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
