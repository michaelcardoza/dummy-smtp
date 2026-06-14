<script lang="ts">
  import type { Message } from '../types'
  import { messages } from '../stores/messages.svelte'
  import { ago, datetime } from '../utils/time'
  import { copy } from '../utils/clipboard'
  import { download } from '../utils/download'
  import { downloadBaseName } from '../utils/filename'
  import { formatHtml } from '../utils/html'
  import { highlight } from '../utils/highlight'
  import Detections from './Detections.svelte'
  import Icon from './Icon.svelte'
  import Button from './Button.svelte'
  import { confirm } from '../stores/confirm.svelte'

  function confirmDelete() {
    confirm.ask({
      title: 'Delete message',
      message: "Delete this message? This can't be undone.",
      confirmLabel: 'delete',
      onconfirm: () => messages.remove(message.id),
    })
  }

  let { message }: { message: Message } = $props()

  let downloadName = $derived(downloadBaseName(message))
  let highlightedSource = $derived(
    message.htmlBody ? highlight(formatHtml(message.htmlBody), 'xml') : '',
  )
  let highlightedJson = $derived(highlight(JSON.stringify(message, null, 2), 'json'))
  let rawHeaders = $derived(message.raw.split(/\r?\n\r?\n/)[0] ?? '')

  type Tab = 'html' | 'text' | 'source' | 'raw' | 'headers' | 'json'
  let activeTab = $state<Tab>('html')

  const tabs: [Tab, string][] = [
    ['html', 'HTML'],
    ['text', 'Text'],
    ['source', 'Source'],
    ['raw', 'Raw'],
    ['headers', 'Headers'],
    ['json', 'JSON'],
  ]

  type PreviewWidth = 'mobile' | 'tablet' | 'desktop'
  let previewWidth = $state<PreviewWidth>('desktop')
  let previewExpanded = $state(false)
  let previewBackground = $state<'light' | 'dark'>('light')

  const widthOptions: [PreviewWidth, string][] = [
    ['mobile', 'Mobile'],
    ['tablet', 'Tablet'],
    ['desktop', 'Desktop'],
  ]
  const widthClasses: Record<PreviewWidth, string> = {
    mobile: 'max-w-sm',
    tablet: 'max-w-3xl',
    desktop: 'max-w-full',
  }

  // Default to the HTML tab when a message has an HTML part, otherwise Text.
  $effect(() => {
    message.id
    activeTab = message.htmlBody ? 'html' : 'text'
  })
</script>

