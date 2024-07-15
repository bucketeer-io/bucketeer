import tailwindCssForm from "@tailwindcss/forms";
import type { Config } from "tailwindcss";
import { fontFamily } from "tailwindcss/defaultTheme";
import plugin from "tailwindcss/plugin";

const screens = {
  xxs: "320px",
  xs: "444px",
  sm: "600px",
  md: "900px",
  lg: "1200px",
  xl: "1440px",
};

const colors = {
  transparent: "transparent",
  white: "#FFFFFF",
  overlay: "rgba(18, 18, 18, 0.25)",
  primary: {
    600: "#115876",
    500: "#156E93",
    300: "#73A8BF",
    200: "#AACBD8",
    100: "#DAE8EE",
    50: "#EDF4F6",
  },
  secondary: {
    600: "#568A7E",
    500: "#6BAD9E",
    300: "#A7CEC5",
    200: "#CAE1DC",
    100: "#E8F2F0",
    50: "#F3F9F7",
  },
  purple: {
    600: "#915991",
    500: "#B56FB5",
    300: "#D3A9D3",
    200: "#E4CBE4",
    100: "#F3E8F3",
    50: "#F9F4F9",
  },
};

const theme = {
  screens,
  colors,
  boxShadow: {
    DEFAULT: "0px 2px 2px rgba(136, 135, 135, 0.25)",
    none: "none",
  },
  fontFamily: {
    sans: ["Poppins", ...fontFamily.sans],
  },
  extend: {
    animation: {
      fade: "150ms cubic-bezier(0.16, 1, 0.3, 1) 0s 1 normal none running fade",
      zoom: "150ms cubic-bezier(0.16, 1, 0.3, 1) 0s 1 normal none running zoom",
    },
    backdropBlur: {
      xxs: "1px",
      xs: "2px",
    },
  },
};

const container = plugin(({ addComponents }) => {
  const components = {
    ".container": {
      width: "100%",
      maxWidth: "100%",
      paddingLeft: "16px",
      paddingRight: "16px",
      "@screen xs": {
        paddingLeft: "24px",
        paddingRight: "24px",
      },
      "@screen sm": {
        paddingLeft: "24px",
        paddingRight: "24px",
      },
      "@screen md": {
        paddingLeft: "24px",
        paddingRight: "24px",
      },
      "@screen lg": {
        paddingLeft: "0",
        paddingRight: "0",
        marginLeft: "auto",
        marginRight: "auto",
      },
      "@screen xl": {
        maxWidth: "1440px",
      },
    },
    ".container-none": {
      width: "auto",
      maxWidth: "none",
      paddingLeft: "unset",
      paddingRight: "unset",
      marginLeft: "unset",
      marginRight: "unset",
    },
  };

  addComponents(components);
});

