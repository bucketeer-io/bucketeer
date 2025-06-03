import react from '@vitejs/plugin-react';
import Unfonts from 'unplugin-fonts/vite';
import { defineConfig } from 'vite';
import svgr from 'vite-plugin-svgr';
import viteTsconfigPaths from 'vite-tsconfig-paths';
import { TanStackRouterVite } from '@tanstack/router-plugin/vite'

// https://vitejs.dev/config/
export default defineConfig({
  base: '/v3/',
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
    TanStackRouterVite({ target: 'react', autoCodeSplitting: true }),
    react(),
    svgr(),
    viteTsconfigPaths(),
    Unfonts({
      custom: {
        families: [
          {
            name: 'Sofia Pro',
            src: './src/assets/fonts/sofiapro/*.ttf'
          },
          {
            name: 'FiraCode',
            src: './src/assets/fonts/firacode/*.ttf'
          }
        ]
      }
    })
  ]
});
