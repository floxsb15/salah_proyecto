<template>
  <Dialog v-model:visible="visible" modal header="Modificar Vehiculo" :style="{ width: '56rem' }">
    <Form
      :key="formKey"
      v-slot="$form"
      :resolver="resolver"
      :initialValues="initialValues"
      @submit="onFormSubmit"
      class="flex max-h-[78vh] w-full flex-col gap-4 overflow-y-auto pr-1"
    >
      <div class="grid grid-cols-2 gap-2 rounded-md bg-gray-100 p-1 text-sm font-medium text-gray-600">
        <button
          type="button"
          class="rounded px-3 py-2 text-center transition-colors"
          :class="pasoActual === 1 ? 'bg-white text-gray-900 shadow-sm' : 'hover:bg-white/70'"
          @click="pasoActual = 1"
        >
          Informacion del vehiculo
        </button>
        <button
          type="button"
          class="rounded px-3 py-2 text-center transition-colors"
          :class="pasoActual === 2 ? 'bg-white text-gray-900 shadow-sm' : 'hover:bg-white/70'"
          @click="pasoActual = 2"
        >
          Especificaciones tecnicas
        </button>
      </div>

      <div v-show="pasoActual === 1" class="grid grid-cols-1 gap-4 lg:grid-cols-[1fr_16rem]">
        <div class="grid grid-cols-1 gap-3 md:grid-cols-2">
          <div class="flex flex-col gap-1">
            <label for="marca"> Marca </label>
            <Select
              id="marca"
              name="marca"
              v-model="initialValues.marca"
              :options="MarcasDisponibles"
              option-label="nombre"
              option-value="nombre"
              placeholder="Seleccione marca"
              filter
              fluid
              size="small"
            />
            <Message v-if="$form.marca?.invalid" severity="error" size="small" variant="simple">
              {{ $form.marca.error?.message }}
            </Message>
          </div>

          <div class="flex flex-col gap-1">
            <label for="modelo"> Modelo </label>
            <InputText id="modelo" name="modelo" v-model="initialValues.modelo" placeholder="Ej. Minibus" fluid size="small" />
            <Message v-if="$form.modelo?.invalid" severity="error" size="small" variant="simple">
              {{ $form.modelo.error?.message }}
            </Message>
          </div>

          <div class="flex flex-col gap-1">
            <label for="anio"> Anio </label>
            <Select
              id="anio"
              name="anio"
              v-model="initialValues.anio"
              :options="AniosDisponibles"
              option-label="valor"
              option-value="valor"
              placeholder="Seleccione anio"
              fluid
              size="small"
            />
            <Message v-if="$form.anio?.invalid" severity="error" size="small" variant="simple">
              {{ $form.anio.error?.message }}
            </Message>
          </div>

          <div class="flex flex-col gap-1">
            <label for="version"> Version </label>
            <InputText id="version" name="version" v-model="initialValues.version" placeholder="Ej. 4x2 Standard" fluid size="small" />
          </div>

          <div class="flex flex-col gap-1">
            <label for="categoria"> Categoria </label>
            <Select
              id="categoria"
              name="id_categoria"
              :options="Categorias"
              v-model="initialValues.id_categoria"
              placeholder="Seleccione una categoria"
              option-label="nombre"
              option-value="id"
              fluid
              size="small"
            />
            <Message v-if="$form.id_categoria?.invalid" severity="error" size="small" variant="simple">
              {{ $form.id_categoria.error?.message }}
            </Message>
          </div>

          <div class="flex flex-col gap-1">
            <label for="segmento"> Segmento </label>
            <Select
              id="segmento"
              name="segmento"
              :options="SegmentosFiltrados"
              v-model="initialValues.id_segmento"
              placeholder="Seleccione un segmento"
              option-label="nombre"
              option-value="id"
              show-clear
              fluid
              size="small"
              :disabled="!initialValues.id_categoria"
            />
          </div>

          <div class="flex flex-col gap-1">
            <label for="precio"> Precio </label>
            <InputNumber
              id="precio"
              name="precio"
              v-model="initialValues.precio"
              placeholder="Ingrese un precio"
              fluid
              size="small"
            />
            <Message v-if="$form.precio?.invalid" severity="error" size="small" variant="simple">
              {{ $form.precio.error?.message }}
            </Message>
          </div>

          <div v-if="esAdmin" class="flex flex-col gap-1">
            <label for="precio_compra"> Precio compra </label>
            <InputNumber
              id="precio_compra"
              name="precio_compra"
              v-model="initialValues.precio_compra"
              mode="currency"
              currency="USD"
              locale="es-BO"
              :min="0"
              fluid
              size="small"
            />
            <Message v-if="$form.precio_compra?.invalid" severity="error" size="small" variant="simple">
              {{ $form.precio_compra.error?.message }}
            </Message>
          </div>

          <div v-if="esAdmin" class="flex flex-col gap-1">
            <label for="margen_ganancia"> Margen ganancia </label>
            <InputNumber
              id="margen_ganancia"
              :model-value="margenGanancia"
              mode="currency"
              currency="USD"
              locale="es-BO"
              disabled
              fluid
              size="small"
            />
          </div>

          <div class="flex flex-col gap-1">
            <label for="cantidad_disponible"> Cantidad disponible </label>
            <InputNumber
              id="cantidad_disponible"
              name="cantidad_disponible"
              v-model="initialValues.cantidad_disponible"
              :min="0"
              :useGrouping="false"
              fluid
              size="small"
            />
            <Message v-if="$form.cantidad_disponible?.invalid" severity="error" size="small" variant="simple">
              {{ $form.cantidad_disponible.error?.message }}
            </Message>
          </div>

          <div class="flex flex-col gap-1">
            <label for="estado"> Estado </label>
            <Select
              id="estado"
              name="estado"
              :options="Estados"
              v-model="initialValues.estado"
              placeholder="Seleccione un estado"
              fluid
              size="small"
            />
            <Message v-if="$form.estado?.invalid" severity="error" size="small" variant="simple">
              {{ $form.estado.error?.message }}
            </Message>
          </div>
        </div>

        <div class="flex flex-col gap-2">
          <label for="fotos"> Fotos </label>
          <div class="flex min-h-60 flex-col items-center justify-center gap-3 rounded-md border border-dashed border-gray-300 bg-gray-50 p-3">
            <div v-if="srcs.length" class="grid w-full grid-cols-2 gap-2">
              <img
                v-for="(src, index) in srcs"
                :key="index"
                :src="src"
                alt="Imagen del vehiculo"
                class="h-24 w-full rounded-md object-cover shadow-sm"
              />
            </div>
            <i v-else class="pi pi-car text-5xl text-gray-400"></i>
            <FileUpload
              id="fotos"
              name="fotos"
              mode="basic"
              @select="onFileSelect"
              severity="secondary"
              class="p-button-outlined"
              choose-label="Cambiar imagenes"
              accept="image/*"
              :multiple="true"
              :file-limit="5"
              custom-upload
              auto
            />
            <small class="text-center text-xs text-gray-500">Puede subir de 1 a 5 fotos. Si selecciona nuevas, reemplazan las actuales.</small>
            <Message v-if="$form.fotos?.invalid" severity="error" size="small" variant="simple">
              {{ $form.fotos.error?.message }}
            </Message>
          </div>
        </div>
      </div>

      <div v-show="pasoActual === 2" class="grid grid-cols-1 gap-3 md:grid-cols-2">
        <div class="flex flex-col gap-1">
          <label for="tipo_techo"> Tipo de techo </label>
          <Select id="tipo_techo" name="tipo_techo" v-model="initialValues.tipo_techo" :options="TiposTecho" placeholder="Seleccione tipo de techo" show-clear fluid size="small" />
        </div>

        <div v-if="initialValues.tipo_techo === 'Otro'" class="flex flex-col gap-1">
          <label for="tipo_techo_otro"> Especifique tipo de techo </label>
          <InputText id="tipo_techo_otro" v-model="tipoTechoOtro" placeholder="Ingrese tipo de techo" fluid size="small" />
        </div>

        <div class="flex flex-col gap-1">
          <label for="combustible"> Combustible </label>
          <Select id="combustible" name="combustible" v-model="initialValues.combustible" :options="Combustibles" placeholder="Seleccione combustible" show-clear fluid size="small" />
        </div>

        <div class="flex flex-col gap-1">
          <label for="traccion"> Traccion </label>
          <Select id="traccion" name="traccion" v-model="initialValues.traccion" :options="Tracciones" placeholder="Seleccione traccion" show-clear fluid size="small" />
        </div>

        <div class="flex flex-col gap-1">
          <label for="transmision"> Transmision </label>
          <Select id="transmision" name="transmision" v-model="initialValues.transmision" :options="Transmisiones" placeholder="Seleccione transmision" show-clear fluid size="small" />
        </div>

        <div class="flex flex-col gap-1">
          <label for="asientos"> Asientos </label>
          <InputNumber id="asientos" name="asientos" v-model="initialValues.asientos" :min="1" placeholder="Cantidad" fluid size="small" />
        </div>

        <div class="flex flex-col gap-1">
          <label for="garantia_anios"> Garantia en anos </label>
          <InputNumber id="garantia_anios" v-model="garantiaAnios" :min="0" placeholder="Ej. 3" fluid size="small" />
        </div>

        <div class="flex flex-col gap-1">
          <label for="garantia_km"> Garantia en km recorridos </label>
          <InputNumber id="garantia_km" v-model="garantiaKm" :min="0" :useGrouping="false" placeholder="Ej. 100000" fluid size="small" />
        </div>

        <div class="flex flex-col gap-1 md:col-span-2">
          <label for="equipamiento"> Equipamiento </label>
          <textarea
            id="equipamiento"
            name="equipamiento"
            v-model="initialValues.equipamiento"
            rows="5"
            placeholder="Camara de retroceso, pantalla tactil, aire acondicionado..."
            class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm outline-none focus:border-yellow-400"
          ></textarea>
        </div>
      </div>

      <div class="flex justify-end gap-2 border-t border-gray-100 pt-3">
        <Button v-if="pasoActual === 2" type="button" label="Atras" severity="secondary" @click="pasoActual = 1" />
        <Button v-if="pasoActual === 1" type="button" label="Siguiente" icon="pi pi-arrow-right" icon-pos="right" @click="pasoActual = 2" />
        <Button v-else type="submit" label="Modificar" icon="pi pi-check" />
      </div>
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

