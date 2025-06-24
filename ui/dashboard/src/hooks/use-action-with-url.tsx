import { useCallback, useMemo } from 'react';
import { useLocation, useNavigate, useParams } from 'react-router-dom';
import { ID_NEW } from 'constants/routing';
import { AnyObject } from 'yup';

interface Props {
  idKey?: string;
  addPath?: string;
  closeModalPath?: string;
}

const useActionWithURL = ({ idKey = '*', addPath, closeModalPath }: Props) => {
  const { [idKey]: id, ['*']: path, ...params } = useParams();
  const navigate = useNavigate();
  const location = useLocation();
  const { state } = location;

  const isAdd = useMemo(() => path === ID_NEW, [path]);
  const isClone = useMemo(() => path?.includes('clone'), [path]);
  const isEdit = useMemo(
    () => path && !isAdd && !isClone,
    [path, isAdd, isClone]
  );

  const onOpenAddModal = useCallback(
    () => navigate(addPath || `${location.pathname}/${ID_NEW}`),
    [addPath, location]
  );

  const onOpenEditModal = useCallback(
    (path: string, state?: AnyObject) =>
      navigate(path, {
        state
      }),
    []
  );

  const onCloseActionModal = useCallback(
    (path?: string) => {
      if (closeModalPath || path) navigate(String(closeModalPath || path));
    },
    [closeModalPath]
  );

  return {
    id,
    isAdd,
    isEdit,
    isClone,
    state,
    onOpenAddModal,
    onOpenEditModal,
    onCloseActionModal,
    params
  };
};

export default useActionWithURL;
