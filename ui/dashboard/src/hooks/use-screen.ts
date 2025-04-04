import { useMediaQuery } from 'react-responsive';

export const useScreen = () => {
  const fromXLScreen = useMediaQuery({ minWidth: 1440 });
  const from2XLScreen = useMediaQuery({ minWidth: 1600 });
  const from3XLScreen = useMediaQuery({ minWidth: 1920 });
  const from4XLScreen = useMediaQuery({ minWidth: 2560 });

  return {
    fromXLScreen,
    from2XLScreen,
    from3XLScreen,
    from4XLScreen
  };
};
