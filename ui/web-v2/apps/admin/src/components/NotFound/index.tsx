import { PAGE_PATH_ROOT } from '@/constants/routing';
import { messages } from '@/lang/messages';
import { FC } from 'react';
import { useIntl } from 'react-intl';
import { Link } from 'react-router-dom';

export const NotFound: FC = () => {
  const { formatMessage: f } = useIntl();

  return (
    <main className="grid min-h-full place-items-center py-24 px-6 sm:py-32 lg:px-8">
      <div className="text-center">
        <p className="text-base font-semibold text-primary">404</p>
        <h1 className="mt-4 text-3xl font-bold tracking-tight sm:text-5xl">
          {f(messages.notFound.title)}
        </h1>
        <p className="mt-6 text-base leading-7 text-gray-600">
          {f(messages.notFound.description)}
        </p>
        <Link to={PAGE_PATH_ROOT} className="btn-submit mt-10">
          {f(messages.notFound.goBackHome)}
        </Link>
      </div>
    </main>
  );
};
