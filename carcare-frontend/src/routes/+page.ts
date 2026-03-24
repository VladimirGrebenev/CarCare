import { redirect } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
  // Redirect to welcome page for first-time visitors
  // This ensures unauthenticated users always see the welcome flow
  throw redirect(302, '/welcome');
};
