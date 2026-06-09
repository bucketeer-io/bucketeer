import { useForm, FormProvider } from 'react-hook-form';
import { Trans } from 'react-i18next';
import {
  IconArrowDownwardFilled,
  IconArrowUpwardFilled,
  IconUndoOutlined
} from 'react-icons-material-design';
import { Fragment } from 'react/jsx-runtime';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { useTranslation } from 'i18n';
import { GripVertical } from 'lucide-react';
import { Feature, UserSegment } from '@types';
import { IconClose, IconDropZone, IconInfo } from '@icons';
import Icon from 'components/icon';
import { Tooltip } from 'components/tooltip';
import Card from '../../elements/card';
import { RuleSchema, TargetingSchema } from '../form-schema';
import { DiscardChangesType } from '../types';
import RuleForm from './rule';
import SegmentVariation from './variation';

export interface DragOverlayCardProps {
  segment: RuleSchemaFields;
  segmentIndex: number;
  feature: Feature;
  features: Feature[];
  userSegments: UserSegment[];
  sdkAttributeKeys: string[];
}

interface RuleSchemaFields extends RuleSchema {
  segmentId: string;
}

export interface SortableCardProps {
  segment: RuleSchemaFields;
  segmentIndex: number;
  segmentRules: RuleSchemaFields[];
  feature: Feature;
  features: Feature[];
  userSegments: UserSegment[];
  sdkAttributeKeys: string[];
  ghostHeight?: number | null;
  editSegmentRule: (index: number) => boolean;
  handleChangeIndexRule: (type: 'increase' | 'decrease', index: number) => void;
  handleDiscardChanges: (type: DiscardChangesType, index?: number) => void;
  segmentRulesRemove: (index: number) => void;
}

export interface CardContentProps extends Omit<
  SortableCardProps,
  'isDisableAddIndividualRules' | 'isDisableAddPrerequisite' | 'onAddRule'
> {
  displayIndex?: number;
  dragHandleProps?: {
    attributes: ReturnType<typeof useSortable>['attributes'];
    listeners: ReturnType<typeof useSortable>['listeners'];
  };
}

export const CardContent = ({
  segmentIndex,
  displayIndex,
  segmentRules,
  feature,
  features,
  userSegments,
  sdkAttributeKeys,
  editSegmentRule,
  handleChangeIndexRule,
  handleDiscardChanges,
  segmentRulesRemove,
  dragHandleProps
}: CardContentProps) => {
  const { t } = useTranslation(['table', 'form']);

  return (
    <Card>
      <div className="w-full h-8 flex items-center justify-between">
        <div className="flex items-center gap-x-2">
          {segmentRules.length > 1 && (
            <div
              className="flex-center cursor-grab active:cursor-grabbing text-gray-400 hover:text-gray-600 dark:text-dark-gray-200 dark:hover:text-dark-gray-400 touch-none"
              {...dragHandleProps?.attributes}
              {...dragHandleProps?.listeners}
            >
              <GripVertical size={16} />
            </div>
          )}
          <p className="typo-para-medium leading-5 text-gray-700 dark:text-dark-gray-400">
            <Trans
              i18nKey={'table:feature-flags.rule-index'}
              values={{ index: (displayIndex ?? segmentIndex) + 1 }}
            />
          </p>
          <Tooltip
            align="start"
            alignOffset={-68}
            content={t('form:targeting.tooltip.custom-rules')}
            trigger={
              <div className="flex-center size-fit">
                <Icon icon={IconInfo} size={'xxs'} color="gray-500" />
              </div>
            }
            className="max-w-[400px]"
          />
        </div>
        <div className="flex items-center gap-x-2">
          {editSegmentRule(segmentIndex) && (
            <div
              className="flex-center h-8 w-8 px-2 rounded-md cursor-pointer group border border-gray-300 hover:border-gray-800 dark:border-dark-black-700 dark:hover:border-dark-purple-300"
              onClick={() =>
                handleDiscardChanges(DiscardChangesType.CUSTOM, segmentIndex)
              }
            >
              <Icon
                icon={IconUndoOutlined}
                size={'sm'}
                className="flex-center text-gray-500 group-hover:text-gray-700 dark:text-dark-gray-200 dark:group-hover:text-dark-gray-400"
              />
            </div>
          )}
          {segmentRules.length > 1 && (
            <div className="flex items-center gap-x-1">
              {segmentIndex !== segmentRules.length - 1 && (
                <div
                  className="flex-center group cursor-pointer"
                  onClick={() =>
                    handleChangeIndexRule('increase', segmentIndex)
                  }
                >
                  <Icon
                    icon={IconArrowDownwardFilled}
                    size={'sm'}
                    className="text-gray-500 group-hover:text-gray-700 dark:text-dark-gray-200 dark:group-hover:text-dark-gray-400"
                  />
                </div>
              )}
              {segmentIndex !== 0 && (
                <div
                  className="flex-center group cursor-pointer"
                  onClick={() =>
                    handleChangeIndexRule('decrease', segmentIndex)
                  }
                >
                  <Icon
                    icon={IconArrowUpwardFilled}
                    size={'sm'}
                    className="text-gray-500 group-hover:text-gray-700 dark:text-dark-gray-200 dark:group-hover:text-dark-gray-400"
                  />
                </div>
              )}
            </div>
          )}
          <div
            className="flex-center cursor-pointer group"
            onClick={() => segmentRulesRemove(segmentIndex)}
          >
            <Icon
              icon={IconClose}
              size={'sm'}
              className="flex-center text-gray-500 group-hover:text-gray-700 dark:text-dark-gray-200 dark:group-hover:text-dark-gray-400"
            />
          </div>
        </div>
      </div>
      <Fragment>
        <RuleForm
          feature={feature}
          features={features}
          segmentIndex={segmentIndex}
          userSegments={userSegments}
          sdkAttributeKeys={sdkAttributeKeys}
        />
        <SegmentVariation
          feature={feature}
          segmentIndex={segmentIndex}
          segmentRules={segmentRules}
        />
      </Fragment>
    </Card>
  );
};

