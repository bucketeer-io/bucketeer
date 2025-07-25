import { demoSiteStatusFetcher } from '@api/auth';
import { useQuery } from '@tanstack/react-query';

export const DEMO_SITE_STATUS_KEY = 'demo-site-status';

export const useQueryDemoSiteStatus = () => {
  const query = useQuery({
    queryKey: [DEMO_SITE_STATUS_KEY],
    queryFn: async () => {
      return demoSiteStatusFetcher();
    },
    gcTime: 0
  });
  return query;
};
