import { useEffect, useState } from 'react';
// import { useToggleOpen } from 'hooks/use-toggle-open';
// import { Project } from '@types';
import PageLayout from 'elements/page-layout';
import PageContent from './page-content';

const PageLoader = () => {
  const [isLoading, setIsLoading] = useState(true);
  //   const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
  //     useToggleOpen(false);
  //   const [isOpenEditModal, onOpenEditModal, onCloseEditModal] =
  //     useToggleOpen(false);

  //   const isEmpty = collection?.projects.length === 0;

  useEffect(() => {
    const timerId = setTimeout(() => {
      setIsLoading(false);
    }, 1000);
    return () => clearTimeout(timerId);
  }, []);

  return <>{isLoading ? <PageLayout.LoadingState /> : <PageContent />}</>;
};

export default PageLoader;
