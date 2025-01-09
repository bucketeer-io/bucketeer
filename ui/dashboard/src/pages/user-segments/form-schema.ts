import * as yup from 'yup';
import {
  SEGMENT_MAX_FILE_SIZE,
  SEGMENT_SUPPORTED_FORMATS
} from 'pages/user-segments/constants';

export const formSchema = yup.object().shape({
  name: yup.string().required(),
  description: yup.string(),
  id: yup.string(),
  userIds: yup.string(),
  file: yup
    .mixed()
    .nullable()
    .test('fileSize', 'The maximum size of the file is 2MB', value => {
      return !value || (value as File)?.size <= SEGMENT_MAX_FILE_SIZE;
    })
    .test('fileType', 'The file format is not supported', value => {
      return !value || SEGMENT_SUPPORTED_FORMATS.includes((value as File).type);
    })
});
