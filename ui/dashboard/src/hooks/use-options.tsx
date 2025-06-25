import { useMemo } from 'react';
import { getLanguage, i18n } from 'i18n';
import {
  IconFlagJSON,
  IconFlagNumber,
  IconFlagString,
  IconFlagSwitch
} from '@icons';
import { StatusFilterType } from 'pages/feature-flags/types';

export interface FilterOption {
  value: string | number | undefined;
  label: string;
  filterValue?: string | boolean | number | string[];
}

export enum FilterTypes {
  ENABLED = 'enabled',
  STATUS = 'status',
  HAS_EXPERIMENT = 'hasExperiment',
  HAS_PREREQUISITES = 'hasPrerequisites',
  MAINTAINER = 'maintainer',
  TAGS = 'tags',
  ROLE = 'role',
  STATUSES = 'statuses',
  HAS_RULE = 'hasFeatureFlagAsRule',
  IN_USE = 'in-use',
  NOT_IN_USE = 'not-in-use'
}

const useOptions = () => {
  const language = getLanguage();
  const translation = (key: string, ns = 'common') => i18n.t(`${ns}:${key}`);

  const filterEnabledOptions: FilterOption[] = useMemo(
    () => [
      {
        value: FilterTypes.ENABLED,
        label: translation('enabled'),
        filterValue: ''
      }
    ],
    [language]
  );

  const filterStatusOptions: FilterOption[] = useMemo(
    () => [
      {
        value: FilterTypes.STATUS,
        label: translation('status'),
        filterValue: ''
      }
    ],
    [language]
  );

  const enabledOptions: FilterOption[] = useMemo(
    () => [
      {
        value: 'yes',
        label: translation('yes')
      },
      {
        value: 'no',
        label: translation('no')
      }
    ],
    [language]
  );

  const experimentStatusOptions: FilterOption[] = useMemo(
    () => [
      {
        value: 'WAITING',
        label: translation('experiment.waiting', 'table')
      },
      {
        value: 'RUNNING',
        label: translation('experiment.running', 'table')
      },
      {
        value: 'STOPPED',
        label: translation('experiment.stopped', 'table')
      },
      {
        value: 'FORCE_STOPPED',
        label: translation('experiment.force-stopped', 'table')
      }
    ],
    [language]
  );

  const flagFilterOptions: FilterOption[] = useMemo(
    () => [
      {
        value: FilterTypes.HAS_EXPERIMENT,
        label: translation('has-experiment'),
        filterValue: ''
      },
      {
        value: FilterTypes.HAS_PREREQUISITES,
        label: translation('has-prerequisites'),
        filterValue: ''
      },
      {
        value: FilterTypes.MAINTAINER,
        label: translation('maintainer'),
        filterValue: ''
      },
      {
        value: FilterTypes.ENABLED,
        label: translation('enabled'),
        filterValue: ''
      },
      {
        value: FilterTypes.TAGS,
        label: translation('tags'),
        filterValue: []
      },
      {
        value: FilterTypes.STATUS,
        label: translation('status'),
        filterValue: ''
      },
      {
        value: FilterTypes.HAS_RULE,
        label: translation('has-flag-as-rule'),
        filterValue: ''
      }
    ],
    [language]
  );

  const booleanOptions: FilterOption[] = useMemo(
    () => [
      {
        value: 1,
        label: translation('yes')
      },
      {
        value: 0,
        label: translation('no')
      }
    ],
    [language]
  );

  const roleOptions: FilterOption[] = useMemo(
    () => [
      {
        value: '1',
        label: translation('member')
      },
      {
        value: '2',
        label: translation('admin')
      },
      {
        value: '3',
        label: translation('owner')
      }
    ],
    [language]
  );

  const memberFilterOptions: FilterOption[] = useMemo(
    () => [
      {
        value: FilterTypes.ENABLED,
        label: translation('enabled'),
        filterValue: ''
      },
      {
        value: FilterTypes.ROLE,
        label: translation('role'),
        filterValue: ''
      }
    ],
    [language]
  );

  const experimentFilterOptions: FilterOption[] = useMemo(
    () => [
      {
        value: FilterTypes.STATUSES,
        label: translation('status'),
        filterValue: []
      },
      {
        value: FilterTypes.MAINTAINER,
        label: translation('maintainer'),
        filterValue: ''
      }
    ],
    [language]
  );

  const flagStatusOptions = useMemo(
    () => [
      {
        value: StatusFilterType.NEW,
        label: translation('new')
      },
      {
        value: StatusFilterType.ACTIVE,
        label: translation('active')
      },
      {
        value: StatusFilterType.NO_ACTIVITY,
        label: translation('no-activity')
      }
    ],
    [language]
  );

  const segmentStatusOptions: FilterOption[] = useMemo(
    () => [
      {
        value: FilterTypes.IN_USE,
        label: translation('in-use')
      },
      {
        value: FilterTypes.NOT_IN_USE,
        label: translation('not-in-use')
      }
    ],
    [language]
  );

  const flagTypeOptions = useMemo(
    () => [
      {
        label: translation('boolean', 'form'),
        value: 'BOOLEAN',
        icon: IconFlagSwitch
      },
      {
        label: translation('string', 'form'),
        value: 'STRING',
        icon: IconFlagString
      },
      {
        label: translation('number', 'form'),
        value: 'NUMBER',
        icon: IconFlagNumber
      },
      {
        label: 'JSON',
        value: 'JSON',
        icon: IconFlagJSON
      }
    ],
    [language]
  );

  const organizationRoles = useMemo(
    () => [
      {
        value: 'Organization_MEMBER',
        label: translation('member'),
        description: translation('member-role-description')
      },
      {
        value: 'Organization_ADMIN',
        label: translation('admin'),
        description: translation('admin-role-description')
      }
    ],
    [language]
  );

  const flagSortByOptions = useMemo(
    () => [
      {
        label: translation('name'),
        value: 'NAME'
      },
      {
        label: translation('tags'),
        value: 'TAGS'
      },
      {
        label: translation('created-at', 'table'),
        value: 'CREATED_AT'
      },
      {
        label: translation('updated-at', 'table'),
        value: 'UPDATED_AT'
      },
      {
        label: translation('enabled'),
        value: 'ENABLED'
      }
    ],
    [language]
  );

  const flagSortDirectionOptions = useMemo(
    () => [
      {
        label: translation('sort-asc'),
        value: 'ASC'
      },
      {
        label: translation('sort-desc'),
        value: 'DESC'
      }
    ],
    [language]
  );

  const environmentRoleOptions = useMemo(
    () => [
      {
        value: 'Environment_VIEWER',
        label: translation('viewer')
      },
      {
        value: 'Environment_EDITOR',
        label: translation('editor')
      }
    ],
    [language]
  );

  return {
    filterEnabledOptions,
    filterStatusOptions,
    enabledOptions,
    experimentStatusOptions,
    flagFilterOptions,
    booleanOptions,
    roleOptions,
    memberFilterOptions,
    experimentFilterOptions,
    flagStatusOptions,
    segmentStatusOptions,
    flagTypeOptions,
    organizationRoles,
    flagSortByOptions,
    flagSortDirectionOptions,
    environmentRoleOptions
  };
};

export default useOptions;
