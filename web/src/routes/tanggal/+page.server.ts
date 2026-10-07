import { api, query, requireUser } from '#lib/server/api.js';
import { error } from '@sveltejs/kit';
export const load = async (e: any) => {
  return {
    from: e.url.searchParams.get('from') || '',
    to: e.url.searchParams.get('to') || '',
    category: e.url.searchParams.get('category') || '',
    period: e.url.searchParams.get('period') || 'upcoming',
  };
};
