<template>
  <Toast />

  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <h2 class="text-2xl font-bold">Historial General de Ventas</h2>
        <p class="text-sm text-gray-500">Lista de ventas registradas por todos los vendedores.</p>
      </div>

      <span class="p-input-icon-left">
        <i class="pi pi-search" />
        <InputText v-model="searchQuery" placeholder="Buscar..." size="small" />
      </span>
    </div>

    <div class="grid grid-cols-1 gap-3 md:grid-cols-3">
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Ventas</p>
        <p class="text-2xl font-bold text-gray-900">{{ filteredVentas.length }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Total vendido</p>
        <p class="text-2xl font-bold text-gray-900">$ {{ formatPrecio(totalVendido) }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Vendedores</p>
        <p class="text-2xl font-bold text-gray-900">{{ totalVendedores }}</p>
      </div>
    </div>

    <DataTable
      :value="filteredVentas"
      :loading="loading"
      tableStyle="min-width: 72rem"
      size="small"
      stripedRows
      removableSort
      paginator
      :rows="10"
      :rowsPerPageOptions="[10, 20, 50]"
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
      currentPageReportTemplate="Mostrando {first} a {last} de {totalRecords} ventas"
    >
      <template #empty>
        <p class="p-4 text-center">No hay ventas registradas.</p>
      </template>

      <Column field="id" header="ID" sortable />
      <Column field="fecha" header="Fecha" sortable>
        <template #body="slotProps">
          {{ formatFecha(slotProps.data.fecha) }}
        </template>
      </Column>
      <Column field="fecha_venta" header="Hora" sortable>
        <template #body="slotProps">
          {{ formatHora(slotProps.data.fecha_venta) }}
        </template>
      </Column>
      <Column field="vendedor" header="Vendedor" sortable>
        <template #body="slotProps">
          {{ slotProps.data.vendedor || 'Sin vendedor' }}
        </template>
      </Column>
      <Column field="cliente" header="Cliente" sortable>
        <template #body="slotProps">
          <div>
            <p class="font-medium text-gray-900">{{ slotProps.data.cliente || 'Sin nombre' }}</p>
            <p class="text-xs text-gray-500">CI/NIT: {{ slotProps.data.ci_cliente || 'N/A' }}</p>
          </div>
        </template>
      </Column>
      <Column field="vehiculo" header="Vehiculo" sortable>
        <template #body="slotProps">
          <div>
            <p class="font-medium text-gray-900">{{ slotProps.data.vehiculo }}</p>
            <p class="text-xs text-gray-500">{{ slotProps.data.categoria || 'Sin categoria' }}</p>
          </div>
        </template>
      </Column>
      <Column field="tipo_venta" header="Tipo" sortable />
      <Column field="estado_venta" header="Venta" sortable>
        <template #body="slotProps">
          <Tag :value="slotProps.data.estado_venta" :severity="severityVenta(slotProps.data)" />
        </template>
      </Column>
      <Column field="estado_pago" header="Pago" sortable />
      <Column field="metodo_pago" header="Metodo" sortable />
      <Column field="detalle_pago" header="Detalle pago">
        <template #body="slotProps">
          {{ formatDetallePago(slotProps.data) }}
        </template>
      </Column>
      <Column field="tipo_cambio_usado" header="TC" sortable>
        <template #body="slotProps">
          {{ formatPrecio(slotProps.data.tipo_cambio_usado) }}
        </template>
      </Column>
      <Column field="monto_bob_calculado" header="Total BOB" sortable>
        <template #body="slotProps">
          Bs {{ formatPrecio(slotProps.data.monto_bob_calculado) }}
        </template>
      </Column>
      <Column field="estado_entrega" header="Entrega" sortable />
      <Column field="cantidad" header="Cant." sortable />
      <Column field="precio_total" header="Total" sortable>
        <template #body="slotProps">
          $ {{ formatPrecio(slotProps.data.precio_total) }}
        </template>
      </Column>
      <Column field="cuota_inicial" header="Inicial/Reserva" sortable>
        <template #body="slotProps">
          {{ esCreditoOReserva(slotProps.data.tipo_venta) ? '$ ' + formatPrecio(slotProps.data.cuota_inicial) : '-' }}
        </template>
      </Column>
      <Column field="saldo" header="Saldo/Resta" sortable>
        <template #body="slotProps">
          {{ esCreditoOReserva(slotProps.data.tipo_venta) ? '$ ' + formatPrecio(slotProps.data.saldo) : '-' }}
        </template>
      </Column>
      <Column field="fecha_vencimiento_proforma" header="Vence" sortable>
        <template #body="slotProps">
          <span :class="slotProps.data.proforma_vencida ? 'font-semibold text-red-600' : ''">
            {{ formatFecha(slotProps.data.fecha_vencimiento_proforma) }}
          </span>
        </template>
      </Column>
      <Column field="observacion" header="Observacion">
        <template #body="slotProps">
          {{ slotProps.data.observacion || 'Sin observacion' }}
        </template>
      </Column>
      <Column header="PDF">
        <template #body="slotProps">
          <Button icon="pi pi-file-pdf" size="small" severity="danger" text rounded aria-label="Descargar PDF" @click="descargarVenta(slotProps.data.id)" />
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { server } from '~/server/server';
import Button from 'primevue/button';
import Column from 'primevue/column';
import DataTable from 'primevue/datatable';
import InputText from 'primevue/inputtext';
import Tag from 'primevue/tag';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';

definePageMeta({ layout: 'menu-ventas' });

const toast = useToast();
const { descargarPDFVenta } = useVentaPdf();
const ventas = ref<any[]>([]);
const searchQuery = ref('');
const loading = ref(true);

const filteredVentas = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  if (!query) {
    return ventas.value;
  }

  return ventas.value.filter((venta: any) =>
    (venta.id?.toString() || '').includes(query) ||
    (venta.fecha?.toLowerCase() || '').includes(query) ||
    (venta.fecha_venta?.toLowerCase() || '').includes(query) ||
    (venta.vendedor?.toLowerCase() || '').includes(query) ||
    (venta.cliente?.toLowerCase() || '').includes(query) ||
    (venta.ci_cliente?.toLowerCase() || '').includes(query) ||
    (venta.vehiculo?.toLowerCase() || '').includes(query) ||
    (venta.categoria?.toLowerCase() || '').includes(query) ||
    (venta.tipo_venta?.toLowerCase() || '').includes(query) ||
    (venta.estado_venta?.toLowerCase() || '').includes(query) ||
    (venta.estado_pago?.toLowerCase() || '').includes(query) ||
    (venta.metodo_pago?.toLowerCase() || '').includes(query) ||
    (venta.detalle_pago?.toLowerCase() || '').includes(query) ||
    (venta.estado_entrega?.toLowerCase() || '').includes(query)
  );
});

