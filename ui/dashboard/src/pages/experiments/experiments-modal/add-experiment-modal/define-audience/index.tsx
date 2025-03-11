import { DefineAudienceField } from '..';
import DefineAudienceAdvanced from './define-audience-advanced';
import DefineAudienceAmount from './define-audience-amount';
import DefineAudienceRule from './define-audience-rule';

export interface DefineAudienceProps {
  field: DefineAudienceField;
}

const DefineAudience = ({ field }: DefineAudienceProps) => {
  return (
    <>
      <DefineAudienceRule field={field} />
      <DefineAudienceAmount field={field} />
      <DefineAudienceAdvanced field={field} />
    </>
  );
};

export default DefineAudience;
