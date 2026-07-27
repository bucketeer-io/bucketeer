import { createContext, ReactNode, useContext } from 'react';
import { useWalkthrough } from 'hooks/use-walkthrough';
import ConnectSdkModal from 'elements/connect-sdk-modal';

interface WalkthroughContextType {
  startWalkthrough: (options?: { withWelcome?: boolean }) => void;
}

const WalkthroughContext = createContext<WalkthroughContextType | undefined>(
  undefined
);

export const WalkthroughProvider = ({ children }: { children: ReactNode }) => {
  const { startWalkthrough, closeWalkthrough, createdFlagId, isSdkModalOpen } =
    useWalkthrough();

  return (
    <WalkthroughContext.Provider value={{ startWalkthrough }}>
      {children}
      {isSdkModalOpen && (
        <ConnectSdkModal
          isOpen={isSdkModalOpen}
          flagId={createdFlagId}
          onClose={closeWalkthrough}
        />
      )}
    </WalkthroughContext.Provider>
  );
};

export const useWalkthroughContext = (): WalkthroughContextType => {
  const context = useContext(WalkthroughContext);
  if (!context) {
    throw new Error(
      'useWalkthroughContext must be used within WalkthroughProvider'
    );
  }
  return context;
};
