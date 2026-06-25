import { FilterOption, FilterTypes } from 'hooks/use-options';
import { DropdownOption, DropdownValue } from 'components/dropdown';

/**
 * How a field's value is edited in the UI.
 * - `boolean`     plain yes/no dropdown
 * - `enum`        single-select from a fixed option list
 * - `multiselect` multi-select from a fixed option list
 * - `searchable`  single/multi searchable dropdown backed by fetched options
 * - `searchable-paginated` DropdownMenuWithSearch backed by a paginated loader
 *                 (currently the maintainer/account picker)
 */
export type FilterValueKind =
  | 'boolean'
  | 'enum'
  | 'multiselect'
  | 'searchable'
  | 'searchable-paginated';

/**
 * Runtime value/state a field needs from a data source. Returned by a field's
 * `useData` hook so the generic modal does not need to know about queries.
 */
export interface FilterFieldData {
  /** Options to render in the value dropdown. */
  options: DropdownOption[];
  isLoading?: boolean;
  /** Resolve a stored value to its display label (e.g. email/tag/env name). */
  getLabel?: (value: FilterOption['filterValue']) => string;
  /** Pagination/search wiring, only used by `searchable-paginated`. */
  hasMore?: boolean;
  isLoadingMore?: boolean;
  isSearching?: boolean;
  loadMore?: () => void;
  onSearchChange?: (value: string) => void;
}

export interface FilterFieldContext {
  /** Whether this field is currently selected (gates data fetching). */
  enabled: boolean;
  /**
   * The field's current edited value, for data sources that need to preload the
   * selected entry (e.g. the maintainer account loader's `preloadEmails`).
   */
  value?: FilterOption['filterValue'];
}

/**
 * Definition of a single filter type a page supports.
 *
 * The page-specific serialization quirks live entirely in `toFilter` /
 * `fromFilter`, keeping the generic modal free of any page knowledge.
 */
export interface FilterFieldDef<F> {
  type: FilterTypes;
  /** Translation key for the filter-type label (namespace `common`). */
  labelKey: string;
  valueKind: FilterValueKind;

  /**
   * Lazy data source for the value dropdown. Called once per render with
   * whether the field is active, so the underlying query can be gated.
   * Fixed-option fields (boolean/enum/multiselect) can supply options here too.
   */
  useData?: (ctx: FilterFieldContext) => FilterFieldData;

  /** Initial empty value when this field is freshly selected. */
  emptyValue: FilterOption['filterValue'];

  /** Map this field's edited value INTO the page filters object on submit. */
  toFilter: (filterValue: FilterOption['filterValue']) => Partial<F>;
  /**
   * Read this field's value FROM the page filters object on open. Return
   * `undefined` when the filter is not present (field is not pre-selected).
   */
  fromFilter: (filters: Partial<F>) => FilterOption['filterValue'] | undefined;
}

export interface FilterModalConfig<F> {
  /** `single` hides the add-filter button and per-row remove icons. */
  mode: 'single' | 'multi';
  fields: FilterFieldDef<F>[];
  /**
   * Keys reset to `undefined` on submit before applying selected filters, so
   * de-selected filters are cleared. Defaults to each field's resulting keys.
   */
  defaultFilters?: Partial<F>;
  /**
   * Extra keys merged into every submitted filters object (e.g. experiments'
   * `isFilter: true` marker). Applied last, after default + selected filters.
   */
  submitExtra?: Partial<F>;
  /**
   * Gate for hydrating selected filters from `filters` on open. Return `false`
   * to skip hydration (e.g. experiments only hydrate when `isFilter` is set).
   * Defaults to always hydrate.
   */
  shouldHydrate?: (filters: Partial<F>) => boolean;
}

export type { DropdownOption, DropdownValue };
