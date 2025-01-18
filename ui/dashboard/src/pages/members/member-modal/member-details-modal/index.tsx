import { useQueryTags } from '@queries/tags';
import { getCurrentEnvironment, useAuth } from 'auth';
import { LIST_PAGE_SIZE } from 'constants/app';
import { useTranslation } from 'i18n';
import { Account } from '@types';
import { joinName } from 'utils/name';
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
  const { t } = useTranslation(['common', 'form']);

  const { data: collection, isLoading } = useFetchEnvironments({
    organizationId: currentEnvironment.organizationId
  });

  const { data: tagCollection } = useQueryTags({
    params: {
      cursor: String(0),
      pageSize: LIST_PAGE_SIZE,
      environmentId: currentEnvironment.id,
      entityType: 'ACCOUNT'
    }
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
          <div className="flex flex-col w-full gap-y-1">
            <p className="typo-para-small text-gray-600">{t('tags')}</p>
            <div className="flex items-center flex-wrap w-full max-w-full gap-2">
              {member?.tags.map(tagId => (
                <Tag
                  key={tagId}
                  tooltipCls={'!max-w-[450px]'}
                  tagId={tagId}
                  maxSize={487}
                  value={tagList?.find(tag => tag.id === tagId)?.name || tagId}
                />
              ))}
            </div>
          </div>
          <Divider />
          <h3 className="typo-head-bold-small text-gray-800">
            {t(`form:env-access`)}
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
                  {env?.role.split('_')[1]}
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
