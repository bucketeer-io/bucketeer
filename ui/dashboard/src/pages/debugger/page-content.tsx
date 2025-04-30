import { useCallback, useMemo, useState } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { debuggerEvaluate } from '@api/debugger';
import { yupResolver } from '@hookform/resolvers/yup';
import { useQueryFeatures } from '@queries/features';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useToast } from 'hooks';
import { Evaluation } from '@types';
import { ExpandOrCollapse } from 'pages/audit-logs/types';
import Form from 'components/form';
import PageLayout from 'elements/page-layout';
import AddDebuggerForm from './add-debugger-form';
import DebuggerResults from './debugger-results';
import { addDebuggerFormSchema, AddDebuggerFormType } from './form-schema';
import { EvaluationFeature } from './types';

export type GroupByType = 'FLAG' | 'USER';

const PageContent = () => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  const [isShowResults, setIsShowResults] = useState(false);
  const [evaluations, setEvaluations] = useState<Evaluation[]>([]);
  const [groupBy, setGroupBy] = useState<GroupByType>('FLAG');
  const [expandedItems, setExpandedItems] = useState<number[]>([]);
  const [expandOrCollapseAllState, setExpandOrCollapseAllState] =
    useState<ExpandOrCollapse>(ExpandOrCollapse.COLLAPSE);

  const { errorNotify } = useToast();

  const { data: featureCollection } = useQueryFeatures({
    params: {
      cursor: String(0),
      environmentId: currentEnvironment.id
    }
  });

  const features = featureCollection?.features || [];

  const defaultValues = useMemo(
    () => ({
      flags: [''],
      userIds: [],
      attributes: [
        {
          key: '',
          value: ''
        }
      ]
    }),
    []
  );

  const form = useForm({
    resolver: yupResolver(addDebuggerFormSchema),
    defaultValues: {
      ...defaultValues
    },
    mode: 'onChange'
  });

  const getGroupByEvaluateFeatures = useCallback(
    ({
      evaluationsVariable,
      groupByVariable
    }: {
      evaluationsVariable: Evaluation[];
      groupByVariable: GroupByType;
    }) => {
      const isFlag = groupByVariable === 'FLAG';
      const data = new Map();

      evaluationsVariable.forEach(item => {
        const groupByField = isFlag ? item.featureId : item.userId;
        data.set(groupByField, [
          ...(data.get(groupByField) || []),
          {
            ...item,
            feature: features.find(feature => feature.id === item.featureId)
          }
        ]);
      });
      const results: EvaluationFeature[][] = [];
      data.forEach(evaluations => results.push(evaluations));
      return results;
    },
    [groupBy, features, evaluations]
  );

  const groupByEvaluateFeatures = useMemo(
    () =>
      getGroupByEvaluateFeatures({
        evaluationsVariable: evaluations,
        groupByVariable: groupBy
      }),
    [evaluations, groupBy]
  );

  const onSubmit = useCallback(
    async (values: AddDebuggerFormType) => {
      try {
        const dataMap = new Map();
        values?.attributes?.forEach(item => dataMap.set(item.key, item.value));

        const userData: { [key: string]: string } = {};
        dataMap?.forEach((value, key) => (userData[key] = value));

        const resp = await debuggerEvaluate({
          environmentId: currentEnvironment.id,
          featureIds: values.flags,
          users: values.userIds.map(item => ({
            id: item,
            data: userData
          }))
        });
        setEvaluations(resp.evaluations);
        onToggleExpandAll({
          evaluationsVariable: resp.evaluations
        });
        setIsShowResults(true);
      } catch (error) {
        errorNotify(error);
      }
    },
    [currentEnvironment, groupBy, features]
  );

  const onChangeGroupBy = useCallback(
    (value: GroupByType) => {
      if (value !== groupBy) {
        const evaluationData = getGroupByEvaluateFeatures({
          evaluationsVariable: evaluations,
          groupByVariable: value
        });
        setExpandOrCollapseAllState(ExpandOrCollapse.EXPAND);
        setExpandedItems(
          Array.from({ length: evaluationData.length }, (_, i) => i)
        );
      }

      setGroupBy(value);
    },
    [evaluations, groupBy, expandOrCollapseAllState]
  );

  const onToggleExpandItem = useCallback(
    (index: number) => {
      const isExistedItem = expandedItems.includes(index);
      const newExpandedItems = isExistedItem
        ? expandedItems.filter(item => item !== index)
        : [...expandedItems, index];

      const expandState =
        newExpandedItems.length < groupByEvaluateFeatures.length
          ? ExpandOrCollapse.COLLAPSE
          : ExpandOrCollapse.EXPAND;

      setExpandOrCollapseAllState(expandState);
      setExpandedItems(newExpandedItems);
    },
    [expandedItems, groupByEvaluateFeatures]
  );

  const onToggleExpandAll = useCallback(
    ({ evaluationsVariable }: { evaluationsVariable?: Evaluation[] }) => {
      if (expandOrCollapseAllState === ExpandOrCollapse.EXPAND) {
        setExpandedItems([]);
        return setExpandOrCollapseAllState(ExpandOrCollapse.COLLAPSE);
      }
      const evaluationData = getGroupByEvaluateFeatures({
        evaluationsVariable: evaluationsVariable || evaluations,
        groupByVariable: groupBy
      });

      setExpandedItems(
        Array.from({ length: evaluationData.length }, (_, i) => i)
      );
      setExpandOrCollapseAllState(ExpandOrCollapse.EXPAND);
    },
    [evaluations, groupBy, features, expandOrCollapseAllState]
  );

  return (
    <PageLayout.Content className={!isShowResults ? 'pt-0' : ''}>
      {!isShowResults ? (
        <FormProvider {...form}>
          <Form onSubmit={form.handleSubmit(onSubmit)}>
            <AddDebuggerForm
              isLoading={form.formState.isSubmitting}
              evaluations={evaluations}
              onCancel={() => setIsShowResults(true)}
            />
          </Form>
        </FormProvider>
      ) : (
        <DebuggerResults
          groupBy={groupBy}
          isExpandAll={expandOrCollapseAllState === ExpandOrCollapse.EXPAND}
          expandedItems={expandedItems}
          groupByEvaluateFeatures={groupByEvaluateFeatures}
          onEditFields={() => setIsShowResults(false)}
          onResetFields={() => {
            setIsShowResults(false);
            form.reset({
              ...defaultValues
            });
            setEvaluations([]);
            setExpandedItems([]);
            setGroupBy('FLAG');
            setExpandOrCollapseAllState(ExpandOrCollapse.COLLAPSE);
          }}
          onChangeGroupBy={onChangeGroupBy}
          onToggleExpandItem={onToggleExpandItem}
          onToggleExpandAll={() =>
            onToggleExpandAll({
              evaluationsVariable: evaluations
            })
          }
        />
      )}
    </PageLayout.Content>
  );
};

export default PageContent;
