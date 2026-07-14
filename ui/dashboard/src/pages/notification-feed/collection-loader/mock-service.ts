import { markdownToText } from '../markdown-content';
import {
  FeedNotification,
  FeedPage,
  FeedQuery,
  NotificationDraft,
  NotificationLocalization,
  NotificationLocalizationInput,
  NotificationReadRow,
  NotificationRecord,
  NotificationStatus,
  PublishNotificationInput
} from '../types';

// Temporary in-memory data source, shaped after the backend SQL schema:
//   notifications  -> `notification`              (status + audit + timestamps)
//   localizations  -> `notification_localization` (per-language title/tags/body)
//   reads          -> `notification_read`         (per-viewer read state)
//
// UI components never see these rows directly — `project()` resolves a record
// into the current language and folds in the current viewer's read state.
// Replace with real gRPC/react-query calls once the backend endpoints exist.

const HOUR = 60 * 60 * 1000;
const DAY = 24 * HOUR;

const now = Date.now();

// English display name -> localized names, so seed localizations stay in sync
// across languages. Colors come from the English preset.
const TAGS = {
  announcement: { color: '#3B82F6', en: 'Announcement', ja: 'お知らせ' },
  maintenance: { color: '#F97316', en: 'Maintenance', ja: 'メンテナンス' },
  feature: { color: '#8B5CF6', en: 'Feature', ja: '新機能' },
  update: { color: '#6366F1', en: 'Update', ja: 'アップデート' }
} as const;

type TagKey = keyof typeof TAGS;

const tag = (key: TagKey, lang: 'en' | 'ja') => ({
  name: TAGS[key][lang],
  color: TAGS[key].color
});

const SEED_EMAIL = 'admin@bucketeer.io';

// Seed content, authored once in both languages. Each entry becomes one
// `notification` record plus an `en` and `ja` localization.
const SEED: {
  id: string;
  status: NotificationStatus;
  publishedAgo: number; // ms before `now`; also used as createdAt offset
  tags: TagKey[];
  en: { title: string; content: string };
  ja: { title: string; content: string };
}[] = [
  {
    id: '1',
    status: NotificationStatus.PUBLISHED,
    publishedAgo: 2 * HOUR,
    tags: ['announcement', 'maintenance'],
    en: {
      title: 'Holiday support hours',
      content:
        'Our support team will have reduced availability during the upcoming holidays. Response times may be slower than usual.\n\n**Impact**\n\n- Slower response times\n- Limited live support\n\n> Thanks for your understanding and patience!'
    },
    ja: {
      title: '年末年始のサポート時間',
      content:
        '年末年始の期間中、サポートチームの対応時間が短縮されます。通常よりも回答に時間がかかる場合があります。\n\n**影響**\n\n- 回答までの時間が長くなります\n- ライブサポートが制限されます\n\n> ご理解とご協力をお願いいたします。'
    }
  },
  {
    id: '2',
    status: NotificationStatus.PUBLISHED,
    publishedAgo: DAY,
    tags: ['feature'],
    en: {
      title: 'New dashboard navigation update',
      content:
        'We are rolling out a refreshed navigation experience next week. Keep an eye out for the new layout.'
    },
    ja: {
      title: 'ダッシュボードのナビゲーション刷新',
      content:
        '来週、刷新されたナビゲーションを順次公開します。新しいレイアウトにご注目ください。'
    }
  },
  {
    id: '3',
    status: NotificationStatus.PUBLISHED,
    publishedAgo: 4 * DAY,
    tags: ['update'],
    en: {
      title: 'Deprecating legacy SDK v1',
      content:
        'SDK v1 will reach end of life in three months. Please upgrade to v2 to continue receiving updates.'
    },
    ja: {
      title: '旧 SDK v1 の提供終了',
      content:
        'SDK v1 は 3 か月後にサポートを終了します。引き続き更新を受け取るには v2 へのアップグレードをお願いします。'
    }
  },
  {
    id: '4',
    status: NotificationStatus.PUBLISHED,
    publishedAgo: 5 * DAY,
    tags: ['maintenance'],
    en: {
      title: 'Scheduled maintenance',
      content:
        'We will be performing scheduled maintenance on Sunday from 02:00–04:00 UTC.'
    },
    ja: {
      title: '定期メンテナンスのお知らせ',
      content: '日曜日の 02:00〜04:00 (UTC) に定期メンテナンスを実施します。'
    }
  },
  {
    id: '5',
    status: NotificationStatus.PUBLISHED,
    publishedAgo: 7 * DAY,
    tags: ['feature'],
    en: {
      title: 'New analytics export',
      content:
        'You can now export experiment data in CSV format directly from the insights page.'
    },
    ja: {
      title: '分析データのエクスポート機能',
      content:
        'インサイトページからエクスペリメントのデータを CSV 形式で直接エクスポートできるようになりました。'
    }
  },
  {
    id: '6',
    status: NotificationStatus.DRAFT,
    publishedAgo: 3 * HOUR,
    tags: ['announcement'],
    en: {
      title: 'Upcoming pricing update (draft)',
      content:
        '# Pricing update\n\nWe are refreshing our pricing next quarter.\n\n- Simpler tiers\n- Clearer usage limits'
    },
    ja: {
      title: '料金改定のお知らせ（下書き）',
      content:
        '# 料金改定\n\n来四半期に料金体系を見直します。\n\n- よりシンプルなプラン\n- 分かりやすい利用上限'
    }
  }
];