{#snippet copyButton(value: string)}
  <button onclick={() => copy(value)} class="ml-2 px-1.5 flex items-center text-sm text-neutral-600 hover:text-accent" title="copy">
    <Icon name="copy" class="text-base" />
  </button>
{/snippet}

<article class="flex h-full flex-col">
  {#if !(previewExpanded && activeTab === 'html')}
  <header class="border-b-2 border-neutral-900 bg-neutral-950 px-8 py-6">
    <div class="flex items-start justify-between gap-4">
      <div class="min-w-0">
        <h2 class="mb-1.5 truncate text-xl font-semibold text-neutral-100">
          {message.subject || '(no subject)'}
        </h2>
        <div class="mb-3 flex items-center gap-1.5 text-xs text-neutral-600">
          <span class="uppercase tracking-wider">id</span>
          <span class="select-all text-neutral-500">{message.id}</span>
          {@render copyButton(message.id)}
        </div>
        <dl class="space-y-1 text-neutral-500">
          <div class="flex items-center">
            <dt class="inline-block w-16 font-semibold text-neutral-300">from</dt>
            {message.from}{@render copyButton(message.from)}
          </div>
          <div class="flex items-center">
            <dt class="inline-block w-16 font-semibold text-neutral-300">to</dt>
            {message.to.join(', ')}{@render copyButton(message.to.join(', '))}
          </div>
          <div class="flex items-center">
            <dt class="inline-block w-16 font-semibold text-neutral-300">cc</dt>
            <span class="text-neutral-600">(none)</span>
          </div>
          <div class="flex items-center">
            <dt class="inline-block w-16 font-semibold text-neutral-300">reply</dt>
            <span class="text-neutral-600">{message.from}</span>
          </div>
          <div class="flex items-center">
            <dt class="inline-block w-16 font-semibold text-neutral-300">date</dt>
            {datetime(message.createdAt)}
            <span class="text-neutral-600">· {ago(message.createdAt)}</span>
          </div>
        </dl>
      </div>
      <div class="flex shrink-0 gap-2">
        <Button
          icon="download"
          onclick={() => download(`${downloadName}.eml`, message.raw, 'message/rfc822')}
        >
          .eml
        </Button>
        <Button
          icon="code"
          onclick={() =>
            download(`${downloadName}.json`, JSON.stringify(message, null, 2), 'application/json')}
        >
          .json
        </Button>
        <Button variant="danger" icon="trash" onclick={confirmDelete}>delete</Button>
      </div>
    </div>

    {#if message.attachments?.length}
      <div class="mt-4 flex flex-wrap gap-2">
        {#each message.attachments as attachment}
          <button class="flex items-center gap-2 border-2 border-neutral-900 bg-neutral-950 px-3 py-1 text-sm text-neutral-300 hover:border-accent hover:text-accent">
            <Icon name="paperclip" />
            {attachment.filename} <span class="text-neutral-500">{attachment.content_type}</span>
          </button>
        {/each}
      </div>
    {/if}
  </header>

  <Detections {message} />
  {/if}

  <div class="flex gap-2 border-b-2 border-neutral-900 bg-neutral-950 px-8 text-sm">
    {#each tabs as [id, label]}
      <button
        onclick={() => (activeTab = id)}
        class="-mb-px border-b-2 px-4 py-4 {activeTab === id
          ? 'border-accent text-accent glow'
          : 'border-transparent text-neutral-500 hover:text-neutral-300'}"
      >
        {label}
      </button>
    {/each}
    {#if activeTab === 'html' && message.htmlBody}
      <div class="ml-auto flex items-center gap-3 self-center">
        <div class="flex h-9 border-2 border-neutral-900">
          {#each widthOptions as [width, label], i}
            <button
              onclick={() => (previewWidth = width)}
              class="flex h-full items-center gap-1.5 px-3 text-xs {previewWidth === width
                ? 'bg-accent text-black'
                : 'text-neutral-400 hover:text-neutral-100'} {i < widthOptions.length - 1
                ? 'border-r-2 border-neutral-900'
                : ''}"
            >
              <Icon name={width} class="text-lg" /> {label}
            </button>
          {/each}
        </div>
        <button
          onclick={() => (previewBackground = previewBackground === 'light' ? 'dark' : 'light')}
          title={previewBackground === 'light'
            ? 'Light background — click for dark'
            : 'Dark background — click for light'}
          aria-label="toggle preview background"
          class="inline-flex h-9 w-9 items-center justify-center border-2 border-neutral-700 hover:border-accent {previewBackground ===
          'light'
            ? 'bg-white text-black'
            : 'bg-neutral-900 text-neutral-400'}"
        >
          <Icon name={previewBackground === 'light' ? 'lightbulb' : 'lightbulb-off'} class="text-base" />
        </button>
        <button
          onclick={() => (previewExpanded = !previewExpanded)}
          title={previewExpanded ? 'collapse' : 'expand'}
          class={previewExpanded ? 'text-accent glow' : 'text-neutral-500 hover:text-accent'}
        >
          <Icon name="expand" />
        </button>
      </div>
    {/if}
  </div>

  <div class="flex-1 overflow-y-auto p-8">
    {#if activeTab === 'html'}
      {#if message.htmlBody}
        <iframe
          title="html"
          sandbox=""
          srcdoc={message.htmlBody}
          class="relative z-40 mx-auto block h-full w-full {widthClasses[previewWidth]} border-2 border-neutral-900 {previewBackground === 'light' ? 'bg-white' : 'bg-neutral-900'}"
        ></iframe>
      {:else}
        <div class="text-neutral-600">(no HTML part)</div>
      {/if}
    {:else if activeTab === 'text'}
      {#if message.textBody}
        <pre class="whitespace-pre-wrap break-words leading-relaxed text-neutral-200">{message.textBody}</pre>
      {:else}
        <div class="text-neutral-600">(no text part)</div>
      {/if}
    {:else if activeTab === 'source'}
      <pre class="whitespace-pre-wrap break-words leading-relaxed text-neutral-400">{@html highlightedSource}</pre>
    {:else if activeTab === 'raw'}
      <pre class="whitespace-pre-wrap break-words leading-relaxed text-neutral-400">{message.raw}</pre>
    {:else if activeTab === 'headers'}
      <pre class="whitespace-pre-wrap break-words leading-relaxed text-neutral-400">{rawHeaders}</pre>
    {:else}
      <pre class="whitespace-pre-wrap break-words leading-relaxed text-neutral-400">{@html highlightedJson}</pre>
    {/if}
  </div>
</article>
