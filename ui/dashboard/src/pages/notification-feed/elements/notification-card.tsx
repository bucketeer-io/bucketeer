import { ReactNode } from 'react';
import { cn } from 'utils/style';
import Button from 'components/button';

interface NotificationCardProps {
  active?: boolean;
  onClick?: () => void;
  header: ReactNode;
  children: ReactNode;
  footer: ReactNode;
  // Lets a caller that already supplies its own bordered container (e.g. a
  // row with a checkbox prefix) drop this card's own border, so the two
  // don't render as two separate boxes.
  bordered?: boolean;
}

// Shared "Button styled as a card" shell for a notification/draft list item:
// a full-width, left-aligned clickable card with a header row, a tags row,
// and a footer line. Used by both the feed row and the draft card, which
// differ only in what they put in each slot (e.g. a checkbox prefix, an
// unread dot, or an active-selection border).
const NotificationCard = ({
  active,
  onClick,
  header,
  children,
  footer,
  bordered = true
}: NotificationCardProps) => (
  <Button
    type="button"
    variant="text"
    onClick={onClick}
    className={cn(
      // Override the Button base (centered, nowrap, fixed height) so it reads
      // as a full-width, left-aligned card.
      'flex h-auto w-full flex-col items-start justify-start gap-2 whitespace-normal text-left transition-colors',
      bordered
        ? cn(
            'rounded-lg border p-4',
            active
              ? 'border-primary-500 shadow-border-primary-500'
              : 'border-gray-200 hover:border-gray-300'
          )
        : 'p-0'
    )}
  >
    {header}
    <div className="flex flex-wrap items-center gap-2">{children}</div>
    {footer}
  </Button>
);

export default NotificationCard;
