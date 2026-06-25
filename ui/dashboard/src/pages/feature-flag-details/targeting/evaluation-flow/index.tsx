import { ReactNode } from 'react';
import { cn } from 'utils/style';
import {
  IconFlagSwitch,
  IconMember,
  IconPrerequisite,
  IconSetting,
  IconUserOutlined
} from '@icons';
import Icon from 'components/icon';

/**
 * Static labels for the evaluation flow captions. Intentionally not localised
 * — these read as semantic markers (`IF` / `ELSE IF` / `OTHERWISE`) and stay
 * the same across UI languages.
 */
export const FLOW_LABELS = {
  required: 'Required',
  if: 'If',
  elseIf: 'Else if',
  otherwise: 'Otherwise'
} as const;

/**
 * Horizontal offset (in px) for nodes/markers placed on the EvaluationFlow
 * spine but rendered outside `<FlowNode>` (e.g. the clickable plus inside
 * `<AddRule>`). Kept here so the value stays in lock-step with the spine x
 * position (14 px from the wrapper's left) and the FlowStep gutter (pl-14 →
 * 56 px). Reused via import to avoid magic numbers in callers.
 *
 * Math: spine_x − node_radius − pl_left = 14 − 10 − 56 = −52
 */
export const ADD_NODE_LEFT_OFFSET_PX = -52;

export type FlowKind =
  | 'start'
  | 'gate'
  | 'gate-off'
  | 'prerequisite'
  | 'individual'
  | 'rule'
  | 'add'
  | 'default';

/**
 * Outer container for the targeting evaluation flow.
 * Draws a continuous vertical "spine" on the left edge of the content column
 * so that every step (audience traffic, flag switch, prerequisites, rules,
 * default rule) reads as a single top-to-bottom evaluation pipeline.
 *
 * Steps should be wrapped in `<FlowStep>` (directly or nested inside another
 * component such as `<TargetSegmentRule>`) so their nodes align with the
 * spine. Non-FlowStep children render in the column without a spine node.
 */
export const EvaluationFlow = ({
  children,
  muted = false,
  className
}: {
  children: ReactNode;
  /** When true (e.g. flag is OFF), render the spine in a muted tone. */
  muted?: boolean;
  className?: string;
}) => {
  return (
    // Spine sits in the gutter to the left of the cards. Cards keep their
    // original horizontal position; only the spine moves inward, so 28px-wide
    // nodes have their left edge flush with the Targeting tab underline.
    <div
      className={cn(
        'relative w-full pl-14',
        'flex flex-col items-stretch gap-y-8',
        className
      )}
    >
      <div
        aria-hidden
        className={cn(
          'absolute top-3 bottom-3 w-px pointer-events-none rounded-full',
          muted ? 'bg-gray-200' : 'bg-gray-300'
        )}
        style={{ left: '14px' }}
      />
      {children}
    </div>
  );
};

interface FlowStepProps {
  kind: FlowKind;
  /** Caption shown between the node and the content (e.g. IF / ELSE IF / OTHERWISE). */
  stepLabel?: string;
  /**
   * Where to anchor the node vertically inside the step. 'top' aligns with the
   * card header (default for cards), 'center' for thin one-line rows.
   */
  align?: 'top' | 'center';
  /** Visual tone for the node and label. Default flows from kind. */
  tone?: 'default' | 'muted' | 'accent';
  children: ReactNode;
  className?: string;
}

/**
 * A single step in the evaluation flow. Wraps content in a relative container
 * and renders a `<FlowNode>` that anchors visually on the spine.
 */
