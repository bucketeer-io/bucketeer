import { useMemo } from 'react';
import { useLocation, useNavigate, useParams } from 'react-router-dom';
import { ID_NEW } from 'constants/routing';
import { AnyObject } from 'yup';
import { useToast } from './use-toast';

interface Props {
  idKey?: string;
  addPath?: string;
  closeModalPath?: string;
}

const useActionWithURL = ({ idKey = '*', addPath, closeModalPath }: Props) => {
  const { [idKey]: action, ...params } = useParams();
  const navigate = useNavigate();
  const { notify } = useToast();
  const location = useLocation();

  const { state } = location;
  const isAdd = useMemo(() => action === ID_NEW, [action]);
  const isEdit = useMemo(() => action && !isAdd, [action, isAdd]);

  const onOpenAddModal = () =>
    navigate(addPath || `${location.pathname}/${ID_NEW}`);

  const onOpenEditModal = (path: string, state?: AnyObject) =>
    navigate(path, {
      state
    });

  const onCloseActionModal = (path?: string) => {
    if (closeModalPath || path) navigate(String(closeModalPath || path));
  };

  const errorToast = (error: Error) =>
    notify({
      toastType: 'toast',
      messageType: 'error',
      message: error?.message || 'Something went wrong.'
    });

  return {
    id: action,
    isAdd,
    isEdit,
    state,
    onOpenAddModal,
    onOpenEditModal,
    onCloseActionModal,
    errorToast,
    params
  };
};

export default useActionWithURL;
