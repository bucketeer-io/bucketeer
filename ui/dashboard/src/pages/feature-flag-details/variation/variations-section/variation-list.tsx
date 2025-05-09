import { useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { useQueryRollouts } from '@queries/rollouts';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { IconInfo } from '@icons';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';
import { VariationProps } from '..';
import { VariationForm } from '../form-schema';
import Variations from './variations';

const VariationList = ({ feature }: VariationProps) => {
  const { t } = useTranslation(['common', 'form', 'table']);

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { data: rolloutCollection } = useQueryRollouts({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment?.id,
      featureIds: [feature?.id]
    },
    enabled: !!currentEnvironment?.id && !!feature?.id
  });

  const rollouts = rolloutCollection?.progressiveRollouts || [];

  const { control, watch } = useFormContext<VariationForm>();

  const offVariation = watch('offVariation');
  const variations = watch('variations');

  const offVariationId = useMemo(() => {
    const variation = variations.find(item => item.id === offVariation);
    return variation?.id || '';
  }, [offVariation, [...variations]]);

  return (
    <>
      <Form.Field
        control={control}
        name="variations"
        render={() => (
          <Form.Item className="flex flex-col w-full py-0">
            <Form.Control>
              <Variations feature={feature} rollouts={rollouts} />
            </Form.Control>
          </Form.Item>
        )}
      />

      <Form.Field
        control={control}
        name={'offVariation'}
        render={({ field }) => (
          <Form.Item className="pt-6 pb-0">
            <Form.Label required className="relative w-fit mb-6">
              {t('form:off-variation')}
              <Tooltip
                content={t('table:feature-flags.off-variation-tooltip')}
                trigger={
                  <div className="flex-center size-fit absolute top-0 -right-6">
                    <Icon icon={IconInfo} size="xs" color="gray-500" />
                  </div>
                }
                className="max-w-[310px]"
              />
            </Form.Label>
            <Form.Control>
              <DropdownMenu>
                <DropdownMenuTrigger
                  label={
                    variations.find(item => item.id === offVariationId)?.name ||
                    ''
                  }
                  isExpand
                />
                <DropdownMenuContent align="start">
                  {variations?.map((item, index) => (
                    <DropdownMenuItem
                      {...field}
                      key={index}
                      label={item.name || item.value}
                      value={item.id}
                      onSelectOption={value => field.onChange(value)}
                    />
                  ))}
                </DropdownMenuContent>
              </DropdownMenu>
            </Form.Control>
          </Form.Item>
        )}
      />
    </>
  );
};

export default VariationList;
