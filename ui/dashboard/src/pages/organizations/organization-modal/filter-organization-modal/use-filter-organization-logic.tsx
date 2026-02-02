import { useEffect, useMemo, useState } from 'react';
import useOptions, { FilterOption, FilterTypes } from 'hooks/use-options';
import { isNotEmpty } from 'utils/data-type';
import { OrganizationFilters } from 'pages/organizations/types';

const useFilterOrganizationLogic = (
  onSubmit: (v: Partial<OrganizationFilters>) => void,
  filters?: Partial<OrganizationFilters>
) => {
  const { enabledOptions, filterEnabledOptions } = useOptions();
  const [selectedFilterType, setSelectedFilterType] = useState<FilterOption>();
  const [selectedValue, setSelectedValue] = useState<FilterOption>();
  const isDisabledSubmitBtn = useMemo(
    () => !selectedFilterType || !selectedValue,
    [selectedFilterType, selectedValue]
  );

  const onConfirmHandler = () => {
    switch (selectedFilterType?.value) {
      case FilterTypes.ENABLED:
        if (selectedValue) {
          onSubmit({
            disabled: selectedValue?.value === 'no'
          });
        }
        return;
    }
  };

  useEffect(() => {
    if (isNotEmpty(filters?.disabled)) {
      setSelectedFilterType(filterEnabledOptions[0]);
      setSelectedValue(enabledOptions[filters?.disabled ? 1 : 0]);
    } else {
      setSelectedFilterType(undefined);
      setSelectedValue(undefined);
    }
  }, [filters]);

  return {
    isDisabledSubmitBtn,
    onConfirmHandler,
    selectedFilterType,
    setSelectedFilterType,
    selectedValue,
    setSelectedValue,
    filterEnabledOptions,
    enabledOptions
  };
};

export default useFilterOrganizationLogic;
