import { useId, useLayoutEffect, useRef, useState } from 'react';
import { GroupBase, mergeStyles, SelectInstance } from 'react-select';
import ReactCreatableSelect, { CreatableProps } from 'react-select/creatable';
import * as Popover from '@radix-ui/react-popover';
import { useTranslation } from 'i18n';
import { cn } from 'utils/style';
import Spinner from 'components/spinner';
import { getVisibleValueCount } from './compact-values';
import type { Option } from './index';

type Props = CreatableProps<Option, true, GroupBase<Option>> & {
  'aria-describedby'?: string;
};

export const CompactSelect = (props: Props) => {
  const { t } = useTranslation(['form', 'common']);
  const editorTitleId = useId();
  const editorId = useId();
  const [editorOpen, setEditorOpen] = useState(false);
  const [internalValue, setInternalValue] = useState(props.defaultValue ?? []);
  const selectedValues =
    props.value === undefined ? internalValue : props.value;
  const values: readonly Option[] = selectedValues
    ? Array.isArray(selectedValues)
      ? selectedValues
      : [selectedValues as Option]
    : [];
  const triggerRef = useRef<HTMLButtonElement>(null);
  const valuesRef = useRef<HTMLSpanElement>(null);
  const editorRef = useRef<SelectInstance<Option, true>>(null);
  const restoringFocus = useRef(false);
  const interactedOutside = useRef(false);
  const [width, setWidth] = useState(0);
  const visibleCount = getVisibleValueCount(width, values.length);
  const hiddenCount = values.length - visibleCount;
  const open = editorOpen && !props.isDisabled;

  useLayoutEffect(() => {
    const container = valuesRef.current;
    if (!container) return;
    setWidth(container.clientWidth);
    const observer = new ResizeObserver(([entry]) => {
      setWidth(entry.contentRect.width);
    });
    observer.observe(container);
    return () => observer.disconnect();
  }, []);

  const openEditor = () => {
    if (props.isDisabled || restoringFocus.current) return;
    interactedOutside.current = false;
    setEditorOpen(true);
  };

  return (
    <Popover.Root open={open} onOpenChange={setEditorOpen}>
      <Popover.Anchor asChild>
        <button
          ref={triggerRef}
          id={props.inputId ?? props.id}
          type="button"
          disabled={props.isDisabled}
          autoFocus={props.autoFocus}
          tabIndex={props.tabIndex}
          aria-haspopup="dialog"
          aria-expanded={open}
          aria-controls={open ? editorId : undefined}
          aria-label={props['aria-label'] ?? t('form:feature-flags.values')}
          aria-labelledby={props['aria-labelledby']}
          aria-describedby={props['aria-describedby']}
          aria-invalid={props['aria-invalid']}
          className={cn(
            'flex items-center gap-2 w-full min-w-0 h-12 px-2 rounded-lg border border-gray-300 bg-white text-left text-gray-700 typo-para-small focus-visible:outline focus-visible:outline-2 focus-visible:outline-primary-500 disabled:cursor-not-allowed disabled:bg-gray-100 disabled:text-gray-500',
            props.className
          )}
          onFocus={openEditor}
          onClick={openEditor}
        >
          <span
            ref={valuesRef}
            className="flex items-center gap-1 flex-1 min-w-0 overflow-hidden"
          >
            {values.length ? (
              <>
                {values.slice(0, visibleCount).map(option => {
                  const color =
                    typeof option.color === 'string' ? option.color : undefined;
                  return (
                    <span
                      key={option.value}
                      title={option.label}
                      className="min-w-0 max-w-[140px] truncate rounded bg-primary-100 px-2 py-1 text-primary-600"
                      style={
                        color
                          ? { color, backgroundColor: `${color}1A` }
                          : undefined
                      }
                    >
                      {option.label}
                    </span>
                  );
                })}
                {hiddenCount > 0 && (
                  <span className="shrink-0 px-1 font-medium text-primary-500">
                    +{hiddenCount}
                  </span>
                )}
              </>
            ) : (
              <span className="truncate text-gray-400">
                {props.placeholder || t('form:feature-flags.values')}
              </span>
            )}
          </span>
          {props.isLoading && <Spinner className="size-5 shrink-0" />}
        </button>
      </Popover.Anchor>
      <Popover.Portal>
        <Popover.Content
          id={editorId}
          aria-labelledby={editorTitleId}
          align="start"
          sideOffset={8}
          className="z-[100] w-[480px] max-w-[calc(100vw-32px)] max-h-[var(--radix-popover-content-available-height)] overflow-y-auto rounded-lg bg-white p-4 shadow-dropdown"
          onOpenAutoFocus={event => {
            event.preventDefault();
            editorRef.current?.focus();
          }}
          onInteractOutside={event => {
            if (triggerRef.current?.contains(event.target as Node)) {
              event.preventDefault();
            } else {
              interactedOutside.current = true;
            }
          }}
          onCloseAutoFocus={event => {
            event.preventDefault();
            if (interactedOutside.current) return;
            // Return keyboard focus without reopening the editor.
            restoringFocus.current = true;
            triggerRef.current?.focus();
            restoringFocus.current = false;
          }}
        >
          <div className="flex items-center justify-between gap-4 mb-3">
            <span id={editorTitleId} className="typo-para-small font-medium">
              {t('form:feature-flags.values')}
            </span>
            <Popover.Close className="typo-para-small text-primary-500">
              {t('common:close')}
            </Popover.Close>
          </div>
          <ReactCreatableSelect
            {...props}
            value={selectedValues}
            onChange={(newValue, actionMeta) => {
              if (props.value === undefined) setInternalValue(newValue);
              props.onChange?.(newValue, actionMeta);
            }}
            ref={editorRef}
            className="w-full"
            id={undefined}
            inputId={undefined}
            instanceId={undefined}
            name={undefined}
            aria-label={t('form:feature-flags.values')}
            aria-labelledby={undefined}
            menuPortalTarget={null}
            menuPosition="absolute"
            styles={mergeStyles(props.styles ?? {}, {
              valueContainer: base => ({
                ...base,
                maxHeight: 200,
                overflowY: 'auto'
              }),
              menu: base => ({ ...base, position: 'relative' })
            })}
          />
        </Popover.Content>
      </Popover.Portal>
    </Popover.Root>
  );
};
