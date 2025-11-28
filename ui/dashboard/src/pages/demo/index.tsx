import { useMemo } from 'react';
import { authenticationUrl } from '@api/auth';
import { useQueryDemoSiteStatus } from '@queries/demo-site-status';
import { useAuth } from 'auth';
import { urls } from 'configs';
import { setCookieState } from 'cookie';
import { useSubmit } from 'hooks';
import { getLanguage, Language, setLanguage, useTranslation } from 'i18n';
import { clearDemoTokenStorage } from 'storage/demo-token';
import { cn } from 'utils/style';
import { IconEnglishFlag, IconGoogle } from '@icons';
import { languageList } from 'pages/members/member-modal/add-member-modal';
import AuthWrapper from 'pages/signin/elements/auth-wrapper';
import Button from 'components/button';
import Dropdown from 'components/dropdown';
import Icon from 'components/icon';
import FormLoading from 'elements/form-loading';

const AccessDemoPage = () => {
  const { t } = useTranslation(['common', 'auth', 'message']);
  const language = getLanguage();
  const { demoGoogleAuthError } = useAuth();

  const { data: demoSiteStatusData, isLoading } = useQueryDemoSiteStatus();

  const isDemoSiteEnabled = demoSiteStatusData?.isDemoSiteEnabled;
  const optionLanguages = useMemo(
    () =>
      languageList?.map(item => ({
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
              {item.label || 'English'}
            </div>
          </div>
        ),
        value: item.value
      })),
    [languageList]
  );

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

            <Dropdown
              value={language}
              options={optionLanguages}
              onChange={value => setLanguage(value as Language)}
              className="bg-transparent !shadow-none !border-none"
              wrapTriggerStyle="!w-fit"
              contentClassName="[&>div>button]:!cursor-pointer"
            />
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
