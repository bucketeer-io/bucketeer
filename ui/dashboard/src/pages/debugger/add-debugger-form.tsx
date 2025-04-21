import { useTranslation } from 'i18n';
import Button from 'components/button';
import DebuggerAttributes from './debugger-attributes';
import DebuggerFlags from './debugger-flags';
import DebuggerUserIds from './debugger-user-ids';

const AddDebuggerForm = ({ isLoading }: { isLoading: boolean }) => {
  const { t } = useTranslation(['common']);
  return (
    <div className="flex flex-col w-full gap-y-6 p-6">
      <DebuggerFlags />
      <DebuggerUserIds />
      <DebuggerAttributes />
      <Button className="w-fit" loading={isLoading}>
        {t('evaluate')}
      </Button>
    </div>
  );
};

export default AddDebuggerForm;
