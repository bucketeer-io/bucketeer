import { useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { Feature } from '@types';
import { cn } from 'utils/style';
import { IconInfo } from '@icons';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import Switch from 'components/switch';
import { Tooltip } from 'components/tooltip';
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

  const options = feature.variations.map((item, index) => ({
    label: (
      <div className="flex items-center gap-x-2">
        <FlagVariationPolygon index={index} />
        {item.name || item.value}
      </div>
    ),
    value: item.id
  }));

  return (
    <div
      className={cn(
        'flex items-center justify-between w-full p-5 rounded-lg shadow-card-secondary',
        {
          'p-4': !enabledWatch
        }
      )}
    >
      <Form.Field
        control={control}
        name="enabled"
        render={({ field }) => (
          <Form.Item className="w-full py-0">
            <Form.Control>
              <div className="flex items-center w-full gap-x-2 typo-para-medium text-gray-700">
                <Trans
                  i18nKey={`form:targeting.flag-switch-${field.value ? 'on' : 'off'}`}
                  components={{
                    switch: (
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
            <Form.Item>
              <Form.Control>
                <DropdownMenu>
                  <DropdownMenuTrigger
                    label={
                      options.find(item => item.value === field.value)?.label ||
                      ''
                    }
                    disabled={!editable}
                  />
                  <DropdownMenuContent align="end">
                    {options.map((item, index) => (
                      <DropdownMenuItem
                        key={index}
                        label={item.label}
                        value={item.value}
                        onSelectOption={field.onChange}
                      />
                    ))}
                  </DropdownMenuContent>
                </DropdownMenu>
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
