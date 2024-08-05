import { PropsWithChildren } from 'react';

export type MenuGroupProps = PropsWithChildren<object>;

const MenuGroup = ({ children }: MenuGroupProps) => {
  return <ul>{children}</ul>;
};

export default MenuGroup;
