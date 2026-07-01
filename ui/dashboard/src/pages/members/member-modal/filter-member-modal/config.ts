import { useQueryTeams } from '@queries/teams';
import { getCurrentEnvironment, useAuth } from 'auth';
import { FilterTypes } from 'hooks/use-options';
import { i18n } from 'i18n';
import { MembersFilters } from 'pages/members/types';
import { FilterModalConfig } from 'elements/filter-modal/types';

const t = (key: string) => i18n.t(`common:${key}`);

// Matches `booleanOptions` in use-options (numeric 1 = yes / 0 = no).
const booleanOptions = () => [
  { value: 1, label: t('yes') },
  { value: 0, label: t('no') }
];

const roleOptions = () => [
  { value: '1', label: t('member') },
  { value: '2', label: t('admin') },
  { value: '3', label: t('owner') }
];

export const memberFilterConfig: FilterModalConfig<MembersFilters> = {
  mode: 'multi',
  fields: [
    {
      type: FilterTypes.ENABLED,
      labelKey: 'enabled',
      valueKind: 'boolean',
      emptyValue: '',
      useData: () => ({ options: booleanOptions() }),
      // Stored filterValue is 1 (yes) / 0 (no); `disabled` is the inverse.
      toFilter: filterValue => ({ disabled: !filterValue }),
      fromFilter: filters =>
        filters.disabled === undefined ? undefined : filters.disabled ? 0 : 1
    },
    {
      type: FilterTypes.ROLE,
      labelKey: 'role',
      valueKind: 'enum',
      emptyValue: '',
      useData: () => ({ options: roleOptions() }),
      toFilter: filterValue => ({ organizationRole: Number(filterValue) }),
      fromFilter: filters => filters.organizationRole?.toString()
    },
    {
      type: FilterTypes.TEAMS,
      labelKey: 'teams',
      valueKind: 'searchable',
      emptyValue: [],
      useData: ({ enabled }) => {
        const { consoleAccount } = useAuth();
        const currentEnvironment = getCurrentEnvironment(consoleAccount!);
        const { data, isLoading } = useQueryTeams({
          params: {
            cursor: String(0),
            organizationId: currentEnvironment.organizationId
          },
          enabled
        });
        const teams = data?.teams || [];
        const options = teams.map(item => ({
          label: item.name,
          value: item.name
        }));
        return {
          options,
          isLoading,
          getLabel: value =>
            (Array.isArray(value) &&
              teams.length &&
              value
                .map(item => teams.find(team => team.name === item)?.name)
                .join(', ')) ||
            ''
        };
      },
      toFilter: filterValue => ({ teams: filterValue as string[] }),
      fromFilter: filters => filters.teams
    }
  ]
};
