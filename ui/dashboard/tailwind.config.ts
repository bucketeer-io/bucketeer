import tailwindCssForm from '@tailwindcss/forms';
import type { Config } from 'tailwindcss';
import plugin from 'tailwindcss/plugin';

const screens = {
  xxs: '320px',
  xs: '444px',
  sm: '600px',
  md: '900px',
  lg: '1200px',
  xl: '1440px',
  xxl: '1600px'
};

const colors = {
  transparent: 'transparent',
  white: '#FFFFFF',
  'primary-additional': '#EFEFFE',
  overlay: 'rgba(0, 0, 0, 0.8)',
  primary: {
    900: '#292C4C',
    800: '#352F5E',
    700: '#40316F',
    600: '#4C3481',
    500: '#573792',
    400: '#6746A4',
    300: '#9A87BE',
    200: '#BCAFD3',
    100: '#E8E4F1',
    50: '#F9F8FB'
  },
  gray: {
    900: '#1E293B',
    800: '#334155',
    700: '#475569',
    600: '#64748B',
    500: '#94A3B8',
    400: '#CBD5E1',
    300: '#E2E8F0',
    200: '#F1F5F9',
    100: '#FAFAFC',
    50: '#FFFFFF'
  },
  accent: {
    pink: {
      900: '#463154',
      800: '#6F396D',
      700: '#974285',
      600: '#C04A9E',
      500: '#E439AC',
      400: '#E961BD',
      300: '#EF88CD',
      200: '#F4B0DE',
      100: '#FAD7EE',
      50: '#FDECF7'
    },
    green: {
      900: '#25473C',
      800: '#2C653E',
      700: '#32833F',
      600: '#39A141',
      500: '#40BF42',
      400: '#66CC68',
      300: '#8CD98E',
      200: '#B3E5B3',
      100: '#D9F2D9',
      50: '#ECF9ED'
    },
    red: {
      900: '#472537',
      800: '#702233',
      700: '#991E2E',
      600: '#C21B2A',
      500: '#EB1726',
      400: '#EF4551',
      300: '#F3747D',
      200: '#F7A2A8',
      100: '#FBD1D4',
      50: '#FEF0F2'
    },
    blue: {
      900: '#23405D',
      800: '#29577F',
      700: '#2E6EA0',
      600: '#3485C2',
      500: '#399CE4',
      400: '#61B0E9',
      300: '#88C4EF',
      200: '#B0D7F4',
      100: '#D7EBFA',
      50: '#ECF6FD'
    },
    orange: {
      900: '#4A403F',
      800: '#765743',
      700: '#A26D46',
      600: '#CE844A',
      500: '#FA9B4E',
      400: '#FBAF71',
      300: '#FCC395',
      200: '#FDD7B8',
      100: '#FEEBDC',
      50: '#FFF6EE'
    },
    yellow: {
      900: '#725201',
      800: '#A17401',
      700: '#C68F02',
      600: '#E4A502',
      500: '#FFB802',
      400: '#FFE072',
      300: '#FFEBA1',
      200: '#FFF3C6',
      100: '#FFFAD6',
      50: '#FDFBE8'
    }
  }
};

