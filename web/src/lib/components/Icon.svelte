<script module lang="ts">
  import envelope from 'pixelarticons/svg/mail.svg?raw'
  import search from 'pixelarticons/svg/search.svg?raw'
  import code from 'pixelarticons/svg/brackets-angle.svg?raw'
  import refresh from 'pixelarticons/svg/reload.svg?raw'
  import cog from 'pixelarticons/svg/settings-cog.svg?raw'
  import times from 'pixelarticons/svg/close.svg?raw'
  import checkCircle from 'pixelarticons/svg/check.svg?raw'
  import trash from 'pixelarticons/svg/trash.svg?raw'
  import download from 'pixelarticons/svg/download.svg?raw'
  import paperclip from 'pixelarticons/svg/attachment.svg?raw'
  import copy from 'pixelarticons/svg/copy.svg?raw'
  import hashtag from 'pixelarticons/svg/hash.svg?raw'
  import link from 'pixelarticons/svg/link.svg?raw'
  import externalLink from 'pixelarticons/svg/external-link.svg?raw'
  import bell from 'pixelarticons/svg/bell.svg?raw'
  import expand from 'pixelarticons/svg/expand.svg?raw'
  import mobile from 'pixelarticons/svg/smartphone.svg?raw'
  import tablet from 'pixelarticons/svg/tablet.svg?raw'
  import desktop from 'pixelarticons/svg/monitor.svg?raw'
  import send from 'pixelarticons/svg/mail-right.svg?raw'
  import power from 'pixelarticons/svg/power.svg?raw'
  import lightbulb from 'pixelarticons/svg/lightbulb.svg?raw'
  import lightbulbOff from 'pixelarticons/svg/lightbulb-off.svg?raw'

  const github =
    '<svg viewBox="0 0 24 24" fill="currentColor"><path d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12"/></svg>'

  const icons = {
    envelope,
    search,
    code,
    refresh,
    cog,
    times,
    'check-circle': checkCircle,
    trash,
    download,
    paperclip,
    copy,
    hashtag,
    link,
    'external-link': externalLink,
    bell,
    expand,
    mobile,
    tablet,
    desktop,
    github,
    send,
    power,
    lightbulb,
    'lightbulb-off': lightbulbOff,
  } satisfies Record<string, string>

  export type IconName = keyof typeof icons
</script>

<script lang="ts">
  let { name, class: cls = '' }: { name: IconName; class?: string } = $props()
</script>

<span class="icon {cls}">{@html icons[name] ?? ''}</span>
