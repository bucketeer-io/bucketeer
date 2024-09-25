import { useEffect, useState } from 'react';
import { IconMoreVertOutlined } from 'react-icons-material-design';
import { useParams } from 'react-router-dom';
import { PAGE_PATH_ORGANIZATIONS } from 'constants/routing';
import PageDetailHeader from 'containers/page-details-header';
import { OrganizationDetailsContent } from 'containers/pages/organization-details';
import { projectTab, settingTab, userTab } from 'helpers/tab';
import { useAddQuery, useQuery } from 'hooks';
import { Popover } from 'components/popover';

const OrganizationDetails = () => {
  const query = useQuery();
  const params = useParams();
  const { organizationId } = params;
  const { addQuery } = useAddQuery();
  const tab = query.get('tab');

  const tabs = [projectTab, userTab, settingTab];

  const [targetTab, setTargetTab] = useState(tab || tabs[0].value);

  const handleChangeTab = (value: string) => {
    setTargetTab(value);
    addQuery(query, { tab: value });
  };

  useEffect(() => {
    if (!tab) handleChangeTab(projectTab.value);
  }, [tab]);

  return (
    <div className="flex flex-col size-full overflow-auto">
      <PageDetailHeader
        title="Organization Name 1"
        description="Created 21 hours ago"
        navigateRoute={PAGE_PATH_ORGANIZATIONS}
        tabs={tabs}
        targetTab={targetTab}
        titleActions={<Popover options={[]} icon={IconMoreVertOutlined} />}
        onSelectTab={handleChangeTab}
      />
      <OrganizationDetailsContent
        targetTab={targetTab}
        organizationId={organizationId}
      />
    </div>
  );
};

export default OrganizationDetails;
