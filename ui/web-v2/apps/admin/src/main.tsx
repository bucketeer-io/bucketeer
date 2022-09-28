import { render } from 'react-dom';
import { RawIntlProvider } from 'react-intl';
import { Provider } from 'react-redux';
import { Router } from 'react-router-dom';

import { history } from './history';
import './styles/styles.css';
import { intl } from './lang';
import { App } from './pages/index';
import { store } from './store';

async function run() {
  render(
    <Provider store={store}>
      <RawIntlProvider value={intl}>
        <Router history={history}>
          <App />
        </Router>
      </RawIntlProvider>
    </Provider>,
    document.getElementById('app')
  );
}

run();
