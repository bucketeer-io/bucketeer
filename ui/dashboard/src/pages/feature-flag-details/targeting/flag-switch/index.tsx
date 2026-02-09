import { useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { Feature } from '@types';
import { IconInfo } from '@icons';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Dropdown from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Switch from 'components/switch';
import { Tooltip } from 'components/tooltip';
import DisabledButtonTooltip from 'elements/disabled-button-tooltip';
import { TargetingSchema } from '../form-schema';

const FlagSwitch = ({
  feature,
  editable,
  setIsShowRules
}: {
  feature: Feature;
  editable: boolean;
  setIsShowRules: (value: boolean) => void;
}) => {
  const { t } = useTranslation(['form', 'common']);
  const { control, watch } = useFormContext<TargetingSchema>();

  const enabledWatch = watch('enabled');

  const options = useMemo(
    () =>
      feature.variations.map((item, index) => ({
        label: (
          <div className="flex items-center w-full gap-x-2">
            <FlagVariationPolygon index={index} />
            <p className="truncate">{item.name || item.value}</p>
          </div>
        ),
        value: item.id
      })),
    [feature]
  );

  return (
    <div className="flex flex-wrap items-center justify-between w-full p-2 sm:p-5 rounded-lg shadow-card-secondary">
      <Form.Field
        control={control}
        name="enabled"
        render={({ field }) => (
          <Form.Item className="w-fit py-0">
            <Form.Control>
              <div className="flex items-center w-full gap-x-2 typo-para-medium text-gray-700">
                <Trans
                  i18nKey={`form:targeting.flag-switch-${field.value ? 'on' : 'off'}`}
                  components={{
                    switch: (
                      <DisabledButtonTooltip
                        align="start"
                        hidden={editable}
                        trigger={
                          <div className="w-fit flex items-center">
                            <Switch
                              className="-mb-1"
                              disabled={!editable}
                              checked={!!field.value}
                              onCheckedChange={checked => {
                                field.onChange(checked, {
                                  shouldDirty: false
                                });
                                setIsShowRules(checked);
                              }}
                            />
                          </div>
                        }
                      />
                    )
                  }}
                />
                <Tooltip
                  content={t(
                    `targeting.tooltip.flag-${field.value ? 'on' : 'off'}`
                  )}
                  trigger={
                    <div className="flex-center size-fit -mb-1">
                      <Icon icon={IconInfo} size="xxs" color="gray-500" />
                    </div>
                  }
                  className="max-w-[400px]"
                />
              </div>
            </Form.Control>
          </Form.Item>
        )}
      />
      {!enabledWatch && (
        <Form.Field
          name="offVariation"
          control={control}
          render={({ field }) => (
            <Form.Item className="w-fit">
              <Form.Control>
                <Dropdown
                  value={field.value}
                  onChange={field.onChange}
                  options={options}
                  disabled={!editable}
                  className="max-w-[400px] w-full truncate"
                  contentClassName="max-w-[400px]"
                  alignContent="end"
                />
              </Form.Control>
              <Form.Message />
            </Form.Item>
          )}
        />
      )}
    </div>
  );
};

export default FlagSwitch;
