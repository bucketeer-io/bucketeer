import { useCallback } from 'react';
import { useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import {
  IconArrowDownwardFilled,
  IconArrowUpwardFilled
} from 'react-icons-material-design';
import { Fragment } from 'react/jsx-runtime';
import { useQueryUserSegments } from '@queries/user-segments';
import { getCurrentEnvironment, useAuth } from 'auth';
import { LIST_PAGE_SIZE } from 'constants/app';
import { useTranslation } from 'i18n';
import { Feature } from '@types';
import { IconClose, IconInfo } from '@icons';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';
import { TargetingDivider } from '..';
import Card from '../../elements/card';
import AddRule from '../add-rule';
import { RuleSchema, TargetingSchema } from '../form-schema';
import { RuleCategory } from '../types';
import { getDefaultRolloutStrategy } from '../utils';
import RuleForm from './rule';
import SegmentVariation from './variation';

interface RuleSchemaFields extends RuleSchema {
  segmentId: string;
}

interface Props {
  feature: Feature;
  features: Feature[];
  segmentRules: RuleSchemaFields[];
  isDisableAddIndividualRules: boolean;
  isDisableAddPrerequisite: boolean;
  onAddRule: (rule: RuleCategory) => void;
  segmentRulesRemove: (index: number) => void;
  segmentRulesSwap: (indexA: number, indexB: number) => void;
}

const TargetSegmentRule = ({
  feature,
  features,
  segmentRules,
  isDisableAddIndividualRules,
  isDisableAddPrerequisite,
  onAddRule,
  segmentRulesRemove,
  segmentRulesSwap
}: Props) => {
  const { t } = useTranslation(['table', 'form']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { data: segmentCollection } = useQueryUserSegments({
    params: {
      cursor: String(0),
      pageSize: LIST_PAGE_SIZE,
      environmentId: currentEnvironment.id
    }
  });

  const userSegments = segmentCollection?.segments || [];

  const methods = useFormContext<TargetingSchema>();

  const { watch } = methods;

  const segmentRulesWatch = watch('segmentRules');

  const handleChangeIndexRule = useCallback(
    (type: 'increase' | 'decrease', currentIndex: number) => {
      segmentRulesSwap(
        currentIndex,
        type === 'increase' ? currentIndex + 1 : currentIndex - 1
      );
    },

    [segmentRulesWatch, segmentRulesSwap]
  );

  return (
    segmentRules.length > 0 && (
      <div className="flex flex-col w-full">
        {segmentRules.map((segment, segmentIndex) => (
          <div key={segment?.segmentId} className="flex flex-col w-full">
            {segmentIndex !== 0 && (
              <>
                <TargetingDivider />
                <AddRule
                  isDisableAddIndividualRules={isDisableAddIndividualRules}
                  isDisableAddPrerequisite={isDisableAddPrerequisite}
                  onAddRule={onAddRule}
                />
                <TargetingDivider />
              </>
            )}
            <Card>
              <div className="flex items-center justify-between w-full">
                <div className="flex items-center gap-x-2">
                  <p className="typo-para-medium leading-5 text-gray-700">
                    <Trans
                      i18nKey={'table:feature-flags.rule-index'}
                      values={{
                        index: segmentIndex + 1
                      }}
                    />
                  </p>
                  <Tooltip
                    align="start"
                    alignOffset={-68}
                    content={t('form:targeting.tooltip.custom-rules')}
                    trigger={
                      <div className="flex-center size-fit">
                        <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
                      </div>
                    }
                    className="max-w-[400px]"
                  />
                </div>

                <div className="flex items-center gap-x-2">
                  <div
                    className="flex-center cursor-pointer group"
                    onClick={() => segmentRulesRemove(segmentIndex)}
                  >
                    <Icon
                      icon={IconClose}
                      size={'sm'}
                      className="flex-center text-gray-500 group-hover:text-gray-700"
                    />
                  </div>
                  {segmentRules.length > 1 && (
                    <div className="flex items-center gap-x-1">
                      {segmentIndex !== segmentRules.length - 1 && (
                        <div
                          className="flex-center group cursor-pointer"
                          onClick={() =>
                            handleChangeIndexRule('increase', segmentIndex)
                          }
                        >
                          <Icon
                            icon={IconArrowDownwardFilled}
                            size={'sm'}
                            className="text-gray-500 group-hover:text-gray-700"
                          />
                        </div>
                      )}
                      {segmentIndex !== 0 && (
                        <div
                          className="flex-center group cursor-pointer"
                          onClick={() =>
                            handleChangeIndexRule('decrease', segmentIndex)
                          }
                        >
                          <Icon
                            icon={IconArrowUpwardFilled}
                            size={'sm'}
                            className="text-gray-500 group-hover:text-gray-700"
                          />
                        </div>
                      )}
                    </div>
                  )}
                </div>
              </div>
              <Fragment>
                <RuleForm
                  feature={feature}
                  features={features}
                  segmentIndex={segmentIndex}
                  userSegments={userSegments}
                />
                <SegmentVariation
                  feature={feature}
                  defaultRolloutStrategy={getDefaultRolloutStrategy(feature)}
                  segmentIndex={segmentIndex}
                  segmentRules={segmentRules}
                />
              </Fragment>
            </Card>
          </div>
        ))}
      </div>
    )
  );
};

export default TargetSegmentRule;
