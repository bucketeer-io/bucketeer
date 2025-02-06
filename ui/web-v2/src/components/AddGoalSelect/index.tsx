import { FC, Fragment, memo, useCallback, useState } from 'react';
import ReactSelect, { components } from 'react-select';
import { Option } from '../Select';
import { useIntl } from 'react-intl';
import { PlusIcon, XIcon } from '@heroicons/react/outline';
import { messages } from '../../lang/messages';
import { useDispatch } from 'react-redux';
import { AppDispatch } from '../../store';
import { useCurrentEnvironment } from '../../modules/me';
import { useForm, useFormContext } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { addGoalModalFormSchema } from '../../pages/goal/formSchema';
import { createGoal, listGoals } from '../../modules/goals';
import { Goal } from '../../proto/experiment/goal_pb';
import { ListGoalsRequest } from '../../proto/experiment/service_pb';
import { Dialog, Transition } from '@headlessui/react';

export interface AddGoalSelectProps {
  name: string;
  options: Option[];
  disabled?: boolean;
  clearable?: boolean;
  isLoading?: boolean;
  isMulti?: boolean;
  isSearchable?: boolean;
  value?: Option | Option[];
  className?: string;
  onChange: ((option: Option) => void) | ((option: Option[]) => void);
  placeholder?: string;
  connectionType?: Goal.ConnectionTypeMap[keyof Goal.ConnectionTypeMap];
}

export const AddGoalSelect: FC<AddGoalSelectProps> = memo(
  ({
    name,
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
    connectionType
  }) => {
    const textColor = '#3F3F46';
    const textColorDisabled = '#6B7280';
    const backgroundColor = 'white';
    const backgroundColorDisabled = '#F3F4F6';
    const borderColor = '#D1D5DB';
    const fontSize = '0.875rem';
    const lineHeight = '1.25rem';
    const minHeight = '2.5rem';

    const [isAddGoalOpen, setIsAddGoalOpen] = useState(false);

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

    const openAddGoalModal = () => {
      setIsAddGoalOpen(true);
    };

    return (
      <>
        <AddGoalModal
          name={name}
          open={isAddGoalOpen}
          setOpen={setIsAddGoalOpen}
          isMulti={isMulti}
          connectionType={connectionType}
        />
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
        />
      </>
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
  name: string;
  open: boolean;
  setOpen: React.Dispatch<React.SetStateAction<boolean>>;
  isMulti?: boolean;
  connectionType: Goal.ConnectionTypeMap[keyof Goal.ConnectionTypeMap];
}

const AddGoalModal: FC<AddGoalModalProps> = ({
  name,
  open,
  setOpen,
  isMulti,
  connectionType
}) => {
  const { formatMessage: f } = useIntl();
  const dispatch = useDispatch<AppDispatch>();
  const currentEnvironment = useCurrentEnvironment();

  const {
    register,
    handleSubmit,
    formState: { errors, isValid, isSubmitting },
    reset: resetAdd
  } = useForm({
    resolver: yupResolver(addGoalModalFormSchema),
    defaultValues: {
      id: '',
      name: '',
      description: ''
    },
    mode: 'onChange'
  });

  const methods = useFormContext();
  const { watch, setValue } = methods;

  const watchGoalIds = watch('goalIds') || [];

  const handleCreateGoal = useCallback(
    async (data) => {
      dispatch(
        createGoal({
          environmentId: currentEnvironment.id,
          id: data.id,
          name: data.name,
          description: data.description,
          connectionType
        })
      ).then(() => {
        setOpen(false);
        resetAdd();
        dispatch(
          listGoals({
            environmentId: currentEnvironment.id,
            pageSize: 99999,
            cursor: '',
            searchKeyword: '',
            status: null,
            orderBy: ListGoalsRequest.OrderBy.NAME,
            orderDirection: ListGoalsRequest.OrderDirection.ASC,
            connectionType
          })
        ).then(() => {
          setValue(name, isMulti ? [...watchGoalIds, data.id] : data.id, {
            shouldValidate: true
          });
        });
      });
    },
    [dispatch, watchGoalIds]
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
