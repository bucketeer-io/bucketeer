import { useMemo, useState, useCallback, useEffect } from 'react';
import { useFormContext } from 'react-hook-form';
import { useQueryFeatures } from '@queries/features';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { debounce } from 'lodash';
import { Feature } from '@types';
import { IconPlus, IconTrash } from '@icons';
import Button from 'components/button';
import Form from 'components/form';
import Icon from 'components/icon';
import DropdownMenuWithSearch from 'elements/dropdown-with-search';
import FeatureFlagStatus from 'elements/feature-flag-status';
import { AddDebuggerFormType } from './form-schema';

const DebuggerFlags = ({
  feature,
  isOnTargeting
}: {
  feature?: Feature;
  isOnTargeting?: boolean;
}) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { t } = useTranslation(['common', 'form']);
  const { control, watch, setValue } = useFormContext<AddDebuggerFormType>();

  const [searchQuery, setSearchQuery] = useState('');
  const [selectedFlagsCache, setSelectedFlagsCache] = useState<
    Map<string, Feature>
  >(new Map());

  const { data: flagCollection, isLoading } = useQueryFeatures({
    params: {
      cursor: String(0),
      // Always use pageSize: 0 to fetch all flags
      // This ensures users can browse all flags and see all search results
      // Virtual scrolling (maxOptions: 15) handles rendering performance
      pageSize: 0,
      searchKeyword: searchQuery,
      environmentId: currentEnvironment.id,
      archived: false
    }
  });

  const flags = flagCollection?.features || [];

  // Update cache with newly fetched flags
  useEffect(() => {
    if (flags.length > 0) {
      setSelectedFlagsCache(prev => {
        const newCache = new Map(prev);
        flags.forEach(flag => {
          newCache.set(flag.id, flag);
        });
        return newCache;
      });
    }
  }, [flags]);

  const flagsSelected: string[] = [...watch('flags')];

  // Merge cached selected flags with current search results
  const allAvailableFlags = useMemo(() => {
    const flagsMap = new Map(flags.map(flag => [flag.id, flag]));

    // Always add selected flags from cache for label display purposes
    flagsSelected.forEach(selectedId => {
      if (selectedId && !flagsMap.has(selectedId)) {
        const cachedFlag = selectedFlagsCache.get(selectedId);
        if (cachedFlag) {
          flagsMap.set(selectedId, cachedFlag);
        }
      }
    });

    return Array.from(flagsMap.values());
  }, [flags, flagsSelected, selectedFlagsCache]);

  const flagOptions = useMemo(
    () =>
      allAvailableFlags.map(item => ({
        label: item.name,
        value: item.id,
        enabled: item.enabled,
        disabled: flagsSelected.includes(item.id)
      })),
    [allAvailableFlags, flagsSelected]
  );

  const flagsRemaining = useMemo(() => {
    // Get IDs of flags actually returned by the API (not from cache)
    const apiReturnedFlagIds = new Set(flags.map(flag => flag.id));

    return flagOptions.filter(item => {
      // Filter out the current feature if in targeting mode
      if (item.value === feature?.id) return false;

      // If actively searching and flag is selected but NOT in API results,
      // hide it from dropdown to show accurate search results
      if (searchQuery && item.disabled && !apiReturnedFlagIds.has(item.value)) {
        return false;
      }

      return true;
    });
  }, [flagOptions, flags, feature, searchQuery]);

  const debouncedSearch = useMemo(
    () => debounce((value: string) => setSearchQuery(value), 300),
    []
  );

  const handleSearchChange = useCallback(
    (value: string) => {
      if (value === '') {
        // Clear search immediately without debounce
        debouncedSearch.cancel();
        setSearchQuery('');
      } else {
        // Debounce for actual searches
        debouncedSearch(value);
      }
    },
    [debouncedSearch]
  );

  // Cleanup debounced function on unmount
  useEffect(() => {
    return () => {
      debouncedSearch.cancel();
    };
  }, [debouncedSearch]);

  const isDisabledAddBtn = useMemo(() => {
    const totalFlagCount = flagCollection?.totalCount
      ? parseInt(flagCollection.totalCount)
      : 0;

    // Disable if no remaining flags in current view OR all flags are selected
    return !flagsRemaining.length || flagsSelected?.length >= totalFlagCount;
  }, [flagsRemaining, flagsSelected, flagCollection]);

  return (
    <>
      <div className="flex flex-col w-full gap-y-6">
        {flagsSelected.map((_, index) => (
          <Form.Field
            name={`flags.${index}`}
            key={index}
            control={control}
            render={({ field }) => (
              <Form.Item className="py-0">
                <Form.Label required>{t('flag')}</Form.Label>
                <Form.Control>
                  <div className="flex items-center w-full gap-x-4">
                    <DropdownMenuWithSearch
                      label={
                        flagOptions.find(flag => flag.value === field.value)
                          ?.label || ''
                      }
                      isExpand
                      disabled={isOnTargeting}
                      isLoading={isLoading}
                      placeholder={t('form:experiments.select-flag')}
                      options={flagsRemaining}
                      triggerClassName={
                        flagsSelected.length > 1
                          ? 'max-w-[calc(100%-36px)]'
                          : ''
                      }
                      additionalElement={item => (
                        <FeatureFlagStatus
                          status={t(
                            item.enabled
                              ? 'form:experiments.on'
                              : 'form:experiments.off'
                          )}
                          enabled={item.enabled as boolean}
                        />
                      )}
                      onSelectOption={value => field.onChange(value)}
                      onSearchChange={handleSearchChange}
                    />
                    {flagsSelected.length > 1 && (
                      <Button
                        type="button"
                        variant="grey"
                        className="size-5"
                        onClick={() =>
                          setValue(
                            'flags',
                            flagsSelected.filter((_, i) => i !== index)
                          )
                        }
                      >
                        <Icon icon={IconTrash} size="sm" />
                      </Button>
                    )}
                  </div>
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
        ))}
        {!isOnTargeting && (
          <Button
            type="button"
            variant="text"
            className="w-fit px-0 h-6"
            disabled={isDisabledAddBtn}
            onClick={() => setValue('flags', [...flagsSelected, ''])}
          >
            <Icon icon={IconPlus} size="md" />
            {t('form:add-flag')}
          </Button>
        )}
      </div>
    </>
  );
};

export default DebuggerFlags;
