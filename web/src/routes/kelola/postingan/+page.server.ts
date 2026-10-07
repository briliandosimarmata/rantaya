import { api, query, requireUser } from '#lib/server/api.js';
import { error } from '@sveltejs/kit';
export const load = async (e: any) => {
  const { user } = await e.parent();
  return { posts: await api(e, '/posts' + query({ organizer: user.organizer_id })) };
};
