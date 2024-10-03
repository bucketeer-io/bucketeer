import { IconAddOutlined } from 'react-icons-material-design';
import EmptyState from 'elements/empty-state';

export const EmptyCollection = ({ onAdd }: { onAdd: () => void }) => {
  return (
    <EmptyState.Root variant="no-data" size="md">
      <EmptyState.Illustration />
      <EmptyState.Body>
        <EmptyState.Title>{`No organizations yet!`}</EmptyState.Title>
        <EmptyState.Description>{`Your created organizations will appear here.`}</EmptyState.Description>
      </EmptyState.Body>
      <EmptyState.Actions>
        <EmptyState.ActionButton variant="primary" onClick={onAdd}>
          <IconAddOutlined />
          {`Add organization`}
        </EmptyState.ActionButton>
      </EmptyState.Actions>
    </EmptyState.Root>
  );
};
