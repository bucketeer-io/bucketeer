import type { AxiosInstance } from 'axios';
import MockAdapter from 'axios-mock-adapter';
import type {
  NotificationCenterFeedItem,
  NotificationCenterLocalization,
  NotificationCenterTag
} from '@types';
import { NotificationCenterStatus } from '@types';

// Backend for the notification center hasn't shipped yet (see
// proto/notification/service.proto — request/response messages are still
// empty). This intercepts the exact axios calls the real API modules make and
// answers with in-memory fake data, so the UI can be built end-to-end against
// the real request/response shapes now. Delete this file (and its `use()`
// call in axios-client.ts) once the backend implements these endpoints —
// nothing else needs to change, since the API modules, queries, and
// components all only ever talk to `axiosClient`.

const HOUR = 60 * 60 * 1000;
const DAY = 24 * HOUR;
const now = Date.now();
const toEpochSeconds = (ms: number) => String(Math.floor(ms / 1000));

const TAGS: Record<string, { color: string; en: string; ja: string }> = {
  announcement: { color: '#3B82F6', en: 'Announcement', ja: 'お知らせ' },
  maintenance: { color: '#F97316', en: 'Maintenance', ja: 'メンテナンス' },
  feature: { color: '#8B5CF6', en: 'Feature', ja: '新機能' },
  update: { color: '#6366F1', en: 'Update', ja: 'アップデート' }
};

const tag = (
  key: keyof typeof TAGS,
  lang: 'en' | 'ja'
): NotificationCenterTag => ({
  name: TAGS[key][lang],
  color: TAGS[key].color
});

const SEED_EMAIL = 'admin@bucketeer.io';

