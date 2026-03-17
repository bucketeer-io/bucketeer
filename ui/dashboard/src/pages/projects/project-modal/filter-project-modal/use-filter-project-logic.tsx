import { useEffect, useMemo, useState } from 'react';
import useOptions, { FilterOption, FilterTypes } from 'hooks/use-options';
import { isNotEmpty } from 'utils/data-type';
import { ProjectFilters } from 'pages/projects/types';

export type FilterProps = {
  onSubmit: (v: Partial<ProjectFilters>) => void;

  filters?: Partial<ProjectFilters>;
};

const useFilterProjectLogic = (
  onSubmit: (v: Partial<ProjectFilters>) => void,

  filters?: Partial<ProjectFilters>
) => {
  const { enabledOptions, filterEnabledOptions } = useOptions();
  const [selectedFilterType, setSelectedFilterType] = useState<FilterOption>();
  const [valueOption, setValueOption] = useState<FilterOption>();

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

  const isDisabledSubmitBtn = useMemo(
    () => !selectedFilterType || !valueOption,
    [selectedFilterType, valueOption]
  );

  const handleChangeOption = (value: string) => {
    const selected = filterEnabledOptions.find(item => item.value === value);
    setSelectedFilterType(selected);
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

  return {
    enabledOptions,
    filterEnabledOptions,
    selectedFilterType,
    valueOption,
    isDisabledSubmitBtn,
    setSelectedFilterType,
    setValueOption,
    onConfirmHandler,
    handleChangeOption
  };
};

export default useFilterProjectLogic;
