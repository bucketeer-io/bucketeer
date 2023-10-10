import { Dialog } from '@headlessui/react';
import React, { FC, memo } from 'react';
import { useFormContext, Controller } from 'react-hook-form';
import { useIntl } from 'react-intl';

import { messages } from '../../lang/messages';
import { useCurrentEnvironment, useEnvironments } from '../../modules/me';
import { Select, Option } from '../Select';

export interface FeatureCloneFormProps {
  onSubmit: () => void;
  onCancel: () => void;
}

export const FeatureCloneForm: FC<FeatureCloneFormProps> = memo(
  ({ onSubmit, onCancel }) => {
    const currenEnvironment = useCurrentEnvironment();
    const environments = useEnvironments();
    const { formatMessage: f } = useIntl();
    const methods = useFormContext();
    const {
      control,
      getValues,
      formState: { errors, isValid, isSubmitted },
    } = methods;
    const filteredOptions = environments.filter((environment) => {
      return environment.id != currenEnvironment.id;
    });
    const environmentOptions = filteredOptions.map((environment) => {
      return {
        value: environment.id,
        label: environment.name,
      };
    });
    return (
      <div className="w-[500px]">
        <form className="flex flex-col">
          <div className="flex-1 h-0">
            <div className="py-6 px-4 bg-primary">
              <div className="flex items-center justify-between">
                <Dialog.Title className="text-lg font-medium text-white">
                  {f(messages.feature.clone.header.title)}
                </Dialog.Title>
              </div>
              <div className="mt-1">
                <p className="text-sm text-indigo-300">
                  {f(messages.feature.clone.header.description)}
                </p>
              </div>
            </div>
            <div
              className="
                flex-1
                flex flex-col
                justify-between
              "
            >
              <div
                className="
                  space-y-6 px-5 pt-6 pb-5
                  flex flex-col
                "
              >
                <div className="">
                  <label htmlFor="featureName">
                    <span className="input-label">
                      {f(messages.input.featureFlag)}
                    </span>
                  </label>
                  <div className="mt-1">
                    <input
                      type="text"
                      name="featureName"
                      id="featureName"
                      className="input-text w-full"
                      defaultValue={
                        getValues().feature && getValues().feature.name
                      }
                      disabled={true}
                    />
                  </div>
                </div>
                <div className="">
                  <label htmlFor="originEnvironmentId">
                    <span className="input-label">
                      {f(messages.input.originEnvironment)}
                    </span>
                  </label>
                  <div className="mt-1">
                    <input
                      type="text"
                      name="originEnvironmentName"
                      id="originEnvironmentName"
                      className="input-text w-full"
                      defaultValue={currenEnvironment.name}
                      disabled={true}
                    />
                  </div>
                </div>
                <div className="">
                  <label htmlFor="destinationEnvironmentId" className="block">
                    <span className="input-label">
                      {f(messages.input.destinationEnvironment)}
                    </span>
                  </label>
                  <div className="mt-1">
                    <Controller
                      name="destinationEnvironmentId"
                      control={control}
                      render={({ field }) => {
                        return (
                          <Select
                            onChange={(o: Option) => field.onChange(o.value)}
                            options={environmentOptions}
                            disabled={isSubmitted}
                          />
                        );
                      }}
                    />
                    <p className="input-error">
                      {errors.destinationEnvironmentId && (
                        <span role="alert">
                          {errors.destinationEnvironmentId.message}
                        </span>
                      )}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div className="flex-shrink-0 px-4 py-4 flex justify-end">
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
              className="btn-submit"
              disabled={!isValid || isSubmitted}
              onClick={onSubmit}
            >
              {f(messages.button.submit)}
            </button>
          </div>
        </form>
      </div>
    );
  }
);
