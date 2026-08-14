<template>
  <Toast />

  <div class="flex flex-col gap-4">
    <div class="flex items-center justify-between gap-3">
      <div>
        <h2 class="text-2xl font-bold">Completar pedido</h2>
        <p class="text-sm text-gray-500">Registra el pago final del pedido recibido.</p>
      </div>
      <Button label="Volver" icon="pi pi-arrow-left" size="small" severity="secondary" @click="router.push(historialPath)" />
    </div>

    <div v-if="loading" class="grid grid-cols-1 gap-4 lg:grid-cols-2">
      <Skeleton height="18rem" />
      <Skeleton height="18rem" />
    </div>

    <div v-else-if="!pedido" class="flex flex-col items-center justify-center gap-2 py-16 text-gray-500">
      <i class="pi pi-exclamation-triangle text-4xl"></i>
      <p>No se encontro el pedido seleccionado.</p>
    </div>

    <div v-else class="grid grid-cols-1 gap-4 xl:grid-cols-[24rem_1fr]">
      <section class="rounded-lg border border-gray-200 bg-white p-4">
        <h3 class="text-lg font-semibold text-gray-900">{{ pedido.vehiculo || vehiculoSolicitado }}</h3>
        <p class="text-sm text-gray-500">{{ pedido.cliente }} / {{ pedido.pais_origen }}</p>

        <div class="mt-4 grid grid-cols-2 gap-3 text-sm">
          <div class="rounded-md bg-gray-50 p-3">
            <span class="text-gray-500">Precio estimado</span>
            <strong class="block text-gray-900">$ {{ formatPrecio(pedido.precio_estimado_usd) }}</strong>
          </div>
          <div class="rounded-md bg-gray-50 p-3">
            <span class="text-gray-500">Adelanto pagado</span>
            <strong class="block text-gray-900">$ {{ formatPrecio(totalPagadoUSD) }}</strong>
          </div>
          <div class="rounded-md bg-gray-50 p-3">
            <span class="text-gray-500">Saldo pendiente</span>
            <strong class="block text-gray-900">$ {{ formatPrecio(pedido.saldo_pendiente_usd) }}</strong>
          </div>
          <div class="rounded-md bg-gray-50 p-3">
            <span class="text-gray-500">Estado</span>
            <strong class="block text-gray-900">{{ pedido.estado }}</strong>
          </div>
        </div>
      </section>

      <section class="rounded-lg border border-gray-200 bg-white p-4">
        <form class="grid grid-cols-1 gap-4 lg:grid-cols-2" @submit.prevent="completarPedido">
          <div class="flex flex-col gap-1">
            <label for="tipo_cambio">Tipo de cambio del dia</label>
            <InputNumber id="tipo_cambio" v-model="form.tipo_cambio" :min="0" :minFractionDigits="2" :maxFractionDigits="4" suffix=" Bs/USD" fluid size="small" />
          </div>
          <div class="flex flex-col gap-1">
            <label>Pago final requerido</label>
            <div class="rounded-md border border-gray-200 bg-gray-50 px-3 py-2 text-sm">
              <strong>$ {{ formatPrecio(pedido.saldo_pendiente_usd) }}</strong>
              <span class="ml-2 text-gray-500">Bs {{ formatPrecio(saldoFinalBOB) }}</span>
            </div>
          </div>

          <section class="flex flex-col gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2">
            <div class="flex items-center justify-between gap-3">
              <div>
                <h3 class="text-sm font-semibold text-gray-900">Detalle de pago final</h3>
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

          <div class="flex flex-col gap-1 lg:col-span-2">
            <label for="observacion">Observacion</label>
            <Textarea id="observacion" v-model="form.observacion" rows="3" auto-resize fluid />
          </div>

          <div class="flex justify-end gap-2 border-t border-gray-100 pt-4 lg:col-span-2">
            <Button label="Cancelar" severity="secondary" type="button" @click="router.push(historialPath)" />
            <Button label="Completar pedido" icon="pi pi-check" severity="success" type="submit" :loading="saving" />
          </div>
        </form>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { server } from '~/server/server';
import Button from 'primevue/button';
import InputNumber from 'primevue/inputnumber';
import Select from 'primevue/select';
import Skeleton from 'primevue/skeleton';
import Textarea from 'primevue/textarea';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';

const props = withDefaults(defineProps<{ historialPath?: string }>(), {
  historialPath: '/ventas/pedidos'
});

