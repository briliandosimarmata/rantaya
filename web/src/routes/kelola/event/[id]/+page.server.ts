import { api, query, requireUser } from '#lib/server/api.js';
import { error } from '@sveltejs/kit';
export const load = async (e: any) => {
  return { event: e.params.id === 'baru' ? null : await api(e, '/events/' + e.params.id) };
};
