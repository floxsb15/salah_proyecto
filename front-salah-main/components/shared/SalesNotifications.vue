<template>
  <div class="relative" ref="containerRef">
    <button
      type="button"
      class="relative inline-flex h-10 w-10 items-center justify-center rounded-md border border-neutral-700 bg-neutral-900 transition-colors hover:bg-red-900/60"
      aria-label="Notificaciones de ventas"
      @click="toggleOpen"
    >
      <i class="pi pi-bell text-yellow-400"></i>
      <span
        v-if="ventasPendientes.length"
        class="absolute -right-1 -top-1 inline-flex min-h-5 min-w-5 items-center justify-center rounded-full bg-red-500 px-1 text-xs font-bold text-white"
      >
        {{ badgeCount }}
      </span>
    </button>

    <div
      v-if="open"
      class="absolute right-0 top-12 z-[60] w-96 max-w-[calc(100vw-2rem)] overflow-hidden rounded-lg border border-gray-200 bg-white text-gray-900 shadow-xl"
    >
      <div class="flex items-center justify-between border-b border-gray-100 px-4 py-3">
        <div>
          <h3 class="text-sm font-semibold">Ventas</h3>
          <p class="text-xs text-gray-500">{{ subtitle }}</p>
        </div>
        <div class="flex items-center gap-1">
          <button
            v-if="ventasPendientes.length"
            type="button"
            class="inline-flex h-8 items-center gap-2 rounded-md px-2 text-xs font-medium text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900"
            @click="marcarTodasComoLeidas"
          >
            <i class="pi pi-check"></i>
            <span>Todo leido</span>
          </button>
          <button
            type="button"
            class="inline-flex h-8 w-8 items-center justify-center rounded-md text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-900"
            aria-label="Actualizar ventas"
            @click="obtenerVentas"
          >
            <i class="pi pi-refresh text-xs" :class="{ 'pi-spin': loading }"></i>
          </button>
        </div>
      </div>

      <div v-if="loading" class="flex flex-col gap-2 p-4">
        <Skeleton v-for="item in 3" :key="item" height="3rem" />
      </div>

      <div v-else-if="ventasRecientesPendientes.length === 0" class="px-4 py-8 text-center text-sm text-gray-500">
        No hay ventas pendientes.
      </div>

      <div v-else class="max-h-96 overflow-y-auto">
        <div
          v-for="venta in ventasRecientesPendientes"
          :key="venta.id"
          class="border-b border-gray-100 px-4 py-3 last:border-b-0"
        >
          <div class="flex items-start justify-between gap-3">
            <div class="min-w-0">
              <p class="truncate text-sm font-semibold text-gray-900">{{ venta.vehiculo || 'Vehiculo' }}</p>
              <p class="truncate text-xs text-gray-500">{{ venta.cliente || 'Cliente no registrado' }}</p>
              <p class="mt-1 text-xs text-gray-500">
                {{ formatFecha(venta.fecha) }} · {{ venta.vendedor || 'Sin vendedor' }}
              </p>
            </div>
            <div class="shrink-0 text-right">
              <p class="text-sm font-bold text-gray-900">$ {{ formatPrecio(venta.precio_total) }}</p>
              <Tag :value="venta.proforma_vencida ? 'Proforma vencida' : venta.estado_venta" :severity="severityVenta(venta)" />
              <button
                type="button"
                class="mt-2 inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900"
                @click="marcarComoLeida(venta.id)"
              >
                <i class="pi pi-check"></i>
                <span>Leido</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="border-t border-gray-100 px-4 py-3">
        <NuxtLink
          :to="historialPath"
          class="inline-flex w-full items-center justify-center gap-2 rounded-md bg-neutral-950 px-3 py-2 text-sm font-medium text-white transition-colors hover:bg-red-900"
          @click="open = false"
        >
          <i class="pi pi-history text-yellow-400"></i>
          <span>{{ actionLabel }}</span>
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { server } from '~/server/server';
import Skeleton from 'primevue/skeleton';
import Tag from 'primevue/tag';

