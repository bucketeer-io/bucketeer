import { Suspense } from 'react';
import ReactDOM from 'react-dom/client';
import { Toaster } from 'react-hot-toast';
import { ThemeProvider } from 'hooks/use-theme';
import 'unfonts.css';
import App from 'app';
import './index.css';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

root.render(
  <Suspense>
    <ThemeProvider>
      <App />
      <Toaster />
    </ThemeProvider>
  </Suspense>
);
