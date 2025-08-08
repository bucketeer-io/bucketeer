import { getCurrentEnvironment, useAuth, useAuthAccess } from 'auth';
import { useTranslation } from 'i18n';
import { Account } from '@types';
import { joinName } from 'utils/name';
import { useFetchTags } from 'pages/members/collection-loader';
import { useFetchEnvironments } from 'pages/project-details/environments/collection-loader/use-fetch-environments';
import Divider from 'components/divider';
import SlideModal from 'components/modal/slide';
import Spinner from 'components/spinner';
import { Tag } from 'elements/expandable-tag';

interface MemberDetailsModalProps {
  isOpen: boolean;
  member: Account;
  onClose: () => void;
}

const MemberDetailsModal = ({
  isOpen,
  member,
  onClose
}: MemberDetailsModalProps) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { isOrganizationAdmin } = useAuthAccess();
  const { t } = useTranslation(['common', 'form']);

  const { data: collection, isLoading } = useFetchEnvironments({
    organizationId: currentEnvironment.organizationId
  });

  const { data: tagCollection } = useFetchTags({
    organizationId: currentEnvironment.organizationId
  });

  const environments = collection?.environments || [];
  const tagList = tagCollection?.tags || [];

  return (
    <SlideModal
      isOpen={isOpen}
      title={t(`form:member-details`)}
      onClose={onClose}
    >
      {isLoading ? (
        <div className="w-full flex-center py-12">
          <Spinner />
        </div>
      ) : (
        <div className="flex flex-col w-full p-5 gap-y-5">
          <h3 className="typo-head-bold-small text-gray-800">
            {t(`form:general-info`)}
          </h3>
          <div className="flex items-start w-full gap-x-4">
            <div className="flex-1">
              <p className="typo-para-small text-gray-600">{t('name')}</p>
              <p className="text-gray-700 mt-1 typo-para-medium">
                {joinName(member.firstName, member.lastName) || member.name}
              </p>
            </div>
            <div className="flex-1">
              <p className="typo-para-small text-gray-600">{t('email')}</p>
              <p className="text-gray-700 mt-1 typo-para-medium break-all">
                {member.email}
              </p>
            </div>
          </div>
          <div className="flex items-start w-full gap-x-4">
            <div className="flex flex-col w-full gap-y-1 flex-1">
              <p className="typo-para-small text-gray-600">{t('teams')}</p>
              <div className="flex items-center flex-wrap w-full max-w-full gap-2">
                {member?.teams.map(item => (
                  <Tag
                    key={item}
                    tooltipCls={'!max-w-[450px]'}
                    tagId={item}
                    maxSize={487}
                    value={tagList?.find(tag => tag.id === item)?.name || item}
                  />
                ))}
              </div>
            </div>
            <div className="flex-1">
              <p className="typo-para-small text-gray-600">{t('role')}</p>
              <p className="text-gray-700 mt-1 typo-para-medium break-all">
                {t(
                  String(member.organizationRole).split('_')[1]?.toLowerCase()
                )}
              </p>
            </div>
          </div>
          <Divider />
          <h3 className="typo-head-bold-small text-gray-800">
            {t(
              isOrganizationAdmin ? `form:env-admin-access` : `form:env-access`
            )}
          </h3>
          {(member.environmentRoles || []).map((env, index) => (
            <div className="flex items-start w-full gap-x-4" key={index}>
              <div className="flex-1">
                <p className="typo-para-small text-gray-600">
                  {t('environment')}
                </p>
                <p className="text-gray-700 mt-1 typo-para-medium">
                  {
                    environments?.find(item => item.id === env.environmentId)
                      ?.name
                  }
                </p>
              </div>
              <div className="flex-1">
                <p className="typo-para-small text-gray-600">{t('role')}</p>
                <p className="text-gray-700 mt-1 capitalize typo-para-medium">
                  {env?.role === 'Environment_EDITOR'
                    ? t('editor')
                    : t('viewer')}
                </p>
              </div>
            </div>
          ))}
        </div>
      )}
    </SlideModal>
  );
};

export default MemberDetailsModal;
