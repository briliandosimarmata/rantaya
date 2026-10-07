<script lang="ts">
  import type { ArtsEvent } from '#lib/types.js';
  import { recordClick } from '#lib/api.js';
  import { ended } from '#lib/format.js';
  import Icon from './Icon.svelte';
  let {
    event,
    saved,
    bookmark,
    sharing,
    canReview = true,
  } = $props<{
    event: ArtsEvent;
    saved: boolean;
    bookmark: () => void;
    sharing: () => void;
    canReview?: boolean;
  }>();
</script>

<div class="event-actions">
  {#if !ended(event) || canReview}<a
      class="button event-primary" class:lime={!ended(event)} class:primary={ended(event)}
      href={'/event/' + event.slug + (ended(event) ? '/ulasan' : '/tiket')}
      onclick={() => { if (!ended(event)) recordClick('ticket',event.id); }}
      >{#if !ended(event)}<Icon name="ticket" />{/if}{ended(event) ? 'Tulis ulasan' : 'Beli tiket'}</a
    >{/if}
  <button
    class="icon-button event-secondary"
    class:active={saved}
    onclick={bookmark}
    aria-pressed={saved}
    aria-label={saved ? 'Batalkan simpan event' : 'Simpan event'}
    ><Icon name={saved ? 'check' : 'bookmark'} /></button
  >
  <button class="icon-button event-secondary" onclick={sharing} aria-label="Bagikan event"
    ><Icon name="share" /></button
  >
</div>
