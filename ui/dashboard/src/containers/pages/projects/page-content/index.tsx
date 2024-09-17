import { useMemo, useState } from 'react';
import { IconAddOutlined } from 'react-icons-material-design';
import Filter from 'containers/filter';
import TableContent from 'containers/table-content';
import { projectsHeader } from 'helpers/layouts/header-table';
import { commonTabs } from 'helpers/layouts/tab';
import { projectsMockData } from 'helpers/mock/table-data';
import { Button } from 'components/button';
import Icon from 'components/icon';
import Tab from 'components/tab';
import { TableRows } from 'components/table';
import { StatusTagType } from 'components/tag/status-tag';

export const ProjectsContent = () => {
  const [targetTab, setTargetTab] = useState(commonTabs[0].value);

  const rows: TableRows = useMemo(() => {
    return projectsMockData.map(e => [
      {
        text: e.projectName,
        type: e?.status ? 'flag' : 'title',
        status: e?.status as StatusTagType
      },
      {
        text: e.userName,
        description: e.email,
        type: 'text'
      },
      {
        text: e.envQuantity,
        type: 'text'
      },
      {
        text: e.createdAt,
        type: 'text'
      },
      {
        type: 'toggle'
      },
      {
        type: 'icon'
      }
    ]);
  }, [projectsMockData]);

  return (
    <div className="py-8 px-6">
      <Filter
        additionalActions={
          <Button className="flex-1 lg:flex-none">
            <Icon icon={IconAddOutlined} size="sm" />
            New Project
          </Button>
        }
      />
      <div className="mt-6">
        <Tab
          options={commonTabs}
          value={targetTab}
          onSelect={value => setTargetTab(value)}
        />
        <TableContent
          headers={projectsHeader}
          rows={rows}
          emptyTitle="No registered projects"
          emptyDescription="There are no registered projects. Add a new one to start managing."
          emptyActions={
            <div className="flex-center">
              <Button className="w-[164px]">
                <Icon icon={IconAddOutlined} size="sm" />
                New Project
              </Button>
            </div>
          }
          className="mt-5"
        />
      </div>
    </div>
  );
};
