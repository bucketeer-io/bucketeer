import type { ReactNode } from 'react';
import { cva } from 'class-variance-authority';
import { cn } from 'utils/style';
import { InputGroupContext } from './context';

interface InputGroupProps {
  addon: ReactNode;
  children: ReactNode;
  className?: string;
  addonSlot?: 'left' | 'right';
  addonSize?: 'sm' | 'md' | 'lg';
  addonClassName?: string;
}

const inputGroupVariants = cva(['relative'], {
  variants: {
    addonSize: {
      sm: 'w-[28px]',
      md: 'w-[45px]',
      lg: 'w-[60px]'
    },
    addonSlot: {
      left: 'left-3 flex justify-start',
      right: 'right-3 flex justify-end'
    }
  }
});

const InputGroup = ({
  addon,
  addonSlot = 'left',
  addonSize = 'md',
  children,
  className,
  addonClassName
}: InputGroupProps) => {
  return (
    <InputGroupContext.Provider value={{ addonSlot, addonSize }}>
      <div className={cn(inputGroupVariants({ addonSize }), className)}>
        {children}
        <div
          className={cn(
            inputGroupVariants({ addonSlot }),
            'typo-para-medium absolute top-1/2 -translate-y-1/2',
            'flex items-center text-center text-gray-500',
            addonClassName
          )}
        >
          {addon}
        </div>
      </div>
    </InputGroupContext.Provider>
  );
};

export default InputGroup;
