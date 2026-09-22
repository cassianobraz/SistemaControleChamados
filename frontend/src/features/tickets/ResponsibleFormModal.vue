<script setup lang="ts">
import { reactive, watch } from 'vue'

import BaseButton from '@/components/ui/BaseButton.vue'
import BaseModal from '@/components/ui/BaseModal.vue'
import type { Responsible } from '@/types/responsible'

export interface ResponsibleFormValues {
  name: string
  email: string
  active: boolean
}

const props = defineProps<{
  responsible: Responsible | null
  submitting: boolean
  errorMessage: string
}>()

const open = defineModel<boolean>({ required: true })

const emit = defineEmits<{
  submit: [values: ResponsibleFormValues]
}>()

function blankForm(): ResponsibleFormValues {
  return { name: '', email: '', active: true }
}

const form = reactive<ResponsibleFormValues>(blankForm())

watch(
  open,
  (isOpen) => {
    if (!isOpen) return

    if (props.responsible) {
      form.name = props.responsible.name
      form.email = props.responsible.email
      form.active = props.responsible.active
      return
    }

    Object.assign(form, blankForm())
  },
  { immediate: true },
)

function handleSubmit() {
  emit('submit', { ...form })
}
</script>

<template>
  <BaseModal v-model="open" :title="responsible ? 'Editar responsável' : 'Novo responsável'">
    <form id="responsible-form" class="space-y-4" @submit.prevent="handleSubmit">
      <div>
        <label class="mb-1 block text-xs font-medium text-slate-600" for="responsible-name">Nome</label>
        <input
          id="responsible-name"
          v-model="form.name"
          type="text"
          required
          placeholder="Ex.: Ana Souza"
          class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus:border-brand-500 focus:ring-1 focus:ring-brand-500 focus:outline-none"
        />
      </div>

      <div>
        <label class="mb-1 block text-xs font-medium text-slate-600" for="responsible-email">E-mail</label>
        <input
          id="responsible-email"
          v-model="form.email"
          type="email"
          placeholder="ana@empresa.com"
          class="w-full rounded-md border border-slate-300 px-3 py-2 text-sm focus:border-brand-500 focus:ring-1 focus:ring-brand-500 focus:outline-none"
        />
      </div>

      <label v-if="responsible" class="flex items-center gap-2 text-sm text-slate-700">
        <input v-model="form.active" type="checkbox" class="rounded border-slate-300 text-brand-600 focus:ring-brand-500" />
        Ativo para distribuição automática
      </label>

      <p v-if="errorMessage" class="rounded-md bg-red-50 px-3 py-2 text-sm text-red-700">{{ errorMessage }}</p>
    </form>

    <template #footer>
      <BaseButton variant="secondary" type="button" :disabled="submitting" @click="open = false">
        Cancelar
      </BaseButton>
      <BaseButton type="submit" form="responsible-form" :disabled="submitting">
        {{ submitting ? 'Salvando...' : 'Salvar' }}
      </BaseButton>
    </template>
  </BaseModal>
</template>
