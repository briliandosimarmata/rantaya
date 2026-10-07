import { api, query, requireUser } from '#lib/server/api.js';
import { error } from '@sveltejs/kit';
export const load = async (e: any) => {
  const role = e.params.role;
  if (!['customer', 'organizer', 'admin'].includes(role)) error(404, 'Halaman tidak ditemukan.');
  return {
    role,
    next: e.url.searchParams.get('next') || '',
    loginError: e.url.searchParams.get('error') || '',
  };
};