const SEED: {
  id: string;
  status: NotificationCenterStatus;
  publishedAgo: number;
  tags: (keyof typeof TAGS)[];
  en: { title: string; content: string };
  ja: { title: string; content: string };
}[] = [
  {
    id: '1',
    status: NotificationCenterStatus.PUBLISHED,
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
    status: NotificationCenterStatus.PUBLISHED,
    publishedAgo: DAY,
    tags: ['feature'],
    en: {
      title: 'New dashboard navigation update',
      content:
        'We are rolling out a refreshed navigation experience next week. See the [changelog](https://bucketeer.io/changelog) for the full list of changes.'
    },
    ja: {
      title: 'ダッシュボードのナビゲーション刷新',
      content:
        '来週、刷新されたナビゲーションを順次公開します。詳細は[更新履歴](https://bucketeer.io/changelog)をご覧ください。'
    }
  },
  {
    id: '3',
    status: NotificationCenterStatus.PUBLISHED,
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
    status: NotificationCenterStatus.PUBLISHED,
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
    status: NotificationCenterStatus.PUBLISHED,
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
    status: NotificationCenterStatus.PUBLISHED,
    publishedAgo: 9 * DAY,
    tags: ['update'],
    en: {
      title: 'API rate limit increase',
      content:
        'Default API rate limits have been raised for all plans. [Read the docs](https://docs.bucketeer.io/api/rate-limits) for the new thresholds.'
    },
    ja: {
      title: 'API レート制限の引き上げ',
      content:
        '全プランのデフォルト API レート制限を引き上げました。新しい上限値は[ドキュメント](https://docs.bucketeer.io/api/rate-limits)をご確認ください。'
    }
  },
  {
    id: '7',
    status: NotificationCenterStatus.PUBLISHED,
    publishedAgo: 10 * DAY,
    tags: ['feature'],
    en: {
      title: 'Segment targeting improvements',
      content:
        'You can now combine multiple segment rules with AND/OR logic when targeting a feature flag.'
    },
    ja: {
      title: 'セグメントターゲティングの改善',
      content:
        '機能フラグのターゲティング時に、複数のセグメントルールを AND/OR で組み合わせられるようになりました。'
    }
  },
  {
    id: '8',
    status: NotificationCenterStatus.PUBLISHED,
    publishedAgo: 11 * DAY,
    tags: ['maintenance'],
    en: {
      title: 'Database failover drill',
      content:
        'We will run a routine database failover drill on Tuesday. No downtime is expected.'
    },
    ja: {
      title: 'データベースフェイルオーバー訓練',
      content:
        '火曜日に定期的なデータベースフェイルオーバー訓練を実施します。ダウンタイムは発生しない見込みです。'
    }
  },
  {
    id: '9',
    status: NotificationCenterStatus.PUBLISHED,
    publishedAgo: 12 * DAY,
    tags: ['announcement'],
    en: {
      title: 'New pricing page live',
      content: 'Our updated pricing page is now live with clearer plan tiers.'
    },
    ja: {
      title: '新しい料金ページを公開',
      content: 'より分かりやすいプラン構成の新しい料金ページを公開しました。'
    }
  },
  {
    id: '10',
    status: NotificationCenterStatus.PUBLISHED,
    publishedAgo: 13 * DAY,
    tags: ['update'],
    en: {
      title: 'Audit log retention extended',
      content:
        'Audit log retention has been extended from 30 to 90 days on Enterprise plans.'
    },
    ja: {
      title: '監査ログの保存期間を延長',
      content:
        'Enterprise プランの監査ログ保存期間を 30 日から 90 日に延長しました。'
    }
  },
  {
    id: '11',
    status: NotificationCenterStatus.PUBLISHED,
    publishedAgo: 14 * DAY,
    tags: ['feature'],
    en: {
      title: 'Bulk flag import',
      content:
        'You can now bulk import feature flags from a CSV file on the Features page.'
    },
    ja: {
      title: 'フラグの一括インポート',
      content:
        'Features ページから CSV ファイルによる機能フラグの一括インポートができるようになりました。'
    }
  },
  {
    id: '12',
    status: NotificationCenterStatus.PUBLISHED,
    publishedAgo: 15 * DAY,
    tags: ['maintenance'],
    en: {
      title: 'TLS certificate rotation',
      content:
        'We rotated TLS certificates for all API endpoints. No action is required on your end.'
    },
    ja: {
      title: 'TLS 証明書のローテーション',
      content:
        '全 API エンドポイントの TLS 証明書をローテーションしました。特別な対応は不要です。'
    }
  },
  {
    id: '13',
    status: NotificationCenterStatus.PUBLISHED,
    publishedAgo: 16 * DAY,
    tags: ['announcement'],
    en: {
      title: 'Team workspace limits raised',
      content:
        'The maximum number of members per team workspace has been increased.'
    },
    ja: {
      title: 'チームワークスペースの上限を引き上げ',
      content: 'チームワークスペースあたりの最大メンバー数を引き上げました。'
    }
  },
  {
    id: '14',
    status: NotificationCenterStatus.PUBLISHED,
    publishedAgo: 17 * DAY,
    tags: ['update'],
    en: {
      title: 'Webhook retry policy updated',
      content:
        'Failed webhook deliveries are now retried with exponential backoff for up to 24 hours.'
    },
    ja: {
      title: 'Webhook 再試行ポリシーの更新',
      content:
        'Webhook 配信の失敗時、最大 24 時間まで指数バックオフで再試行するようになりました。'
    }
  },
  {
    id: '15',
    status: NotificationCenterStatus.PUBLISHED,
    publishedAgo: 18 * DAY,
    tags: ['feature'],
    en: {
      title: 'Experiment results export to Slack',
      content:
        'Experiment result summaries can now be posted directly to a Slack channel.'
    },
    ja: {
      title: '実験結果の Slack エクスポート',
      content:
        '実験結果のサマリーを Slack チャンネルへ直接投稿できるようになりました。'
    }
  },
  {
    id: '16',
    status: NotificationCenterStatus.DRAFT,
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
  },
  {
    id: '17',
    status: NotificationCenterStatus.DRAFT,
    publishedAgo: 6 * HOUR,
    tags: ['feature'],
    en: {
      title: 'Onboarding flow revamp (draft)',
      content:
        'Draft notes for the upcoming onboarding flow revamp. Still gathering feedback.'
    },
    ja: {
      title: 'オンボーディングフローの刷新（下書き）',
      content:
        'オンボーディングフロー刷新に関する下書きメモです。現在フィードバックを収集中です。'
    }
  },
  {
    id: '18',
    status: NotificationCenterStatus.DRAFT,
    publishedAgo: 30 * HOUR,
    tags: ['maintenance'],
    en: {
      title: 'Q3 maintenance window plan (draft)',
      content: 'Draft plan for the Q3 maintenance window. Dates TBD.'
    },
    ja: {
      title: 'Q3 メンテナンス計画（下書き）',
      content: 'Q3 のメンテナンス計画に関する下書きです。日程は未定です。'
    }
  }
];

