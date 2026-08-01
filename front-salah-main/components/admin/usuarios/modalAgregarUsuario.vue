<template>
  <Dialog
    v-model:visible="visible"
    modal
    :showHeader="false"
    :style="{ width: '62rem', maxWidth: '96vw' }"
    :pt="{ content: { class: 'salah-dialog-content' } }"
  >
    <div class="salah-user-modal">
      <header class="modal-header">
        <div class="header-icon" aria-hidden="true">
          <svg viewBox="0 0 24 24">
            <path d="M20 21a8 8 0 0 0-16 0" />
            <circle cx="12" cy="7" r="4" />
          </svg>
        </div>

        <div class="header-copy">
          <h2>Nuevo Usuario</h2>
          <p>Personal de Salah Motors</p>
          <span></span>
        </div>

        <button type="button" class="close-button" aria-label="Cerrar" @click="visible = false">
          <svg viewBox="0 0 24 24">
            <path d="M18 6 6 18" />
            <path d="m6 6 12 12" />
          </svg>
        </button>
      </header>

      <Form
        v-slot="$form"
        :resolver="resolver"
        :initialValues="initialValues"
        @submit="onFormSubmit"
        class="user-form"
      >
        <section class="form-card personal-card">
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
              <label for="ci">CI</label>
              <div class="field-control">
                <svg viewBox="0 0 24 24" class="field-icon">
                  <rect width="18" height="14" x="3" y="5" rx="2" />
                  <path d="M7 10h4" />
                  <path d="M7 14h7" />
                </svg>
                <InputText id="ci" name="ci" v-model="initialValues.ci" placeholder="Documento" class="salah-input" :class="{ 'p-invalid': ciExistente }" />
              </div>
              <Message v-if="$form.ci?.invalid" severity="error" size="small" variant="simple">
                {{ $form.ci.error?.message }}
              </Message>
              <Message v-if="ciExistente" severity="error" size="small" variant="simple">
                Este CI ya existe
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
                <InputText id="direccion" name="direccion" v-model="initialValues.direccion" placeholder="Direccion del usuario" class="salah-input" />
              </div>
              <Message v-if="$form.direccion?.invalid" severity="error" size="small" variant="simple">
                {{ $form.direccion.error?.message }}
              </Message>
            </div>
          </div>
        </section>

        <section class="form-card access-card">
          <div class="section-title">
            <span aria-hidden="true">
              <svg viewBox="0 0 24 24">
                <path d="M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2 0 4.5-1.2 6.24-2.72a1.17 1.17 0 0 1 1.52 0C14.51 3.81 17 5 19 5a1 1 0 0 1 1 1Z" />
              </svg>
            </span>
            <h3>Acceso</h3>
          </div>

          <div class="fields-grid">
            <div class="field">
              <label for="usuario">Usuario</label>
              <div class="field-control">
                <svg viewBox="0 0 24 24" class="field-icon">
                  <path d="M16 8a6 6 0 1 0-2.2 4.64" />
                  <circle cx="12" cy="8" r="2" />
                  <path d="M16 8v3a2 2 0 0 0 4 0V8a8 8 0 1 0-4.74 7.3" />
                </svg>
                <InputText id="usuario" name="usuario" v-model="initialValues.usuario" placeholder="Usuario" class="salah-input" :class="{ 'p-invalid': usuarioExistente }" />
              </div>
              <Message v-if="$form.usuario?.invalid" severity="error" size="small" variant="simple">
                {{ $form.usuario.error?.message }}
              </Message>
              <Message v-if="usuarioExistente" severity="error" size="small" variant="simple">
                Este usuario ya existe
              </Message>
            </div>

            <div class="field">
              <label for="rol">Rol</label>
              <div class="field-control">
                <svg viewBox="0 0 24 24" class="field-icon">
                  <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10Z" />
                  <path d="m9 12 2 2 4-4" />
                </svg>
                <Select id="rol" name="rol" :options="Roles" v-model="initialValues.rol" placeholder="Seleccione un rol" class="salah-select" />
              </div>
              <Message v-if="$form.rol?.invalid" severity="error" size="small" variant="simple">
                {{ $form.rol.error?.message }}
              </Message>
            </div>

            <div class="field">
              <label for="contra">Contrasena</label>
              <div class="field-control">
                <svg viewBox="0 0 24 24" class="field-icon">
                  <rect width="18" height="11" x="3" y="11" rx="2" />
                  <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                </svg>
                <InputText id="contra" name="contra" :type="showPassword ? 'text' : 'password'" v-model="initialValues.contra" placeholder="Contrasena" class="salah-input with-action" />
                <button type="button" class="field-action" aria-label="Mostrar contrasena" @click="showPassword = !showPassword">
                  <i class="pi" :class="showPassword ? 'pi-eye-slash' : 'pi-eye'"></i>
                </button>
              </div>
              <Message v-if="$form.contra?.invalid" severity="error" size="small" variant="simple">
                {{ $form.contra.error?.message }}
              </Message>
            </div>

            <div class="field">
              <label for="confirmarContra">Confirmar contrasena</label>
              <div class="field-control">
                <svg viewBox="0 0 24 24" class="field-icon">
                  <path d="M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2 0 4.5-1.2 6.24-2.72a1.17 1.17 0 0 1 1.52 0C14.51 3.81 17 5 19 5a1 1 0 0 1 1 1Z" />
                  <path d="m9 12 2 2 4-4" />
                </svg>
                <InputText id="confirmarContra" name="confirmarContra" :type="showConfirmPassword ? 'text' : 'password'" v-model="initialValues.confirmarContra" placeholder="Confirmar contrasena" class="salah-input with-action" />
                <button type="button" class="field-action" aria-label="Mostrar confirmacion" @click="showConfirmPassword = !showConfirmPassword">
                  <i class="pi" :class="showConfirmPassword ? 'pi-eye-slash' : 'pi-eye'"></i>
                </button>
              </div>
              <Message v-if="$form.confirmarContra?.invalid" severity="error" size="small" variant="simple">
                {{ $form.confirmarContra.error?.message }}
              </Message>
            </div>
          </div>
        </section>

        <section class="form-card photo-card">
          <div class="section-title">
            <span aria-hidden="true">
              <svg viewBox="0 0 24 24">
                <path d="M14.5 4h-5L7 7H4a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3l-2.5-3Z" />
                <circle cx="12" cy="13" r="3" />
              </svg>
            </span>
            <h3>Fotografia</h3>
          </div>

          <button
            type="button"
            class="photo-dropzone"
            :class="{ 'is-dragging': isDragging }"
            @click="triggerPhotoInput"
            @dragenter.prevent="isDragging = true"
            @dragover.prevent="isDragging = true"
            @dragleave.prevent="isDragging = false"
            @drop.prevent="onPhotoDrop"
          >
            <input ref="photoInput" type="file" accept="image/*" class="photo-input" @change="onFileSelect" />
            <span class="photo-preview" aria-hidden="true">
              <img v-if="src" :src="src" alt="" />
              <svg v-else viewBox="0 0 24 24">
                <path d="M12 5v14" />
                <path d="m19 12-7 7-7-7" />
              </svg>
            </span>
            <span class="photo-copy">
              <strong>{{ src ? 'Fotografia seleccionada' : 'Subir fotografia' }}</strong>
              <small>Arrastra una imagen o haz clic para seleccionar</small>
            </span>
          </button>
          <Message v-if="$form.foto?.invalid" severity="error" size="small" variant="simple">
            {{ $form.foto.error?.message }}
          </Message>
        </section>

        <footer class="modal-actions">
          <Button type="submit" label="Registrar usuario" class="salah-submit" />
        </footer>
      </Form>
    </div>
  </Dialog>
