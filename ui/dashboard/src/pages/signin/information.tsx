import { FormProvider, SubmitHandler, useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import { UserInfoForm } from '@types';
import { Button } from 'components/button';
import Form from 'components/form';
import Input from 'components/input';
import AuthWrapper from './elements/auth-wrapper';

const formSchema = yup.object().shape({
  first_name: yup
    .string()
    .required()
    .min(2, 'The first name you have provided must have at least 2 characters'),
  last_name: yup
    .string()
    .required()
    .min(2, 'The last name you have provided must have at least 2 characters'),
  language: yup.string().required()
});

const UserInformation = () => {
  const form = useForm({
    resolver: yupResolver(formSchema),
    defaultValues: {
      first_name: '',
      last_name: '',
      language: ''
    }
  });

  const onSubmit: SubmitHandler<UserInfoForm> = values => {
    console.log(values);
  };

  return (
    <AuthWrapper>
      <div className="grid gap-10">
        <div>
          <h1 className="text-gray-900 typo-head-bold-huge mb-4">
            {`Enter Your Information`}
          </h1>
          <p className="text-gray-600 typo-para-medium">
            {`Enter your information to complete your account creation.`}
          </p>
        </div>

        <FormProvider {...form}>
          <Form onSubmit={form.handleSubmit(onSubmit)} className="mt-8">
            <Form.Field
              control={form.control}
              name="first_name"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required>{`First Name`}</Form.Label>
                  <Form.Control>
                    <Input placeholder="Enter your first name" {...field} />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            <Form.Field
              control={form.control}
              name="last_name"
              render={({ field }) => (
                <Form.Item>
                  <Form.Label required>{`Last Name`}</Form.Label>
                  <Form.Control>
                    <Input placeholder="Enter your last name" {...field} />
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            {/* <Select
              placeholder="Select your Language"
              options={[]}
              label="Language"
              className="mt-4"
              defaultValue={0}
              required
            /> */}
            <Button type="submit" className="mt-8 w-full">
              {`Sign In`}
            </Button>
          </Form>
        </FormProvider>
      </div>
    </AuthWrapper>
  );
};

export default UserInformation;
