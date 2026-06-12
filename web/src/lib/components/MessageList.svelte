<script lang="ts">
  import { messages } from '../stores/messages.svelte'
  import { ago, dayLabel } from '../utils/time'
  import { SMTP_HOST, SMTP_PORT } from '../config'
  import Icon from './Icon.svelte'

  function sizeOf(raw: string): string {
    const kb = raw.length / 1024
    return kb < 1 ? `${raw.length} B` : `${kb.toFixed(1)} KB`
  }
</script>

<nav class="w-96 shrink-0 overflow-y-auto border-r-2 border-neutral-900 bg-neutral-950">
  {#each messages.filtered as m, i (m.id)}
    {#if i === 0 || dayLabel(m.created_at) !== dayLabel(messages.filtered[i - 1].created_at)}
      <div class="border-b border-neutral-900 bg-neutral-950 px-4 py-1.5 text-xs font-bold uppercase tracking-widest text-accent glow">
        {dayLabel(m.created_at)}
      </div>
    {/if}
    <button
      onclick={() => messages.select(m.id)}
      class="block w-full cursor-pointer border-b border-l-2 border-neutral-900 px-4 py-3.5 text-left hover:bg-neutral-900 {messages.selectedId === m.id
        ? 'border-l-accent bg-neutral-900'
        : 'border-l-transparent'}"
    >
      <div class="flex items-center gap-2.5">
        {#if i === 0}
          <span class="h-2 w-2 shrink-0 bg-accent"></span>
        {/if}
        <div class="min-w-0 truncate font-semibold text-neutral-100">{m.subject || '(no subject)'}</div>
      </div>
      <div class="truncate text-sm text-neutral-500">
        <span class="text-neutral-600">from</span> {m.from}
      </div>
      <div class="truncate text-sm text-neutral-500">
        <span class="text-neutral-600">to</span> {m.to.join(', ')}
      </div>
      <div class="mt-1 flex items-center gap-2.5 text-sm text-neutral-600">
        <span>{sizeOf(m.raw)}</span>
        {#if m.attachments?.length}
          <span class="flex items-center gap-1">
            <Icon name="paperclip" />{m.attachments.length}
          </span>
        {/if}
        <span class="ml-auto shrink-0" title={m.created_at}>{ago(m.created_at)}</span>
      </div>
    </button>
  {:else}
    <div class="p-12 text-center leading-relaxed text-neutral-500">
      <Icon name="envelope" class="mb-4 block text-4xl text-neutral-700" />
      <div class="mb-3 text-neutral-400">No messages yet<span class="blink text-accent">_</span></div>
      <div class="text-sm">Point your app's SMTP at:</div>
      <div class="mt-4 inline-block border-2 border-neutral-900 bg-neutral-950 px-4 py-3 text-left text-sm text-neutral-400">
        host: <span class="text-accent glow">{SMTP_HOST}</span><br />
        port: <span class="text-accent glow">{SMTP_PORT}</span>
      </div>
    </div>
  {/each}
</nav>
