<template>
  <Dialog
    v-model:visible="visible"
    modal
    :showHeader="false"
    :style="{ width: '54rem', maxWidth: '96vw' }"
    :pt="{ content: { class: 'salah-dialog-content' } }"
  >
    <div class="salah-client-modal">
      <header class="modal-header">
        <div class="header-icon" aria-hidden="true">
          <svg viewBox="0 0 24 24">
            <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
            <circle cx="9" cy="7" r="4" />
            <path d="M22 21v-2a4 4 0 0 0-3-3.87" />
            <path d="M16 3.13a4 4 0 0 1 0 7.75" />
          </svg>
        </div>
        <div class="header-copy">
          <h2>Editar Cliente</h2>
          <p>Registro comercial de Salah Motors</p>
          <span></span>
        </div>
        <button type="button" class="close-button" aria-label="Cerrar" @click="visible = false">
          <svg viewBox="0 0 24 24">
            <path d="M18 6 6 18" />
            <path d="m6 6 12 12" />
          </svg>
        </button>
      </header>

      <Form v-slot="$form" ref="validations" :resolver="resolver" :initialValues="initialValues" @submit="onFormSubmit" class="client-form">
        <section class="form-card">
          <div class="section-title">
            <span aria-hidden="true">
              <svg viewBox="0 0 24 24">
                <rect width="18" height="14" x="3" y="5" rx="2" />
                <path d="M7 10h4" />
                <path d="M7 14h7" />
                <circle cx="17" cy="12" r="1" />
              </svg>
            </span>
            <h3>Informacion Personal</h3>
          </div>

          <div class="fields-grid">
            <div class="field">
              <label for="nombre">Nombre</label>
              <div class="field-control">
                <svg viewBox="0 0 24 24" class="field-icon">
                  <path d="M20 21a8 8 0 0 0-16 0" />
                  <circle cx="12" cy="7" r="4" />
                </svg>
                <InputText id="nombre" name="nombre" v-model="initialValues.nombre" placeholder="Nombre" class="salah-input" />
              </div>
              <Message v-if="$form.nombre?.invalid" severity="error" size="small" variant="simple">
                {{ $form.nombre.error?.message }}
              </Message>
            </div>

            <div class="field">
              <label for="apellido">Apellido</label>
              <div class="field-control">
                <svg viewBox="0 0 24 24" class="field-icon">
                  <path d="M16 7a4 4 0 0 1-8 0" />
                  <path d="M12 14c-4 0-7 2-7 5v2h14v-2c0-3-3-5-7-5Z" />
                </svg>
                <InputText id="apellido" name="apellido" v-model="initialValues.apellido" placeholder="Apellido" class="salah-input" />
              </div>
              <Message v-if="$form.apellido?.invalid" severity="error" size="small" variant="simple">
                {{ $form.apellido.error?.message }}
              </Message>
            </div>

            <div class="field">
              <label for="ci">CI/NIT</label>
              <div class="field-control">
                <svg viewBox="0 0 24 24" class="field-icon">
                  <rect width="18" height="14" x="3" y="5" rx="2" />
                  <path d="M7 10h4" />
                  <path d="M7 14h7" />
                </svg>
                <InputText id="ci" name="ci" v-model="initialValues.ci" placeholder="CI/NIT" class="salah-input" :class="{ 'p-invalid': ciExistente }" />
              </div>
              <Message v-if="$form.ci?.invalid" severity="error" size="small" variant="simple">
                {{ $form.ci.error?.message }}
              </Message>
              <Message v-if="ciExistente" severity="error" size="small" variant="simple">
                Este CI/NIT ya existe
              </Message>
            </div>

            <div class="field">
              <label for="celular">Celular</label>
              <div class="field-control">
                <svg viewBox="0 0 24 24" class="field-icon">
                  <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.8 19.8 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6A19.8 19.8 0 0 1 2.1 4.18 2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.12.9.32 1.77.59 2.61a2 2 0 0 1-.45 2.11L8 9.69a16 16 0 0 0 6.31 6.31l1.25-1.25a2 2 0 0 1 2.11-.45c.84.27 1.71.47 2.61.59A2 2 0 0 1 22 16.92Z" />
                </svg>
                <InputText id="celular" name="celular" v-model="initialValues.celular" placeholder="Numero de celular" maxlength="8" class="salah-input" />
              </div>
              <Message v-if="$form.celular?.invalid" severity="error" size="small" variant="simple">
                {{ $form.celular.error?.message }}
              </Message>
            </div>

            <div class="field field-wide">
              <label for="direccion">Direccion</label>
              <div class="field-control">
                <svg viewBox="0 0 24 24" class="field-icon">
                  <path d="M20 10c0 5-8 12-8 12S4 15 4 10a8 8 0 1 1 16 0Z" />
                  <circle cx="12" cy="10" r="3" />
                </svg>
                <InputText id="direccion" name="direccion" v-model="initialValues.direccion" placeholder="Direccion del cliente" class="salah-input" />
              </div>
              <Message v-if="$form.direccion?.invalid" severity="error" size="small" variant="simple">
                {{ $form.direccion.error?.message }}
              </Message>
            </div>

            <div class="field field-wide">
              <label for="estado">Estado</label>
              <div class="field-control">
                <svg viewBox="0 0 24 24" class="field-icon">
                  <path d="M12 2v20" />
                  <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7H14a3.5 3.5 0 0 1 0 7H6" />
                </svg>
                <Select id="estado" name="estado" :options="Estados" v-model="initialValues.estado" placeholder="Seleccione un estado" class="salah-select" />
              </div>
              <Message v-if="$form.estado?.invalid" severity="error" size="small" variant="simple">
                {{ $form.estado.error?.message }}
              </Message>
            </div>
          </div>
        </section>

        <footer class="modal-actions">
          <Button type="submit" label="Guardar cambios" class="salah-submit" />
        </footer>
      </Form>
    </div>
  </Dialog>
