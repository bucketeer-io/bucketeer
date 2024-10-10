import { useState } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import Filter from 'containers/filter';
import { commonTabs } from 'helpers/tab';
import { usePartialState } from 'hooks';
import { useTranslation } from 'i18n';
import Button from 'components/button';
import Icon from 'components/icon';
import Tab from 'components/tab';
import PageLayout from 'elements/page-layout';
import CollectionLoader from './collection-loader';

const PageContent = ({ onAdd }: { onAdd: () => void }) => {
  const { t } = useTranslation(['common']);

  const [targetTab, setTargetTab] = useState(commonTabs[0].value);
  const [filters, setFilters] = usePartialState({
    searchQuery: ''
  });

  return (
    <PageLayout.Content>
      <Filter
        action={
          <Button className="flex-1 lg:flex-none" onClick={onAdd}>
            <Icon icon={IconAddOutlined} size="sm" />
            {t(`new-org`)}
          </Button>
        }
        searchValue={filters.searchQuery}
        onSearchChange={v => setFilters({ ...filters, searchQuery: v })}
      />
      <div className="mt-6">
        <Tab
          options={commonTabs}
          value={targetTab}
          onSelect={value => setTargetTab(value)}
        />
        <CollectionLoader />
      </div>
    </PageLayout.Content>
  );
};

export default PageContent;
