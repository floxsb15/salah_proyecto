<template>
  <Toast />

  <div class="flex flex-col gap-4">
    <div class="flex items-center justify-between gap-3">
      <div>
        <h2 class="text-2xl font-bold">Completar reserva</h2>
        <p class="text-sm text-gray-500">Finaliza el pago pendiente del vehiculo reservado.</p>
      </div>
      <Button label="Volver" icon="pi pi-arrow-left" size="small" severity="secondary" @click="router.push(historialPath)" />
    </div>

    <div v-if="loading" class="grid grid-cols-1 gap-4 lg:grid-cols-2">
      <Skeleton height="18rem" />
      <Skeleton height="18rem" />
    </div>

    <div v-else-if="!reserva" class="flex flex-col items-center justify-center gap-2 py-16 text-gray-500">
      <i class="pi pi-exclamation-triangle text-4xl"></i>
      <p>No se encontro la reserva seleccionada.</p>
    </div>

    <div v-else class="grid grid-cols-1 gap-4 xl:grid-cols-[24rem_1fr]">
      <section class="rounded-lg border border-gray-200 bg-white p-4">
        <div class="flex items-start justify-between gap-3">
          <div>
            <h3 class="text-lg font-semibold text-gray-900">{{ reserva.vehiculo }}</h3>
            <p v-if="esReservaPedido" class="text-sm font-medium text-orange-600">Vehiculo a pedido / {{ reserva.estado_venta }}</p>
          </div>
          <Tag v-if="esReservaPedido" value="A pedido" severity="warning" />
        </div>
        <p class="text-sm text-gray-500">{{ reserva.cliente }} · CI/NIT: {{ reserva.ci_cliente || 'N/A' }}</p>

        <div v-if="esReservaPedido" class="mt-3 rounded-md border border-orange-200 bg-orange-50 p-3 text-sm text-gray-700">
          <p class="font-semibold text-gray-900">{{ vehiculoPedido }}</p>
          <p v-if="reserva.pedido_version" class="text-xs text-gray-600">{{ reserva.pedido_version }}</p>
          <p class="text-xs text-gray-600">{{ detallePedido }}</p>
        </div>

        <div class="mt-4 grid grid-cols-2 gap-3 text-sm">
          <div class="rounded-md bg-gray-50 p-3">
            <span class="text-gray-500">Precio total</span>
            <strong class="block text-gray-900">$ {{ formatPrecio(reserva.precio_total) }}</strong>
          </div>
          <div class="rounded-md bg-gray-50 p-3">
            <span class="text-gray-500">Reserva pagada</span>
            <strong class="block text-gray-900">$ {{ formatPrecio(reserva.cuota_inicial) }}</strong>
          </div>
          <div class="rounded-md bg-gray-50 p-3">
            <span class="text-gray-500">Saldo pendiente</span>
            <strong class="block text-gray-900">$ {{ formatPrecio(reserva.saldo) }}</strong>
          </div>
          <div class="rounded-md bg-gray-50 p-3">
            <span class="text-gray-500">Vence proforma</span>
            <strong class="block text-gray-900">{{ formatFecha(reserva.fecha_vencimiento_proforma) }}</strong>
          </div>
        </div>
      </section>

      <section class="rounded-lg border border-gray-200 bg-white p-4">
        <form class="grid grid-cols-1 gap-4 lg:grid-cols-2" @submit.prevent @keydown.enter.prevent>
          <section v-if="esReservaPedido && !reserva.id_vehiculo" class="grid grid-cols-1 gap-3 rounded-md border border-orange-200 bg-orange-50 p-3 lg:col-span-2 md:grid-cols-[1fr_auto]">
            <div class="flex flex-col gap-1">
              <label for="vehiculo_importado">Vehiculo importado en inventario</label>
              <Select id="vehiculo_importado" v-model="form.id_vehiculo" :options="vehiculos" option-label="nombreCompleto" option-value="id" placeholder="Seleccione el vehiculo creado" filter fluid size="small" />
            </div>
            <div class="flex items-end">
              <Button label="Crear vehiculo" icon="pi pi-car" severity="secondary" type="button" size="small" @click="vehiculoVisible = true" />
            </div>
          </section>

          <div class="flex flex-col gap-1 lg:col-span-2">
            <label for="tipo_pago">Tipo de pago</label>
            <Select id="tipo_pago" v-model="form.tipo_pago" :options="tiposPago" fluid size="small" />
          </div>

          <div class="flex flex-col gap-1">
            <label for="tipo_cambio">Tipo de cambio del dia</label>
            <InputNumber id="tipo_cambio" v-model="form.tipo_cambio" :min="0" :minFractionDigits="2" :maxFractionDigits="4" suffix=" Bs/USD" fluid size="small" />
          </div>
          <div class="flex flex-col gap-1">
            <label>Saldo restante</label>
            <div class="rounded-md border border-gray-200 bg-gray-50 px-3 py-2 text-sm">
              <strong>$ {{ formatPrecio(saldoRestanteUSD) }}</strong>
              <span class="ml-2 text-gray-500">Bs {{ formatPrecio(saldoRestanteBOB) }}</span>
            </div>
          </div>

          <section v-if="esContado" class="flex flex-col gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2">
            <div class="flex items-center justify-between gap-3">
              <div>
                <h3 class="text-sm font-semibold text-gray-900">Detalle de pago</h3>
                <p class="text-xs text-gray-500">Pagado: $ {{ formatPrecio(pagoEquivalenteUSD) }} / Bs {{ formatPrecio(pagoEquivalenteBOB) }}</p>
              </div>
              <Button label="Agregar pago" icon="pi pi-plus" size="small" severity="secondary" type="button" @click="agregarPago" />
            </div>
            <div v-for="(pago, index) in form.pagos" :key="pago.key" class="grid grid-cols-1 gap-2 md:grid-cols-[9rem_13rem_1fr_2.5rem]">
              <Select v-model="pago.moneda" :options="monedasPago" placeholder="Moneda" fluid size="small" />
              <Select v-model="pago.metodo" :options="metodosPago" placeholder="Tipo de pago" fluid size="small" />
              <InputNumber v-model="pago.monto" mode="decimal" :min="0" :minFractionDigits="2" :maxFractionDigits="2" placeholder="Monto" fluid size="small" />
              <Button icon="pi pi-trash" severity="danger" text rounded type="button" aria-label="Eliminar pago" @click="eliminarPago(index)" />
            </div>
            <p v-if="pagoExcedeSaldo" class="text-sm font-semibold text-red-600">El pago supera el saldo pendiente.</p>
          </section>

          <template v-else>
            <div class="flex flex-col gap-1">
              <label for="tipo_credito">Tipo de credito</label>
              <Select id="tipo_credito" v-model="form.tipo_credito" :options="tiposCredito" option-label="label" option-value="value" fluid size="small" />
            </div>
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
            <div class="grid grid-cols-1 gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2 md:grid-cols-3">
              <div><span class="text-xs font-semibold uppercase text-gray-500">Monto financiado</span><strong class="block text-gray-900">$ {{ formatPrecio(saldoRestanteUSD) }}</strong></div>
              <div><span class="text-xs font-semibold uppercase text-gray-500">Monto por cuota</span><strong class="block text-gray-900">$ {{ formatPrecio(montoCuotaCredito) }}</strong></div>
              <div><span class="text-xs font-semibold uppercase text-gray-500">Saldo Bs</span><strong class="block text-gray-900">Bs {{ formatPrecio(saldoRestanteBOB) }}</strong></div>
            </div>
            <template v-if="esCreditoBancario">
              <div class="flex flex-col gap-1"><label for="referencia_bancaria">Referencia bancaria</label><InputText id="referencia_bancaria" v-model="form.referencia_bancaria" placeholder="Texto libre" fluid size="small" /></div>
              <div class="flex flex-col gap-1"><label for="estado_desembolso">Estado desembolso</label><Select id="estado_desembolso" v-model="form.estado_desembolso" :options="estadosDesembolso" fluid size="small" /></div>
            </template>
            <template v-if="esCreditoDirecto">
              <div class="flex items-center gap-2"><Checkbox id="tiene_respaldo" v-model="form.tiene_respaldo" binary /><label for="tiene_respaldo">Cliente con respaldo/garantia</label></div>
              <div class="flex flex-col gap-1"><label for="tipo_garantia">Tipo de garantia</label><InputText id="tipo_garantia" v-model="form.tipo_garantia" placeholder="Ej. garante, inmueble, vehiculo" fluid size="small" /></div>
              <div class="flex flex-col gap-1 lg:col-span-2"><label for="documento_garantia">Documento garantia</label><InputText id="documento_garantia" v-model="form.documento_garantia" placeholder="Referencia o descripcion" fluid size="small" /></div>
              <div class="flex flex-col gap-1 lg:col-span-2"><label for="datos_garante">Datos del garante/respaldo</label><Textarea id="datos_garante" v-model="form.datos_garante" rows="3" auto-resize fluid /></div>
            </template>
          </template>

          <div class="flex flex-col gap-1">
            <label for="estado_entrega">Estado de entrega</label>
            <Select id="estado_entrega" v-model="form.estado_entrega" :options="estadosEntrega" fluid size="small" />
          </div>
          <div class="flex flex-col gap-1">
            <label for="fecha_entrega">Fecha de entrega</label>
            <InputText id="fecha_entrega" v-model="form.fecha_entrega" type="date" fluid size="small" />
          </div>

          <div class="flex flex-col gap-1 lg:col-span-2">
            <label for="observacion">Observacion</label>
            <Textarea id="observacion" v-model="form.observacion" rows="3" auto-resize fluid />
          </div>

          <div class="flex justify-end gap-2 border-t border-gray-100 pt-4 lg:col-span-2">
            <Button label="Cancelar" severity="secondary" type="button" @click="router.push(historialPath)" />
            <Button label="Completar pago" icon="pi pi-check" type="button" :loading="saving" :disabled="requiereVehiculoImportado" @click="completarReserva" />
          </div>
        </form>
      </section>
    </div>

    <modalAgregarVehiculo
      v-if="vehiculoVisible"
      :open="vehiculoVisible"
      @close="vehiculoVisible = false"
      @update="obtenerVehiculos"
      @success="toast.add({ severity: 'success', summary: 'Vehiculo agregado', life: 3000 })"
      @error="mostrarError('Error al agregar vehiculo', $event)"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { server } from '~/server/server';
