import dayjs from 'dayjs';
import isSameOrBefore from 'dayjs/plugin/isSameOrBefore';
import { createRoot } from 'react-dom/client';
import { RawIntlProvider } from 'react-intl';
import { Provider } from 'react-redux';
import { Router } from 'react-router-dom';

import { history } from './history';
import './styles/styles.css';
import { intl } from './lang';
import { getSelectedLanguage } from './lang/getSelectedLanguage';
import { App } from './pages/index';
import { store } from './store';

dayjs.extend(isSameOrBefore);

document.documentElement.setAttribute('lang', getSelectedLanguage());

async function run() {
  console.log('✅ run');
  const container = document.getElementById('app');
  const root = createRoot(container);
  root.render(
    <Provider store={store}>
      <RawIntlProvider value={intl}>
        <Router history={history}>
          <App />
        </Router>
      </RawIntlProvider>
    </Provider>
  );
}

run();
