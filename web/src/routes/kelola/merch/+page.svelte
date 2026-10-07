<script lang="ts">
  import { money } from '#lib/format.js';
  import Empty from '#lib/components/Empty.svelte';
  let { data } = $props();
</script>
<div class="between" style="margin-bottom:20px"><h2>Merchandise</h2><a class="button primary" href="/kelola/merch/baru">Tambah produk</a></div>
{#if data.products.length}<div class="card table-wrap"><table class="table"><thead><tr><th>Produk</th><th>Harga</th><th>Ketersediaan</th><th>Tindakan</th></tr></thead><tbody>
{#each data.products as p}<tr><td><a href={'/merch/' + p.id}><b>{p.name}</b></a></td><td>{money(p.price)}</td><td>{p.availability}</td><td><a class="button slim" href={'/kelola/merch/' + p.id}>Edit</a></td></tr>{/each}
</tbody></table></div>{:else}<Empty title="Belum ada merchandise" body="Tambahkan produk dari ruang komunitasmu." />{/if}