import Button from 'primevue/button';
import Checkbox from 'primevue/checkbox';
import InputNumber from 'primevue/inputnumber';
import InputText from 'primevue/inputtext';
import Select from 'primevue/select';
import Skeleton from 'primevue/skeleton';
import Tag from 'primevue/tag';
import Textarea from 'primevue/textarea';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';
import modalAgregarVehiculo from '~/components/admin/vehiculos/modalAgregarProducto.vue';

const props = withDefaults(defineProps<{
  historialPath?: string;
}>(), {
  historialPath: '/ventas/historial-reservas'
});

const route = useRoute();
const router = useRouter();
const toast = useToast();
const { descargarPDFVenta } = useVentaPdf();
const reserva = ref<any>(null);
const vehiculos = ref<any[]>([]);
const loading = ref(true);
const saving = ref(false);
const vehiculoVisible = ref(false);
const monedasPago = ref(['USD', 'BOB']);
const metodosPago = ref(['QR', 'Transferencia', 'Efectivo']);
const estadosEntrega = ref(['Pendiente', 'Entregado']);
const tiposPago = ref(['Contado', 'Credito']);
const tiposCredito = ref([
  { label: 'Credito directo', value: 'credito_directo' },
  { label: 'Credito bancario', value: 'credito_bancario' }
]);
const frecuenciasPago = ref(['mensual', 'quincenal', 'semanal']);
const estadosDesembolso = ref(['Pendiente', 'Aprobado', 'Desembolsado']);
const form = reactive({
  tipo_pago: 'Contado',
  tipo_cambio: null as number | null,
  pagos: [crearPago()],
  monto_pago: 0,
  metodo_pago: 'Efectivo',
  estado_entrega: 'Pendiente',
  fecha_entrega: '',
  id_vehiculo: null as number | null,
  observacion: '',
  tipo_credito: 'credito_directo',
  numero_cuotas: 12,
  fecha_inicio_credito: new Date().toISOString().slice(0, 10),
  frecuencia_pago: 'mensual',
  tiene_respaldo: false,
  tipo_garantia: '',
  documento_garantia: '',
  datos_garante: '',
  referencia_bancaria: '',
  estado_desembolso: 'Pendiente'
});

