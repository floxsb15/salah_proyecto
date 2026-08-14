<template>
  <Dialog
    :visible="visible"
    modal
    header="Generar proforma vehicular"
    :style="{ width: 'min(980px, 96vw)' }"
    @update:visible="emitClose"
  >
    <div v-if="vehiculo" class="flex flex-col gap-5">
      <section class="grid grid-cols-1 gap-4 lg:grid-cols-[1.1fr_0.9fr]">
        <div class="flex flex-col gap-3">
          <div class="flex items-center justify-between gap-3">
            <h3 class="text-base font-semibold text-gray-900">Informacion del cliente</h3>
            <SelectButton v-model="modoCliente" :options="modosCliente" option-label="label" option-value="value" size="small" />
          </div>

          <div v-if="modoCliente === 'existente'" class="field-group">
            <label for="cliente_proforma">Cliente</label>
            <Select
              id="cliente_proforma"
              v-model="form.id_cliente"
              :options="clientes"
              option-label="nombreCompleto"
              option-value="id"
              placeholder="Buscar cliente"
              filter
              show-clear
              fluid
              size="small"
              :loading="loadingClientes"
            />
            <small v-if="errors.id_cliente" class="text-red-600">{{ errors.id_cliente }}</small>
          </div>

          <div class="grid grid-cols-1 gap-3 md:grid-cols-2">
            <div class="field-group md:col-span-2">
              <label for="cliente_nombre">Nombre</label>
              <InputText id="cliente_nombre" v-model.trim="form.cliente_nombre" size="small" :disabled="modoCliente === 'existente' && !!form.id_cliente" />
              <small v-if="errors.cliente_nombre" class="text-red-600">{{ errors.cliente_nombre }}</small>
            </div>
            <div class="field-group">
              <label for="cliente_telefono">Telefono</label>
              <InputText id="cliente_telefono" v-model.trim="form.cliente_telefono" size="small" />
              <small v-if="errors.cliente_telefono" class="text-red-600">{{ errors.cliente_telefono }}</small>
            </div>
            <div class="field-group">
              <label for="cliente_direccion">Direccion</label>
              <InputText id="cliente_direccion" v-model.trim="form.cliente_direccion" size="small" />
            </div>
          </div>
        </div>

        <div class="rounded-lg border border-gray-200 bg-gray-50 p-4">
          <h3 class="text-base font-semibold text-gray-900">Vehiculo</h3>
          <dl class="mt-3 grid grid-cols-2 gap-x-4 gap-y-2 text-sm">
            <dt class="text-gray-500">Marca</dt>
            <dd class="font-medium text-gray-900">{{ vehiculo.marca || 'N/A' }}</dd>
            <dt class="text-gray-500">Modelo</dt>
            <dd class="font-medium text-gray-900">{{ vehiculo.modelo || 'N/A' }}</dd>
            <dt class="text-gray-500">Anio</dt>
            <dd class="font-medium text-gray-900">{{ vehiculo.anio || 'N/A' }}</dd>
            <dt class="text-gray-500">Garantia</dt>
            <dd class="font-medium text-gray-900">{{ vehiculo.garantia || 'N/A' }}</dd>
            <dt class="text-gray-500">Equipamiento</dt>
            <dd class="font-medium text-gray-900">{{ vehiculo.equipamiento || 'N/A' }}</dd>
            <dt class="text-gray-500">Tecnicas</dt>
            <dd class="font-medium text-gray-900">{{ especificacionesVehiculo }}</dd>
          </dl>
        </div>
      </section>

      <section class="grid grid-cols-1 gap-3 md:grid-cols-3">
        <div class="field-group">
          <label for="modalidad">Modalidad</label>
          <Select id="modalidad" v-model="form.modalidad" :options="modalidades" size="small" fluid />
        </div>
        <div class="field-group">
          <label for="precio_unidad">Precio unidad USD</label>
          <InputNumber
            id="precio_unidad"
            v-model="form.precio_unidad"
            mode="currency"
            currency="USD"
            locale="en-US"
            :min="precioCatalogo"
            :min-fraction-digits="2"
            :max-fraction-digits="2"
            fluid
            size="small"
          />
          <small v-if="errors.precio_unidad" class="text-red-600">{{ errors.precio_unidad }}</small>
          <small v-else class="text-gray-500">Minimo catalogo: $ {{ formatPrecio(precioCatalogo) }}</small>
        </div>
        <div class="field-group">
          <label for="cantidad">Cantidad</label>
          <InputNumber id="cantidad" v-model="form.cantidad" :min="1" :max="cantidadMaxima" show-buttons fluid size="small" />
          <small v-if="errors.cantidad" class="text-red-600">{{ errors.cantidad }}</small>
        </div>
        <div class="field-group">
          <label>Precio total</label>
          <div class="readonly-money">$ {{ formatPrecio(precioTotal) }}</div>
        </div>
        <div class="field-group">
          <label for="cuota_inicial">Cuota inicial USD</label>
          <InputNumber
            id="cuota_inicial"
            v-model="form.cuota_inicial"
            mode="currency"
            currency="USD"
            locale="en-US"
            :min="0"
            :max="precioTotal"
            :min-fraction-digits="2"
            :max-fraction-digits="2"
            fluid
            size="small"
          />
          <small v-if="errors.cuota_inicial" class="text-red-600">{{ errors.cuota_inicial }}</small>
        </div>
        <div class="field-group">
          <label>Saldo</label>
          <div class="readonly-money">$ {{ formatPrecio(saldo) }}</div>
        </div>
        <div class="field-group">
          <label for="validez_dias">Validez</label>
          <InputNumber id="validez_dias" v-model="form.validez_dias" suffix=" dias" :min="1" :max="90" show-buttons fluid size="small" />
          <small v-if="errors.validez_dias" class="text-red-600">{{ errors.validez_dias }}</small>
        </div>
      </section>

      <Message v-if="errors.general" severity="error" size="small">{{ errors.general }}</Message>
    </div>

    <template #footer>
      <Button label="Cancelar" icon="pi pi-times" severity="secondary" text @click="emitClose" />
      <Button label="Generar PDF" icon="pi pi-file-pdf" :loading="saving" @click="generarProforma" />
    </template>
  </Dialog>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import InputNumber from 'primevue/inputnumber';
