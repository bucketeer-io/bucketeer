import { authenticationUrl } from '@api/auth';
import { useQueryDemoSiteStatus } from '@queries/demo-site-status';
import { useAuth } from 'auth';
import { urls } from 'configs';
import { setCookieState } from 'cookie';
import { useSubmit } from 'hooks';
import { useTranslation, getLanguage, Language, setLanguage } from 'i18n';
import { clearDemoTokenStorage } from 'storage/demo-token';
import { cn } from 'utils/style';
import { IconGoogle, IconEnglishFlag } from '@icons';
import { languageList } from 'pages/members/member-modal/add-member-modal';
import AuthWrapper from 'pages/signin/elements/auth-wrapper';
import Button from 'components/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import Icon from 'components/icon';
import FormLoading from 'elements/form-loading';

const AccessDemoPage = () => {
  const { t } = useTranslation(['common', 'auth', 'message']);
  const language = getLanguage();
  const { demoGoogleAuthError } = useAuth();

  const { data: demoSiteStatusData, isLoading } = useQueryDemoSiteStatus();

  const isDemoSiteEnabled = demoSiteStatusData?.isDemoSiteEnabled;

  const { onSubmit: onGoogleLoginHandler, submitting } = useSubmit(() => {
    const state = `${Date.now()}`;
    setCookieState(state);
    clearDemoTokenStorage();

    return authenticationUrl({
      state,
      redirectUrl: urls.AUTH_DEMO_REDIRECT,
      type: 2 // Google auth type
    }).then(response => {
      if (response.url) {
        window.location.href = response.url;
      }
    });
  });

  return (
    <AuthWrapper>
      {isLoading ? (
        <FormLoading />
      ) : (
        <>
          <div className="flex items-center justify-between gap-2 mt-6">
            <h1 className="text-gray-900 typo-head-bold-huge">
              {t('auth:demo')}
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
          {demoGoogleAuthError && (
            <div className="typo-para-medium text-accent-red-500 first-letter:uppercase mt-6">
              {demoGoogleAuthError.split('environment:')}
            </div>
          )}
          <div
            className={cn('text-gray-600 typo-para-medium mt-8', {
              'text-accent-red-500': !isDemoSiteEnabled
            })}
          >
            {t(
              isDemoSiteEnabled
                ? 'message:demo-available'
                : 'message:demo-not-available'
            )}
          </div>
          {isDemoSiteEnabled && (
            <Button
              loading={submitting}
              onClick={onGoogleLoginHandler}
              variant={'secondary-2'}
              className="w-full mt-6"
            >
              <Icon icon={IconGoogle} />
              {t('auth:sign-in-with-google')}
            </Button>
          )}
        </>
      )}
    </AuthWrapper>
  );
};

export default AccessDemoPage;
