import { FC, memo, useCallback, useEffect } from 'react';
import { Controller, useFieldArray, useFormContext } from 'react-hook-form';

import { Select } from '../Select';
import { useCurrentEnvironment } from '../../modules/me';
import { AppDispatch } from '../../store';
import { shallowEqual, useDispatch, useSelector } from 'react-redux';
import {
  listFeatures,
  selectAll as selectAllFeatures
} from '../../modules/features';
import { ListFeaturesRequest } from '../../proto/feature/service_pb';
import { AppState } from '../../modules';
import { Feature } from '../../proto/feature/feature_pb';
import { PlusIcon, TrashIcon } from '@heroicons/react/outline';
import { v4 as uuid } from 'uuid';

export interface DebuggerEvaluateFormProps {
  onSubmit: () => void;
}

export const DebuggerEvaluateForm: FC<DebuggerEvaluateFormProps> = memo(
  ({ onSubmit }) => {
    const currentEnvironment = useCurrentEnvironment();
    const dispatch = useDispatch<AppDispatch>();

    const methods = useFormContext();
    const {
      register,
      control,
      formState: { errors, isValid, isSubmitting }
    } = methods;

    const {
      append: appendUserAttributes,
      remove: removeUserAttribute,
      fields: userAttributes
    } = useFieldArray({
      control,
      name: 'userAttributes'
    });

    const features = useSelector<AppState, Feature.AsObject[]>(
      (state) => selectAllFeatures(state.features),
      shallowEqual
    );

    useEffect(() => {
      dispatch(
        listFeatures({
          environmentNamespace: currentEnvironment.id,
          pageSize: 0,
          cursor: '',
          tags: [],
          searchKeyword: null,
          maintainerId: null,
          orderBy: ListFeaturesRequest.OrderBy.DEFAULT,
          orderDirection: ListFeaturesRequest.OrderDirection.ASC,
          archived: false
        })
      );
    }, []);

    const handleAddAttribute = () => {
      appendUserAttributes({
        id: uuid(),
        key: '',
        value: ''
      });
    };

    const handleDeleteAttribute = useCallback((index) => {
      removeUserAttribute(index);
    }, []);

    return (
      <form className="flex flex-col space-y-6">
        <div className="">
          <label htmlFor="flag">
            <span className="input-label">Flag</span>
          </label>
          <Controller
            name="flag"
            control={control}
            render={({ field }) => {
              const selectedOptions = features
                .filter((feature) => field.value.includes(feature.id))
                .map((feature) => ({
                  label: feature.name,
                  value: feature.id
                }));

              return (
                <Select
                  isMulti
                  options={features.map((feature) => ({
                    label: feature.name,
                    value: feature.id
                  }))}
                  value={selectedOptions} // Ensure the value prop is set correctly
                  onChange={(selected) => {
                    const newValues = selected
                      ? selected.map((o) => o.value)
                      : [];
                    field.onChange(newValues);
                  }}
                  placeholder="Select a feature flag"
                  disabled={isSubmitting}
                />
              );
            }}
          />
          <p className="input-error">
            {errors.flag && <span role="alert">{errors.flag.message}</span>}
          </p>
        </div>
        <div className="">
          <label htmlFor="userId">
            <span className="input-label">User ID</span>
          </label>
          <div className="mt-1">
            <input
              {...register('userId')}
              placeholder="Enter a user ID"
              type="text"
              id="userId"
              className="input-text w-full"
              disabled={isSubmitting}
            />
            <p className="input-error">
              {errors.userId && (
                <span role="alert">{errors.userId.message}</span>
              )}
            </p>
          </div>
        </div>
        <div className="space-y-4">
          {userAttributes.length > 0 && (
            <div>
              <div>
                <p>User Attributes</p>
                <p className="text-sm">
                  Nunc vulputate libero et velit interdum, ac aliquet odio
                  mattis.
                </p>
              </div>
              <div>
                {userAttributes.map((attr, index) => (
                  <div key={attr.id} className="flex space-x-4 mt-4 items-end">
                    <div className="flex flex-col flex-1">
                      <label htmlFor="key">
                        <span className="input-label">Key</span>
                      </label>
                      <input
                        {...register(`userAttributes.${index}.key`)}
                        type="text"
                        id="key"
                        className="input-text w-full"
                      />
                    </div>
                    <div className="flex flex-col flex-1">
                      <label htmlFor="value">
                        <span className="input-label">Value</span>
                      </label>
                      <input
                        {...register(`userAttributes.${index}.value`)}
                        type="text"
                        id="value"
                        className="input-text w-full"
                      />
                    </div>
                    <TrashIcon
                      width={18}
                      className="cursor-pointer text-gray-400 mb-3"
                      onClick={() => handleDeleteAttribute(index)}
                    />
                  </div>
                ))}
              </div>
            </div>
          )}
          <button
            className="flex whitespace-nowrap space-x-2 text-primary max-w-min py-2 items-center"
            type="button"
            onClick={handleAddAttribute}
          >
            <PlusIcon width={18} />
            <span>Add Attribute</span>
          </button>
        </div>
        <div className="flex">
          <button
            type="button"
            className="btn-submit"
            disabled={!isValid || isSubmitting}
            onClick={onSubmit}
          >
            Evaluate
          </button>
        </div>
      </form>
    );
  }
);
