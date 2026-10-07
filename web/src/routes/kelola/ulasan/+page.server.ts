import { api, query, requireUser } from '#lib/server/api.js';
import { error } from '@sveltejs/kit';
export const load = async (e: any) => {
  const { user } = await e.parent();
  return { reviews: await api(e, '/reviews' + query({ organizer: user.organizer_id })) };
};
