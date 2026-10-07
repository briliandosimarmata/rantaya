import { api, query, requireUser } from '#lib/server/api.js';
export const load = async (e: any) => {
  await requireUser(e);
  const { city, organizers } = await e.parent();
  const [items, nearby] = await Promise.all([api(e, '/bookmarks'), api(e, '/events' + query({ city, period: 'upcoming' }))]);
  const eventIDs = [...new Set(items.filter((x: any)=>x.kind==='post' && x.item).flatMap((x: any)=>x.item.mentions.filter((m: any)=>m.kind==='event').map((m: any)=>m.id)))];
  const relatedEvents = await Promise.all(eventIDs.map(id=>api(e,'/events/'+id).catch((err: any)=>{if([403,404].includes(err.status))return null;throw err;})));
  return { items, nearby, organizers, relatedEvents: relatedEvents.filter(Boolean) };
};
