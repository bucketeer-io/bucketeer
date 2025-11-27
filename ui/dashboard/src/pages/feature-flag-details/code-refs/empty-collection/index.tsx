import { IconAddOutlined, IconLaunchOutlined } from 'react-icons-material-design';
import { Trans } from 'react-i18next';
import { Link } from 'react-router-dom';
import { DOCUMENTATION_LINKS } from 'constants/documentation-links';
import { useTranslation } from 'i18n';
import Icon from 'components/icon';
import EmptyState, { EmptyStateProps } from 'elements/empty-state';

const EmptyCollection = ({
  variant,
  onAdd
}: {
  variant: EmptyStateProps['variant'];
  onAdd?: () => void;
}) => {
  const { t } = useTranslation(['table', 'common']);

  const isEmpty = variant === 'no-data';

  return (
    <EmptyState.Root variant={variant} size="lg" className="mt-10">
      <EmptyState.Illustration />
      <EmptyState.Body>
        <EmptyState.Title>
          {t(`code-refs.${isEmpty ? 'empty' : 'enable'}`)}
        </EmptyState.Title>
        <EmptyState.Description>
          {isEmpty ? (
            <Trans
              i18nKey={'table:code-refs.empty-desc'}
              components={{
                comp: (
                  <Link
                    to={DOCUMENTATION_LINKS.FLAG_CODE_REFERENCE}
                    target="_blank"
                    className="inline-flex items-center gap-x-1 text-primary-500 underline"
                  />
                ),
                icon: <Icon icon={IconLaunchOutlined} size="sm" />
              }}
            />
          ) : (
            t(`code-refs.enable-desc`)
          )}
        </EmptyState.Description>
      </EmptyState.Body>
      {!isEmpty && onAdd && (
        <EmptyState.Actions>
          <EmptyState.ActionButton
            isNeedAdminAccess
            type={'new'}
            variant="primary"
            onClick={onAdd}
          >
            <IconAddOutlined />
            {t(`common:create-api-key`)}
          </EmptyState.ActionButton>
        </EmptyState.Actions>
      )}
    </EmptyState.Root>
  );
};

export default EmptyCollection;
