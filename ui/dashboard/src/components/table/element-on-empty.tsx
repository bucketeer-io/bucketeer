import { PropsWithChildren } from 'react';

const ElementOnEmpty = ({ children }: PropsWithChildren) => {
  return (
    <div className="min-h-[526px] grid place-items-center">{children}</div>
  );
};

export default ElementOnEmpty;
