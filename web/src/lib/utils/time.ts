import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'

dayjs.extend(relativeTime)

export function ago(iso: string): string {
  return dayjs(iso).fromNow()
}

export function datetime(iso: string): string {
  return dayjs(iso).format('MMM D, YYYY HH:mm')
}

export function fileStamp(iso: string): string {
  return dayjs(iso).format('YYYY-MM-DD_HHmm')
}

export function dayLabel(iso: string): string {
  const d = dayjs(iso)
  if (d.isSame(dayjs(), 'day')) return 'Today'
  if (d.isSame(dayjs().subtract(1, 'day'), 'day')) return 'Yesterday'
  return d.format('MMM D, YYYY')
}