export const DragOverlayCard = ({
  segment,
  segmentIndex,
  feature,
  features,
  userSegments,
  sdkAttributeKeys
}: DragOverlayCardProps) => {
  // Mount RuleForm/SegmentVariation in an isolated FormProvider seeded with a
  // frozen snapshot at index 0, so their useFieldArray / Form.Field hooks never
  // touch the live parent form store.
  const snapshotMethods = useForm<TargetingSchema>({
    defaultValues: {
      segmentRules: [segment]
    } as unknown as TargetingSchema
  });

  const noop = () => false;

  return (
    <FormProvider {...snapshotMethods}>
      <CardContent
        segment={segment}
        segmentIndex={0}
        displayIndex={segmentIndex}
        segmentRules={[segment]}
        feature={feature}
        features={features}
        userSegments={userSegments}
        sdkAttributeKeys={sdkAttributeKeys}
        editSegmentRule={noop}
        handleChangeIndexRule={() => {}}
        handleDiscardChanges={() => {}}
        segmentRulesRemove={() => {}}
      />
    </FormProvider>
  );
};

const SortableCard = ({
  segment,
  segmentIndex,
  segmentRules,
  feature,
  features,
  userSegments,
  sdkAttributeKeys,
  ghostHeight,
  editSegmentRule,
  handleChangeIndexRule,
  handleDiscardChanges,
  segmentRulesRemove
}: SortableCardProps) => {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging
  } = useSortable({ id: segment.segmentId });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition
  };

  return (
    <div ref={setNodeRef} style={style} className="flex flex-col w-full">
      {isDragging ? (
        <div
          className="flex items-center justify-center gap-x-2 w-full rounded-lg border-2 border-dashed border-primary-500 bg-primary-50 text-primary-500 dark:border-dark-purple-300 dark:bg-dark-purple-100 dark:text-dark-purple-700"
          style={{ height: ghostHeight ?? undefined }}
        >
          <Icon icon={IconDropZone} className="w-[120px] h-[120px]" />
        </div>
      ) : (
        <CardContent
          segment={segment}
          segmentIndex={segmentIndex}
          segmentRules={segmentRules}
          feature={feature}
          features={features}
          userSegments={userSegments}
          sdkAttributeKeys={sdkAttributeKeys}
          editSegmentRule={editSegmentRule}
          handleChangeIndexRule={handleChangeIndexRule}
          handleDiscardChanges={handleDiscardChanges}
          segmentRulesRemove={segmentRulesRemove}
          dragHandleProps={{ attributes, listeners }}
        />
      )}
    </div>
  );
};

export default SortableCard;