// ----------------------------------------------------------------------------
// In-memory stores (mirroring the three tables)
// ----------------------------------------------------------------------------

let notifications: NotificationRecord[] = SEED.map(s => ({
  id: s.id,
  status: s.status,
  createdBy: SEED_EMAIL,
  lastEditedBy: SEED_EMAIL,
  publishedBy:
    s.status === NotificationStatus.PUBLISHED ? SEED_EMAIL : undefined,
  publishedAt:
    s.status === NotificationStatus.PUBLISHED ? now - s.publishedAgo : 0,
  createdAt: now - s.publishedAgo,
  updatedAt: now - s.publishedAgo
}));

let localizations: NotificationLocalization[] = SEED.flatMap(s => [
  {
    notificationId: s.id,
    language: 'en',
    tags: s.tags.map(k => tag(k, 'en')),
    title: s.en.title,
    content: s.en.content
  },
  {
    notificationId: s.id,
    language: 'ja',
    tags: s.tags.map(k => tag(k, 'ja')),
    title: s.ja.title,
    content: s.ja.content
  }
]);

let reads: NotificationReadRow[] = [];

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

const delay = <T>(value: T, ms = 300): Promise<T> =>
  new Promise(resolve => setTimeout(() => resolve(value), ms));

// Resolves the localization for `lang`, falling back to English so a missing
// translation never leaves the UI blank.
const localizationFor = (
  notificationId: string,
  lang: string
): NotificationLocalization | undefined => {
  const forId = localizations.filter(l => l.notificationId === notificationId);
  return (
    forId.find(l => l.language === lang) ??
    forId.find(l => l.language === 'en') ??
    forId[0]
  );
};

// All language versions of a record, shaped for the publish form so editing a
// draft can load and edit every language.
const localizationsFor = (
  notificationId: string
): NotificationLocalizationInput[] =>
  localizations
    .filter(l => l.notificationId === notificationId)
    .map(l => ({
      language: l.language,
      title: l.title,
      content: l.content,
      tags: l.tags
    }));

const isRead = (notificationId: string, email: string): boolean =>
  reads.some(r => r.notificationId === notificationId && r.email === email);

// Projects a record into the current language + viewer's read state.
const project = (
  record: NotificationRecord,
  email: string,
  lang: string
): FeedNotification => {
  const loc = localizationFor(record.id, lang);
  return {
    id: record.id,
    title: loc?.title ?? '',
    content: loc?.content ?? '',
    tags: loc?.tags ?? [],
    read: isRead(record.id, email),
    status: record.status,
    publishedAt: record.publishedAt,
    createdAt: record.createdAt,
    updatedAt: record.updatedAt,
    createdBy: record.createdBy,
    lastEditedBy: record.lastEditedBy,
    localizations: localizationsFor(record.id)
  };
};

// NOTE: `environmentId` is accepted to mirror the eventual gRPC signature. The
// SQL schema has no environment column, so the in-memory store is global for
// now; this is flagged for the backend design.

// ----------------------------------------------------------------------------
// Reads
// ----------------------------------------------------------------------------

// Published notifications for the given tab, with filter/sort/paginate applied
// server-side. Returns the page plus totals so the UI can render "Showing
// 1–10 of N" and the Unread/Read tab counts.
export const fetchFeed = (
  environmentId: string,
  email: string,
  lang: string,
  query: FeedQuery
): Promise<FeedPage> => {
  void environmentId;

  const published = notifications
    .filter(n => n.status === NotificationStatus.PUBLISHED)
    .map(n => project(n, email, lang));

  const unreadCount = published.filter(n => !n.read).length;
  const readCount = published.length - unreadCount;

  const q = (query.searchQuery ?? '').trim().toLowerCase();
  const filtered = published
    .filter(n => n.read === query.read)
    .filter(n => (query.from ? n.createdAt >= query.from : true))
    .filter(n => (query.to ? n.createdAt <= query.to : true))
    .filter(n =>
      q
        ? n.title.toLowerCase().includes(q) ||
          markdownToText(n.content).toLowerCase().includes(q)
        : true
    )
    .sort((a, b) =>
      query.sort === 'oldest'
        ? a.createdAt - b.createdAt
        : b.createdAt - a.createdAt
    );

  const start = (query.page - 1) * query.pageSize;
  const items = filtered.slice(start, start + query.pageSize);

  return delay({ items, total: filtered.length, unreadCount, readCount });
};