const typography = plugin(({ addComponents }) => {
  const components = {
    ".typo-display": {
      fontWeight: "700",
      fontSize: "48px",
      lineHeight: "80px",
      letterSpacing: "0.01em",
    },
    ".typo-header-bold-huge": {
      fontWeight: "700",
      fontSize: "32px",
      lineHeight: "52px",
      letterSpacing: "0.01em",
    },
    ".typo-header-light-huge": {
      fontWeight: "500",
      fontSize: "32px",
      lineHeight: "52px",
      letterSpacing: "0.01em",
    },
    ".typo-header-bold-big": {
      fontWeight: "700",
      fontSize: "24px",
      lineHeight: "40px",
      letterSpacing: "0.01em",
    },
    ".typo-header-light-big": {
      fontWeight: "500",
      fontSize: "24px",
      lineHeight: "40px",
      letterSpacing: "0.01em",
    },
    ".typo-header-bold-medium": {
      fontWeight: "700",
      fontSize: "20px",
      lineHeight: "32px",
      letterSpacing: "0.01em",
    },
    ".typo-header-light-medium": {
      fontWeight: "500",
      fontSize: "20px",
      lineHeight: "32px",
      letterSpacing: "0.01em",
    },
    ".typo-header-bold-small": {
      fontWeight: "700",
      fontSize: "16px",
      lineHeight: "24px",
      letterSpacing: "0.01em",
    },
    ".typo-header-light-small": {
      fontWeight: "500",
      fontSize: "16px",
      lineHeight: "24px",
      letterSpacing: "0.01em",
    },
    ".typo-header-bold-tiny": {
      fontWeight: "700",
      fontSize: "14px",
      lineHeight: "24px",
      letterSpacing: "0.01em",
    },
    ".typo-header-light-tiny": {
      fontWeight: "500",
      fontSize: "14px",
      lineHeight: "24px",
      letterSpacing: "0.01em",
    },
    ".typo-header-bold-tiniest": {
      fontWeight: "700",
      fontSize: "12px",
      lineHeight: "20px",
      letterSpacing: "0.01em",
    },
    ".typo-header-light-tiniest": {
      fontWeight: "500",
      fontSize: "12px",
      lineHeight: "20px",
      letterSpacing: "0.01em",
    },
    ".typo-label-giant": {
      fontWeight: "600",
      fontSize: "24px",
      lineHeight: "36px",
      letterSpacing: "0.01em",
    },
    ".typo-label-huge": {
      fontWeight: "600",
      fontSize: "20px",
      lineHeight: "32px",
      letterSpacing: "0.01em",
    },
    ".typo-label-big": {
      fontWeight: "600",
      fontSize: "16px",
      lineHeight: "28px",
      letterSpacing: "0.01em",
    },
    ".typo-label-medium": {
      fontWeight: "600",
      fontSize: "14px",
      lineHeight: "24px",
      letterSpacing: "0.01em",
    },
    ".typo-label-small": {
      fontWeight: "600",
      fontSize: "12px",
      lineHeight: "20px",
      letterSpacing: "0.01em",
    },
    ".typo-label-tiny": {
      fontWeight: "600",
      fontSize: "10px",
      lineHeight: "16px",
      letterSpacing: "0.01em",
    },
    ".typo-body-huge": {
      fontWeight: "400",
      fontSize: "18px",
      lineHeight: "32px",
      letterSpacing: "0.01em",
    },
    ".typo-body-big": {
      fontWeight: "400",
      fontSize: "16px",
      lineHeight: "28px",
      letterSpacing: "0.01em",
    },
    ".typo-body-medium": {
      fontWeight: "400",
      fontSize: "14px",
      lineHeight: "24px",
      letterSpacing: "0.01em",
    },
    ".typo-body-small": {
      fontWeight: "400",
      fontSize: "12px",
      lineHeight: "20px",
      letterSpacing: "0.01em",
    },
    ".typo-body-tiny": {
      fontWeight: "400",
      fontSize: "10px",
      lineHeight: "16px",
      letterSpacing: "0.01em",
    },
  };

  addComponents(components);
});

const iconSize = plugin(({ addComponents }) => {
  const components = {
    ".icon-size-xxs": {
      width: "16px",
      height: "16px",
      fontSize: "16px",
    },
    ".icon-size-xs": {
      width: "18px",
      height: "18px",
      fontSize: "18px",
    },
    ".icon-size-sm": {
      width: "20px",
      height: "20px",
      fontSize: "20px",
    },
    ".icon-size-md": {
      width: "24px",
      height: "24px",
      fontSize: "24px",
    },
    ".icon-size-lg": {
      width: "28px",
      height: "28px",
      fontSize: "28px",
    },
    ".icon-size-xl": {
      width: "32px",
      height: "32px",
      fontSize: "32px",
    },
    ".icon-size-2xl": {
      width: "40px",
      height: "40px",
      fontSize: "40px",
    },
    ".icon-size-3xl": {
      width: "60px",
      height: "60px",
      fontSize: "60px",
    },
  };

  addComponents(components);
});

export default {
  content: ["./src/**/*.{js,ts,jsx,tsx}"],
  theme,
  plugins: [tailwindCssForm, container, typography, iconSize],
  corePlugins: {
    container: false,
  },
  future: {
    hoverOnlyWhenSupported: true,
  },
} satisfies Config;
