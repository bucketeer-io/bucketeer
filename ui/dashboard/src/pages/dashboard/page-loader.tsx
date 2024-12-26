import { useEffect, useState } from 'react';
import PageLayout from 'elements/page-layout';
import PageContent from './page-content';

const PageLoader = () => {
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const timerId = setTimeout(() => {
      setIsLoading(false);
    }, 1000);
    return () => clearTimeout(timerId);
  }, []);

  return <>{isLoading ? <PageLayout.LoadingState /> : <PageContent />}</>;
};

export default PageLoader;
