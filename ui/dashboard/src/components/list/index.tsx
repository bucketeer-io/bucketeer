import ListGroup from './list-group';
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
      <ListGroup>
        {options.map((item, index) => (
          <ListItem key={index} {...item} />
        ))}
      </ListGroup>
    </div>
  );
};

List.Group = ListGroup;
List.Item = ListItem;
List.Title = ListTitle;

export default List;
