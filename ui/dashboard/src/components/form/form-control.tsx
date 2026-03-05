import React from 'react';
import { Slot } from '@radix-ui/react-slot';
import { useFormField } from 'components/form';

const FormControl = React.forwardRef<
  React.ElementRef<typeof Slot>,
  React.ComponentPropsWithoutRef<typeof Slot>
>(({ ...props }, ref) => {
  const { error, formItemId, formMessageId } = useFormField();

  const ariaDescribedBy = error ? formMessageId : undefined;

  return (
    <Slot
      ref={ref}
      id={formItemId}
      aria-describedby={ariaDescribedBy}
      aria-invalid={!!error}
      {...props}
    />
  );
});

FormControl.displayName = 'FormControl';

export default FormControl;
