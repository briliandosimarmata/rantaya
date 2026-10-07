import { api, query, requireUser } from '#lib/server/api.js';
import { error } from '@sveltejs/kit';
export const load = async (e: any) => {
  await requireUser(e, 'admin');
  return { reports: await api(e, '/admin/reports') };
};
