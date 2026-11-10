import { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { IconCloseFilled } from 'react-icons-material-design';
import Icon from 'components/icon';
import { Popover } from 'components/popover';

interface Props {
  labels: string[];
  text: ReactNode;
  onRemove?: (index: number) => void;
}

// Text shown below a collapsed multi-value input, opening a list of every value
const OverflowValuesPopover = ({ labels, text, onRemove }: Props) => {
  const { t } = useTranslation(['common']);

  return (
    <div className="mt-0.5 flex items-center">
      <Popover
        align="start"
        trigger={
          <div>
            <span className="typo-para-small font-medium text-primary-500">
              {text}
            </span>
          </div>
        }
        triggerCls="w-fit justify-start"
        className="flex flex-col w-[300px] gap-y-0.5 p-0 overflow-hidden"
      >
        <div className="flex flex-wrap gap-1 max-h-[220px] overflow-y-auto small-scroll p-2">
          {labels.map((label, index) => (
            <div
              key={index}
              className="flex items-center max-w-full gap-x-2 px-1.5 rounded bg-primary-100"
            >
              <span
                title={label}
                className="truncate py-1 typo-para-small text-primary-500"
              >
                {label}
              </span>
              {onRemove && (
                <button
                  type="button"
                  aria-label={t('remove-value', { value: label })}
                  className="flex-center w-3 min-w-3 self-stretch text-primary-500 hover:text-gray-900"
                  onClick={() => onRemove(index)}
                >
                  <Icon
                    icon={IconCloseFilled}
                    className="h-full w-3 pointer-events-none"
                  />
                </button>
              )}
            </div>
          ))}
        </div>
      </Popover>
    </div>
  );
};

export default OverflowValuesPopover;
