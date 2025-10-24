import { memo, useMemo } from 'react';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { DomainEventEntityMap } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { AuditLogsFilters } from 'pages/audit-logs/types';
import Dropdown from 'components/dropdown';

interface Props {
  entityType?: DomainEventEntityMap;
  isSystemAdmin: boolean;
  onChangeFilters: (filters: Partial<AuditLogsFilters>) => void;
}

const EntityTypeDropdown = memo(
  ({ entityType, isSystemAdmin, onChangeFilters }: Props) => {
    const { t } = useTranslation(['common', 'form', 'table']);

    const options = useMemo(() => {
      const {
        FEATURE,
        GOAL,
        EXPERIMENT,
        SEGMENT,
        ACCOUNT,
        APIKEY,
        AUTOOPS_RULE,
        PROGRESSIVE_ROLLOUT,
        PUSH,
        ADMIN_SUBSCRIPTION,
        SUBSCRIPTION
      } = DomainEventEntityMap;

      return [
        {
          labelKey: 'source-type.feature-flag',
          value: FEATURE
        },
        {
          labelKey: 'goal',
          value: GOAL
        },
        {
          labelKey: 'source-type.experiment',
          value: EXPERIMENT
        },
        {
          labelKey: 'navigation.user-segment',
          value: SEGMENT
        },
        {
          labelKey: 'account',
          value: ACCOUNT
        },
        {
          labelKey: 'source-type.api-key',
          value: APIKEY
        },
        {
          labelKey: 'source-type.auto-operation',
          value: AUTOOPS_RULE
        },
        {
          labelKey: 'source-type.progressive-rollout',
          value: PROGRESSIVE_ROLLOUT
        },
        {
          labelKey: 'source-type.push',
          value: PUSH
        },
        {
          labelKey: 'source-type.subscription',
          value: isSystemAdmin ? ADMIN_SUBSCRIPTION : SUBSCRIPTION
        }
      ];
    }, [isSystemAdmin]);

    const entityLabel = useMemo(() => {
      const labelKey = options.find(
        item => item.value === Number(entityType)
      )?.labelKey;

      return labelKey ? t(labelKey) : '';
    }, [options, entityType, isSystemAdmin]);
    const optionCustom = options.map(opt => ({
      value: opt.value,
      label: t(opt.labelKey)
    }));
    const label = useMemo(
      () =>
        isNotEmpty(entityType) ? (
          <Trans
            i18nKey={'form:kind-filter-value'}
            values={{
              value: entityLabel
            }}
          />
        ) : (
          ''
        ),
      [entityType]
    );
    const placeholder = (
      <Trans
        i18nKey={'form:kind-filter-value'}
        values={{
          value: t('table:code-refs.all')
        }}
      />
    );
    return (
      <Dropdown
        options={optionCustom}
        value={Number(entityType)}
        labelCustom={label}
        cleanable
        placeholder={placeholder}
        onClear={() =>
          onChangeFilters({
            entityType: undefined
          })
        }
        onChange={value =>
          onChangeFilters({
            entityType: +value
          })
        }
        className="max-w-[175px] xxl:max-w-fit [&>div>p]:!text-gray-700"
        contentClassName="min-w-[180px]"
        alignContent="end"
      />
    );
  }
);

export default EntityTypeDropdown;
