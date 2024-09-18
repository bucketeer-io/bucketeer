import { useMemo, useState } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import { ProjectFetcherParams } from '@api/project';
import { useQueryProjects } from '@queries/projects';
import { LIST_PAGE_SIZE } from 'constants/app';
import Filter from 'containers/filter';
import TableContent from 'containers/table-content';
import { projectsHeader } from 'helpers/layouts/header-table';
import { commonTabs } from 'helpers/layouts/tab';
import { TableRows } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { Button } from 'components/button';
import Icon from 'components/icon';
import Spinner from 'components/spinner';
import Tab from 'components/tab';

export const ProjectsContent = () => {
  const formatDateTime = useFormatDateTime();
  const [targetTab, setTargetTab] = useState(commonTabs[0].value);

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

  const rows: TableRows = useMemo(() => {
    const projects = data?.projects || [];
    return projects.map(item => [
      {
        text: item.name,
        type: item.trial ? 'flag' : 'title',
        status: 'new'
      },
      {
        text: '-',
        type: 'text'
      },
      {
        text: '-',
        type: 'text'
      },
      {
        text: formatDateTime(item.createdAt),
        type: 'text'
      },
      {
        type: 'toggle'
      },
      {
        type: 'icon'
      }
    ]);
  }, [data]);

  return (
    <div className="py-8 px-6">
      <Filter
        additionalActions={
          <Button className="flex-1 lg:flex-none">
            <Icon icon={IconAddOutlined} size="sm" />
            {`New Project`}
          </Button>
        }
      />
      <div className="mt-6">
        <Tab
          options={commonTabs}
          value={targetTab}
          onSelect={value => setTargetTab(value)}
        />
        {isLoading ? (
          <div className="pt-20 flex items-center justify-center">
            <Spinner />
          </div>
        ) : (
          <TableContent
            headers={projectsHeader}
            rows={rows}
            emptyTitle="No registered projects"
            emptyDescription="There are no registered projects. Add a new one to start managing."
            emptyActions={
              <div className="flex-center">
                <Button className="w-fit">
                  <Icon icon={IconAddOutlined} size="sm" />
                  {`New Project`}
                </Button>
              </div>
            }
            className="mt-5"
          />
        )}
      </div>
    </div>
  );
};
