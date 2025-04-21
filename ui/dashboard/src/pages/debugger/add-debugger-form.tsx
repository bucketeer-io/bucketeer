import { useTranslation } from 'i18n';
import Button from 'components/button';
import DebuggerAttributes from './debugger-attributes';
import DebuggerFlags from './debugger-flags';
import DebuggerUserIds from './debugger-user-ids';

const AddDebuggerForm = () => {
  const { t } = useTranslation(['common']);
  return (
    <div className="flex flex-col w-full gap-y-6">
      <DebuggerFlags />
      <DebuggerUserIds />
      <DebuggerAttributes />
      <Button className="w-fit">{t('evaluate')}</Button>
    </div>
  );
};

export default AddDebuggerForm;
