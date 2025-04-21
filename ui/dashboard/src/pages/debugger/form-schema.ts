import * as yup from 'yup';

export const addDebuggerFormSchema = yup.object().shape({
  flags: yup
    .array()
    .of(yup.string().required('This field is required.'))
    .required('This field is required.'),
  userIds: yup
    .array()
    .of(yup.string().required('This field is required.'))
    .required('This field is required.'),
  attributes: yup.array().of(
    yup.object().shape({
      key: yup.string(),
      value: yup.string()
    })
  )
});

export type AddDebuggerFormType = yup.InferType<typeof addDebuggerFormSchema>;
