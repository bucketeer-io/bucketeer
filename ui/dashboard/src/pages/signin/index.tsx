import { useNavigate } from 'react-router-dom';
import { authenticationUrl } from '@api/auth';
import { useAuth } from 'auth';
import { urls } from 'configs';
import { DEMO_SIGN_IN_ENABLED } from 'configs';
import { PAGE_PATH_AUTH_SIGNIN } from 'constants/routing';
import { setCookieState } from 'cookie';
import { useSubmit } from 'hooks';
import { useTranslation, getLanguage, Language, setLanguage } from 'i18n';
import { IconEmail, IconEnglishFlag, IconGoogle } from '@icons';
import { languageList } from 'pages/members/member-modal/add-member-modal';
import Button from 'components/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
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

  return (
    <AuthWrapper>
      <div className="grid gap-6">
        <div className="flex items-center justify-between gap-2">
          <h1 className="text-gray-900 typo-head-bold-huge">
            {t(`auth:sign-in-to-bucketeer`)}
          </h1>
          <DropdownMenu>
            <DropdownMenuTrigger
              trigger={
                <div className="flex items-center justify-between w-full">
                  <div className="flex items-center w-full gap-x-2">
                    <div className="flex-center size-fit mt-0.5">
                      <Icon
                        color="primary-50"
                        size="sm"
                        icon={
                          languageList.find(item => item.value === language)
                            ?.icon || IconEnglishFlag
                        }
                      />
                    </div>
                    {languageList.find(item => item.value === language)
                      ?.label || 'English'}
                  </div>
                </div>
              }
              className="bg-transparent !shadow-none !border-none"
            />
            <DropdownMenuContent side="bottom" align="start">
              {languageList?.map((item, index) => (
                <DropdownMenuItem
                  key={index}
                  label={item.label}
                  value={item.value}
                  icon={item?.icon}
                  iconElement={
                    item?.icon ? (
                      <div className="flex-center size-fit mt-0.5">
                        <Icon size="sm" icon={item?.icon} />
                      </div>
                    ) : null
                  }
                  className="[&>div>button]:!cursor-pointer"
                  onSelectOption={value => setLanguage(value as Language)}
                />
              ))}
            </DropdownMenuContent>
          </DropdownMenu>
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
