import { useNavigate } from 'react-router-dom';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_FEATURES } from 'constants/routing';
import { useTranslation } from 'i18n';
import PageDetailsHeader from 'elements/page-details-header';
import PageContent from './page-content';

const PageLoader = () => {
  const { t } = useTranslation(['common']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const navigate = useNavigate();
  return (
    <div className="w-full min-h-screen !max-w-[1192px]">
      <PageDetailsHeader
        onBack={() =>
          navigate(`/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}`)
        }
        title={t('new-flag')}
      />
      <PageContent />
    </div>
  );
};

export default PageLoader;
