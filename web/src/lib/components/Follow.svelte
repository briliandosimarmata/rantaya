<script lang="ts">
  import { getContext } from 'svelte';
  import type { AppContext } from '#lib/types.js';
  import { mutate } from '#lib/api.js';
  import Icon from './Icon.svelte';
  let { id, small = false } = $props<{ id: string; small?: boolean }>();
  const ctx = getContext<AppContext>('app');
  let busy = $state(false);
  let active = $derived(ctx.account.follows.includes(id));
  async function toggle() {
    if (!ctx.requireLogin()) return;
    busy = true;
    try {
      await mutate('/organizers/' + id + '/follow', { active: !active });
      await ctx.refresh();
    } catch (e: any) {
      ctx.toast(e.message);
    } finally {
      busy = false;
    }
  }
</script>

<button
  class:primary={!active}
  class={small ? 'follow-mini' : 'button'}
  class:followed={active}
  onclick={toggle}
  disabled={busy}
  aria-pressed={active}>{#if !small}<Icon name={active ? 'check' : 'plus'} />{/if}{active ? 'Diikuti' : 'Ikuti'}</button
>
