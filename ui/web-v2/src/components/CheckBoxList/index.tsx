import { FC, memo, useEffect, useState } from 'react';

import { CheckBox } from '../CheckBox';
import { useIntl } from 'react-intl';
import { messages } from '../../lang/messages';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import { AppState } from '../../modules';
import { Tag } from '../../proto/tag/tag_pb';
import { listTags, selectAll as selectAllTags } from '../../modules/tags';
import { Subscription } from '../../proto/notification/subscription_pb';
import { Select } from '../Select';
import { Controller, useFormContext } from 'react-hook-form';
import { AppDispatch } from '../../store';
import { ListTagsRequest } from '../../proto/tag/service_pb';
import { useCurrentEnvironment } from '../../modules/me';
import { HoverPopover } from '../HoverPopover';
import { classNames } from '../../utils/css';
import { InformationCircleIcon } from '@heroicons/react/outline';

export interface Option {
  value: string;
  label: string;
  description?: string;
}

export interface CheckBoxListProps {
  onChange: (values: string[]) => void;
  options: Option[];
  defaultValues?: Option[];
  disabled?: boolean;
}

export const CheckBoxList: FC<CheckBoxListProps> = memo(
  ({ onChange, options, defaultValues, disabled }) => {
    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const {
      control,
      formState: { errors }
    } = methods;

    const dispatch = useDispatch<AppDispatch>();
    const currentEnvironment = useCurrentEnvironment();

    const tagsList = useSelector<AppState, Tag.AsObject[]>(
      (state) => selectAllTags(state.tags),
      shallowEqual
    );

    const featureFlagTagsList = tagsList.filter(
      (tag) => tag.entityType === Tag.EntityType.FEATURE_FLAG
    );

    const [checkedItems] = useState(() => {
      const items = new Map();
      defaultValues &&
        defaultValues.forEach((item) => {
          items.set(item.value, item.value);
        });
      return items;
    });

    useEffect(() => {
      dispatch(
        listTags({
          environmentId: currentEnvironment.id,
          pageSize: 0,
          cursor: '',
          orderBy: ListTagsRequest.OrderBy.DEFAULT,
          orderDirection: ListTagsRequest.OrderDirection.ASC,
          searchKeyword: null
        })
      );
    }, [dispatch]);

    const handleOnChange = (value: string, checked: boolean) => {
      if (checked) {
        checkedItems.set(value, value);
      } else {
        checkedItems.delete(value);
      }
      const valueList = [];
      checkedItems.forEach((v) => {
        valueList.push(v);
      });
      onChange(valueList);
    };

    return (
      <div>
        <fieldset className="border-t border-b border-gray-300">
          <div className="divide-y divide-gray-300">
            {options.map((item, index) => {
              return (
                <div key={item.label}>
                  <div className="relative flex items-start py-4">
                    <div className="min-w-0 flex-1 text-sm">
                      <label htmlFor={`id_${index}`} key={`key_${index}`}>
                        <p className="text-sm font-medium text-gray-700">
                          {item.label}
                        </p>
                        {item.description && (
                          <p className="text-sm text-gray-500">
                            {item.description}
                          </p>
                        )}
                      </label>
                    </div>
                    <div className="ml-3 flex items-center h-5">
                      <CheckBox
                        id={`id_${index}`}
                        value={item.value}
                        onChange={handleOnChange}
                        defaultChecked={checkedItems.has(item.value)}
                        disabled={disabled}
                      />
                    </div>
                  </div>
                  {item.value ===
                    Subscription.SourceType.DOMAIN_EVENT_FEATURE.toString() &&
                    checkedItems.has(
                      Subscription.SourceType.DOMAIN_EVENT_FEATURE.toString()
                    ) && (
                      <div className="mb-4">
                        <div className="flex space-x-1 items-center">
                          <label htmlFor="tags">
                            <span className="input-label">
                              {f(messages.tags.title)}
                            </span>
                          </label>
                          <HoverPopover
                            render={() => {
                              return (
                                <div
                                  className={classNames(
                                    'border shadow-sm bg-gray-900 text-white p-1',
                                    'text-xs rounded whitespace-normal break-words w-80'
                                  )}
                                >
                                  {f(messages.notification.tagsTooltipMessage)}
                                </div>
                              );
                            }}
                          >
                            <div
                              className={classNames(
                                'hover:text-gray-500 mb-[2px]'
                              )}
                            >
                              <InformationCircleIcon
                                className="w-5 h-5 text-gray-400"
                                aria-hidden="true"
                              />
                            </div>
                          </HoverPopover>
                        </div>
                        <Controller
                          name="featureFlagTagsList"
                          control={control}
                          render={({ field }) => {
                            return (
                              <Select
                                isMulti
                                value={field.value?.map((tag: string) => {
                                  return {
                                    value: tag,
                                    label: tag
                                  };
                                })}
                                options={featureFlagTagsList.map((tag) => ({
                                  label: tag.name,
                                  value: tag.name
                                }))}
                                onChange={(options: Option[]) => {
                                  field.onChange(options.map((o) => o.value));
                                }}
                                closeMenuOnSelect={false}
                                placeholder={f(messages.tags.tagsPlaceholder)}
                                disabled={disabled}
                              />
                            );
                          }}
                        />
                        {errors.tags && (
                          <p className="input-error">
                            {errors.tags && (
                              <span role="alert">{errors.tags.message}</span>
                            )}
                          </p>
                        )}
                      </div>
                    )}
                </div>
              );
            })}
          </div>
        </fieldset>
      </div>
    );
  }
);
