<script setup lang="ts">
import BaseBadge from '@/components/ui/BaseBadge.vue'
import BaseButton from '@/components/ui/BaseButton.vue'
import BaseModal from '@/components/ui/BaseModal.vue'
import type { Ticket } from '@/types/ticket'
import { formatDateTime } from '@/utils/format'
import { priorityLabel, priorityTone, statusLabel, statusTone } from '@/utils/ticketMeta'

const props = defineProps<{
  ticket: Ticket | null
}>()

const open = defineModel<boolean>({ required: true })

const emit = defineEmits<{
  edit: [ticket: Ticket]
}>()

function handleEdit() {
  if (props.ticket) emit('edit', props.ticket)
}
</script>

<template>
  <BaseModal v-model="open" :title="ticket?.title ?? 'Chamado'">
    <div v-if="ticket" class="space-y-4">
      <div class="flex flex-wrap gap-2">
        <BaseBadge :label="priorityLabel(ticket.priority)" :tone="priorityTone(ticket.priority)" />
        <BaseBadge :label="statusLabel(ticket.status)" :tone="statusTone(ticket.status)" />
      </div>

      <div>
        <h3 class="mb-1 text-xs font-medium tracking-wide text-slate-500 uppercase">Descrição</h3>
        <p class="text-sm whitespace-pre-wrap text-slate-700">{{ ticket.description }}</p>
      </div>

      <dl class="grid grid-cols-2 gap-4 border-t border-slate-100 pt-4 text-sm">
        <div>
          <dt class="text-xs font-medium text-slate-500">Responsável</dt>
          <dd class="mt-0.5 text-slate-700">{{ ticket.responsible_name }}</dd>
        </div>
        <div>
          <dt class="text-xs font-medium text-slate-500">Aberto em</dt>
          <dd class="mt-0.5 text-slate-700">{{ formatDateTime(ticket.created_at) }}</dd>
        </div>
        <div class="col-span-2">
          <dt class="text-xs font-medium text-slate-500">Última atualização</dt>
          <dd class="mt-0.5 text-slate-700">{{ formatDateTime(ticket.updated_at) }}</dd>
        </div>
      </dl>
    </div>

    <template #footer>
      <BaseButton variant="secondary" type="button" @click="open = false">Fechar</BaseButton>
      <BaseButton type="button" @click="handleEdit">Editar</BaseButton>
    </template>
  </BaseModal>
</template>
