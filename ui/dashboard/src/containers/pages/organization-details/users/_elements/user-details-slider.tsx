import { useTranslation } from 'i18n';
import { Account, Environment } from '@types';
import Divider from 'components/divider';
import SlideModal from 'components/modal/slide';
import Text from 'components/table/table-row-items/text';

type Props = {
  isOpenSlider: boolean;
  accountSelected?: Account;
  environmentData?: Environment[];
  setIsOpenSlider: (value: boolean) => void;
  setAccountSelected: (account: Account | undefined) => void;
};

const UserDetailsSlider = ({
  isOpenSlider,
  accountSelected,
  environmentData,
  setIsOpenSlider,
  setAccountSelected
}: Props) => {
  const { t } = useTranslation(['common', 'form']);

  return (
    <SlideModal
      isOpen={isOpenSlider}
      title={t(`form:user-details`)}
      onClose={() => {
        setIsOpenSlider(false);
        setAccountSelected(undefined);
      }}
    >
      <div className="flex flex-col w-full p-5 gap-y-5">
        <h3 className="typo-head-bold-small text-gray-800">
          {t(`form:general-info`)}
        </h3>
        <div className="flex items-center w-full gap-x-4">
          <div className="flex flex-1">
            <Text
              text={accountSelected?.name}
              description={t('name')}
              isReverse
              descClassName="text-sm text-gray-600"
              className="line-clamp-1"
            />
          </div>
          <div className="flex flex-1">
            <Text
              text={accountSelected?.email}
              description={t('email')}
              isReverse
              descClassName="text-sm text-gray-600"
            />
          </div>
        </div>
        <Divider />
        <h3 className="typo-head-bold-small text-gray-800">
          {t(`form:env-access`)}
        </h3>
        {accountSelected?.environmentRoles?.map((env, index) => (
          <div className="flex items-center w-full gap-x-4" key={index}>
            <div className="flex flex-1">
              <Text
                text={
                  environmentData?.find(item => item.id === env.environmentId)
                    ?.name
                }
                description={t('name')}
                isReverse
                descClassName="text-sm text-gray-600"
                className="line-clamp-1"
              />
            </div>
            <div className="flex flex-1">
              <Text
                text={env?.role}
                description={t('role')}
                isReverse
                descClassName="text-sm text-gray-600"
              />
            </div>
          </div>
        ))}
      </div>
    </SlideModal>
  );
};

export default UserDetailsSlider;
