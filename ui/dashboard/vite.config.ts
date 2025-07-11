import react from '@vitejs/plugin-react';
import Unfonts from 'unplugin-fonts/vite';
import { defineConfig } from 'vite';
import svgr from 'vite-plugin-svgr';
import viteTsconfigPaths from 'vite-tsconfig-paths';

// https://vitejs.dev/config/
export default defineConfig({
  base: '/',
  preview: {
    port: 8000,
    open: true
  },
  server: {
    port: 8000,
    open: true
  },
  build: {
    outDir: 'build'
  },
  plugins: [
    react(),
    svgr(),
    viteTsconfigPaths(),
    Unfonts({
      custom: {
        families: [
          {
            name: 'Sofia Pro',
            src: './src/assets/fonts/sofiapro/*.woff2'
          },
          {
            name: 'FiraCode',
            src: './src/assets/fonts/firacode/*.woff2'
          },
          {
            name: 'Noto Sans JP',
            src: './src/assets/fonts/noto-sans-jp/*.woff2'
          }
        ],
        preload: true,
        display: 'swap'
      }
    })
  ]
});
