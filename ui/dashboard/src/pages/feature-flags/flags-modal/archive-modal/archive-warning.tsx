import { Trans } from 'react-i18next';
import { IconToastWarning } from '@icons';
import Icon from 'components/icon';

const ArchiveWarning = () => {
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
          values={{ days: `7 days` }}
          components={{ text: <span /> }}
        />
      </div>
    </div>
  );
};

export default ArchiveWarning;