import InputText from 'primevue/inputtext';
import Message from 'primevue/message';
import Select from 'primevue/select';
import SelectButton from 'primevue/selectbutton';
import { useToast } from 'primevue/usetoast';
import { server } from '~/server/server';

const props = defineProps<{
  visible: boolean;
  vehiculo: any | null;
}>();

const emit = defineEmits<{
  close: [];
  generated: [proforma: any];
}>();

const toast = useToast();
const clientes = ref<any[]>([]);
const loadingClientes = ref(false);
const saving = ref(false);
const modoCliente = ref<'existente' | 'nuevo'>('existente');

const modosCliente = [
  { label: 'Cliente existente', value: 'existente' },
  { label: 'Ingresar nuevo', value: 'nuevo' }
];
const modalidades = ['Almacen', 'Pedido', 'Consignacion', 'Otro'];

const form = reactive({
  id_cliente: null as number | null,
  cliente_nombre: '',
  cliente_direccion: '',
  cliente_telefono: '',
  modalidad: 'Almacen',
  precio_unidad: 0,
  cantidad: 1,
  cuota_inicial: 0,
  validez_dias: 10
});

const errors = reactive<Record<string, string>>({});

const precioCatalogo = computed(() => Number(props.vehiculo?.precio_usd || props.vehiculo?.precio || 0));
const cantidadMaxima = computed(() => Math.max(Number(props.vehiculo?.cantidad_disponible || 1), 1));
const precioTotal = computed(() => Number(form.precio_unidad || 0) * Number(form.cantidad || 0));
const saldo = computed(() => Math.max(precioTotal.value - Number(form.cuota_inicial || 0), 0));
const especificacionesVehiculo = computed(() => {
  if (!props.vehiculo) {
    return 'N/A';
  }
  return [props.vehiculo.combustible, props.vehiculo.traccion, props.vehiculo.transmision, props.vehiculo.asientos ? `${props.vehiculo.asientos} asientos` : '']
    .filter(Boolean)
    .join(' / ') || 'N/A';
});

watch(() => props.visible, async (visible) => {
  if (!visible) {
    return;
  }
  resetForm();
  await obtenerClientes();
});

watch(() => form.id_cliente, (id) => {
  if (!id || modoCliente.value !== 'existente') {
    return;
  }
  const cliente = clientes.value.find(c => c.id === id);
  if (!cliente) {
    return;
  }
  form.cliente_nombre = cliente.nombreCliente;
  form.cliente_direccion = cliente.direccion || '';
  form.cliente_telefono = cliente.celular || '';
});

watch(modoCliente, () => {
  form.id_cliente = null;
  form.cliente_nombre = '';
  form.cliente_direccion = '';
  form.cliente_telefono = '';
});

