import { useMediaQuery } from 'react-responsive';

export const useScreen = () => {
  // true when width >= 600px (tablet/desktop); use !fromMobileScreen to detect mobile
  const fromMobileScreen = useMediaQuery({ minWidth: 600 });
  const fromTabletScreen = useMediaQuery({ minWidth: 900 });
  const fromXLScreen = useMediaQuery({ minWidth: 1440 });
  const from2XLScreen = useMediaQuery({ minWidth: 1600 });
  const from3XLScreen = useMediaQuery({ minWidth: 1920 });
  const from4XLScreen = useMediaQuery({ minWidth: 2560 });
  const lessThanXLScreen = useMediaQuery({ maxWidth: 1439 });
  const isMobile = !fromMobileScreen;

  return {
    isMobile,
    fromMobileScreen,
    fromTabletScreen,
    fromXLScreen,
    from2XLScreen,
    from3XLScreen,
    from4XLScreen,
    lessThanXLScreen
  };
};
