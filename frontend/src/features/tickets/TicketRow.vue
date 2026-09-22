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
  <tr class="cursor-pointer border-b border-slate-100 last:border-0 hover:bg-slate-50" @click="emit('view', ticket)">
    <td class="max-w-56 px-3 py-3">
      <p class="truncate font-medium text-slate-900">{{ ticket.title }}</p>
      <p class="truncate text-sm text-slate-500">{{ ticket.description }}</p>
    </td>
    <td class="px-3 py-3">
      <BaseBadge :label="priorityLabel(ticket.priority)" :tone="priorityTone(ticket.priority)" />
    </td>
    <td class="px-3 py-3">
      <BaseBadge :label="statusLabel(ticket.status)" :tone="statusTone(ticket.status)" />
    </td>
    <td class="px-3 py-3" @click.stop>
      <select
        v-model="selectedResponsibleId"
        class="w-32 rounded-md border border-slate-300 px-2 py-1.5 text-xs focus:border-brand-500 focus:ring-1 focus:ring-brand-500 focus:outline-none"
      >
        <option v-for="responsible in responsibles" :key="responsible.id" :value="responsible.id">
          {{ responsible.name }}
        </option>
      </select>
    </td>
    <td class="px-3 py-3 text-sm whitespace-nowrap text-slate-500">{{ formatDateTime(ticket.created_at) }}</td>
    <td class="px-3 py-3 text-right whitespace-nowrap" @click.stop>
      <BaseButton variant="ghost" size="sm" @click="emit('edit', ticket)">Editar</BaseButton>
      <BaseButton variant="ghost" size="sm" @click="emit('assign', ticket.id, '')">Auto</BaseButton>
    </td>
  </tr>
</template>
