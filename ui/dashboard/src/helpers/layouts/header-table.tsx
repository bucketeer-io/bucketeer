import { TableHeaders } from '@types';

export const orgDetailsUserHeader: TableHeaders = [
  {
    text: 'NAME',
    sort: true
  },
  {
    text: 'ROLE',
    sort: true
  },
  {
    text: 'ENVIRONMENT',
    sort: true
  },
  {
    text: 'CREATED AT',
    sort: true
  },
  { type: 'empty' }
];

export const orgDetailsEnvHeader: TableHeaders = [
  {
    text: 'NAME',
    sort: true
  },
  {
    text: 'PROJECT',
    sort: true
  },
  {
    text: 'CREATED AT',
    sort: true
  },
  {
    text: 'STATE',
    sort: true
  },
  {
    type: 'empty'
  }
];

export const orgDetailsProjectHeader: TableHeaders = [
  {
    text: 'NAME',
    sort: true
  },
  {
    text: 'MAINTAINER',
    sort: true
  },
  {
    text: 'ENVIRONMENT',
    sort: true
  },
  {
    text: 'CREATED AT',
    sort: true
  },
  { type: 'empty' }
];

export const orgHeader: TableHeaders = [
  {
    text: 'NAME',
    sort: true,
    width: '40%'
  },
  {
    text: 'PROJECTS',
    sort: true
  },
  {
    text: 'ENVIRONMENTS',
    sort: true
  },
  {
    text: 'USERS',
    sort: true
  },
  {
    text: 'CREATED AT',
    sort: true
  },
  { type: 'empty' }
];

export const projectsHeader: TableHeaders = [
  {
    text: 'NAME',
    sort: true
  },
  {
    text: 'MAINTAINER',
    sort: true
  },
  {
    text: 'ENVIRONMENTS',
    sort: true
  },
  {
    text: 'FLAGS',
    sort: true
  },
  {
    text: 'CREATED AT',
    sort: true
  },

  { type: 'empty' }
];

export const prjDetailsEnvHeader: TableHeaders = [
  {
    text: 'NAME',
    sort: true
  },
  {
    text: 'PROJECT',
    sort: true
  },
  {
    text: 'STATUS',
    sort: true
  },
  {
    text: 'CREATED AT',
    sort: true
  },
  {
    type: 'empty'
  }
];

export const membersHeader: TableHeaders = [
  {
    type: 'checkbox',
    width: 'fit-content'
  },
  {
    text: 'NAME',
    sort: true
  },
  {
    text: 'ROLE',
    sort: true
  },
  {
    text: 'ENVIRONMENT',
    sort: true
  },
  {
    text: 'LAST SEEN',
    sort: true
  },
  {
    text: 'STATE',
    sort: true
  },
  { type: 'empty' }
];
