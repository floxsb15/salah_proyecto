<template>
  <Dialog v-model:visible="visible" modal header="Agregar Nuevo Segmento" :style="{ width: '25rem' }">
    <Form
      v-slot="$form"
      :resolver="resolver"
      :initialValues="initialValues"
      @submit="onFormSubmit"
      class="flex flex-col gap-2 w-full"
    >
      <div class="flex flex-col gap-1">
        <label for="categoria"> Categoria </label>
        <Select
          id="categoria"
          name="id_categoria"
          v-model="initialValues.id_categoria"
          :options="categorias"
          option-label="nombre"
          option-value="id"
          placeholder="Seleccione una categoria"
          fluid
          size="small"
        />
        <Message v-if="$form.id_categoria?.invalid" severity="error" size="small" variant="simple">
          {{ $form.id_categoria.error?.message }}
        </Message>
      </div>

      <div class="flex flex-col gap-1">
        <label for="nombre"> Nombre </label>
        <InputText
          id="nombre"
          name="nombre"
          v-model="initialValues.nombre"
          placeholder="Ingrese nombre"
          fluid
          size="small"
          :class="{ 'p-invalid': nombreExistente }"
        />
        <Message v-if="$form.nombre?.invalid" severity="error" size="small" variant="simple">
          {{ $form.nombre.error?.message }}
        </Message>
        <Message v-if="nombreExistente" severity="error" size="small" variant="simple">
          Este segmento ya existe en la categoria seleccionada
        </Message>
      </div>

      <div class="flex flex-col gap-1">
        <label for="descripcion"> Descripcion </label>
        <InputText
          id="descripcion"
          name="descripcion"
          v-model="initialValues.descripcion"
          placeholder="Ingrese la descripcion"
          fluid
          size="small"
        />
        <Message v-if="$form.descripcion?.invalid" severity="error" size="small" variant="simple">
          {{ $form.descripcion.error?.message }}
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
const categorias = ref<any[]>([])
const segmentos = ref<any[]>([])
const nombreExistente = ref(false)

const initialValues = reactive({
  id_categoria: 0,
  nombre: '',
  descripcion: '',
  estado: 'Activo',
})

onMounted(async () => {
  const [resCategorias, resSegmentos] = await Promise.all([
    $fetch(server.HOST + '/api/v1/categorias-vehiculos', { method: 'GET' }),
    $fetch(server.HOST + '/api/v1/segmentos-vehiculos', { method: 'GET' })
  ])
  categorias.value = Array.isArray(resCategorias) ? resCategorias : []
  segmentos.value = Array.isArray(resSegmentos) ? resSegmentos : []
})

watch(() => props.open, (newValue) => {
  visible.value = newValue
})

watch(visible, (newValue) => {
  if (!newValue) { emit('close') }
})

const resolver = ref(zodResolver(
  z.object({
    id_categoria: z.number().min(1, { message: 'Categoria requerida.' }),
    nombre: z.string().min(1, { message: 'Nombre requerido.' }),
    descripcion: z.string().min(1, { message: 'Descripcion requerida.' }),
    estado: z.enum(['Activo', 'Inactivo'], { message: 'Debe seleccionar un estado valido.' })
  })
))

async function onFormSubmit({ valid }: any) {
  if (!valid) {
    return
  }

  const nombreSegmento = initialValues.nombre.trim().toLowerCase()
  const segmentoExistente = segmentos.value.find(segmento =>
    segmento.nombre?.toLowerCase() === nombreSegmento &&
    Number(segmento.id_categoria) === Number(initialValues.id_categoria)
  )

  if (segmentoExistente) {
    nombreExistente.value = true
    return
  }

  nombreExistente.value = false

  try {
    await $fetch(server.HOST + '/api/v1/segmentos-vehiculos', {
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
