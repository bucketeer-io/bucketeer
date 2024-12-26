import { cn } from 'utils/style';
import Button from 'components/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from 'components/dropdown';
import { Options } from './chart-header';

type TimeSwitchProps = {
  timeOptions: Options[];
  timeValue: string;
  onSelectTimeOption: (value: string) => void;
};

export type ViewActionsProps = TimeSwitchProps & {
  dropdownValue?: string;
  dropdownPlaceholder?: string;
  dropdownOptions?: Options[];
  onSelectDropdownOption?: (value: string) => void;
};

const TimeSwitches = ({
  timeOptions,
  timeValue,
  onSelectTimeOption
}: TimeSwitchProps) => {
  return (
    <div className="flex items-center">
      {timeOptions?.map((option, index) => (
        <Button
          key={index}
          variant={'text'}
          className={cn(
            'w-14 h-12 border border-gray-300 first:rounded-l-lg [&:not(:first-child)]:border-l-0 last:rounded-r-lg text-gray-500',
            {
              'text-primary-500': timeValue === option.value
            }
          )}
          onClick={() => onSelectTimeOption(option.value)}
        >
          {option.label}
        </Button>
      ))}
    </div>
  );
};

const ViewActions = ({
  timeValue,
  timeOptions,
  dropdownPlaceholder = 'All',
  dropdownValue,
  dropdownOptions,
  onSelectTimeOption,
  onSelectDropdownOption
}: ViewActionsProps) => {
  return (
    <div className="flex items-center w-fit gap-x-4">
      <TimeSwitches
        timeOptions={timeOptions}
        timeValue={timeValue}
        onSelectTimeOption={onSelectTimeOption}
      />
      {dropdownOptions && (
        <DropdownMenu>
          <DropdownMenuTrigger
            placeholder={dropdownPlaceholder}
            label={
              dropdownOptions?.find(item => item.value === dropdownValue)?.label
            }
            variant="secondary"
            className="w-[140px]"
          />
          <DropdownMenuContent
            className="w-[140px] min-w-[140px]"
            align="start"
          >
            {dropdownOptions?.map((item, index) => (
              <DropdownMenuItem
                key={index}
                value={item.value}
                label={item.label}
                onSelectOption={value => {
                  if (onSelectDropdownOption)
                    onSelectDropdownOption(value as string);
                }}
              />
            ))}
          </DropdownMenuContent>
        </DropdownMenu>
      )}
    </div>
  );
};

export default ViewActions;
