import ListItem, { ListItemProps } from './list-item';
import ListTitle from './list-title';

export type ListProps = {
  title: string;
  options?: ListItemProps[];
};

const List = ({ title, options = [] }: ListProps) => {
  return (
    <div>
      <ListTitle text={title} />
      {options.map((item, index) => (
        <ListItem key={index} {...item} />
      ))}
    </div>
  );
};

List.Item = ListItem;
List.Title = ListTitle;

export default List;