</template>

<script setup lang="ts">
import { zodResolver } from '@primevue/forms/resolvers/zod';
import { z } from 'zod';
import { server } from '~/server/server';

interface Props { open : boolean }
const props = defineProps<Props>()
const emit = defineEmits(['close', 'success', 'update', 'error'])
const visible = ref(props.open)

const ciExistente = ref(false)
const usuarioExistente = ref(false)
const usuarios = ref<any[]>([])
const showPassword = ref(false)
const showConfirmPassword = ref(false)
const src = ref<string>()
const isDragging = ref(false)
const photoInput = ref<HTMLInputElement | null>(null)

const initialValues = reactive({
  nombre: '',
  apellido: '',
  ci: '',
  celular: '',
  direccion: '',
  foto: null as File | null,
  usuario: '',
  contra: '',
  confirmarContra: '',
  rol: '',
  estado: 'Activo'
})
const Roles = ref(['admin', 'encargado de ventas', 'vendedor', 'contador'])

onMounted(async () => {
  const res:any[] = await $fetch(server.HOST + '/api/v1/usuarios', {
    method: 'GET'
  })
  usuarios.value = res
})

watch(() => props.open, (newValue) => {
  visible.value = newValue
})

watch(visible, (newValue) => {
  if (!newValue) { emit('close') }
})

