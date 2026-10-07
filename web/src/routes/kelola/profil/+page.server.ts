import { api, query, requireUser } from '#lib/server/api.js';
import { error } from '@sveltejs/kit';
export const load = async (e: any) => {
  const { user } = await e.parent();
  return { organizer: await api(e, '/organizers/' + user.organizer_id) };
};
