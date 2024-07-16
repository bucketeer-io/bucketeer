import tailwindCssForm from '@tailwindcss/forms';
import type { Config } from 'tailwindcss';
import { fontFamily } from 'tailwindcss/defaultTheme';
import plugin from 'tailwindcss/plugin';

const screens = {
  xxs: '320px',
  xs: '444px',
  sm: '600px',
  md: '900px',
  lg: '1200px',
  xl: '1440px'
};

const colors = {
  transparent: 'transparent',
  white: '#FFFFFF',
  overlay: 'rgba(18, 18, 18, 0.25)',
  primary: {
    600: '#492F7A',
    500: '#573792',
    300: '#5F4295',
    200: '#E4CBE4',
    100: '#F3E8F3',
    50: '#F9F4F9'
  },
  secondary: {
    600: '#568A7E',
    500: '#6BAD9E',
    300: '#A7CEC5',
    200: '#CAE1DC',
    100: '#E8F2F0',
    50: '#F3F9F7'
  },
  red: {
    600: '#BF3E3A',
    500: '#EB1726',
    300: '#F59592',
    200: '#F9BFBD',
    100: '#FCE3E2',
    50: '#FEF1F1'
  },
  light: {
    600: '#A1A7B4',
    500: '#B8BDC6',
    300: '#D5D7DD',
    200: '#EDEEF0',
    100: '#F6F6F8',
    50: '#FFFFFF'
  },
  dark: {
    600: '#121212',
    500: '#292A2D',
    300: '#42454D',
    200: '#525660',
    100: '#6E7481',
    50: '#8A91A1'
  }
};

const theme = {
  screens,
  colors,
  boxShadow: {
    DEFAULT: '0px 2px 2px rgba(136, 135, 135, 0.25)',
    'border-primary-500': `inset 0 0 0 1px ${colors.primary[500]}`,
    'border-primary-600': `inset 0 0 0 1px ${colors.primary[600]}`,
    'border-red-500': `inset 0 0 0 1px ${colors.red[500]}`,
    'border-light-500': `inset 0 0 0 1px ${colors.light[500]}`,
    'border-dark-200': `inset 0 0 0 1px ${colors.dark[200]}`,
    none: 'none'
  },
  fontFamily: {
    sans: ['Poppins', ...fontFamily.sans]
  },
  extend: {
    animation: {
      fade: '150ms cubic-bezier(0.16, 1, 0.3, 1) 0s 1 normal none running fade',
      zoom: '150ms cubic-bezier(0.16, 1, 0.3, 1) 0s 1 normal none running zoom'
    },
    keyframes: {
      fade: {
        '0%': {
          opacity: '0'
        },
        '100%': {
          opacity: '1'
        }
      }
    },
    backdropBlur: {
      xxs: '1px',
      xs: '2px'
    }
  }
};

const container = plugin(({ addComponents }) => {
  const components = {
    '.container': {
      width: '100%',
      maxWidth: '100%',
      paddingLeft: '16px',
      paddingRight: '16px',
      '@screen xs': {
        paddingLeft: '24px',
        paddingRight: '24px'
      },
      '@screen sm': {
        paddingLeft: '24px',
        paddingRight: '24px'
      },
      '@screen md': {
        paddingLeft: '24px',
        paddingRight: '24px'
      },
      '@screen lg': {
        paddingLeft: '0',
        paddingRight: '0',
        marginLeft: 'auto',
        marginRight: 'auto'
      },
      '@screen xl': {
        maxWidth: '1440px'
      }
    },
    '.container-none': {
      width: 'auto',
      maxWidth: 'none',
      paddingLeft: 'unset',
      paddingRight: 'unset',
      marginLeft: 'unset',
      marginRight: 'unset'
    }
  };

  addComponents(components);
});

const typography = plugin(({ addComponents }) => {
  const components = {
    '.typo-display': {
      fontWeight: '700',
      fontSize: '48px',
      lineHeight: '80px',
      letterSpacing: '0.01em'
    },
    '.typo-body-huge': {
      fontWeight: '400',
      fontSize: '20px',
      lineHeight: '32px',
      letterSpacing: '0.01em'
    },
    '.typo-body-big': {
      fontWeight: '400',
      fontSize: '18px',
      lineHeight: '28px',
      letterSpacing: '0.01em'
    },
    '.typo-body-medium': {
      fontWeight: '400',
      fontSize: '16px',
      lineHeight: '24px',
      letterSpacing: '0.01em'
    },
    '.typo-body-small': {
      fontWeight: '400',
      fontSize: '14px',
      lineHeight: '20px',
      letterSpacing: '0.01em'
    },
    '.typo-body-tiny': {
      fontWeight: '400',
      fontSize: '12px',
      lineHeight: '16px',
      letterSpacing: '0.01em'
    }
  };

  addComponents(components);
});

const iconSize = plugin(({ addComponents }) => {
  const components = {
    '.icon-size-xxs': {
      width: '16px',
      height: '16px',
      fontSize: '16px'
    },
    '.icon-size-xs': {
      width: '18px',
      height: '18px',
      fontSize: '18px'
    },
    '.icon-size-sm': {
      width: '20px',
      height: '20px',
      fontSize: '20px'
    },
    '.icon-size-md': {
      width: '24px',
      height: '24px',
      fontSize: '24px'
    },
    '.icon-size-lg': {
      width: '28px',
      height: '28px',
      fontSize: '28px'
    },
    '.icon-size-xl': {
      width: '32px',
      height: '32px',
      fontSize: '32px'
    },
    '.icon-size-2xl': {
      width: '40px',
      height: '40px',
      fontSize: '40px'
    },
    '.icon-size-3xl': {
      width: '60px',
      height: '60px',
      fontSize: '60px'
    }
  };

  addComponents(components);
});

export default {
  content: ['./src/**/*.{js,ts,jsx,tsx}'],
  theme,
  plugins: [tailwindCssForm, container, typography, iconSize],
  corePlugins: {
    container: false
  },
  future: {
    hoverOnlyWhenSupported: true
  }
} satisfies Config;