const resolver = ref(zodResolver(
  z.object({
    nombre: z.string().min(1, { message: 'Nombre requerido.' }),
    apellido: z.string().min(1, {message: 'Apellido requerido.' }),
    ci: z.string()
      .min(7, {message: 'CI requerido entre 7 a 9 caracteres.'})
      .max(9, {message: 'CI requerido entre 7 a 9 caracteres.'}),
    celular: z.string().length(8, { message: 'El numero de celular debe tener 8 digitos.' }),
    direccion: z.string().min(1, {message: 'Direccion requerida.'}),
    foto: z
      .instanceof(File)
      .nullable()
      .refine((file) => !file || file.type.startsWith('image/'), {
        message: 'El archivo debe ser una imagen valida.',
      })
      .refine((file) => !file || file.size <= 2 * 1024 * 1024, {
        message: 'La imagen no debe superar los 2MB.',
      }),
    usuario: z.string().min(1, {message: 'Usuario requerido.'}),
    contra: z.string()
      .min(12, { message: 'La contrasena debe tener al menos 12 caracteres.' })
      .regex(/^(?=.*[a-zA-Z])(?=.*[0-9])/, { message: 'La contrasena debe contener letras y numeros.' }),
    confirmarContra: z.string()
      .min(12, { message: 'La contrasena debe tener al menos 12 caracteres.' })
      .regex(/^(?=.*[a-zA-Z])(?=.*[0-9])/, { message: 'La contrasena debe contener letras y numeros.' }),
    rol: z.string().min(1, { message: 'Rol requerido.' })
  }).refine(data => data.contra === data.confirmarContra, {
    message: 'Las contrasenas no coinciden',
    path: ['confirmarContra'],
  })
))

async function onFormSubmit({ valid } : any ) {
  if (valid) {
    const ci = initialValues.ci.trim().toLowerCase()
    const username = initialValues.usuario.trim().toLowerCase()

    const ciDuplicado = usuarios.value.find(u => u.ci.toLowerCase() === ci)
    const usuarioDuplicado = usuarios.value.find(u => u.usuario.toLowerCase() === username)

    ciExistente.value = !!ciDuplicado
    usuarioExistente.value = !!usuarioDuplicado

    if (ciDuplicado || usuarioDuplicado) {
      return
    }

    try {
      const formData = new FormData()
      formData.append('nombre', initialValues.nombre)
      formData.append('apellido', initialValues.apellido)
      formData.append('ci', initialValues.ci)
      formData.append('celular', initialValues.celular)
      formData.append('direccion', initialValues.direccion)
      formData.append('usuario', initialValues.usuario)
      formData.append('contra', initialValues.contra)
      formData.append('rol', initialValues.rol)
      formData.append('estado', initialValues.estado)
      if (initialValues.foto) {
        formData.append('foto', initialValues.foto)
      }

      await $fetch(server.HOST + '/api/v1/usuarios', {
        method: 'POST',
        body: formData
      })
      emit('update')
      emit('success')
      visible.value = false
    } catch (err) {
      console.error(err)
      emit('error')
    }
  }
}

function triggerPhotoInput() {
  photoInput.value?.click()
}