const props = defineProps<{
  userData: {
    id?: number;
    nombre?: string;
    rol?: string;
  };
}>();

const open = ref(false);
const loading = ref(false);
const ventas = ref<any[]>([]);
const leidas = ref<number[]>([]);
const containerRef = ref<HTMLElement | null>(null);

const esVendedor = computed(() => props.userData?.rol === 'vendedor');
const storageKey = computed(() => `ventas_leidas_${props.userData?.rol || 'usuario'}_${props.userData?.id || 'general'}`);
const ventasPendientes = computed(() => ventas.value.filter(venta => !leidas.value.includes(Number(venta.id))));
const ventasRecientesPendientes = computed(() => ventasPendientes.value.slice(0, 6));
const badgeCount = computed(() => ventasPendientes.value.length > 99 ? '99+' : ventasPendientes.value.length);
const subtitle = computed(() => esVendedor.value ? 'Tus ventas pendientes' : 'Ventas pendientes');
const historialPath = computed(() => {
  if (props.userData?.rol === 'admin') {
    return '/admin/historial-ventas';
  }

  return esVendedor.value ? '/ventas/historial-ventas' : '/ventas/historial-general';
});
const actionLabel = computed(() => 'Ver historial');

onMounted(async () => {
  document.addEventListener('click', onDocumentClick);
  cargarLeidas();
  await obtenerVentas();
});

watch(
  () => [props.userData?.id, props.userData?.rol],
  async () => {
    cargarLeidas();
    await obtenerVentas();
  }
);

onBeforeUnmount(() => {
  document.removeEventListener('click', onDocumentClick);
});

async function toggleOpen() {
  open.value = !open.value;
  if (open.value) {
    await obtenerVentas();
  }
}

async function obtenerVentas() {
  loading.value = true;
  cargarLeidas();

  try {
    const query = esVendedor.value && props.userData?.id ? { id_usuario: props.userData.id } : undefined;
    const res = await $fetch(server.HOST + '/api/v1/ventas', {
      method: 'GET',
      query
    });

    ventas.value = Array.isArray(res) ? res : [];
  } catch (err) {
    console.error('Error al cargar notificaciones de ventas:', err);
    ventas.value = [];
  } finally {
    loading.value = false;
  }
}

function cargarLeidas() {
  try {
    const stored = localStorage.getItem(storageKey.value);
    const parsed = stored ? JSON.parse(stored) : [];
    leidas.value = Array.isArray(parsed) ? parsed.map((id: any) => Number(id)).filter(Boolean) : [];
  } catch (err) {
    console.error('Error al cargar ventas leidas:', err);
    leidas.value = [];
  }
}

function guardarLeidas() {
  localStorage.setItem(storageKey.value, JSON.stringify(leidas.value));
}

function marcarComoLeida(id: number) {
  const ventaId = Number(id);
  if (!ventaId || leidas.value.includes(ventaId)) {
    return;
  }

  leidas.value = [...leidas.value, ventaId];
  guardarLeidas();
}

function marcarTodasComoLeidas() {
  const idsPendientes = ventasPendientes.value.map(venta => Number(venta.id)).filter(Boolean);
  leidas.value = [...new Set([...leidas.value, ...idsPendientes])];
  guardarLeidas();
}

function onDocumentClick(event: MouseEvent) {
  if (!containerRef.value || containerRef.value.contains(event.target as Node)) {
    return;
  }

  open.value = false;
}

function formatPrecio(precio: number) {
  return Number(precio || 0).toLocaleString('es-BO', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  });
}

function formatFecha(fecha: string) {
  if (!fecha) {
    return 'N/A';
  }

  return new Date(`${fecha}T00:00:00`).toLocaleDateString('es-BO');
}

function severityVenta(venta: any) {
  if (venta.estado_venta === 'Completada') return 'success';
  if (venta.estado_venta === 'Anulada' || venta.proforma_vencida) return 'danger';
  return 'warning';
}
</script>
