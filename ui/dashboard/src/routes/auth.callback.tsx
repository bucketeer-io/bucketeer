import { createFileRoute } from '@tanstack/react-router';
import { AuthCallbackPage } from 'auth';
import { z } from 'zod';

export const Route = createFileRoute('/auth/callback')({
  validateSearch: z.object({
    state: z.number(),
    code: z.string(),
    scope: z.string(),
    authuser: z.number(),
    prompt: z.string()
  }),
  component: AuthCallbackPage
});