function onPhotoDrop(event: DragEvent) {
  isDragging.value = false
  const file = event.dataTransfer?.files?.[0]
  if (file) {
    setPhoto(file)
  }
}

function onFileSelect(event: Event | any) {
  const file = event?.files?.[0] || (event.target as HTMLInputElement)?.files?.[0]
  if (file) {
    setPhoto(file)
  }
}

function setPhoto(file: File) {
  initialValues.foto = file

  const reader = new FileReader()
  reader.onload = (e) => {
    src.value = String(e.target?.result || '')
  }
  reader.readAsDataURL(file)
}
</script>

<style scoped>
:deep(.salah-dialog-content) {
  padding: 0;
  border-radius: 20px;
  overflow: hidden;
  background: #ffffff;
}

:deep(.p-dialog) {
  border-radius: 20px;
  box-shadow: 0 24px 80px rgba(13, 13, 13, 0.32);
}

.salah-user-modal {
  background: #f7f7f7;
  color: #0d0d0d;
}

.modal-header {
  position: relative;
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px 28px;
  background: linear-gradient(135deg, #0d0d0d 0%, #202020 100%);
  color: #ffffff;
}

.header-icon,
.section-title span {
  display: grid;
  place-items: center;
  flex: 0 0 auto;
}

.header-icon {
  width: 52px;
  height: 52px;
  border-radius: 16px;
  background: #ffd700;
  color: #0d0d0d;
  box-shadow: 0 12px 30px rgba(255, 215, 0, 0.24);
}

svg {
  width: 22px;
  height: 22px;
  fill: none;
  stroke: currentColor;
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.header-copy h2 {
  margin: 0;
  font-size: 1.55rem;
  font-weight: 800;
  letter-spacing: 0;
}

.header-copy p {
  margin: 2px 0 10px;
  color: rgba(255, 255, 255, 0.72);
  font-size: 0.92rem;
}

.header-copy span {
  display: block;
  width: 96px;
  height: 3px;
  border-radius: 999px;
  background: #ffd700;
}

.close-button {
  margin-left: auto;
  width: 42px;
  height: 42px;
  border: 1px solid rgba(255, 255, 255, 0.14);
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.08);
  color: #ffffff;
  cursor: pointer;
  transition: background-color 0.2s ease, transform 0.2s ease;
}

.close-button:hover {
  background: rgba(255, 255, 255, 0.16);
  transform: translateY(-1px);
}

.close-button svg {
  margin: auto;
}

.user-form {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 320px;
  gap: 18px;
  padding: 20px;
}

.form-card {
  border: 1px solid rgba(13, 13, 13, 0.08);
  border-radius: 18px;
  background: #ffffff;
  box-shadow: 0 18px 45px rgba(13, 13, 13, 0.08);
}

.personal-card,
.access-card {
  padding: 18px;
}

.personal-card,
.access-card,
.modal-actions {
  grid-column: 1;
}

.photo-card {
  grid-column: 2;
  grid-row: 1 / span 2;
  padding: 18px;
  align-self: stretch;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
}

.section-title span {
  width: 36px;
  height: 36px;
  border-radius: 12px;
  background: rgba(255, 215, 0, 0.18);
  color: #0d0d0d;
}

.section-title h3 {
  margin: 0;
  font-size: 1rem;
  font-weight: 800;
  letter-spacing: 0;
}

.fields-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
}

.field {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 6px;
}

.field-wide {
  grid-column: 1 / -1;
}

.field label {
  color: #2a2a2a;
  font-size: 0.79rem;
  font-weight: 700;
  text-transform: uppercase;
}

.field-control {
  position: relative;
  display: flex;
  align-items: center;
}

.field-icon {
  position: absolute;
  left: 14px;
  z-index: 1;
  width: 19px;
  height: 19px;
  color: #6c6c6c;
  pointer-events: none;
}

.salah-input,
:deep(.salah-select) {
  width: 100%;
}

:deep(.salah-input),
:deep(.salah-select .p-select-label) {
  min-height: 50px;
}

:deep(.salah-input),
:deep(.salah-select) {
  border: 1px solid #d8d8d8;
  border-radius: 12px;
  background: #ffffff;
  color: #0d0d0d;
  box-shadow: none;
  transition: border-color 0.2s ease, box-shadow 0.2s ease, background-color 0.2s ease;
}

