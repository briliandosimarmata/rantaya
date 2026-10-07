import { api, query, requireUser } from '#lib/server/api.js';
import { error } from '@sveltejs/kit';
export const load = async (e: any) => {
  await requireUser(e);
  return {
    event: e.url.searchParams.get('event') || '',
    organizer: e.url.searchParams.get('organizer') || '',
  };
};
