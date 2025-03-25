import { useMediaQuery } from 'react-responsive';

export const useScreen = () => {
  const isXLScreen = useMediaQuery({ minWidth: 1440, maxWidth: 1599 });
  const isXXLScreen = useMediaQuery({ minWidth: 1600 });

  return {
    isXLScreen,
    isXXLScreen
  };
};
