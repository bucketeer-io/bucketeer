import { useCallback, useEffect, useMemo, useState } from 'react';
import { useQueryFeatures } from '@queries/features';
import { useTranslation } from 'i18n';
import { debounce } from 'lodash';
import { LIST_PAGE_SIZE } from 'constants/app';
import { getCurrentEnvIdStorage } from 'storage/environment';
import DropdownMenuWithSearch from 'elements/dropdown-with-search';

interface FlagSelectorProps {
  selectedFlagId: string | undefined;
  onSelectFlag: (flagId: string | undefined) => void;
}

const FlagSelector = ({ selectedFlagId, onSelectFlag }: FlagSelectorProps) => {
  const { t } = useTranslation(['ai-chat']);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedFlagLabel, setSelectedFlagLabel] = useState('');
  const environmentId = getCurrentEnvIdStorage() || '';

  const { data: flagCollection, isLoading } = useQueryFeatures({
    params: {
      cursor: String(0),
      pageSize: LIST_PAGE_SIZE,
      searchKeyword: searchQuery,
      environmentId,
      archived: false
    }
  });

  const flagOptions = useMemo(
    () =>
      (flagCollection?.features || []).map(f => ({
        label: f.name,
        value: f.id,
        enabled: f.enabled
      })),
    [flagCollection]
  );

  // Cache selected flag label so it persists across searches
  useEffect(() => {
    if (selectedFlagId) {
      const found = flagOptions.find(f => f.value === selectedFlagId);
      if (found) {
        setSelectedFlagLabel(found.label as string);
      }
    } else {
      setSelectedFlagLabel('');
    }
  }, [flagOptions, selectedFlagId]);

  // Debounce search to avoid excessive API calls (matches debugger pattern)
  const debouncedSearch = useMemo(
    () => debounce((value: string) => setSearchQuery(value), 300),
    []
  );

  useEffect(() => {
    return () => {
      debouncedSearch.cancel();
    };
  }, [debouncedSearch]);

  const handleSearchChange = useCallback(
    (value: string) => {
      if (value === '') {
        debouncedSearch.cancel();
        setSearchQuery('');
      } else {
        debouncedSearch(value);
      }
    },
    [debouncedSearch]
  );

  return (
    <DropdownMenuWithSearch
      label={selectedFlagLabel}
      isExpand
      isLoading={isLoading}
      placeholder={t('ai-chat:flag-selector.placeholder')}
      inputPlaceholder={t('ai-chat:flag-selector.search-placeholder')}
      options={flagOptions}
      itemSelected={selectedFlagId}
      showClear={!!selectedFlagId}
      onSelectOption={value => onSelectFlag(value as string)}
      onSearchChange={handleSearchChange}
      onClear={() => onSelectFlag(undefined)}
    />
  );
};

export default FlagSelector;
