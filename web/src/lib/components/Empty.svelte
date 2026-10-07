<script lang="ts">
  import Icon from './Icon.svelte';
  let { title, body = '', href = '', label = 'Kembali ke beranda' } = $props();
</script>

<div class="empty">
  <Icon name="chat" size={32} />
  <h2>{title}</h2>
  <p>{body}</p>
  {#if href}<a class="button primary" {href}>{label}</a>{/if}
</div>
