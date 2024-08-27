import React from 'react';
import { cn } from 'utils/style';
import { FormItemContext } from 'components/form';


const FormItem = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
  
>(({ className, ...props }, ref) => {
  const id = React.useId();

  return (
    <FormItemContext.Provider value={{ id }}>
      <div ref={ref} className={cn('py-3', className)} {...props} /> 
    </FormItemContext.Provider>
  );
});

export default FormItem;  
