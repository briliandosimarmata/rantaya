<script lang="ts">
  import ProfileForm from '#lib/components/ProfileForm.svelte';
  let { data } = $props();
  import SEO from '#lib/components/SEO.svelte';
</script>

<SEO title="Kenalan dulu" privatePage />
<div class="page-heading">
  <h1>Kenalan dulu, yuk.</h1>
  <p>Lengkapi profil supaya komunitas bisa mengenalmu. Boleh dilewati dan dilengkapi nanti.</p>
</div>
<section class="narrow card pad"><ProfileForm onboarding next={data.next} /></section>
