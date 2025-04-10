import type { ReactElement } from 'react';
import { ButtonProps } from 'components/button';
import EmptyState from 'elements/empty-state';

type CollectionEmptyProps<Data extends object> = {
  data: Data[];
  empty: ReactElement;
  isFilter?: boolean;
  searchQuery?: string;
  description?: string;
  buttonText?: string;
  buttonVariant?: ButtonProps['variant'];
  onClear?: () => void;
};

export const NoResultsCollection = ({
  description,
  buttonVariant = 'secondary',
  buttonText,
  onClear
}: {
  description?: string;
  buttonVariant?: ButtonProps['variant'];
  buttonText?: string;
  onClear?: () => void;
}) => {
  return (
    <EmptyState.Root variant="no-search" size="lg">
      <EmptyState.Illustration />
      <EmptyState.Body>
        <EmptyState.Title>{`No results found`}</EmptyState.Title>
        <EmptyState.Description>
          {description || `We couldn't find what you're looking for`}
        </EmptyState.Description>
      </EmptyState.Body>
      {onClear && (
        <EmptyState.Actions>
          <EmptyState.ActionButton variant={buttonVariant} onClick={onClear}>
            {buttonText || `Clear search & filters`}
          </EmptyState.ActionButton>
        </EmptyState.Actions>
      )}
    </EmptyState.Root>
  );
};

const CollectionEmpty = <Data extends object>({
  data,
  empty,
  isFilter,
  searchQuery,
  buttonText,
  buttonVariant,
  description,
  onClear
}: CollectionEmptyProps<Data>) => {
  if (data.length === 0) {
    return searchQuery || isFilter ? (
      <NoResultsCollection
        onClear={onClear}
        buttonText={buttonText}
        buttonVariant={buttonVariant}
        description={description}
      />
    ) : (
      <div className="h-full flex-center">{empty}</div>
    );
  }
};

export default CollectionEmpty;