const formKey = ref(0)
const pasoActual = ref(1)
const Estados = ref(['Activo', 'Inactivo'])
const TiposTecho = ref(['Techo rigido', 'Techo solar', 'Techo panoramico', 'Otro'])
const Combustibles = ref(['Gasolina', 'Diesel', 'GNV', 'GLP', 'Hibrido', 'Electrico'])
const Tracciones = ref(['4x2', '4x4', 'AWD', 'FWD', 'RWD'])
const Transmisiones = ref(['Manual', 'Automatica', 'CVT', 'Electrica'])
const srcs = ref<string[]>([])
const tipoTechoOtro = ref('')
const garantiaAnios = ref<number | null>(null)
const garantiaKm = ref<number | null>(null)

const initialValues = reactive({
  precio: 0,
  precio_compra: 0,
  cantidad_disponible: 0,
  id_categoria: 0,
  id_segmento: null as number | null,
  marca: '',
  modelo: '',
  anio: 0,
  version: '',
  tipo_techo: '',
  combustible: '',
  traccion: '',
  transmision: '',
  asientos: null as number | null,
  garantia: '',
  equipamiento: '',
  estado: '',
  fotos: [] as Array<File | string>
})
const Categorias = ref<any[]>([])
const Segmentos = ref<any[]>([])
const Marcas = ref<any[]>([])
const Anios = ref<any[]>([])
const SegmentosFiltrados = computed(() => {
  return Segmentos.value.filter(segmento => Number(segmento.id_categoria) === Number(initialValues.id_categoria))
})
const MarcasDisponibles = computed(() => {
  const marcas = Marcas.value.filter(marca => marca.estado === 'Activo')
  if (initialValues.marca && !marcas.some(marca => marca.nombre === initialValues.marca)) {
    return [{ id: 0, nombre: initialValues.marca, estado: 'Activo' }, ...marcas]
  }
  return marcas
})
const AniosDisponibles = computed(() => {
  const anios = Anios.value.filter(anio => anio.estado === 'Activo')
  if (initialValues.anio && !anios.some(anio => Number(anio.valor) === Number(initialValues.anio))) {
    return [{ id: 0, valor: initialValues.anio, estado: 'Activo' }, ...anios]
  }
  return anios
})
const esAdmin = computed(() => getActorUserRol() === 'admin')
const margenGanancia = computed(() => Number(initialValues.precio || 0) - Number(initialValues.precio_compra || 0))

