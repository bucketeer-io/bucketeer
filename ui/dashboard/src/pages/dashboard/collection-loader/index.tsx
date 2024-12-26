import { useTranslation } from 'i18n';
import Header from 'elements/charts/header';
import CollectionEmpty from 'elements/collection/collection-empty';
import { DataTable } from 'elements/data-table';
import PageLayout from 'elements/page-layout';
import {
  TempMemberType,
  useColumns
} from '../collection-layout/data-collection';
import { EmptyCollection } from '../collection-layout/empty-collection';

const CollectionLoader = () => {
  const { t } = useTranslation(['common']);
  const columns = useColumns();

  const isError = false;
  const emptyState = <CollectionEmpty data={[]} empty={<EmptyCollection />} />;
  const data: TempMemberType[] = [
    {
      account: {
        email: 'johntucker@email.com',
        name: 'John Tucker',
        firstName: 'John',
        lastName: 'Tucker'
      },
      total_flags_created: 5
    },
    {
      account: {
        email: 'ruthguerra@email.com',
        name: 'Ruth Guerra',
        firstName: 'Ruth',
        lastName: 'Guerra'
      },
      total_flags_created: 5
    },
    {
      account: {
        email: 'ruthguerra@email.com',
        name: 'willyamdeppout@email.com',
        firstName: 'Willyam',
        lastName: 'Deppout'
      },
      total_flags_created: 5
    }
  ];

  return isError ? (
    <PageLayout.ErrorState />
  ) : (
    <div className="flex flex-col w-full min-w-[700px] h-fit bg-white border border-gray-200 rounded-2xl">
      <Header title={t('top-members')} />
      <DataTable
        data={data}
        columns={columns}
        emptyCollection={emptyState}
        rowClassName="!shadow-none [&>td]:!border-b [&>td]:!rounded-none [&>td]:last:!border-b-0"
      />
    </div>
  );
};

export default CollectionLoader;
