import type { ColumnDef } from '@tanstack/react-table';
import primaryAvatar from 'assets/avatars/primary.svg';
import { useTranslation } from 'i18n';
import { joinName } from 'utils/name';
import { AvatarImage } from 'components/avatar';

type TempAccountType = {
  name: string;
  firstName: string;
  lastName: string;
  email: string;
};

export type TempMemberType = {
  account: TempAccountType;
  total_flags_created: number;
};

export const useColumns = (): ColumnDef<TempMemberType>[] => {
  const { t } = useTranslation(['common', 'table']);

  return [
    {
      accessorKey: 'name',
      header: `${t('name')}`,
      size: 500,
      cell: ({ row }) => {
        const member = row.original;

        return (
          <div className="flex gap-2">
            <AvatarImage image={primaryAvatar} alt="member-avatar" />
            <div className="flex flex-col gap-0.5">
              <button
                // onClick={() => onActions(account)}
                className="underline text-primary-500 typo-para-medium text-left"
              >
                {joinName(
                  member?.account?.firstName,
                  member?.account?.lastName
                ) || member?.account?.name}
              </button>

              <div className="typo-para-medium text-gray-700">
                {member?.account?.email}
              </div>
            </div>
          </div>
        );
      }
    },
    {
      accessorKey: 'total_flags_created',
      header: `${t('table:total_flags_created')}`,
      size: 500
    }
  ];
};
