import { Language } from 'i18n';
import { NotificationTag } from './types';

export const DRAFTS_PAGE_SIZE = 5;

// Preset tags offered in the publish form. Names are localized per language;
// the color is shared. The suggestions the author sees follow the active
// language tab, mirroring how tags are stored per localization.
interface TagPreset {
  color: string;
  names: Record<string, string>; // language code -> display name
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

// Returns the preset tags with names resolved for `language`, falling back to
// English when a translation is missing.
export const getTagPresets = (language: string): NotificationTag[] =>
  TAG_PRESETS_SOURCE.map(preset => ({
    name: preset.names[language] ?? preset.names[Language.ENGLISH],
    color: preset.color
  }));
