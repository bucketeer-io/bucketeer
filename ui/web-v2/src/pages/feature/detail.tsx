import { SerializedError } from '@reduxjs/toolkit';
import { FC, memo, useEffect } from 'react';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import {
  NavLink,
  Route,
  Switch,
  Redirect,
  useRouteMatch,
  useParams
} from 'react-router-dom';

import { FeatureHeader } from '../../components/FeatureHeader';
import {
  PAGE_PATH_CODE_REFS,
  PAGE_PATH_FEATURES,
  PAGE_PATH_FEATURE_AUTOOPS,
  PAGE_PATH_FEATURE_EVALUATION,
  PAGE_PATH_FEATURE_EXPERIMENTS,
  PAGE_PATH_FEATURE_HISTORY,
  PAGE_PATH_FEATURE_SETTING,
  PAGE_PATH_FEATURE_TARGETING,
  PAGE_PATH_FEATURE_TRIGGER,
  PAGE_PATH_FEATURE_VARIATION,
  PAGE_PATH_NEW,
  PAGE_PATH_ROOT
} from '../../constants/routing';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import {
  selectById as selectFeatureById,
  getFeature
} from '../../modules/features';
import { useCurrentEnvironment } from '../../modules/me';
import { Feature } from '../../proto/feature/feature_pb';
import { AppDispatch } from '../../store';

import { FeatureAutoOpsPage } from './autoops';
import { FeatureEvaluationPage } from './evaluation';
import { FeatureExperimentsPage } from './experiments';
import { FeatureHistoryPage } from './history';
import { FeatureCodeRefsPage } from './coderefs';
import { FeatureSettingsPage } from './settings';
import { FeatureTargetingPage } from './targeting';
import { FeatureTriggerPage } from './triggers';
import { FeatureVariationsPage } from './variations';
import {
  listCodeRefs,
  selectAll as selectAllCodeRefs
} from '../../modules/codeRefs';
import { CodeReference } from '../../proto/coderef/code_reference_pb';

