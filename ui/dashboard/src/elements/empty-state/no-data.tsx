import EmptyState, { type EmptyStateProps } from 'elements/empty-state';

interface NoDataStateProps {
  size?: EmptyStateProps['size'];
  title?: string;
  description?: string;
  onAdd?: () => void;
}

export const NoDataState = ({
  size = 'md',
  title = `No Data`,
  description = `Your data will appear here!`,
  onAdd
}: NoDataStateProps) => {
  return (
    <EmptyState.Root variant="no-data" size={size}>
      <EmptyState.Illustration />
      <EmptyState.Body>
        <EmptyState.Title>{title}</EmptyState.Title>
        {description && (
          <EmptyState.Description>{description}</EmptyState.Description>
        )}
        {onAdd && (
          <EmptyState.Actions>
            <EmptyState.ActionButton type={'new'} onClick={onAdd}>
              {`Add`}
            </EmptyState.ActionButton>
          </EmptyState.Actions>
        )}
      </EmptyState.Body>
    </EmptyState.Root>
  );
};
