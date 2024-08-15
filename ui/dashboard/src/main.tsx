import { Suspense } from 'react';
import ReactDOM from 'react-dom/client';
import 'unfonts.css';
import App from 'app';
import './index.css';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

root.render(
  <Suspense>
    <App />
  </Suspense>
);
