import primaryAvatar from 'assets/avatars/primary.svg';
import { AvatarImage } from 'components/avatar';
import StatusTag, { StatusTagType } from 'components/tag/status-tag';
import { FlagProps } from './flag';
import Text from './text';
import Title from './title';

export type MemberProps = FlagProps & {
  status?: StatusTagType;
  statusLabel?: string;
};

const Member = ({ text, description, status, statusLabel }: MemberProps) => {
  return (
    <div className="flex gap-2">
      <div className="bg-primary-50 size-8 rounded-full flex-center">
        <AvatarImage size={'md'} image={primaryAvatar} />
      </div>
      {!status ? (
        <Title
          text={text}
          description={description}
          descClassName={'text-gray-700'}
        />
      ) : (
        <div className="flex flex-col gap-2">
          <Text text={description} />
          <StatusTag variant={status} label={statusLabel} />
        </div>
      )}
    </div>
  );
};

export default Member;
