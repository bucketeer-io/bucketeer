import { IconInfoFilled } from '@icons';
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
        <Button icon={IconInfoFilled} iconSlot="left">{`Button 1`}</Button>
        <Button variant="text" icon={IconAddRound} iconSlot="left">
          {`Text button`}
        </Button>
      </div>

      <div className="mt-8 flex items-center gap-6">
        <div className="typo-head-bold-huge">{`Heading H1`}</div>
        <div className="typo-head-bold-big">{`Heading H2`}</div>
        <div className="typo-head-bold-medium">{`Heading H3`}</div>
        <div className="typo-head-bold-small">{`Heading H4`}</div>
      </div>

      <div className="mt-4 flex items-center gap-6">
        <div className="typo-para-big">{`Paragraph LG`}</div>
        <div className="typo-para-medium">{`Paragraph MD`}</div>
        <div className="typo-para-small">{`Paragraph SM`}</div>
        <div className="typo-para-tiny">{`Paragraph XS`}</div>
      </div>
    </div>
  );
};

export default DashboardPage;
