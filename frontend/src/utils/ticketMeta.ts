import type { Priority, Status } from '@/types/ticket'

export type Tone = 'slate' | 'blue' | 'amber' | 'green' | 'red'

const STATUS_META: Record<Status, { label: string; tone: Tone }> = {
  aberto: { label: 'Aberto', tone: 'blue' },
  em_andamento: { label: 'Em andamento', tone: 'amber' },
  resolvido: { label: 'Resolvido', tone: 'green' },
  fechado: { label: 'Fechado', tone: 'slate' },
}

const PRIORITY_META: Record<Priority, { label: string; tone: Tone }> = {
  baixa: { label: 'Baixa', tone: 'slate' },
  media: { label: 'Média', tone: 'amber' },
  alta: { label: 'Alta', tone: 'red' },
}

export const STATUS_OPTIONS = Object.entries(STATUS_META).map(([value, meta]) => ({
  value: value as Status,
  label: meta.label,
}))

export const PRIORITY_OPTIONS = Object.entries(PRIORITY_META).map(([value, meta]) => ({
  value: value as Priority,
  label: meta.label,
}))

export function statusLabel(status: Status): string {
  return STATUS_META[status].label
}

export function statusTone(status: Status): Tone {
  return STATUS_META[status].tone
}

export function priorityLabel(priority: Priority): string {
  return PRIORITY_META[priority].label
}

export function priorityTone(priority: Priority): Tone {
  return PRIORITY_META[priority].tone
}
