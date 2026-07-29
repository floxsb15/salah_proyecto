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
          <div class="flex flex-col gap-1">
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

          <div class="flex flex-col gap-1">
            <label for="metodo_pago">Metodo de pago</label>
            <Select id="metodo_pago" v-model="form.metodo_pago" :options="metodosPago" fluid size="small" />
          </div>

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
import { onMounted, reactive, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { server } from '~/server/server';
import Button from 'primevue/button';
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
const metodosPago = ref(['QR', 'Transferencia', 'Efectivo']);
const estadosEntrega = ref(['Pendiente', 'Entregado']);
const form = reactive({
  monto_pago: 0,
  metodo_pago: 'Efectivo',
  estado_entrega: 'Pendiente',
  fecha_entrega: '',
  observacion: ''
});

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
    form.metodo_pago = res.metodo_pago || 'Efectivo';
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar reserva', life: 3000 });
  } finally {
    loading.value = false;
  }
}

async function completarReserva() {
  if (!reserva.value) return;
  const saldo = Number(reserva.value.saldo || 0);
  if (Number(form.monto_pago || 0) !== saldo) {
    toast.add({ severity: 'warn', summary: 'Debe pagar el saldo completo', detail: `Saldo pendiente: $ ${formatPrecio(saldo)}`, life: 4000 });
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
        monto_pago: Number(form.monto_pago || 0),
        id_usuario_pago: Number(userId || 0),
        metodo_pago: form.metodo_pago,
        estado_entrega: form.estado_entrega,
        fecha_entrega: form.fecha_entrega,
        observacion: form.observacion
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

function formatPrecio(precio: number) {
  return Number(precio || 0).toLocaleString('es-BO', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

function formatFecha(fecha: string) {
  if (!fecha) return 'N/A';
  return new Date(`${fecha}T00:00:00`).toLocaleDateString('es-BO');
}
</script>