onMounted( async () => {
  try {
    const resValue:any = await $fetch(server.HOST + '/api/v1/vehiculos/' + props.id, {
      method: 'GET',
      headers: getAuthHeaders()
    })
    initialValues.precio = resValue.precio
    initialValues.precio_compra = Number(resValue.precio_compra || 0)
    initialValues.cantidad_disponible = Number(resValue.cantidad_disponible || 0)
    initialValues.id_categoria = resValue.id_categoria
    initialValues.marca = resValue.marca || ''
    initialValues.modelo = resValue.modelo || ''
    initialValues.anio = resValue.anio || 0
    initialValues.version = resValue.version || ''
    setTipoTecho(resValue.tipo_techo || '')
    initialValues.combustible = resValue.combustible || ''
    initialValues.traccion = resValue.traccion || ''
    initialValues.transmision = resValue.transmision || ''
    initialValues.asientos = resValue.asientos || null
    initialValues.garantia = resValue.garantia || ''
    setGarantia(initialValues.garantia)
    initialValues.equipamiento = resValue.equipamiento || ''
    initialValues.estado = resValue.estado
    srcs.value = Array.isArray(resValue.imagenes) && resValue.imagenes.length
      ? resValue.imagenes
      : (resValue.imagen && resValue.imagen !== 'N/A' ? [resValue.imagen] : [])
    initialValues.fotos = [...srcs.value]
    const categoria = resValue.categoria

    const [resCategorias, resSegmentos, resMarcas, resAnios]: any[] = await Promise.all([
      $fetch(server.HOST + '/api/v1/categorias-vehiculos', { method: 'GET' }),
      $fetch(server.HOST + '/api/v1/segmentos-vehiculos', { method: 'GET' }),
      $fetch(server.HOST + '/api/v1/marcas-vehiculos', { method: 'GET' }),
      $fetch(server.HOST + '/api/v1/anios-vehiculos', { method: 'GET' })
    ])
    Categorias.value = Array.isArray(resCategorias) ? resCategorias : []
    Segmentos.value = Array.isArray(resSegmentos) ? resSegmentos : []
    Marcas.value = Array.isArray(resMarcas) ? resMarcas : []
    Anios.value = Array.isArray(resAnios) ? resAnios : []

    const match = Categorias.value.find(cat => cat.nombre === categoria)
    initialValues.id_categoria = match ? match.id : 0
    initialValues.id_segmento = resValue.id_segmento || null
    formKey.value += 1
  } catch(err) {
    console.error(err)
  }
})

