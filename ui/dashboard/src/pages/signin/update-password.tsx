import { useCallback, useState } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import {
  IconRemoveRedEyeOutlined,
  IconVisibilityOffOutlined
} from 'react-icons-material-design';
import { useNavigate } from 'react-router-dom';
import { updatePassword } from '@api/auth/setup-password';
import { yupResolver } from '@hookform/resolvers/yup';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_AUTH_SIGNIN } from 'constants/routing';
import { useToast } from 'hooks';
import useFormSchema, { FormSchemaProps } from 'hooks/use-form-schema';
import { useTranslation } from 'i18n';
import * as yup from 'yup';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import InputGroup from 'components/input-group';

const formSchema = ({ requiredMessage, translation }: FormSchemaProps) =>
  yup.object().shape({
    currentPassword: yup
      .string()
      .required(requiredMessage)
      .min(
        2,
        translation('message:validation.name-at-least-characters', {
          count: 2,
          name: translation('auth:current-password').toLowerCase()
        })
      ),
    newPassword: yup
      .string()
      .required(requiredMessage)
      .min(
        2,
        translation('message:validation.name-at-least-characters', {
          count: 2,
          name: translation('auth:new-password').toLowerCase()
        })
      ),
    confirmPassword: yup
      .string()
      .oneOf(
        [yup.ref('newPassword'), ''],
        translation('auth:password-confirm-not-match')
      )
      .required(requiredMessage)
  });

interface UserInfoForm {
  currentPassword: string;
  newPassword: string;
}
interface UpdatePasswordI {
  className?: string;
  onBack?: () => void;
  onCancel?: () => void;
}

const UpdatePassword = ({
  className,
  onBack,
  onCancel,
  ...props
}: UpdatePasswordI) => {
  const { t } = useTranslation(['auth', 'form', 'common', 'message']);
  const { consoleAccount } = useAuth();
  const [showCurrentPassword, setShowCurrentPassword] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [showconfirmPassword, setShowconfirmPassword] = useState(false);
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const { notify, errorNotify } = useToast();
  const navigate = useNavigate();

  const form = useForm({
    resolver: yupResolver(useFormSchema(formSchema)),
    defaultValues: {
      currentPassword: '',
      newPassword: ''
    },
    mode: 'onSubmit'
  });

  // const PasswordAddonAction = ({
  //   type
  // }: {
  //   type: 'password' | 'confirm-password' | 'current-password';
  // }) => (
  //   <Button
  //     key={type}
  //     type="button"
  //     variant="grey"
  //     className="text-gray-500 size-6"
  //     onClick={() =>
  //       type === 'password'
  //         ? setShowPassword(!showPassword)
  //         : type === 'current-password'
  //           ? setShowCurrentPassword(!showCurrentPassword)
  //           : setShowconfirmPassword(!showconfirmPassword)
  //     }
  //   >
  //     <Icon
  //       icon={
  //         (showPassword && type === 'password') ||
  //         (showCurrentPassword && type === 'current-password') ||
  //         (showconfirmPassword && type === 'confirm-password')
  //           ? IconVisibilityOffOutlined
  //           : IconRemoveRedEyeOutlined
  //       }
  //       size="sm"
  //     />
  //   </Button>
  // );

  const onSubmit: SubmitHandler<UserInfoForm> = useCallback(
    async values => {
      try {
        if (consoleAccount) {
          const { currentPassword, newPassword } = values;
          const resp = await updatePassword({
            currentPassword: currentPassword,
            newPassword: newPassword
          });

          if (resp) {
            notify({
              message: t('message:collection-action-success', {
                collection: t('auth:password'),
                action: t('common:updated')
              })
            });
            navigate(PAGE_PATH_AUTH_SIGNIN);
          }
        }
      } catch (error) {
        errorNotify(error);
      }
    },
    [consoleAccount, currentEnvironment]
  );

  return (
    <div className={className} {...props}>
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(onSubmit)}>
          <div className="p-5">
            <Form.Field
              control={form.control}
              name="currentPassword"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required>{t(`current-password`)}</Form.Label>
                  <Form.Control>
                    <InputGroup
                      className="w-full"
                      addonSlot="right"
                      addon={
                        <Button
                          type="button"
                          variant="grey"
                          className="text-gray-500 size-6"
                          onClick={() =>
                            setShowCurrentPassword(!showCurrentPassword)
                          }
                        >
                          <Icon
                            icon={
                              showCurrentPassword
                                ? IconVisibilityOffOutlined
                                : IconRemoveRedEyeOutlined
                            }
                            size="sm"
                          />
                        </Button>
                      }
                    >
                      <Input
                        type={showCurrentPassword ? 'text' : 'password'}
                        placeholder={t(`current-password-placeholder`)}
                        {...field}
                        onChange={value => field.onChange(value)}
                      />
                    </InputGroup>
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name="newPassword"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required>{t(`new-password`)}</Form.Label>
                  <Form.Control>
                    <InputGroup
                      className="w-full"
                      addonSlot="right"
                      addon={
                        <Button
                          type="button"
                          variant="grey"
                          className="text-gray-500 size-6"
                          onClick={() => setShowPassword(!showPassword)}
                        >
                          <Icon
                            icon={
                              showPassword
                                ? IconVisibilityOffOutlined
                                : IconRemoveRedEyeOutlined
                            }
                            size="sm"
                          />
                        </Button>
                      }
                    >
                      <Input
                        type={showPassword ? 'text' : 'password'}
                        placeholder={t(`new-password-placeholder`)}
                        {...field}
                        onChange={value => field.onChange(value)}
                      />
                    </InputGroup>
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name="confirmPassword"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required>{t(`password-confirm`)}</Form.Label>
                  <Form.Control>
                    <InputGroup
                      className="w-full"
                      addonSlot="right"
                      addon={
                        <Button
                          type="button"
                          variant="grey"
                          className="text-gray-500 size-6"
                          onClick={() =>
                            setShowconfirmPassword(!showconfirmPassword)
                          }
                        >
                          <Icon
                            icon={
                              showconfirmPassword
                                ? IconVisibilityOffOutlined
                                : IconRemoveRedEyeOutlined
                            }
                            size="sm"
                          />
                        </Button>
                      }
                    >
                      <Input
                        type={showconfirmPassword ? 'text' : 'password'}
                        placeholder={t(`password-confirm-placeholder`)}
                        {...field}
                        onChange={value => field.onChange(value)}
                      />
                    </InputGroup>
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <div className="w-full flex justify-start">
              <Button
                type="button"
                className="bg-transparent border-0 text-primary-500 hover:bg-transparent p-0"
                onClick={() => onBack?.()}
              >
                {t('common:change-profile')}
              </Button>
            </div>
          </div>
          <ButtonBar
            secondaryButton={
              <Button disabled={!form.formState.isValid}>
                {t(`common:save`)}
              </Button>
            }
            primaryButton={
              <Button type="button" variant="secondary" onClick={onCancel}>
                {t(`common:cancel`)}
              </Button>
            }
          />
        </Form>
      </FormProvider>
    </div>
  );
};

export default UpdatePassword;
