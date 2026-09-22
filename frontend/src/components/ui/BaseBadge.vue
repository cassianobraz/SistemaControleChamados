<script setup lang="ts">
import { computed } from 'vue'

import type { Tone } from '@/utils/ticketMeta'

const props = defineProps<{
  label: string
  tone: Tone
}>()

const toneClasses = computed<Record<Tone, string>>(() => ({
  slate: 'bg-slate-100 text-slate-700',
  blue: 'bg-blue-100 text-blue-700',
  amber: 'bg-amber-100 text-amber-800',
  green: 'bg-green-100 text-green-700',
  red: 'bg-red-100 text-red-700',
}))
</script>

<template>
  <span
    class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium whitespace-nowrap"
    :class="toneClasses[props.tone]"
  >
    {{ props.label }}
  </span>
</template>
