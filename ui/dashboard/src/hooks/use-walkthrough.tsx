import { useCallback, useEffect, useRef, useState } from 'react';
import { useLocation, useNavigate } from 'react-router';
import { apiKeysFetcher } from '@api/api-key';
import { getCurrentEnvironment, useAuth } from 'auth';
import { WALKTHROUGH_ENABLED } from 'configs';
import { PAGE_PATH_APIKEYS, PAGE_PATH_FEATURES } from 'constants/routing';
import { WALKTHROUGH_TARGETS } from 'constants/walkthrough';
import { driver, type Config, type Driver, type DriveStep } from 'driver.js';
import 'driver.js/dist/driver.css';
import { useTranslation } from 'i18n';
import {
  clearWalkthroughPendingStorage,
  getWalkthroughPendingStorage
} from 'storage/walkthrough';

type WalkthroughStage =
  | 'idle'
  | 'flag-tour'
  | 'await-flag-created'
  | 'apikey-tour'
  | 'await-apikey-created'
  | 'sdk-modal';

const tourTarget = (id: string) => `[data-tour="${id}"]`;

const TARGETING_PATH_REGEX = /\/features\/([^/]+)\/targeting/;

// Submitting is reserved for the dedicated submit step; textarea keeps Enter.
const blockEnterSubmit = (event: KeyboardEvent) => {
  const target = event.target as HTMLElement | null;
  if (
    event.key === 'Enter' &&
    target?.tagName !== 'TEXTAREA' &&
    target?.closest(tourTarget(WALKTHROUGH_TARGETS.APIKEY_FORM))
  ) {
    event.preventDefault();
  }
};
const AWAIT_APIKEY_TIMEOUT = 10 * 60 * 1000;
const AWAIT_APIKEY_POLL_INTERVAL = 2000;

const hasUsableSDKKey = async (
  organizationId: string,
  environmentId: string
) => {
  try {
    const collection = await apiKeysFetcher({
      cursor: String(0),
      organizationId,
      environmentIds: [environmentId]
    });
    return (
      collection?.apiKeys?.some(
        key => ['SDK_CLIENT', 'SDK_SERVER'].includes(key.role) && !key.disabled
      ) ?? false
    );
  } catch {
    return false;
  }
};

