import { useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { useQueryFeatures } from '@queries/features';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { IconPlus, IconTrash } from '@icons';
import Button from 'components/button';
import Form from 'components/form';
import Icon from 'components/icon';
import DropdownMenuWithSearch from 'elements/dropdown-with-search';
import FeatureFlagStatus from 'elements/feature-flag-status';
import { AddDebuggerFormType } from './form-schema';

const DebuggerFlags = () => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { t } = useTranslation(['common', 'form']);
  const { control, watch, setValue } = useFormContext<AddDebuggerFormType>();

  const { data: flagCollection } = useQueryFeatures({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment.id
    }
  });

  const flags = flagCollection?.features || [];

  const flagsSelected: string[] = [...watch('flags')];
  const flagOptions = useMemo(
    () =>
      flags.map(item => ({
        label: item.name,
        value: item.id,
        enabled: item.enabled
      })),
    [flags]
  );

  const flagsRemaining = useMemo(() => {
    return flagOptions.filter(item => !flagsSelected.includes(item.value));
  }, [flagsSelected, flagOptions, flags]);

  const isDisabledAddBtn = useMemo(
    () => !flagsRemaining.length || flagsSelected?.length === flags.length,
    [flagsRemaining, flagsSelected, flags]
  );

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
      </div>
    </>
  );
};

export default DebuggerFlags;
