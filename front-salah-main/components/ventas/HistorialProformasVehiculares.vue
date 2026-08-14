<template>
  <Toast />

  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <h2 class="text-2xl font-bold">Historial de Proformas</h2>
        <p class="text-sm text-gray-500">Proformas vehiculares generadas desde el catalogo de ventas.</p>
      </div>

      <span class="p-input-icon-left">
        <i class="pi pi-search" />
        <InputText v-model="searchQuery" placeholder="Buscar..." size="small" />
      </span>
    </div>

    <div class="grid grid-cols-1 gap-3 md:grid-cols-3">
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Proformas</p>
        <p class="text-2xl font-bold text-gray-900">{{ filteredProformas.length }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Total ofertado</p>
        <p class="text-2xl font-bold text-gray-900">$ {{ formatPrecio(totalOfertado) }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Vendedores</p>
        <p class="text-2xl font-bold text-gray-900">{{ totalVendedores }}</p>
      </div>
    </div>

    <DataTable
      :value="filteredProformas"
      :loading="loading"
      tableStyle="min-width: 72rem"
      size="small"
      stripedRows
      removableSort
      paginator
      :rows="10"
      :rowsPerPageOptions="[10, 20, 50]"
    >
      <template #empty>
        <p class="p-4 text-center">No hay proformas registradas.</p>
      </template>

      <Column field="id" header="ID" sortable />
      <Column field="fecha" header="Fecha" sortable>
        <template #body="slotProps">{{ formatFechaHora(slotProps.data.fecha) }}</template>
      </Column>
      <Column field="cliente_nombre" header="Cliente" sortable>
        <template #body="slotProps">
          <div>
            <p class="font-medium text-gray-900">{{ slotProps.data.cliente_nombre || 'Sin nombre' }}</p>
            <p class="text-xs text-gray-500">{{ slotProps.data.cliente_telefono || 'Sin telefono' }}</p>
          </div>
        </template>
      </Column>
      <Column field="vehiculo" header="Vehiculo" sortable />
      <Column field="vendedor" header="Vendedor" sortable>
        <template #body="slotProps">{{ slotProps.data.vendedor || 'N/A' }}</template>
      </Column>
      <Column field="modalidad" header="Modalidad" sortable />
      <Column field="precio_unidad" header="Precio unidad" sortable>
        <template #body="slotProps">$ {{ formatPrecio(slotProps.data.precio_unidad) }}</template>
      </Column>
      <Column field="cantidad" header="Cant." sortable />
      <Column field="precio_total" header="Total" sortable>
        <template #body="slotProps">$ {{ formatPrecio(slotProps.data.precio_total) }}</template>
      </Column>
      <Column field="cuota_inicial" header="Inicial" sortable>
        <template #body="slotProps">$ {{ formatPrecio(slotProps.data.cuota_inicial) }}</template>
      </Column>
      <Column field="saldo" header="Saldo" sortable>
        <template #body="slotProps">$ {{ formatPrecio(slotProps.data.saldo) }}</template>
      </Column>
      <Column field="fecha_vencimiento" header="Vence" sortable>
        <template #body="slotProps">{{ formatFecha(slotProps.data.fecha_vencimiento) }}</template>
      </Column>
      <Column header="PDF">
        <template #body="slotProps">
          <Button icon="pi pi-file-pdf" size="small" severity="danger" text rounded aria-label="Descargar PDF" @click="descargarProforma(slotProps.data.id)" />
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import Button from 'primevue/button';
import Column from 'primevue/column';
import DataTable from 'primevue/datatable';
import InputText from 'primevue/inputtext';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';
import { server } from '~/server/server';

const toast = useToast();
const proformas = ref<any[]>([]);
const searchQuery = ref('');
const loading = ref(true);

const filteredProformas = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  if (!query) {
    return proformas.value;
  }

  return proformas.value.filter((proforma: any) =>
    (proforma.id?.toString() || '').includes(query) ||
    (proforma.fecha?.toLowerCase() || '').includes(query) ||
    (proforma.cliente_nombre?.toLowerCase() || '').includes(query) ||
    (proforma.cliente_telefono?.toLowerCase() || '').includes(query) ||
    (proforma.vehiculo?.toLowerCase() || '').includes(query) ||
    (proforma.vendedor?.toLowerCase() || '').includes(query) ||
    (proforma.modalidad?.toLowerCase() || '').includes(query)
  );
});

const totalOfertado = computed(() => filteredProformas.value.reduce((total: number, proforma: any) => total + Number(proforma.precio_total || 0), 0));
const totalVendedores = computed(() => new Set(filteredProformas.value.map((proforma: any) => proforma.vendedor).filter(Boolean)).size);

onMounted(async () => {
  await obtenerProformas();
});

async function obtenerProformas() {
  loading.value = true;
  try {
    const res = await $fetch(server.HOST + '/api/v1/proformas-vehiculares', { method: 'GET' });
    proformas.value = Array.isArray(res) ? res : [];
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar proformas', life: 3000 });
  } finally {
    loading.value = false;
  }
}

async function descargarProforma(idProforma: number | string) {
  try {
    const response = await fetch(`${server.HOST}/api/v1/reportes/proformas-vehiculares/${idProforma}`, { method: 'GET' });
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
  } catch (err: any) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al descargar PDF', detail: err?.message, life: 4000 });
  }
}

function formatPrecio(precio: number) {
  return Number(precio || 0).toLocaleString('es-BO', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

function formatFecha(fecha: string) {
  if (!fecha) {
    return 'N/A';
  }
  return new Date(`${fecha}T00:00:00`).toLocaleDateString('es-BO');
}

function formatFechaHora(fecha: string) {
  if (!fecha) {
    return 'N/A';
  }
  const [dia, hora] = fecha.split(' ');
  return `${formatFecha(dia)} ${hora?.slice(0, 5) || ''}`.trim();
}
</script>
