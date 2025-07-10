import { OrganizationDetailsResponse } from '@api/organization';
import { QueryObserverResult } from '@tanstack/react-query';
import { Organization } from '@types';
import PageLayout from 'elements/page-layout';
import PageContent from './page-content';

interface Props {
  isLoading: boolean;
  isError: boolean;
  organization?: Organization;
  onRetry: () => Promise<
    QueryObserverResult<OrganizationDetailsResponse, Error>
  >;
}

const PageLoader = ({ isLoading, isError, organization, onRetry }: Props) => {
  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isError || !organization ? (
        <PageLayout.ErrorState onRetry={onRetry} />
      ) : (
        <PageContent organization={organization} />
      )}
    </>
  );
};

export default PageLoader;
