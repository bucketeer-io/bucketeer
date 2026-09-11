import { Fragment, useCallback, useMemo, useState } from 'react';
import { useFieldArray, useFormContext } from 'react-hook-form';
import { Trans } from 'react-i18next';
import {
  IconArrowDownwardFilled,
  IconArrowUpwardFilled
} from 'react-icons-material-design';
import {
  DndContext,
  DragEndEvent,
  DragStartEvent,
  KeyboardSensor,
  PointerSensor,
  closestCenter,
  useSensor,
  useSensors
} from '@dnd-kit/core';
import {
  SortableContext,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy
} from '@dnd-kit/sortable';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { useQueryAttributeKeys } from '@queries/attribute-keys';
import { getCurrentEnvironment, useAuth } from 'auth';
import useOptions from 'hooks/use-options';
import { useTranslation } from 'i18n';
import compact from 'lodash/compact';
import { GripVertical } from 'lucide-react';
import { cn } from 'utils/style';
import { IconClose, IconPlus } from '@icons';
import Card from 'pages/feature-flag-details/elements/card';
import { RuleClauseType } from 'pages/feature-flag-details/targeting/types';
import { UserSegmentForm } from 'pages/user-segments/types';
import {
  createSegmentRule,
  SegmentRuleFormValue
} from 'pages/user-segments/utils';
import Button from 'components/button';
import Icon from 'components/icon';
import RuleClausesForm, { SituationOption } from 'elements/rule-clauses-form';

const RULES_FIELD_NAME = 'rules';

interface SortableRuleCardProps {
  ruleKey: string;
  ruleIndex: number;
  rulesCount: number;
  situationOptions: SituationOption[];
  usedAttributeKeys: string[];
  sdkAttributeKeys: string[];
  onChangeIndex: (type: 'increase' | 'decrease', index: number) => void;
  onRemove: (index: number) => void;
}

const SortableRuleCard = ({
  ruleKey,
  ruleIndex,
  rulesCount,
  situationOptions,
  usedAttributeKeys,
  sdkAttributeKeys,
  onChangeIndex,
  onRemove
}: SortableRuleCardProps) => {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging
  } = useSortable({ id: ruleKey });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={cn('flex flex-col w-full', {
        'opacity-50': isDragging
      })}
    >
      <Card>
        <div className="w-full h-8 flex items-center justify-between">
          <div className="flex items-center gap-x-2">
            {rulesCount > 1 && (
              <div
                className="flex-center cursor-grab active:cursor-grabbing text-gray-400 dark:text-dark-gray-200 hover:text-gray-600 dark:hover:text-dark-gray-400 touch-none"
                {...attributes}
                {...listeners}
              >
                <GripVertical size={16} />
              </div>
            )}
            <p className="typo-para-medium leading-5 text-gray-700 dark:text-dark-gray-300">
              <Trans
                i18nKey={'table:feature-flags.rule-index'}
                values={{ index: ruleIndex + 1 }}
              />
            </p>
          </div>
          <div className="flex items-center gap-x-2">
            {rulesCount > 1 && (
              <div className="flex items-center gap-x-1">
                {ruleIndex !== rulesCount - 1 && (
                  <div
                    className="flex-center group cursor-pointer"
                    onClick={() => onChangeIndex('increase', ruleIndex)}
                  >
                    <Icon
                      icon={IconArrowDownwardFilled}
                      size={'sm'}
                      className="text-gray-500 dark:text-dark-gray-200 group-hover:text-gray-700 dark:group-hover:text-dark-gray-400"
                    />
                  </div>
                )}
                {ruleIndex !== 0 && (
                  <div
                    className="flex-center group cursor-pointer"
                    onClick={() => onChangeIndex('decrease', ruleIndex)}
                  >
                    <Icon
                      icon={IconArrowUpwardFilled}
                      size={'sm'}
                      className="text-gray-500 dark:text-dark-gray-200 group-hover:text-gray-700 dark:group-hover:text-dark-gray-400"
                    />
                  </div>
                )}
              </div>
            )}
            <div
              className="flex-center cursor-pointer group"
              onClick={() => onRemove(ruleIndex)}
            >
              <Icon
                icon={IconClose}
                size={'sm'}
                className="flex-center text-gray-500 dark:text-dark-gray-200 group-hover:text-gray-700 dark:group-hover:text-dark-gray-400"
              />
            </div>
          </div>
        </div>
        <RuleClausesForm
          rulesFieldName={RULES_FIELD_NAME}
          ruleIndex={ruleIndex}
          situationOptions={situationOptions}
          usedAttributeKeys={usedAttributeKeys}
          sdkAttributeKeys={sdkAttributeKeys}
        />
      </Card>
    </div>
  );
};

