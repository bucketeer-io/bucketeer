import { useEffect, useState } from 'react';
import { useTranslation } from 'i18n';
import { isNotEmpty } from 'utils/data-type';
import { ExperimentFilters } from 'pages/experiments/types';
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
  onSubmit: (v: Partial<ExperimentFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<ExperimentFilters>;
};

export interface Option {
  value: string;
  label: string;
}

export enum FilterTypes {
  ARCHIVED = 'archived',
  FINISHED = 'finished'
}

export const filterOptions: Option[] = [
  {
    value: FilterTypes.ARCHIVED,
    label: 'Archived'
  },
  {
    value: FilterTypes.FINISHED,
    label: 'Finished'
  }
];

export const archiveOptions: Option[] = [
  {
    value: 'yes',
    label: 'Yes'
  },
  {
    value: 'no',
    label: 'No'
  }
];

const FilterExperimentModal = ({
  onSubmit,
  isOpen,
  onClose,
  onClearFilters,
  filters
}: FilterProps) => {
  const { t } = useTranslation(['common']);
  const [selectedFilterType, setSelectedFilterType] = useState<Option>();
  const [selectedValue, setSelectedValue] = useState<Option>();

  const onConfirmHandler = () => {
    switch (selectedFilterType?.value) {
      case FilterTypes.ARCHIVED:
        if (selectedValue) {
          onSubmit({
            archived: selectedValue?.value === 'yes',
            statuses: undefined
          });
        }
        return;
      case FilterTypes.FINISHED:
        if (selectedValue) {
          onSubmit({
            archived: undefined,
            statuses: selectedValue?.value === 'yes' ? 'RUNNING' : undefined
          });
        }
        return;
    }
  };

  useEffect(() => {
    if (isNotEmpty(filters?.archived || filters?.statuses)) {
      setSelectedFilterType(
        filters?.archived ? filterOptions[0] : filterOptions[1]
      );
      setSelectedValue(
        archiveOptions[filters?.archived || filters?.statuses ? 0 : 1]
      );
    } else {
      setSelectedFilterType(undefined);
      setSelectedValue(undefined);
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
              label={selectedValue?.label}
              variant="secondary"
              disabled={!selectedFilterType}
              className="w-full"
            />
            <DropdownMenuContent className="w-[235px]" align="start">
              {archiveOptions.map((item, index) => (
                <DropdownMenuItem
                  key={index}
                  value={item.value}
                  label={item.label}
                  onSelectOption={() => setSelectedValue(item)}
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

export default FilterExperimentModal;
