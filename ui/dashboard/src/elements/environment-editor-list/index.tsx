import { useCallback, useMemo } from 'react';
import { getEditorEnvironments, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { onFormatEnvironments } from 'utils/function';
import DropdownMenuWithSearch, {
  DropdownMenuWithSearchProps
} from 'elements/dropdown-with-search';

interface Props extends Omit<DropdownMenuWithSearchProps, 'options'> {
  value: string | string[];
  selectedValues?: string[];
  currentEnvironmentId?: string;
}

const EnvironmentEditorList = ({
  value,
  placeholder,
  itemSize = 60,
  maxOptions = 10,
  selectedValues,
  currentEnvironmentId,
  ...props
}: Props) => {
  const { t } = useTranslation(['form', 'common']);
  const { consoleAccount } = useAuth();
  const { editorEnvironments, projects } = getEditorEnvironments(
    consoleAccount!
  );
  const { formattedEnvironments } = onFormatEnvironments(editorEnvironments);

  const environmentOptions = useMemo(() => {
    const options = formattedEnvironments
      .filter(env => env.id !== currentEnvironmentId)
      .map(item => ({
        label: `${item.name}`,
        value: item.id,
        description: `${t('common:source-type.project')}: ${projects.find(project => project.id === item.projectId)?.name}`,
        projectId: item.projectId
      }));
    return options;
  }, [formattedEnvironments, currentEnvironmentId, projects, t]);

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

  const getEnvLabel = useCallback(
    (value: string) => {
      const selectedEnv = environmentOptions.find(env => env.value === value);
      return selectedEnv
        ? `${selectedEnv.label} (${selectedEnv.description})`
        : value;
    },
    [environmentOptions]
  );

  const environmentLabel = useMemo(() => {
    if (Array.isArray(value)) {
      return value.map(item => getEnvLabel(item)).join(', ');
    }
    if (typeof value === 'string') {
      return getEnvLabel(value);
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
      maxOptions={maxOptions}
      itemClassName="!p-1.5 !mb-0"
      {...props}
    />
  );
};

export default EnvironmentEditorList;
