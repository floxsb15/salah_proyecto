<template>
  <Toast />

  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <h2 class="text-2xl font-bold">Catalogo de Vehiculos</h2>

      <div class="flex flex-col gap-2 sm:flex-row sm:items-center">
        <span class="p-input-icon-left">
          <i class="pi pi-search" />
          <InputText v-model="searchQuery" placeholder="Buscar..." size="small" />
        </span>

        <Select
          v-model="categoriaSeleccionada"
          :options="categorias"
          placeholder="Categoria"
          show-clear
          size="small"
          class="min-w-48"
        />

        <Select
          v-model="segmentoSeleccionado"
          :options="segmentos"
          placeholder="Segmento"
          show-clear
          size="small"
          class="min-w-48"
        />
      </div>
    </div>

    <div v-if="loading" class="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
      <Skeleton v-for="item in 6" :key="item" height="18rem" />
    </div>

    <div v-else-if="filteredVehiculos.length === 0" class="flex flex-col items-center justify-center gap-2 py-16 text-gray-500">
      <i class="pi pi-car text-4xl"></i>
      <p>No hay vehiculos para mostrar.</p>
    </div>

    <div v-else class="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
      <article
        v-for="vehiculo in filteredVehiculos"
        :key="vehiculo.id"
        class="overflow-hidden rounded-lg border border-gray-200 bg-white shadow-sm"
      >
        <div class="relative flex aspect-[16/10] items-center justify-center bg-gray-100">
          <img
            v-if="imagenesVehiculo(vehiculo).length"
            :src="imagenActual(vehiculo)"
            :alt="etiquetaVehiculo(vehiculo)"
            class="h-full w-full object-cover"
          >
          <i v-else class="pi pi-car text-5xl text-gray-400"></i>
          <div
            v-if="imagenesVehiculo(vehiculo).length > 1"
            class="absolute bottom-2 left-0 right-0 flex justify-center gap-1"
          >
            <span
              v-for="(_, index) in imagenesVehiculo(vehiculo)"
              :key="index"
              class="h-1.5 w-1.5 rounded-full bg-white/70 shadow"
              :class="index === indiceImagenActual(vehiculo) ? 'opacity-100' : 'opacity-40'"
            ></span>
          </div>
        </div>

        <div class="flex flex-col gap-3 p-4">
          <div class="flex items-start justify-between gap-3">
            <div class="min-w-0">
              <h3 class="truncate text-lg font-semibold text-gray-900">{{ etiquetaVehiculo(vehiculo) }}</h3>
              <p class="text-sm text-gray-500">
                {{ vehiculo.categoria || 'Sin categoria' }}
                <span v-if="vehiculo.segmento"> / {{ vehiculo.segmento }}</span>
              </p>
              <p v-if="vehiculo.version" class="text-xs text-gray-500">
                {{ vehiculo.version }}
              </p>
            </div>
            <Tag :value="vehiculo.estado" :severity="vehiculo.estado === 'Activo' ? 'success' : 'danger'" />
          </div>

          <div v-if="vehiculo.combustible || vehiculo.traccion || vehiculo.asientos" class="grid grid-cols-3 gap-2 text-xs text-gray-600">
            <span v-if="vehiculo.combustible">{{ vehiculo.combustible }}</span>
            <span v-if="vehiculo.traccion">{{ vehiculo.traccion }}</span>
            <span v-if="vehiculo.asientos">{{ vehiculo.asientos }} asientos</span>
          </div>

          <div class="flex items-center justify-between border-t border-gray-100 pt-3">
            <span class="text-sm text-gray-500">Precio</span>
            <span class="text-xl font-bold text-gray-900">$ {{ formatPrecio(vehiculo.precio) }}</span>
          </div>
          <div class="grid grid-cols-2 gap-2 text-sm">
            <div class="rounded-md bg-gray-50 p-2">
              <span class="text-gray-500">Compra</span>
              <strong class="block text-gray-900">$ {{ formatPrecio(vehiculo.precio_compra) }}</strong>
            </div>
            <div class="rounded-md bg-gray-50 p-2">
              <span class="text-gray-500">Margen</span>
              <strong class="block text-gray-900">$ {{ formatPrecio(vehiculo.margen_ganancia) }}</strong>
            </div>
          </div>
          <div class="flex items-center justify-between text-sm">
            <span class="text-gray-500">Disponible</span>
            <Tag :value="String(vehiculo.cantidad_disponible || 0)" :severity="Number(vehiculo.cantidad_disponible || 0) > 0 ? 'info' : 'danger'" />
          </div>

          <Button
            label="Registrar venta"
            icon="pi pi-shopping-cart"
            size="small"
            class="w-full"
            :disabled="vehiculo.estado !== 'Activo' || Number(vehiculo.cantidad_disponible || 0) <= 0"
            @click="irACompra(vehiculo.id)"
          />
        </div>
      </article>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import { server } from '~/server/server';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import Select from 'primevue/select';
