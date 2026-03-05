import React from 'react';
import { useFormContext } from 'react-hook-form';
import { cn } from 'utils/style';
import { FormFieldContext, FormItemContext } from 'components/form';

const FormMessage = React.forwardRef<
  HTMLParagraphElement,
  React.HTMLAttributes<HTMLParagraphElement>
>(({ className, children, ...props }, ref) => {
  const fieldContext = React.useContext(FormFieldContext);
  const itemContext = React.useContext(FormItemContext);
  const formContext = useFormContext();

  const error =
    fieldContext && formContext
      ? formContext.getFieldState(fieldContext.name, formContext.formState)
          .error
      : undefined;

  const formMessageId = itemContext
    ? `${itemContext.id}-form-item-message`
    : undefined;

  const errorMessage = error?.message ? String(error.message) : null;

  if (!errorMessage && !children) {
    return null;
  }

  return (
    <p
      ref={ref}
      id={formMessageId}
      className={cn('typo-para-small mt-0.5', className, {
        'text-accent-red-500': !!errorMessage,
        'text-gray-500': !errorMessage
      })}
      {...props}
    >
      {errorMessage ?? children}
    </p>
  );
});

FormMessage.displayName = 'FormMessage';

export default FormMessage;
