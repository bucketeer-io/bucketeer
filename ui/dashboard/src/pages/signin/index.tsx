import { useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { authenticationUrl } from '@api/auth';
import { useAuth } from 'auth';
import { DEMO_SIGN_IN_ENABLED, urls } from 'configs';
import { PAGE_PATH_AUTH_SIGNIN } from 'constants/routing';
import { setCookieState } from 'cookie';
import { useSubmit } from 'hooks';
import { getLanguage, Language, setLanguage, useTranslation } from 'i18n';
import { IconEmail, IconEnglishFlag, IconGoogle } from '@icons';
import { languageList } from 'pages/members/member-modal/add-member-modal';
import Button from 'components/button';
import Dropdown from 'components/dropdown';
import Icon from 'components/icon';
import AuthWrapper from './elements/auth-wrapper';

const SignIn = () => {
  const { t } = useTranslation(['auth']);
  const navigate = useNavigate();
  const { isGoogleAuthError, setIsGoogleAuthError } = useAuth();
  const language = getLanguage();

  const { onSubmit: onGoogleLoginHandler, submitting } = useSubmit(() => {
    const state = `${Date.now()}`;
    setCookieState(state);
    setIsGoogleAuthError(false);

    return authenticationUrl({
      state,
      redirectUrl: urls.AUTH_REDIRECT,
      type: 2 // Google auth type
    }).then(response => {
      if (response.url) {
        window.location.href = response.url;
      }
    });
  });

  const optionLanguages = useMemo(() => {
    return languageList.map(item => ({
      label: (
        <div className="flex items-center justify-between w-full">
          <div className="flex items-center w-full gap-x-2">
            <div className="flex-center size-fit mt-0.5">
              <Icon
                color="primary-50"
                size="sm"
                icon={item?.icon || IconEnglishFlag}
              />
            </div>

            {item?.label || 'English'}
          </div>
        </div>
      ),
      value: item.value
    }));
  }, [languageList]);
  return (
    <AuthWrapper>
      <div className="grid gap-6">
        <div className="flex items-center justify-between gap-2">
          <h1 className="text-gray-900 typo-head-bold-huge">
            {t(`auth:sign-in-to-bucketeer`)}
          </h1>
          <Dropdown
            value={language}
            options={optionLanguages}
            onChange={value => setLanguage(value as Language)}
            className="w-fit bg-transparent !shadow-none !border-none"
            wrapTriggerStyle="!w-fit"
            menuContentSide="bottom"
            itemClassName="[&>div>button]:!cursor-pointer"
          />
        </div>
        {isGoogleAuthError && (
          <p className="text-accent-red-500 typo-para-medium">
            {t(`error-message.invalid-google-auth`)}
          </p>
        )}
      </div>
      <div className="flex flex-col gap-4 mt-10">
        {DEMO_SIGN_IN_ENABLED && (
          <Button
            variant="secondary-2"
            onClick={() => {
              navigate(PAGE_PATH_AUTH_SIGNIN);
              setIsGoogleAuthError(false);
            }}
          >
            <Icon icon={IconEmail} />
            {t(`auth:sign-in-with-email`)}
          </Button>
        )}
        <Button
          loading={submitting}
          onClick={onGoogleLoginHandler}
          variant={'secondary-2'}
        >
          <Icon icon={IconGoogle} />
          {t(`auth:sign-in-with-google`)}
        </Button>
        {/* <Button variant={'secondary-2'}>
          <Icon icon={IconGithub} />
          {`Sign in With Github`}
        </Button>
        <Button variant={'secondary-2'}>
          <Icon icon={IconKey} />
          {`Sign in With SSO`}
        </Button> */}
      </div>
    </AuthWrapper>
  );
};

export default SignIn;
