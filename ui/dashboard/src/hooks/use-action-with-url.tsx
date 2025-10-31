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
  const lastPath = useMemo(() => {
    const paths = location.pathname.split('/').filter(Boolean);
    return paths[paths.length - 1];
  }, [location.pathname]);

  const isAdd = useMemo(() => lastPath === ID_NEW, [lastPath]);
  const isClone = useMemo(() => path?.includes('clone'), [path]);
  const isEdit = useMemo(
    () => id && path && !isAdd && !isClone,
    [path, isAdd, isClone, idKey]
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
