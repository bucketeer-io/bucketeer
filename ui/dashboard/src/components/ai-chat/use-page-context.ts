import { useMemo } from 'react';
import { useLocation, useParams } from 'react-router-dom';
import {
  PAGE_PATH_FEATURES,
  PAGE_PATH_EXPERIMENTS,
  PAGE_PATH_USER_SEGMENTS,
  PAGE_PATH_FEATURE_AUTOOPS
} from 'constants/routing';
import { PageContext, PageType } from '@types';

export const usePageContext = (): PageContext => {
  const location = useLocation();
  const params = useParams();

  return useMemo(() => {
    const pathname = location.pathname;
    let pageType: PageType | '' = '';
    let featureId: string | undefined;

    if (pathname.includes(PAGE_PATH_FEATURES)) {
      if (params.flagId) {
        pageType = 'targeting';
        featureId = params.flagId;
      } else {
        pageType = 'feature_flags';
      }
    } else if (pathname.includes(PAGE_PATH_EXPERIMENTS)) {
      pageType = 'experiments';
    } else if (pathname.includes(PAGE_PATH_USER_SEGMENTS)) {
      pageType = 'segments';
    } else if (pathname.includes(PAGE_PATH_FEATURE_AUTOOPS)) {
      pageType = 'autoops';
    }

    return { pageType, featureId };
  }, [location.pathname, params.flagId]);
};
