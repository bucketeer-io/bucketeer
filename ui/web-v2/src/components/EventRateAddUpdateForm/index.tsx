import { AppState } from '../../modules';
import { useCurrentEnvironment, useIsEditable } from '../../modules/me';
import { AutoOpsRule, OpsType } from '../../proto/autoops/auto_ops_rule_pb';
import {
  createGoal,
  listGoals,
  selectAll as selectAllGoals
} from '../../modules/goals';
import { AppDispatch } from '../../store';
import { PlusIcon, XIcon } from '@heroicons/react/outline';
import React, {
  FC,
  Fragment,
  memo,
  useCallback,
  useEffect,
  useState
} from 'react';
import { Controller, useForm, useFormContext } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';

import { messages } from '../../lang/messages';
import { OperationForm } from '../../pages/feature/formSchema';
import { Option, Select } from '../Select';
import { Goal } from '../../proto/experiment/goal_pb';
import { ListGoalsRequest } from '../../proto/experiment/service_pb';
import { classNames } from '../../utils/css';
import { yupResolver } from '@hookform/resolvers/yup';
import { addFormSchema } from '../../pages/goal/formSchema';
import { Dialog, Transition } from '@headlessui/react';
import { OpsEventRateClause, ActionType } from '../../proto/autoops/clause_pb';
import ReactSelect, { components } from 'react-select';
import {
  CreateAutoOpsRuleCommand,
  ChangeOpsEventRateClauseCommand
} from '../../proto/autoops/command_pb';
import {
  createAutoOpsRule,
  updateAutoOpsRule,
  UpdateAutoOpsRuleParams
} from '../../modules/autoOpsRules';
import { operatorOptions } from '../../pages/feature/autoops';

export interface EventRateAddUpdateFormProps {
  onCancel: () => void;
  featureId: string;
  autoOpsRule?: AutoOpsRule.AsObject;
  isActiveTabSelected: boolean;
  variationOptions: Option[];
  onSubmit: () => void;
}

