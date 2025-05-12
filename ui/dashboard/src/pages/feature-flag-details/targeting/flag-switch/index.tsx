import { useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { cn, getVariationColor } from 'utils/style';
import { IconInfo } from '@icons';
import { Polygon } from 'pages/experiment-details/elements/header-details';
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

const FlagSwitch = () => {
  const { t } = useTranslation(['form', 'common']);
  const { control, watch } = useFormContext<TargetingSchema>();

  const enabledWatch = watch('enabled');

  const options = [
    {
      label: t('common:true'),
      value: 1
    },
    {
      label: t('common:false'),
      value: 0
    }
  ];

  return (
    <div
      className={cn('flex items-center w-full p-5 rounded-lg shadow-card', {
        'p-4': !enabledWatch
      })}
    >
      <Form.Field
        control={control}
        name="enabled"
        render={({ field }) => (
          <Form.Item className="w-full py-0">
            <Form.Control>
              <div className="flex items-center w-full justify-between">
                <div className="flex items-center w-full gap-x-2 typo-para-medium text-gray-700">
                  <Trans
                    i18nKey={`form:targeting.flag-switch-${field.value ? 'on' : 'off'}`}
                    components={{
                      switch: (
                        <Switch
                          className="-mb-1"
                          checked={!!field.value}
                          onCheckedChange={field.onChange}
                        />
                      )
                    }}
                  />
                  <Tooltip
                    content={t('targeting.all-audience-traffic')}
                    trigger={
                      <div className="flex-center size-fit -mb-1">
                        <Icon icon={IconInfo} size="xxs" color="gray-500" />
                      </div>
                    }
                  />
                </div>
                {!enabledWatch && (
                  <DropdownMenu>
                    <DropdownMenuTrigger
                      trigger={
                        <div className="flex items-center gap-x-2">
                          <div className="ml-0.5">
                            <Polygon
                              style={{
                                background: getVariationColor(
                                  field.value ? 0 : 1
                                )
                              }}
                            />
                          </div>
                          <p className="capitalize">{`${field.value}`}</p>
                        </div>
                      }
                    />
                    <DropdownMenuContent>
                      {options.map((item, index) => (
                        <DropdownMenuItem
                          key={index}
                          iconElement={
                            <Polygon
                              style={{
                                background: getVariationColor(index)
                              }}
                            />
                          }
                          label={item.label}
                          value={item.value}
                          onSelectOption={field.onChange}
                        />
                      ))}
                    </DropdownMenuContent>
                  </DropdownMenu>
                )}
              </div>
            </Form.Control>
          </Form.Item>
        )}
      />
    </div>
  );
};

export default FlagSwitch;