export const fetchDrafts = (
  environmentId: string,
  email: string,
  lang: string
): Promise<NotificationDraft[]> => {
  void environmentId;
  const drafts = notifications
    .filter(n => n.status === NotificationStatus.DRAFT)
    .sort((a, b) => b.updatedAt - a.updatedAt)
    .map(record => {
      const view = project(record, email, lang);
      return {
        id: record.id,
        title: view.title,
        content: view.content,
        tags: view.tags,
        status: record.status,
        createdAt: record.createdAt,
        updatedAt: record.updatedAt,
        createdBy: record.createdBy,
        lastEditedBy: record.lastEditedBy,
        localizations: localizationsFor(record.id)
      };
    });
  return delay(drafts);
};

// ----------------------------------------------------------------------------
// Read-state mutations (per viewer, keyed by email)
// ----------------------------------------------------------------------------

const upsertRead = (notificationId: string, email: string) => {
  if (!isRead(notificationId, email)) {
    reads = [...reads, { notificationId, email, readAt: Date.now() }];
  }
};

export const markAsRead = (
  environmentId: string,
  id: string,
  email: string
): Promise<void> => {
  void environmentId;
  upsertRead(id, email);
  return delay(undefined);
};

export const markManyAsRead = (
  environmentId: string,
  ids: string[],
  email: string
): Promise<void> => {
  void environmentId;
  ids.forEach(id => upsertRead(id, email));
  return delay(undefined);
};

export const markAllAsRead = (
  environmentId: string,
  email: string
): Promise<void> => {
  void environmentId;
  notifications
    .filter(n => n.status === NotificationStatus.PUBLISHED)
    .forEach(n => upsertRead(n.id, email));
  return delay(undefined);
};

// ----------------------------------------------------------------------------
// Writes (publish flow) — create a record + its per-language localizations
// ----------------------------------------------------------------------------

// Replaces the localization rows for `id` with the ones in the request. Mirrors
// the backend writing one `notification_localization` row per provided language.
const writeLocalizations = (id: string, input: PublishNotificationInput) => {
  localizations = localizations.filter(l => l.notificationId !== id);
  input.localizations.forEach(loc => {
    localizations.push({
      notificationId: id,
      language: loc.language,
      tags: loc.tags,
      title: loc.title,
      content: loc.content
    });
  });
};

export const publishNotification = (
  environmentId: string,
  email: string,
  input: PublishNotificationInput
): Promise<FeedNotification> => {
  void environmentId;
  const ts = Date.now();
  const id = `${ts}`;
  const record: NotificationRecord = {
    id,
    status: NotificationStatus.PUBLISHED,
    createdBy: email,
    lastEditedBy: email,
    publishedBy: email,
    publishedAt: ts,
    createdAt: ts,
    updatedAt: ts
  };
  notifications = [record, ...notifications];
  writeLocalizations(id, input);
  return delay(project(record, email, 'en'));
};

export const saveDraft = (
  environmentId: string,
  email: string,
  input: PublishNotificationInput
): Promise<NotificationDraft> => {
  void environmentId;
  const ts = Date.now();
  const id = `draft-${ts}`;
  const record: NotificationRecord = {
    id,
    status: NotificationStatus.DRAFT,
    createdBy: email,
    lastEditedBy: email,
    publishedAt: 0,
    createdAt: ts,
    updatedAt: ts
  };
  notifications = [record, ...notifications];
  writeLocalizations(id, input);
  const view = project(record, email, 'en');
  return delay({
    id,
    title: view.title,
    content: view.content,
    tags: view.tags,
    status: record.status,
    createdAt: record.createdAt,
    updatedAt: record.updatedAt,
    createdBy: record.createdBy,
    lastEditedBy: record.lastEditedBy,
    localizations: localizationsFor(id)
  });
};

// Updates an existing notification in place: rewrites its localizations and
// applies the given status. Editing a draft keeps the same id (no duplicate);
// publishing an edited draft promotes the same record to PUBLISHED.
export const updateNotification = (
  environmentId: string,
  email: string,
  id: string,
  input: PublishNotificationInput
): Promise<FeedNotification> => {
  void environmentId;
  const ts = Date.now();
  const wasPublished = input.status === NotificationStatus.PUBLISHED;
  notifications = notifications.map(n =>
    n.id === id
      ? {
          ...n,
          status: input.status,
          lastEditedBy: email,
          updatedAt: ts,
          publishedBy: wasPublished ? email : n.publishedBy,
          publishedAt: wasPublished && !n.publishedAt ? ts : n.publishedAt
        }
      : n
  );
  writeLocalizations(id, input);
  const record = notifications.find(n => n.id === id)!;
  return delay(project(record, email, 'en'));
};
