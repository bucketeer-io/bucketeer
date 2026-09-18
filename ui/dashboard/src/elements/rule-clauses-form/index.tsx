import { useCallback, useEffect, useMemo, useState } from 'react';
import { useFieldArray, useFormContext } from 'react-hook-form';
import { getCurrentEnvironment, useAuth } from 'auth';
import useOptions from 'hooks/use-options';
import { getLanguage, Language, useTranslation } from 'i18n';
import difference from 'lodash/difference';
import uniq from 'lodash/uniq';
import { v4 as uuid } from 'uuid';
import { Feature, FeatureRuleClauseOperator, UserSegment } from '@types';
import { IconPlus } from '@icons';
import { RuleClauseType } from 'pages/feature-flag-details/targeting/types';
import Button from 'components/button';
import Icon from 'components/icon';
import ClauseRow from './clause-row';

interface ClauseValue {
  id?: string;
  type?: string;
  attribute?: string;
  operator?: string;
  values?: string[];
}

export type SituationOption = {
  label: string;
  value: RuleClauseType;
};

interface Props {
  /** Name of the rules field array in the enclosing form. */
  rulesFieldName: string;
  ruleIndex: number;
  /** Clause types offered in the context-kind dropdown. */
  situationOptions: SituationOption[];
  /** Attribute keys already used in the enclosing form (suggested options). */
  usedAttributeKeys?: string[];
  feature?: Feature;
  features?: Feature[];
  sdkAttributeKeys: string[];
}

export const getSegmentSummary = (
  segment: UserSegment,
  t: (key: string) => string
) => {
  const userCount = Number(segment.includedUserCount) || 0;
  const ruleCount = segment.rules?.length || 0;
  const userLabel = t(
    `common:${userCount === 1 ? 'user' : 'users'}`
  ).toLowerCase();
  const ruleLabel = t(
    `common:${ruleCount === 1 ? 'rule' : 'rules'}`
  ).toLowerCase();
  return `${userCount} ${userLabel} · ${ruleCount} ${ruleLabel}`;
};

const RuleClausesForm = ({
  rulesFieldName,
  ruleIndex,
  situationOptions,
  usedAttributeKeys = [],
  feature,
  features,
  sdkAttributeKeys
}: Props) => {
  const { t } = useTranslation(['form', 'common', 'table']);
  const { conditionerCompareOptions, conditionerDateOptions } = useOptions();

  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const isLanguageJapanese = getLanguage() === Language.JAPANESE;

  const [createdOptionList, setCreatedOptionList] = useState<string[]>([]);

  const methods = useFormContext();
  const {
    control,
    formState: { errors },
    watch,
    setValue
  } = methods;

  const clausesName = `${rulesFieldName}.${ruleIndex}.clauses`;
  const clausesWatch: ClauseValue[] = watch(clausesName);
  const {
    fields: clauses,
    append,
    remove
  } = useFieldArray({
    control,
    name: clausesName,
    keyName: 'clauseId'
  });

  const formatClauses = (clausesWatch ?? []).map(item => ({
    ...item,
    clauseId: clauses.find(
      clause => (clause as unknown as ClauseValue).id === item.id
    )?.clauseId
  }));

  const flagOptions = useMemo(() => {
    const flagsSelected = (clausesWatch ?? [])
      .filter(item => item.type === RuleClauseType.FEATURE_FLAG)
      ?.map(item => item.attribute);

    return (features || [])
      .filter(
        item => ![...(flagsSelected ?? []), feature?.id]?.includes(item.id)
      )
      .map(item => ({
        label: item.name,
        value: item.id,
        enabled: item.enabled
      }));
  }, [features, [...(clausesWatch ?? [])], feature]);

  const getFieldName = (name: string, index: number) =>
    `${clausesName}.${index}.${name}`;

  const handleChangeConditioner = useCallback(
    (
      value: RuleClauseType,
      index: number,
      onChange: (value: RuleClauseType) => void
    ) => {
      let _value = '';
      switch (value) {
        case RuleClauseType.COMPARE:
          _value = FeatureRuleClauseOperator.EQUALS;
          break;
        case RuleClauseType.SEGMENT:
          _value = FeatureRuleClauseOperator.SEGMENT;
          break;
        case RuleClauseType.FEATURE_FLAG:
          _value = FeatureRuleClauseOperator.FEATURE_FLAG;
          break;
        case RuleClauseType.DATE:
          _value = FeatureRuleClauseOperator.BEFORE;
          break;
        default:
          break;
      }

      setValue(getFieldName('operator', index), _value, { shouldDirty: true });
      const currentType = watch(getFieldName('type', index));
      if (currentType !== value) {
        setValue(getFieldName('values', index), [], { shouldDirty: true });
        setValue(getFieldName('attribute', index), '', { shouldDirty: true });
      }

      onChange(value);
    },
    [clauses, ruleIndex]
  );

  useEffect(() => {
    setCreatedOptionList(
      difference(uniq(usedAttributeKeys), sdkAttributeKeys).sort()
    );
  }, [usedAttributeKeys, sdkAttributeKeys]);

  return (
    <>
      <div className="flex flex-col w-full gap-y-4">
        {formatClauses?.map((clause, clauseIndex) => (
          <ClauseRow
            key={clause.clauseId ?? clauseIndex}
            clause={clause}
            clauseIndex={clauseIndex}
            clausesName={clausesName}
            clausesLength={formatClauses.length}
            situationOptions={situationOptions}
            environmentId={currentEnvironment.id}
            environmentUrlCode={currentEnvironment.urlCode}
            isLanguageJapanese={isLanguageJapanese}
            control={control}
            errors={errors}
            feature={feature}
            features={features}
            flagOptions={flagOptions}
            createdOptionList={createdOptionList}
            sdkAttributeKeys={sdkAttributeKeys}
            conditionerDateOptions={conditionerDateOptions}
            conditionerCompareOptions={conditionerCompareOptions}
            onCreateOption={value =>
              setCreatedOptionList(prev => [...prev, value].sort())
            }
            onChangeConditioner={handleChangeConditioner}
            onRemoveClause={remove}
            getFieldName={getFieldName}
          />
        ))}
      </div>
      <Button
        type="button"
        variant={'text'}
        className="w-fit gap-x-2 h-6 !p-0"
        onClick={() =>
          append({
            id: uuid(),
            type: RuleClauseType.COMPARE,
            attribute: '',
            operator: FeatureRuleClauseOperator.EQUALS,
            values: []
          })
        }
      >
        <Icon
          icon={IconPlus}
          color="primary-500"
          className="flex-center"
          size={'sm'}
        />{' '}
        {t('form:feature-flags.add-condition')}
      </Button>
    </>
  );
};

export default RuleClausesForm;
