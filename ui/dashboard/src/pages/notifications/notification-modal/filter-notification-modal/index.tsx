import { useEffect, useMemo, useState } from 'react';
import useOptions, { FilterOption, FilterTypes } from 'hooks/use-options';
import { useTranslation } from 'i18n';
import { isNotEmpty } from 'utils/data-type';
import { APIKeysFilters } from 'pages/api-keys/types';
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
  onSubmit: (v: Partial<APIKeysFilters>) => void;
  isOpen: boolean;
  onClose: () => void;
  onClearFilters: () => void;
  filters?: Partial<APIKeysFilters>;
};

const FilterAPIKeyModal = ({
  onSubmit,
  isOpen,
  onClose,
  onClearFilters,
  filters
}: FilterProps) => {
  const { t } = useTranslation(['common']);
  const { enabledOptions, filterEnabledOptions } = useOptions();
  const [selectedFilterType, setSelectedFilterType] = useState<FilterOption>();
  const [valueOption, setValueOption] = useState<FilterOption>();

  const isDisabledSubmitBtn = useMemo(
    () => !selectedFilterType || !valueOption,
    [selectedFilterType, valueOption]
  );

  const onConfirmHandler = () => {
    switch (selectedFilterType?.value) {
      case FilterTypes.ENABLED:
        if (valueOption?.value) {
          onSubmit({
            disabled: valueOption?.value === 'no'
          });
        }
        return;
    }
  };

  useEffect(() => {
    if (isNotEmpty(filters?.disabled)) {
      setSelectedFilterType(filterEnabledOptions[0]);
      setValueOption(enabledOptions[filters?.disabled ? 1 : 0]);
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
              {filterEnabledOptions.map((item, index) => (
                <DropdownMenuItem
                  key={index}
                  value={item.value as string}
                  label={item.label}
                  onSelectOption={() => setSelectedFilterType(item)}
                />
              ))}
            </DropdownMenuContent>
          </DropdownMenu>
          <p className="typo-para-medium text-gray-600">Is</p>
          <DropdownMenu>
            <DropdownMenuTrigger
              placeholder={t(`select-value`)}
              label={valueOption?.label}
              disabled={!selectedFilterType}
              variant="secondary"
              className="w-full"
            />
            <DropdownMenuContent className="w-[235px]" align="start">
              {enabledOptions.map((item, index) => (
                <DropdownMenuItem
                  key={index}
                  value={item.value as string}
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

export default FilterAPIKeyModal;