export const FlowStep = ({
  kind,
  stepLabel,
  align = 'top',
  tone,
  children,
  className
}: FlowStepProps) => {
  const resolvedTone =
    tone ?? (kind === 'add' || kind === 'default' ? 'muted' : 'default');

  return (
    <div className={cn('relative w-full', className)}>
      <FlowNode kind={kind} align={align} tone={resolvedTone} />
      {stepLabel && (
        <span
          className={cn(
            'absolute z-10 uppercase tracking-wide typo-head-bold-tiny',
            'px-1.5 py-0.5 rounded select-none whitespace-nowrap',
            // Caption floats fully in the gap above each card so it reads as a
            // label *for* the card rather than overlapping its top border.
            align === 'center'
              ? 'top-1/2 -translate-y-1/2 left-0'
              : '-top-5 left-0',
            resolvedTone === 'muted'
              ? 'bg-gray-100 text-gray-500'
              : 'bg-primary-50 text-primary-500'
          )}
        >
          {stepLabel}
        </span>
      )}
      {children}
    </div>
  );
};

interface FlowNodeProps {
  kind: FlowKind;
  align: 'top' | 'center';
  tone: 'default' | 'muted' | 'accent';
}

const FlowNode = ({ kind, align, tone }: FlowNodeProps) => {
  const baseClasses =
    'absolute z-10 flex items-center justify-center rounded-full ring-4 ring-white';

  const alignClasses =
    align === 'center' ? 'top-1/2 -translate-y-1/2' : 'top-4';

  const toneClasses =
    tone === 'muted'
      ? 'bg-white border border-gray-300 text-gray-500'
      : tone === 'accent'
        ? 'bg-primary-500 text-white'
        : 'bg-white border border-primary-300 text-primary-500';

  let content: ReactNode = null;
  let sizeClass = 'size-7';
  let extraClass = '';

  switch (kind) {
    case 'start':
      sizeClass = 'size-7';
      extraClass = 'bg-primary-500 text-white border-0';
      content = <Icon icon={IconMember} size="xxs" className="!text-white" />;
      break;
    case 'gate':
      sizeClass = 'size-7';
      extraClass = 'bg-primary-50 border-primary-300';
      content = <Icon icon={IconFlagSwitch} size="xxs" color="primary-500" />;
      break;
    case 'gate-off':
      sizeClass = 'size-7';
      extraClass = 'bg-gray-100 border-gray-400 text-gray-500';
      content = <Icon icon={IconFlagSwitch} size="xxs" color="gray-500" />;
      break;
    case 'prerequisite':
      sizeClass = 'size-7';
      content = <Icon icon={IconPrerequisite} size="xxs" color="primary-500" />;
      break;
    case 'individual':
      sizeClass = 'size-7';
      content = <Icon icon={IconUserOutlined} size="xxs" color="primary-500" />;
      break;
    case 'rule':
      sizeClass = 'size-7';
      // Matches the icon used for "Custom Rule" in the AddRule dropdown.
      content = <Icon icon={IconSetting} size="xxs" color="primary-500" />;
      break;
    case 'add':
      // The spine plus is rendered inside <AddRule>'s dropdown trigger so the
      // entire plus circle is clickable (same affordance as "+ Add Rule").
      // FlowStep with kind="add" therefore emits no spine node of its own.
      return null;
    case 'default':
      sizeClass = 'size-7';
      extraClass = 'border-2';
      content = (
        <span className="size-2 rounded-full bg-gray-400" aria-hidden />
      );
      break;
  }

  // Spine is at x=14px from EvaluationFlow's left edge. FlowStep starts at
  // x=56px (pl-14), so to center a 28px-wide node on the spine its left edge
  // sits at -56px relative to FlowStep (14 - 14 - 56 = -56). The 'add' kind
  // returns early above (its plus is rendered inside <AddRule>), so it is
  // intentionally excluded from this map.
  const nodeOffsetByKind: Record<Exclude<FlowKind, 'add'>, number> = {
    start: -56,
    gate: -56,
    'gate-off': -56,
    prerequisite: -56,
    individual: -56,
    rule: -56,
    default: -56
  };

  return (
    <span
      aria-hidden
      className={cn(
        baseClasses,
        sizeClass,
        toneClasses,
        alignClasses,
        extraClass
      )}
      style={{ left: `${nodeOffsetByKind[kind]}px` }}
    >
      {content}
    </span>
  );
};
