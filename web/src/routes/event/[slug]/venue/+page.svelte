<script lang="ts">
  import SEO from '#lib/components/SEO.svelte';
  let { data } = $props();
  let zoom = $state(100);
</script>
<SEO title={'Layout venue ' + data.event.title} />
<section class="wide-page"><div class="venue-toolbar"><div><h1>Layout venue</h1><p class="meta">{data.event.title} · {data.event.venue}</p></div><div class="row"><button class="icon-button" onclick={() => zoom = Math.max(100,zoom - 50)} aria-label="Perkecil layout">−</button><button class="icon-button" onclick={() => zoom = Math.min(300,zoom + 50)} aria-label="Perbesar layout">+</button><button class="button slim" onclick={() => zoom = 100}>Reset</button></div></div>
  {#if data.event.layout_url}<div class="venue-canvas"><img src={data.event.layout_url} alt={'Layout venue dan tempat duduk ' + data.event.title} style={'width:' + zoom + '%'} /></div>{:else}<p>Pengelola belum mengunggah layout venue.</p>{/if}
  <p class="meta" style="margin-top:14px">Perbesar lalu geser untuk melihat detail. Layout ini panduan venue; tidak memilih kursi.</p>
</section>
