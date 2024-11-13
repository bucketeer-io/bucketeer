import { IconAddOutlined } from 'react-icons-material-design';
import { EnvironmentRoleItem } from '@api/account/account-creator';
import { useTranslation } from 'i18n';
import { Environment, EnvironmentRoleType } from '@types';
import { IconTrash } from '@icons';
import Button from 'components/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Form from 'components/form';
import Icon from 'components/icon';

interface environmentRoleOption {
  value: EnvironmentRoleType;
  label: string;
}

const environmentRoles: environmentRoleOption[] = [
  {
    value: 'Environment_EDITOR',
    label: 'Editor'
  },
  {
    value: 'Environment_VIEWER',
    label: 'Viewer'
  }
];

const EnvironmentRoles = ({
  environments,
  memberEnvironments,
  setMemberEnvironments
}: {
  environments: Environment[];
  memberEnvironments: EnvironmentRoleItem[];
  setMemberEnvironments: (v: EnvironmentRoleItem[]) => void;
}) => {
  const { t } = useTranslation(['common', 'form']);

  const selectedEnvs = memberEnvironments.map(item => item.environmentId);
  const environmentsOptions = environments.filter(
    item => !selectedEnvs.includes(item.id)
  );

  const onAddEnvironment = () => {
    memberEnvironments.push({
      environmentId: '',
      role: 'Environment_UNASSIGNED'
    });
    setMemberEnvironments([...memberEnvironments]);
  };

  const onDeleteEnvironment = (itemIndex: number) => {
    setMemberEnvironments(
      memberEnvironments.filter((_item, index) => itemIndex !== index)
    );
  };

  return (
    <>
      <p className="text-gray-800 typo-head-bold-small">{t('environment')}</p>
      <Button
        onClick={onAddEnvironment}
        variant="text"
        type="button"
        className="my-1"
      >
        <Icon icon={IconAddOutlined} />
        {t(`add-environment`)}
      </Button>

      {memberEnvironments.map((environment, envIndex) => (
        <div key={envIndex} className="flex items-center w-full pb-4 gap-x-4">
          <div className="flex-1">
            <Form.Label required>{t('environment')}</Form.Label>
            <DropdownMenu>
              <DropdownMenuTrigger
                placeholder={t(`form:select-environment`)}
                label={
                  environments.find(
                    item => item.id === environment.environmentId
                  )?.name
                }
                variant="secondary"
                className="w-full"
              />
              <DropdownMenuContent className="w-[310px]" align="start">
                {environmentsOptions.map((item, index) => (
                  <DropdownMenuItem
                    key={index}
                    value={item.id}
                    label={item.name}
                    onSelectOption={() => {
                      memberEnvironments[envIndex].environmentId = item.id;
                      setMemberEnvironments([...memberEnvironments]);
                    }}
                  />
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
          </div>

          <div className="w-[140px]">
            <Form.Label required>{t('role')}</Form.Label>
            <DropdownMenu>
              <DropdownMenuTrigger
                placeholder={t(`form:select-role`)}
                label={
                  environmentRoles.find(item => item.value === environment.role)
                    ?.label
                }
                variant="secondary"
                className="w-full"
              />
              <DropdownMenuContent className="min-w-[140px]" align="start">
                {environmentRoles.map((item, index) => (
                  <DropdownMenuItem
                    key={index}
                    value={item.value}
                    label={item.label}
                    onSelectOption={() => {
                      memberEnvironments[envIndex].role = item.value;
                      setMemberEnvironments([...memberEnvironments]);
                    }}
                  />
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
          </div>

          <Button
            variant="text"
            size="icon"
            type="button"
            className="p-0 size-5 mt-5"
            onClick={() => onDeleteEnvironment(envIndex)}
          >
            <Icon icon={IconTrash} size="sm" color="gray-600" />
          </Button>
        </div>
      ))}
    </>
  );
};

export default EnvironmentRoles;
