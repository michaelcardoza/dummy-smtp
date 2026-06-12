<script lang="ts">
  import Icon from "./Icon.svelte";
  import { toasts } from '../stores/toasts.svelte'
</script>

<div class="fixed bottom-12 right-5 z-50 flex flex-col gap-2">
  {#each toasts.items as t (t.id)}
    <div class="flex items-center gap-2 border-2 bg-neutral-900 px-4 py-2 text-sm text-neutral-200 {t.error
      ? 'border-red-900/60'
      : 'border-accent/40'}">
      <Icon name={t.error ? 'times' : 'bell'} class="glow {t.error ? 'text-red-400' : 'text-accent'}" />
      {t.text}
    </div>
  {/each}
</div>