export const FeatureDetailPage: FC = memo(() => {
  const dispatch = useDispatch<AppDispatch>();
  const { url } = useRouteMatch();
  const currentEnvironment = useCurrentEnvironment();
  const { featureId } = useParams<{ featureId: string }>();
  const [feature] = useSelector<
    AppState,
    [Feature.AsObject | undefined, SerializedError | null]
  >(
    (state) => [
      selectFeatureById(state.features, featureId),
      state.features.getFeatureError
    ],
    shallowEqual
  );

  const codeRefs = useSelector<AppState, CodeReference.AsObject[]>(
    (state) => selectAllCodeRefs(state.codeRefs),
    shallowEqual
  );

  useEffect(() => {
    if (featureId) {
      dispatch(
        getFeature({
          environmentId: currentEnvironment.id,
          id: featureId
        })
      );
      dispatch(
        listCodeRefs({
          environmentId: currentEnvironment.id,
          featureId: featureId,
          pageSize: 0,
          fileExtension: null,
          repositoryBranch: null,
          repositoryType: null
        })
      );
    }
  }, [featureId, dispatch, currentEnvironment]);

  if (!feature) {
    return <div>loading</div>;
  }

  return (
    <div className="bg-white h-full">
      <div className="pt-5 px-10">
        <FeatureHeader featureId={featureId} />
        <div className="hidden sm:block">
          <nav className="-mb-px flex space-x-8" aria-label="Tabs">
            {createTabs().map((tab, idx) => {
              return (
                <NavLink
                  key={idx}
                  className="
                      tab-item
                      flex items-center
                      border-transparent
                      text-gray-500
                      hover:text-gray-700
                      whitespace-nowrap py-4 px-1 border-b-2
                      font-medium text-sm"
                  to={`${PAGE_PATH_ROOT}${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}/${featureId}${tab.to}`}
                >
                  {tab.message}
                  {tab.to === PAGE_PATH_CODE_REFS &&
                    (codeRefs.length === 0 ? (
                      <div className="rounded-sm bg-[#F3F9FD] text-[#399CE4] px-2 py-[6px] text-sm inline-block ml-3">
                        New
                      </div>
                    ) : (
                      <span className="ml-1">({codeRefs.length})</span>
                    ))}
                </NavLink>
              );
            })}
          </nav>
        </div>
      </div>
      <div className="border-b border-gray-300"></div>
      <Switch>
        <Route
          exact
          path={`${url}`}
          component={() => (
            <Redirect to={`${url}${PAGE_PATH_FEATURE_TARGETING}`} />
          )}
        />
        <Route exact path={`${url}${PAGE_PATH_FEATURE_TARGETING}`}>
          <FeatureTargetingPage featureId={featureId} />
        </Route>
        <Route exact path={`${url}${PAGE_PATH_FEATURE_SETTING}`}>
          <FeatureSettingsPage featureId={featureId} />
        </Route>
        <Route exact path={`${url}${PAGE_PATH_FEATURE_EVALUATION}`}>
          <FeatureEvaluationPage featureId={featureId} />
        </Route>
        <Route
          exact
          path={[
            `${url}${PAGE_PATH_FEATURE_EXPERIMENTS}`,
            `${url}${PAGE_PATH_FEATURE_EXPERIMENTS}${PAGE_PATH_NEW}`
          ]}
        >
          <FeatureExperimentsPage featureId={featureId} />
        </Route>
        <Route exact path={`${url}${PAGE_PATH_FEATURE_TRIGGER}`}>
          <FeatureTriggerPage featureId={featureId} />
        </Route>
        <Route exact path={`${url}${PAGE_PATH_FEATURE_VARIATION}`}>
          <FeatureVariationsPage featureId={featureId} />
        </Route>
        <Route
          exact
          path={[
            `${url}${PAGE_PATH_FEATURE_AUTOOPS}`,
            `${url}${PAGE_PATH_FEATURE_AUTOOPS}/:operationType`
          ]}
        >
          <FeatureAutoOpsPage featureId={featureId} />
        </Route>
        <Route exact path={`${url}${PAGE_PATH_FEATURE_HISTORY}`}>
          <FeatureHistoryPage featureId={featureId} />
        </Route>
        <Route exact path={`${url}${PAGE_PATH_CODE_REFS}`}>
          <FeatureCodeRefsPage featureId={featureId} />
        </Route>
      </Switch>
    </div>
  );
});

export interface TabItem {
  readonly message: string;
  readonly to: string;
}

const createTabs = (): Array<TabItem> => {
  return [
    {
      message: intl.formatMessage(messages.feature.tab.targeting),
      to: PAGE_PATH_FEATURE_TARGETING
    },
    {
      message: intl.formatMessage(messages.feature.tab.variations),
      to: PAGE_PATH_FEATURE_VARIATION
    },
    {
      message: intl.formatMessage(messages.feature.tab.autoOps),
      to: PAGE_PATH_FEATURE_AUTOOPS
    },
    {
      message: intl.formatMessage(messages.feature.tab.triggers),
      to: PAGE_PATH_FEATURE_TRIGGER
    },
    {
      message: intl.formatMessage(messages.feature.tab.experiments),
      to: PAGE_PATH_FEATURE_EXPERIMENTS
    },
    {
      message: intl.formatMessage(messages.feature.tab.evaluation),
      to: PAGE_PATH_FEATURE_EVALUATION
    },
    {
      message: intl.formatMessage(messages.feature.tab.history),
      to: PAGE_PATH_FEATURE_HISTORY
    },
    {
      message: intl.formatMessage(messages.feature.tab.codeRefs),
      to: PAGE_PATH_CODE_REFS
    },
    {
      message: intl.formatMessage(messages.feature.tab.settings),
      to: PAGE_PATH_FEATURE_SETTING
    }
  ];
};
