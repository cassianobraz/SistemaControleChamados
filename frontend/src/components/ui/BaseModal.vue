<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'

const open = defineModel<boolean>({ required: true })

defineProps<{ title: string }>()

function close() {
  open.value = false
}

function handleKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape' && open.value) close()
}

onMounted(() => window.addEventListener('keydown', handleKeydown))
onUnmounted(() => window.removeEventListener('keydown', handleKeydown))
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="open"
        class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 p-4"
        @click.self="close"
      >
        <div class="modal-panel w-full max-w-lg rounded-lg bg-white shadow-xl" role="dialog" aria-modal="true">
          <div class="flex items-center justify-between border-b border-slate-200 px-5 py-4">
            <h2 class="text-base font-semibold text-slate-900">{{ title }}</h2>
            <button
              type="button"
              class="rounded-md p-1 text-slate-400 hover:bg-slate-100 hover:text-slate-600"
              aria-label="Fechar"
              @click="close"
            >
              ✕
            </button>
          </div>

          <div class="max-h-[70vh] overflow-y-auto px-5 py-4">
            <slot />
          </div>

          <div v-if="$slots.footer" class="flex justify-end gap-2 border-t border-slate-200 px-5 py-4">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .modal-panel,
.modal-leave-active .modal-panel {
  transition: transform 0.2s ease;
}

.modal-enter-from .modal-panel,
.modal-leave-to .modal-panel {
  transform: scale(0.96) translateY(8px);
}
</style>
