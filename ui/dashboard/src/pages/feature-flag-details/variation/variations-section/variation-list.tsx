import { useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { useQueryAutoOpsRules } from '@queries/auto-ops-rules';
import { useQueryFeatures } from '@queries/features';
import { useQueryRollouts } from '@queries/rollouts';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import isNil from 'lodash/isNil';
import { IconInfo } from '@icons';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import Dropdown from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';
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
      environmentId: currentEnvironment?.id,
      featureIds: [feature?.id]
    },
    enabled: !isNil(currentEnvironment?.id) && !!feature?.id
  });

  const rollouts = rolloutCollection?.progressiveRollouts || [];

  const { data: operationCollection } = useQueryAutoOpsRules({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment?.id,
      featureIds: [feature?.id]
    }
  });

  const { data: collection } = useQueryFeatures({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment.id
    },
    enabled: !!currentEnvironment
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
    label: (
      <div className="flex items-center gap-x-2 pl-0.5">
        <FlagVariationPolygon index={index} />
        <span className="typo-para-medium text-gray-700">
          {item.name || item.value}
        </span>
      </div>
    ),
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
                features={collection?.features || []}
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
              <Dropdown
                options={variationOptions}
                value={field.value}
                onChange={field.onChange}
                trigger={
                  offVariationId ? (
                    <div className="flex items-center gap-x-2 pl-0.5 w-0 flex-1 typo-para-medium text-gray-700">
                      <FlagVariationPolygon
                        index={variations.findIndex(
                          v => v.id === offVariationId
                        )}
                      />
                      <p className="truncate">
                        {variations.find(v => v.id === offVariationId)?.name ||
                          variations.find(v => v.id === offVariationId)?.value}
                      </p>
                    </div>
                  ) : undefined
                }
                disabled={isRunningExperiment || !editable}
              />
            </Form.Control>
          </Form.Item>
        )}
      />
    </>
  );
};

export default VariationList;
