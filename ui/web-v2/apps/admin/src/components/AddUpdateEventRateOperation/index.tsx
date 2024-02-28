import { messages } from '@/lang/messages';
import { AppState } from '@/modules';
import { useCurrentEnvironment, useIsEditable } from '@/modules/me';
import { addFormSchema } from '@/pages/goal/formSchema';
import { OpsEventRateClause } from '@/proto/autoops/clause_pb';
import { Goal } from '@/proto/experiment/goal_pb';
import { ListGoalsRequest } from '@/proto/experiment/service_pb';
import { AppDispatch } from '@/store';
import { classNames } from '@/utils/css';
import { Dialog, Transition } from '@headlessui/react';
import { ExclamationCircleIcon, XIcon } from '@heroicons/react/outline';
import { yupResolver } from '@hookform/resolvers/yup';
import React, {
  FC,
  Fragment,
  memo,
  useCallback,
  useEffect,
  useState,
} from 'react';
import { Controller, useForm, useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import ReactSelect, { components } from 'react-select';

import {
  createGoal,
  listGoals,
  selectAll as selectAllGoals,
} from '../../modules/goals';
import { operatorOptions } from '../FeatureAutoOpsRulesForm';
import { Select, Option } from '../Select';

interface AddUpdateEventRateOperationProps {
  isSeeDetailsSelected: boolean;
  variationOptions: Option[];
  featureId: string;
}

export const AddUpdateEventRateOperation: FC<AddUpdateEventRateOperationProps> =
  memo(({ isSeeDetailsSelected, variationOptions, featureId }) => {
    const { formatMessage: f } = useIntl();
    const editable = useIsEditable();
    const dispatch = useDispatch<AppDispatch>();
    const currentEnvironment = useCurrentEnvironment();

    const [isAddGoalOpen, setIsAddGoalOpen] = useState(false);

    const methods = useFormContext<any>();
    const {
      control,
      formState: { errors },
      register,
    } = methods;

    const goals = useSelector<AppState, Goal.AsObject[]>(
      (state) => selectAllGoals(state.goals),
      shallowEqual
    );

    useEffect(() => {
      dispatch(
        listGoals({
          environmentNamespace: currentEnvironment.id,
          pageSize: 99999,
          cursor: '',
          searchKeyword: '',
          status: null,
          orderBy: ListGoalsRequest.OrderBy.NAME,
          orderDirection: ListGoalsRequest.OrderDirection.ASC,
        })
      );
    }, [dispatch, featureId, currentEnvironment]);

    const goalOptions = goals.map((goal) => {
      return {
        value: goal.id,
        label: goal.id,
      };
    });

    return (
      <div className="mt-4 space-y-2 pl-1">
        <AddGoalModal open={isAddGoalOpen} setOpen={setIsAddGoalOpen} />
        <div className="px-4 py-1 text-pink-500 bg-pink-50 inline-block">
          If
        </div>
        <div>
          <span className="input-label">{f(messages.feature.variation)}</span>
          <Controller
            name="eventRate.variation"
            control={control}
            render={({ field }) => (
              <Select
                onChange={(o: Option) => {
                  field.onChange(o.value);
                }}
                options={variationOptions}
                disabled={!editable || isSeeDetailsSelected}
                value={variationOptions.find((o) => o.value === field.value)}
              />
            )}
          />
        </div>
        <div>
          <span className="input-label">
            {f(messages.autoOps.opsEventRateClause.goal)}
          </span>
          <Controller
            name="eventRate.goal"
            control={control}
            render={({ field }) => (
              <AddGoalSelect
                onChange={(o: Option) => field.onChange(o.value)}
                options={goalOptions}
                disabled={!editable || isSeeDetailsSelected}
                value={goalOptions.find((o) => o.value === field.value)}
                openAddGoalModal={() => setIsAddGoalOpen(true)}
              />
            )}
          />
        </div>
        <div className="grid grid-cols-3 gap-3">
          <div>
            <span className="input-label">{f(messages.autoOps.condition)}</span>
            <Controller
              name="eventRate.operator"
              control={control}
              render={({ field }) => (
                <Select
                  onChange={(o: Option) => field.onChange(o.value)}
                  options={operatorOptions}
                  disabled={!editable || isSeeDetailsSelected}
                  value={operatorOptions.find((o) => o.value === field.value)}
                />
              )}
            />
          </div>
          <div>
            <span className="input-label">{f(messages.autoOps.threshold)}</span>
            <div className="flex">
              <input
                {...register('eventRate.threadsholdRate')}
                type="number"
                min="0"
                max="100"
                className={classNames(
                  'w-full',
                  errors.eventRate?.threadsholdRate
                    ? 'input-text-error'
                    : 'input-text'
                )}
                placeholder={''}
                required
                disabled={!editable || isSeeDetailsSelected}
              />
              <span
                className={classNames(
                  'px-1 py-1 inline-flex items-center bg-gray-100',
                  'rounded-r border border-l-0 border-gray-300 text-gray-600'
                )}
              >
                {'%'}
              </span>
            </div>
          </div>
          <div>
            <span className="input-label">
              {f(messages.autoOps.opsEventRateClause.minCount)}
            </span>
            <div>
              <input
                {...register('eventRate.minCount')}
                type="number"
                min="0"
                className={classNames(
                  'w-full',
                  errors.eventRate?.minCount ? 'input-text-error' : 'input-text'
                )}
                disabled={!editable || isSeeDetailsSelected}
              />
            </div>
          </div>
        </div>
        <div>
          {errors.eventRate?.threadsholdRate?.message && (
            <p className="input-error">
              <span role="alert">
                {errors.eventRate?.threadsholdRate?.message}
              </span>
            </p>
          )}
          {errors.eventRate?.minCount?.message && (
            <p className="input-error">
              <span role="alert">{errors.eventRate?.minCount?.message}</span>
            </p>
          )}
        </div>
      </div>
    );
  });

export interface AddGoalSelectProps {
  options: Option[];
  disabled?: boolean;
  clearable?: boolean;
  isLoading?: boolean;
  isMulti?: boolean;
  isSearchable?: boolean;
  value?: Option;
  className?: string;
  onChange: ((option: Option) => void) | ((option: Option[]) => void);
  placeholder?: string;
  openAddGoalModal: () => void;
}

const AddGoalSelect: FC<AddGoalSelectProps> = memo(
  ({
    disabled,
    className,
    clearable,
    isLoading,
    isMulti,
    isSearchable,
    onChange,
    options,
    value,
    placeholder,
    openAddGoalModal,
  }) => {
    const textColor = '#3F3F46';
    const textColorDisabled = '#6B7280';
    const backgroundColor = 'white';
    const backgroundColorDisabled = '#F3F4F6';
    const borderColor = '#D1D5DB';
    const fontSize = '0.875rem';
    const lineHeight = '1.25rem';
    const minHeight = '2.5rem';
    const colourStyles = {
      control: (styles, { isDisabled }) => ({
        ...styles,
        backgroundColor: isDisabled ? backgroundColorDisabled : backgroundColor,
        borderColor: borderColor,
        '&:hover': {
          borderColor: borderColor,
        },
        fontSize: fontSize,
        lineHeight: lineHeight,
        minHeight: minHeight,
        '*': {
          boxShadow: 'none !important',
        },
      }),
      option: (styles, { isFocused, isSelected }) => {
        return {
          ...styles,
          backgroundColor: isFocused
            ? backgroundColorDisabled
            : isSelected
            ? backgroundColor
            : null,
          overflow: 'hidden',
          textOverflow: 'ellipsis',
          whiteSpace: 'nowrap',
          color: textColor,
          ':active': {
            backgroundColor: backgroundColor,
          },
        };
      },
      menu: (base) => ({
        ...base,
        fontSize: fontSize,
        lineHeight: lineHeight,
        color: textColor,
      }),
      multiValueLabel: (base, { isDisabled }) => ({
        ...base,
        color: isDisabled ? textColorDisabled : textColor,
      }),
      singleValue: (base, { isDisabled }) => ({
        ...base,
        color: isDisabled ? textColorDisabled : textColor,
      }),
    };
    return (
      <ReactSelect
        options={options}
        className={className}
        classNamePrefix="add-goal-react-select"
        styles={colourStyles}
        components={{
          Option: (props) => (
            <CustomOption {...props} openAddGoalModal={openAddGoalModal} />
          ),
        }}
        isDisabled={isLoading || disabled}
        isClearable={clearable}
        isMulti={isMulti}
        isSearchable={isSearchable}
        isLoading={isLoading}
        placeholder={placeholder ? placeholder : ''}
        value={value}
        onChange={onChange}
        openAddGoalModal={openAddGoalModal}
      />
    );
  }
);

const CustomOption = ({ children, ...props }) => {
  const isLastOption =
    props.options[props.options.length - 1]?.value === props.data.value;

  if (isLastOption) {
    return (
      <div>
        <div
          {...props.innderRef}
          {...props.innerProps}
          className="px-3 py-2 hover:bg-[#F3F4F6]"
        >
          {children}
        </div>
        <div
          onClick={props.openAddGoalModal}
          className="text-center py-[10px] cursor-pointer border-t hover:bg-[#F3F4F6]"
        >
          Add New Goal
        </div>
      </div>
    );
  }

  return <components.Option {...props}>{children}</components.Option>;
};

interface AddGoalModalProps {
  open: boolean;
  setOpen: React.Dispatch<React.SetStateAction<boolean>>;
}

const AddGoalModal: FC<AddGoalModalProps> = ({ open, setOpen }) => {
  const { formatMessage: f } = useIntl();
  const dispatch = useDispatch<AppDispatch>();
  const currentEnvironment = useCurrentEnvironment();
  const methods = useFormContext();

  const { setValue } = methods;

  const {
    register,
    handleSubmit,
    formState: { errors, isValid, isSubmitting },
    reset: resetAdd,
  } = useForm({
    resolver: yupResolver(addFormSchema),
    defaultValues: {
      id: '',
      name: '',
      description: '',
    },
    mode: 'onChange',
  });

  const handleCreateGoal = useCallback(
    async (data) => {
      dispatch(
        createGoal({
          environmentNamespace: currentEnvironment.id,
          id: data.id,
          name: data.name,
          description: data.description,
        })
      ).then(() => {
        setOpen(false);
        resetAdd();
        dispatch(
          listGoals({
            environmentNamespace: currentEnvironment.id,
            pageSize: 99999,
            cursor: '',
            searchKeyword: '',
            status: null,
            orderBy: ListGoalsRequest.OrderBy.NAME,
            orderDirection: ListGoalsRequest.OrderDirection.ASC,
          })
        );
        setValue('eventRate.goal', data.id);
      });
    },
    [dispatch]
  );

  return (
    <Transition.Root show={open} as={Fragment}>
      <Dialog as="div" className="relative z-50" onClose={setOpen}>
        <Transition.Child
          as={Fragment}
          enter="ease-out duration-300"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in duration-200"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
        </Transition.Child>
        <form className="fixed inset-0 z-10 overflow-y-auto">
          <div className="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
            <Transition.Child
              as={Fragment}
              enter="ease-out duration-300"
              enterFrom="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
              enterTo="opacity-100 translate-y-0 sm:scale-100"
              leave="ease-in duration-200"
              leaveFrom="opacity-100 translate-y-0 sm:scale-100"
              leaveTo="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
            >
              <div className="relative transform overflow-hidden rounded-lg bg-white text-left shadow-xl transition-all w-[542px]">
                <div className="flex items-center justify-between px-4 py-5 border-b">
                  <p className="text-xl font-medium">Create Goal</p>
                  <XIcon
                    width={20}
                    className="text-gray-400 cursor-pointer"
                    onClick={() => setOpen(false)}
                  />
                </div>
                <div className="p-4 space-y-4">
                  <p className="font-bold">General Information</p>
                  <div className="space-y-1">
                    <label
                      htmlFor="name"
                      className="flex space-x-2 items-center"
                    >
                      <span className="input-label">{f({ id: 'name' })}</span>
                      <ExclamationCircleIcon width={18} />
                    </label>
                    <input
                      {...register('name')}
                      type="text"
                      name="name"
                      id="name"
                      className="input-text w-full"
                    />
                    <p className="input-error">
                      {errors.name && (
                        <span role="alert">{errors.name.message}</span>
                      )}
                    </p>
                  </div>
                  <div className="space-y-1">
                    <label htmlFor="id" className="flex space-x-2 items-center">
                      <span className="input-label">{f({ id: 'id' })}</span>
                      <ExclamationCircleIcon width={18} />
                    </label>
                    <input
                      {...register('id')}
                      type="text"
                      name="id"
                      id="id"
                      className="input-text w-full"
                    />
                    <p className="input-error">
                      {errors.id && (
                        <span role="alert">{errors.id.message}</span>
                      )}
                    </p>
                  </div>
                  <div className="space-y-1">
                    <label htmlFor="description" className="block">
                      <span className="input-label">
                        {f(messages.description)}
                      </span>
                      <span className="input-label-optional">
                        {' '}
                        ({f(messages.input.optional)})
                      </span>
                    </label>
                    <textarea
                      {...register('description')}
                      id="description"
                      name="description"
                      rows={5}
                      className="input-text w-full break-all"
                    />
                    <p className="input-error">
                      {errors.description && (
                        <span role="alert">{errors.description.message}</span>
                      )}
                    </p>
                  </div>
                </div>
                <div className="p-4 flex justify-end border-t space-x-4">
                  <button
                    type="button"
                    className="btn-cancel"
                    disabled={false}
                    onClick={() => setOpen(false)}
                  >
                    {f(messages.button.cancel)}
                  </button>
                  <button
                    type="button"
                    className="btn-submit"
                    disabled={!isValid || isSubmitting}
                    onClick={handleSubmit(handleCreateGoal)}
                  >
                    New Goal
                  </button>
                </div>
              </div>
            </Transition.Child>
          </div>
        </form>
      </Dialog>
    </Transition.Root>
  );
};

interface OpsEventRateClauseSchema {
  variation: string;
  goal: string;
  minCount: number;
  threadsholdRate: number;
  operator: string;
}

export function createOpsEventRateClause(
  oerc: OpsEventRateClauseSchema
): OpsEventRateClause {
  const clause = new OpsEventRateClause();
  clause.setVariationId(oerc.variation);
  clause.setGoalId(oerc.goal);
  clause.setMinCount(oerc.minCount);
  clause.setThreadsholdRate(oerc.threadsholdRate / 100);
  clause.setOperator(createOpsEventRateOperator(oerc.operator));
  return clause;
}

export function createOpsEventRateOperator(
  value: string
): OpsEventRateClause.OperatorMap[keyof OpsEventRateClause.OperatorMap] {
  if (value === OpsEventRateClause.Operator.GREATER_OR_EQUAL.toString()) {
    return OpsEventRateClause.Operator.GREATER_OR_EQUAL;
  }
  return OpsEventRateClause.Operator.LESS_OR_EQUAL;
}
