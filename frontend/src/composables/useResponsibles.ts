import { readonly, ref } from 'vue'

import * as responsiblesApi from '@/api/responsibles'
import { extractErrorMessage } from '@/api/http'
import type { CreateResponsiblePayload, Responsible, UpdateResponsiblePayload } from '@/types/responsible'

export type ActionResult = { ok: true } | { ok: false; message: string }

export function useResponsibles() {
  const _responsibles = ref<Responsible[]>([])
  const _loading = ref(false)
  const _error = ref('')

  async function fetchResponsibles(): Promise<void> {
    _loading.value = true
    _error.value = ''

    try {
      _responsibles.value = await responsiblesApi.listResponsibles()
    } catch (error) {
      _error.value = extractErrorMessage(error)
    } finally {
      _loading.value = false
    }
  }

  async function runAction(action: () => Promise<unknown>): Promise<ActionResult> {
    try {
      await action()
      await fetchResponsibles()
      return { ok: true }
    } catch (error) {
      return { ok: false, message: extractErrorMessage(error) }
    }
  }

  function createResponsible(payload: CreateResponsiblePayload): Promise<ActionResult> {
    return runAction(() => responsiblesApi.createResponsible(payload))
  }

  function updateResponsible(id: string, payload: UpdateResponsiblePayload): Promise<ActionResult> {
    return runAction(() => responsiblesApi.updateResponsible(id, payload))
  }

  function deleteResponsible(id: string): Promise<ActionResult> {
    return runAction(() => responsiblesApi.deleteResponsible(id))
  }

  return {
    responsibles: readonly(_responsibles),
    loading: readonly(_loading),
    error: readonly(_error),
    fetchResponsibles,
    createResponsible,
    updateResponsible,
    deleteResponsible,
  }
}
