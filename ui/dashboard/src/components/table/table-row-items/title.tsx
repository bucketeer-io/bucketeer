import { forwardRef, Ref } from 'react';
import Text, { TextProps } from './text';

export type TitleProps = TextProps;

const Title = forwardRef(
  (
    { text, description, sub, className, descClassName }: TitleProps,
    ref: Ref<HTMLDivElement>
  ) => {
    return (
      <Text
        ref={ref}
        text={text}
        description={description}
        sub={sub}
        isLink
        className={className}
        descClassName={descClassName}
      />
    );
  }
);
export default Title;
