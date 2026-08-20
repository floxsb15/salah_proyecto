<template>
  <template v-if="mostrarDetallePago">
    <div class="flex flex-col gap-1">
      <label for="numero_cuotas">Numero de cuotas</label>
      <InputNumber id="numero_cuotas" v-model="form.numero_cuotas" :min="1" :useGrouping="false" fluid size="small" />
    </div>

    <div class="flex flex-col gap-1">
      <label for="fecha_inicio_credito">Fecha inicio pagos</label>
      <InputText id="fecha_inicio_credito" v-model="form.fecha_inicio_credito" type="date" fluid size="small" />
    </div>

    <div class="flex flex-col gap-1">
      <label for="frecuencia_pago">Frecuencia</label>
      <Select id="frecuencia_pago" v-model="form.frecuencia_pago" :options="frecuenciasPago" fluid size="small" />
    </div>

    <div class="grid grid-cols-1 gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2 md:grid-cols-4">
      <div>
        <span class="text-xs font-semibold uppercase text-gray-500">Monto financiado</span>
        <strong class="block text-gray-900">$ {{ formatPrecio(montoFinanciado) }}</strong>
      </div>
      <div>
        <span class="text-xs font-semibold uppercase text-gray-500">Monto por cuota</span>
        <strong class="block text-gray-900">$ {{ formatPrecio(montoCuota) }}</strong>
      </div>
      <div>
        <span class="text-xs font-semibold uppercase text-gray-500">Resta</span>
        <strong class="block text-gray-900">$ {{ formatPrecio(saldo) }}</strong>
      </div>
      <div>
        <span class="text-xs font-semibold uppercase text-gray-500">Pagado inicial</span>
        <strong class="block text-gray-900">$ {{ formatPrecio(pagoDetalleUSD) }}</strong>
      </div>
    </div>

    <section class="flex flex-col gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2">
      <div class="flex items-center justify-between gap-3">
        <div>
          <h3 class="text-sm font-semibold text-gray-900">Detalle de pago inicial</h3>
          <p class="text-xs text-gray-500">Pagado: $ {{ formatPrecio(pagoDetalleUSD) }} / Bs {{ formatPrecio(pagoDetalleBOB) }}</p>
        </div>
        <Button label="Agregar pago" icon="pi pi-plus" size="small" severity="secondary" type="button" @click="agregarPago" />
      </div>

      <div v-for="(pago, index) in form.pagos" :key="pago.key" class="grid grid-cols-1 gap-2 md:grid-cols-[9rem_13rem_1fr_2.5rem]">
        <Select v-model="pago.moneda" :options="monedasPago" placeholder="Moneda" fluid size="small" />
        <Select v-model="pago.metodo" :options="metodosPago" placeholder="Metodo" fluid size="small" />
        <InputNumber v-model="pago.monto" mode="decimal" :min="0" :minFractionDigits="2" :maxFractionDigits="2" placeholder="Monto" fluid size="small" />
        <Button icon="pi pi-trash" severity="danger" text rounded type="button" aria-label="Eliminar pago" @click="eliminarPago(index)" />
      </div>

      <p v-if="pagoExcedeTotal" class="text-sm font-semibold text-red-600">El pagado supera el total de la venta.</p>
    </section>

    <div class="grid grid-cols-1 gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2 md:grid-cols-3">
      <div><span class="text-xs font-semibold uppercase text-gray-500">Pagado equivalente USD</span><strong class="block text-gray-900">$ {{ formatPrecio(pagoDetalleUSD) }}</strong></div>
      <div><span class="text-xs font-semibold uppercase text-gray-500">Saldo USD</span><strong class="block text-gray-900">$ {{ formatPrecio(saldoDetalleUSD) }}</strong></div>
      <div><span class="text-xs font-semibold uppercase text-gray-500">Saldo BOB</span><strong class="block text-gray-900">Bs {{ formatPrecio(saldoDetalleBOB) }}</strong></div>
    </div>
  </template>

  <div class="flex flex-col gap-1">
    <label for="estado_entrega">Estado de entrega</label>
    <Select id="estado_entrega" v-model="form.estado_entrega" :options="estadosEntrega" :disabled="deshabilitarEntrega" fluid size="small" />
  </div>

  <div class="flex flex-col gap-1">
    <label for="fecha_entrega">Fecha de entrega</label>
    <InputText id="fecha_entrega" v-model="form.fecha_entrega" type="date" fluid size="small" />
  </div>

  <template v-if="esCreditoBancario">
    <div class="flex flex-col gap-1">
      <label for="referencia_bancaria">Referencia bancaria</label>
      <InputText id="referencia_bancaria" v-model="form.referencia_bancaria" placeholder="Texto libre" fluid size="small" />
    </div>
    <div class="flex flex-col gap-1">
      <label for="estado_desembolso">Estado desembolso</label>
      <Select id="estado_desembolso" v-model="form.estado_desembolso" :options="estadosDesembolso" fluid size="small" />
    </div>
  </template>

  <template v-if="esCreditoDirecto">
    <div class="flex items-center gap-2">
      <Checkbox id="tiene_respaldo" v-model="form.tiene_respaldo" binary />
      <label for="tiene_respaldo">Cliente con respaldo/garantia</label>
    </div>
    <div class="flex flex-col gap-1">
      <label for="tipo_garantia">Tipo de garantia</label>
      <InputText id="tipo_garantia" v-model="form.tipo_garantia" placeholder="Ej. garante, inmueble, vehiculo" fluid size="small" />
    </div>
    <div class="flex flex-col gap-1 lg:col-span-2">
      <label for="documento_garantia">Documento o imagen de la garantia</label>
      <div class="flex flex-col gap-2 rounded-md border border-dashed border-gray-300 bg-gray-50 p-3">
        <FileUpload id="documento_garantia" name="documento_garantia" mode="basic" choose-label="Subir garantia" accept="image/*,application/pdf" severity="secondary" class="p-button-outlined" custom-upload auto @select="$emit('documento-select', $event)" @clear="$emit('documento-clear')" @remove="$emit('documento-clear')" />
        <small class="text-xs text-gray-500">{{ documentoGarantiaNombre || 'Formatos: imagen o PDF, maximo 10MB.' }}</small>
      </div>
    </div>
    <div class="flex flex-col gap-1 lg:col-span-2">
      <label for="datos_garante">Datos del garante/respaldo</label>
      <Textarea id="datos_garante" v-model="form.datos_garante" rows="3" auto-resize fluid />
    </div>
  </template>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import Button from 'primevue/button';