const theme = {
  screens,
  colors,
  boxShadow: {
    DEFAULT: '0px 2px 2px rgba(136, 135, 135, 0.25)',
    menu: '0px 8px 12px rgba(0, 0, 0, 0.08)',
    card: '0px 4px 8px 1px rgba(0, 0, 0, 0.1)',
    'card-secondary': '0px 4px 13px 2px rgba(0, 0, 0, 0.1)',
    dropdown: '0px 4px 8px rgba(35, 35, 35, 0.1)',
    'border-primary-500': `inset 0 0 0 1px ${colors.primary[500]}`,
    'border-primary-600': `inset 0 0 0 1px ${colors.primary[600]}`,
    'border-primary-700': `inset 0 0 0 1px ${colors.primary[700]}`,
    'border-gray-200': `inset 0 0 0 1px ${colors.gray[200]}`,
    'border-gray-300': `inset 0 0 0 1px ${colors.gray[300]}`,
    'border-gray-400': `inset 0 0 0 1px ${colors.gray[400]}`,
    'border-gray-500': `inset 0 0 0 1px ${colors.gray[500]}`,
    'border-accent-red-500': `inset 0 0 0 1px ${colors.accent.red[500]}`,
    none: 'none',
    tooltip:
      'rgba(29, 29, 29, 0.35) 0px 10px 38px -10px, rgba(29, 29, 29, 0.20) 0px 10px 20px -15px'
  },
  fontFamily: {
    'sofia-pro': ['Sofia Pro', 'sans-serif'],
    'fira-code': ['FiraCode', 'monospace']
  },
  extend: {
    animation: {
      fade: '150ms cubic-bezier(0.16, 1, 0.3, 1) 0s 1 normal none running fade',
      zoom: '150ms cubic-bezier(0.16, 1, 0.3, 1) 0s 1 normal none running zoom',
      'slide-left':
        '150ms cubic-bezier(0.16, 1, 0.3, 1) 0s 1 normal none running slide-left',
      'slide-up':
        '150ms cubic-bezier(0.16, 1, 0.3, 1) 0s 1 normal none running slide-up'
    },
    keyframes: {
      fade: {
        '0%': {
          opacity: '0'
        },
        '100%': {
          opacity: '1'
        }
      },
      zoom: {
        '0%': {
          opacity: '0',
          transform: 'scale(0.96)'
        },
        '100%': {
          opacity: '1',
          transform: 'scale(1)'
        }
      },
      'slide-left': {
        '0%': {
          opacity: '0',
          right: '-100%'
        },
        '100%': {
          opacity: '1',
          right: '0',
          transform: 'scale(1)'
        }
      },
      'slide-up': {
        '0%': {
          opacity: '0',
          bottom: '-100%'
        },
        '100%': {
          opacity: '1',
          bottom: '0',
          transform: 'scale(1)'
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
    },
    'size-120': {
      width: '120px',
      height: '120px'
    },
    '.flex-center': {
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center'
    }
  };

  addComponents(components);
});

const typography = plugin(({ addComponents }) => {
  const components = {
    '.typo-head-bold-huge': {
      fontWeight: '700',
      fontSize: '24px',
      lineHeight: '32px',
      letterSpacing: '0.01em'
    },
    '.typo-head-semi-huge': {
      fontWeight: '600',
      fontSize: '24px',
      lineHeight: '32px',
      letterSpacing: '0.01em'
    },
    '.typo-head-light-huge': {
      fontWeight: '500',
      fontSize: '24px',
      lineHeight: '32px',
      letterSpacing: '0.01em'
    },
    '.typo-head-bold-big': {
      fontWeight: '700',
      fontSize: '20px',
      lineHeight: '28px',
      letterSpacing: '0.01em'
    },
    '.typo-head-semi-big': {
      fontWeight: '600',
      fontSize: '20px',
      lineHeight: '28px',
      letterSpacing: '0.01em'
    },
    '.typo-head-light-big': {
      fontWeight: '500',
      fontSize: '20px',
      lineHeight: '24px',
      letterSpacing: '0.01em'
    },
    '.typo-head-bold-medium': {
      fontWeight: '700',
      fontSize: '18px',
      lineHeight: '24px',
      letterSpacing: '0.01em'
    },
    '.typo-head-semi-medium': {
      fontWeight: '600',
      fontSize: '18px',
      lineHeight: '24px',
      letterSpacing: '0.01em'
    },
    '.typo-head-light-medium': {
      fontWeight: '500',
      fontSize: '18px',
      lineHeight: '24px',
      letterSpacing: '0.01em'
    },
    '.typo-head-bold-small': {
      fontWeight: '700',
      fontSize: '16px',
      lineHeight: '20px',
      letterSpacing: '0.01em'
    },
    '.typo-head-semi-small': {
      fontWeight: '600',
      fontSize: '16px',
      lineHeight: '20px',
      letterSpacing: '0.01em'
    },
    '.typo-head-light-small': {
      fontWeight: '500',
      fontSize: '16px',
      lineHeight: '20px',
      letterSpacing: '0.01em'
    },
    '.typo-head-bold-tiny': {
      fontWeight: '800',
      fontSize: '10px',
      lineHeight: '10px',
      letterSpacing: '0.01em'
    },
    '.typo-para-big': {
      fontWeight: '400',
      fontSize: '18px',
      lineHeight: '28px',
      letterSpacing: '0.01em'
    },
    '.typo-para-medium': {
      fontWeight: '400',
      fontSize: '16px',
      lineHeight: '24px',
      letterSpacing: '0.01em'
    },
    '.typo-para-small': {
      fontWeight: '400',
      fontSize: '14px',
      lineHeight: '20px',
      letterSpacing: '0.01em'
    },
    '.typo-para-tiny': {
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
    '.icon-size-fit': {
      width: 'fit-content',
      height: 'fit-content'
    },
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
      width: '36px',
      height: '36px',
      fontSize: '36px'
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
  mode: 'jit',
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
