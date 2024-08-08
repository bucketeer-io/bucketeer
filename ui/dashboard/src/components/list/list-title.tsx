export type ListTitleProps = {
  text: string;
};

const ListTitle = ({ text }: ListTitleProps) => {
  return <h3 className="typo-head-bold-medium h-10 text-gray-700">{text}</h3>;
};

export default ListTitle;
