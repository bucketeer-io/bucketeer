import { FormProvider, useForm } from 'react-hook-form';
import {
  IconAddRound,
  IconEditOutlined,
  IconPersonRound
} from 'react-icons-material-design';
import { yupResolver } from '@hookform/resolvers/yup';
import primaryAvatar from 'assets/avatars/primary.svg';
import { useToggleOpen } from 'hooks';
import * as yup from 'yup';
import { IconGoal } from '@icons';
import { AvatarIcon, AvatarImage } from 'components/avatar';
import { Badge } from 'components/badge';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import Divider from 'components/divider';
import Form from 'components/form';
import Icon from 'components/icon';
import Input from 'components/input';
import DialogModal from 'components/modal/dialog';
import SlideModal from 'components/modal/slide';

const formSchema = yup.object().shape({
  username: yup
    .string()
    .required()
    .min(2, 'Username must be at least 2 characters.')
});

const DashboardPage = () => {
  const [openModal, onOpenModal, onCloseModal] = useToggleOpen(false);
  const [openSlider, onOpenSlider, onCloseSlider] = useToggleOpen(false);

  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      username: ''
    }
  });

  function onSubmit(values: { username?: string }) {
    console.log(values);
  }

  return (
    <div className="p-10">
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
        <div className="typo-head-bold-huge font-sofia">{`Heading H1`}</div>
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
        <AvatarImage size="xl" image={primaryAvatar} />
        <AvatarImage size="md" image={primaryAvatar} />
        <AvatarImage size="sm" image={primaryAvatar} />
        <AvatarIcon icon={IconPersonRound} size="md" />

        <Badge>{'1'}</Badge>
        <Badge variant="secondary">{'1'}</Badge>
      </div>
      <div className="mt-8 flex flex-col gap-6">
        <Divider />

        <div className="flex gap-6">
          <Button onClick={onOpenModal} variant="secondary">{`Modal`}</Button>
          <Button onClick={onOpenSlider} variant="secondary">{`Slider`}</Button>
          <Input className="w-fit" />
          <Input className="w-fit" disabled value="Disabled" />
        </div>

        <SlideModal
          title={'New Environment'}
          isOpen={openSlider}
          onClose={onCloseSlider}
        >
          <div className="py-8 px-5 flex flex-col gap-6 items-center justify-center">
            <div className="typo-para-big text-gray-700 px-20 text-center">
              {`This experiment has the following goals connected to it`}
            </div>
          </div>

          <div className="absolute bottom-0 bg-gray-50 w-full rounded-b-lg">
            <ButtonBar
              primaryButton={<Button variant="secondary">{`Cancel`}</Button>}
              secondaryButton={<Button>{`Create Goal`}</Button>}
            />
          </div>
        </SlideModal>

        <DialogModal
          title={'Goals Connected'}
          isOpen={openModal}
          onClose={onCloseModal}
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
      <FormProvider {...form}>
        <Form onSubmit={form.handleSubmit(onSubmit)} className="mt-4">
          <Form.Field
            control={form.control}
            name="username"
            render={({ field }) => (
              <Form.Item>
                <Form.Label>{`Username`}</Form.Label>
                <Form.Control>
                  <Input placeholder="text" {...field} />
                </Form.Control>
                <Form.Message />
              </Form.Item>
            )}
          />
          <Button type="submit">Submit</Button>
        </Form>
      </FormProvider>
    </div>
  );
};

export default DashboardPage;