const esReservaPedido = computed(() => reserva.value?.tipo_reserva === 'pedido');
const requiereVehiculoImportado = computed(() => esReservaPedido.value && !reserva.value?.id_vehiculo && !form.id_vehiculo);
const vehiculoPedido = computed(() => [reserva.value?.pedido_marca, reserva.value?.pedido_modelo, reserva.value?.pedido_anio].filter(Boolean).join(' ') || reserva.value?.vehiculo || 'Vehiculo a pedido');
const detallePedido = computed(() => {
  if (!reserva.value) return '';
  return [
    reserva.value.pedido_color ? `Color: ${reserva.value.pedido_color}` : '',
    reserva.value.pedido_pais_origen ? `Origen: ${reserva.value.pedido_pais_origen}` : '',
    reserva.value.pedido_llegada_estimada ? `Llegada: ${reserva.value.pedido_llegada_estimada}` : '',
    reserva.value.pedido_proveedor ? `Proveedor: ${reserva.value.pedido_proveedor}` : ''
  ].filter(Boolean).join(' / ');
});
const saldoRestanteUSD = computed(() => Number(reserva.value?.saldo || 0));
const tipoCambio = computed(() => Number(form.tipo_cambio || 0));
const saldoRestanteBOB = computed(() => roundMoney(saldoRestanteUSD.value * tipoCambio.value));
const pagoUSDDirecto = computed(() => form.pagos.filter((pago: any) => pago.moneda === 'USD').reduce((total: number, pago: any) => total + Number(pago.monto || 0), 0));
const pagoBOBDirecto = computed(() => form.pagos.filter((pago: any) => pago.moneda === 'BOB').reduce((total: number, pago: any) => total + Number(pago.monto || 0), 0));
const pagoEquivalenteUSD = computed(() => roundMoney(pagoUSDDirecto.value + (tipoCambio.value > 0 ? pagoBOBDirecto.value / tipoCambio.value : 0)));
const pagoEquivalenteBOB = computed(() => roundMoney(pagoEquivalenteUSD.value * tipoCambio.value));
const pagoExcedeSaldo = computed(() => pagoEquivalenteUSD.value > saldoRestanteUSD.value);
const esContado = computed(() => form.tipo_pago === 'Contado');
const esCreditoDirecto = computed(() => !esContado.value && form.tipo_credito === 'credito_directo');
const esCreditoBancario = computed(() => !esContado.value && form.tipo_credito === 'credito_bancario');
const montoCuotaCredito = computed(() => Number(form.numero_cuotas || 0) > 0 ? roundMoney(saldoRestanteUSD.value / Number(form.numero_cuotas || 1)) : 0);

