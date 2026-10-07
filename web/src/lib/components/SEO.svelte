<script lang="ts">
  import { page } from '$app/state';
  let {
    title,
    description = 'Temukan pertunjukan, ikuti komunitas, dan teruskan percakapan di Karawang.',
    image = '',
    privatePage = false,
  } = $props<{ title: string; description?: string; image?: string; privatePage?: boolean }>();
  let canonical = $derived(page.url.origin + page.url.pathname);
</script>

<svelte:head
  ><title>{title} · Rantaya</title><meta
    name="description"
    content={description.slice(0, 160)}
  /><link rel="canonical" href={canonical} /><meta property="og:title" content={title} /><meta
    property="og:description"
    content={description.slice(0, 160)}
  /><meta property="og:url" content={canonical} /><meta
    property="og:type"
    content="website"
  />{#if image}<meta
      property="og:image"
      content={new URL(image, page.url.origin).href}
    />{/if}{#if privatePage}<meta name="robots" content="noindex,nofollow" />{/if}</svelte:head
>