export const EventRateAddUpdateForm: FC<EventRateAddUpdateFormProps> = memo(
  ({
    onCancel,
    featureId,
    autoOpsRule,
    isActiveTabSelected,
    variationOptions,
    onSubmit
  }) => {
    const dispatch = useDispatch<AppDispatch>();
    const currentEnvironment = useCurrentEnvironment();

    const { formatMessage: f } = useIntl();
    const editable = useIsEditable();

    const methods = useFormContext<OperationForm>();
    const {
      handleSubmit,
      control,
      formState: { isValid, isSubmitting, errors },
      register,
      setValue
    } = methods;

    const isSeeDetailsSelected = autoOpsRule && !isActiveTabSelected;

    const title = () => {
      if (isSeeDetailsSelected) {
        return f(messages.autoOps.operationDetails);
      } else {
        return autoOpsRule
          ? f(messages.autoOps.updateOperation)
          : f(messages.autoOps.createOperation);
      }
    };

    const [isAddGoalOpen, setIsAddGoalOpen] = useState(false);

    const goals = useSelector<AppState, Goal.AsObject[]>(
      (state) => selectAllGoals(state.goals),
      shallowEqual
    );

    useEffect(() => {
      dispatch(
        listGoals({
          environmentId: currentEnvironment.id,
          pageSize: 0,
          cursor: '',
          searchKeyword: '',
          status: null,
          orderBy: ListGoalsRequest.OrderBy.NAME,
          orderDirection: ListGoalsRequest.OrderDirection.ASC
        })
      );
    }, [dispatch, featureId, currentEnvironment]);

    useEffect(() => {
      if (autoOpsRule) {
        const opsEventRateClause = OpsEventRateClause.deserializeBinary(
          autoOpsRule.clausesList[0].clause.value as Uint8Array
        ).toObject();

        setValue('eventRate.variation', opsEventRateClause.variationId);
        setValue('eventRate.goal', opsEventRateClause.goalId);
        setValue('eventRate.operator', opsEventRateClause.operator.toString());
        setValue(
          'eventRate.threadsholdRate',
          Math.round(opsEventRateClause.threadsholdRate * 100)
        );
        setValue('eventRate.minCount', opsEventRateClause.minCount);
      }
    }, [autoOpsRule]);

    const handleOnSubmit = useCallback(
      (data) => {
        if (autoOpsRule) {
          const changeOpsEventRateClauseCommands: ChangeOpsEventRateClauseCommand[] =
            [];

          const command = new ChangeOpsEventRateClauseCommand();
          command.setId(autoOpsRule.clausesList[0].id);
          command.setOpsEventRateClause(
            createOpsEventRateClause(data.eventRate)
          );
          changeOpsEventRateClauseCommands.push(command);

          const param: UpdateAutoOpsRuleParams = {
            environmentId: currentEnvironment.id,
            id: autoOpsRule.id,
            changeOpsEventRateClauseCommands
          };

          dispatch(updateAutoOpsRule(param)).then(() => onSubmit());
        } else {
          const command = new CreateAutoOpsRuleCommand();
          command.setFeatureId(featureId);
          command.setOpsType(OpsType.EVENT_RATE);
          command.setOpsEventRateClausesList([
            createOpsEventRateClause(data.eventRate)
          ]);

          dispatch(
            createAutoOpsRule({
              environmentId: currentEnvironment.id,
              command: command
            })
          ).then(() => onSubmit());
        }
      },
      [autoOpsRule]
    );

    const goalOptions = goals.map((goal) => {
      return {
        value: goal.id,
        label: goal.id
      };
    });

    return (
      <div className="w-[530px] h-full overflow-hidden">
        <form className="flex flex-col h-full overflow-hidden">
          <div className="h-full flex flex-col overflow-hidden">
            <div className="flex items-center justify-between px-4 py-5 border-b">
              <p className="text-xl font-medium">{title()}</p>
              <XIcon
                width={20}
                className="text-gray-400 cursor-pointer"
                onClick={onCancel}
              />
            </div>
            <div className="px-4 h-full flex flex-col overflow-hidden">
              <div className="flex py-5">
                <p className="font-bold">{f(messages.autoOps.eventRate)}</p>
              </div>
              <div className="h-full flex flex-col overflow-hidden space-y-4 px-1">
                <div className="flex divide-x divide-primary space-x-3">
                  <AddGoalModal
                    open={isAddGoalOpen}
                    setOpen={setIsAddGoalOpen}
                  />
                  <div className="px-4 py-1 text-pink-500 bg-pink-50 inline-block self-center">
                    If
                  </div>
                  <div className="space-y-3 pl-3">
                    <div>
                      <span className="input-label">
                        {f(messages.feature.variation)}
                      </span>
                      <Controller
                        name="eventRate.variation"
                        control={control}
                        render={({ field }) => (
                          <Select
                            isSearchable={false}
                            onChange={(o: Option) => {
                              field.onChange(o.value);
                            }}
                            options={variationOptions}
                            disabled={!editable || isSeeDetailsSelected}
                            value={variationOptions.find(
                              (o) => o.value === field.value
                            )}
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
                            value={goalOptions.find(
                              (o) => o.value === field.value
                            )}
                            openAddGoalModal={() => setIsAddGoalOpen(true)}
                          />
                        )}
                      />
                    </div>
                    <div className="grid grid-cols-3 gap-3">
                      <div>
                        <span className="input-label">
                          {f(messages.autoOps.condition)}
                        </span>
                        <Controller
                          name="eventRate.operator"
                          control={control}
                          render={({ field }) => (
                            <Select
                              onChange={(o: Option) => field.onChange(o.value)}
                              options={operatorOptions}
                              disabled={!editable || isSeeDetailsSelected}
                              value={operatorOptions.find(
                                (o) => o.value === field.value
                              )}
                            />
                          )}
                        />
                      </div>
                      <div>
                        <span className="input-label">
                          {f(messages.autoOps.threshold)}
                        </span>
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
                              errors.eventRate?.minCount
                                ? 'input-text-error'
                                : 'input-text'
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
                          <span role="alert">
                            {errors.eventRate?.minCount?.message}
                          </span>
                        </p>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div className="flex-shrink-0 px-4 py-4 flex justify-end border-t">
            <div className="mr-3">
              <button
                type="button"
                className="btn-cancel"
                disabled={false}
                onClick={onCancel}
              >
                {f(messages.button.cancel)}
              </button>
            </div>
            <button
              type="button"
              className="btn-submit-gradient"
              disabled={!isValid || isSubmitting || isSeeDetailsSelected}
              onClick={handleSubmit(handleOnSubmit)}
            >
              {f(messages.button.submit)}
            </button>
          </div>
        </form>
      </div>
    );
  }
);

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
    openAddGoalModal
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
          borderColor: borderColor
        },
        fontSize: fontSize,
        lineHeight: lineHeight,
        minHeight: minHeight,
        '*': {
          boxShadow: 'none !important'
        }
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
            backgroundColor: backgroundColor
          }
        };
      },
      menu: (base) => ({
        ...base,
        fontSize: fontSize,
        lineHeight: lineHeight,
        color: textColor
      }),
      multiValueLabel: (base, { isDisabled }) => ({
        ...base,
        color: isDisabled ? textColorDisabled : textColor
      }),
      singleValue: (base, { isDisabled }) => ({
        ...base,
        color: isDisabled ? textColorDisabled : textColor
      })
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
          NoOptionsMessage: () => (
            <AddNewGoalOption openAddGoalModal={openAddGoalModal} />
          )
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
        <AddNewGoalOption openAddGoalModal={props.openAddGoalModal} />
      </div>
    );
  }

  return <components.Option {...props}>{children}</components.Option>;
};

