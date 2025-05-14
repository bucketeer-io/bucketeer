import { useCallback, useState } from 'react';
import { useQueryTriggers } from '@queries/triggers';
import { getCurrentEnvironment, useAuth } from 'auth';
import { useTranslation } from 'i18n';
import { Feature, TriggerItemType } from '@types';
import { IconPlus } from '@icons';
import Button from 'components/button';
import Icon from 'components/icon';
import Card from 'elements/card';
import FormLoading from 'elements/form-loading';
import CreateTriggerForm from '../create-trigger-form';
import { TriggerAction } from '../types';
import TriggerItem from './trigger-item';

const TriggerList = ({ feature }: { feature: Feature }) => {
  const { t } = useTranslation(['table']);
  const { consoleAccount } = useAuth();
  const currentEnvironment = getCurrentEnvironment(consoleAccount!);
  const [triggerNewlyCreated, setTriggerNewlyCreated] = useState<
    TriggerItemType | undefined
  >(undefined);
  const [isShowCreateForm, setIsShowCreateForm] = useState(false);

  const { data: triggerCollection, isLoading } = useQueryTriggers({
    params: {
      environmentId: currentEnvironment.id,
      featureId: feature.id,
      cursor: String(0)
    }
  });

  const triggers = triggerCollection?.flagTriggers || [];

  const onActions = useCallback(
    (trigger: TriggerItemType, action: TriggerAction) => {
      console.log({ trigger, action });
    },
    []
  );

  return (
    <Card className="gap-y-6">
      <p className="typo-head-bold-small text-gray-800">{t('trigger.title')}</p>
      <p className="typo-para-small text-gray-500">
        {t('trigger.description')}
      </p>
      {isLoading ? (
        <FormLoading />
      ) : (
        <>
          {triggers.map((trigger, index) => (
            <TriggerItem
              key={index}
              trigger={trigger}
              triggerNewlyCreated={triggerNewlyCreated}
              onActions={action => onActions(trigger, action)}
            />
          ))}
          {isShowCreateForm ? (
            <CreateTriggerForm
              featureId={feature.id}
              environmentId={currentEnvironment.id}
              onCancel={() => setIsShowCreateForm(false)}
              setTriggerNewlyCreated={setTriggerNewlyCreated}
            />
          ) : (
            <Button
              variant="text"
              className="h-8 w-fit p-0"
              onClick={() => setIsShowCreateForm(true)}
            >
              <Icon icon={IconPlus} size="md" />
              {t('trigger.add-trigger')}
            </Button>
          )}
        </>
      )}
    </Card>
  );
};

export default TriggerList;
