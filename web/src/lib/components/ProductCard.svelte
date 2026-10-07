<script lang="ts">
  import { money } from '#lib/format.js';
  let { product: p } = $props<{ product: any }>();
</script>
<article class="card product-card">
  <a href={'/merch/' + p.id}><img class="product-image" src={p.image_url} alt={p.name + ', foto produk contoh'} loading="lazy" width="480" height="360" /></a>
  <div class="product-body">
    <span class="pill">{p.availability}</span>
    <a href={'/merch/' + p.id}><h3>{p.name}</h3></a>
    <a class="meta" href={'/ruang/' + p.organizer.slug}>{p.organizer.name}</a>
    <div class="between" style="margin-top:12px"><span class="price">{money(p.price)}</span></div>
    <a class="button full" href={'/merch/' + p.id}>Lihat produk</a>
  </div>
</article>
