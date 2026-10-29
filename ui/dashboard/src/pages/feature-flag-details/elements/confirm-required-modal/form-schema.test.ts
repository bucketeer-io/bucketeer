import { describe, expect, it, vi } from 'vitest';
import {
  formSchema,
  SCHEDULE_TYPE_SCHEDULE,
  SCHEDULE_TYPE_UPDATE_NOW
} from './form-schema';

vi.mock('i18n', () => ({
  i18n: { t: (key: string) => key }
}));

const nowInSeconds = Math.floor(Date.now() / 1000);
const inOneHour = String(nowInSeconds + 3600);
const oneHourAgo = String(nowInSeconds - 3600);

const isValid = (values: Record<string, unknown>) =>
  formSchema.isValid({
    requireComment: false,
    resetSampling: false,
    comment: '',
    scheduleType: SCHEDULE_TYPE_UPDATE_NOW,
    scheduleAt: inOneHour,
    ...values
  });

describe('confirm-required-modal formSchema comment requirement', () => {
  it('requires a comment for an immediate update when the environment requires one', async () => {
    expect(
      await isValid({
        requireComment: true,
        scheduleType: SCHEDULE_TYPE_UPDATE_NOW
      })
    ).toBe(false);
    expect(
      await isValid({
        requireComment: true,
        scheduleType: SCHEDULE_TYPE_UPDATE_NOW,
        comment: 'enable for the launch'
      })
    ).toBe(true);
  });

  it('requires a comment for a scheduled update when the environment requires one', async () => {
    expect(
      await isValid({
        requireComment: true,
        scheduleType: SCHEDULE_TYPE_SCHEDULE
      })
    ).toBe(false);
    expect(
      await isValid({
        requireComment: true,
        scheduleType: SCHEDULE_TYPE_SCHEDULE,
        comment: 'enable for the launch'
      })
    ).toBe(true);
  });

  it('does not require a comment when the environment does not require one', async () => {
    expect(
      await isValid({
        requireComment: false,
        scheduleType: SCHEDULE_TYPE_UPDATE_NOW
      })
    ).toBe(true);
    expect(
      await isValid({
        requireComment: false,
        scheduleType: SCHEDULE_TYPE_SCHEDULE
      })
    ).toBe(true);
  });

  it('rejects a scheduled update whose time is not in the future', async () => {
    expect(
      await isValid({
        scheduleType: SCHEDULE_TYPE_SCHEDULE,
        scheduleAt: oneHourAgo
      })
    ).toBe(false);
    expect(
      await isValid({ scheduleType: SCHEDULE_TYPE_SCHEDULE, scheduleAt: '' })
    ).toBe(false);
  });
});