onMounted(async () => {
  await cargarReserva();
});

async function cargarReserva() {
  loading.value = true;
  try {
    const id = route.query.id;
    if (!id) {
      reserva.value = null;
      return;
    }
    const res: any = await $fetch(server.HOST + '/api/v1/ventas/' + id, { method: 'GET' });
    if (res?.tipo_venta !== 'Reserva') {
      reserva.value = null;
      return;
    }
    reserva.value = res;
    form.monto_pago = Number(res.saldo || 0);
    form.tipo_cambio = Number(res.tipo_cambio_usado || 0) || null;
    form.pagos.splice(0, form.pagos.length, { ...crearPago(), monto: Number(res.saldo || 0) });
    form.metodo_pago = res.metodo_pago || 'Efectivo';
    form.id_vehiculo = res.id_vehiculo || null;
    if (res.tipo_reserva === 'pedido' && !res.id_vehiculo) {
      await obtenerVehiculos();
    }
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar reserva', life: 3000 });
  } finally {
    loading.value = false;
  }
}

async function completarReserva() {
  if (!reserva.value) return;
  if (requiereVehiculoImportado.value) {
    toast.add({ severity: 'warn', summary: 'Registre o seleccione el vehiculo importado', life: 4000 });
    return;
  }
  if (tipoCambio.value <= 0) {
    toast.add({ severity: 'warn', summary: 'Ingrese tipo de cambio', life: 3000 });
    return;
  }
  if (esContado.value) {
    if (pagoEquivalenteUSD.value !== saldoRestanteUSD.value || form.pagos.some((pago: any) => !pago.moneda || !pago.metodo || Number(pago.monto || 0) <= 0)) {
      toast.add({ severity: 'warn', summary: 'Debe completar el saldo con un detalle de pago valido', detail: `$ ${formatPrecio(saldoRestanteUSD.value)}`, life: 4000 });
      return;
    }
  } else if (!validarCredito()) {
    return;
  }

  saving.value = true;
  try {
    const user = localStorage.getItem('user');
    const userId = user ? JSON.parse(user)?.id : null;
    if (!userId) {
      toast.add({ severity: 'warn', summary: 'No se encontro usuario en sesion', life: 3000 });
      return;
    }
    await $fetch(server.HOST + `/api/v1/ventas/${reserva.value.id}/completar-reserva`, {
      method: 'PATCH',
      body: {
        tipo_pago: form.tipo_pago.toLowerCase(),
        tipo_cambio: tipoCambio.value,
        pagos: esContado.value ? form.pagos.map((pago: any) => ({ moneda: pago.moneda, metodo: pago.metodo, monto: Number(pago.monto || 0) })) : [],
        monto_pago: saldoRestanteUSD.value,
        id_usuario_pago: Number(userId || 0),
        id_vehiculo: form.id_vehiculo,
        metodo_pago: form.metodo_pago,
        estado_entrega: form.estado_entrega,
        fecha_entrega: form.fecha_entrega,
        observacion: form.observacion,
        tipo_credito: form.tipo_credito,
        numero_cuotas: Number(form.numero_cuotas || 0),
        fecha_inicio_credito: form.fecha_inicio_credito,
        frecuencia_pago: form.frecuencia_pago,
        tiene_respaldo: form.tiene_respaldo,
        tipo_garantia: form.tipo_garantia,
        documento_garantia: form.documento_garantia,
        datos_garante: form.datos_garante,
        referencia_bancaria: form.referencia_bancaria,
        estado_desembolso: form.estado_desembolso
      }
    });
    toast.add({ severity: 'success', summary: 'Reserva completada', detail: 'El pago del vehiculo fue completado.', life: 3000 });
    await descargarPDFVenta(reserva.value.id);
    setTimeout(() => router.push(props.historialPath), 900);
  } catch (err: any) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al completar reserva', detail: err?.data || err?.message, life: 4000 });
  } finally {
    saving.value = false;
  }
}