</template>

<script setup lang="ts">
import { zodResolver } from '@primevue/forms/resolvers/zod';
import { z } from 'zod';
import { server } from '~/server/server';

interface Props { open : boolean, id: number }
const props = defineProps<Props>()
const emit = defineEmits(['close', 'success', 'update', 'error'])
const visible = ref(props.open)
const validations = ref()

const ciExistente = ref(false)
const clientes = ref<any[]>([])
const Estados = ref(['Activo', 'Inactivo'])

const initialValues = reactive({
  nombre: '',
  apellido: '',
  ci: '',
  celular: '',
  direccion: '',
  estado: ''
})

onMounted(async () => {
  try {
    const resValue:any = await $fetch(server.HOST + '/api/v1/clientes/' + props.id, { method: 'GET' })

    initialValues.nombre = resValue.nombre || ''
    initialValues.apellido = resValue.apellido || ''
    initialValues.ci = resValue.ci || ''
    initialValues.celular = resValue.celular || ''
    initialValues.direccion = resValue.direccion || ''
    initialValues.estado = resValue.estado || ''

    validations.value?.reset()

    const res = await $fetch(server.HOST + '/api/v1/clientes', { method: 'GET' })
    clientes.value = Array.isArray(res) ? res : []
  } catch(err) {
    console.error(err)
  }
})

watch(() => props.open, (newValue) => {
  visible.value = newValue
})

watch(visible, (newValue) => {
  if (!newValue) { emit('close') }
})

watch(() => initialValues.ci, () => {
  ciExistente.value = false
})

const resolver = ref(zodResolver(
  z.object({
    nombre: z.string().min(1, { message: 'Nombre requerido.' }),
    apellido: z.string().min(1, {message: 'Apellido requerido.' }),
    ci: z.string().min(1, {message: 'CI/NIT requerido.'}),
    celular: z.string().length(8, { message: 'El numero de celular debe tener 8 digitos.' }),
    direccion: z.string().min(1, {message: 'Direccion requerida.'}),
    estado: z.enum(['Activo', 'Inactivo'], { message: 'Debe seleccionar un estado valido.' })
  })
))

async function onFormSubmit({ valid } : any ) {
  if (valid) {
    initialValues.ci = initialValues.ci.trim()
    const ci = normalizarCi(initialValues.ci)
    const ciDuplicado = clientes.value.find(c => normalizarCi(c.ci) === ci && c.id !== props.id)

    ciExistente.value = !!ciDuplicado
    if (ciDuplicado) {
      return
    }

    try {
      await $fetch(server.HOST + '/api/v1/clientes/' + props.id, {
        method: 'PUT',
        body: initialValues
      })
      emit('update')
      emit('success')
      visible.value = false
    } catch(err: any) {
      console.error(err)
      if (err?.response?.status === 409 || err?.statusCode === 409) {
        ciExistente.value = true
        return
      }
      emit('error')
    }
  }
}

function normalizarCi(ci: string | null | undefined) {
  return String(ci || '').trim().toLowerCase()
}
</script>

