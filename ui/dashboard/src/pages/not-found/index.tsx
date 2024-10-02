import { useNavigate } from 'react-router-dom';
import { PAGE_PATH_ROOT } from 'constants/routing';
import Button from 'components/button';

const NotFoundPage = () => {
  const navigate = useNavigate();
  return (
    <main className="grid min-h-full place-items-center py-24 px-6 sm:py-32 lg:px-8">
      <div className="text-center">
        <p className="text-xl font-semibold text-primary-500">{`404`}</p>
        <h1 className="mt-4 text-3xl font-bold tracking-tight sm:text-5xl">
          {`Page not found`}
        </h1>
        <p className="mt-6 text-base leading-7 text-gray-600">
          {`Sorry, we couldn't find the page you're looking for.`}
        </p>
        <Button
          onClick={() => navigate(PAGE_PATH_ROOT, { replace: true })}
          className="mt-8"
        >
          {`Go back home`}
        </Button>
      </div>
    </main>
  );
};

export default NotFoundPage;