const SegmentRules = ({ disabled }: { disabled?: boolean }) => {
  const { t } = useTranslation(['form', 'common', 'table']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const { situationOptions } = useOptions();

  const { control, watch } = useFormContext<UserSegmentForm>();
  const { fields, append, remove, swap, move } = useFieldArray({
    control,
    name: RULES_FIELD_NAME,
    keyName: 'key'
  });

  const { data: keysCollection } = useQueryAttributeKeys({
    params: {
      environmentId: currentEnvironment.id
    }
  });
  const sdkAttributeKeys = keysCollection?.userAttributeKeys || [];

  const rulesWatch = watch(RULES_FIELD_NAME) as
    SegmentRuleFormValue[] | undefined;

  // Segment rules reject the SEGMENT and FEATURE_FLAG operators server-side;
  // only attribute comparisons and date conditions are offered.
  const segmentSituationOptions = useMemo(
    () =>
      (situationOptions as SituationOption[]).filter(option =>
        [RuleClauseType.COMPARE, RuleClauseType.DATE].includes(option.value)
      ),
    [situationOptions]
  );

  const usedAttributeKeys = useMemo(
    () =>
      compact(
        (rulesWatch || [])
          .flatMap(rule => rule?.clauses || [])
          .map(clause => clause.attribute)
      ),
    [rulesWatch]
  );

  const [, setActiveDragId] = useState<string | null>(null);

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: { distance: 5 }
    }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates
    })
  );

  const handleDragStart = useCallback((event: DragStartEvent) => {
    setActiveDragId(event.active.id as string);
  }, []);

  const handleDragEnd = useCallback(
    (event: DragEndEvent) => {
      const { active, over } = event;
      setActiveDragId(null);
      if (!over || active.id === over.id) return;
      const fromIndex = fields.findIndex(field => field.key === active.id);
      const toIndex = fields.findIndex(field => field.key === over.id);
      if (fromIndex !== -1 && toIndex !== -1) {
        move(fromIndex, toIndex);
      }
    },
    [fields, move]
  );

  const handleChangeIndex = useCallback(
    (type: 'increase' | 'decrease', currentIndex: number) => {
      swap(
        currentIndex,
        type === 'increase' ? currentIndex + 1 : currentIndex - 1
      );
    },
    [swap]
  );

  return (
    <div
      className={cn('flex flex-col w-full gap-y-4', {
        'pointer-events-none opacity-50': disabled
      })}
    >
      {fields.length > 0 && (
        <DndContext
          sensors={sensors}
          collisionDetection={closestCenter}
          onDragStart={handleDragStart}
          onDragEnd={handleDragEnd}
        >
          <SortableContext
            items={fields.map(field => field.key)}
            strategy={verticalListSortingStrategy}
          >
            <div className="flex flex-col w-full gap-y-3">
              {fields.map((field, ruleIndex) => (
                <Fragment key={field.key}>
                  {ruleIndex !== 0 && (
                    <div className="flex-center self-center w-[42px] h-[26px] rounded-[3px] typo-para-small leading-[14px] bg-gray-200 dark:bg-dark-black-700 text-gray-600 dark:text-dark-gray-200">
                      {t('common:or')}
                    </div>
                  )}
                  <SortableRuleCard
                    ruleKey={field.key}
                    ruleIndex={ruleIndex}
                    rulesCount={fields.length}
                    situationOptions={segmentSituationOptions}
                    usedAttributeKeys={usedAttributeKeys}
                    sdkAttributeKeys={sdkAttributeKeys}
                    onChangeIndex={handleChangeIndex}
                    onRemove={remove}
                  />
                </Fragment>
              ))}
            </div>
          </SortableContext>
        </DndContext>
      )}
      <Button
        type="button"
        variant={'text'}
        className="w-fit gap-x-2 h-6 !p-0"
        onClick={() => append(createSegmentRule())}
      >
        <Icon
          icon={IconPlus}
          color="primary-500"
          className="flex-center"
          size={'sm'}
        />{' '}
        {t('form:segment-rules.add-rule')}
      </Button>
    </div>
  );
};

export default SegmentRules;