const route = useRoute();
const router = useRouter();
const toast = useToast();
const pedido = ref<any>(null);
const loading = ref(true);
const saving = ref(false);
const monedasPago = ref(['USD', 'BOB']);
const metodosPago = ref(['Efectivo', 'QR', 'Transferencia', 'Tarjeta']);
const form = reactive({
  tipo_cambio: null as number | null,
  pagos: [crearPago()],
  observacion: ''
});

const vehiculoSolicitado = computed(() => [pedido.value?.marca, pedido.value?.modelo, pedido.value?.anio].filter(Boolean).join(' ') || 'Vehiculo solicitado');
const tipoCambio = computed(() => Number(form.tipo_cambio || 0));
const totalPagadoUSD = computed(() => {
  const tc = Number(pedido.value?.tipo_cambio_usado || 0);
  return roundMoney(Number(pedido.value?.adelanto_pagado_usd || 0) + (tc > 0 ? Number(pedido.value?.adelanto_pagado_bob || 0) / tc : 0));
});
const saldoFinalUSD = computed(() => Number(pedido.value?.saldo_pendiente_usd || 0));
const saldoFinalBOB = computed(() => roundMoney(saldoFinalUSD.value * tipoCambio.value));
const pagoUSDDirecto = computed(() => form.pagos.filter((pago: any) => pago.moneda === 'USD').reduce((total: number, pago: any) => total + Number(pago.monto || 0), 0));
const pagoBOBDirecto = computed(() => form.pagos.filter((pago: any) => pago.moneda === 'BOB').reduce((total: number, pago: any) => total + Number(pago.monto || 0), 0));
const pagoEquivalenteUSD = computed(() => roundMoney(pagoUSDDirecto.value + (tipoCambio.value > 0 ? pagoBOBDirecto.value / tipoCambio.value : 0)));
const pagoEquivalenteBOB = computed(() => roundMoney(pagoEquivalenteUSD.value * tipoCambio.value));
const pagoExcedeSaldo = computed(() => pagoEquivalenteUSD.value > saldoFinalUSD.value);

onMounted(async () => {
  await cargarPedido();
});

async function cargarPedido() {
  loading.value = true;
  try {
    const id = route.query.id;
    if (!id) {
      pedido.value = null;
      return;
    }
    pedido.value = await $fetch(server.HOST + '/api/v1/pedidos/' + id, { method: 'GET' });
    form.tipo_cambio = Number(pedido.value?.tipo_cambio_usado || 0) || null;
    form.pagos.splice(0, form.pagos.length, {
      ...crearPago(),
      monto: Number(pedido.value?.saldo_pendiente_usd || 0)
    });
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar pedido', life: 3000 });
  } finally {
    loading.value = false;
  }
}

async function completarPedido() {
  if (!pedido.value) return;
  if (tipoCambio.value <= 0) {
    toast.add({ severity: 'warn', summary: 'Ingrese tipo de cambio', life: 3000 });
    return;
  }
  if (pagoEquivalenteUSD.value !== saldoFinalUSD.value) {
    toast.add({ severity: 'warn', summary: 'Debe pagar el saldo completo', detail: `$ ${formatPrecio(saldoFinalUSD.value)}`, life: 4000 });
    return;
  }
  if (form.pagos.some((pago: any) => !pago.moneda || !pago.metodo || Number(pago.monto || 0) <= 0)) {
    toast.add({ severity: 'warn', summary: 'Complete moneda, tipo de pago y monto en cada fila', life: 4000 });
    return;
  }
  saving.value = true;
  try {
    await $fetch(server.HOST + `/api/v1/pedidos/${pedido.value.id}/completar`, {
      method: 'PATCH',
      body: {
        tipo_cambio: tipoCambio.value,
        pagos: form.pagos.map((pago: any) => ({
          moneda: pago.moneda,
          metodo: pago.metodo,
          monto: Number(pago.monto || 0)
        })),
        observacion: form.observacion
      }
    });
    toast.add({ severity: 'success', summary: 'Pedido completado', life: 3000 });
    setTimeout(() => router.push(props.historialPath), 700);
  } catch (err: any) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al completar pedido', detail: err?.data || err?.message, life: 4000 });
  } finally {
    saving.value = false;
  }
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

function formatPrecio(precio: number) {
  return Number(precio || 0).toLocaleString('es-BO', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

function roundMoney(value: number) {
  return Math.round(Number(value || 0) * 100) / 100;
}
</script>
