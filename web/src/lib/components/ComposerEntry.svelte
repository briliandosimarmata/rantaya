<script lang="ts">
  import { getContext } from 'svelte';
  import type { AppContext } from '#lib/types.js';
  import Avatar from './Avatar.svelte';
  import Icon from './Icon.svelte';
  let { href = '/tulis' } = $props<{ href?: string }>();
  const ctx = getContext<AppContext>('app');
</script>
<div class="card compose-entry">
  <Avatar name={ctx.user?.name || '?'} src={ctx.user?.avatar_url || ''} size={'lime ' + (ctx.user?.role === 'organizer' ? 'square' : '')} />
  <a class="composer-start" {href}>Mau cerita apa hari ini?</a>
  <a class="icon-button" {href} aria-label="Buat postingan"><Icon name="edit" /></a>
</div>