function validarCredito() {
  if (!form.tipo_credito || Number(form.numero_cuotas || 0) <= 0 || !form.fecha_inicio_credito) {
    toast.add({ severity: 'warn', summary: 'Complete los datos del credito', life: 3000 });
    return false;
  }
  if (esCreditoDirecto.value && (!form.tiene_respaldo || !form.tipo_garantia.trim() || !form.datos_garante.trim())) {
    toast.add({ severity: 'warn', summary: 'Credito directo incompleto', detail: 'Debe registrar respaldo, tipo de garantia y datos del garante.', life: 4000 });
    return false;
  }
  return true;
}

function crearPago() {
  return { key: `${Date.now()}-${Math.random().toString(36).slice(2)}`, moneda: 'USD', metodo: 'Efectivo', monto: 0 };
}

function agregarPago() {
  form.pagos.push(crearPago());
}

function eliminarPago(index: number) {
  form.pagos.splice(index, 1);
  if (form.pagos.length === 0) agregarPago();
}

async function obtenerVehiculos() {
  try {
    const res = await $fetch(server.HOST + '/api/v1/vehiculos', { method: 'GET' });
    vehiculos.value = Array.isArray(res) ? res.map((vehiculo: any) => ({
      ...vehiculo,
      nombreCompleto: etiquetaVehiculo(vehiculo)
    })) : [];
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar vehiculos', life: 3000 });
  }
}

function etiquetaVehiculo(vehiculo: any) {
  return [vehiculo.marca, vehiculo.modelo, vehiculo.anio].filter(Boolean).join(' ') || vehiculo.nombre || 'Vehiculo';
}

function mostrarError(summary: string, err: any) {
  toast.add({ severity: 'error', summary, detail: err?.data || err?.message, life: 4000 });
}

function formatPrecio(precio: number) {
  return Number(precio || 0).toLocaleString('es-BO', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

function formatFecha(fecha: string) {
  if (!fecha) return 'N/A';
  return new Date(`${fecha}T00:00:00`).toLocaleDateString('es-BO');
}

function roundMoney(value: number) {
  return Math.round(Number(value || 0) * 100) / 100;
}
</script>
