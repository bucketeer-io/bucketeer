import { defineMessage } from 'react-intl';

export const messages = {
  adminProject: {
    action: {
      convertProject: defineMessage({
        id: 'adminProject.action.convertProject',
        defaultMessage: 'Convert to paid version',
      }),
    },
    add: {
      header: {
        title: defineMessage({
          id: 'adminProject.add.header.title',
          defaultMessage: 'Create a project',
        }),
        description: defineMessage({
          id: 'adminProject.add.header.description',
          defaultMessage:
            'You can manage multiple different user projects by using projects.',
        }),
      },
    },
    confirm: {
      convertProjectTitle: defineMessage({
        id: 'adminProject.confirm.convertProject.title',
        defaultMessage: 'Convert project',
      }),
      convertProjectDescription: defineMessage({
        id: 'adminProject.confirm.convertProject.description',
        defaultMessage:
          'Are you sure you want to convert the {projectId} project to paid version?',
      }),
      enableTitle: defineMessage({
        id: 'adminProject.confirm.enable.title',
        defaultMessage: 'Enable project',
      }),
      enableDescription: defineMessage({
        id: 'adminProject.confirm.enable.description',
        defaultMessage:
          'Are you sure you want to enable the {projectId} project?',
      }),
      disableTitle: defineMessage({
        id: 'adminProject.confirm.disable.title',
        defaultMessage: 'Disable project',
      }),
      disableDescription: defineMessage({
        id: 'adminProject.confirm.disable.description',
        defaultMessage:
          'Are you sure you want to disable the {projectId} project?',
      }),
    },
    creator: defineMessage({
      id: 'adminProject.creator',
      defaultMessage: 'Creator',
    }),
    filter: {
      enabled: defineMessage({
        id: 'adminProject.filter.enabled',
        defaultMessage: 'Enabled',
      }),
    },
    list: {
      header: {
        title: defineMessage({
          id: 'adminProject.list.header.title',
          defaultMessage: 'Projects',
        }),
        description: defineMessage({
          id: 'adminProject.list.header.description',
          defaultMessage:
            'On this page, you can check all projects. Select a project to update or click on the Add button to add a new one.',
        }),
      },
      noResult: {
        searchKeyword: defineMessage({
          id: 'adminProject.list.noResult.searchKeyword',
          defaultMessage: 'ID and email',
        }),
      },
    },
    search: {
      placeholder: defineMessage({
        id: 'adminProject.search.placeholder',
        defaultMessage: 'ID and email',
      }),
    },
    sort: {
      nameAz: defineMessage({
        id: 'adminProject.sort.nameAz',
        defaultMessage: 'Name A-Z',
      }),
      idZa: defineMessage({
        id: 'adminProject.sort.nameZa',
        defaultMessage: 'Name Z-A',
      }),
      newest: defineMessage({
        id: 'adminProject.sort.newest',
        defaultMessage: 'Newest',
      }),
      oldest: defineMessage({
        id: 'adminProject.sort.oldest',
        defaultMessage: 'Oldest',
      }),
    },
    trialPeriod: defineMessage({
      id: 'adminProject.trialPeriod',
      defaultMessage: 'Trial period',
    }),
    update: {
      header: {
        title: defineMessage({
          id: 'adminProject.update.header.title',
          defaultMessage: 'Update the project',
        }),
        description: defineMessage({
          id: 'adminProject.update.header.description',
          defaultMessage:
            'You can manage multiple different user projects by using projects.',
        }),
      },
    },
  },
  adminSettings: {
    list: {
      header: {
        title: defineMessage({
          id: 'adminSettings.list.header.title',
          defaultMessage: 'Admin Settings',
        }),
        description: defineMessage({
          id: 'adminSettings.list.header.description',
          defaultMessage:
            'On this page, you can check all admin settings. Select a tab to manage the settings.',
        }),
      },
    },
    tab: {
      auditLogs: defineMessage({
        id: 'adminSettings.tab.auditLogs',
        defaultMessage: 'Audit Logs',
      }),
      notifications: defineMessage({
        id: 'adminSettings.tab.notifications',
        defaultMessage: 'Notifications',
      }),
      projects: defineMessage({
        id: 'adminSettings.tab.projects',
        defaultMessage: 'Projects',
      }),
      environments: defineMessage({
        id: 'adminSettings.tab.environments',
        defaultMessage: 'Environments',
      }),
    },
  },
  adminAuditLog: {
    list: {
      header: {
        description: defineMessage({
          id: 'adminAuditLog.list.header.description',
          defaultMessage: 'On this tab you can check all admin audit logs.',
        }),
      },
    },
  },
  adminEnvironment: {
    add: {
      header: {
        title: defineMessage({
          id: 'adminEnvironment.add.header.title',
          defaultMessage: 'Create an environment',
        }),
        description: defineMessage({
          id: 'adminEnvironment.add.header.description',
          defaultMessage:
            "You can manage your feature flag's development lifecycle, from local development through production.",
        }),
      },
    },
    filter: {
      project: defineMessage({
        id: 'adminEnvironment.filter.project',
        defaultMessage: 'Project',
      }),
    },
    list: {
      header: {
        title: defineMessage({
          id: 'adminEnvironment.list.header.title',
          defaultMessage: 'Environments',
        }),
        description: defineMessage({
          id: 'adminEnvironment.list.header.description',
          defaultMessage:
            'On this tab, you can check all environments. Select an environment to update or click on the Add button to add a new one.',
        }),
      },
      noResult: {
        searchKeyword: defineMessage({
          id: 'adminEnvironment.list.noResult.searchKeyword',
          defaultMessage: 'ID and description',
        }),
      },
    },
    search: {
      placeholder: defineMessage({
        id: 'adminEnvironment.search.placeholder',
        defaultMessage: 'ID and description',
      }),
    },
    sort: {
      nameAz: defineMessage({
        id: 'adminEnvironment.sort.nameAz',
        defaultMessage: 'Name A-Z',
      }),
      nameZa: defineMessage({
        id: 'adminEnvironment.sort.nameZa',
        defaultMessage: 'Name Z-A',
      }),
      newest: defineMessage({
        id: 'adminEnvironment.sort.newest',
        defaultMessage: 'Newest',
      }),
      oldest: defineMessage({
        id: 'adminEnvironment.sort.oldest',
        defaultMessage: 'Oldest',
      }),
    },
    update: {
      header: {
        title: defineMessage({
          id: 'adminEnvironment.update.header.title',
          defaultMessage: 'Update the environment',
        }),
        description: defineMessage({
          id: 'adminEnvironment.update.header.description',
          defaultMessage:
            "You can manage your feature flag's development lifecycle, from local development through production.",
        }),
      },
    },
  },
  autoOps: {
    rule: defineMessage({
      id: 'autoOps.rule',
      defaultMessage: 'Rule',
    }),
    operation: defineMessage({
      id: 'autoOps.operation',
      defaultMessage: 'Operation',
    }),
    operationType: defineMessage({
      id: 'autoOps.operationType',
      defaultMessage: 'Operation type',
    }),
    enableFeatureType: defineMessage({
      id: 'autoOps.enableFeatureType',
      defaultMessage: 'Enable feature',
    }),
    disableFeatureType: defineMessage({
      id: 'autoOps.disableFeatureType',
      defaultMessage: 'Disable feature',
    }),
    clauseType: defineMessage({
      id: 'autoOps.clauseType',
      defaultMessage: 'Rule type',
    }),
    eventRateClauseType: defineMessage({
      id: 'autoOps.eventRateClauseType',
      defaultMessage: 'Event rate',
    }),
    datetimeClauseType: defineMessage({
      id: 'autoOps.datetimeClauseType',
      defaultMessage: 'Schedule',
    }),
    opsEventRateClause: {
      featureVersion: defineMessage({
        id: 'autoOps.opsEventRateClause.featureVersion',
        defaultMessage: 'Feature version',
      }),
      minCount: defineMessage({
        id: 'autoOps.opsEventRateClause.minCount',
        defaultMessage: 'Minimum count',
      }),
      goal: defineMessage({
        id: 'autoOps.opsEventRateClause.goal',
        defaultMessage: 'Goal',
      }),
    },
    datetimeClause: {
      datetime: defineMessage({
        id: 'autoOps.datetime.datetime',
        defaultMessage: 'Date time',
      }),
    },
    condition: defineMessage({
      id: 'autoOps.condition',
      defaultMessage: 'Condition',
    }),
    threshold: defineMessage({
      id: 'autoOps.threshold',
      defaultMessage: 'Threshold',
    }),
    enable: defineMessage({
      id: 'autoOps.enable',
      defaultMessage: 'Enable',
    }),
    killSwitch: defineMessage({
      id: 'autoOps.killSwitch',
      defaultMessage: 'Kill Switch',
    }),
    schedule: defineMessage({
      id: 'autoOps.schedule',
      defaultMessage: 'Schedule',
    }),
    eventRate: defineMessage({
      id: 'autoOps.eventRate',
      defaultMessage: 'Event Rate',
    }),
    startDate: defineMessage({
      id: 'autoOps.startDate',
      defaultMessage: 'Start Date',
    }),
    infoBlocks: {
      title: defineMessage({
        id: 'infoBlocks.title',
        defaultMessage:
          'You can safely switch a flag on and off by using auto operations',
      }),
      scheduleInfo: defineMessage({
        id: 'infoBlocks.scheduleInfo',
        defaultMessage: 'Schedule a flag to turn on or off',
      }),
      killSwitch: defineMessage({
        id: 'infoBlocks.killSwitch',
        defaultMessage: 'Kill Switch',
      }),
      killSwitchInfo: defineMessage({
        id: 'infoBlocks.killSwitchInfo',
        defaultMessage: 'Turn off automatically a flag based on KPI events',
      }),
      progressiveRollout: defineMessage({
        id: 'infoBlocks.progressiveRollout',
        defaultMessage: 'Progressive Rollout',
      }),
      progressiveRolloutInfo: defineMessage({
        id: 'infoBlocks.progressiveRolloutInfo',
        defaultMessage: 'Coming soon',
      }),
    },
    editOperation: defineMessage({
      id: 'autoOps.editOperation',
      defaultMessage: 'Edit Operation',
    }),
    operationDetails: defineMessage({
      id: 'autoOps.operationDetails',
      defaultMessage: 'Operation Details',
    }),
    deleteOperation: defineMessage({
      id: 'autoOps.deleteOperation',
      defaultMessage: 'Delete Operation',
    }),
    createAnOperation: defineMessage({
      id: 'autoOps.createAnOperation',
      defaultMessage: 'Create An Operation',
    }),
    updateAnOperation: defineMessage({
      id: 'autoOps.updateAnOperation',
      defaultMessage: 'Update An Operation',
    }),
    minimumGoalCount: defineMessage({
      id: 'autoOps.minimumGoalCount',
      defaultMessage: 'Minimum Goal Count',
    }),
    totalGoalCountEvents: defineMessage({
      id: 'autoOps.totalGoalCountEvents',
      defaultMessage: 'Total Goal Count Events',
    }),
    currentEventRate: defineMessage({
      id: 'autoOps.currentEventRate',
      defaultMessage: 'Current Event Rate',
    }),
    enableOperation: defineMessage({
      id: 'autoOps.enableOperation',
      defaultMessage: 'Enable Operation',
    }),
    killSwitchOperation: defineMessage({
      id: 'autoOps.killSwitchOperation',
      defaultMessage: 'Kill Switch Operation',
    }),
    progressInformation: defineMessage({
      id: 'autoOps.progressInformation',
      defaultMessage: 'Progress information',
    }),
    active: defineMessage({
      id: 'autoOps.active',
      defaultMessage: 'Active',
    }),
    completed: defineMessage({
      id: 'autoOps.completed',
      defaultMessage: 'Completed',
    }),
    goalCount: defineMessage({
      id: 'autoOps.goalCount',
      defaultMessage: 'Goal Count',
    }),
    evaluationCount: defineMessage({
      id: 'autoOps.evaluationCount',
      defaultMessage: 'Evaluation Count',
    }),
  },
  trigger: {
    documentation: defineMessage({
      id: 'trigger.documentation',
      defaultMessage: 'documentation',
    }),
    description: defineMessage({
      id: 'trigger.description',
      defaultMessage:
        'Use triggers to turn a flag on or off remotely. See the {link}',
    }),
    addTrigger: defineMessage({
      id: 'trigger.addTrigger',
      defaultMessage: 'Add Trigger',
    }),
    triggerType: defineMessage({
      id: 'trigger.triggerType',
      defaultMessage: 'Type',
    }),
    action: defineMessage({
      id: 'trigger.action',
      defaultMessage: 'Action',
    }),
    triggerURL: defineMessage({
      id: 'trigger.triggerURL',
      defaultMessage: 'Trigger URL',
    }),
    triggeredTimes: defineMessage({
      id: 'trigger.triggeredTimes',
      defaultMessage: 'Triggered Times',
    }),
    lastTriggered: defineMessage({
      id: 'trigger.lastTriggered',
      defaultMessage: 'Last Triggered',
    }),
    turnTheFlagON: defineMessage({
      id: 'trigger.turnTheFlagON',
      defaultMessage: 'Turn the flag ON',
    }),
    turnTheFlagOFF: defineMessage({
      id: 'trigger.turnTheFlagOFF',
      defaultMessage: 'Turn the flag OFF',
    }),
    editDescription: defineMessage({
      id: 'trigger.editDescription',
      defaultMessage: 'Edit description',
    }),
    enableTrigger: defineMessage({
      id: 'trigger.enableTrigger',
      defaultMessage: 'Enable Trigger',
    }),
    disableTrigger: defineMessage({
      id: 'trigger.disableTrigger',
      defaultMessage: 'Disable Trigger',
    }),
    resetURL: defineMessage({
      id: 'trigger.resetURL',
      defaultMessage: 'Reset URL',
    }),
    deleteTrigger: defineMessage({
      id: 'trigger.deleteTrigger',
      defaultMessage: 'Delete Trigger',
    }),
    triggerUrlTitle: defineMessage({
      id: 'trigger.triggerUrlTitle',
      defaultMessage: 'Copy and store this URL.',
    }),
    triggerUrlDescription: defineMessage({
      id: 'trigger.triggerUrlDescription',
      defaultMessage: 'Once you leave this page, the URL will be hidden.',
    }),
    deleteTriggerDialogTitle: defineMessage({
      id: 'trigger.deleteTriggerDialogTitle',
      defaultMessage: 'Delete Trigger',
    }),
    deleteTriggerDialogMessage: defineMessage({
      id: 'trigger.deleteTriggerDialogMessage',
      defaultMessage: 'The trigger will be deleted permanently.',
    }),
    deleteTriggerDialogBtnLabel: defineMessage({
      id: 'trigger.deleteTriggerDialogBtnTxt',
      defaultMessage: 'Delete',
    }),
    resetTriggerDialogTitle: defineMessage({
      id: 'trigger.resetTriggerDialogTitle',
      defaultMessage: 'Reset Trigger URL',
    }),
    resetTriggerDialogMessage: defineMessage({
      id: 'trigger.resetTriggerDialogMessage',
      defaultMessage:
        'The current URL will become invalid. Ensure that you copy and store the new URL.',
    }),
    resetTriggerDialogBtnLabel: defineMessage({
      id: 'trigger.resetTriggerDialogBtnLabel',
      defaultMessage: 'Reset',
    }),
    updated: defineMessage({
      id: 'trigger.updated',
      defaultMessage: 'Updated',
    }),
  },
  maintainer: defineMessage({
    id: 'maintainer',
    defaultMessage: 'Maintainer',
  }),
  filter: {
    filter: defineMessage({
      id: 'filter.filter',
      defaultMessage: 'Filter',
    }),
    add: defineMessage({
      id: 'filter.add',
      defaultMessage: 'Add filter',
    }),
  },
  notFound: {
    title: defineMessage({
      id: 'notFound.title',
      defaultMessage: 'Page not found',
    }),
    description: defineMessage({
      id: 'notFound.description',
      defaultMessage: `Sorry, we couldn't find the page you're looking for.`,
    }),
    goBackHome: defineMessage({
      id: 'notFound.goBackHome',
      defaultMessage: 'Go back home',
    }),
  },
  error: defineMessage({
    id: 'error',
    defaultMessage: 'Error',
  }),
  warning: defineMessage({
    id: 'warning',
    defaultMessage: 'Warning',
  }),
  success: defineMessage({
    id: 'success',
    defaultMessage: 'Success',
  }),
  total: defineMessage({
    id: 'total',
    defaultMessage: 'Total',
  }),
  copy: {
    copied: defineMessage({
      id: 'copy.copied',
      defaultMessage: 'Copied!',
    }),
    copyToClipboard: defineMessage({
      id: 'copy.copyToClipboard',
      defaultMessage: 'Copy to clipboard',
    }),
  },
  sort: defineMessage({
    id: 'sort',
    defaultMessage: 'Sort',
  }),
  environment: {
    select: {
      label: defineMessage({
        id: 'environment.select.label',
        defaultMessage: '(Project) Environment',
      }),
    },
  },
  created: defineMessage({
    id: 'created',
    defaultMessage: 'Created',
  }),
  description: defineMessage({
    id: 'description',
    defaultMessage: 'Description',
  }),
  yes: defineMessage({
    id: 'yes',
    defaultMessage: 'Yes',
  }),
  no: defineMessage({
    id: 'no',
    defaultMessage: 'No',
  }),
  enabled: defineMessage({
    id: 'enabled',
    defaultMessage: 'Enabled',
  }),
  disabled: defineMessage({
    id: 'disabled',
    defaultMessage: 'Disabled',
  }),
  close: defineMessage({
    id: 'close',
    defaultMessage: 'Close',
  }),
  seeMore: defineMessage({
    id: 'seeMore',
    defaultMessage: 'See more',
  }),
  action: defineMessage({
    id: 'action',
    defaultMessage: 'Action',
  }),
  all: defineMessage({
    id: 'all',
    defaultMessage: 'All',
  }),
  show: defineMessage({
    id: 'show',
    defaultMessage: 'Show',
  }),
  mostRecent: defineMessage({
    id: 'mostRecent',
    defaultMessage: 'Most Recent',
  }),
  button: {
    archive: defineMessage({
      id: 'button.archive',
      defaultMessage: 'Archive',
    }),
    copyFlags: defineMessage({
      id: 'button.copyFlags',
      defaultMessage: 'Copy flags',
    }),
    add: defineMessage({
      id: 'button.add',
      defaultMessage: 'Add',
    }),
    edit: defineMessage({
      id: 'button.edit',
      defaultMessage: 'Edit',
    }),
    save: defineMessage({
      id: 'button.save',
      defaultMessage: 'Save',
    }),
    result: defineMessage({
      id: 'button.result',
      defaultMessage: 'Result',
    }),
    clearAll: defineMessage({
      id: 'button.clearAll',
      defaultMessage: 'Clear All',
    }),
    addVariation: defineMessage({
      id: 'button.addVariation',
      defaultMessage: 'Add variation',
    }),
    addRule: defineMessage({
      id: 'button.addRule',
      defaultMessage: 'Add rule',
    }),
    addOperation: defineMessage({
      id: 'button.addOperation',
      defaultMessage: 'Add operation',
    }),
    addCondition: defineMessage({
      id: 'button.addCondition',
      defaultMessage: 'Add condition',
    }),
    cancel: defineMessage({
      id: 'button.cancel',
      defaultMessage: 'Cancel',
    }),
    saveWithComment: defineMessage({
      id: 'button.saveWithComment',
      defaultMessage: 'Save with comment',
    }),
    submit: defineMessage({
      id: 'button.submit',
      defaultMessage: 'Submit',
    }),
    enable: defineMessage({
      id: 'button.enable',
      defaultMessage: 'Enable',
    }),
    disable: defineMessage({
      id: 'button.disable',
      defaultMessage: 'Disable',
    }),
    schedule: defineMessage({
      id: 'button.schedule',
      defaultMessage: 'Schedule',
    }),
  },
  account: {
    confirm: {
      enableTitle: defineMessage({
        id: 'account.confirm.enable.title',
        defaultMessage: 'Enable account',
      }),
      enableDescription: defineMessage({
        id: 'account.confirm.enable.description',
        defaultMessage: 'Are you sure you want to enable this account?',
      }),
      disableTitle: defineMessage({
        id: 'account.confirm.disable.title',
        defaultMessage: 'Disable account',
      }),
      disableDescription: defineMessage({
        id: 'account.confirm.disable.description',
        defaultMessage: 'Are you sure you want to disable this account?',
      }),
    },
    filter: {
      role: defineMessage({
        id: 'account.filter.role',
        defaultMessage: 'Role',
      }),
      enabled: defineMessage({
        id: 'account.filter.enabled',
        defaultMessage: 'Enabled',
      }),
    },
    add: {
      header: {
        title: defineMessage({
          id: 'account.add.header.title',
          defaultMessage: 'Create an account',
        }),
        description: defineMessage({
          id: 'account.add.header.description',
          defaultMessage:
            'The account is required to access the admin console. The account has three roles: viewer, editor, and owner.',
        }),
      },
    },
    update: {
      header: {
        title: defineMessage({
          id: 'account.update.header.title',
          defaultMessage: 'Update the account',
        }),
        description: defineMessage({
          id: 'account.update.header.description',
          defaultMessage:
            'The account is required to access the admin console. The account has three roles: viewer, editor, and owner.',
        }),
      },
    },
    list: {
      header: {
        title: defineMessage({
          id: 'account.list.header.title',
          defaultMessage: 'Accounts',
        }),
        description: defineMessage({
          id: 'account.list.header.description',
          defaultMessage:
            'On this page, you can check all accounts for this environment. Select an account to manage the role settings or click on the Add button to add a new one.',
        }),
      },
      noData: {
        description: defineMessage({
          id: 'account.list.noData.description',
          defaultMessage:
            'You can add new team members, disable, or manage access controls for members by setting roles.',
        }),
      },
      noResult: {
        searchKeyword: defineMessage({
          id: 'account.list.noResult.searchKeyword',
          defaultMessage: 'email',
        }),
      },
    },
    role: {
      viewer: defineMessage({
        id: 'account.role.viewer',
        defaultMessage: 'Viewer',
      }),
      editor: defineMessage({
        id: 'account.role.editor',
        defaultMessage: 'Editor',
      }),
      owner: defineMessage({
        id: 'account.role.owner',
        defaultMessage: 'Owner',
      }),
    },
    search: {
      placeholder: defineMessage({
        id: 'account.search.placeholder',
        defaultMessage: 'Email',
      }),
    },
    sort: {
      emailAz: defineMessage({
        id: 'account.sort.emailAz',
        defaultMessage: 'Email A-Z',
      }),
      emailZa: defineMessage({
        id: 'account.sort.emailZa',
        defaultMessage: 'Email Z-A',
      }),
      newest: defineMessage({
        id: 'account.sort.newest',
        defaultMessage: 'Newest',
      }),
      oldest: defineMessage({
        id: 'account.sort.oldest',
        defaultMessage: 'Oldest',
      }),
    },
  },
  apiKey: {
    add: {
      header: {
        title: defineMessage({
          id: 'apiKey.add.header.title',
          defaultMessage: 'Create an API Key',
        }),
        description: defineMessage({
          id: 'apiKey.add.header.description',
          defaultMessage:
            'The API key is required for the client SDK to access the server API.',
        }),
      },
    },
    update: {
      header: {
        title: defineMessage({
          id: 'apiKey.update.header.title',
          defaultMessage: 'Update the API Key',
        }),
        description: defineMessage({
          id: 'apiKey.update.header.description',
          defaultMessage:
            'The API key is required for the client SDK to access the server API.',
        }),
      },
    },
    list: {
      header: {
        title: defineMessage({
          id: 'apiKey.list.header.title',
          defaultMessage: 'API Keys',
        }),
        description: defineMessage({
          id: 'apiKey.list.header.description',
          defaultMessage:
            'On this page, you can check all API Keys for this environment. Select an API Key to manage the settings or click on the Add button to add a new one.',
        }),
      },
      noData: {
        description: defineMessage({
          id: 'apiKey.list.noData.description',
          defaultMessage:
            'You can add an API Key to allow requests from the client SDK.',
        }),
      },
      noResult: {
        searchKeyword: defineMessage({
          id: 'apiKey.list.noResult.searchKeyword',
          defaultMessage: 'name',
        }),
      },
    },
    confirm: {
      enableTitle: defineMessage({
        id: 'apiKey.confirm.enable.title',
        defaultMessage: 'Enable API Key',
      }),
      enableDescription: defineMessage({
        id: 'apiKey.confirm.enable.description',
        defaultMessage: 'Are you sure you want to enable this API Key?',
      }),
      disableTitle: defineMessage({
        id: 'apiKey.confirm.disable.title',
        defaultMessage: 'Disable API Key',
      }),
      disableDescription: defineMessage({
        id: 'apiKey.confirm.disable.description',
        defaultMessage: 'Are you sure you want to disable this API Key?',
      }),
    },
    search: {
      placeholder: defineMessage({
        id: 'apiKey.search.placeholder',
        defaultMessage: 'Name',
      }),
    },
    filter: {
      enabled: defineMessage({
        id: 'apiKey.filter.enabled',
        defaultMessage: 'Enabled',
      }),
    },
  },
  auditLog: {
    list: {
      header: {
        title: defineMessage({
          id: 'auditLog.list.header.title',
          defaultMessage: 'Audit Logs',
        }),
        description: defineMessage({
          id: 'auditLog.list.header.description',
          defaultMessage:
            'On this page, you can check all audit logs for this environment.',
        }),
      },
      noData: {
        description: defineMessage({
          id: 'auditLog.list.noData.description',
          defaultMessage:
            'The history will be created when you add or edit something on the Admin Console.',
        }),
      },
      noResult: {
        searchKeyword: defineMessage({
          id: 'auditLog.list.noResult.searchKeyword',
          defaultMessage: 'email',
        }),
      },
    },
    filter: {
      dates: defineMessage({
        id: 'auditLog.filter.dates',
        defaultMessage: 'Dates',
      }),
      type: defineMessage({
        id: 'auditLog.filter.type',
        defaultMessage: 'Type',
      }),
      clear: defineMessage({
        id: 'auditLog.filter.clear',
        defaultMessage: 'Clear',
      }),
      apply: defineMessage({
        id: 'auditLog.filter.apply',
        defaultMessage: 'Apply',
      }),
      cancel: defineMessage({
        id: 'auditLog.filter.cancel',
        defaultMessage: 'Cancel',
      }),
    },
    search: {
      placeholder: defineMessage({
        id: 'auditLog.search.placeholder',
        defaultMessage: 'Email',
      }),
    },
    sort: {
      newest: defineMessage({
        id: 'auditLog.sort.newest',
        defaultMessage: 'Newest',
      }),
      oldest: defineMessage({
        id: 'auditLog.sort.oldest',
        defaultMessage: 'Oldest',
      }),
    },
  },
  goal: {
    action: {
      archive: defineMessage({
        id: 'goal.action.archive',
        defaultMessage: 'Archive',
      }),
    },
    add: {
      header: {
        title: defineMessage({
          id: 'goal.add.header.title',
          defaultMessage: 'Create a goal',
        }),
        description: defineMessage({
          id: 'goal.add.header.description',
          defaultMessage:
            'The goal lets you measure user behaviors affected by your feature flags in experiments.',
        }),
      },
    },
    confirm: {
      archiveTitle: defineMessage({
        id: 'goal.confirm.archive.title',
        defaultMessage: 'Archive goal',
      }),
      archiveDescription: defineMessage({
        id: 'goal.confirm.archive.description',
        defaultMessage:
          'We recommend removing the code references to "{goalId}" from your application before archiving.',
      }),
    },
    update: {
      header: {
        title: defineMessage({
          id: 'goal.update.header.title',
          defaultMessage: 'Update the goal',
        }),
        description: defineMessage({
          id: 'goal.update.header.description',
          defaultMessage:
            'The goal lets you measure user behaviors affected by your feature flags in experiments.',
        }),
      },
    },
    list: {
      header: {
        title: defineMessage({
          id: 'goal.list.header.title',
          defaultMessage: 'Goals',
        }),
        description: defineMessage({
          id: 'goal.list.header.description',
          defaultMessage:
            'Use this page to see all goals in this environment. Select a goal to manage settings.',
        }),
      },
      noData: {
        description: defineMessage({
          id: 'goal.list.noData.description',
          defaultMessage:
            'Goals are the metrics used to measure the effectiveness of a Feature Flag.',
        }),
      },
      noResult: {
        searchKeyword: defineMessage({
          id: 'goal.list.noResult.searchKeyword',
          defaultMessage: 'name and description',
        }),
      },
    },
    filter: {
      status: defineMessage({
        id: 'goal.filter.status',
        defaultMessage: 'Status',
      }),
      archived: defineMessage({
        id: 'goal.filter.archived',
        defaultMessage: 'Archived',
      }),
    },
    status: {
      status: defineMessage({
        id: 'goal.status.status',
        defaultMessage: 'Status',
      }),
      inUse: defineMessage({
        id: 'goal.status.inUse',
        defaultMessage: 'in use',
      }),
      notInUse: defineMessage({
        id: 'goal.status.notInUse',
        defaultMessage: 'not in use',
      }),
    },
  },
  experiment: {
    action: {
      archive: defineMessage({
        id: 'experiment.action.archive',
        defaultMessage: 'Archive',
      }),
      archiveTooltip: defineMessage({
        id: 'experiment.action.archiveTooltip',
        defaultMessage: 'Please stop the experiment before archiving.',
      }),
    },
    confirm: {
      archiveTitle: defineMessage({
        id: 'experiment.confirm.archive.title',
        defaultMessage: 'Archive experiment',
      }),
      archiveDescription: defineMessage({
        id: 'experiment.confirm.archive.description',
        defaultMessage:
          'Are you sure you want to archive the {experimentName} experiment?',
      }),
    },
    filter: {
      archived: defineMessage({
        id: 'experiment.filter.archived',
        defaultMessage: 'Archived',
      }),
      maintainer: defineMessage({
        id: 'experiment.filter.maintainer',
        defaultMessage: 'Maintainer',
      }),
      status: defineMessage({
        id: 'experiment.filter.status',
        defaultMessage: 'Status',
      }),
    },
    maintainer: defineMessage({
      id: 'experiment.maintainer',
      defaultMessage: 'Maintainer',
    }),
    result: {
      noData: {
        errorMessage: defineMessage({
          id: 'experiment.result.noData.errorMessage',
          defaultMessage: 'The data is not ready. Please come back later.',
        }),
        experimentResult: defineMessage({
          id: 'experiment.result.noData.experimentResult',
          defaultMessage: 'Experiment result',
        }),
        description: defineMessage({
          id: 'experiment.result.noData.description',
          defaultMessage: 'The result is created when the experiment starts.',
        }),
      },
      variation: {
        label: defineMessage({
          id: 'experiment.result.variation.label',
          defaultMessage: 'Variation',
        }),
      },
      goals: {
        label: defineMessage({
          id: 'experiment.result.goals.label',
          defaultMessage: 'Goal total',
        }),
        helpText: defineMessage({
          id: 'experiment.result.goals.helpText',
          defaultMessage:
            'The total number of goal events fired by the client.',
        }),
      },
      goalUser: {
        label: defineMessage({
          id: 'experiment.result.goalUser.label',
          defaultMessage: 'Goal user',
        }),
        helpText: defineMessage({
          id: 'experiment.result.goalUser.helpText',
          defaultMessage:
            'The number of unique users who fired the goal event. The count will not increase if the same user reaches the goal event multiple times.',
        }),
      },
      evaluationUser: {
        label: defineMessage({
          id: 'experiment.result.evaluationUser.label',
          defaultMessage: 'Evaluation user',
        }),
        helpText: defineMessage({
          id: 'experiment.result.evaluationUser.helpText',
          defaultMessage:
            'The number of unique users for which variations have been returned. The number of users actually assigned to the feature flag variation. It does not include offline users or potential new users.',
        }),
      },
      valueSum: {
        label: defineMessage({
          id: 'experiment.result.valueSum.label',
          defaultMessage: 'Value total',
        }),
        helpText: defineMessage({
          id: 'experiment.result.valueSum.helpText',
          defaultMessage:
            'The total number of values assigned to a goal event. This value is different for each goal.',
        }),
      },
      valuePerUser: {
        label: defineMessage({
          id: 'experiment.result.valuePerUser.label',
          defaultMessage: 'Value/User',
        }),
        helpText: defineMessage({
          id: 'experiment.result.valuePerUser.helpText',
          defaultMessage:
            'The total number of values assigned to the goal event per user. It is calculated as (the sum of the numbers assigned to the goal event / the number of unique users who fired the goal event).',
        }),
      },
      conversionRate: {
        label: defineMessage({
          id: 'experiment.result.conversionRate.label',
          defaultMessage: 'Conversion rate',
        }),
        helpText: defineMessage({
          id: 'experiment.result.conversionRate.helpText',
          defaultMessage:
            'Calculated as (number of unique users who fired the goal event / number of unique users for whom a variation was returned).',
        }),
      },
      improvement: {
        label: defineMessage({
          id: 'experiment.result.improvement.label',
          defaultMessage: 'Improvement',
        }),
        helpText: defineMessage({
          id: 'experiment.result.improvement.helpText',
          defaultMessage:
            'A measure of improvement in an indicator related to variation compared to a baseline (also called a control group). It is calculated by comparing the range of values for variation with the range of values for baseline.',
        }),
      },
      probabilityToBeatBaseline: {
        label: defineMessage({
          id: 'experiment.result.probabilityToBeatBaseline.label',
          defaultMessage: 'Probability to beat baseline',
        }),
        helpText: defineMessage({
          id: 'experiment.result.probabilityToBeatBaseline.helpText',
          defaultMessage:
            'Estimated likelihood of exceeding baseline (also known as a control group). A criterion of at least 95% is recommended.',
        }),
      },
      probabilityToBest: {
        label: defineMessage({
          id: 'experiment.result.probabilityToBest.label',
          defaultMessage: 'Probability to best',
        }),
        helpText: defineMessage({
          id: 'experiment.result.probabilityToBest.helpText',
          defaultMessage:
            'Possibility of being the best variation. Possibility of being presumed to outperform all other variations. We recommend a criterion of at least 95%.',
        }),
      },
    },
    add: {
      header: {
        title: defineMessage({
          id: 'experiment.add.header.title',
          defaultMessage: 'Create a experiment',
        }),
        description: defineMessage({
          id: 'experiment.add.header.description',
          defaultMessage:
            'Get started by filling in the information below to create your new experiment.',
        }),
      },
    },
    update: {
      header: {
        title: defineMessage({
          id: 'experiment.update.header.title',
          defaultMessage: 'Update a experiment',
        }),
        description: defineMessage({
          id: 'experiment.update.header.description',
          defaultMessage:
            'Fill in the information below to update your experiment.',
        }),
      },
    },
    stop: {
      dialog: {
        title: defineMessage({
          id: 'experiment.stop.dialog.title',
          defaultMessage: 'Confirm',
        }),
        description: defineMessage({
          id: 'experiment.stop.dialog.description',
          defaultMessage: 'Do you really stop an experiment?',
        }),
      },
      button: defineMessage({
        id: 'experiment.stop.stopExperiment',
        defaultMessage: 'Stop experiment',
      }),
    },
    list: {
      header: {
        title: defineMessage({
          id: 'experiment.list.header.title',
          defaultMessage: 'Experiments',
        }),
        description: defineMessage({
          id: 'experiment.list.header.description',
          defaultMessage:
            'Use this page to see all experiments in this environment. Select an experiment to manage settings and display the results.',
        }),
      },
      noData: {
        description: defineMessage({
          id: 'experiment.list.noData.description',
          defaultMessage:
            'By using Experiments, you can improve web page loading time, test new features, etc.',
        }),
      },
      noResult: {
        searchKeyword: defineMessage({
          id: 'experiment.list.noResult.searchKeyword',
          defaultMessage: 'name and description',
        }),
      },
    },
    feature: defineMessage({
      id: 'experiment.feature',
      defaultMessage: 'Feature flag',
    }),
    baselineVariation: defineMessage({
      id: 'experiment.baselineVariation',
      defaultMessage: 'Baseline variation',
    }),
    goalIds: defineMessage({
      id: 'experiment.goalIds',
      defaultMessage: 'Goals',
    }),
    startAt: defineMessage({
      id: 'experiment.startAt',
      defaultMessage: 'Start at',
    }),
    stopAt: defineMessage({
      id: 'experiment.stopAt',
      defaultMessage: 'Stop at',
    }),
    period: defineMessage({
      id: 'experiment.period',
      defaultMessage: 'Period',
    }),
    search: {
      placeholder: defineMessage({
        id: 'experiment.search.placeholder',
        defaultMessage: 'Name, Description',
      }),
    },
    status: {
      status: defineMessage({
        id: 'experiment.status.status',
        defaultMessage: 'Status',
      }),
      waiting: defineMessage({
        id: 'experiment.status.waiting',
        defaultMessage: 'Waiting',
      }),
      running: defineMessage({
        id: 'experiment.status.running',
        defaultMessage: 'Running',
      }),
      stopped: defineMessage({
        id: 'experiment.status.stopped',
        defaultMessage: 'Finished',
      }),
      forceStopped: defineMessage({
        id: 'experiment.status.forceStopped',
        defaultMessage: 'Stopped',
      }),
    },
  },
  feature: {
    action: {
      archive: defineMessage({
        id: 'feature.action.archive',
        defaultMessage: 'Archive',
      }),
      clone: defineMessage({
        id: 'feature.action.clone',
        defaultMessage: 'Clone',
      }),
      unarchive: defineMessage({
        id: 'feature.action.unarchive',
        defaultMessage: 'Unarchive',
      }),
    },
    id: defineMessage({
      id: 'feature.id',
      defaultMessage: 'ID',
    }),
    rule: defineMessage({
      id: 'feature.rule',
      defaultMessage: 'Rule rollout percentage',
    }),
    clause: {
      type: {
        compare: defineMessage({
          id: 'feature.clause.type.compare',
          defaultMessage: 'Compare',
        }),
        segment: defineMessage({
          id: 'feature.clause.type.segment',
          defaultMessage: 'User segment',
        }),
        date: defineMessage({
          id: 'feature.clause.type.date',
          defaultMessage: 'Date',
        }),
      },
      operator: {
        equal: defineMessage({
          id: 'feature.clause.operator.equal',
          defaultMessage: '=',
        }),
        greaterOrEqual: defineMessage({
          id: 'feature.clause.operator.greaterOrEqual',
          defaultMessage: '>=',
        }),
        greater: defineMessage({
          id: 'feature.clause.operator.greater',
          defaultMessage: '>',
        }),
        less: defineMessage({
          id: 'feature.clause.operator.less',
          defaultMessage: '<',
        }),
        lessOrEqual: defineMessage({
          id: 'feature.clause.operator.lessOrEqual',
          defaultMessage: '<=',
        }),
        in: defineMessage({
          id: 'feature.clause.operator.in',
          defaultMessage: 'contains',
        }),
        startWith: defineMessage({
          id: 'feature.clause.operator.startWith',
          defaultMessage: 'starts with',
        }),
        endWith: defineMessage({
          id: 'feature.clause.operator.endWith',
          defaultMessage: 'ends with',
        }),
        before: defineMessage({
          id: 'feature.clause.operator.before',
          defaultMessage: 'before',
        }),
        after: defineMessage({
          id: 'feature.clause.operator.after',
          defaultMessage: 'after',
        }),
        segment: defineMessage({
          id: 'feature.clause.operator.segment',
          defaultMessage: 'is included in',
        }),
      },
    },
    strategy: {
      selectRolloutPercentage: defineMessage({
        id: 'feature.strategy.selectRolloutPercentage',
        defaultMessage: 'Select rollout percentage',
      }),
    },
    flagStatus: {
      new: defineMessage({
        id: 'feature.flagStatus.new',
        defaultMessage: 'New',
      }),
      receivingRequests: defineMessage({
        id: 'feature.flagStatus.receivingRequests',
        defaultMessage: 'Receiving requests',
      }),
      inactive: defineMessage({
        id: 'feature.flagStatus.inactive',
        defaultMessage: 'Inactive',
      }),
    },
    targetingDescription: defineMessage({
      id: 'feature.targetingDescription',
      defaultMessage:
        'Enable targeting settings. You can configure targeting users, complex rules, default strategy, and off variation.',
    }),
    flagIsPrerequisite: defineMessage({
      id: 'feature.flagIsPrerequisite',
      defaultMessage:
        'This flag is a prerequisite of {length} other flag{length, plural, one {} other {s}}.',
    }),
    flagIsPrerequisiteDescription: defineMessage({
      id: 'feature.flagIsPrerequisiteDescription',
      defaultMessage:
        'Changes to the targeting rules may affect the variations served by the flag{length, plural, one {} other {s}} below.',
    }),
    prerequisites: defineMessage({
      id: 'feature.prerequisites',
      defaultMessage: 'Prerequisites',
    }),
    addPrerequisites: defineMessage({
      id: 'feature.addPrerequisites',
      defaultMessage: 'Add Prerequisites',
    }),
    selectFlag: defineMessage({
      id: 'feature.selectFlag',
      defaultMessage: 'Select a feature flag',
    }),
    selectVariation: defineMessage({
      id: 'feature.selectVariation',
      defaultMessage: 'Select a variation',
    }),
    targetingUsers: defineMessage({
      id: 'feature.targetings',
      defaultMessage: 'Individual targeting',
    }),
    addUser: defineMessage({
      id: 'feature.addUser',
      defaultMessage: 'Add user {userId}',
    }),
    addUserIds: defineMessage({
      id: 'feature.addUserIds',
      defaultMessage: 'Add user ids',
    }),
    alreadyTargeted: defineMessage({
      id: 'feature.alreadyTargeted',
      defaultMessage: 'Already targeted',
    }),
    alreadyTargetedInVariation: defineMessage({
      id: 'feature.alreadyTargetedInVariation',
      defaultMessage: '"{userId}" is already targeted in "{variationName}"',
    }),
    updateComment: defineMessage({
      id: 'feature.updateComment',
      defaultMessage: 'Comment for update',
    }),
    resetRandomSampling: defineMessage({
      id: 'feature.resetRandomSampling',
      defaultMessage: 'Reset random sampling',
    }),
    variationType: defineMessage({
      id: 'feature.variationType',
      defaultMessage: 'Flag type',
    }),
    type: {
      boolean: defineMessage({
        id: 'feature.type.boolean',
        defaultMessage: 'boolean',
      }),
      string: defineMessage({
        id: 'feature.type.string',
        defaultMessage: 'string',
      }),
      number: defineMessage({
        id: 'feature.type.number',
        defaultMessage: 'number',
      }),
      json: defineMessage({
        id: 'feature.type.json',
        defaultMessage: 'json',
      }),
    },
    status: defineMessage({
      id: 'feature.status',
      defaultMessage: 'status',
    }),
    variation: defineMessage({
      id: 'feature.variation',
      defaultMessage: 'variation',
    }),
    defaultStrategy: defineMessage({
      id: 'feature.defaultStrategy',
      defaultMessage: 'Default strategy',
    }),
    onVariation: defineMessage({
      id: 'feature.onVariation',
      defaultMessage: 'on variation',
    }),
    offVariation: defineMessage({
      id: 'feature.offVariation',
      defaultMessage: 'off variation',
    }),
    variationSettings: {
      defaultStrategy: defineMessage({
        id: 'feature.variationSettings.defaultStrategy',
        defaultMessage:
          'This variation cannot be deleted because it is used in the default variation settings.',
      }),
      offVariation: defineMessage({
        id: 'feature.variationSettings.offVariation',
        defaultMessage:
          'This variation cannot be deleted because it is used in the off-variation settings.',
      }),
      bothVariations: defineMessage({
        id: 'feature.variationSettings.bothVariations',
        defaultMessage:
          'This variation cannot be deleted because it is used in the default and the off-variation settings.',
      }),
    },
    evaluation: {
      last30Days: defineMessage({
        id: 'feature.last30Days',
        defaultMessage: 'Last 30 days',
      }),
      last14Days: defineMessage({
        id: 'feature.last14Days',
        defaultMessage: 'Last 14 days',
      }),
      last7Days: defineMessage({
        id: 'feature.last7Days',
        defaultMessage: 'Last 7 days',
      }),
      last24Hours: defineMessage({
        id: 'feature.last24Hours',
        defaultMessage: 'Last 24 hours',
      }),
    },
    confirm: {
      title: defineMessage({
        id: 'feature.confirm.title',
        defaultMessage: 'Confirmation required',
      }),
      description: defineMessage({
        id: 'feature.confirm.description',
        defaultMessage:
          'This will make changes to the flag and increment the version.',
      }),
      archiveTitle: defineMessage({
        id: 'feature.confirm.archiveTitle',
        defaultMessage: 'Archive Feature Flag',
      }),
      archiveDescription: defineMessage({
        id: 'feature.confirm.archiveDescription',
        defaultMessage:
          'This will archive and return the default value defined in your code for all users. We recommend removing the code references to "{featureId}" from your application before archiving.',
      }),
      unarchiveTitle: defineMessage({
        id: 'feature.confirm.unarchiveTitle',
        defaultMessage: 'Unarchive Feature Flag',
      }),
      unarchiveDescription: defineMessage({
        id: 'feature.confirm.unarchiveDescription',
        defaultMessage:
          'Are you sure you want to unarchive the feature flag "{featureId}"?',
      }),
      flagUsedAsPrerequisite: defineMessage({
        id: 'feature.confirm.flagUsedAsPrerequisite',
        defaultMessage: `You can't archive while other flags use this flag as a
          prerequisite.`,
      }),
      flagIsActive: defineMessage({
        id: 'feature.confirm.flagIsActive',
        defaultMessage: 'It is receiving one more requests in the last 7 days.',
      }),
      enableNow: defineMessage({
        id: 'feature.confirm.enableNow',
        defaultMessage: 'Enable now',
      }),
      disableNow: defineMessage({
        id: 'feature.confirm.disableNow',
        defaultMessage: 'Disable now',
      }),
      schedule: defineMessage({
        id: 'feature.confirm.schedule',
        defaultMessage: 'Schedule',
      }),
      selectDate: defineMessage({
        id: 'feature.confirm.selectDate',
        defaultMessage: 'Select date',
      }),
      scheduleInfo: defineMessage({
        id: 'feature.confirm.scheduleInfo',
        defaultMessage:
          'You can update or delete the schedule on the Auto Operations tab on the Feature Flag details page.',
      }),
    },
    list: {
      header: {
        title: defineMessage({
          id: 'feature.list.header.title',
          defaultMessage: 'Feature Flags',
        }),
        description: defineMessage({
          id: 'feature.list.header.description',
          defaultMessage:
            'Use this page to see all feature flags in this project. Select a flag to manage the environment-specific targeting and rollout rules.',
        }),
      },
      active: defineMessage({
        id: 'feature.list.active',
        defaultMessage: 'Active',
      }),
      archive: defineMessage({
        id: 'feature.list.archive',
        defaultMessage: 'Archive',
      }),
      noData: {
        description: defineMessage({
          id: 'feature.list.noData.description',
          defaultMessage:
            'Create feature flags to manage who sees your features.',
        }),
      },
      noResult: {
        searchKeyword: defineMessage({
          id: 'feature.list.noResult.searchKeyword',
          defaultMessage: 'ID, name and description',
        }),
      },
    },
    filter: {
      maintainer: defineMessage({
        id: 'feature.filter.maintainer',
        defaultMessage: 'Maintainer',
      }),
      hasExperiment: defineMessage({
        id: 'feature.filter.hasExperiment',
        defaultMessage: 'Has experiment',
      }),
      enabled: defineMessage({
        id: 'feature.filter.enabled',
        defaultMessage: 'Enabled',
      }),
      archived: defineMessage({
        id: 'feature.filter.archived',
        defaultMessage: 'Archived',
      }),
      tags: defineMessage({
        id: 'feature.filter.tags',
        defaultMessage: 'Tags',
      }),
      tagsPlaceholder: defineMessage({
        id: 'feature.filter.tags.placeholder',
        defaultMessage: 'Select one or more tags',
      }),
      hasPrerequisites: defineMessage({
        id: 'feature.filter.hasPrerequisites',
        defaultMessage: 'Has prerequisites',
      }),
    },
    sort: {
      nameAz: defineMessage({
        id: 'feature.sort.nameAz',
        defaultMessage: 'Name A-Z',
      }),
      nameZa: defineMessage({
        id: 'feature.sort.nameZa',
        defaultMessage: 'Name Z-A',
      }),
      newest: defineMessage({
        id: 'feature.sort.newest',
        defaultMessage: 'Newest',
      }),
      oldest: defineMessage({
        id: 'feature.sort.oldest',
        defaultMessage: 'Oldest',
      }),
    },
    tab: {
      autoOps: defineMessage({
        id: 'feature.tab.autoOps',
        defaultMessage: 'Auto Operations',
      }),
      triggers: defineMessage({
        id: 'feature.tab.triggers',
        defaultMessage: 'Triggers',
      }),
      evaluation: defineMessage({
        id: 'feature.tab.evaluation',
        defaultMessage: 'Evaluation',
      }),
      experiments: defineMessage({
        id: 'feature.tab.experiments',
        defaultMessage: 'Experiments',
      }),
      history: defineMessage({
        id: 'feature.tab.history',
        defaultMessage: 'History',
      }),
      settings: defineMessage({
        id: 'feature.tab.settings',
        defaultMessage: 'Settings',
      }),
      targeting: defineMessage({
        id: 'feature.tab.targeting',
        defaultMessage: 'Targeting',
      }),
      variations: defineMessage({
        id: 'feature.tab.variations',
        defaultMessage: 'Variations',
      }),
    },
    search: {
      placeholder: defineMessage({
        id: 'feature.search.placeholder',
        defaultMessage: 'ID, Name, Description',
      }),
    },
    add: {
      header: {
        title: defineMessage({
          id: 'feature.add.header.title',
          defaultMessage: 'Create a feature flag',
        }),
        description: defineMessage({
          id: 'feature.add.header.description',
          defaultMessage:
            'Get started by filling in the information below to create your new feature flag.',
        }),
      },
    },
    clone: {
      header: {
        title: defineMessage({
          id: 'feature.clone.header.title',
          defaultMessage: 'Clone feature flag',
        }),
        description: defineMessage({
          id: 'feature.clone.header.description',
          defaultMessage:
            'It will copy the full targeting configuration, including on/off variation from the original flag to the new flag.',
        }),
      },
    },
    successMessages: {
      schedule: defineMessage({
        id: 'successMessages.schedule',
        defaultMessage: 'Schedule has been configured',
      }),
      flagEnabled: defineMessage({
        id: 'successMessages.flagEnabled',
        defaultMessage: 'Flag has been enabled',
      }),
      flagDisabled: defineMessage({
        id: 'successMessages.flagDisabled',
        defaultMessage: 'Flag has been disabled',
      }),
    },
  },
  input: {
    originEnvironment: defineMessage({
      id: 'input.originEnvironment',
      defaultMessage: 'Origin environment',
    }),
    destinationEnvironment: defineMessage({
      id: 'input.destinationEnvironment',
      defaultMessage: 'Destination environment',
    }),
    error: {
      maxLength: defineMessage({
        id: 'input.error.maxLength',
        defaultMessage:
          'The maximum length for this field is {max} characters.',
      }),
      minSelectOptionLength: defineMessage({
        id: 'input.error.minSelectOptionLength',
        defaultMessage: 'Must select at least one option.',
      }),
      invalidEmailAddress: defineMessage({
        id: 'input.error.invalidEmailAddress',
        defaultMessage: 'Invalid email address.',
      }),
      invalidEmailDomain: defineMessage({
        id: 'input.error.invalidEmailDomain',
        defaultMessage: 'Invalid email domain.',
      }),
      invalidId: defineMessage({
        id: 'input.error.invalidId',
        defaultMessage:
          "Invalid ID. ID must only contain lowercase letters, numbers or '-', and must start with an alphanumeric.",
      }),
      invalidName: defineMessage({
        id: 'input.error.invalidName',
        defaultMessage:
          "Invalid name. Name must only contain lowercase letters, numbers or '-', and must start with an alphanumeric.",
      }),
      invalidUrlCode: defineMessage({
        id: 'input.error.invalidUrlCode',
        defaultMessage:
          "Invalid URL code. URL code must only contain lowercase letters, numbers or '-', and must start with an alphanumeric.",
      }),
      required: defineMessage({
        id: 'input.error.required',
        defaultMessage: 'This is required.',
      }),
      mustBeUnique: defineMessage({
        id: 'input.error.mustBeUnique',
        defaultMessage: 'This must be unique.',
      }),
      notNumber: defineMessage({
        id: 'input.error.notNumber',
        defaultMessage: 'This must be a number.',
      }),
      not100Percentage: defineMessage({
        id: 'input.error.not100Percentage',
        defaultMessage: 'Total should be 100%.',
      }),
      notJson: defineMessage({
        id: 'input.error.notJson',
        defaultMessage: 'Invalid JSON.',
      }),
      notLaterThanCurrentTime: defineMessage({
        id: 'input.error.notLaterThanCurrentTime',
        defaultMessage: 'This must be later than the current time.',
      }),
      notLaterThanOrEqualDays: defineMessage({
        id: 'input.error.notLaterThanOrEqualDays',
        defaultMessage: 'This must be later than or equal to {days} days ago.',
      }),
      notLaterThanStartAt: defineMessage({
        id: 'input.error.notLaterThanStartAt',
        defaultMessage: 'This must be later than the start at.',
      }),
      notLessThanOrEquals30Days: defineMessage({
        id: 'input.error.notLessThanOrEquals30Days',
        defaultMessage: 'The period must be less than or equals to 30 days.',
      }),
    },
    optional: defineMessage({
      id: 'input.optional',
      defaultMessage: '(optional)',
    }),
    name: defineMessage({
      id: 'input.name',
      defaultMessage: 'Name',
    }),
    email: defineMessage({
      id: 'input.email',
      defaultMessage: 'Email',
    }),
    featureFlag: defineMessage({
      id: 'input.featureFlag',
      defaultMessage: 'Feature Flag',
    }),
    projectId: defineMessage({
      id: 'input.projectId',
      defaultMessage: 'Project ID',
    }),
    role: defineMessage({
      id: 'input.role',
      defaultMessage: 'Role',
    }),
  },
  id: defineMessage({
    id: 'id',
    defaultMessage: 'ID',
  }),
  noData: {
    title: defineMessage({
      id: 'noData.title',
      defaultMessage: 'There are no {title} yet.',
    }),
  },
  noResult: {
    title: defineMessage({
      id: 'noResult.title',
      defaultMessage: 'No {title} match. You can try this:',
    }),
    searchByKeyword: defineMessage({
      id: 'noResult.searchByKeyword',
      defaultMessage: 'Search by {keyword}',
    }),
    changeFilterSelection: defineMessage({
      id: 'noResult.changeFilterSelection',
      defaultMessage: 'Change your filter selection',
    }),
    checkTypos: defineMessage({
      id: 'noResult.checkTypos',
      defaultMessage: 'Check for typos',
    }),
    dateRange: {
      title: defineMessage({
        id: 'noResult.dateRange.title',
        defaultMessage: 'No entries',
      }),
      description: defineMessage({
        id: 'noResult.dateRange.description',
        defaultMessage:
          'There are no entries for these dates. Please choose a different date and try again.',
      }),
    },
  },
  name: defineMessage({
    id: 'name',
    defaultMessage: 'Name',
  }),
  notification: {
    add: {
      header: {
        title: defineMessage({
          id: 'notification.add.header.title',
          defaultMessage: 'Create a notification',
        }),
        description: defineMessage({
          id: 'notification.add.header.description',
          defaultMessage:
            'A notification lets you know when someone adds or updates something on the admin console and operational tasks status.',
        }),
      },
    },
    confirm: {
      deleteTitle: defineMessage({
        id: 'notification.confirm.delete.title',
        defaultMessage: 'Delete notification',
      }),
      deleteDescription: defineMessage({
        id: 'notification.confirm.delete.description',
        defaultMessage:
          'The {notificationName} notification will be deleted permanently.',
      }),
      enableTitle: defineMessage({
        id: 'notification.confirm.enable.title',
        defaultMessage: 'Enable notification',
      }),
      enableDescription: defineMessage({
        id: 'notification.confirm.enable.description',
        defaultMessage:
          'Are you sure you want to enable the {notificationName} notification?',
      }),
      disableTitle: defineMessage({
        id: 'notification.confirm.disable.title',
        defaultMessage: 'Disable notification',
      }),
      disableDescription: defineMessage({
        id: 'notification.confirm.disable.description',
        defaultMessage:
          'Are you sure you want to disable the {notificationName} notification?',
      }),
    },
    filter: {
      enabled: defineMessage({
        id: 'notification.filter.enabled',
        defaultMessage: 'Enabled',
      }),
    },
    filterOptions: {
      enabled: defineMessage({
        id: 'notification.filterOptions.enabled',
        defaultMessage: 'Enabled',
      }),
      disabled: defineMessage({
        id: 'notification.filterOptions.disabled',
        defaultMessage: 'Disabled',
      }),
    },
    list: {
      header: {
        title: defineMessage({
          id: 'notification.list.header.title',
          defaultMessage: 'Notification',
        }),
        description: defineMessage({
          id: 'notification.list.header.description',
          defaultMessage:
            'Select a notification to manage the settings or click on the Add button to add a new one.',
        }),
      },
      noData: {
        description: defineMessage({
          id: 'notification.list.noData.description',
          defaultMessage:
            'You can receive notifications when operations such as additions and changes are made on the admin console, or the status of operational tasks.',
        }),
      },
      noResult: {
        searchKeyword: defineMessage({
          id: 'notification.list.noResult.searchKeyword',
          defaultMessage: 'Name',
        }),
      },
    },
    slackIncomingWebhookUrl: defineMessage({
      id: 'notification.slackIncomingWebhookUrl',
      defaultMessage: 'Slack incoming webhook URL',
    }),
    search: {
      placeholder: defineMessage({
        id: 'notification.search.placeholder',
        defaultMessage: 'Name',
      }),
    },
    sort: {
      nameAz: defineMessage({
        id: 'notification.sort.nameAz',
        defaultMessage: 'Name A-Z',
      }),
      nameZa: defineMessage({
        id: 'notification.sort.nameZa',
        defaultMessage: 'Name Z-A',
      }),
      newest: defineMessage({
        id: 'notification.sort.newest',
        defaultMessage: 'Newest',
      }),
      oldest: defineMessage({
        id: 'notification.sort.oldest',
        defaultMessage: 'Oldest',
      }),
    },
    update: {
      header: {
        title: defineMessage({
          id: 'notification.update.header.title',
          defaultMessage: 'Update the notification',
        }),
        description: defineMessage({
          id: 'notification.update.header.description',
          defaultMessage:
            'A notification lets you know when someone adds or updates something on the admin console and operational tasks status.',
        }),
      },
    },
  },
  segment: {
    action: {
      download: defineMessage({
        id: 'segment.action.download',
        defaultMessage: 'Download user list',
      }),
      delete: defineMessage({
        id: 'segment.action.delete',
        defaultMessage: 'Delete segment',
      }),
    },
    confirm: {
      deleteTitle: defineMessage({
        id: 'segment.confirm.delete.title',
        defaultMessage: 'Delete segment',
      }),
      deleteDescription: defineMessage({
        id: 'segment.confirm.delete.description',
        defaultMessage:
          'The {segmentName} segment will be deleted permanently.',
      }),
      cannotDelete: defineMessage({
        id: 'segment.confirm.delete.cannotDelete',
        defaultMessage:
          "The {segmentName} segment can't be deleted because {length} {length, plural, one {flag is} other {flags are}} using it.",
      }),
    },
    filter: {
      status: defineMessage({
        id: 'segment.filter.status',
        defaultMessage: 'Status of use',
      }),
    },
    filterOptions: {
      inUse: defineMessage({
        id: 'segment.filterOptions.inUse',
        defaultMessage: 'In use',
      }),
      notInUse: defineMessage({
        id: 'segment.filterOptions.notInUse',
        defaultMessage: 'Not in use',
      }),
    },
    add: {
      header: {
        title: defineMessage({
          id: 'segment.add.header.title',
          defaultMessage: 'Create a segment',
        }),
        description: defineMessage({
          id: 'segment.add.header.description',
          defaultMessage:
            'User segment allows you to manage all user targets for a single feature flag variation. You can use it to make changes to a large number of users or to test beta features on a small number of users.',
        }),
      },
    },
    update: {
      header: {
        title: defineMessage({
          id: 'segment.update.header.title',
          defaultMessage: 'Update the segment',
        }),
        description: defineMessage({
          id: 'segment.update.header.description',
          defaultMessage:
            'User segment allows you to manage all user targets for a single feature flag variation. You can use it to make changes to a large number of users or to test beta features on a small number of users.',
        }),
      },
      userId: defineMessage({
        id: 'segment.update.userId',
        defaultMessage:
          "The user ID list can't be updated because {length} {length, plural, one {flag is} other {flags are}} using it. Remove the segment from the flag before updating it.",
      }),
    },
    uploading: {
      title: defineMessage({
        id: 'segment.uploading.title',
        defaultMessage: 'Upload in progress',
      }),
      message: defineMessage({
        id: 'segment.uploading.message',
        defaultMessage:
          "Segments can't be updated until the user list has been uploaded.",
      }),
    },
    list: {
      header: {
        title: defineMessage({
          id: 'segment.list.header.title',
          defaultMessage: 'Segments',
        }),
        description: defineMessage({
          id: 'segment.list.header.description',
          defaultMessage:
            'On this page, you can check all segments for this environment. Select a segment to manage the settings or click on the Add button to add a new one.',
        }),
      },
      noData: {
        description: defineMessage({
          id: 'segment.list.noData.description',
          defaultMessage:
            'You can create a user segment to manage all user targets for a single feature flag variation.',
        }),
      },
      noResult: {
        searchKeyword: defineMessage({
          id: 'segment.list.noResult.searchKeyword',
          defaultMessage: 'Name and description',
        }),
      },
    },
    select: {
      noData: {
        description: defineMessage({
          id: 'segment.select.noData.description',
          defaultMessage:
            'Please add user segments on the user segment list page.',
        }),
      },
    },
    fileUpload: {
      browseFiles: defineMessage({
        id: 'segment.fileUpload.browseFiles',
        defaultMessage: 'Browse files',
      }),
      fileFormat: defineMessage({
        id: 'segment.fileUpload.fileFormat',
        defaultMessage: 'Accepted file type: .csv and .txt (Max size: 2MB)',
      }),
      userList: defineMessage({
        id: 'segment.fileUpload.userList',
        defaultMessage: 'List of user IDs',
      }),
      fileMaxSize: defineMessage({
        id: 'segment.fileUpload.fileMaxSize',
        defaultMessage: 'The maximum size of the file is 1MB',
      }),
      unsupportedType: defineMessage({
        id: 'segment.fileUpload.unsupportedType',
        defaultMessage: 'The file format is not supported',
      }),
      fileSize: defineMessage({
        id: 'segment.fileUpload.fileSize',
        defaultMessage: '{fileSize} bytes',
      }),
      uploadInProgress: defineMessage({
        id: 'segment.fileUpload.uploadInProgress',
        defaultMessage: 'The file cannot be updated due to upload in progress',
      }),
    },
    search: {
      placeholder: defineMessage({
        id: 'segment.search.placeholder',
        defaultMessage: 'Name and Description',
      }),
    },
    sort: {
      nameAz: defineMessage({
        id: 'segment.sort.nameAz',
        defaultMessage: 'Name A-Z',
      }),
      nameZa: defineMessage({
        id: 'segment.sort.nameZa',
        defaultMessage: 'Name Z-A',
      }),
      newest: defineMessage({
        id: 'segment.sort.newest',
        defaultMessage: 'Newest',
      }),
      oldest: defineMessage({
        id: 'segment.sort.oldest',
        defaultMessage: 'Oldest',
      }),
    },
    status: {
      uploading: defineMessage({
        id: 'segment.status.uploading',
        defaultMessage: 'UPLOADING',
      }),
      uploadFailed: defineMessage({
        id: 'segment.status.uploadFailed',
        defaultMessage: 'UPLOAD FAILED',
      }),
    },
    userCount: defineMessage({
      id: 'segment.userCount',
      defaultMessage: 'users',
    }),
    enterUserIdsPlaceholder: defineMessage({
      id: 'segment.enterUserIdsPlaceholder',
      defaultMessage:
        'Enter IDs separated by commas (E.g., userId1, userId2, userId3)',
    }),
  },
  push: {
    add: {
      header: {
        title: defineMessage({
          id: 'push.add.header.title',
          defaultMessage: 'Create a push',
        }),
        description: defineMessage({
          id: 'push.add.header.description',
          defaultMessage:
            'By using the push feature, the SDK can be updated in real-time. Every time a feature flag is updated, a notification is sent to the client.',
        }),
      },
    },
    confirm: {
      deleteTitle: defineMessage({
        id: 'push.confirm.delete.title',
        defaultMessage: 'Delete push',
      }),
      deleteDescription: defineMessage({
        id: 'push.confirm.delete.description',
        defaultMessage: 'The {pushName} push will be deleted permanently.',
      }),
    },
    input: {
      fcmApiKey: defineMessage({
        id: 'push.input.fcmApiKey',
        defaultMessage: 'Firebase Cloud Messaging API Key',
      }),
    },
    list: {
      header: {
        description: defineMessage({
          id: 'push.list.header.description',
          defaultMessage:
            'Select a push to manage the settings or click on the Add button to add a new one.',
        }),
        title: defineMessage({
          id: 'push.list.header.title',
          defaultMessage: 'Push',
        }),
      },
      noData: {
        description: defineMessage({
          id: 'push.list.noData.description',
          defaultMessage:
            'You can create a push to update the SDK client in real-time. Every time a feature flag is updated, a notification is sent to the client.',
        }),
      },
      noResult: {
        searchKeyword: defineMessage({
          id: 'push.list.noResult.searchKeyword',
          defaultMessage: 'Name',
        }),
      },
    },
    search: {
      placeholder: defineMessage({
        id: 'push.search.placeholder',
        defaultMessage: 'Name',
      }),
    },
    sort: {
      nameAz: defineMessage({
        id: 'push.sort.nameAz',
        defaultMessage: 'Name A-Z',
      }),
      nameZa: defineMessage({
        id: 'push.sort.nameZa',
        defaultMessage: 'Name Z-A',
      }),
      newest: defineMessage({
        id: 'push.sort.newest',
        defaultMessage: 'Newest',
      }),
      oldest: defineMessage({
        id: 'push.sort.oldest',
        defaultMessage: 'Oldest',
      }),
    },
    update: {
      header: {
        title: defineMessage({
          id: 'push.update.header.title',
          defaultMessage: 'Update the push',
        }),
        description: defineMessage({
          id: 'push.update.header.description',
          defaultMessage:
            'By using the push feature, the SDK can be updated in real-time. Every time a feature flag is updated, a notification is sent to the client.',
        }),
      },
    },
  },
  settings: {
    list: {
      header: {
        title: defineMessage({
          id: 'settings.list.header.title',
          defaultMessage: 'Settings',
        }),
        description: defineMessage({
          id: 'settings.list.header.description',
          defaultMessage:
            'On this page, you can check all settings for this environment. Select a tab to manage the settings.',
        }),
      },
    },
    tab: {
      pushes: defineMessage({
        id: 'settings.tab.pushes',
        defaultMessage: 'Pushes',
      }),
      notifications: defineMessage({
        id: 'settings.tab.notifications',
        defaultMessage: 'Notifications',
      }),
      webhooks: defineMessage({
        id: 'settings.tab.webhooks',
        defaultMessage: 'Webhooks',
      }),
    },
  },
  sideMenu: {
    adminSettings: defineMessage({
      id: 'sideMenu.adminSettings',
      defaultMessage: 'Admin Settings',
    }),
    featureFlags: defineMessage({
      id: 'sideMenu.featureFlags',
      defaultMessage: 'Feature Flags',
    }),
    experiments: defineMessage({
      id: 'sideMenu.experiments',
      defaultMessage: 'Experiments',
    }),
    goals: defineMessage({
      id: 'sideMenu.goals',
      defaultMessage: 'Goals',
    }),
    apiKeys: defineMessage({
      id: 'sideMenu.apiKeys',
      defaultMessage: 'API Keys',
    }),
    userSegments: defineMessage({
      id: 'sideMenu.userSegments',
      defaultMessage: 'User Segments',
    }),
    user: defineMessage({
      id: 'sideMenu.user',
      defaultMessage: 'Users',
    }),
    auditLog: defineMessage({
      id: 'sideMenu.auditLog',
      defaultMessage: 'Audit Logs',
    }),
    accounts: defineMessage({
      id: 'sideMenu.accounts',
      defaultMessage: 'Accounts',
    }),
    documentation: defineMessage({
      id: 'sideMenu.documentation',
      defaultMessage: 'Documentation',
    }),
    settings: defineMessage({
      id: 'sideMenu.settings',
      defaultMessage: 'Settings',
    }),
    logout: defineMessage({
      id: 'sideMenu.logout',
      defaultMessage: 'Logout',
    }),
  },
  tags: defineMessage({
    id: 'tags',
    defaultMessage: 'Tags',
  }),
  readMore: defineMessage({
    id: 'readMore',
    defaultMessage: 'Read more',
  }),
  reason: {
    reason: defineMessage({
      id: 'reason.reason',
      defaultMessage: 'Evaluation reason',
    }),
    target: defineMessage({
      id: 'reason.target',
      defaultMessage: 'Target',
    }),
    rule: defineMessage({
      id: 'reason.rule',
      defaultMessage: 'Rule',
    }),
    client: defineMessage({
      id: 'reason.client',
      defaultMessage: 'Client',
    }),
    offVariation: defineMessage({
      id: 'reason.offVariation',
      defaultMessage: 'Off variation',
    }),
  },
  sourceType: {
    account: defineMessage({
      id: 'sourceType.account',
      defaultMessage: 'Account',
    }),
    accountDescription: defineMessage({
      id: 'sourceType.accountDescription',
      defaultMessage: 'Get notified when someone adds or updates an account',
    }),
    adminSubscription: defineMessage({
      id: 'sourceType.adminSubscription',
      defaultMessage: 'Subscription',
    }),
    adminNotification: defineMessage({
      id: 'sourceType.adminNotification',
      defaultMessage: 'Notification',
    }),
    adminNotificationDescription: defineMessage({
      id: 'sourceType.adminNotificationDescription',
      defaultMessage:
        'Get notified when someone adds or updates a notification',
    }),
    apiKey: defineMessage({
      id: 'sourceType.apiKey',
      defaultMessage: 'API Key',
    }),
    apiKeyDescription: defineMessage({
      id: 'sourceType.apiKeyDescription',
      defaultMessage: 'Get notified when someone adds or updates an API Key',
    }),
    autoOps: defineMessage({
      id: 'sourceType.autoOps',
      defaultMessage: 'Auto-Ops',
    }),
    autoOpsDescription: defineMessage({
      id: 'sourceType.autoOpsDescription',
      defaultMessage: 'Get notified when the Auto-Ops is triggered',
    }),
    autoOperation: defineMessage({
      id: 'sourceType.autoOperation',
      defaultMessage: 'Auto Operation',
    }),
    progressiveRollout: defineMessage({
      id: 'sourceType.progressiveRollout',
      defaultMessage: 'Progressive Rollout',
    }),
    subscription: defineMessage({
      id: 'sourceType.subscription',
      defaultMessage: 'Subscription',
    }),
    environment: defineMessage({
      id: 'sourceType.environment',
      defaultMessage: 'Environment',
    }),
    environmentDescription: defineMessage({
      id: 'sourceType.environmentDescription',
      defaultMessage:
        'Get notified when someone adds or updates an environment',
    }),
    experiment: defineMessage({
      id: 'sourceType.experiment',
      defaultMessage: 'Experiment',
    }),
    experimentDescription: defineMessage({
      id: 'sourceType.experimentDescription',
      defaultMessage: 'Get notified when someone adds or updates an experiment',
    }),
    featureFlag: defineMessage({
      id: 'sourceType.featureFlag',
      defaultMessage: 'Feature Flag',
    }),
    featureFlagDescription: defineMessage({
      id: 'sourceType.featureFlagDescription',
      defaultMessage:
        'Get notified when someone adds or updates a feature flag',
    }),
    goal: defineMessage({
      id: 'sourceType.goal',
      defaultMessage: 'Goal',
    }),
    goalDescription: defineMessage({
      id: 'sourceType.goalDescription',
      defaultMessage: 'Get notified when someone adds or updates a goal',
    }),
    mauCount: defineMessage({
      id: 'sourceType.mauCount',
      defaultMessage: 'Monthly Active Users count',
    }),
    mauCountDescription: defineMessage({
      id: 'sourceType.mauCountDescription',
      defaultMessage:
        'Get notified the monthly active users count on the first day of every month',
    }),
    notification: defineMessage({
      id: 'sourceType.notification',
      defaultMessage: 'Notification',
    }),
    notificationDescription: defineMessage({
      id: 'sourceType.notificationDescription',
      defaultMessage:
        'Get notified when someone adds or updates a notification',
    }),
    project: defineMessage({
      id: 'sourceType.project',
      defaultMessage: 'Project',
    }),
    projectDescription: defineMessage({
      id: 'sourceType.projectDescription',
      defaultMessage: 'Get notified when someone adds or updates a project',
    }),
    push: defineMessage({
      id: 'sourceType.push',
      defaultMessage: 'Push',
    }),
    pushDescription: defineMessage({
      id: 'sourceType.pushDescription',
      defaultMessage: 'Get notified when someone adds or updates a push',
    }),
    runningExperiments: defineMessage({
      id: 'sourceType.runningExperiments',
      defaultMessage: 'Running Experiments',
    }),
    runningExperimentsDescription: defineMessage({
      id: 'sourceType.runningExperimentsDescription',
      defaultMessage: 'Get notified daily of the list of running experiments',
    }),
    segment: defineMessage({
      id: 'sourceType.segment',
      defaultMessage: 'User Segment',
    }),
    segmentDescription: defineMessage({
      id: 'sourceType.segmentDescription',
      defaultMessage:
        'Get notified when someone adds or updates a user segment',
    }),
    staleFeatureFlag: defineMessage({
      id: 'sourceType.staleFeatureFlag',
      defaultMessage: 'Stale feature flag',
    }),
    staleFeatureFlagDescription: defineMessage({
      id: 'sourceType.staleFeatureFlagDescription',
      defaultMessage: 'Get notified when a feature flag becomes stale',
    }),
  },
  type: defineMessage({
    id: 'type',
    defaultMessage: 'Type',
  }),
  urlCode: defineMessage({
    id: 'urlCode',
    defaultMessage: 'URL code',
  }),
  fullStop: defineMessage({
    id: 'fullStop',
    defaultMessage: '.',
  }),
  select: defineMessage({
    id: 'select',
    defaultMessage: 'Select',
  }),
};
