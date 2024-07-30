import { useState } from 'react';
import {
  IconAddRound,
  IconEditOutlined,
  IconPersonRound
} from 'react-icons-material-design';
import { IconGoal } from '@icons';
import { AvatarIcon, AvatarImage } from 'components/avatar';
import { Badge } from 'components/badge';
import { Button } from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';

const DashboardPage = () => {
  const [open, setOpen] = useState(false);

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
        <Button>
          <Icon icon={IconAddRound} /> {`Button 1`}
        </Button>
        <Button variant="text">
          <Icon icon={IconAddRound} /> {`Text button`}
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

      <div className="mt-6 flex items-center gap-6">
        <Button size="icon">
          <Icon icon={IconEditOutlined} />
        </Button>
        <Button variant="secondary" size="icon">
          <Icon icon={IconEditOutlined} />
        </Button>
        <Button variant="grey" size="icon">
          <Icon icon={IconEditOutlined} />
        </Button>
        <Button variant="secondary" size="icon-sm">
          <Icon icon={IconEditOutlined} size="sm" />
        </Button>
        <Button variant="grey" size="icon-sm">
          <Icon icon={IconEditOutlined} size="sm" />
        </Button>
      </div>

      <div className="mt-10 flex items-center gap-6">
        <AvatarImage size="xl" image="./assets/avatars/primary.svg" />
        <AvatarImage size="md" image="./assets/avatars/primary.svg" />
        <AvatarImage size="sm" image="./assets/avatars/primary.svg" />
        <AvatarIcon icon={IconPersonRound} size="md" />

        <Badge>{'1'}</Badge>
        <Badge variant="secondary">{'1'}</Badge>
      </div>
      <div className="mt-8 flex flex-col gap-6">
        <Divider />

        <div className="flex">
          <Button
            onClick={() => setOpen(true)}
            variant="secondary"
          >{`Modal`}</Button>
        </div>

        <DialogModal
          title={'Goals Connected'}
          isOpen={open}
          onClose={() => setOpen(false)}
        >
          <div className="py-8 px-5 flex flex-col gap-6 items-center justify-center">
            <IconGoal />
            <div className="typo-para-big text-gray-700 px-20 text-center">
              {`This experiment has the following goals connected to it:`}
            </div>
            <div className="w-full rounded px-4 py-3 bg-gray-100">
              <div className="typo-para-medium">
                <span className="text-gray-700 mr-2">{`1.`}</span>
                <span className="text-primary-500 underline">
                  {`This is a big name for the first goal name`}
                </span>
              </div>
              <div className="typo-para-medium mt-3">
                <span className="text-gray-700 mr-2">{`2.`}</span>
                <span className="text-primary-500 underline">
                  {`This is a big name for the second goal name`}
                </span>
              </div>
            </div>
          </div>
          <ButtonBar primaryButton={<Button>{`Close`}</Button>} />
        </DialogModal>
      </div>
    </div>
  );
};

export default DashboardPage;
