<script setup lang="ts">
import { computed } from 'vue'

import BaseBadge from '@/components/ui/BaseBadge.vue'
import BaseButton from '@/components/ui/BaseButton.vue'
import type { Responsible } from '@/types/responsible'
import type { Ticket } from '@/types/ticket'
import { formatDateTime } from '@/utils/format'
import { priorityLabel, priorityTone, statusLabel, statusTone } from '@/utils/ticketMeta'

const props = defineProps<{
  ticket: Ticket
  responsibles: readonly Responsible[]
}>()

const emit = defineEmits<{
  view: [ticket: Ticket]
  edit: [ticket: Ticket]
  assign: [ticketId: string, responsibleId: string]
}>()

const selectedResponsibleId = computed({
  get: () => props.ticket.responsible_id,
  set: (value: string) => emit('assign', props.ticket.id, value),
})
</script>

<template>
  <li class="cursor-pointer space-y-3 p-4 hover:bg-slate-50" @click="emit('view', ticket)">
    <div class="flex items-start justify-between gap-3">
      <div class="min-w-0">
        <p class="truncate font-medium text-slate-900">{{ ticket.title }}</p>
        <p class="mt-0.5 line-clamp-2 text-sm text-slate-500">{{ ticket.description }}</p>
      </div>
      <div class="flex shrink-0 flex-col items-end gap-1">
        <BaseBadge :label="priorityLabel(ticket.priority)" :tone="priorityTone(ticket.priority)" />
        <BaseBadge :label="statusLabel(ticket.status)" :tone="statusTone(ticket.status)" />
      </div>
    </div>

    <div class="flex flex-wrap items-center justify-between gap-2" @click.stop>
      <span class="text-xs text-slate-400">{{ formatDateTime(ticket.created_at) }}</span>
      <select
        v-model="selectedResponsibleId"
        class="rounded-md border border-slate-300 px-2 py-1.5 text-xs focus:border-brand-500 focus:ring-1 focus:ring-brand-500 focus:outline-none"
      >
        <option v-for="responsible in responsibles" :key="responsible.id" :value="responsible.id">
          {{ responsible.name }}
        </option>
      </select>
    </div>

    <div class="flex justify-end gap-2 border-t border-slate-100 pt-3" @click.stop>
      <BaseButton variant="ghost" size="sm" @click="emit('edit', ticket)">Editar</BaseButton>
      <BaseButton variant="ghost" size="sm" @click="emit('assign', ticket.id, '')">Auto</BaseButton>
    </div>
  </li>
</template>
