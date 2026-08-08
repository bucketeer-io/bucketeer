import { FunctionComponent } from 'react';
import { useTranslation } from 'i18n';
import { X } from 'lucide-react';
import { cn } from 'utils/style';
import { IconInfo } from '@icons';
import Button from 'components/button';
import Dropdown from 'components/dropdown';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';

export interface LanguageTabField {
  id: string;
  language: string;
}

interface LanguageMeta {
  label: string;
  englishName: string;
  icon: FunctionComponent;
}

interface LanguageTabsProps {
  fields: LanguageTabField[];
  activeLanguage: string;
  availableToAdd: string[];
  canRemove: boolean;
  languageMeta: Record<string, LanguageMeta>;
  onSelect: (language: string) => void;
  onAdd: (language: string) => void;
  onRemove: (index: number, language: string) => void;
}

// The "Languages" field of the publish form: one tab per authored language,
// each removable (while more than one remains), plus a dropdown to add any
// language not already present.
const LanguageTabs = ({
  fields,
  activeLanguage,
  availableToAdd,
  canRemove,
  languageMeta,
  onSelect,
  onAdd,
  onRemove
}: LanguageTabsProps) => {
  const { t } = useTranslation(['form']);

  const label = (language: string) => languageMeta[language]?.label ?? language;
  const icon = (language: string) => languageMeta[language]?.icon;

  return (
    <div className="flex flex-col gap-2">
      <div className="flex items-center gap-1">
        <label className="typo-para-medium text-gray-700">
          {t('form:languages')}
        </label>
        <Tooltip
          content={t('form:languages-info')}
          trigger={
            <span className="flex text-gray-400">
              <Icon icon={IconInfo} size="xxs" />
            </span>
          }
        />
      </div>
      <div className="flex items-center gap-2">
        <div className="flex items-center gap-2">
          {fields.map((field, index) => (
            <div
              key={field.id}
              className={cn(
                'flex items-center gap-2 rounded-md border px-3 py-[11px] min-w-[127px] typo-para-medium transition-colors',
                field.language === activeLanguage
                  ? 'border-primary-500 text-primary-500'
                  : 'border-gray-300 text-gray-600 hover:text-gray-900'
              )}
            >
              <Button
                type="button"
                variant="text"
                className="h-auto gap-2 px-0 text-current"
                onClick={() => onSelect(field.language)}
              >
                {icon(field.language) && (
                  <Icon icon={icon(field.language)!} size="xs" />
                )}
                {label(field.language)}
              </Button>
              {canRemove && (
                <Button
                  type="button"
                  variant="text"
                  aria-label={t('form:remove-language')}
                  onClick={() => onRemove(index, field.language)}
                  className="h-auto px-0 text-gray-400 hover:text-gray-600"
                >
                  <X size={14} />
                </Button>
              )}
            </div>
          ))}
        </div>

        <Dropdown
          className="w-[200px] py-[11px]"
          disabled={availableToAdd.length <= 0}
          placeholder={t('form:add-language')}
          options={availableToAdd.map(lang => ({
            value: lang,
            label: `${label(lang)} (${languageMeta[lang]?.englishName ?? lang})`,
            icon: icon(lang)
          }))}
          onChange={value => onAdd(String(value))}
        />
      </div>
    </div>
  );
};

export default LanguageTabs;
