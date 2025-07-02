import { useState } from 'react';
import { useParams } from 'react-router-dom';
import { usePartialState, useToggleOpen } from 'hooks';
import pickBy from 'lodash/pickBy';
import { Account } from '@types';
import { isNotEmpty } from 'utils/data-type';
import { useSearchParams } from 'utils/search-params';
import MemberDetailsModal from 'pages/members/member-modal/member-details-modal';
import Filter from 'elements/filter';
import TableListContainer from 'elements/table-list-container';
import { OrganizationMembersFilters } from '../types';
import CollectionLoader from './collection-loader';

const OrganizationMembers = () => {
  const { organizationId } = useParams();

  const { searchOptions, onChangSearchParams } = useSearchParams();
  const searchFilters: Partial<OrganizationMembersFilters> = searchOptions;

  const [selectedMember, setSelectedMember] = useState<Account>();

  const [isOpenDetailsModal, onOpenDetailsModal, onCloseDetailsModal] =
    useToggleOpen(false);

  const defaultFilters = {
    page: 1,
    orderBy: 'CREATED_AT',
    orderDirection: 'DESC',
    ...searchFilters
  } as OrganizationMembersFilters;

  const [filters, setFilters] =
    usePartialState<OrganizationMembersFilters>(defaultFilters);

  const onChangeFilters = (values: Partial<OrganizationMembersFilters>) => {
    const options = pickBy({ ...filters, ...values }, v => isNotEmpty(v));
    onChangSearchParams(options);
    setFilters({ ...values });
  };

  const filterParams = { ...filters, organizationId };

  return (
    <>
      <Filter
        isShowDocumentation={false}
        searchValue={filters.searchQuery}
        onSearchChange={searchQuery => onChangeFilters({ searchQuery })}
      />
      <TableListContainer>
        <CollectionLoader
          filters={filterParams}
          setFilters={onChangeFilters}
          onActions={member => {
            setSelectedMember(member);
            onOpenDetailsModal();
          }}
        />
      </TableListContainer>

      {isOpenDetailsModal && (
        <MemberDetailsModal
          isOpen={isOpenDetailsModal}
          onClose={onCloseDetailsModal}
          member={selectedMember!}
        />
      )}
    </>
  );
};

export default OrganizationMembers;
