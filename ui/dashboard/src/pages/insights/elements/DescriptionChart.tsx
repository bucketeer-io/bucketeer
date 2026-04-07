import { ReactNode } from 'react';

interface ChartDescriptionProps {
  title: ReactNode;
  notes: ReactNode[];
}

const ChartDescription = ({ title, notes }: ChartDescriptionProps) => {
  return (
    <div className="w-fit">
      <p className="typo-para-medium font-bold">{title}</p>
      <div className="mt-2 space-y-1">
        {notes.map((note, i) => (
          <p key={i} className="[&>b]:inline-block [&>b]:min-w-[40px]">
            {note}
          </p>
        ))}
      </div>
    </div>
  );
};

export default ChartDescription;
