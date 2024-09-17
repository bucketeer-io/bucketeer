import primaryAvatar from 'assets/avatars/primary.svg';
import { AvatarImage } from 'components/avatar';
import { FlagProps } from './flag';
import Title from './title';

export type MemberProps = FlagProps;

const Member = ({ text, description }: MemberProps) => {
  return (
    <div className="flex gap-2">
      <div className="bg-primary-50 size-8 rounded-full flex-center">
        <AvatarImage size={'md'} image={primaryAvatar} />
      </div>
      <Title text={text} description={description} />
    </div>
  );
};

export default Member;
