import { useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { useQueryAutoOpsRules } from '@queries/auto-ops-rules';
import { useQueryRollouts } from '@queries/rollouts';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { isNil } from 'lodash';
import { checkEnvironmentEmptyId } from 'utils/function';
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
import VariationLabel from 'elements/variation-label';
import { VariationProps } from '..';
import { VariationForm } from '../form-schema';
import Variations from './variations';

const VariationList = ({
  feature,
  isRunningExperiment,
  editable
}: VariationProps) => {
  const { t } = useTranslation(['common', 'form', 'table']);

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { data: rolloutCollection } = useQueryRollouts({
    params: {
      cursor: String(0),
      environmentId: checkEnvironmentEmptyId(currentEnvironment?.id),
      featureIds: [feature?.id]
    },
    enabled: !isNil(currentEnvironment?.id) && !!feature?.id
  });

  const rollouts = rolloutCollection?.progressiveRollouts || [];

  const { data: operationCollection } = useQueryAutoOpsRules({
    params: {
      cursor: String(0),
      environmentId: checkEnvironmentEmptyId(currentEnvironment?.id),
      featureIds: [feature?.id]
    }
  });

  const autoOps = operationCollection?.autoOpsRules || [];
  const eventRateOperations = autoOps?.filter(
    item =>
      item.opsType === 'EVENT_RATE' &&
      ['WAITING', 'RUNNING'].includes(item.autoOpsStatus)
  );

  const { control, watch } = useFormContext<VariationForm>();

  const offVariation = watch('offVariation');
  const variations = watch('variations');

  const variationOptions = variations.map((item, index) => ({
    label: <VariationLabel label={item.name || item.value} index={index} />,
    value: item.id
  }));

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
              <Variations
                feature={feature}
                rollouts={rollouts}
                isRunningExperiment={isRunningExperiment}
                eventRateOperations={eventRateOperations}
                editable={editable}
              />
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
                    variationOptions.find(item => item.value === offVariationId)
                      ?.label || ''
                  }
                  isExpand
                  disabled={isRunningExperiment || !editable}
                />
                <DropdownMenuContent align="start">
                  {variationOptions?.map((item, index) => (
                    <DropdownMenuItem
                      {...field}
                      key={index}
                      label={item.label}
                      value={item.value}
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
