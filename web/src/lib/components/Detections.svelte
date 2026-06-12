<script lang="ts">
  import Icon from "./Icon.svelte";
  import type { Message } from '../types'
  import { extractCodes, extractLinks } from '../utils/extract'
  import { copy } from '../utils/clipboard'

  let { message }: { message: Message } = $props()

  let codes = $derived(extractCodes(message))
  let links = $derived(extractLinks(message))
</script>

{#if codes.length || links.length}
  <div class="space-y-2 border-b-2 border-neutral-900 bg-neutral-950 px-8 py-3">
    {#if codes.length}
      <div class="flex flex-wrap items-center gap-2">
        <span class="text-xs font-bold uppercase tracking-widest text-neutral-600">
          <span class="text-accent glow">&gt;</span> codes
        </span>
        {#each codes as c}
          <button
            onclick={() => copy(c)}
            class="flex items-center gap-1.5 border-2 border-neutral-900 px-2.5 py-1 text-sm text-accent glow hover:border-accent"
            title="copy"
          >
            <Icon name="hashtag" />{c}
          </button>
        {/each}
      </div>
    {/if}
    {#if links.length}
      <div class="flex flex-wrap items-center gap-2">
        <span class="text-xs font-bold uppercase tracking-widest text-neutral-600">
          <span class="text-accent glow">&gt;</span> links
        </span>
        {#each links as l}
          <span class="flex items-center gap-2 border-2 border-neutral-900 px-2.5 py-1 text-sm text-neutral-300">
            <Icon name="link" class="text-neutral-600" />
            <span class="max-w-xs truncate">{l}</span>
            <button onclick={() => copy(l)} aria-label="copy" class="text-neutral-600 hover:text-accent">
              <Icon name="copy" />
            </button>
            <a href={l} target="_blank" rel="noreferrer" aria-label="open" class="text-neutral-600 hover:text-accent">
              <Icon name="external-link" />
            </a>
          </span>
        {/each}
      </div>
    {/if}
  </div>
{/if}
