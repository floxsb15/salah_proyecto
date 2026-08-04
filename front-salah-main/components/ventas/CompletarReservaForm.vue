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
        <h3 class="text-lg font-semibold text-gray-900">{{ reserva.vehiculo }}</h3>
        <p class="text-sm text-gray-500">{{ reserva.cliente }} · CI/NIT: {{ reserva.ci_cliente || 'N/A' }}</p>

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
          <div class="flex flex-col gap-1 lg:col-span-2">
            <label for="tipo_pago_saldo">Tipo de pago del saldo</label>
            <Select id="tipo_pago_saldo" v-model="form.tipo_pago_saldo" :options="tiposPagoSaldo" option-label="label" option-value="value" fluid size="small" />
          </div>

          <template v-if="esPagoContado">
            <div class="flex flex-col gap-1">
              <label for="metodo_pago">Metodo de pago</label>
              <Select id="metodo_pago" v-model="form.metodo_pago" :options="metodosPago" fluid size="small" />
            </div>

            <div v-if="!esPagoMixto" class="flex flex-col gap-1">
              <label for="monto_pago">Pago restante</label>
              <InputNumber
                id="monto_pago"
                v-model="form.monto_pago"
                mode="currency"
                currency="USD"
                locale="es-BO"
                :min="0"
                :max="Number(reserva.saldo || 0)"
                fluid
                size="small"
              />
            </div>

            <section v-else class="flex flex-col gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2">
              <div class="flex items-center justify-between gap-3">
                <h3 class="text-sm font-semibold text-gray-900">Detalle de pago</h3>
                <Button label="Agregar pago" icon="pi pi-plus" size="small" severity="secondary" type="button" @click="agregarPago" />
              </div>

              <div v-for="(pago, index) in form.pagos" :key="pago.key" class="grid grid-cols-1 gap-2 md:grid-cols-[9rem_13rem_1fr_2.5rem]">
                <Select v-model="pago.moneda" :options="monedasPago" placeholder="Moneda" fluid size="small" />
                <Select v-model="pago.metodo" :options="metodosPagoLinea" placeholder="Metodo" fluid size="small" />
                <InputNumber v-model="pago.monto" mode="decimal" :min="0" :minFractionDigits="2" :maxFractionDigits="2" placeholder="Monto" fluid size="small" />
                <Button icon="pi pi-trash" severity="danger" text rounded type="button" aria-label="Eliminar pago" @click="eliminarPago(index)" />
              </div>

              <p v-if="pagoMixtoExcedeSaldo" class="text-sm font-semibold text-red-600">
                El pagado supera el saldo pendiente.
              </p>
            </section>

            <div class="grid grid-cols-1 gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2 md:grid-cols-2">
              <div>
                <span class="text-xs font-semibold uppercase text-gray-500">Pagado equivalente</span>
                <strong class="block text-gray-900">$ {{ formatPrecio(pagoEquivalenteUSD) }}</strong>
                <span class="block text-sm font-semibold text-gray-700">Bs {{ formatPrecio(pagoEquivalenteBOB) }}</span>
              </div>
              <div>
                <span class="text-xs font-semibold uppercase text-gray-500">Saldo pendiente</span>
                <strong class="block text-gray-900">$ {{ formatPrecio(saldoRestanteUSD) }}</strong>
                <span class="block text-sm font-semibold text-gray-700">Bs {{ formatPrecio(saldoRestanteBOB) }}</span>
              </div>
            </div>
          </template>

          <template v-if="esCreditoDirecto">
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
                <FileUpload
                  id="documento_garantia"
                  name="documento_garantia"
                  mode="basic"
                  choose-label="Subir garantia"
                  accept="image/*,application/pdf"
                  severity="secondary"
                  class="p-button-outlined"
                  custom-upload
                  auto
                  @select="onDocumentoGarantiaSelect"
                  @clear="limpiarDocumentoGarantia"
                  @remove="limpiarDocumentoGarantia"
                />
                <small class="text-xs text-gray-500">{{ documentoGarantiaNombre || 'Formatos: imagen o PDF, maximo 10MB.' }}</small>
              </div>
            </div>

            <div class="flex flex-col gap-1 lg:col-span-2">
              <label for="datos_garante">Datos del garante/respaldo</label>
              <Textarea id="datos_garante" v-model="form.datos_garante" rows="3" auto-resize fluid />
            </div>

            <div class="grid grid-cols-1 gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2 md:grid-cols-2">
              <div>
                <span class="text-xs font-semibold uppercase text-gray-500">Monto financiado</span>
                <strong class="block text-gray-900">$ {{ formatPrecio(saldoPendiente) }}</strong>
              </div>
              <div>
                <span class="text-xs font-semibold uppercase text-gray-500">Monto por cuota</span>
                <strong class="block text-gray-900">$ {{ formatPrecio(montoCuota) }}</strong>
              </div>
            </div>
          </template>

          <template v-if="esCreditoBancario">
            <div class="flex flex-col gap-1">
              <label for="banco_entidad_financiera">Banco/entidad financiera</label>
              <InputText id="banco_entidad_financiera" v-model="form.banco_entidad_financiera" placeholder="Nombre del banco" fluid size="small" />
            </div>

            <div class="flex flex-col gap-1">
              <label for="estado_tramite_bancario">Estado del tramite</label>
              <Select id="estado_tramite_bancario" v-model="form.estado_tramite_bancario" :options="estadosTramiteBancario" fluid size="small" />
            </div>

            <div class="flex flex-col gap-1">
              <label for="fecha_estimada_desembolso">Fecha estimada de desembolso</label>
              <InputText id="fecha_estimada_desembolso" v-model="form.fecha_estimada_desembolso" type="date" fluid size="small" />
            </div>

            <div class="rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2">
              <span class="text-xs font-semibold uppercase text-gray-500">Monto a financiar por banco</span>
              <strong class="block text-gray-900">$ {{ formatPrecio(saldoPendiente) }}</strong>
            </div>
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
            <Button label="Completar pago" icon="pi pi-check" type="button" :loading="saving" @click="completarReserva" />
          </div>
        </form>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { server } from '~/server/server';
