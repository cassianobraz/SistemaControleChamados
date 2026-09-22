<script setup lang="ts">
import { computed } from 'vue'

import BaseButton from '@/components/ui/BaseButton.vue'

const props = defineProps<{
  page: number
  totalPages: number
  disabled?: boolean
}>()

const emit = defineEmits<{
  change: [page: number]
}>()

const isFirstPage = computed(() => props.page <= 1)
const isLastPage = computed(() => props.page >= props.totalPages)
</script>

<template>
  <div v-if="totalPages > 1" class="flex items-center justify-between gap-3">
    <BaseButton
      variant="secondary"
      :disabled="disabled || isFirstPage"
      @click="emit('change', props.page - 1)"
    >
      Anterior
    </BaseButton>

    <span class="text-sm text-slate-600">Página {{ page }} de {{ totalPages }}</span>

    <BaseButton variant="secondary" :disabled="disabled || isLastPage" @click="emit('change', props.page + 1)">
      Próxima
    </BaseButton>
  </div>
</template>
