import { useMemo } from 'react';
import { useFormContext } from 'react-hook-form';
import { IconAddOutlined } from 'react-icons-material-design';
import { EnvironmentRoleItem } from '@api/account/account-creator';
import useOptions from 'hooks/use-options';
import { getLanguage, Language, useTranslation } from 'i18n';
import { Environment, EnvironmentRoleType } from '@types';
import { cn } from 'utils/style';
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
import EnvironmentEditorList from 'elements/environment-editor-list';
import { AddMemberForm } from '.';
import { EditMemberForm } from '../edit-member-modal';

const EnvironmentRoles = ({
  environments
}: {
  environments: Environment[];
}) => {
  const { t } = useTranslation(['common', 'form']);
  const { environmentRoleOptions } = useOptions();
  const isJapaneseLanguage = getLanguage() === Language.JAPANESE;

  const methods = useFormContext<AddMemberForm | EditMemberForm>();
  const { control, watch, setValue } = methods;

  const environmentRolesWatch: EnvironmentRoleItem[] =
    watch('environmentRoles') || [];

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
        <div
          key={envIndex}
          className="flex items-start w-full max-w-full gap-x-4"
        >
          <div
            className={cn('flex w-full max-w-[310px]', {
              'max-w-[290px]': isJapaneseLanguage
            })}
          >
            <Form.Field
              control={control}
              name={`environmentRoles.${envIndex}.environmentId`}
              render={({ field }) => (
                <Form.Item className="flex flex-col w-full py-2">
                  <Form.Label required>{t('environment')}</Form.Label>
                  <Form.Control>
                    <EnvironmentEditorList
                      align="start"
                      value={field.value}
                      selectedValues={selectedEnvs}
                      onSelectOption={value =>
                        setValue(
                          `environmentRoles.${envIndex}.environmentId`,
                          value as string,
                          {
                            shouldDirty: true,
                            shouldValidate: true
                          }
                        )
                      }
                      triggerClassName="max-w-full truncate"
                    />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
          </div>

          <div
            className={cn('w-[140px] min-w-[140px] h-full', {
              'w-[160px] min-w-[160px]': isJapaneseLanguage
            })}
          >
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
                        className={
                          isJapaneseLanguage ? 'min-w-[170px]' : 'min-w-[140px]'
                        }
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

          <Button
            variant="grey"
            size="icon"
            type="button"
            disabled={environmentRolesWatch.length <= 1}
            className="p-0 size-5 min-w-5 mt-5 self-center"
            onClick={() => onDeleteEnvironment(envIndex)}
          >
            <Icon icon={IconTrash} size="sm" />
          </Button>
        </div>
      ))}
    </>
  );
};

export default EnvironmentRoles;