export const useWalkthrough = () => {
  const { t } = useTranslation(['common']);
  const navigate = useNavigate();
  const { pathname } = useLocation();
  const { consoleAccount } = useAuth();
  const driverRef = useRef<Driver | null>(null);
  const cancelledRef = useRef(false);
  const [stage, setStage] = useState<WalkthroughStage>('idle');
  const [createdFlagId, setCreatedFlagId] = useState('');

  useEffect(
    () => () => {
      driverRef.current?.destroy();
      driverRef.current = null;
    },
    []
  );

  const closeWalkthrough = useCallback(() => {
    driverRef.current?.destroy();
    setStage('idle');
    setCreatedFlagId('');
  }, []);

  const buildDriverConfig = useCallback(
    (
      tourStage: WalkthroughStage,
      lastStepTarget: string,
      excludeSkipTargets: string[] = []
    ): Config => ({
      showProgress: true,
      overlayOpacity: 0.6,
      stagePadding: 6,
      stageRadius: 8,
      // Steps may target elements on a page the previous step navigates to.
      waitForElement: 5000,
      overlayClickBehavior: () => undefined,
      allowKeyboardControl: false,
      nextBtnText: t('walkthrough.next'),
      doneBtnText: t('walkthrough.done'),
      onPopoverRender: (popover, { state }) => {
        const activeElement = state.activeStep?.element;
        if (
          typeof activeElement === 'string' &&
          excludeSkipTargets.includes(activeElement)
        ) {
          return;
        }
        const skipButton = document.createElement('button');
        skipButton.type = 'button';
        skipButton.classList.add(
          'driver-popover-footer-btn',
          'driver-popover-skip-btn'
        );
        skipButton.innerText = t('walkthrough.skip');
        skipButton.addEventListener('click', () => {
          cancelledRef.current = true;
          driverRef.current?.destroy();
        });
        popover.footerButtons.insertBefore(skipButton, popover.nextButton);
      },
      onCloseClick: (_element, _step, { driver: driverInstance }) => {
        cancelledRef.current = true;
        driverInstance.destroy();
      },
      onDestroyed: (_element, step) => {
        document.removeEventListener('keydown', blockEnterSubmit, true);
        // Ending on the last step's action click hands off to the next stage.
        if (!cancelledRef.current && step?.element === lastStepTarget) {
          setStage(
            tourStage === 'flag-tour'
              ? 'await-flag-created'
              : 'await-apikey-created'
          );
        } else {
          setStage(current => (current === tourStage ? 'idle' : current));
        }
      }
    }),
    [t]
  );

  const startApiKeyTour = useCallback(() => {
    if (!consoleAccount) return;
    const currentEnvironment = getCurrentEnvironment(consoleAccount);

    cancelledRef.current = false;
    setStage('apikey-tour');
    driverRef.current?.destroy();
    driverRef.current = driver({
      ...buildDriverConfig(
        'apikey-tour',
        tourTarget(WALKTHROUGH_TARGETS.SUBMIT_APIKEY_BUTTON),
        [tourTarget(WALKTHROUGH_TARGETS.SUBMIT_APIKEY_BUTTON)]
      ),
      steps: [
        {
          // No element: shown as a centered dialog.
          popover: {
            title: t('walkthrough.flag-created.title'),
            description: t('walkthrough.flag-created.description'),
            showButtons: ['next', 'close'],
            onNextClick: (_element, _step, { driver: driverInstance }) => {
              navigate(`/${currentEnvironment.urlCode}${PAGE_PATH_APIKEYS}`);
              driverInstance.moveNext();
            }
          }
        },
        {
          element: tourTarget(WALKTHROUGH_TARGETS.CREATE_APIKEY_BUTTON),
          advanceOnClick: true,
          popover: {
            title: t('walkthrough.create-apikey-button.title'),
            description: t('walkthrough.create-apikey-button.description'),
            side: 'bottom',
            showButtons: ['close']
          }
        },
        {
          element: tourTarget(WALKTHROUGH_TARGETS.APIKEY_FORM),
          disableActiveInteraction: false,
          onHighlighted: () => {
            document.addEventListener('keydown', blockEnterSubmit, true);
          },
          onDeselected: () => {
            document.removeEventListener('keydown', blockEnterSubmit, true);
          },
          popover: {
            title: t('walkthrough.apikey-form.title'),
            description: t('walkthrough.apikey-form.description'),
            side: 'left',
            showButtons: ['next', 'close']
          }
        },
        {
          element: tourTarget(WALKTHROUGH_TARGETS.SUBMIT_APIKEY_BUTTON),
          advanceOnClick: true,
          popover: {
            title: t('walkthrough.submit-apikey-button.title'),
            description: t('walkthrough.submit-apikey-button.description'),
            side: 'top',
            showButtons: ['previous', 'close']
          }
        }
      ]
    });
    driverRef.current.drive();
  }, [buildDriverConfig, consoleAccount, navigate, t]);

  // The created flag id comes from the targeting page URL after submit.
  useEffect(() => {
    if (stage !== 'await-flag-created') return;
    const match = pathname.match(TARGETING_PATH_REGEX);
    if (!match) return;

    setCreatedFlagId(match[1]);
    startApiKeyTour();
  }, [pathname, stage, startApiKeyTour]);

  // Poll until the key exists and the "API key created" dialog is closed.
  useEffect(() => {
    if (stage !== 'await-apikey-created' || !consoleAccount) return;
    const currentEnvironment = getCurrentEnvironment(consoleAccount);
    const startedAt = Date.now();
    let checking = false;
    const interval = setInterval(async () => {
      if (Date.now() - startedAt > AWAIT_APIKEY_TIMEOUT) {
        setStage('idle');
        return;
      }
      if (checking) return;
      checking = true;
      try {
        const hasKey = await hasUsableSDKKey(
          currentEnvironment.organizationId,
          currentEnvironment.id
        );
        const isDialogOpen = !!document.querySelector('[role="dialog"]');
        if (hasKey && !isDialogOpen) {
          setStage('sdk-modal');
        }
      } finally {
        checking = false;
      }
    }, AWAIT_APIKEY_POLL_INTERVAL);
    return () => clearInterval(interval);
  }, [stage, consoleAccount]);

  const startWalkthrough = useCallback(
    (options?: { withWelcome?: boolean }) => {
      if (!consoleAccount) return;
      const currentEnvironment = getCurrentEnvironment(consoleAccount);
      navigate(`/${currentEnvironment.urlCode}${PAGE_PATH_FEATURES}`);

      cancelledRef.current = false;
      setStage('flag-tour');
      driverRef.current?.destroy();
      const welcomeStep: DriveStep = {
        // No element: shown as a centered dialog.
        popover: {
          title: t('walkthrough.welcome.title'),
          description: t('walkthrough.welcome.description'),
          nextBtnText: t('walkthrough.welcome.start'),
          showButtons: ['next', 'close']
        }
      };
      driverRef.current = driver({
        ...buildDriverConfig(
          'flag-tour',
          tourTarget(WALKTHROUGH_TARGETS.SUBMIT_FLAG_BUTTON)
        ),
        steps: [
          ...(options?.withWelcome ? [welcomeStep] : []),
          {
            element: tourTarget(WALKTHROUGH_TARGETS.FEATURE_FLAGS_MENU),
            disableActiveInteraction: true,
            popover: {
              title: t('walkthrough.feature-flags-menu.title'),
              description: t('walkthrough.feature-flags-menu.description'),
              side: 'right',
              showButtons: ['next', 'close']
            }
          },
          {
            element: tourTarget(WALKTHROUGH_TARGETS.CREATE_FLAG_BUTTON),
            advanceOnClick: true,
            popover: {
              title: t('walkthrough.create-flag-button.title'),
              description: t('walkthrough.create-flag-button.description'),
              side: 'bottom',
              showButtons: ['close']
            }
          },
          {
            element: tourTarget(WALKTHROUGH_TARGETS.FLAG_GENERAL_INFO),
            disableActiveInteraction: false,
            popover: {
              title: t('walkthrough.flag-general-info.title'),
              description: t('walkthrough.flag-general-info.description'),
              side: 'right',
              showButtons: ['next', 'close']
            }
          },
          {
            element: tourTarget(WALKTHROUGH_TARGETS.FLAG_VARIATIONS),
            disableActiveInteraction: false,
            popover: {
              title: t('walkthrough.flag-variations.title'),
              description: t('walkthrough.flag-variations.description'),
              side: 'right',
              showButtons: ['next', 'previous', 'close']
            }
          },
          {
            element: tourTarget(WALKTHROUGH_TARGETS.SUBMIT_FLAG_BUTTON),
            advanceOnClick: true,
            popover: {
              title: t('walkthrough.submit-flag-button.title'),
              description: t('walkthrough.submit-flag-button.description'),
              side: 'top',
              showButtons: ['previous', 'close']
            }
          }
        ]
      });
      driverRef.current.drive();
    },
    [buildDriverConfig, consoleAccount, navigate, t]
  );

  // Auto-start once for first-time users (flag set at login).
  const hasAutoStartedRef = useRef(false);
  useEffect(() => {
    if (hasAutoStartedRef.current || !WALKTHROUGH_ENABLED || !consoleAccount)
      return;
    if (getWalkthroughPendingStorage()) {
      hasAutoStartedRef.current = true;
      clearWalkthroughPendingStorage();
      startWalkthrough({ withWelcome: true });
    }
  }, [startWalkthrough, consoleAccount]);

  return {
    startWalkthrough,
    closeWalkthrough,
    createdFlagId,
    isSdkModalOpen: stage === 'sdk-modal'
  };
};
