import { readdirSync, readFileSync, statSync } from 'fs';
import { dirname, join, resolve } from 'path';
import { fileURLToPath } from 'url';
import { describe, expect, it } from 'vitest';
import { resolveInvalidationKeys, URL_TO_KEYS } from './cache-invalidation-map';

/**
 * The repo root for this dashboard package, computed relative to this file.
 * Tests resolve `@api` and `@queries` paths from here so the test stays robust
 * regardless of where vitest is invoked from. The package is ESM
 * (`"type": "module"`), so we derive the path from `import.meta.url` instead
 * of relying on the CommonJS `__dirname` shim.
 */
const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const SRC_ROOT = resolve(__dirname, '..');
const API_ROOT = resolve(SRC_ROOT, '@api');
const QUERIES_ROOT = resolve(SRC_ROOT, '@queries');

/**
 * Endpoints that are intentionally excluded from cache invalidation.
 *
 * Add a new entry here ONLY if the endpoint genuinely should not invalidate
 * any client-side cache (auth handshakes, read-only debug endpoints, etc.).
 * Anything that mutates user-visible data MUST be covered by a rule in
 * `URL_TO_KEYS` instead.
 */
const NON_INVALIDATING_ENDPOINTS: ReadonlySet<string> = new Set([
  // Auth handshakes – they mutate session state, not cached data.
  '/v1/auth/signin',
  '/v1/auth/refresh_token',
  '/v1/auth/exchange_token',
  '/v1/auth/authentication_url',
  '/v1/auth/switch_organization',
  '/v1/exchange_demo_token',
  // Debug/inspection endpoints – they don't change persisted state.
  '/v1/debug_evaluate_features'
]);

/**
 * Query keys defined in `src/@queries/*.ts` that are read-only and have no
 * known mutating endpoint, so they are NOT expected to appear in
 * `URL_TO_KEYS`.
 */
const NON_INVALIDATED_QUERY_KEYS: ReadonlySet<string> = new Set([
  'audit-log-details', // immutable historical data; once fetched, never changes
  'ai-chat-suggestions', // server-side AI suggestions; not user-mutable
  'demo-site-status', // server-side feature flag; polled, not mutated
  'insights-monthly-summary',
  'insights-latency',
  'insights-requests',
  'insights-evaluations',
  'insights-error-rates'
]);

/**
 * Files under `src/@api/` that are allowed to use a non-axios HTTP transport
 * (native `fetch`, `XMLHttpRequest`, etc.). Such requests do NOT pass through
 * the cache-invalidation interceptor, so each entry must be a read-only
 * endpoint or have a documented reason to be exempt.
 *
 * Add a new entry here ONLY after confirming the endpoint does not need to
 * invalidate any client-side cache. Anything that mutates user-visible data
 * MUST be migrated to `@api/axios-client` (or another axios instance with the
 * interceptor installed) instead.
 */
const NON_AXIOS_TRANSPORT_FILES: ReadonlySet<string> = new Set([
  // POST /v1/aichat/chat – uses native fetch for SSE/ReadableStream support
  // (axios cannot consume SSE streams). The response is a streamed AI
  // suggestion, not a persisted mutation; the corresponding query key
  // `ai-chat-suggestions` is in NON_INVALIDATED_QUERY_KEYS above.
  'src/@api/ai-chat/chat-streamer.ts'
]);

/**
 * Matches HTTP-method calls (`.get(`, `.post(`, etc.) on an axios chain. We
 * deliberately match `.method(` rather than `axiosClient.method(` because the
 * codebase formats requests as multi-line chains, e.g.
 *
 *     return axiosClient
 *       .post('/v1/foo', body)
 *       .then(...);
 *
 * To prevent false positives from unrelated objects (e.g. `array.delete()` or
 * `Map.get()`), the scan is gated on `AXIOS_CLIENT_IMPORT_REGEX` below: only
 * files that actually import the shared `axiosClient` — or create their own
 * instance and install the cache-invalidation interceptor on it — are
 * considered. Both cases route through the same invalidation logic.
 */
