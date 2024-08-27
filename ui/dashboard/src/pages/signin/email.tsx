import { useState } from 'react';
import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import {
  IconRemoveRedEyeOutlined,
  IconVisibilityOffOutlined
} from 'react-icons-material-design';
import { Link, useNavigate } from 'react-router-dom';
import { signIn } from '@api/auth';
import { yupResolver } from '@hookform/resolvers/yup';
import { useAuth } from 'auth';
import { PAGE_PATH_ROOT } from 'constants/routing';
import * as yup from 'yup';
import { SignInForm } from '@types';
import { IconBackspace } from '@icons';
import { Button } from 'components/button';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import InputGroup from 'components/input-group';
import AuthWrapper from './elements/auth-wrapper';

const formSchema = yup.object().shape({
  email: yup.string().email().required(),
  password: yup
    .string()
    .required()
    .min(4, 'The password you have provided must have at least 4 characters')
});

const SignInWithEmail = () => {
  const { syncSignIn } = useAuth();
  const navigate = useNavigate();

  const [showPassword, setShowPassword] = useState(false);
  const [showAuthError, setShowAuthError] = useState(false);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      email: '',
      password: ''
    }
  });

  const onSubmit: SubmitHandler<SignInForm> = values => {
    setShowAuthError(false);
    return signIn(values)
      .then(response => {
        syncSignIn(response.token);
        navigate(PAGE_PATH_ROOT);
      })
      .catch(error => {
        if (error) {
          setShowAuthError(true);
        }
      });
  };

  const PasswordAddonAction = () => (
    <Button
      type="button"
      variant="grey"
      className="text-gray-500 size-6"
      onClick={() => setShowPassword(!showPassword)}
    >
      <Icon
        icon={
          showPassword ? IconVisibilityOffOutlined : IconRemoveRedEyeOutlined
        }
        size="sm"
      />
    </Button>
  );

  return (
    <AuthWrapper>
      <Button
        variant="secondary-2"
        onClick={() => navigate(PAGE_PATH_ROOT)}
        className="p-2 h-auto"
      >
        <Icon icon={IconBackspace} size="sm" />
      </Button>
      <h1 className="text-gray-900 typo-head-bold-huge mt-8">{`Sign in`}</h1>
      <p className="text-gray-600 typo-para-medium mt-4">
        {`To access our Demo site, please sign in using the follow in information.`}
      </p>
      <div className="text-gray-600 typo-para-medium mt-6">
        <p>{`Email: demo@bucketeer.io`}</p>
        <p>{`Password: demo`}</p>
      </div>

      {showAuthError && (
        <p className="text-accent-red-500 typo-para-medium mt-6">
          {`Wrong email or password. Try again or`}
          <Link to={PAGE_PATH_ROOT}>
            <span className="underline ml-1">{`create an account.`}</span>
          </Link>
        </p>
      )}

      <FormProvider {...form}>
        <Form
          onSubmit={form.handleSubmit(onSubmit)}
          onChange={() => setShowAuthError(false)}
          className="mt-8"
        >
          <Form.Field
            control={form.control}
            name="email"
            render={({ field }) => (
              <Form.Item>
                <Form.Label>{`Email`}</Form.Label>
                <Form.Control>
                  <Input placeholder="Email" {...field} />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Form.Field
            control={form.control}
            name="password"
            render={({ field }) => (
              <Form.Item>
                <Form.Label>{`Password`}</Form.Label>
                <Form.Control>
                  <InputGroup
                    className="w-full"
                    addonSlot="right"
                    addon={<PasswordAddonAction />}
                  >
                    <Input
                      type={showPassword ? 'text' : 'password'}
                      placeholder="Password"
                      {...field}
                    />
                  </InputGroup>
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Button
            type="submit"
            loading={form.formState.isSubmitting}
            className="mt-8 w-full"
          >
            {`Sign In`}
          </Button>
        </Form>
      </FormProvider>
    </AuthWrapper>
  );
};

export default SignInWithEmail;
