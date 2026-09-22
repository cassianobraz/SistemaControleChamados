<script setup lang="ts">
import BaseButton from '@/components/ui/BaseButton.vue'
import type { Responsible } from '@/types/responsible'
import type { Priority, SortOrder, Status } from '@/types/ticket'
import { PRIORITY_OPTIONS, STATUS_OPTIONS } from '@/utils/ticketMeta'

defineProps<{
  responsibles: readonly Responsible[]
}>()

const emit = defineEmits<{
  apply: []
}>()

const status = defineModel<Status | ''>('status', { required: true })
const priority = defineModel<Priority | ''>('priority', { required: true })
const responsibleId = defineModel<string>('responsibleId', { required: true })
const search = defineModel<string>('search', { required: true })
const sort = defineModel<SortOrder>('sort', { required: true })
</script>

<template>
  <form
    class="space-y-3 rounded-lg border border-slate-200 bg-white p-4 shadow-sm"
    @submit.prevent="emit('apply')"
  >
    <div>
      <label class="mb-1 block text-xs font-medium text-slate-600" for="filter-search">Buscar</label>
      <input
        id="filter-search"
        v-model="search"
        type="search"
        placeholder="Título ou descrição..."
        class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus:border-brand-500 focus:ring-1 focus:ring-brand-500 focus:outline-none"
      />
    </div>

    <div class="flex flex-wrap items-end gap-3">
      <div class="min-w-36 flex-1">
        <label class="mb-1 block text-xs font-medium text-slate-600" for="filter-status">Status</label>
        <select
          id="filter-status"
          v-model="status"
          class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus:border-brand-500 focus:ring-1 focus:ring-brand-500 focus:outline-none"
        >
          <option value="">Todos</option>
          <option v-for="option in STATUS_OPTIONS" :key="option.value" :value="option.value">
            {{ option.label }}
          </option>
        </select>
      </div>

      <div class="min-w-36 flex-1">
        <label class="mb-1 block text-xs font-medium text-slate-600" for="filter-priority">Prioridade</label>
        <select
          id="filter-priority"
          v-model="priority"
          class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus:border-brand-500 focus:ring-1 focus:ring-brand-500 focus:outline-none"
        >
          <option value="">Todas</option>
          <option v-for="option in PRIORITY_OPTIONS" :key="option.value" :value="option.value">
            {{ option.label }}
          </option>
        </select>
      </div>

      <div class="min-w-40 flex-1">
        <label class="mb-1 block text-xs font-medium text-slate-600" for="filter-responsible">Responsável</label>
        <select
          id="filter-responsible"
          v-model="responsibleId"
          class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus:border-brand-500 focus:ring-1 focus:ring-brand-500 focus:outline-none"
        >
          <option value="">Todos</option>
          <option v-for="responsible in responsibles" :key="responsible.id" :value="responsible.id">
            {{ responsible.name }}
          </option>
        </select>
      </div>

      <div class="min-w-44 flex-1">
        <label class="mb-1 block text-xs font-medium text-slate-600" for="filter-sort">Ordenar por</label>
        <select
          id="filter-sort"
          v-model="sort"
          class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus:border-brand-500 focus:ring-1 focus:ring-brand-500 focus:outline-none"
        >
          <option value="created_at_desc">Mais recentes</option>
          <option value="created_at_asc">Mais antigos</option>
          <option value="priority_desc">Maior prioridade</option>
        </select>
      </div>

      <BaseButton type="submit" class="w-full sm:w-auto">Filtrar</BaseButton>
    </div>
  </form>
</template>