:deep(.salah-input) {
  padding-left: 44px;
}

:deep(.salah-input.with-action) {
  padding-right: 48px;
}

:deep(.salah-select) {
  min-height: 50px;
  padding-left: 34px;
}

:deep(.salah-select .p-select-label) {
  display: flex;
  align-items: center;
}

:deep(.salah-input:enabled:focus),
:deep(.salah-select:not(.p-disabled).p-focus) {
  border-color: #ffd700;
  box-shadow: 0 0 0 4px rgba(255, 215, 0, 0.18);
}

:deep(.p-invalid) {
  border-color: #e30613;
}

.field-action {
  position: absolute;
  right: 6px;
  display: grid;
  place-items: center;
  width: 38px;
  height: 38px;
  border: 0;
  border-radius: 10px;
  background: transparent;
  color: #555555;
  cursor: pointer;
}

.field-action:hover {
  background: rgba(255, 215, 0, 0.18);
  color: #0d0d0d;
}

.photo-dropzone {
  display: flex;
  min-height: 286px;
  width: 100%;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 18px;
  border: 1.5px dashed rgba(13, 13, 13, 0.22);
  border-radius: 18px;
  background: linear-gradient(180deg, #ffffff 0%, #fafafa 100%);
  color: #0d0d0d;
  cursor: pointer;
  text-align: center;
  transition: border-color 0.2s ease, box-shadow 0.2s ease, transform 0.2s ease;
}

.photo-dropzone:hover,
.photo-dropzone.is-dragging {
  border-color: #ffd700;
  box-shadow: 0 14px 36px rgba(255, 215, 0, 0.16);
  transform: translateY(-1px);
}

.photo-input {
  display: none;
}

.photo-preview {
  display: grid;
  place-items: center;
  width: 112px;
  height: 112px;
  overflow: hidden;
  border: 3px solid #ffd700;
  border-radius: 24px;
  background: #0d0d0d;
  color: #ffd700;
  box-shadow: 0 18px 40px rgba(13, 13, 13, 0.18);
}

.photo-preview svg {
  width: 34px;
  height: 34px;
}

.photo-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.photo-copy {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.photo-copy strong {
  font-size: 1rem;
  font-weight: 800;
}

.photo-copy small {
  max-width: 220px;
  color: #696969;
  font-size: 0.82rem;
  line-height: 1.35;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
}

:deep(.salah-submit) {
  min-height: 50px;
  min-width: 210px;
  border: 0;
  border-radius: 12px;
  background: #ffd700;
  color: #0d0d0d;
  font-weight: 800;
  box-shadow: 0 14px 28px rgba(255, 215, 0, 0.26);
  transition: background-color 0.2s ease, box-shadow 0.2s ease, transform 0.2s ease;
}

:deep(.salah-submit:hover) {
  background: #e6c200;
  color: #0d0d0d;
  box-shadow: 0 18px 34px rgba(255, 215, 0, 0.32);
  transform: translateY(-1px);
}

:deep(.p-message-text) {
  font-size: 0.78rem;
}

@media (max-width: 860px) {
  .modal-header {
    padding: 20px;
  }

  .user-form {
    grid-template-columns: 1fr;
    padding: 14px;
  }

  .personal-card,
  .access-card,
  .photo-card,
  .modal-actions {
    grid-column: 1;
  }

  .photo-card {
    grid-row: auto;
  }

  .fields-grid {
    grid-template-columns: 1fr;
  }

  .photo-dropzone {
    min-height: 220px;
  }

  .modal-actions {
    justify-content: stretch;
  }

  :deep(.salah-submit) {
    width: 100%;
  }
}

@media (max-width: 520px) {
  .modal-header {
    align-items: flex-start;
  }

  .header-icon {
    width: 46px;
    height: 46px;
    border-radius: 14px;
  }

  .header-copy h2 {
    font-size: 1.3rem;
  }

  .close-button {
    position: absolute;
    top: 14px;
    right: 14px;
  }
}
</style>
