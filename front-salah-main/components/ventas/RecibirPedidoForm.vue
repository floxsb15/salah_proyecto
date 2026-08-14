<template>
  <Toast />

  <div class="flex flex-col gap-4">
    <div class="flex items-center justify-between gap-3">
      <div>
        <h2 class="text-2xl font-bold">Recibir pedido</h2>
        <p class="text-sm text-gray-500">Vincula el vehiculo importado al inventario y marca el pedido como recibido.</p>
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
        <div class="flex items-start justify-between gap-3">
          <div>
            <h3 class="text-lg font-semibold text-gray-900">{{ vehiculoSolicitado }}</h3>
            <p class="text-sm text-gray-500">{{ pedido.cliente }} / {{ pedido.pais_origen }}</p>
          </div>
          <Tag :value="pedido.estado" :severity="estadoSeverity(pedido.estado)" />
        </div>

        <div class="mt-4 grid grid-cols-2 gap-3 text-sm">
          <div class="rounded-md bg-gray-50 p-3">
            <span class="text-gray-500">Precio estimado</span>
            <strong class="block text-gray-900">$ {{ formatPrecio(pedido.precio_estimado_usd) }}</strong>
          </div>
          <div class="rounded-md bg-gray-50 p-3">
            <span class="text-gray-500">Llegada estimada</span>
            <strong class="block text-gray-900">{{ formatFecha(pedido.fecha_llegada_estimada) }}</strong>
          </div>
          <div class="rounded-md bg-gray-50 p-3">
            <span class="text-gray-500">Adelanto</span>
            <strong class="block text-gray-900">$ {{ formatPrecio(totalPagadoUSD) }}</strong>
          </div>
          <div class="rounded-md bg-gray-50 p-3">
            <span class="text-gray-500">Saldo</span>
            <strong class="block text-gray-900">$ {{ formatPrecio(pedido.saldo_pendiente_usd) }}</strong>
          </div>
        </div>
      </section>

      <section class="rounded-lg border border-gray-200 bg-white p-4">
        <form class="grid grid-cols-1 gap-4 lg:grid-cols-2" @submit.prevent="recibirPedido">
          <section class="grid grid-cols-1 gap-3 rounded-md border border-yellow-200 bg-yellow-50 p-3 lg:col-span-2 md:grid-cols-[1fr_auto]">
            <div class="flex flex-col gap-1">
              <label for="vehiculo">Vehiculo en inventario</label>
              <Select id="vehiculo" v-model="form.id_vehiculo" :options="vehiculos" option-label="nombreCompleto" option-value="id" placeholder="Seleccione el vehiculo importado" filter fluid size="small" />
            </div>
            <div class="flex items-end">
              <Button label="Crear vehiculo" icon="pi pi-car" severity="secondary" type="button" size="small" @click="vehiculoVisible = true" />
            </div>
          </section>

          <div class="flex flex-col gap-1 lg:col-span-2">
            <label for="observacion">Observacion</label>
            <Textarea id="observacion" v-model="form.observacion" rows="3" auto-resize fluid />
          </div>

          <div class="flex justify-end gap-2 border-t border-gray-100 pt-4 lg:col-span-2">
            <Button label="Cancelar" severity="secondary" type="button" @click="router.push(historialPath)" />
            <Button label="Confirmar recepcion" icon="pi pi-check" severity="success" type="submit" :loading="saving" />
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
import Select from 'primevue/select';
import Skeleton from 'primevue/skeleton';
import Tag from 'primevue/tag';
import Textarea from 'primevue/textarea';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';
import modalAgregarVehiculo from '~/components/admin/vehiculos/modalAgregarProducto.vue';

const props = withDefaults(defineProps<{ historialPath?: string }>(), {
  historialPath: '/ventas/pedidos'
});

const route = useRoute();
const router = useRouter();
const toast = useToast();
const pedido = ref<any>(null);
const vehiculos = ref<any[]>([]);
const loading = ref(true);
const saving = ref(false);
const vehiculoVisible = ref(false);
const form = reactive({
  id_vehiculo: null as number | null,
  observacion: ''
});

const vehiculoSolicitado = computed(() => [pedido.value?.marca, pedido.value?.modelo, pedido.value?.anio].filter(Boolean).join(' ') || 'Vehiculo solicitado');
const totalPagadoUSD = computed(() => {
  const tc = Number(pedido.value?.tipo_cambio_usado || 0);
  return roundMoney(Number(pedido.value?.adelanto_pagado_usd || 0) + (tc > 0 ? Number(pedido.value?.adelanto_pagado_bob || 0) / tc : 0));
});

onMounted(async () => {
  await cargarDatos();
});

async function cargarDatos() {
  loading.value = true;
  try {
    await Promise.all([cargarPedido(), obtenerVehiculos()]);
  } finally {
    loading.value = false;
  }
}

async function cargarPedido() {
  const id = route.query.id;
  if (!id) {
    pedido.value = null;
    return;
  }
  try {
    pedido.value = await $fetch(server.HOST + '/api/v1/pedidos/' + id, { method: 'GET' });
    form.id_vehiculo = pedido.value?.id_vehiculo || null;
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar pedido', life: 3000 });
  }
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

async function recibirPedido() {
  if (!pedido.value || !form.id_vehiculo) {
    toast.add({ severity: 'warn', summary: 'Seleccione el vehiculo recibido', life: 3000 });
    return;
  }
  saving.value = true;
  try {
    await $fetch(server.HOST + `/api/v1/pedidos/${pedido.value.id}/recibir`, {
      method: 'PATCH',
      body: {
        id_vehiculo: form.id_vehiculo,
        observacion: form.observacion
      }
    });
    toast.add({ severity: 'success', summary: 'Pedido recibido', life: 3000 });
    setTimeout(() => router.push(props.historialPath), 700);
  } catch (err: any) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al recibir pedido', detail: err?.data || err?.message, life: 4000 });
  } finally {
    saving.value = false;
  }
}

function etiquetaVehiculo(vehiculo: any) {
  return [vehiculo.marca, vehiculo.modelo, vehiculo.anio].filter(Boolean).join(' ') || vehiculo.nombre || 'Vehiculo';
}

function estadoSeverity(estado: string) {
  if (estado === 'Completado') return 'success';
  if (estado === 'Recibido') return 'warning';
  if (String(estado || '').includes('trans')) return 'info';
  return 'secondary';
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

function mostrarError(summary: string, err: any) {
  toast.add({ severity: 'error', summary, detail: err?.data || err?.message, life: 4000 });
}
</script>