import Button from 'primevue/button';
import Checkbox from 'primevue/checkbox';
import FileUpload from 'primevue/fileupload';
import InputNumber from 'primevue/inputnumber';
import InputText from 'primevue/inputtext';
import Select from 'primevue/select';
import Skeleton from 'primevue/skeleton';
import Textarea from 'primevue/textarea';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';

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
const loading = ref(true);
const saving = ref(false);
const documentoGarantia = ref<File | null>(null);
const documentoGarantiaNombre = ref('');
const tiposPagoSaldo = ref([
  { label: 'Contado', value: 'contado' },
  { label: 'Credito directo', value: 'credito_directo' },
  { label: 'Credito bancario', value: 'credito_bancario' }
]);
const metodosPago = ref(['QR', 'Transferencia', 'Tarjeta', 'Efectivo', 'Mixto']);
const metodosPagoLinea = ref(['QR', 'Transferencia', 'Tarjeta', 'Efectivo']);
const monedasPago = ref(['USD', 'BOB']);
const estadosEntrega = ref(['Pendiente', 'Entregado']);
const frecuenciasPago = ref(['mensual', 'quincenal', 'semanal']);
const estadosTramiteBancario = ref(['Pendiente', 'En evaluacion', 'Aprobado', 'Rechazado']);
const form = reactive({
  tipo_pago_saldo: 'contado',
  monto_pago: 0,
  metodo_pago: 'Efectivo',
  pagos: [
    crearPago()
  ],
  estado_entrega: 'Pendiente',
  fecha_entrega: '',
  numero_cuotas: 12,
  fecha_inicio_credito: new Date().toISOString().slice(0, 10),
  frecuencia_pago: 'mensual',
  tiene_respaldo: false,
  tipo_garantia: '',
  datos_garante: '',
  banco_entidad_financiera: '',
  estado_tramite_bancario: 'Pendiente',
  fecha_estimada_desembolso: '',
  observacion: ''
});

