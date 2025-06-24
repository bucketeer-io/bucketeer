import { i18n } from 'i18n';
import { StatusFilterType } from 'pages/feature-flags/types';

const translation = (key: string, ns = 'common') => i18n.t(`${ns}:${key}`);

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

export const filterEnabledOptions: FilterOption[] = [
  {
    value: FilterTypes.ENABLED,
    label: translation('enabled'),
    filterValue: ''
  }
];

export const filterStatusOptions: FilterOption[] = [
  {
    value: FilterTypes.STATUS,
    label: translation('status'),
    filterValue: ''
  }
];

export const enabledOptions: FilterOption[] = [
  {
    value: 'yes',
    label: translation('yes')
  },
  {
    value: 'no',
    label: translation('no')
  }
];

export const experimentStatusOptions: FilterOption[] = [
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
];

export const flagFilterOptions: FilterOption[] = [
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
];

export const booleanOptions: FilterOption[] = [
  {
    value: 1,
    label: translation('yes')
  },
  {
    value: 0,
    label: translation('no')
  }
];

export const roleOptions: FilterOption[] = [
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
];

export const memberFilterOptions: FilterOption[] = [
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
];

export const experimentFilterOptions: FilterOption[] = [
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
];

export const flagStatusOptions = [
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
];

export const segmentStatusOptions: FilterOption[] = [
  {
    value: FilterTypes.IN_USE,
    label: translation('in-use')
  },
  {
    value: FilterTypes.NOT_IN_USE,
    label: translation('not-in-use')
  }
];