function resetForm() {
  clearErrors();
  modoCliente.value = 'existente';
  form.id_cliente = null;
  form.cliente_nombre = '';
  form.cliente_direccion = '';
  form.cliente_telefono = '';
  form.modalidad = 'Almacen';
  form.precio_unidad = precioCatalogo.value;
  form.cantidad = 1;
  form.cuota_inicial = 0;
  form.validez_dias = 10;
}

async function obtenerClientes() {
  loadingClientes.value = true;
  try {
    const res: any = await $fetch(server.HOST + '/api/v1/clientes', { method: 'GET' });
    const activos = Array.isArray(res) ? res.filter((cliente: any) => cliente.estado === 'Activo') : [];
    clientes.value = activos.map((cliente: any) => ({
      ...cliente,
      nombreCliente: `${cliente.nombre || ''} ${cliente.apellido || ''}`.trim() || 'Sin nombre',
      nombreCompleto: `${cliente.nombre || ''} ${cliente.apellido || ''} - ${cliente.ci || 'Sin CI'}`.trim()
    }));
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar clientes', life: 3000 });
  } finally {
    loadingClientes.value = false;
  }
}

function validarFormulario() {
  clearErrors();
  if (modoCliente.value === 'existente' && !form.id_cliente) {
    errors.id_cliente = 'Seleccione un cliente';
  }
  if (!form.cliente_nombre.trim()) {
    errors.cliente_nombre = 'Nombre requerido';
  }
  if (form.cliente_telefono && !/^[0-9+\-\s()]{6,30}$/.test(form.cliente_telefono)) {
    errors.cliente_telefono = 'Telefono no valido';
  }
  if (Number(form.precio_unidad || 0) < precioCatalogo.value) {
    errors.precio_unidad = 'No puede ser menor al precio del catalogo';
  }
  if (Number(form.cantidad || 0) < 1) {
    errors.cantidad = 'Cantidad requerida';
  }
  if (Number(form.cantidad || 0) > cantidadMaxima.value) {
    errors.cantidad = 'No puede superar el stock disponible';
  }
  if (Number(form.cuota_inicial || 0) < 0 || Number(form.cuota_inicial || 0) > precioTotal.value) {
    errors.cuota_inicial = 'La cuota debe estar entre 0 y el precio total';
  }
  if (Number(form.validez_dias || 0) < 1) {
    errors.validez_dias = 'Validez requerida';
  }

  return Object.keys(errors).length === 0;
}

async function generarProforma() {
  if (!props.vehiculo || !validarFormulario()) {
    return;
  }

  saving.value = true;
  try {
    const proforma: any = await $fetch(server.HOST + '/api/v1/proformas-vehiculares', {
      method: 'POST',
      body: {
        id_cliente: modoCliente.value === 'existente' ? form.id_cliente : null,
        id_vehiculo: props.vehiculo.id,
        cliente_nombre: form.cliente_nombre,
        cliente_direccion: form.cliente_direccion,
        cliente_telefono: form.cliente_telefono,
        modalidad: form.modalidad,
        precio_unidad: Number(form.precio_unidad || 0),
        cantidad: Number(form.cantidad || 1),
        cuota_inicial: Number(form.cuota_inicial || 0),
        validez_dias: Number(form.validez_dias || 10)
      }
    });

    await descargarPDF(proforma.id);
    toast.add({ severity: 'success', summary: 'Proforma generada', life: 3000 });
    emit('generated', proforma);
    emitClose();
  } catch (err: any) {
    const message = err?.data?.message || err?.message || 'No se pudo generar la proforma';
    errors.general = message;
    toast.add({ severity: 'error', summary: message, life: 4000 });
  } finally {
    saving.value = false;
  }
}

async function descargarPDF(idProforma: number | string) {
  const response = await fetch(`${server.HOST}/api/v1/reportes/proformas-vehiculares/${idProforma}`, {
    method: 'GET'
  });

  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || 'No se pudo generar el PDF');
  }

  const blob = await response.blob();
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = `proforma_vehicular_${idProforma}.pdf`;
  document.body.appendChild(link);
  link.click();
  link.remove();
  window.URL.revokeObjectURL(url);
}

function clearErrors() {
  Object.keys(errors).forEach((key) => {
    delete errors[key];
  });
}

function emitClose() {
  emit('close');
}

function formatPrecio(precio: number) {
  return Number(precio || 0).toLocaleString('es-BO', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  });
}
</script>

<style scoped>
.field-group {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.field-group label {
  font-size: 0.82rem;
  font-weight: 600;
  color: #374151;
}

.readonly-money {
  min-height: 2.35rem;
  display: flex;
  align-items: center;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  background: #f9fafb;
  padding: 0.35rem 0.65rem;
  font-weight: 700;
  color: #111827;
}
</style>
