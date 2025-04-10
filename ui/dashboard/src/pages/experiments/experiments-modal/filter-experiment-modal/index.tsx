import { useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'i18n';
import { ExperimentStatus } from '@types';
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
  STATUS = 'status'
}

export const filterOptions: Option[] = [
  {
    value: FilterTypes.STATUS,
    label: 'Status'
  }
];

export const statusOptions: Option[] = [
  {
    value: 'WAITING',
    label: 'Waiting'
  },
  {
    value: 'RUNNING',
    label: 'Running'
  },
  {
    value: 'STOPPED',
    label: 'Stopped'
  },
  {
    value: 'FORCE_STOPPED',
    label: 'Force Stopped'
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
  const [selectedStatuses, setSelectedStatuses] = useState<ExperimentStatus[]>(
    []
  );

  const isDisabledSubmitBtn = useMemo(
    () => !selectedStatuses.length,
    [selectedStatuses]
  );

  const onConfirmHandler = () =>
    onSubmit({
      archived: undefined,
      isFilter: true,
      statuses: selectedStatuses
    });

  useEffect(() => {
    if (
      isNotEmpty(
        (filters?.isFilter || filters?.filterBySummary) &&
          filters?.statuses?.length
      )
    ) {
      setSelectedStatuses(
        typeof filters!.statuses === 'string'
          ? [filters!.statuses]
          : (filters!.statuses as ExperimentStatus[])
      );
    } else {
      setSelectedStatuses([]);
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
              label={filterOptions[0].label}
              variant="secondary"
              className="w-full"
            />
            <DropdownMenuContent className="w-[235px]" align="start">
              {filterOptions.map((item, index) => (
                <DropdownMenuItem
                  key={index}
                  value={item.value}
                  label={item.label}
                />
              ))}
            </DropdownMenuContent>
          </DropdownMenu>
          <p className="typo-para-medium text-gray-600">{`is`}</p>
          <DropdownMenu>
            <DropdownMenuTrigger
              placeholder={t(`select-value`)}
              label={
                selectedStatuses?.length
                  ? selectedStatuses
                      ?.join(', ')
                      ?.replace('_', ' ')
                      ?.toLowerCase()
                  : ''
              }
              variant="secondary"
              className="w-full capitalize"
            />
            <DropdownMenuContent className="w-[235px]" align="start">
              {statusOptions.map((item, index) => (
                <DropdownMenuItem
                  isSelected={selectedStatuses.includes(
                    item.value as ExperimentStatus
                  )}
                  isMultiselect
                  key={index}
                  value={item.value}
                  label={item.label}
                  onSelectOption={value => {
                    const isExisted = selectedStatuses?.find(
                      item => item === value
                    );
                    setSelectedStatuses(
                      isExisted
                        ? selectedStatuses.filter(item => item !== value)
                        : [...selectedStatuses, value as ExperimentStatus]
                    );
                  }}
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

export default FilterExperimentModal;
