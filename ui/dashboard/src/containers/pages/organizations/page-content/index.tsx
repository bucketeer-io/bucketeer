import { useMemo, useState } from 'react';
import {
  IconAddOutlined,
  IconArchiveOutlined,
  IconEditOutlined
} from 'react-icons-material-design';
import { useNavigate } from 'react-router-dom';
import { OrganizationsFetcherParams } from '@api/organization';
import { useQueryOrganizations } from '@queries/organizations';
import { LIST_PAGE_SIZE } from 'constants/app';
import Filter from 'containers/filter';
import TableContent from 'containers/table-content';
import { orgHeader } from 'helpers/layouts/header-table';
import { commonTabs } from 'helpers/layouts/tab';
import { TableRows } from '@types';
import { useFormatDateTime } from 'utils/date-time';
import { Button } from 'components/button';
import Icon from 'components/icon';
import Spinner from 'components/spinner';
import Tab from 'components/tab';

export const OrganizationsContent = () => {
  const navigate = useNavigate();
  const formatDateTime = useFormatDateTime();

  const [targetTab, setTargetTab] = useState(commonTabs[0].value);

  const organizationParams: OrganizationsFetcherParams = {
    pageSize: LIST_PAGE_SIZE,
    cursor: String(0),
    orderBy: 'DEFAULT',
    orderDirection: 'ASC',
    disabled: false,
    archived: false
  };

  const { data, isLoading } = useQueryOrganizations({
    params: organizationParams
  });

  const rows = useMemo(() => {
    const organizations = data?.Organizations || [];

    return organizations.map(organization => [
      {
        text: organization.name,
        type: 'title',
        width: '40%',
        onClick: () => navigate('/organizations/1')
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
        text: '-',
        type: 'text'
      },
      {
        text: formatDateTime(organization.createdAt),
        type: 'text',
        tooltip: 'Thu, September 21, 2023'
      },
      {
        type: 'icon',
        options: [
          {
            label: 'Edit Organization',
            icon: IconEditOutlined
          },
          {
            label: 'Archive Organization',
            icon: IconArchiveOutlined
          }
        ]
      }
    ]);
  }, [data]);

  return (
    <div className="py-8 px-6">
      <Filter
        additionalActions={
          <Button className="flex-1 lg:flex-none">
            <Icon icon={IconAddOutlined} size="sm" />
            {`New Organization`}
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
            headers={orgHeader}
            rows={rows as TableRows}
            emptyTitle="No registered organizations"
            emptyDescription="There are no registered organizations. Add a new one to start managing."
            emptyActions={
              <div className="flex justify-center">
                <Button className="w-fit">
                  <Icon icon={IconAddOutlined} size="sm" />
                  {`New Organization`}
                </Button>
              </div>
            }
          />
        )}
      </div>
    </div>
  );
};
