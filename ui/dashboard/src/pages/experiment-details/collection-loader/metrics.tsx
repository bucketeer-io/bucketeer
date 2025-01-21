import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { isNotEmptyObject } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import Evaluation from './evaluation';

type MetricsTab = 'EVALUATION' | 'CONVERSION';

type MetricFilters = {
  tab: MetricsTab;
  [key: string]: MetricsTab | string;
};

const Metrics = () => {
  const { t } = useTranslation(['common']);
  const { searchOptions, onChangSearchParams } = useSearchParams();
  const defaultFilters: MetricFilters = {
    tab: 'EVALUATION',
    ...searchOptions
  };
  const [filters, setFilters] = useState(defaultFilters);

  useEffect(() => {
    if (isNotEmptyObject(searchOptions)) {
      setFilters(defaultFilters);
    }
  }, [searchOptions]);

  return (
    <div className="w-fit min-w-full p-5 shadow-card rounded-lg bg-white">
      <p className="text-gray-800 typo-head-bold-small">{t('metrics')}</p>
      <Tabs
        className="flex-1 flex h-full flex-col mt-6"
        value={filters.tab}
        onValueChange={tab =>
          onChangSearchParams({
            tab
          })
        }
      >
        <TabsList>
          <TabsTrigger value="EVALUATION">{t(`evaluation`)}</TabsTrigger>
          <TabsTrigger value="CONVERSION">{t(`conversion-rate`)}</TabsTrigger>
        </TabsList>

        <TabsContent value={filters.tab}>
          {filters.tab === 'EVALUATION' && <Evaluation />}
          {filters.tab === 'CONVERSION' && <div>CONVERSION</div>}
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default Metrics;