import Skeleton from 'primevue/skeleton';
import Tag from 'primevue/tag';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';

definePageMeta({ layout: 'menu-admin' });

const toast = useToast();
const router = useRouter();
const vehiculos = ref<any[]>([]);
const searchQuery = ref('');
const topSearch = useTopSearch();
const categoriaSeleccionada = ref<string | null>(null);
const segmentoSeleccionado = ref<string | null>(null);
const loading = ref(true);
const indiceCarrusel = ref(0);
let carruselTimer: ReturnType<typeof setInterval> | null = null;

const categorias = computed(() => {
  const values = vehiculos.value
    .map(v => v.categoria)
    .filter(Boolean);
  return [...new Set(values)].sort((a, b) => a.localeCompare(b));
});

const segmentos = computed(() => {
  const values = vehiculos.value
    .filter(v => !categoriaSeleccionada.value || v.categoria === categoriaSeleccionada.value)
    .map(v => v.segmento)
    .filter(Boolean);
  return [...new Set(values)].sort((a, b) => a.localeCompare(b));
});

watch(categoriaSeleccionada, () => {
  if (segmentoSeleccionado.value && !segmentos.value.includes(segmentoSeleccionado.value)) {
    segmentoSeleccionado.value = null;
  }
});

const filteredVehiculos = computed(() => {
  const query = (topSearch.value || searchQuery.value).trim().toLowerCase();

  return vehiculos.value.filter(v => {
    const matchCategoria = !categoriaSeleccionada.value || v.categoria === categoriaSeleccionada.value;
    const matchSegmento = !segmentoSeleccionado.value || v.segmento === segmentoSeleccionado.value;
    const matchBusqueda =
      !query ||
      etiquetaVehiculo(v).toLowerCase().includes(query) ||
      (v.categoria?.toLowerCase() || '').includes(query) ||
      (v.segmento?.toLowerCase() || '').includes(query) ||
      (v.marca?.toLowerCase() || '').includes(query) ||
      (v.modelo?.toLowerCase() || '').includes(query) ||
      (v.anio?.toString() || '').includes(query) ||
      (v.cantidad_disponible?.toString() || '').includes(query) ||
      (v.version?.toLowerCase() || '').includes(query) ||
      (v.combustible?.toLowerCase() || '').includes(query);

    return matchCategoria && matchSegmento && matchBusqueda;
  });
});

onMounted(async () => {
  await obtenerVehiculos();
  carruselTimer = setInterval(() => {
    indiceCarrusel.value += 1;
  }, 3000);
});

onUnmounted(() => {
  if (carruselTimer) {
    clearInterval(carruselTimer);
  }
});

async function obtenerVehiculos() {
  loading.value = true;

  try {
    const res = await $fetch(server.HOST + '/api/v1/vehiculos', {
      method: 'GET',
      headers: getAuthHeaders()
    });

    vehiculos.value = Array.isArray(res) ? res : [];
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar vehiculos', life: 3000 });
  } finally {
    loading.value = false;
  }
}

function formatPrecio(precio: number) {
  return Number(precio || 0).toLocaleString('es-BO', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  });
}

function irACompra(idVehiculo: number) {
  router.push({
    path: '/admin/venta-vehiculo',
    query: { id: idVehiculo }
  });
}

function etiquetaVehiculo(vehiculo: any) {
  return [vehiculo.marca, vehiculo.modelo, vehiculo.anio].filter(Boolean).join(' ') || vehiculo.nombre || 'Vehiculo';
}

function imagenesVehiculo(vehiculo: any) {
  if (Array.isArray(vehiculo.imagenes) && vehiculo.imagenes.length > 0) {
    return vehiculo.imagenes.filter((imagen: string) => imagen && imagen !== 'N/A');
  }
  return vehiculo.imagen && vehiculo.imagen !== 'N/A' ? [vehiculo.imagen] : [];
}

function indiceImagenActual(vehiculo: any) {
  const imagenes = imagenesVehiculo(vehiculo);
  return imagenes.length > 0 ? indiceCarrusel.value % imagenes.length : 0;
}

function imagenActual(vehiculo: any) {
  const imagenes = imagenesVehiculo(vehiculo);
  return imagenes[indiceImagenActual(vehiculo)] || '';
}

function getAuthHeaders() {
  try {
    const user = localStorage.getItem('user');
    const parsed = user ? JSON.parse(user) : null;
    const headers: Record<string, string> = {};
    if (parsed?.token) {
      headers.Authorization = `Bearer ${parsed.token}`;
    }
    if (parsed?.id) {
      headers['X-User-Id'] = String(parsed.id);
    }
    return headers;
  } catch {
    return {};
  }
}
</script>
