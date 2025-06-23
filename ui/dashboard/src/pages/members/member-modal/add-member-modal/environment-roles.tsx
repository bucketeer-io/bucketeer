import { useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { IconAddOutlined } from 'react-icons-material-design';
import { EnvironmentRoleItem } from '@api/account/account-creator';
import useOptions from 'hooks/use-options';
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
import { AddMemberForm } from '.';
import { EditMemberForm } from '../edit-member-modal';

const EnvironmentRoles = ({
  environments
}: {
  environments: Environment[];
}) => {
  const { t } = useTranslation(['common', 'form']);
  const { environmentRoleOptions } = useOptions();
  const methods = useFormContext<AddMemberForm | EditMemberForm>();
  const { control, watch, setValue } = methods;

  const environmentRolesWatch: EnvironmentRoleItem[] =
    watch('environmentRoles');

  const selectedEnvs = environmentRolesWatch.map(item => item.environmentId);
  const environmentsOptions = environments.filter(
    item => item.id && !selectedEnvs.includes(item.id)
  );

  const isDisabledAddMemberButton = useMemo(
    () =>
      environmentRolesWatch?.length >=
        environments?.filter(item => item.id).length ||
      !environmentsOptions.length,
    [environmentRolesWatch, environments, environmentsOptions]
  );

  const onAddEnvironment = () => {
    const newEnvironmentRoles: EnvironmentRoleItem[] = [
      ...environmentRolesWatch,
      {
        environmentId: '',
        role: 'Environment_UNASSIGNED'
      }
    ];

    setValue('environmentRoles', newEnvironmentRoles, { shouldDirty: true });
  };

  const onDeleteEnvironment = (itemIndex: number) => {
    const newEnvironmentRoles: EnvironmentRoleItem[] =
      environmentRolesWatch.filter((_, index) => index !== itemIndex);
    setValue('environmentRoles', newEnvironmentRoles, { shouldDirty: true });
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

      {environmentRolesWatch?.map((environment, envIndex) => (
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
                              setValue(
                                `environmentRoles.${envIndex}.environmentId`,
                                value as string,
                                {
                                  shouldDirty: true,
                                  shouldValidate: true
                                }
                              );
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
                              setValue(
                                `environmentRoles.${envIndex}.role`,
                                value as EnvironmentRoleType,
                                {
                                  shouldDirty: true,
                                  shouldValidate: true
                                }
                              );
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
          {environmentRolesWatch.length > 1 && (
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
