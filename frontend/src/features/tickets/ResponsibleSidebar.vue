<script setup lang="ts">
import { computed } from 'vue'

import BaseButton from '@/components/ui/BaseButton.vue'
import type { Responsible } from '@/types/responsible'

const props = defineProps<{
  responsibles: readonly Responsible[]
  loading: boolean
}>()

const emit = defineEmits<{
  create: []
  edit: [responsible: Responsible]
  delete: [responsible: Responsible]
}>()

const maxOpenTickets = computed(() => Math.max(1, ...props.responsibles.map((r) => r.open_tickets)))
</script>

<template>
  <aside class="rounded-lg border border-slate-200 bg-white p-4 shadow-sm">
    <div class="mb-3 flex items-center justify-between">
      <h2 class="text-sm font-semibold text-slate-900">Carga da equipe</h2>
      <BaseButton variant="ghost" size="sm" @click="emit('create')">+ Novo</BaseButton>
    </div>

    <p v-if="loading" class="text-sm text-slate-500">Carregando...</p>
    <p v-else-if="responsibles.length === 0" class="text-sm text-slate-500">Nenhum responsável cadastrado.</p>

    <ul v-else class="space-y-4">
      <li v-for="responsible in responsibles" :key="responsible.id">
        <div class="mb-1 flex items-center justify-between gap-2 text-sm">
          <span class="truncate font-medium text-slate-700" :class="{ 'text-slate-400': !responsible.active }">
            {{ responsible.name }}<span v-if="!responsible.active"> (inativo)</span>
          </span>
          <span class="shrink-0 text-slate-500">{{ responsible.open_tickets }} em aberto</span>
        </div>
        <div class="h-1.5 w-full overflow-hidden rounded-full bg-slate-100">
          <div
            class="h-full rounded-full bg-brand-500"
            :style="{ width: `${(responsible.open_tickets / maxOpenTickets) * 100}%` }"
          />
        </div>
        <div class="mt-1 flex justify-end gap-1">
          <BaseButton variant="ghost" size="sm" @click="emit('edit', responsible)">Editar</BaseButton>
          <BaseButton variant="ghost" size="sm" @click="emit('delete', responsible)">Excluir</BaseButton>
        </div>
      </li>
    </ul>

    <p class="mt-4 text-xs text-slate-400">
      A distribuição automática escolhe sempre quem tem menos chamados em aberto no momento.
    </p>
  </aside>
</template>