const AXIOS_CALL_METHOD_REGEX =
  /\.(get|post|put|patch|delete)\b(?:<[^>]*>)?\s*\(/g;
const AXIOS_CLIENT_IMPORT_REGEX =
  /from\s+['"](?:@api\/axios-client|\.{1,2}\/(?:[^'"]+\/)?axios-client)['"]|\binstallCacheInvalidationInterceptor\s*\(/;
const URL_LITERAL_REGEX = /^\s*[`'"]([^`'"]+?)[`'"]/;

type ExtractedCall = {
  method: 'get' | 'post' | 'put' | 'patch' | 'delete';
  url: string;
  file: string;
};

const walk = (dir: string, files: string[] = []): string[] => {
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry);
    const stat = statSync(full);
    if (stat.isDirectory()) {
      walk(full, files);
    } else if (stat.isFile() && /\.(ts|tsx)$/.test(entry)) {
      files.push(full);
    }
  }
  return files;
};

const extractAxiosCalls = (file: string): ExtractedCall[] => {
  const source = readFileSync(file, 'utf8');
  // Skip files that don't actually use the shared axios instance. Only calls
  // routed through `@api/axios-client` go through the response interceptor
  // that triggers cache invalidation, so any `.get/.post(` matches in other
  // files are irrelevant (and could be false positives on unrelated objects).
  if (!AXIOS_CLIENT_IMPORT_REGEX.test(source)) return [];
  const calls: ExtractedCall[] = [];

  for (const match of source.matchAll(AXIOS_CALL_METHOD_REGEX)) {
    const method = match[1] as ExtractedCall['method'];
    const after = source.slice(match.index! + match[0].length);
    const urlMatch = URL_LITERAL_REGEX.exec(after);
    if (!urlMatch) continue;
    let url = urlMatch[1];
    // Strip query string and template-literal interpolation tails so we test
    // the path shape against URL_TO_KEYS' regexes.
    url = url.replace(/\$\{[^}]+\}/g, '').split('?')[0];
    calls.push({ method, url, file });
  }

  return calls;
};

const collectAllApiCalls = (): ExtractedCall[] => {
  const files = walk(API_ROOT).filter(
    f =>
      !f.endsWith('cache-invalidation-map.ts') &&
      !f.endsWith('cache-invalidation-map.test.ts') &&
      !f.endsWith('axios-client.ts')
  );
  return files.flatMap(extractAxiosCalls);
};

