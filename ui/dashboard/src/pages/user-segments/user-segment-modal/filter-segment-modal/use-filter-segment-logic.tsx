import { useEffect, useMemo, useState } from 'react';
import useOptions, { FilterOption, FilterTypes } from 'hooks/use-options';
import { isNotEmpty } from 'utils/data-type';
import { UserSegmentsFilters } from 'pages/user-segments/types';

const useFilterSegmentLogic = (
  filters?: Partial<UserSegmentsFilters>,
  onSubmit?: (v: Partial<UserSegmentsFilters>) => void
) => {
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
          onSubmit?.({
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

  return {
    selectedFilterType,
    valueOption,
    filterStatusOptions,
    segmentStatusOptions,
    isDisabledSubmitBtn,
    setSelectedFilterType,
    setValueOption,
    onConfirmHandler,
    confirm
  };
};

export default useFilterSegmentLogic;
