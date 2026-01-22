import { useMemo } from 'react';
import { useFormContext, useWatch } from 'react-hook-form';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useFeatureFlagsLoader } from 'hooks/use-feature-loading-more';
import { useTranslation } from 'i18n';
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
  const { control, setValue } = useFormContext<AddDebuggerFormType>();

  const flagsSelected: string[] = useWatch({ control, name: 'flags' });

  const {
    remainingFlagOptions,
    data: flagCollection,
    allAvailableFlags,
    isLoadingMore,
    hasMore,
    isInitialLoading,
    loadMore,
    onSearchChange
  } = useFeatureFlagsLoader({
    environmentId: currentEnvironment.id,
    selectedFlagIds: flagsSelected,
    currentFeatureId: feature?.id,
    filterSelected: true
  });

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

  const isDisabledAddBtn = useMemo(() => {
    // API returns totalCount as string; convert to number for comparison
    const totalFlagCount = Number(flagCollection?.totalCount ?? 0);

    // Disable if no remaining flags in current view OR all flags are selected
    return (
      !remainingFlagOptions.length || flagsSelected?.length >= totalFlagCount
    );
  }, [remainingFlagOptions, flagsSelected, flagCollection]);

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
                      isHasMore={hasMore || isLoadingMore}
                      onHasMoreOptions={loadMore}
                      disabled={isOnTargeting}
                      isLoadingMore={isLoadingMore}
                      isLoading={isInitialLoading}
                      placeholder={t('form:experiments.select-flag')}
                      options={remainingFlagOptions}
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
                      onSearchChange={onSearchChange}
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