interface StoreRecord {
  id: string;
  status: NotificationCenterStatus;
  createdBy: string;
  lastEditedBy: string;
  publishedBy?: string;
  publishedAt: number; // ms; 0 while draft
  createdAt: number;
  updatedAt: number;
  localizations: NotificationCenterLocalization[];
}

let notifications: StoreRecord[] = SEED.map(s => ({
  id: s.id,
  status: s.status,
  createdBy: SEED_EMAIL,
  lastEditedBy: SEED_EMAIL,
  publishedBy:
    s.status === NotificationCenterStatus.PUBLISHED ? SEED_EMAIL : undefined,
  publishedAt:
    s.status === NotificationCenterStatus.PUBLISHED ? now - s.publishedAgo : 0,
  createdAt: now - s.publishedAgo,
  updatedAt: now - s.publishedAgo,
  localizations: [
    {
      language: 'en',
      tags: s.tags.map(k => tag(k, 'en')),
      title: s.en.title,
      content: s.en.content
    },
    {
      language: 'ja',
      tags: s.tags.map(k => tag(k, 'ja')),
      title: s.ja.title,
      content: s.ja.content
    }
  ]
}));

// email -> set of read notification ids
const reads = new Map<string, Set<string>>();

const readSetFor = (email: string): Set<string> => {
  if (!reads.has(email)) reads.set(email, new Set());
  return reads.get(email)!;
};

const localizationFor = (
  record: StoreRecord,
  lang: string
): NotificationCenterLocalization | undefined =>
  record.localizations.find(l => l.language === lang) ??
  record.localizations.find(l => l.language === 'en') ??
  record.localizations[0];

const toFeedItem = (
  record: StoreRecord,
  email: string,
  lang: string
): NotificationCenterFeedItem => {
  const loc = localizationFor(record, lang);
  return {
    id: record.id,
    title: loc?.title ?? '',
    content: loc?.content ?? '',
    tags: loc?.tags ?? [],
    read: readSetFor(email).has(record.id),
    status: record.status,
    publishedAt: toEpochSeconds(record.publishedAt),
    createdAt: toEpochSeconds(record.createdAt),
    updatedAt: toEpochSeconds(record.updatedAt),
    createdBy: record.createdBy,
    lastEditedBy: record.lastEditedBy,
    localizations: record.localizations
  };
};

// A single fake viewer identity, since the mock has no access to the
// authenticated user's email (that lives in the JWT, not the request).
const VIEWER_EMAIL = 'viewer@bucketeer.io';

