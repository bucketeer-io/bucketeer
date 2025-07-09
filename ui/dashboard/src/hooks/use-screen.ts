import { useMediaQuery } from 'react-responsive';

export const useScreen = () => {
  const fromTabletScreen = useMediaQuery({ minWidth: 768 });
  const fromXLScreen = useMediaQuery({ minWidth: 1440 });
  const from2XLScreen = useMediaQuery({ minWidth: 1600 });
  const from3XLScreen = useMediaQuery({ minWidth: 1920 });
  const from4XLScreen = useMediaQuery({ minWidth: 2560 });
  const lessThanXLScreen = useMediaQuery({ maxWidth: 1439 });

  return {
    fromTabletScreen,
    fromXLScreen,
    from2XLScreen,
    from3XLScreen,
    from4XLScreen,

    lessThanXLScreen
  };
};