const saldoPendiente = computed(() => Number(reserva.value?.saldo || 0));
const tipoCambio = computed(() => Number(reserva.value?.tipo_cambio_usado || 0));
const esPagoContado = computed(() => form.tipo_pago_saldo === 'contado');
const esCreditoDirecto = computed(() => form.tipo_pago_saldo === 'credito_directo');
const esCreditoBancario = computed(() => form.tipo_pago_saldo === 'credito_bancario');
const esPagoMixto = computed(() => form.metodo_pago === 'Mixto');
const pagoUSDDirecto = computed(() => form.pagos
  .filter((pago: any) => pago.moneda === 'USD')
  .reduce((total: number, pago: any) => total + Number(pago.monto || 0), 0));
const pagoBOBDirecto = computed(() => form.pagos
  .filter((pago: any) => pago.moneda === 'BOB')
  .reduce((total: number, pago: any) => total + Number(pago.monto || 0), 0));
const pagoEquivalenteUSD = computed(() => {
  if (!esPagoContado.value) return 0;
  if (!esPagoMixto.value) return Number(form.monto_pago || 0);
  return pagoUSDDirecto.value + (tipoCambio.value > 0 ? pagoBOBDirecto.value / tipoCambio.value : 0);
});
const pagoEquivalenteBOB = computed(() => {
  if (!esPagoContado.value) return 0;
  if (!esPagoMixto.value) return Number(form.monto_pago || 0) * tipoCambio.value;
  return pagoBOBDirecto.value + (pagoUSDDirecto.value * tipoCambio.value);
});
const saldoRestanteUSD = computed(() => Math.max(saldoPendiente.value - pagoEquivalenteUSD.value, 0));
const saldoRestanteBOB = computed(() => Math.max((saldoPendiente.value * tipoCambio.value) - pagoEquivalenteBOB.value, 0));
const pagoMixtoExcedeSaldo = computed(() => esPagoMixto.value && pagoEquivalenteUSD.value > saldoPendiente.value);
const montoCuota = computed(() => Number(form.numero_cuotas || 0) > 0 ? saldoPendiente.value / Number(form.numero_cuotas || 1) : 0);
const creditosPath = computed(() => props.historialPath.startsWith('/admin') ? '/admin/creditos' : '/ventas/creditos');

onMounted(async () => {
  await cargarReserva();
});

watch(() => form.metodo_pago, () => {
  if (esPagoMixto.value && form.pagos.length === 1 && Number(form.pagos[0].monto || 0) === 0) {
    form.pagos[0].monto = saldoPendiente.value;
  }
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
    form.metodo_pago = res.metodo_pago || 'Efectivo';
    form.pagos = [{ ...crearPago(), monto: Number(res.saldo || 0) }];
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar reserva', life: 3000 });
  } finally {
    loading.value = false;
  }
}

async function completarReserva() {
  if (!reserva.value) return;
  const saldo = saldoPendiente.value;
  if (!validarFormulario(saldo)) return;

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
      body: construirPayload(Number(userId || 0))
    });
    toast.add({ severity: 'success', summary: 'Reserva completada', detail: esPagoContado.value ? 'El pago del vehiculo fue completado.' : 'El saldo fue registrado como credito.', life: 3000 });
    await descargarPDFVenta(reserva.value.id);
    setTimeout(() => router.push(esPagoContado.value ? props.historialPath : creditosPath.value), 900);
  } catch (err: any) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al completar reserva', detail: err?.data || err?.message, life: 4000 });
  } finally {
    saving.value = false;
  }
}

