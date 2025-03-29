import { useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
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

const environmentRoleOptions: environmentRoleOption[] = [
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
  onChangeEnvironments
}: {
  environments: Environment[];
  memberEnvironments: EnvironmentRoleItem[];
  onChangeEnvironments: (v: EnvironmentRoleItem[]) => void;
}) => {
  const { t } = useTranslation(['common', 'form']);

  const methods = useFormContext();
  const { control, watch } = methods;

  const selectedEnvs = memberEnvironments.map(item => item.environmentId);
  const environmentsOptions = environments.filter(
    item => item.id && !selectedEnvs.includes(item.id)
  );

  const environmentRoles = watch('environmentRoles');

  const isDisabledAddMemberButton = useMemo(
    () =>
      environmentRoles?.length >=
        environments?.filter(item => item.id).length ||
      !environmentsOptions.length,
    [environmentRoles, environments, environmentsOptions]
  );

  const onAddEnvironment = () => {
    memberEnvironments.push({
      environmentId: '',
      role: 'Environment_UNASSIGNED'
    });
    onChangeEnvironments([...memberEnvironments]);
  };

  const onDeleteEnvironment = (itemIndex: number) => {
    const environments = memberEnvironments.filter(
      (_item, index) => itemIndex !== index
    );
    onChangeEnvironments([...environments]);
  };

  return (
    <>
      <p className="text-gray-800 typo-head-bold-small">{t('environment')}</p>
      <Button
        onClick={onAddEnvironment}
        variant="text"
        type="button"
        className="my-1"
        disabled={isDisabledAddMemberButton}
      >
        <Icon icon={IconAddOutlined} />
        {t(`add-environment`)}
      </Button>

      {memberEnvironments.map((environment, envIndex) => (
        <div key={envIndex} className="flex items-start w-full gap-x-4">
          <div className="flex-1">
            <Form.Field
              control={control}
              name={`environmentRoles.${envIndex}.environmentId`}
              render={({ field }) => (
                <Form.Item className="py-2">
                  <Form.Label required>{t('environment')}</Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        placeholder={t(`form:select-environment`)}
                        label={
                          // TODO: remove empty id when the backend is fixed
                          environments.find(
                            item =>
                              item.id && item.id === environment.environmentId
                          )?.name
                        }
                        variant="secondary"
                        className="w-full"
                      />
                      <DropdownMenuContent
                        className="w-[310px]"
                        align="start"
                        {...field}
                      >
                        {environmentsOptions.map((item, index) => (
                          <DropdownMenuItem
                            {...field}
                            key={index}
                            value={item.id}
                            label={item.name}
                            onSelectOption={value => {
                              field.onChange(value);
                            }}
                          />
                        ))}
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
          </div>

          <div className="w-[140px] h-full">
            <Form.Field
              control={control}
              name={`environmentRoles.${envIndex}.role`}
              render={({ field }) => (
                <Form.Item className="py-2">
                  <Form.Label required>{t('role')}</Form.Label>
                  <Form.Control>
                    <DropdownMenu>
                      <DropdownMenuTrigger
                        placeholder={t(`form:select-role`)}
                        label={
                          environmentRoleOptions.find(
                            item => item.value === environment.role
                          )?.label
                        }
                        variant="secondary"
                        className="w-full"
                      />
                      <DropdownMenuContent
                        className="min-w-[140px]"
                        align="start"
                        {...field}
                      >
                        {environmentRoleOptions.map((item, index) => (
                          <DropdownMenuItem
                            {...field}
                            key={index}
                            value={item.value}
                            label={item.label}
                            onSelectOption={value => {
                              field.onChange(value);
                            }}
                          />
                        ))}
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
          </div>
          {memberEnvironments.length > 1 && (
            <Button
              variant="text"
              size="icon"
              type="button"
              className="p-0 size-5 mt-5 self-center"
              onClick={() => onDeleteEnvironment(envIndex)}
            >
              <Icon icon={IconTrash} size="sm" color="gray-600" />
            </Button>
          )}
        </div>
      ))}
    </>
  );
};

export default EnvironmentRoles;
