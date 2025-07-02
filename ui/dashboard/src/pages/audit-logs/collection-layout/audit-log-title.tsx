import { Trans } from 'react-i18next';
import { Link, useNavigate } from 'react-router-dom';
import { useTranslation } from 'i18n';
import { DomainEventEntityType } from '@types';
import { getPathName } from '../utils';
import useEntityLocalized from './use-entity-localized';

interface Props {
  isHaveEntityData: boolean;
  entityId?: string;
  username: string;
  action: string;
  entityType: DomainEventEntityType;
  entityName: string;
  urlCode: string;
  additionalText?: string;
}

const AuditLogTitle = ({
  isHaveEntityData,
  entityId,
  username,
  action,
  entityType,
  entityName,
  urlCode,
  additionalText
}: Props) => {
  useTranslation('table');
  const { getEntityLocalized } = useEntityLocalized();
  const navigate = useNavigate();

  return (
    <Trans
      i18nKey={
        isHaveEntityData
          ? 'table:audit-log-title'
          : 'table:audit-log-title-no-entity'
      }
      values={{
        username,
        action,
        entityType: getEntityLocalized(entityType),
        entityName,
        additionalText
      }}
      components={{
        b: <span className="font-bold text-gray-700" />,
        highlight: (
          <Link
            to={getPathName(entityId, entityType) as string}
            onClick={e => {
              e.preventDefault();
              const pathName = getPathName(entityId, entityType);
              if (pathName) navigate(`/${urlCode}${pathName}`);
            }}
            className="text-primary-500 underline truncate"
          />
        )
      }}
    />
  );
};

export default AuditLogTitle;
