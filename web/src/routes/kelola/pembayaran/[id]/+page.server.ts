import { api, query, requireUser } from '#lib/server/api.js';
import { error } from '@sveltejs/kit';
export const load = async (e: any) => {
  return { order: await api(e, '/orders/' + e.params.id) };
};
