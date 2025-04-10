import React, { useMemo } from 'react';
import { Trans } from 'react-i18next';
import { useTranslation } from 'i18n';
import { DomainEventEntityType } from '@types';
import { AuditLogsFilters } from 'pages/audit-logs/types';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';

interface Props {
  entityType?: DomainEventEntityType;
  isSystemAdmin: boolean;
  onChangeFilters: (filters: Partial<AuditLogsFilters>) => void;
}

const EntityTypeDropdown = ({
  entityType,
  isSystemAdmin,
  onChangeFilters
}: Props) => {
  const { t } = useTranslation(['common']);

  const options = useMemo(
    () => [
      {
        labelKey: 'source-type.feature-flag',
        value: 'FEATURE'
      },
      {
        labelKey: 'goal',
        value: 'GOAL'
      },
      {
        labelKey: 'source-type.experiment',
        value: 'EXPERIMENT'
      },
      {
        labelKey: 'navigation.user-segment',
        value: 'SEGMENT'
      },
      {
        labelKey: 'account',
        value: 'ACCOUNT'
      },
      {
        labelKey: 'source-type.api-key',
        value: 'APIKEY'
      },
      {
        labelKey: 'source-type.auto-operation',
        value: 'AUTOOPS_RULE'
      },
      {
        labelKey: 'source-type.progressive-rollout',
        value: 'PROGRESSIVE_ROLLOUT'
      },
      {
        labelKey: 'source-type.push',
        value: 'PUSH'
      },
      {
        labelKey: 'source-type.subscription',
        value: isSystemAdmin ? 'ADMIN_SUBSCRIPTION' : 'SUBSCRIPTION'
      }
    ],
    [isSystemAdmin]
  );
  console.log(entityType)
  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        label={
          entityType ? (
            <Trans
              i18nKey={'form:kind-filter-value'}
              values={{
                value: t(
                  options.find(item => item.value === entityType)?.labelKey ||
                    ''
                )
              }}
            />
          ) : (
            ''
          )
        }
        placeholder={
          <Trans
            i18nKey={'form:kind-filter-value'}
            values={{
              value: 'All'
            }}
          />
        }
      />
      <DropdownMenuContent>
        {options.map((item, index) => (
          <DropdownMenuItem
            key={index}
            label={t(item.labelKey)}
            value={item.value}
            onSelectOption={value =>
              onChangeFilters({
                entityType: value as DomainEventEntityType
              })
            }
          />
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default EntityTypeDropdown;
