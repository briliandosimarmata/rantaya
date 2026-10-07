import { api, query, requireUser } from '#lib/server/api.js';
import { error } from '@sveltejs/kit';
export const load = async (e: any) => {
  return { event: await api(e, '/events/' + encodeURIComponent(e.params.slug)) };
};
