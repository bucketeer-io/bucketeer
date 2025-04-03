import { cn, getVariationColor } from 'utils/style';
import { IconInfo } from '@icons';
import { Polygon } from 'pages/experiment-details/elements/header-details';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';
import NameWithTooltip from 'elements/name-with-tooltip';

export const ResultHeaderCell = ({
  text,
  minSize,
  tooltip,
  isShowIcon = true,
  className,
  isFormatText
}: {
  text: string;
  minSize: number;
  tooltip: string;
  isShowIcon?: boolean;
  className?: string;
  isFormatText?: boolean;
}) => {
  const formatText = isFormatText ? text.replace(' ', '<br />') : text;
  return (
    <div
      className={cn(
        'flex items-center size-fit w-full p-4 pt-0 gap-x-3 text-[13px] leading-[13px] text-gray-500 uppercase relative z-10',
        className
      )}
      style={{
        minWidth: minSize
      }}
    >
      <p
        dangerouslySetInnerHTML={{
          __html: formatText
        }}
      />
      {isShowIcon && tooltip && (
        <Tooltip
          trigger={
            <div className="flex-center size-fit">
              <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
            </div>
          }
          content={tooltip}
          className="max-w-[350px]"
        />
      )}
    </div>
  );
};

export const ResultCell = ({
  variationId,
  value,
  minSize,
  isFirstItem,
  className,
  currentIndex
}: {
  variationId?: string;
  value: string | number | boolean;
  minSize: number;
  isFirstItem?: boolean;
  className?: string;
  currentIndex?: number;
}) => {
  const isBooleanValue = ['true', 'false'].includes(value as string);
  const id = variationId || '';

  return (
    <div
      className={cn(
        'flex items-center size-fit w-full px-4 py-5 gap-x-2 text-gray-500',
        className
      )}
      style={{ minWidth: minSize }}
    >
      {isFirstItem && typeof currentIndex === 'number' && (
        <Polygon
          className="border-none size-3"
          style={{
            background: getVariationColor(currentIndex),
            zIndex: currentIndex
          }}
        />
      )}

      {isFirstItem ? (
        <NameWithTooltip
          id={id}
          maxLines={1}
          content={<NameWithTooltip.Content content={value} id={id} />}
          trigger={
            <NameWithTooltip.Trigger
              id={id}
              name={String(value)}
              maxLines={1}
              haveAction={false}
              className={cn('typo-para-medium text-gray-800', {
                capitalize: isBooleanValue
              })}
            />
          }
        />
      ) : (
        <p
          className={cn('typo-para-medium leading-4 text-gray-800', {
            capitalize: isBooleanValue
          })}
        >
          {String(value)}
        </p>
      )}
    </div>
  );
};
