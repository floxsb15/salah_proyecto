<template>
  <Toast />

  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <h2 class="text-2xl font-bold">Historial de Reservas</h2>
        <p class="text-sm text-gray-500">Reservas registradas y pagos pendientes para completar la venta.</p>
      </div>

      <div class="flex flex-col gap-2 sm:flex-row sm:items-center">
        <SelectButton v-model="estadoFiltro" :options="estadoFiltros" option-label="label" option-value="value" size="small" />
        <span class="p-input-icon-left">
          <i class="pi pi-search" />
          <InputText v-model="searchQuery" placeholder="Buscar..." size="small" />
        </span>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-3 md:grid-cols-3">
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Reservas</p>
        <p class="text-2xl font-bold text-gray-900">{{ filteredReservas.length }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Reservado</p>
        <p class="text-2xl font-bold text-gray-900">$ {{ formatPrecio(totalReservado) }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Saldo pendiente</p>
        <p class="text-2xl font-bold text-gray-900">$ {{ formatPrecio(totalPendiente) }}</p>
      </div>
    </div>

    <DataTable
      :value="filteredReservas"
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
        <p class="p-4 text-center">No hay reservas registradas.</p>
      </template>

      <Column field="id" header="ID" sortable />
      <Column field="fecha" header="Fecha" sortable>
        <template #body="slotProps">{{ formatFecha(slotProps.data.fecha) }}</template>
      </Column>
      <Column field="fecha_venta" header="Hora" sortable>
        <template #body="slotProps">{{ formatHora(slotProps.data.fecha_venta) }}</template>
      </Column>
      <Column v-if="general" field="vendedor" header="Vendedor" sortable />
      <Column field="cliente" header="Cliente" sortable>
        <template #body="slotProps">
          <div>
            <p class="font-medium text-gray-900">{{ slotProps.data.cliente || 'Sin nombre' }}</p>
            <p class="text-xs text-gray-500">CI/NIT: {{ slotProps.data.ci_cliente || 'N/A' }}</p>
          </div>
        </template>
      </Column>
      <Column field="tipo_reserva" header="Tipo" sortable>
        <template #body="slotProps">
          <Tag :value="tipoReservaLabel(slotProps.data)" :severity="tipoReservaSeverity(slotProps.data)" />
        </template>
      </Column>
      <Column field="vehiculo" header="Vehiculo" sortable>
        <template #body="slotProps">
          <div>
            <p class="font-medium text-gray-900">{{ vehiculoReserva(slotProps.data) }}</p>
            <p v-if="slotProps.data.tipo_reserva === 'pedido'" class="text-xs text-gray-500">
              {{ detallePedido(slotProps.data) }}
            </p>
          </div>
        </template>
      </Column>
      <Column field="estado_venta" header="Estado" sortable>
        <template #body="slotProps">
          <Tag :value="slotProps.data.estado_venta" :severity="severityReserva(slotProps.data)" />
        </template>
      </Column>
      <Column field="usuario_pago_reserva" header="Pago aceptado por" sortable>
        <template #body="slotProps">
          {{ slotProps.data.usuario_pago_reserva || '-' }}
        </template>
      </Column>
      <Column field="precio_total" header="Total" sortable>
        <template #body="slotProps">$ {{ formatPrecio(slotProps.data.precio_total) }}</template>
      </Column>
      <Column field="cuota_inicial" header="Reserva pagada" sortable>
        <template #body="slotProps">$ {{ formatPrecio(slotProps.data.cuota_inicial) }}</template>
      </Column>
      <Column field="saldo" header="Resta" sortable>
        <template #body="slotProps">$ {{ formatPrecio(slotProps.data.saldo) }}</template>
      </Column>
      <Column field="fecha_vencimiento_proforma" header="Vence" sortable>
        <template #body="slotProps">
          <span :class="slotProps.data.proforma_vencida ? 'font-semibold text-red-600' : ''">
            {{ formatFecha(slotProps.data.fecha_vencimiento_proforma) }}
          </span>
        </template>
      </Column>
      <Column header="Acciones">
        <template #body="slotProps">
          <div class="flex items-center gap-1">
            <Button
              label="Completar"
              icon="pi pi-check"
              size="small"
              :disabled="slotProps.data.estado_venta === 'Completada' || slotProps.data.estado_venta === 'Anulada' || Number(slotProps.data.saldo || 0) <= 0"
              @click="router.push({ path: completarPath, query: { id: slotProps.data.id } })"
            />
            <Button icon="pi pi-file-pdf" size="small" severity="danger" text rounded aria-label="Descargar PDF" @click="descargarVenta(slotProps.data.id)" />
          </div>
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { server } from '~/server/server';
import Button from 'primevue/button';
import Column from 'primevue/column';
import DataTable from 'primevue/datatable';
import InputText from 'primevue/inputtext';
import SelectButton from 'primevue/selectbutton';
import Tag from 'primevue/tag';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';

const props = withDefaults(defineProps<{
  general?: boolean;
  completarPath?: string;
}>(), {
  general: false,
  completarPath: '/ventas/completar-reserva'
});

const router = useRouter();
const toast = useToast();
const { descargarPDFVenta } = useVentaPdf();
const ventas = ref<any[]>([]);
const searchQuery = ref('');
const estadoFiltro = ref('pendientes');
const estadoFiltros = ref([
  { label: 'Pendientes', value: 'pendientes' },
  { label: 'Completadas', value: 'completadas' },
  { label: 'Todas', value: 'todas' }
]);
const loading = ref(true);

const reservas = computed(() => ventas.value.filter((venta: any) => venta.tipo_venta === 'Reserva'));
const filteredReservas = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  const porEstado = reservas.value.filter((venta: any) => {
    if (estadoFiltro.value === 'completadas') {
      return venta.estado_venta === 'Completada' && venta.estado_pago === 'Pagado completo';
    }
    if (estadoFiltro.value === 'pendientes') {
      return venta.estado_venta !== 'Completada' && venta.estado_venta !== 'Anulada';
    }
    return true;
  });
  if (!query) return porEstado;
  return porEstado.filter((venta: any) =>
    (venta.id?.toString() || '').includes(query) ||
    (venta.fecha?.toLowerCase() || '').includes(query) ||
    (venta.fecha_venta?.toLowerCase() || '').includes(query) ||
    (venta.vendedor?.toLowerCase() || '').includes(query) ||
    (venta.cliente?.toLowerCase() || '').includes(query) ||
    (venta.ci_cliente?.toLowerCase() || '').includes(query) ||
    (venta.vehiculo?.toLowerCase() || '').includes(query) ||
    (venta.tipo_reserva?.toLowerCase() || '').includes(query) ||
    (venta.pedido_marca?.toLowerCase() || '').includes(query) ||
    (venta.pedido_modelo?.toLowerCase() || '').includes(query) ||
    (venta.pedido_color?.toLowerCase() || '').includes(query) ||
    (venta.pedido_pais_origen?.toLowerCase() || '').includes(query) ||
    (venta.pedido_proveedor?.toLowerCase() || '').includes(query) ||
    (venta.usuario_pago_reserva?.toLowerCase() || '').includes(query) ||
    (venta.estado_venta?.toLowerCase() || '').includes(query)
  );
});
const totalReservado = computed(() => filteredReservas.value.reduce((total: number, venta: any) => total + Number(venta.cuota_inicial || 0), 0));
const totalPendiente = computed(() => filteredReservas.value.reduce((total: number, venta: any) => total + Number(venta.saldo || 0), 0));

