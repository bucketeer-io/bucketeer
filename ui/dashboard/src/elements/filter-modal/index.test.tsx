import { ReactNode } from 'react';
import { renderToStaticMarkup } from 'react-dom/server';
import { useAccountsLoader } from 'hooks/use-accounts-loading-more';
import { FilterTypes } from 'hooks/use-options';
import { describe, expect, it, vi } from 'vitest';
import { apiKeyFilterConfig } from 'pages/api-keys/api-key-modal/filter-api-key-modal/config';
import { flagFilterConfig } from 'pages/feature-flags/flags-modal/filter-flag-modal/config';
import { memberFilterConfig } from 'pages/members/member-modal/filter-member-modal/config';
import { notificationFilterConfig } from 'pages/notifications/notification-modal/filter-notification-modal/config';
import { organizationFilterConfig } from 'pages/organizations/organization-modal/filter-organization-modal/config';
import { projectFilterConfig } from 'pages/projects/project-modal/filter-project-modal/config';
import { pushFilterConfig } from 'pages/pushes/push-modal/filter-push-modal/config';
import { userSegmentFilterConfig } from 'pages/user-segments/user-segment-modal/filter-segment-modal/config';
import FilterModal from '.';

vi.mock('@icons', () => ({ IconPlus: () => null, IconTrash: () => null }));
vi.mock('@queries/tags', () => ({ useQueryTags: vi.fn() }));
vi.mock('@queries/teams', () => ({ useQueryTeams: vi.fn() }));
vi.mock('@queries/environments', () => ({ useQueryEnvironments: vi.fn() }));
vi.mock('auth', () => ({
  useAuth: () => ({ consoleAccount: {} }),
  getCurrentEnvironment: () => ({ id: 'env', organizationId: 'org' }),
  getEditorEnvironments: vi.fn()
}));
vi.mock('hooks/use-accounts-loading-more', () => ({
  useAccountsLoader: vi.fn()
}));
vi.mock('i18n', () => ({
  i18n: { t: (key: string) => key },
  getLanguage: () => 'en',
  useTranslation: () => ({ t: (key: string) => key })
}));
vi.mock('components/modal/dialog', () => ({
  default: ({ children }: { children: ReactNode }) => (
    <section>{children}</section>
  )
}));
vi.mock('components/button', () => ({ default: 'button' }));
vi.mock('components/button-bar', () => ({ ButtonBar: () => null }));
vi.mock('components/divider', () => ({ default: () => null }));
vi.mock('components/icon', () => ({ default: () => null }));
vi.mock('components/dropdown', () => ({
  default: ({
    placeholder,
    options
  }: {
    placeholder: string;
    options: { value: string | number; label: string }[];
  }) => (
    <select aria-label={placeholder}>
      {options.map(option => (
        <option key={option.value} value={option.value}>
          {option.label}
        </option>
      ))}
    </select>
  )
}));
vi.mock('elements/dropdown-with-search', () => ({
  default: ({ label }: { label: string }) => <span>{label}</span>
}));

const modalProps = {
  isOpen: true,
  onSubmit: vi.fn(),
  onClose: vi.fn(),
  onClearFilters: vi.fn()
};

describe('filter hydration', () => {
  const fields = [
    ...flagFilterConfig.fields.filter(field => field.valueKind === 'boolean'),
    ...[
      apiKeyFilterConfig,
      memberFilterConfig,
      notificationFilterConfig,
      pushFilterConfig,
      organizationFilterConfig,
      projectFilterConfig
    ].map(config =>
      config.fields.find(field => field.type === FilterTypes.ENABLED)!
    ),
    userSegmentFilterConfig.fields[0]
  ];

  it.each(fields)('ignores empty URL values for $labelKey', field => {
    const key = Object.keys(field.toFilter(undefined))[0];
    for (const value of [undefined, null, '']) {
      expect(field.fromFilter({ [key]: value })).toBeUndefined();
    }
  });

  it.each(fields)('preserves explicit true and false for $labelKey', field => {
    const key = Object.keys(field.toFilter(undefined))[0];
    for (const value of [true, false]) {
      expect(field.toFilter(field.fromFilter({ [key]: value }))).toEqual({
        [key]: value
      });
    }
  });
});

describe('filter value rendering', () => {
  it('renders a newly preloaded maintainer name with unchanged email options', () => {
    const emailOptions = [
      { value: 'first@example.com', label: 'first@example.com' }
    ];
    const maintainer = flagFilterConfig.fields.find(
      field => field.type === FilterTypes.MAINTAINER
    )!;
    const config = {
      ...flagFilterConfig,
      fields: [{ ...maintainer, emptyValue: 'selected@example.com' }]
    };
    const loader = {
      accounts: [],
      emailOptions,
      isLoading: false,
      hasMore: false,
      isLoadingMore: false,
      isInitialLoading: false,
      isSearching: false,
      loadMore: vi.fn(),
      onSearchChange: vi.fn(),
      getAccountLabel: (email: string) => email
    };
    vi.mocked(useAccountsLoader).mockReturnValue(loader);
    expect(
      renderToStaticMarkup(<FilterModal {...modalProps} config={config} />)
    ).toContain('selected@example.com');

    vi.mocked(useAccountsLoader).mockReturnValue({
      ...loader,
      getAccountLabel: () => 'Selected Maintainer'
    });
    expect(
      renderToStaticMarkup(<FilterModal {...modalProps} config={config} />)
    ).toContain('Selected Maintainer');
  });

  it.each([
    organizationFilterConfig,
    projectFilterConfig,
    userSegmentFilterConfig
  ])('keeps the only field available in single mode', config => {
    const html = renderToStaticMarkup(
      <FilterModal {...modalProps} config={config} />
    );
    expect(html).toContain(`<option value="${config.fields[0].type}">`);
  });
});