const totalVendido = computed(() => {
  return filteredVentas.value.reduce((total: number, venta: any) => total + Number(venta.precio_total || 0), 0);
});

const totalVendedores = computed(() => {
  const vendedores = filteredVentas.value
    .map((venta: any) => venta.vendedor)
    .filter(Boolean);
  return new Set(vendedores).size;
});

onMounted(async () => {
  await obtenerHistorialGeneral();
});

async function obtenerHistorialGeneral() {
  loading.value = true;

  try {
    const res = await $fetch(server.HOST + '/api/v1/ventas', {
      method: 'GET'
    });

    ventas.value = filtrarVentasCompletadas(Array.isArray(res) ? res : []);
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar historial general', life: 3000 });
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

function formatFecha(fecha: string) {
  if (!fecha) {
    return 'N/A';
  }

  return new Date(`${fecha}T00:00:00`).toLocaleDateString('es-BO');
}

function formatHora(fechaHora: string) {
  if (!fechaHora) {
    return 'N/A';
  }

  const partes = fechaHora.split(' ');
  return partes[1]?.slice(0, 5) || 'N/A';
}

function severityVenta(venta: any) {
  if (venta.estado_venta === 'Completada' || venta.estado_venta === 'pagado_completo') return 'success';
  if (venta.estado_venta === 'Anulada') return 'danger';
  if (venta.proforma_vencida) return 'danger';
  return 'warning';
}

async function descargarVenta(idVenta: number) {
  try {
    await descargarPDFVenta(idVenta);
  } catch (err: any) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al generar PDF', detail: err?.message, life: 4000 });
  }
}

function filtrarVentasCompletadas(items: any[]) {
  return items.filter((venta: any) =>
    (venta.estado_venta === 'Completada' && venta.estado_pago === 'Pagado completo') ||
    (venta.estado_venta === 'Registrada' && venta.estado_pago === 'Parcial') ||
    venta.estado_venta === 'en_credito' ||
    venta.estado_venta === 'pagado_completo'
  );
}

function esCreditoOReserva(tipoVenta: string) {
  return ['Credito', 'credito_directo', 'credito_bancario', 'Reserva'].includes(tipoVenta);
}

function formatDetallePago(venta: any) {
  if (venta.detalle_pago) {
    return venta.detalle_pago;
  }
  return venta.metodo_pago || 'N/A';
}
</script>
