import { i18n } from 'i18n';

const translation = (key: string, ns = 'common') => i18n.t(`${ns}:${key}`);

export interface FilterOption {
  value: string | number;
  label: string;
}

export enum FilterTypes {
  ENABLED = 'enabled',
  STATUS = 'status',
  HAS_EXPERIMENT = 'hasExperiment',
  HAS_PREREQUISITES = 'hasPrerequisites',
  MAINTAINER = 'maintainer',
  TAGS = 'tags',
  ROLE = 'role'
}

export const filterEnabledOptions: FilterOption[] = [
  {
    value: FilterTypes.ENABLED,
    label: translation('enabled')
  }
];

export const filterStatusOptions: FilterOption[] = [
  {
    value: FilterTypes.STATUS,
    label: translation('status')
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

export const statusOptions: FilterOption[] = [
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
    label: translation('has-experiment')
  },
  {
    value: FilterTypes.HAS_PREREQUISITES,
    label: translation('has-prerequisites')
  },
  {
    value: FilterTypes.MAINTAINER,
    label: translation('maintainer')
  },
  {
    value: FilterTypes.ENABLED,
    label: translation('enabled')
  },
  {
    value: FilterTypes.TAGS,
    label: translation('tags')
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
    label: translation('enabled')
  },
  {
    value: FilterTypes.ROLE,
    label: translation('role')
  }
];
