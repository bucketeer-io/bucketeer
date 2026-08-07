import { Language } from 'i18n';
import { NotificationTag } from './types';

export const DRAFTS_PAGE_SIZE = 5;

interface TagPreset {
  color: string;
  names: Record<string, string>;
}

const TAG_PRESETS_SOURCE: TagPreset[] = [
  {
    color: '#3B82F6',
    names: {
      [Language.ENGLISH]: 'Announcement',
      [Language.JAPANESE]: 'お知らせ'
    }
  },
  {
    color: '#F97316',
    names: {
      [Language.ENGLISH]: 'Maintenance',
      [Language.JAPANESE]: 'メンテナンス'
    }
  },
  {
    color: '#8B5CF6',
    names: { [Language.ENGLISH]: 'Feature', [Language.JAPANESE]: '新機能' }
  },
  {
    color: '#6366F1',
    names: { [Language.ENGLISH]: 'Update', [Language.JAPANESE]: 'アップデート' }
  }
];

export const getTagPresets = (language: string): NotificationTag[] =>
  TAG_PRESETS_SOURCE.map(preset => ({
    name: preset.names[language] ?? preset.names[Language.ENGLISH],
    color: preset.color
  }));

export const TAG_COLOR_SWATCHES: readonly string[] = [
  '#3B82F6',
  '#F97316',
  '#8B5CF6',
  '#6366F1',
  '#10B981',
  '#EC4899',
  '#EAB308',
  '#EF4444'
];
