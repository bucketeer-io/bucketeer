import Text, { TextProps } from './text';

export type TitleProps = TextProps;

const Title = ({ text, description, sub }: TitleProps) => {
  return <Text text={text} description={description} sub={sub} isLink />;
};

export default Title;
