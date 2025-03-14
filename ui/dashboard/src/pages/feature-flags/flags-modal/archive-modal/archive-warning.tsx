import { Trans } from 'react-i18next';
import { IconToastWarning } from '@icons';
import Icon from 'components/icon';

const ArchiveWarning = ({ days }: { days: number }) => {
  return (
    <div className="flex items-center w-full p-4 gap-x-2 rounded-xl bg-accent-yellow-50 text-accent-yellow-500 typo-para-medium">
      <Icon
        icon={IconToastWarning}
        size={'xs'}
        color="accent-yellow-500"
        className="flex-center"
      />
      <div>
        <Trans
          i18nKey={'table:feature-flags.archive-warning'}
          values={{ days: `${days} day${days > 1 ? 's' : ''}` }}
          components={{ text: <span /> }}
        />
      </div>
    </div>
  );
};

export default ArchiveWarning;
