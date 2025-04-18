import { useMemo } from 'react';
import primaryAvatar from 'assets/avatars/primary.svg';
import { AuditLogEditor } from '@types';
import { cn } from 'utils/style';
import { AvatarImage } from 'components/avatar';

const AuditLogAvatar = ({
  editor,
  className
}: {
  editor?: AuditLogEditor;
  className?: string;
}) => {
  const avatarSrc = useMemo(
    () =>
      editor?.avatarImage
        ? `data:${editor?.avatarFileType};base64,${editor?.avatarImage}`
        : primaryAvatar,
    [editor, primaryAvatar]
  );
  return (
    <AvatarImage
      image={avatarSrc}
      alt="member-avatar"
      className={cn('size-10', className)}
    />
  );
};

export default AuditLogAvatar;
