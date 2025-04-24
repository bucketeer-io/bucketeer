import React from 'react';
import { cn } from 'utils/style';
import { useFormField } from 'components/form';

const FormMessage = React.forwardRef<
  HTMLParagraphElement,
  React.HTMLAttributes<HTMLParagraphElement>
>(({ className, children, ...props }, ref) => {
  const { error, formMessageId } = useFormField();
  const body = error?.message ? String(error?.message) : children;

  if (!body) {
    return null;
  }

  return (
    <p
      ref={ref}
      id={formMessageId}
      className={cn('typo-para-small text-accent-red-500 mt-0.5', className)}
      {...props}
    >
      {body}
    </p>
  );
});

export default FormMessage;
