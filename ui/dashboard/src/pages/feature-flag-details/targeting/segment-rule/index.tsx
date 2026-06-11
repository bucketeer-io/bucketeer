import { useCallback, useState } from 'react';
import { useFormContext } from 'react-hook-form';
import {
  DndContext,
  DragEndEvent,
  DragOverlay,
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
import { useQueryAttributeKeys } from '@queries/attribute-keys';
import { useQueryUserSegments } from '@queries/user-segments';
import { getCurrentEnvironment, useAuth } from 'auth';
import { LIST_PAGE_SIZE } from 'constants/app';
import { Feature } from '@types';
import { TargetingDivider } from '..';
import AddRule from '../add-rule';
import { RuleSchema, TargetingSchema } from '../form-schema';
import { DiscardChangesType, RuleCategory } from '../types';
import SortableCard, { DragOverlayCard } from './sortable-card';

interface RuleSchemaFields extends RuleSchema {
  segmentId: string;
}

interface Props {
  feature: Feature;
  features: Feature[];
  segmentRules: RuleSchemaFields[];
  isDisableAddIndividualRules: boolean;
  isDisableAddPrerequisite: boolean;
  onAddRule: (rule: RuleCategory, index?: number) => void;
  segmentRulesRemove: (index: number) => void;
  segmentRulesSwap: (indexA: number, indexB: number) => void;
  segmentRulesMove: (fromIndex: number, toIndex: number) => void;
  handleDiscardChanges: (type: DiscardChangesType, index?: number) => void;
  handleCheckEdit: (type: RuleCategory, index?: number) => boolean;
}

const TargetSegmentRule = ({
  feature,
  features,
  segmentRules,
  isDisableAddIndividualRules,
  isDisableAddPrerequisite,
  onAddRule,
  segmentRulesRemove,
  segmentRulesSwap,
  segmentRulesMove,
  handleDiscardChanges,
  handleCheckEdit
}: Props) => {
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);

  useFormContext<TargetingSchema>();

  const { data: segmentCollection } = useQueryUserSegments({
    params: {
      cursor: String(0),
      pageSize: LIST_PAGE_SIZE,
      environmentId: currentEnvironment.id
    }
  });

  const { data: keysCollection } = useQueryAttributeKeys({
    params: {
      environmentId: currentEnvironment.id
    }
  });

  const userSegments = segmentCollection?.segments || [];
  const sdkAttributeKeys = keysCollection?.userAttributeKeys || [];

  const [activeDragId, setActiveDragId] = useState<string | null>(null);
  const [activeDragHeight, setActiveDragHeight] = useState<number | null>(null);

  const editSegmentRule = (index: number) =>
    handleCheckEdit(RuleCategory.CUSTOM, index);

  const handleChangeIndexRule = useCallback(
    (type: 'increase' | 'decrease', currentIndex: number) => {
      segmentRulesSwap(
        currentIndex,
        type === 'increase' ? currentIndex + 1 : currentIndex - 1
      );
    },
    [segmentRulesSwap]
  );

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
    setActiveDragHeight(null);
  }, []);

  const handleDragEnd = useCallback(
    (event: DragEndEvent) => {
      const { active, over } = event;
      setActiveDragId(null);
      setActiveDragHeight(null);
      if (!over || active.id === over.id) return;
      const fromIndex = segmentRules.findIndex(r => r.segmentId === active.id);
      const toIndex = segmentRules.findIndex(r => r.segmentId === over.id);
      if (fromIndex !== -1 && toIndex !== -1) {
        segmentRulesMove(fromIndex, toIndex);
      }
    },
    [segmentRules, segmentRulesMove]
  );

  const activeDragIndex = activeDragId
    ? segmentRules.findIndex(r => r.segmentId === activeDragId)
    : -1;
  const activeDragSegment =
    activeDragIndex !== -1 ? segmentRules[activeDragIndex] : null;

  const sortableIds = segmentRules.map(r => r.segmentId);

  if (segmentRules.length === 0) return null;

  const sharedCardProps = {
    feature,
    features,
    segmentRules,
    userSegments,
    sdkAttributeKeys,
    editSegmentRule,
    handleChangeIndexRule,
    handleDiscardChanges,
    segmentRulesRemove
  };

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      onDragStart={handleDragStart}
      onDragEnd={handleDragEnd}
    >
      <SortableContext
        items={sortableIds}
        strategy={verticalListSortingStrategy}
      >
        <div className="flex flex-col w-full">
          {segmentRules.map((segment, segmentIndex) => (
            <div key={segment.segmentId} className="flex flex-col w-full">
              {segmentIndex !== 0 && (
                <>
                  <TargetingDivider />
                  <AddRule
                    isDisableAddIndividualRules={isDisableAddIndividualRules}
                    isDisableAddPrerequisite={isDisableAddPrerequisite}
                    onAddRule={onAddRule}
                    indexInsertSegmentRule={segmentIndex}
                    isInsertSegmentRule={true}
                  />
                  <TargetingDivider />
                </>
              )}
              <SortableCard
                segment={segment}
                segmentIndex={segmentIndex}
                ghostHeight={activeDragHeight}
                {...sharedCardProps}
              />
            </div>
          ))}
        </div>
      </SortableContext>
      <DragOverlay>
        {activeDragSegment && activeDragIndex !== -1 && (
          <div
            className="opacity-95 shadow-lg"
            ref={el => {
              if (el && activeDragHeight === null)
                setActiveDragHeight(el.getBoundingClientRect().height);
            }}
          >
            <DragOverlayCard
              segment={activeDragSegment}
              segmentIndex={activeDragIndex}
              feature={feature}
              features={features}
              userSegments={userSegments}
              sdkAttributeKeys={sdkAttributeKeys}
            />
          </div>
        )}
      </DragOverlay>
    </DndContext>
  );
};

export default TargetSegmentRule;
