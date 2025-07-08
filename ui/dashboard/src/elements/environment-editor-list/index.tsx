import { useMemo } from 'react';
import { getEditorEnvironments, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { onFormatEnvironments } from 'utils/function';
import DropdownMenuWithSearch, {
  DropdownMenuWithSearchProps
} from 'elements/dropdown-with-search';

interface Props extends Omit<DropdownMenuWithSearchProps, 'options'> {
  value: string | string[];
  selectedValues?: string[];
}

const EnvironmentEditorList = ({
  value,
  placeholder,
  itemSize = 40,
  selectedValues,
  ...props
}: Props) => {
  const { t } = useTranslation(['form', 'common']);
  const { consoleAccount } = useAuth();
  const { editorEnvironments, projects } = getEditorEnvironments(
    consoleAccount!
  );
  const { formattedEnvironments } = onFormatEnvironments(editorEnvironments);

  const environmentOptions = useMemo(() => {
    const options = formattedEnvironments.map(item => ({
      label: `${item.name} (${t('common:source-type.project')}: ${projects.find(project => project.id === item.projectId)?.name})`,
      value: item.id,
      projectId: item.projectId
    }));
    return options;
  }, [formattedEnvironments]);

  const remainingEnvironmentsOptions = useMemo(() => {
    const remainingEnvironments = environmentOptions.filter(item =>
      selectedValues
        ? !selectedValues.includes(item.value)
        : Array.isArray(value)
          ? !value.includes(item.value)
          : item.value !== value
    );
    return remainingEnvironments;
  }, [environmentOptions, selectedValues]);

  const environmentLabel = useMemo(() => {
    if (Array.isArray(value)) {
      return value
        .map(
          item =>
            environmentOptions.find(env => env.value === item)?.label || item
        )
        .join(', ');
    }
    if (typeof value === 'string') {
      return environmentOptions.find(env => env.value === value)?.label;
    }
    return '';
  }, [value, environmentOptions]);

  return (
    <DropdownMenuWithSearch
      options={remainingEnvironmentsOptions}
      label={environmentLabel}
      placeholder={placeholder || t('select-environment')}
      isMultiselect={Array.isArray(value)}
      selectedOptions={Array.isArray(value) ? value : undefined}
      itemSize={itemSize}
      {...props}
    />
  );
};

export default EnvironmentEditorList;
