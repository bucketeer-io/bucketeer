import type { ReactElement } from 'react';
import EmptyState from 'elements/empty-state';

type CollectionEmptyProps<Data extends object> = {
  data: Data[];
  empty: ReactElement;
  searchQuery?: string;
  onClear?: () => void;
};

export const NoResultsCollection = ({ onClear }: { onClear?: () => void }) => {
  return (
    <EmptyState.Root variant="no-search" size="lg">
      <EmptyState.Illustration />
      <EmptyState.Body>
        <EmptyState.Title>{`No results found`}</EmptyState.Title>
        <EmptyState.Description>{`We couldn't find what you're looking for`}</EmptyState.Description>
      </EmptyState.Body>
      {onClear && (
        <EmptyState.Actions>
          <EmptyState.ActionButton variant="primary" onClick={onClear}>
            {`Clear search & filters`}
          </EmptyState.ActionButton>
        </EmptyState.Actions>
      )}
    </EmptyState.Root>
  );
};

const CollectionEmpty = <Data extends object>({
  data,
  empty,
  searchQuery,
  onClear
}: CollectionEmptyProps<Data>) => {
  if (data.length === 0) {
    return searchQuery ? (
      <NoResultsCollection onClear={onClear} />
    ) : (
      <div className="h-full flex-center">{empty}</div>
    );
  }
};

export default CollectionEmpty;
