import { api, query, requireUser } from '#lib/server/api.js';
import { error } from '@sveltejs/kit';
export const load = async (e: any) => {
  await requireUser(e);
  const kind = e.url.searchParams.get('kind'),
    id = e.url.searchParams.get('id');
  if (!['post', 'review'].includes(kind || '') || !id) error(400, 'Tujuan laporan tidak valid.');
  return { kind, id };
};