watch(visible, (newValue) => {
  if (!newValue) { emit('close') }
})

watch(() => initialValues.id_categoria, () => {
  if (Segmentos.value.length === 0) {
    return
  }
  if (!SegmentosFiltrados.value.some(segmento => Number(segmento.id) === Number(initialValues.id_segmento))) {
    initialValues.id_segmento = null
  }
})

watch(() => initialValues.tipo_techo, (newValue) => {
  if (newValue !== 'Otro') {
    tipoTechoOtro.value = ''
  }
})

const resolver = ref(zodResolver(
  z.object({
    precio: z.number(),
    precio_compra: z.number().min(0, { message: 'Precio de compra no valido.' }),
    cantidad_disponible: z.number().min(0, { message: 'Cantidad no valida.' }),
    modelo: z.string().min(1, { message: 'Modelo requerido.' }),
    marca: z.string().min(1, { message: 'Marca requerida.' }),
    anio: z.number().min(1900, { message: 'Anio requerido.' }),
    id_categoria: z.number().min(1, { message: 'Categoria requerida.' }),
    estado: z.enum(['Activo', 'Inactivo'], { message: 'Estado requerido.' }),
    fotos: z.array(
      z.union([
        z.string(),
        z.instanceof(File, { message: 'Debe seleccionar imagenes.' })
          .refine((file) => file.type.startsWith('image/'), {
            message: 'Todos los archivos deben ser imagenes validas.',
          })
          .refine((file) => file.size <= 2 * 1024 * 1024, {
            message: 'Cada imagen no debe superar los 2MB.',
          })
      ])
    ).max(5, { message: 'Solo se permiten hasta 5 fotos.' }).optional()
  })
))

