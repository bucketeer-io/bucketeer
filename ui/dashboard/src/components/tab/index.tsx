import TabItem, { TabItemProps } from './tab-item';

export type TabProps = {
  options: TabItemProps[];
  value: string;
  onSelect: (value: string) => void;
};

const Tab = ({ options, value, onSelect }: TabProps) => {
  const handleSelect = (tabValue: string) => {
    onSelect(tabValue);
  };

  return (
    <ul className="flex border-b border-gray-300">
      {options.map((option, index) => (
        <TabItem
          key={index}
          {...option}
          selected={option.value === value}
          onClick={handleSelect}
        />
      ))}
    </ul>
  );
};

Tab.Item = TabItem;

export default Tab;
