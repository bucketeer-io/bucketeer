import { useFormContext } from 'react-hook-form';
import { useTranslation } from 'i18n';
import { IconInfo, IconWatch } from '@icons';
import Button from 'components/button';
import { ButtonBar } from 'components/button-bar';
import { ReactDatePicker } from 'components/date-time-picker';
import Form from 'components/form';
import Icon from 'components/icon';
import DialogModal from 'components/modal/dialog';
import { RadioGroup, RadioGroupItem } from 'components/radio';
import TextArea from 'components/textarea';
import { Tooltip } from 'components/tooltip';
import {
  SCHEDULE_TYPE_SCHEDULE,
  SCHEDULE_TYPE_UPDATE_NOW
} from '../../../../elements/confirm-required-modal/form-schema';

export type SaveWithCommentModalProps = {
  isOpen: boolean;
  isRequired: boolean;
  isShowScheduleSelect?: boolean;
  onSubmit: (scheduleType?: string, scheduleAt?: string) => void;
  onClose: () => void;
};

const SaveWithCommentModal = ({
  isOpen,
  isRequired,
  isShowScheduleSelect,
  onClose,
  onSubmit
}: SaveWithCommentModalProps) => {
  const { t } = useTranslation(['common', 'form']);
  const {
    control,
    formState: { isValid, isSubmitting },
    handleSubmit,
    watch
  } = useFormContext();

  const scheduleType = watch('scheduleType');
  const isShowSchedule = scheduleType === SCHEDULE_TYPE_SCHEDULE;

  return (
    <DialogModal
      className="w-[500px]"
      title={t('update-flag')}
      isOpen={isOpen}
      onClose={onClose}
    >
      <div className="flex flex-col w-full items-start px-5 py-8 gap-y-5">
        <Form.Field
          name="comment"
          control={control}
          render={({ field }) => (
            <Form.Item className="py-0 w-full">
              <Form.Label
                required={isRequired && !isShowSchedule}
                optional={!isRequired || isShowSchedule}
              >
                {t('form:comment-for-update')}
              </Form.Label>
              <Form.Control>
                <TextArea
                  {...field}
                  placeholder={t('form:placeholder-comment')}
                />
              </Form.Control>
              <Form.Message />
            </Form.Item>
          )}
        />

        {isShowScheduleSelect && (
          <>
            <Form.Field
              control={control}
              name="scheduleType"
              render={({ field }) => (
                <Form.Item className="flex flex-col w-full py-0 gap-y-4">
                  <Form.Control>
                    <RadioGroup
                      defaultValue={field.value || SCHEDULE_TYPE_UPDATE_NOW}
                      className="flex flex-col w-full gap-y-4"
                      onValueChange={value => {
                        field.onChange(value);
                      }}
                    >
                      <div className="flex items-center gap-x-2">
                        <RadioGroupItem
                          id="settings_active_now"
                          value={SCHEDULE_TYPE_UPDATE_NOW}
                        />
                        <label
                          htmlFor="settings_active_now"
                          className="typo-para-medium leading-4 text-gray-700 cursor-pointer"
                        >
                          {t('update-now')}
                        </label>
                      </div>

                      <div className="flex items-center gap-x-2">
                        <RadioGroupItem
                          id="settings_schedule"
                          value={SCHEDULE_TYPE_SCHEDULE}
                        />
                        <label
                          htmlFor="settings_schedule"
                          className="typo-para-medium leading-4 text-gray-700 cursor-pointer"
                        >
                          {t('form:feature-flags.schedule-the-updates')}
                        </label>
                        <span className="px-2 py-1.5 rounded-[3px] bg-accent-blue-50 text-accent-blue-500 typo-para-small leading-[14px] whitespace-nowrap uppercase">
                          New
                        </span>
                        <Tooltip
                          align="start"
                          content={t(
                            'form:feature-flags.schedule-the-updates-tooltip'
                          )}
                          trigger={
                            <div className="flex-center size-fit">
                              <Icon
                                icon={IconInfo}
                                size="xs"
                                color="gray-500"
                              />
                            </div>
                          }
                          className="max-w-[400px]"
                        />
                      </div>
                    </RadioGroup>
                  </Form.Control>
                  <Form.Message />
                </Form.Item>
              )}
            />
            {isShowSchedule && (
              <Form.Field
                control={control}
                name="scheduleAt"
                render={({ field }) => {
                  const scheduleDate = field.value
                    ? new Date(+field.value * 1000)
                    : null;

                  return (
                    <Form.Item className="py-0">
                      <Form.Control>
                        <div className="flex gap-x-4">
                          <div>
                            <Form.Label required>
                              {t('form:feature-flags.update-date')}
                            </Form.Label>
                            <ReactDatePicker
                              dateFormat="yyyy/MM/dd"
                              minDate={new Date()}
                              selected={scheduleDate}
                              showTimeSelect={false}
                              className="w-[186px]"
                              onChange={date => {
                                if (date) {
                                  if (scheduleDate) {
                                    date.setHours(
                                      scheduleDate.getHours(),
                                      scheduleDate.getMinutes(),
                                      0,
                                      0
                                    );
                                  }
                                  field.onChange(
                                    String(Math.floor(date.getTime() / 1000))
                                  );
                                }
                              }}
                            />
                          </div>
                          <div>
                            <Form.Label required>
                              {t('form:feature-flags.update-time')}
                            </Form.Label>
                            <ReactDatePicker
                              dateFormat="HH:mm"
                              timeFormat="HH:mm"
                              selected={scheduleDate}
                              showTimeSelectOnly={true}
                              className="w-[124px]"
                              onChange={date => {
                                if (date) {
                                  field.onChange(
                                    String(Math.floor(date.getTime() / 1000))
                                  );
                                }
                              }}
                              icon={
                                <Icon
                                  icon={IconWatch}
                                  className="flex-center"
                                />
                              }
                            />
                          </div>
                        </div>
                      </Form.Control>
                      <Form.Message />
                    </Form.Item>
                  );
                }}
              />
            )}
          </>
        )}
      </div>

      <ButtonBar
        secondaryButton={
          <Button
            disabled={!isValid}
            loading={isSubmitting}
            onClick={handleSubmit(() =>
              onSubmit(watch('scheduleType'), watch('scheduleAt'))
            )}
          >
            {t(`submit`)}
          </Button>
        }
        primaryButton={
          <Button onClick={onClose} variant="secondary">
            {t(`cancel`)}
          </Button>
        }
      />
    </DialogModal>
  );
};

export default SaveWithCommentModal;