async function onFormSubmit({ valid } : any ) {
  if( valid ){
    const formData = new FormData()
    try{
      formData.append("precio", String(initialValues.precio))
      formData.append("cantidad_disponible", String(initialValues.cantidad_disponible))
      formData.append("marca", initialValues.marca)
      formData.append("modelo", initialValues.modelo)
      formData.append("anio", String(initialValues.anio))
      formData.append("version", initialValues.version)
      formData.append("tipo_techo", obtenerTipoTecho())
      formData.append("combustible", initialValues.combustible || "")
      formData.append("traccion", initialValues.traccion || "")
      formData.append("transmision", initialValues.transmision || "")
      formData.append("asientos", initialValues.asientos ? String(initialValues.asientos) : "")
      formData.append("garantia", obtenerGarantia())
      formData.append("equipamiento", initialValues.equipamiento)
      formData.append("id_categoria", String(initialValues.id_categoria))
      formData.append("id_segmento", initialValues.id_segmento ? String(initialValues.id_segmento) : "")
      formData.append("estado", initialValues.estado)
      if (esAdmin.value) {
        formData.append("precio_compra", String(initialValues.precio_compra || 0))
      }
      const fotosNuevas = initialValues.fotos.filter((foto): foto is File => foto instanceof File)
      if (fotosNuevas.length > 0) {
        fotosNuevas.forEach((foto) => {
          formData.append("fotos", foto)
        })
      }

      await $fetch(server.HOST + '/api/v1/vehiculos/' + props.id, {
        method: 'PUT',
        body: formData,
        headers: getAuthHeaders()
      })

      emit('update'), emit('success')
      visible.value = false
    } catch(err: any) {
      console.error(err)
      emit('error', err)
    }
  } else {
    pasoActual.value = 1
  }
}

function onFileSelect(event: any) {
  const files = (event.files || []).slice(0, 5)
  initialValues.fotos = files
  srcs.value = []

  files.forEach((file: File) => {
    const reader = new FileReader()
    reader.onload = (e) => {
      if (typeof e.target?.result === 'string') {
        srcs.value.push(e.target.result)
      }
    }
    reader.readAsDataURL(file)
  })
}

function setTipoTecho(tipoTecho: string) {
  if (!tipoTecho || TiposTecho.value.includes(tipoTecho)) {
    initialValues.tipo_techo = tipoTecho
    tipoTechoOtro.value = ''
    return
  }

  initialValues.tipo_techo = 'Otro'
  tipoTechoOtro.value = tipoTecho
}

function obtenerTipoTecho() {
  return initialValues.tipo_techo === 'Otro' ? tipoTechoOtro.value.trim() : initialValues.tipo_techo
}

function setGarantia(garantia: string) {
  const aniosMatch = garantia.match(/(\d+)\s*(?:anos|anios|a[ñn]os|anio|a[ñn]o)/i)
  const kmMatch = garantia.match(/([\d.,]+)\s*km/i)

  garantiaAnios.value = aniosMatch ? Number(aniosMatch[1]) : null
  garantiaKm.value = kmMatch ? Number(kmMatch[1].replace(/\./g, '').replace(',', '.')) : null
}

function obtenerGarantia() {
  const partes = []
  if (garantiaAnios.value !== null) {
    partes.push(`${garantiaAnios.value} anos`)
  }
  if (garantiaKm.value !== null) {
    partes.push(`${garantiaKm.value} km`)
  }
  return partes.length > 0 ? partes.join(' o ') : initialValues.garantia
}

function getAuthHeaders() {
  try {
    const user = localStorage.getItem('user')
    const parsed = user ? JSON.parse(user) : null
    const headers: Record<string, string> = {}
    if (parsed?.token) {
      headers.Authorization = `Bearer ${parsed.token}`
    }
    if (parsed?.id) {
      headers['X-User-Id'] = String(parsed.id)
    }
    return headers
  } catch {
    return {}
  }
}

function getActorUserRol() {
  try {
    const user = localStorage.getItem('user')
    return user ? JSON.parse(user)?.rol || '' : ''
  } catch {
    return ''
  }
}
</script>
