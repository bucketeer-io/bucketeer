import { useState } from 'react';
import { useToggleOpen } from 'hooks';
import PageLayout from 'elements/page-layout';
import { EmptyCollection } from './collection-layout/empty-collection';
import AddGoalModal from './goals-modal/add-goal-modal';
import ConnectionsModal from './goals-modal/connections-modal';
import PageContent from './page-content';
import { GoalActions } from './types';

export type Goal = {
  id: string;
  name: string;
  description: string;
  isInUseStatus: boolean;
  connections: {
    type: string;
    data: {
      id: string;
      name: string;
    }[];
  } | null;
  updatedAt: string;
  createdAt: string;
};

export const mocks: Goal[] = [
  {
    id: 'GOAL_ID_1',
    name: 'Goal 1',
    description: 'Goal 1 Description',
    isInUseStatus: true,
    connections: {
      type: 'experiments',
      data: [
        {
          id: 'experiment_1',
          name: 'Experiment 1'
        }
      ]
    },
    updatedAt: '0',
    createdAt: '0'
  },
  {
    id: 'GOAL_ID_2',
    name: 'Goal 2',
    description: 'Goal 2 Description',
    isInUseStatus: true,
    connections: {
      type: 'operations',
      data: [
        {
          id: 'operation_1',
          name: 'Operation 1'
        },
        {
          id: 'operation_2',
          name: 'Operation 2'
        }
      ]
    },
    updatedAt: '0',
    createdAt: '0'
  },
  {
    id: 'GOAL_ID_3',
    name: 'Goal 3',
    description: 'Goal 3 Description',
    isInUseStatus: false,
    connections: null,
    updatedAt: '0',
    createdAt: '0'
  }
];

export const collection = {
  goals: mocks,
  totalCount: mocks.length
};
const PageLoader = () => {
  // const { t } = useTranslation(['table']);
  // const queryClient = useQueryClient();
  // const { consoleAccount } = useAuth();
  // const currenEnvironment = getCurrentEnvironment(consoleAccount!);

  const [selectedGoal, setSelectedGoal] = useState<Goal>();

  const [isOpenAddModal, onOpenAddModal, onCloseAddModal] =
    useToggleOpen(false);

  const [isOpenConnectionModal, onOpenConnectionModal, onCloseConnectionModal] =
    useToggleOpen(false);

  const onHandleActions = (goal: Goal, type: GoalActions) => {
    setSelectedGoal(goal);
    if (type === 'CONNECTION') {
      return onOpenConnectionModal();
    }
  };
  // const mutationState = useMutation({
  //   mutationFn: async (id: string) => {
  //     return apiKeyUpdater({
  //       id,
  //       environmentId: currenEnvironment.id,
  //       disabled: isDisabling
  //     });
  //   },
  //   onSuccess: () => {
  //     onCloseConfirmModal();
  //     invalidateAPIKeys(queryClient);
  //     mutationState.reset();
  //   }
  // });

  // const onHandleDisable = () => {
  //   if (selectedAPIKey?.id) {
  //     mutationState.mutate(selectedAPIKey.id);
  //   }
  // };
  const isLoading = false;
  const isError = false;
  const isEmpty = collection?.goals.length === 0;

  return (
    <>
      {isLoading ? (
        <PageLayout.LoadingState />
      ) : isError ? (
        <PageLayout.ErrorState onRetry={() => {}} />
      ) : isEmpty ? (
        <PageLayout.EmptyState>
          <EmptyCollection onAdd={onOpenAddModal} />
        </PageLayout.EmptyState>
      ) : (
        <PageContent onAdd={onOpenAddModal} onHandleActions={onHandleActions} />
      )}

      {isOpenAddModal && (
        <AddGoalModal isOpen={isOpenAddModal} onClose={onCloseAddModal} />
      )}
      {isOpenConnectionModal && selectedGoal && (
        <ConnectionsModal
          isOpen={isOpenConnectionModal}
          goal={selectedGoal}
          onClose={onCloseConnectionModal}
        />
      )}
    </>
  );
};

export default PageLoader;
