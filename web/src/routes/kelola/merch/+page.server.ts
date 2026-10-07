import { api, query, requireUser } from '#lib/server/api.js';
import { error } from '@sveltejs/kit';
export const load = async (e: any) => {
  const { user } = await e.parent();
  return { products: await api(e, '/products' + query({ organizer: user.organizer_id })) };
};
