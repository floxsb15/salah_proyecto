<template>
  <Dialog v-model:visible="visible" modal header="Agregar Nuevo Anio" :style="{ width: '25rem' }">
    <Form
      v-slot="$form"
      :resolver="resolver"
      :initialValues="initialValues"
      @submit="onFormSubmit"
      class="flex flex-col gap-2 w-full"
    >
      <div class="flex flex-col gap-1">
        <label for="valor"> Anio </label>
        <InputNumber
          id="valor"
          name="valor"
          v-model="initialValues.valor"
          :min="1900"
          :max="2100"
          placeholder="Ej. 2026"
          fluid
          size="small"
          :class="{ 'p-invalid': anioExistente }"
        />
        <Message v-if="$form.valor?.invalid" severity="error" size="small" variant="simple">
          {{ $form.valor.error?.message }}
        </Message>
        <Message v-if="anioExistente" severity="error" size="small" variant="simple">
          Este anio ya existe
        </Message>
      </div>

      <div class="flex flex-col gap-1">
        <label for="estado"> Estado </label>
        <Select
          id="estado"
          name="estado"
          v-model="initialValues.estado"
          :options="Estados"
          placeholder="Seleccione un estado"
          fluid
          size="small"
        />
        <Message v-if="$form.estado?.invalid" severity="error" size="small" variant="simple">
          {{ $form.estado.error?.message }}
        </Message>
      </div>

      <Button type="submit" label="Agregar" />
    </Form>
  </Dialog>
</template>

<script setup lang="ts">
import { zodResolver } from '@primevue/forms/resolvers/zod';
import { z } from 'zod';
import { server } from '~/server/server';

interface Props { open: boolean }
const props = defineProps<Props>()
const emit = defineEmits(['close', 'success', 'update', 'error'])
const visible = ref(props.open)

const Estados = ref(['Activo', 'Inactivo'])
const anios = ref<any[]>([])
const anioExistente = ref(false)

const initialValues = reactive({
  valor: new Date().getFullYear(),
  estado: 'Activo',
})

onMounted(async () => {
  const res = await $fetch(server.HOST + '/api/v1/anios-vehiculos', { method: 'GET' })
  anios.value = Array.isArray(res) ? res : []
})

watch(() => props.open, (newValue) => {
  visible.value = newValue
})

watch(visible, (newValue) => {
  if (!newValue) { emit('close') }
})

const resolver = ref(zodResolver(
  z.object({
    valor: z.number().min(1900, { message: 'Anio requerido.' }).max(2100, { message: 'Anio no valido.' }),
    estado: z.enum(['Activo', 'Inactivo'], { message: 'Debe seleccionar un estado valido.' })
  })
))

async function onFormSubmit({ valid }: any) {
  if (!valid) {
    return
  }

  const anioExistenteValue = anios.value.find(anio => Number(anio.valor) === Number(initialValues.valor))
  if (anioExistenteValue) {
    anioExistente.value = true
    return
  }

  anioExistente.value = false
  try {
    await $fetch(server.HOST + '/api/v1/anios-vehiculos', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: initialValues
    })
    emit('update')
    emit('success')
    visible.value = false
  } catch (err) {
    console.error(err)
    emit('error')
  }
}
</script>
