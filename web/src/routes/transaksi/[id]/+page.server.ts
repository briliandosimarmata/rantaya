import { api, query, requireUser } from '#lib/server/api.js';
import { error } from '@sveltejs/kit';
export const load = async (e: any) => {
  await requireUser(e, 'customer');
  const order = await api(e, '/orders/' + e.params.id);
  return {
    order,
    methods: await api(e, '/payment-methods' + query({ organizer: order.event.organizer_id })),
  };
};
