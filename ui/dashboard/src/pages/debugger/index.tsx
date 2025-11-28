import { Trans } from 'react-i18next';
import { IconLaunchOutlined } from 'react-icons-material-design';
import { Link } from 'react-router-dom';
import { DOCUMENTATION_LINKS } from 'constants/documentation-links';
import { useTranslation } from 'i18n';
import Icon from 'components/icon';
import PageLayout from 'elements/page-layout';
import PageLoader from './page-loader';

const DebuggerPage = () => {
  const { t } = useTranslation(['common', 'form', 'table']);

  return (
    <PageLayout.Root title={t('navigation.debugger')}>
      <PageLayout.Header>
        <h1 className="text-gray-900 typo-head-bold-huge">
          {t('navigation.debugger')}
        </h1>
        <p className="text-gray-600 mt-3 typo-para-small">
          <Trans
            i18nKey={'common:debugger-subtitle'}
            components={{
              comp: (
                <Link
                  to={DOCUMENTATION_LINKS.DEBUGGER}
                  target="_blank"
                  className="inline-flex items-center gap-x-1 text-primary-500 underline"
                />
              ),
              icon: <Icon icon={IconLaunchOutlined} size="sm" />
            }}
          />
        </p>
      </PageLayout.Header>
      <PageLoader />
    </PageLayout.Root>
  );
};

export default DebuggerPage;
