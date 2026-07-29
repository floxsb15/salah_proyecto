<template>
  <Dialog v-model:visible="visible" modal header="Agregar Nueva Marca" :style="{ width: '25rem' }">
    <Form
      v-slot="$form"
      :resolver="resolver"
      :initialValues="initialValues"
      @submit="onFormSubmit"
      class="flex flex-col gap-2 w-full"
    >
      <div class="flex flex-col gap-1">
        <label for="nombre"> Marca </label>
        <InputText
          id="nombre"
          name="nombre"
          v-model="initialValues.nombre"
          placeholder="Ej. Toyota"
          fluid
          size="small"
          :class="{ 'p-invalid': nombreExistente }"
        />
        <Message v-if="$form.nombre?.invalid" severity="error" size="small" variant="simple">
          {{ $form.nombre.error?.message }}
        </Message>
        <Message v-if="nombreExistente" severity="error" size="small" variant="simple">
          Esta marca ya existe
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
const marcas = ref<any[]>([])
const nombreExistente = ref(false)

const initialValues = reactive({
  nombre: '',
  estado: 'Activo',
})

onMounted(async () => {
  const res = await $fetch(server.HOST + '/api/v1/marcas-vehiculos', { method: 'GET' })
  marcas.value = Array.isArray(res) ? res : []
})

watch(() => props.open, (newValue) => {
  visible.value = newValue
})

watch(visible, (newValue) => {
  if (!newValue) { emit('close') }
})

const resolver = ref(zodResolver(
  z.object({
    nombre: z.string().min(1, { message: 'Marca requerida.' }),
    estado: z.enum(['Activo', 'Inactivo'], { message: 'Debe seleccionar un estado valido.' })
  })
))

async function onFormSubmit({ valid }: any) {
  if (!valid) {
    return
  }

  const nombreMarca = initialValues.nombre.trim().toLowerCase()
  const marcaExistente = marcas.value.find(marca => marca.nombre?.toLowerCase() === nombreMarca)
  if (marcaExistente) {
    nombreExistente.value = true
    return
  }

  nombreExistente.value = false
  try {
    await $fetch(server.HOST + '/api/v1/marcas-vehiculos', {
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