export const installNotificationCenterMockAdapter = (client: AxiosInstance) => {
  const mock = new MockAdapter(client, { delayResponse: 300 });

  mock.onGet(/\/v1\/notifications\?/).reply(config => {
    const params = new URLSearchParams(config.url?.split('?')[1]);
    const read = params.get('read') === 'true';
    const searchKeyword = (params.get('searchKeyword') ?? '').toLowerCase();
    const cursor = Number(params.get('cursor') ?? 0);
    const pageSize = Number(params.get('pageSize') ?? 10);
    const orderDirection = params.get('orderDirection') ?? 'DESC';

    const published = notifications
      .filter(n => n.status === NotificationCenterStatus.PUBLISHED)
      .map(n => toFeedItem(n, VIEWER_EMAIL, 'en'));

    const unreadCount = published.filter(n => !n.read).length;
    const readCount = published.length - unreadCount;

    const filtered = published
      .filter(n => n.read === read)
      .filter(n =>
        searchKeyword
          ? n.title.toLowerCase().includes(searchKeyword) ||
            n.content.toLowerCase().includes(searchKeyword)
          : true
      )
      .sort((a, b) =>
        orderDirection === 'ASC'
          ? Number(a.createdAt) - Number(b.createdAt)
          : Number(b.createdAt) - Number(a.createdAt)
      );

    const page = filtered.slice(cursor, cursor + pageSize);

    return [
      200,
      {
        notifications: page,
        cursor: String(cursor + page.length),
        totalCount: String(filtered.length),
        unreadCount: String(unreadCount),
        readCount: String(readCount)
      }
    ];
  });

  mock.onGet(/\/v1\/notifications\/drafts\?/).reply(config => {
    const params = new URLSearchParams(config.url?.split('?')[1]);
    const cursor = Number(params.get('cursor') ?? 0);
    const pageSize = Number(params.get('pageSize') ?? 50);

    const drafts = notifications
      .filter(n => n.status === NotificationCenterStatus.DRAFT)
      .sort((a, b) => b.updatedAt - a.updatedAt)
      .map(n => toFeedItem(n, VIEWER_EMAIL, 'en'));

    const page = drafts.slice(cursor, cursor + pageSize);

    return [
      200,
      {
        notifications: page,
        cursor: String(cursor + page.length),
        totalCount: String(drafts.length)
      }
    ];
  });

  mock.onGet(/\/v1\/notifications\/unread_count/).reply(() => {
    const unreadCount = notifications
      .filter(n => n.status === NotificationCenterStatus.PUBLISHED)
      .filter(n => !readSetFor(VIEWER_EMAIL).has(n.id)).length;
    return [200, { count: String(unreadCount) }];
  });

  mock.onPost('/v1/notification').reply(config => {
    const body = JSON.parse(config.data);
    const ts = Date.now();
    const record: StoreRecord = {
      id: `draft-${ts}`,
      status: NotificationCenterStatus.DRAFT,
      createdBy: VIEWER_EMAIL,
      lastEditedBy: VIEWER_EMAIL,
      publishedAt: 0,
      createdAt: ts,
      updatedAt: ts,
      localizations: body.localizations
    };
    notifications = [record, ...notifications];
    return [200, { notification: toFeedItem(record, VIEWER_EMAIL, 'en') }];
  });

  mock.onPatch('/v1/notification').reply(config => {
    const body = JSON.parse(config.data);
    const ts = Date.now();
    notifications = notifications.map(n =>
      n.id === body.id
        ? {
            ...n,
            localizations: body.localizations,
            lastEditedBy: VIEWER_EMAIL,
            updatedAt: ts
          }
        : n
    );
    const record = notifications.find(n => n.id === body.id)!;
    return [200, { notification: toFeedItem(record, VIEWER_EMAIL, 'en') }];
  });

  mock.onPost('/v1/notification/publish').reply(config => {
    const body = JSON.parse(config.data);
    const ts = Date.now();

    if (body.id) {
      notifications = notifications.map(n =>
        n.id === body.id
          ? {
              ...n,
              status: NotificationCenterStatus.PUBLISHED,
              localizations: body.localizations,
              lastEditedBy: VIEWER_EMAIL,
              publishedBy: VIEWER_EMAIL,
              publishedAt: n.publishedAt || ts,
              updatedAt: ts
            }
          : n
      );
      const record = notifications.find(n => n.id === body.id)!;
      return [200, { notification: toFeedItem(record, VIEWER_EMAIL, 'en') }];
    }

    const record: StoreRecord = {
      id: `${ts}`,
      status: NotificationCenterStatus.PUBLISHED,
      createdBy: VIEWER_EMAIL,
      lastEditedBy: VIEWER_EMAIL,
      publishedBy: VIEWER_EMAIL,
      publishedAt: ts,
      createdAt: ts,
      updatedAt: ts,
      localizations: body.localizations
    };
    notifications = [record, ...notifications];
    return [200, { notification: toFeedItem(record, VIEWER_EMAIL, 'en') }];
  });

  mock.onDelete(/\/v1\/notification\?/).reply(config => {
    const params = new URLSearchParams(config.url?.split('?')[1]);
    const id = params.get('id');
    notifications = notifications.filter(n => n.id !== id);
    reads.forEach(set => set.delete(id ?? ''));
    return [200, {}];
  });

  mock.onPost('/v1/notifications/mark_as_read').reply(config => {
    const body = JSON.parse(config.data);
    const set = readSetFor(VIEWER_EMAIL);
    (body.ids as string[]).forEach(id => set.add(id));
    return [200, {}];
  });

  mock.onPost('/v1/notifications/mark_all_as_read').reply(() => {
    const set = readSetFor(VIEWER_EMAIL);
    notifications
      .filter(n => n.status === NotificationCenterStatus.PUBLISHED)
      .forEach(n => set.add(n.id));
    return [200, {}];
  });

  // Anything else goes through to the real network as normal.
  mock.onAny().passThrough();
};