function validarFormulario(saldo: number) {
  if (esPagoContado.value) {
    if (roundMoney(pagoEquivalenteUSD.value) !== roundMoney(saldo)) {
      toast.add({ severity: 'warn', summary: 'Debe pagar el saldo completo', detail: `Saldo pendiente: $ ${formatPrecio(saldo)}`, life: 4000 });
      return false;
    }
    if (esPagoMixto.value) {
      const lineaInvalida = form.pagos.find((pago: any) => !pago.moneda || !pago.metodo || Number(pago.monto || 0) <= 0);
      if (lineaInvalida) {
        toast.add({ severity: 'warn', summary: 'Complete moneda, metodo y monto mayor a cero en cada linea', life: 4000 });
        return false;
      }
    }
    return true;
  }

  if (esCreditoDirecto.value) {
    if (Number(form.numero_cuotas || 0) <= 0 || !form.fecha_inicio_credito) {
      toast.add({ severity: 'warn', summary: 'Complete plazo y fecha de inicio del credito', life: 3000 });
      return false;
    }
    if (!form.tiene_respaldo || !form.tipo_garantia.trim() || !form.datos_garante.trim()) {
      toast.add({ severity: 'warn', summary: 'Credito directo bloqueado', detail: 'Debe registrar respaldo, tipo de garantia y datos del garante.', life: 4000 });
      return false;
    }
    if (documentoGarantia.value && documentoGarantia.value.size > 10 * 1024 * 1024) {
      toast.add({ severity: 'warn', summary: 'Documento muy grande', detail: 'El documento de garantia no debe superar los 10MB.', life: 4000 });
      return false;
    }
  }

  if (esCreditoBancario.value && !form.banco_entidad_financiera.trim()) {
    toast.add({ severity: 'warn', summary: 'Ingrese banco o entidad financiera', life: 3000 });
    return false;
  }

  return true;
}

function construirPayload(userId: number) {
  const payload: Record<string, any> = {
    tipo_pago_saldo: form.tipo_pago_saldo,
    monto_pago: esPagoContado.value && !esPagoMixto.value ? Number(form.monto_pago || 0) : Number(pagoEquivalenteUSD.value || 0),
    id_usuario_pago: userId,
    metodo_pago: esPagoContado.value ? form.metodo_pago : 'Efectivo',
    pagos: esPagoContado.value && esPagoMixto.value ? detallePagoPayload() : [],
    estado_entrega: form.estado_entrega,
    fecha_entrega: form.fecha_entrega,
    numero_cuotas: esCreditoDirecto.value ? Number(form.numero_cuotas || 0) : 0,
    fecha_inicio_credito: esCreditoDirecto.value ? form.fecha_inicio_credito : '',
    frecuencia_pago: esCreditoDirecto.value ? form.frecuencia_pago : '',
    tiene_respaldo: esCreditoDirecto.value ? form.tiene_respaldo : false,
    tipo_garantia: esCreditoDirecto.value ? form.tipo_garantia : '',
    datos_garante: esCreditoDirecto.value ? form.datos_garante : '',
    banco_entidad_financiera: esCreditoBancario.value ? form.banco_entidad_financiera : '',
    estado_tramite_bancario: esCreditoBancario.value ? form.estado_tramite_bancario : '',
    fecha_estimada_desembolso: esCreditoBancario.value ? form.fecha_estimada_desembolso : '',
    observacion: form.observacion
  };

  if (!documentoGarantia.value) {
    return payload;
  }

  const formData = new FormData();
  Object.entries(payload).forEach(([key, value]) => {
    formData.append(key, key === 'pagos' ? JSON.stringify(value) : String(value));
  });
  formData.append('documento_garantia', documentoGarantia.value);
  return formData;
}

function crearPago() {
  return {
    key: `${Date.now()}-${Math.random().toString(36).slice(2)}`,
    moneda: 'USD',
    metodo: 'Efectivo',
    monto: 0
  };
}

function agregarPago() {
  form.pagos.push(crearPago());
}

function eliminarPago(index: number) {
  form.pagos.splice(index, 1);
  if (form.pagos.length === 0) {
    agregarPago();
  }
}

function detallePagoPayload() {
  return form.pagos.map((pago: any) => ({
    moneda: pago.moneda,
    metodo: pago.metodo,
    monto: Number(pago.monto || 0)
  }));
}

function roundMoney(value: number) {
  return Math.round(Number(value || 0) * 100) / 100;
}

function onDocumentoGarantiaSelect(event: any) {
  const file = event.files?.[0] || null;
  if (!file) {
    limpiarDocumentoGarantia();
    return;
  }
  documentoGarantia.value = file;
  documentoGarantiaNombre.value = file.name;
}

function limpiarDocumentoGarantia() {
  documentoGarantia.value = null;
  documentoGarantiaNombre.value = '';
}

function formatPrecio(precio: number) {
  return Number(precio || 0).toLocaleString('es-BO', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

function formatFecha(fecha: string) {
  if (!fecha) return 'N/A';
  return new Date(`${fecha}T00:00:00`).toLocaleDateString('es-BO');
}
</script>