import Checkbox from 'primevue/checkbox';
import FileUpload from 'primevue/fileupload';
import InputNumber from 'primevue/inputnumber';
import InputText from 'primevue/inputtext';
import Select from 'primevue/select';
import Textarea from 'primevue/textarea';

const props = withDefaults(defineProps<{
  form: any;
  montoFinanciado: number;
  montoCuota: number;
  saldo: number;
  pagoDetalleUSD: number;
  pagoDetalleBOB: number;
  saldoDetalleUSD: number;
  saldoDetalleBOB: number;
  pagoExcedeTotal: boolean;
  documentoGarantiaNombre?: string;
  deshabilitarEntrega?: boolean;
}>(), {
  documentoGarantiaNombre: '',
  deshabilitarEntrega: false
});

defineEmits<{
  (event: 'documento-select', value: any): void;
  (event: 'documento-clear'): void;
}>();

const form = props.form;
const monedasPago = ['USD', 'BOB'];
const metodosPago = ['QR', 'Transferencia', 'Efectivo'];
const estadosEntrega = ['Pendiente', 'Entregado'];
const frecuenciasPago = ['mensual', 'quincenal', 'semanal'];
const estadosDesembolso = ['Pendiente', 'Desembolsado'];
const esCreditoDirecto = computed(() => form.tipo_venta === 'credito_directo' || form.tipo_credito === 'credito_directo');
const esCreditoBancario = computed(() => form.tipo_venta === 'credito_bancario' || form.tipo_credito === 'credito_bancario');
const mostrarDetallePago = computed(() => form.tipo_venta === 'credito_directo' || form.tipo_venta === 'credito_bancario' || form.tipo_pago === 'Credito');

function crearPago() {
  return { key: `${Date.now()}-${Math.random().toString(36).slice(2)}`, moneda: 'USD', metodo: 'Efectivo', monto: 0 };
}

function agregarPago() {
  if (!Array.isArray(form.pagos)) form.pagos = [];
  form.pagos.push(crearPago());
}

function eliminarPago(index: number) {
  form.pagos.splice(index, 1);
  if (form.pagos.length === 0) agregarPago();
}

function formatPrecio(precio: number) {
  return Number(precio || 0).toLocaleString('es-BO', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}
</script>
