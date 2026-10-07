<script lang="ts">
  import { initials } from '#lib/format.js';
  let { name, src = '', size = '' }: { name: string; src?: string; size?: string } = $props();
</script>

<span class="avatar {size}"
  >{#if src}<img {src} alt="" loading="lazy" width="48" height="48" />{:else}{initials(
      name,
    )}{/if}</span
>
