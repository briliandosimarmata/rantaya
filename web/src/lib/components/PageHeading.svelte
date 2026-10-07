<script lang="ts">
  import { getContext } from 'svelte';
  import type { AppContext } from '#lib/types.js';
  import Icon from './Icon.svelte';
  let { title, description = '', secondary = false, action }: { title: string; description?: string; secondary?: boolean; action?: import('svelte').Snippet } = $props();
  const ctx = getContext<AppContext>('app');
</script>

<div class="page-heading">
  <div class="meta"><Icon name="pin" size={14} /><span>{ctx.city} · Ruang seni lokal</span></div>
  <div class="between"><div>
    {#if secondary}<h2 class="page-title">{title}</h2>{:else}<h1>{title}</h1>{/if}
    {#if description}<p>{description}</p>{/if}
  </div>{#if action}{@render action()}{/if}</div>
</div>