<style scoped>
:deep(.salah-dialog-content) { padding: 0; border-radius: 20px; overflow: hidden; background: #ffffff; }
:deep(.p-dialog) { border-radius: 20px; box-shadow: 0 24px 80px rgba(13, 13, 13, 0.32); }
.salah-client-modal { background: #f7f7f7; color: #0d0d0d; }
.modal-header { position: relative; display: flex; align-items: center; gap: 16px; padding: 24px 28px; background: linear-gradient(135deg, #0d0d0d 0%, #202020 100%); color: #ffffff; }
.header-icon, .section-title span { display: grid; place-items: center; flex: 0 0 auto; }
.header-icon { width: 52px; height: 52px; border-radius: 16px; background: #ffd700; color: #0d0d0d; box-shadow: 0 12px 30px rgba(255, 215, 0, 0.24); }
svg { width: 22px; height: 22px; fill: none; stroke: currentColor; stroke-width: 2; stroke-linecap: round; stroke-linejoin: round; }
.header-copy h2 { margin: 0; font-size: 1.55rem; font-weight: 800; letter-spacing: 0; }
.header-copy p { margin: 2px 0 10px; color: rgba(255, 255, 255, 0.72); font-size: 0.92rem; }
.header-copy span { display: block; width: 96px; height: 3px; border-radius: 999px; background: #ffd700; }
.close-button { margin-left: auto; width: 42px; height: 42px; border: 1px solid rgba(255,255,255,0.14); border-radius: 12px; background: rgba(255,255,255,0.08); color: #ffffff; cursor: pointer; transition: background-color .2s ease, transform .2s ease; }
.close-button:hover { background: rgba(255,255,255,0.16); transform: translateY(-1px); }
.close-button svg { margin: auto; }
.client-form { display: grid; gap: 18px; padding: 20px; }
.form-card { padding: 18px; border: 1px solid rgba(13,13,13,0.08); border-radius: 18px; background: #ffffff; box-shadow: 0 18px 45px rgba(13,13,13,0.08); }
.section-title { display: flex; align-items: center; gap: 10px; margin-bottom: 16px; }
.section-title span { width: 36px; height: 36px; border-radius: 12px; background: rgba(255,215,0,0.18); color: #0d0d0d; }
.section-title h3 { margin: 0; font-size: 1rem; font-weight: 800; letter-spacing: 0; }
.fields-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 14px; }
.field { display: flex; min-width: 0; flex-direction: column; gap: 6px; }
.field-wide { grid-column: 1 / -1; }
.field label { color: #2a2a2a; font-size: .79rem; font-weight: 700; text-transform: uppercase; }
.field-control { position: relative; display: flex; align-items: center; }
.field-icon { position: absolute; left: 14px; z-index: 1; width: 19px; height: 19px; color: #6c6c6c; pointer-events: none; }
.salah-input, :deep(.salah-select) { width: 100%; }
:deep(.salah-input), :deep(.salah-select .p-select-label) { min-height: 50px; }
:deep(.salah-input), :deep(.salah-select) { border: 1px solid #d8d8d8; border-radius: 12px; background: #ffffff; color: #0d0d0d; box-shadow: none; transition: border-color .2s ease, box-shadow .2s ease; }
:deep(.salah-input) { padding-left: 44px; }
:deep(.salah-select) { min-height: 50px; padding-left: 34px; }
:deep(.salah-select .p-select-label) { display: flex; align-items: center; }
:deep(.salah-input:enabled:focus), :deep(.salah-select:not(.p-disabled).p-focus) { border-color: #ffd700; box-shadow: 0 0 0 4px rgba(255,215,0,0.18); }
:deep(.p-invalid) { border-color: #e30613; }
.modal-actions { display: flex; justify-content: flex-end; }
:deep(.salah-submit) { min-height: 50px; min-width: 210px; border: 0; border-radius: 12px; background: #ffd700; color: #0d0d0d; font-weight: 800; box-shadow: 0 14px 28px rgba(255,215,0,0.26); transition: background-color .2s ease, box-shadow .2s ease, transform .2s ease; }
:deep(.salah-submit:hover) { background: #e6c200; color: #0d0d0d; box-shadow: 0 18px 34px rgba(255,215,0,0.32); transform: translateY(-1px); }
:deep(.p-message-text) { font-size: .78rem; }
@media (max-width: 760px) {
  .modal-header { padding: 20px; }
  .client-form { padding: 14px; }
  .fields-grid { grid-template-columns: 1fr; }
  .field-wide { grid-column: auto; }
  .modal-actions { justify-content: stretch; }
  :deep(.salah-submit) { width: 100%; }
}
@media (max-width: 520px) {
  .modal-header { align-items: flex-start; }
  .header-icon { width: 46px; height: 46px; border-radius: 14px; }
  .header-copy h2 { font-size: 1.3rem; }
  .close-button { position: absolute; top: 14px; right: 14px; }
}
</style>
