<script setup lang="ts">
import TicketCard from '@/features/tickets/TicketCard.vue'
import TicketRow from '@/features/tickets/TicketRow.vue'
import type { Responsible } from '@/types/responsible'
import type { Ticket } from '@/types/ticket'

defineProps<{
  tickets: readonly Ticket[]
  responsibles: readonly Responsible[]
  loading: boolean
}>()

const emit = defineEmits<{
  view: [ticket: Ticket]
  edit: [ticket: Ticket]
  assign: [ticketId: string, responsibleId: string]
}>()
</script>

<template>
  <div class="overflow-hidden rounded-lg border border-slate-200 bg-white shadow-sm">
    <div class="hidden overflow-x-auto md:block">
      <table class="w-full text-left text-sm">
        <thead class="bg-slate-50 text-xs tracking-wide text-slate-500 uppercase">
          <tr>
            <th class="px-3 py-3 font-medium">Chamado</th>
            <th class="px-3 py-3 font-medium">Prioridade</th>
            <th class="px-3 py-3 font-medium">Status</th>
            <th class="px-3 py-3 font-medium">Responsável</th>
            <th class="px-3 py-3 font-medium">Aberto em</th>
            <th class="px-3 py-3 font-medium">
              <span class="sr-only">Ações</span>
            </th>
          </tr>
        </thead>
        <tbody>
          <TicketRow
            v-for="ticket in tickets"
            :key="ticket.id"
            :ticket="ticket"
            :responsibles="responsibles"
            @view="(t) => emit('view', t)"
            @edit="(t) => emit('edit', t)"
            @assign="(id, responsibleId) => emit('assign', id, responsibleId)"
          />
        </tbody>
      </table>
    </div>

    <ul class="divide-y divide-slate-100 md:hidden">
      <TicketCard
        v-for="ticket in tickets"
        :key="ticket.id"
        :ticket="ticket"
        :responsibles="responsibles"
        @view="(t) => emit('view', t)"
        @edit="(t) => emit('edit', t)"
        @assign="(id, responsibleId) => emit('assign', id, responsibleId)"
      />
    </ul>

    <p v-if="loading" class="px-4 py-8 text-center text-sm text-slate-500">Carregando chamados...</p>
    <p v-else-if="tickets.length === 0" class="px-4 py-8 text-center text-sm text-slate-500">
      Nenhum chamado encontrado com os filtros atuais.
    </p>
  </div>
</template>
