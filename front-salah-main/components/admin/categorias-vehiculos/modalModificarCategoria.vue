<template>
  <Dialog v-model:visible="visible" modal header="Modificar Categoría de Vehículo" :style="{ width: '25rem' }">
    <Form 
      v-slot="$form" :resolver="resolver" ref="validations"
      :initialValues="initialValues" @submit="onFormSubmit" 
      class="flex flex-col gap-2 w-full">
      
      <div class="flex flex-col gap-1">
        <label for="nombre"> Nombre </label>
        <InputText 
          id="nombre" name="nombre"
          v-model="initialValues.nombre"
          placeholder="ingrese nombre" 
          fluid size="small" 
          :class="{ 'p-invalid': nombreExistente }"
        />
        <Message v-if="$form.nombre?.invalid" severity="error" size="small" variant="simple">
          {{ $form.nombre.error?.message }}
        </Message>
        <Message v-if="nombreExistente" severity="error" size="small" variant="simple">
          Este nombre ya existe
        </Message>
      </div>

      <div class="flex flex-col gap-1">
        <label for="descripcion"> Descripción </label>
        <InputText 
          id="descripcion" name="descripcion"
          v-model="initialValues.descripcion" 
          placeholder="ingrese la descripcion"
          fluid size="small"/>
        <Message v-if="$form.descripcion?.invalid" severity="error" size="small" variant="simple">
          {{ $form.descripcion.error?.message }}
        </Message>
      </div>

      <div class="flex flex-col gap-1">
        <label for="estado"> Estado </label>
        <Select
          id="estado" name="estado"
          v-model="initialValues.estado"
          placeholder="seleccione su estado"
          :options="Estados"
          fluid size="small" />
        <Message v-if="$form.estado?.invalid" severity="error" size="small" variant="simple">
          {{ $form.estado.error?.message }}
        </Message>
      </div>

      <Button type="submit" label="Modificar" />
    </Form>
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

const Estados = ref(['Activo', 'Inactivo'])
const nombreExistente = ref(false)
const categorias = ref<any[]>([])

const initialValues = reactive({ 
  nombre: '',
  descripcion: '',
  estado: '', 
})

onMounted( async () => {
  try {
    const resValue:any = await $fetch(server.HOST + '/api/v1/categorias-vehiculos/' + props.id, {
      method: 'GET'
    })
    initialValues.nombre = resValue.nombre
    initialValues.descripcion = resValue.descripcion
    initialValues.estado = resValue.estado
    validations.value?.reset()

    const res: any[] | null = await $fetch(server.HOST + '/api/v1/categorias-vehiculos', {
      method: 'GET'
    })
    categorias.value = Array.isArray(res) ? res : []
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

const resolver = ref(zodResolver(
  z.object({
    nombre: z.string().min(1, { message: 'Nombre requerido.' }),
    descripcion: z.string().min(1, { message: 'Descripcion requerida.' }),
    estado: z.enum(['Activo', 'Inactivo'], { message: 'Debe seleccionar un estado válido.' })
  })
))

async function onFormSubmit({ valid } : any ) {
  if( valid ){
    const nombreCategoria = initialValues.nombre.trim().toLowerCase()
    const categoriaExistente = categorias.value.find(c => c.nombre?.toLowerCase() === nombreCategoria && c.id !== props.id)

    if (categoriaExistente) {
      nombreExistente.value = true
      return
    }

    nombreExistente.value = false
    try{
      await $fetch(server.HOST + '/api/v1/categorias-vehiculos/' + props.id, {
        method: 'PUT',
        headers: {
          "Content-Type" : "application/json"
        },
        body: initialValues
      })
      emit('update'), emit('success')
      visible.value = false
    } catch(err) {
      console.error(err)
      emit('error')
    }
  }
}
</script>
