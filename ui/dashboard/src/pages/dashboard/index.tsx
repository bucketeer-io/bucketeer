import Button from 'components/button';
import { IconAddRound } from 'react-icons-material-design';

const DashboardPage = () => {
  return (
    <div className="container p-10">
      <div className="pb-4 text-3xl font-bold">{`Design systems`}</div>

      <div className="py-2">{`Button types`}</div>
      <div className="flex items-center gap-6">
        <Button>{`Primary`}</Button>
        <Button variant="secondary">{`Secondary`}</Button>
        <Button variant="negative">{`Negative`}</Button>
        <Button disabled>{`Disabled`}</Button>
        <Button variant="text">{`Text Button`}</Button>
      </div>

      <div className="mt-6 py-2">{`Button sizes`}</div>
      <div className="flex items-center gap-6">
        <Button>{`Button 1`}</Button>
        <Button size="sm">{`Button 1`}</Button>
        <Button size="xs">{`Button 1`}</Button>
      </div>

      <div className="mt-6 py-2">{`Button icons`}</div>
      <div className="flex items-center gap-6">
        <Button icon={IconAddRound} iconSlot="left">{`Button 1`}</Button>
        <Button variant="text" icon={IconAddRound} iconSlot="left">
          {`Text button`}
        </Button>
      </div>
    </div>
  );
};

export default DashboardPage;
