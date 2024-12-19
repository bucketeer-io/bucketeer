import { i18n } from 'i18n';
// import * as pdfjs from 'pdfjs-dist';
// import type { PDFDocumentProxy } from 'pdfjs-dist';

// export const readFile = async (file: File): Promise<pdfjs.PDFDocumentProxy> => {
//   return new Promise((resolve, reject) => {
//     const fileReader = new FileReader();
//     fileReader.readAsBinaryString(file);
//     fileReader.onload = event => {
//       if (event.target) {
//         if (fileReader.result && typeof fileReader.result === 'string') {
//           const loadingTask = pdfjs.getDocument({ data: fileReader.result });
//           loadingTask.promise.then(pdf => resolve(pdf)).catch(reject);
//         }
//       }
//     };
//   });
// };

export function formatBytes(bytes: number, decimals: number = 2) {
  if (!+bytes) return '0 Bytes';

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return `${parseFloat((bytes / k ** i).toFixed(dm))} ${sizes[i]}`;
}

export const readImageFile = (file: File) => {
  return new Promise<string>(resolve => {
    const reader = new FileReader();
    reader.addEventListener(
      'load',
      () => resolve(reader.result as string),
      false
    );
    reader.readAsDataURL(file);
  });
};

export const convertBase64ToFile = async (
  base64Data: string,
  fileName: string = 'file'
) => {
  const base64Response = await fetch(base64Data);
  const blob = await base64Response.blob();
  const file = new File([blob], fileName, {
    type: blob.type
  });
  return file;
};

// const MAX_PAGE_NUMBER = 100;
// export const validateFile = async (file: File) => {
//   // No extra validate if the uploaded file is Image
//   if (file.type.startsWith('image/')) return file;

//   // Validate if the uploaded file is PDF and
//   // * is protected, or
//   // * contains too many pages, or
//   // * including other error cases: concruppted file or for some other reason
//   return readFile(file)
//     .catch(error => {
//       // eslint-disable-next-line prefer-promise-reject-errors
//       return Promise.reject({
//         title: 'Invalid file!',
//         message: error.message
//       });
//     })
//     .then((pdf: PDFDocumentProxy) => {
//       // Check if the file contains too many pages
//       if (pdf.numPages > MAX_PAGE_NUMBER) {
//         // eslint-disable-next-line prefer-promise-reject-errors
//         return Promise.reject({
//           title: i18n.t('file-upload.error-too-many-pages-title'),
//           message: i18n.t('file-upload.error-too-many-pages-message', {
//             pages: MAX_PAGE_NUMBER
//           })
//         });
//       }
//       return file;
//     });
// };