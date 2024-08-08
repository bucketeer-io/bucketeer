import { cn } from 'utils/style';
import ListItem, { ListItemProps } from './list-item';
import ListTitle from './list-title';

export type ListProps = {
  className?: string;
  items: ListItemProps[];
};

const List = ({ className, items = [] }: ListProps) => {
  return (
    <ul className={cn(className)}>
      {items.map((item, index) => (
        <ListItem key={index} {...item} />
      ))}
    </ul>
  );
};

List.Item = ListItem;
List.Title = ListTitle;

export default List;
