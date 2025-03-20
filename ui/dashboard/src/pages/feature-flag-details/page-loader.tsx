import { useNavigate, useParams } from 'react-router-dom';
import { useQueryFeature } from '@queries/feature-details';
import { getCurrentEnvironment, useAuth } from 'auth';
import { PAGE_PATH_FEATURES } from 'constants/routing';
import PageDetailsHeader from 'elements/page-details-header';
import PageLayout from 'elements/page-layout';
import HeaderDetails from './elements/header-details';
import PageContent from './page-content';

export const mockFlags = [
  {
    id: 'flag-1',
    name: 'Flag using boolean',
    type: 'boolean',
    status: 'active',
    tags: ['Android'],
    variations: [],
    disabled: false,
    operations: [],
    createdAt: '1706182987',
    updatedAt: '1706182994'
  },
  {
    id: 'flag-2',
    name: 'Flag using string',
    type: 'string',
    status: 'no_activity',
    tags: ['Web'],
    variations: [],
    disabled: false,
    operations: [],
    createdAt: '1706182987',
    updatedAt: '1706182994'
  },
  {
    id: 'flag-3',
    name: 'Flag using number',
    type: 'number',
    status: 'new',
    tags: ['Android'],
    variations: [],
    disabled: false,
    operations: [],
    createdAt: '1706182987',
    updatedAt: '1706182994'
  },
  {
    id: 'flag-4',
    name: 'Flag using json',
    type: 'json',
    status: 'no_activity',
    tags: ['IOS'],
    variations: [],
    disabled: false,
    operations: [],
    createdAt: '1706182987',
    updatedAt: '1706182994'
  }
];

const PageLoader = () => {
  const params = useParams();

  const navigate = useNavigate();
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const {
    data: collection,
    isLoading,
    isError,
    refetch
  } = useQueryFeature({
    params: {
      environmentId: currentEnvironment?.id,
      id: params?.flagId || ''
    },
    enabled: !!params?.flagId && !!currentEnvironment?.id
  });

  const feature = collection?.feature;

  const isErrorState = isError || !feature;

  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isErrorState ? (
        <PageLayout.ErrorState onRetry={refetch} />
      ) : (
        <>
          <PageDetailsHeader onBack={() => navigate(`${PAGE_PATH_FEATURES}`)}>
            <HeaderDetails feature={feature} />
          </PageDetailsHeader>
          <PageContent feature={feature} />
        </>
      )}
    </>
  );
};

export default PageLoader;