onMounted(async () => {
  await obtenerReservas();
});

async function obtenerReservas() {
  loading.value = true;
  try {
    const query: Record<string, any> = {};
    if (!props.general) {
      const user = localStorage.getItem('user');
      const userId = user ? JSON.parse(user)?.id : null;
      if (!userId) {
        ventas.value = [];
        return;
      }
      query.id_usuario = userId;
    }
    const res = await $fetch(server.HOST + '/api/v1/ventas', { method: 'GET', query });
    ventas.value = Array.isArray(res) ? res : [];
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar reservas', life: 3000 });
  } finally {
    loading.value = false;
  }
}

function formatPrecio(precio: number) {
  return Number(precio || 0).toLocaleString('es-BO', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

function formatFecha(fecha: string) {
  if (!fecha) return 'N/A';
  return new Date(`${fecha}T00:00:00`).toLocaleDateString('es-BO');
}

function formatHora(fechaHora: string) {
  if (!fechaHora) return 'N/A';
  const partes = fechaHora.split(' ');
  return partes[1]?.slice(0, 5) || 'N/A';
}

function severityReserva(venta: any) {
  if (venta.estado_venta === 'Completada') return 'success';
  if (venta.estado_venta === 'Anulada' || venta.proforma_vencida) return 'danger';
  if (venta.estado_venta === 'Importando') return 'info';
  return 'warning';
}

function tipoReservaLabel(venta: any) {
  return venta.tipo_reserva === 'pedido' ? 'A pedido' : 'Stock';
}

function tipoReservaSeverity(venta: any) {
  return venta.tipo_reserva === 'pedido' ? 'warning' : 'info';
}

function vehiculoReserva(venta: any) {
  if (venta.tipo_reserva !== 'pedido') {
    return venta.vehiculo || 'Vehiculo';
  }
  return [venta.pedido_marca, venta.pedido_modelo, venta.pedido_anio].filter(Boolean).join(' ') || venta.vehiculo || 'Vehiculo a pedido';
}

function detallePedido(venta: any) {
  const partes = [
    venta.pedido_color ? `Color: ${venta.pedido_color}` : '',
    venta.pedido_pais_origen ? `Origen: ${venta.pedido_pais_origen}` : '',
    venta.pedido_llegada_estimada ? `Llegada: ${venta.pedido_llegada_estimada}` : ''
  ].filter(Boolean);
  return partes.join(' / ') || venta.pedido_version || 'Importacion pendiente';
}

async function descargarVenta(idVenta: number) {
  try {
    await descargarPDFVenta(idVenta);
  } catch (err: any) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al generar PDF', detail: err?.message, life: 4000 });
  }
}
</script>
