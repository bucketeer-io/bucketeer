import { ExperimentDetailsTab } from '../page-content';
import Metrics from './metrics';

const CollectionLoader = ({
  currentTab
}: {
  currentTab: ExperimentDetailsTab;
}) => {
  if (currentTab === 'results') return <Metrics />;
  if (currentTab === 'settings') return <div>settings</div>;
};

export default CollectionLoader;
