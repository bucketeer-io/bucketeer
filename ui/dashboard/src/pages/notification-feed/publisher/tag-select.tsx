import { StylesConfig } from 'react-select';
import { useTranslation } from 'i18n';
import { CreatableSelect, Option } from 'components/creatable-select';
import { getTagPresets } from '../constants';
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

// Tint each selected chip with its own tag color instead of the default
// purple defined in the shared creatable-select styles.
const tagColorStyles: StylesConfig<Option, boolean> = {
  multiValue: (base, { data }) => {
    const color = (data.color as string) || undefined;
    return {
      ...base,
      backgroundColor: color ? `${color}1A` : base.backgroundColor,
      borderRadius: '4px'
    };
  },
  multiValueLabel: (base, { data }) => ({
    ...base,
    color: (data.color as string) || (base.color as string)
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

const TagSelect = ({ value, language, onChange }: TagSelectProps) => {
  const { t } = useTranslation(['form']);

  const presets = getTagPresets(language);
  const options = presets.map(toOption);
  const selected = value.map(toOption);

  const handleChange = (opts: readonly Option[]) => {
    onChange(
      opts.map(o => {
        const preset = presets.find(p => p.name === o.value);
        // Preset tags keep their color; user-created tags have none and render
        // as a neutral chip.
        return {
          name: o.label,
          color: preset?.color ?? (o.color as string) ?? ''
        };
      })
    );
  };

  return (
    <CreatableSelect
      isMulti
      options={options}
      value={selected}
      styles={tagColorStyles}
      placeholder={t('form:add-tags')}
      onChange={opts => handleChange(opts)}
      onCreateOption={name => onChange([...value, { name, color: '' }])}
    />
  );
};

export default TagSelect;