const QUERY_KEY_REGEX =
  /export const [A-Z][A-Z0-9_]*\s*(?::\s*[^=]+)?=\s*[`'"]([^`'"]+)[`'"]/g;

const collectAllQueryKeys = (): {
  key: string;
  file: string;
}[] => {
  const files = walk(QUERIES_ROOT);
  const keys: { key: string; file: string }[] = [];
  for (const file of files) {
    const source = readFileSync(file, 'utf8');
    for (const match of source.matchAll(QUERY_KEY_REGEX)) {
      keys.push({ key: match[1], file });
    }
  }
  return keys;
};

const ALL_INVALIDATION_KEYS: ReadonlySet<string> = new Set(
  URL_TO_KEYS.flatMap(rule => rule.keys)
);

describe('cache-invalidation-map: every mutating axios endpoint is covered', () => {
  const allCalls = collectAllApiCalls();
  const mutatingCalls = allCalls.filter(c => c.method !== 'get');

  it('finds at least some axios mutation endpoints (sanity check)', () => {
    expect(mutatingCalls.length).toBeGreaterThan(10);
  });

  it('every non-GET axios endpoint either matches URL_TO_KEYS or is allow-listed', () => {
    const uncovered: ExtractedCall[] = [];

    for (const call of mutatingCalls) {
      if (NON_INVALIDATING_ENDPOINTS.has(call.url)) continue;
      const keys = resolveInvalidationKeys(call.url);
      if (keys.length === 0) uncovered.push(call);
    }

    if (uncovered.length > 0) {
      const lines = uncovered.map(
        c =>
          `  - ${c.method.toUpperCase()} ${c.url}\n      (in ${c.file.replace(SRC_ROOT, 'src')})`
      );
      throw new Error(
        [
          'The following mutating endpoints have no rule in `URL_TO_KEYS`',
          'and are not present in `NON_INVALIDATING_ENDPOINTS`:',
          ...lines,
          '',
          'Add a matching entry to `src/@api/cache-invalidation-map.ts` so',
          'related queries get invalidated automatically. If the endpoint is',
          'genuinely read-only or session-only, add it to',
          '`NON_INVALIDATING_ENDPOINTS` in this test file with a comment',
          'explaining why.'
        ].join('\n')
      );
    }
  });

  it('NON_INVALIDATING_ENDPOINTS only contains URLs that actually exist', () => {
    const allUrls = new Set(allCalls.map(c => c.url));
    const stale: string[] = [];
    for (const url of NON_INVALIDATING_ENDPOINTS) {
      if (!allUrls.has(url)) stale.push(url);
    }
    if (stale.length > 0) {
      throw new Error(
        [
          'NON_INVALIDATING_ENDPOINTS references URLs that no longer exist',
          'in `src/@api/`. Remove the stale entries:',
          ...stale.map(u => `  - ${u}`)
        ].join('\n')
      );
    }
  });
});

describe('cache-invalidation-map: every @queries key is reachable', () => {
  const queryKeys = collectAllQueryKeys();

  it('finds at least some query keys (sanity check)', () => {
    expect(queryKeys.length).toBeGreaterThan(10);
  });

  it('every query key is referenced in URL_TO_KEYS or explicitly excluded', () => {
    const orphaned: { key: string; file: string }[] = [];

    for (const { key, file } of queryKeys) {
      if (ALL_INVALIDATION_KEYS.has(key)) continue;
      if (NON_INVALIDATED_QUERY_KEYS.has(key)) continue;
      orphaned.push({ key, file });
    }

    if (orphaned.length > 0) {
      const lines = orphaned.map(
        ({ key, file }) =>
          `  - "${key}"\n      (in ${file.replace(SRC_ROOT, 'src')})`
      );
      throw new Error(
        [
          'The following query keys are defined in `src/@queries/` but are',
          'never invalidated by any rule in `URL_TO_KEYS`. That means data',
          'fetched under these keys can go stale silently after a mutation.',
          '',
          'Add the key to the relevant rule(s) in',
          '`src/@api/cache-invalidation-map.ts`, or — if the key is for',
          'read-only data with no mutating endpoint — add it to',
          '`NON_INVALIDATED_QUERY_KEYS` in this test file with a comment',
          'explaining why:',
          ...lines
        ].join('\n')
      );
    }
  });

  it('NON_INVALIDATED_QUERY_KEYS only contains keys that actually exist', () => {
    const allKeyValues = new Set(queryKeys.map(k => k.key));
    const stale: string[] = [];
    for (const key of NON_INVALIDATED_QUERY_KEYS) {
      if (!allKeyValues.has(key)) stale.push(key);
    }
    if (stale.length > 0) {
      throw new Error(
        [
          'NON_INVALIDATED_QUERY_KEYS references keys that no longer exist',
          'in `src/@queries/`. Remove the stale entries:',
          ...stale.map(k => `  - ${k}`)
        ].join('\n')
      );
    }
  });

  it('every key listed in URL_TO_KEYS is actually used by a query', () => {
    const allKeyValues = new Set(queryKeys.map(k => k.key));
    const unknown: string[] = [];
    for (const key of ALL_INVALIDATION_KEYS) {
      if (!allKeyValues.has(key)) unknown.push(key);
    }
    if (unknown.length > 0) {
      throw new Error(
        [
          'URL_TO_KEYS references keys that are not exported by any module',
          'in `src/@queries/`. Either rename the rule to match the real key',
          'or remove it:',
          ...unknown.map(k => `  - ${k}`)
        ].join('\n')
      );
    }
  });
});

describe('cache-invalidation-map: every axios instance installs the interceptor', () => {
  const AXIOS_CREATE_REGEX = /\baxios\.create\s*\(/;
  const INSTALL_INTERCEPTOR_REGEX =
    /\binstallCacheInvalidationInterceptor\s*\(/;

  const findAxiosInstanceFiles = (): string[] =>
    walk(API_ROOT)
      .filter(f => !f.endsWith('cache-invalidation-interceptor.ts'))
      .filter(f => AXIOS_CREATE_REGEX.test(readFileSync(f, 'utf8')));

  it('finds the shared axios client (sanity check)', () => {
    const files = findAxiosInstanceFiles();
    // We always expect at least the shared client. If this drops to zero the
    // test is misconfigured (e.g. wrong path or regex), so fail loudly.
    expect(
      files.some(f => f.endsWith('axios-client.ts')),
      `Expected to find @api/axios-client.ts among axios.create() sites, got:\n${files.join('\n')}`
    ).toBe(true);
  });

  it('every axios.create() site under src/@api/ installs the cache-invalidation interceptor', () => {
    const offenders: string[] = [];
    for (const file of findAxiosInstanceFiles()) {
      const source = readFileSync(file, 'utf8');
      if (!INSTALL_INTERCEPTOR_REGEX.test(source)) {
        offenders.push(file.replace(SRC_ROOT, 'src'));
      }
    }
    if (offenders.length > 0) {
      throw new Error(
        [
          'The following files create their own axios instance via',
          '`axios.create()` but do NOT call',
          '`installCacheInvalidationInterceptor(...)` on it. Mutations made',
          'through these instances will silently bypass React Query cache',
          'invalidation and leave stale data in the UI:',
          ...offenders.map(f => `  - ${f}`),
          '',
          'Either route the request through `@api/axios-client` (preferred)',
          'or call `installCacheInvalidationInterceptor(client)` from',
          '`@api/cache-invalidation-interceptor` after creating the instance.'
        ].join('\n')
      );
    }
  });
});

describe('cache-invalidation-map: non-axios HTTP transports are explicitly allow-listed', () => {
  // Detects native `fetch(...)` / `new XMLHttpRequest()` / `sendBeacon(...)`
  // usage. These bypass our axios response interceptor entirely, so the
  // cache-invalidation guardrails above do not protect them. The only safe
  // non-axios calls are read-only endpoints captured in
  // `NON_AXIOS_TRANSPORT_FILES`.
  const NON_AXIOS_TRANSPORT_REGEX =
    /\bfetch\s*\(|\bnew\s+XMLHttpRequest\s*\(|\bnavigator\.sendBeacon\s*\(/;

  const findNonAxiosTransportFiles = (): string[] =>
    walk(API_ROOT)
      .filter(f => !f.endsWith('cache-invalidation-map.test.ts'))
      .filter(f => NON_AXIOS_TRANSPORT_REGEX.test(readFileSync(f, 'utf8')));

  it('every non-axios HTTP call site under src/@api/ is in the allowlist', () => {
    const found = findNonAxiosTransportFiles().map(f =>
      f.replace(SRC_ROOT, 'src')
    );
    const offenders = found.filter(f => !NON_AXIOS_TRANSPORT_FILES.has(f));
    if (offenders.length > 0) {
      throw new Error(
        [
          'The following files use a non-axios HTTP transport (e.g.',
          '`fetch`, `XMLHttpRequest`, or `sendBeacon`) and therefore bypass',
          'the cache-invalidation interceptor:',
          ...offenders.map(f => `  - ${f}`),
          '',
          'If the endpoint is genuinely read-only (or has another reason it',
          'cannot use axios), add it to `NON_AXIOS_TRANSPORT_FILES` in this',
          'test file with a comment explaining why. Otherwise, migrate it',
          'to `@api/axios-client` so its mutations participate in cache',
          'invalidation.'
        ].join('\n')
      );
    }
  });

  it('NON_AXIOS_TRANSPORT_FILES only references files that actually exist', () => {
    const found = new Set(
      findNonAxiosTransportFiles().map(f => f.replace(SRC_ROOT, 'src'))
    );
    const stale: string[] = [];
    for (const f of NON_AXIOS_TRANSPORT_FILES) {
      if (!found.has(f)) stale.push(f);
    }
    if (stale.length > 0) {
      throw new Error(
        [
          'NON_AXIOS_TRANSPORT_FILES references files that either no longer',
          'exist or no longer contain a non-axios HTTP call. Remove the',
          'stale entries:',
          ...stale.map(f => `  - ${f}`)
        ].join('\n')
      );
    }
  });
});
