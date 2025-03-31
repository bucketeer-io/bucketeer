import { useEffect, useState } from 'react';
import { useTranslation } from 'i18n';
import { Feature } from '@types';
import { isEmptyObject } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import { Tabs, TabsContent, TabsList, TabsTrigger } from 'components/tabs';
import Filter from 'elements/filter';
import CollectionLoader from './elements/collection-loader';
import OperationActions from './elements/operation-actions';
import { OperationTab } from './types';

const Operations = ({ feature }: { feature: Feature }) => {
  const { t } = useTranslation(['common', 'table']);
  const { searchOptions, onChangSearchParams } = useSearchParams();

  const [currentTab, setCurrentTab] = useState(OperationTab.ACTIVE);

  useEffect(() => {
    if (isEmptyObject(searchOptions)) {
      const tab = searchOptions?.tab;
      setCurrentTab((tab as OperationTab) || OperationTab.ACTIVE);
    }
  }, [searchOptions]);

  return (
    <div className="flex flex-col w-full gap-y-6">
      <Filter
        searchValue=""
        isShowDocumentation={false}
        onSearchChange={() => {}}
        onOpenFilter={() => {}}
        action={<OperationActions />}
      />

      <Tabs
        className="flex-1 flex h-full flex-col"
        value={currentTab}
        onValueChange={value => {
          const tab = value as OperationTab;
          setCurrentTab(tab);
          onChangSearchParams({ tab });
        }}
      >
        <TabsList>
          <TabsTrigger value="ACTIVE">{t(`active`)}</TabsTrigger>
          <TabsTrigger value="COMPLETED">{t(`completed`)}</TabsTrigger>
        </TabsList>

        <TabsContent value={currentTab}>
          <CollectionLoader feature={feature} currentTab={currentTab} />
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default Operations;
