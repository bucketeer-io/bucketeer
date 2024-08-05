import { PropsWithChildren } from 'react';

export type ListGroupProps = PropsWithChildren;

const ListGroup = ({ children }: ListGroupProps) => {
  return <ul>{children}</ul>;
};

export default ListGroup;