const AddNewGoalOption = ({ openAddGoalModal }) => {
  const { formatMessage: f } = useIntl();
  return (
    <div
      onClick={openAddGoalModal}
      className="py-[10px] space-x-2 cursor-pointer border-t hover:bg-[#F3F4F6] text-primary flex items-center justify-center"
    >
      <PlusIcon width={16} />
      <span>{f(messages.goal.addNewGoal)}</span>
    </div>
  );
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

  const { setValue, trigger } = methods;

  const {
    register,
    handleSubmit,
    formState: { errors, isValid, isSubmitting },
    reset: resetAdd
  } = useForm({
    resolver: yupResolver(addFormSchema),
    defaultValues: {
      id: '',
      name: '',
      description: ''
    },
    mode: 'onChange'
  });

  const handleCreateGoal = useCallback(
    async (data) => {
      dispatch(
        createGoal({
          environmentId: currentEnvironment.id,
          id: data.id,
          name: data.name,
          description: data.description
        })
      ).then(() => {
        setOpen(false);
        resetAdd();
        dispatch(
          listGoals({
            environmentId: currentEnvironment.id,
            pageSize: 0,
            cursor: '',
            searchKeyword: '',
            status: null,
            orderBy: ListGoalsRequest.OrderBy.NAME,
            orderDirection: ListGoalsRequest.OrderDirection.ASC
          })
        );
        setValue('eventRate.goal', data.id);
        trigger('eventRate.goal');
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
                  <p className="text-xl font-medium">
                    {f(messages.goal.newGoal)}
                  </p>
                  <XIcon
                    width={20}
                    className="text-gray-400 cursor-pointer"
                    onClick={() => setOpen(false)}
                  />
                </div>
                <div className="p-4 space-y-4">
                  <p className="font-bold">{f(messages.generalInformation)}</p>
                  <div className="space-y-1">
                    <label htmlFor="name" className="flex items-center">
                      <span className="input-label">{f({ id: 'name' })}</span>
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
                    <label htmlFor="id" className="flex items-center">
                      <span className="input-label">
                        {f(messages.goal.goalId)}
                      </span>
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
                        {f(messages.input.optional)}
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
                    {f(messages.goal.createGoal)}
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
  clause.setActionType(ActionType.DISABLE);
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
