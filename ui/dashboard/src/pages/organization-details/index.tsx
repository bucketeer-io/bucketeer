import { useMemo, useState } from 'react';
import {
  IconArchiveOutlined,
  IconEditOutlined
} from 'react-icons-material-design';
import { ProjectFetcherParams } from '@api/project';
import { useQueryProjects } from '@queries/projects';
import { LIST_PAGE_SIZE } from 'constants/app';
import { PAGE_PATH_ORGANIZATIONS } from 'constants/routing';
import PageDetailHeader from 'containers/page-details-header';
import TableContent from 'containers/table-content';
import { projectsHeader } from 'helpers/layouts/header-table';
import {
  environmentTab,
  projectTab,
  settingTab,
  userTab
} from 'helpers/layouts/tab';
import { TableRows } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import Spinner from 'components/spinner';

const OrganizationDetails = () => {
  const formatDateTime = useFormatDateTime();
  const tabs = [projectTab, environmentTab, userTab, settingTab];

  const [targetTab, setTargetTab] = useState(tabs[0].value);

  const projectParams: ProjectFetcherParams = {
    pageSize: LIST_PAGE_SIZE,
    cursor: String(0),
    orderBy: 'DEFAULT',
    orderDirection: 'ASC',
    searchKeyword: '',
    disabled: false,
    organizationIds: []
  };

  const { data, isLoading } = useQueryProjects({
    params: projectParams
  });

  const rows = useMemo(() => {
    const projects = data?.projects || [];
    return projects.map(item => [
      {
        text: item.name,
        type: item.trial ? 'flag' : 'title',
        status: 'new'
      },
      {
        text: item.creatorEmail,
        type: 'text'
      },
      {
        text: item.environmentCount,
        type: 'text'
      },
      {
        text: item.featureFlagCount,
        type: 'text'
      },
      {
        text: formatDateTime(item.createdAt),
        type: 'text'
      },
      {
        type: 'icon',
        options: [
          {
            label: 'Edit Project',
            icon: IconEditOutlined
          },
          {
            label: 'Archive Project',
            icon: IconArchiveOutlined
          }
        ]
      }
    ]);
  }, [data]);

  return (
    <div className="flex flex-col size-full overflow-auto">
      <PageDetailHeader
        title="Organization Name 1"
        description="Created 21 hours ago"
        navigateRoute={PAGE_PATH_ORGANIZATIONS}
        tabs={tabs}
        targetTab={targetTab}
        onSelectTab={setTargetTab}
      />
      {isLoading ? (
        <div className="pt-20 flex items-center justify-center">
          <Spinner />
        </div>
      ) : (
        <TableContent
          headers={projectsHeader}
          rows={rows as TableRows}
          emptyTitle="No registered projects"
          emptyDescription="There are no registered projects. Add a new one to start managing."
          className="px-6 pb-8"
        />
      )}
    </div>
  );
};

export default OrganizationDetails;
