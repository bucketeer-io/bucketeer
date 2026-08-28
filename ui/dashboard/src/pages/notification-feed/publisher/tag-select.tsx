import { useRef } from 'react';
import { MultiValueGenericProps, StylesConfig } from 'react-select';
import { useTheme } from 'hooks/use-theme';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import {
  buildColorStyles,
  CreatableSelect,
  Option
} from 'components/creatable-select';
import { Popover } from 'components/popover';
import { getTagPresets, TAG_COLOR_SWATCHES } from '../constants';
import { NotificationTag } from '../types';

interface TagSelectProps {
  value: NotificationTag[];
  // Active language; the preset tag names shown are localized to it.
  language: string;
  onChange: (tags: NotificationTag[]) => void;
}

const toOption = (tag: NotificationTag): Option => ({
  value: tag.name,
  label: tag.name,
  color: tag.color
});

const contentFont = {
  fontSize: '14px',
  fontFamily: "'Helvetica Neue', Helvetica, Arial, sans-serif"
};

const buildTagColorStyles = (
  isDark: boolean
): StylesConfig<Option, boolean> => {
  const colorStyles = buildColorStyles(isDark);
  return {
    control: (base, props) => ({
      ...colorStyles.control?.(base, props),
      ...contentFont
    }),
    menu: (base, props) => ({
      ...colorStyles.menu?.(base, props),
      ...contentFont
    }),
    multiValue: (base, { data }) => {
      const color = (data.color as string) || undefined;
      return {
        ...base,
        backgroundColor: color ? `${color}1A` : base.backgroundColor,
        borderRadius: '4px',
        padding: '4px'
      };
    },
    multiValueLabel: (base, { data }) => ({
      ...base,
      color: (data.color as string) || (base.color as string),
      padding: 0
    }),
    multiValueRemove: (base, { data }) => {
      const color = (data.color as string) || undefined;
      return {
        ...base,
        color: color ?? (base.color as string),
        ':hover': {
          backgroundColor: color ? `${color}33` : undefined,
          color: color ?? undefined
        }
      };
    }
  };
};

const makeTagColorLabel = (
  onSwatchPick: (name: string, color: string) => void
) => {
  const TagColorLabel = (props: MultiValueGenericProps<Option, true>) => {
    const { t } = useTranslation(['common']);
    const closeRef = useRef<HTMLButtonElement>(null);
    const { data } = props;
    const name = data.value;
    const color = (data.color as string) || '';

    return (
      <Popover
        align="start"
        modal
        closeRef={closeRef}
        trigger={
          <span className="flex items-center gap-1.5 px-2 py-0.5 typo-para-tiny font-medium cursor-pointer">
            <span
              className={cn(
                'size-1.5 shrink-0 rounded-full',
                !color && 'bg-gray-400 dark:bg-dark-gray-200'
              )}
              style={color ? { backgroundColor: color } : undefined}
            />
            {name}
          </span>
        }
      >
        <div className="p-2">
          <p className="mb-2 typo-para-tiny text-gray-500 dark:text-dark-gray-200">
            {t('tag-color')}
          </p>
          <div className="grid grid-cols-4 gap-2">
            {TAG_COLOR_SWATCHES.map(swatch => (
              <button
                key={swatch}
                type="button"
                aria-label={swatch}
                className={cn(
                  'size-6 rounded-full ring-offset-1 dark:ring-offset-dark-black-900 hover:ring-2 hover:ring-gray-300 dark:hover:ring-dark-black-600',
                  color === swatch &&
                    'ring-2 ring-gray-500 dark:ring-dark-gray-200'
                )}
                style={{ backgroundColor: swatch }}
                onClick={() => {
                  onSwatchPick(name, swatch);
                  closeRef.current?.click();
                }}
              />
            ))}
          </div>
        </div>
      </Popover>
    );
  };
  return TagColorLabel;
};

const TagSelect = ({ value, language, onChange }: TagSelectProps) => {
  const { t } = useTranslation(['form']);
  const { theme } = useTheme();
  const tagColorStyles = buildTagColorStyles(theme === 'dark');

  const presets = getTagPresets(language);
  const options = presets.map(toOption);
  const selected = value.map(toOption);

  const handleChange = (opts: readonly Option[]) => {
    onChange(
      opts.map(o => {
        const existing = value.find(tag => tag.name === o.value);
        const preset = presets.find(p => p.name === o.value);
        return {
          name: o.label,
          color: existing?.color ?? preset?.color ?? (o.color as string) ?? ''
        };
      })
    );
  };

  const onSwatchPick = (name: string, color: string) => {
    onChange(value.map(tag => (tag.name === name ? { ...tag, color } : tag)));
  };

  return (
    <CreatableSelect
      isMulti
      options={options}
      value={selected}
      styles={tagColorStyles}
      placeholder={t('form:add-tags')}
      onChange={opts => handleChange(opts)}
      onCreateOption={name =>
        onChange([...value, { name, color: TAG_COLOR_SWATCHES[0] }])
      }
      components={{
        MultiValueLabel: makeTagColorLabel(onSwatchPick)
      }}
    />
  );
};

export default TagSelect;
