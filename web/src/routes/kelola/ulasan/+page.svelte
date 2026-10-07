<script lang="ts">
  import ReviewCard from '#lib/components/ReviewCard.svelte';
  import Empty from '#lib/components/Empty.svelte';
  let { data } = $props();
</script>
<div class="note">Ulasan terhubung ke event asalnya. Pengelola bisa membalas atau melaporkan.</div>
{#each data.reviews as review}<ReviewCard {review} />{:else}<Empty title="Belum ada ulasan" body="Ulasan penonton akan muncul setelah event selesai." />{/each}
