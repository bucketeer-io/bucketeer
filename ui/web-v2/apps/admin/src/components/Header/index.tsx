import { FC, memo } from 'react';

export interface HeaderProps {
  title: string;
  description: string;
}

export const Header: FC<HeaderProps> = memo(({ title, description }) => {
  return (
    <div
      id="header"
      className="bg-white px-10 py-4 text-gray-700 border-b border-gray-300"
    >
      <p className="text-xl">{title}</p>
      <p className="text-sm">{description}</p>
    </div>
  );
});
