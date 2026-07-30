<template>
  <Toast />

  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <h2 class="text-2xl font-bold">Dashboard</h2>
        <p class="text-sm text-gray-500">Resumen general de ventas, clientes, usuarios y vehiculos.</p>
      </div>

      <Button label="Actualizar" icon="pi pi-refresh" size="small" :loading="loading" @click="cargarDashboard" />
    </div>

    <div class="grid grid-cols-1 gap-3 md:grid-cols-2 xl:grid-cols-4">
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <div class="flex items-center justify-between gap-3">
          <div>
            <p class="text-sm text-gray-500">Ventas completadas</p>
            <p class="text-2xl font-bold text-gray-900">{{ ventasCompletadas.length }}</p>
          </div>
          <i class="pi pi-shopping-cart text-2xl text-yellow-500"></i>
        </div>
      </div>

      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <div class="flex items-center justify-between gap-3">
          <div>
            <p class="text-sm text-gray-500">Ingresos USD</p>
            <p class="text-2xl font-bold text-gray-900">$ {{ formatPrecio(totalVendidoUSD) }}</p>
          </div>
          <i class="pi pi-wallet text-2xl text-emerald-600"></i>
        </div>
      </div>

      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <div class="flex items-center justify-between gap-3">
          <div>
            <p class="text-sm text-gray-500">Ingresos BOB</p>
            <p class="text-2xl font-bold text-gray-900">Bs {{ formatPrecio(totalVendidoBOB) }}</p>
          </div>
          <i class="pi pi-money-bill text-2xl text-sky-600"></i>
        </div>
      </div>

      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <div class="flex items-center justify-between gap-3">
          <div>
            <p class="text-sm text-gray-500">Margen potencial</p>
            <p class="text-2xl font-bold text-gray-900">$ {{ formatPrecio(margenInventario) }}</p>
          </div>
          <i class="pi pi-car text-2xl text-red-600"></i>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-3 md:grid-cols-2 xl:grid-cols-4">
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Reservas pendientes</p>
        <p class="text-2xl font-bold text-gray-900">{{ reservasPendientes.length }}</p>
        <p class="text-xs text-gray-500">$ {{ formatPrecio(totalReservadoUSD) }} reservado</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Saldo reservas</p>
        <p class="text-2xl font-bold text-gray-900">$ {{ formatPrecio(saldoReservasUSD) }}</p>
        <p class="text-xs text-gray-500">Bs {{ formatPrecio(saldoReservasBOB) }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Pagos mixtos</p>
        <p class="text-2xl font-bold text-gray-900">{{ ventasMixtas.length }}</p>
        <p class="text-xs text-gray-500">$ {{ formatPrecio(totalPagoMixtoUSD) }} · Bs {{ formatPrecio(totalPagoMixtoBOB) }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Tipo cambio promedio</p>
        <p class="text-2xl font-bold text-gray-900">{{ formatTipoCambio(tipoCambioPromedio) }}</p>
        <p class="text-xs text-gray-500">Ventas con TC: {{ ventasConTipoCambio.length }}</p>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-4 xl:grid-cols-[1fr_24rem]">
      <section class="rounded-lg border border-gray-200 bg-white p-4">
        <div class="mb-4 flex items-center justify-between gap-3">
          <h3 class="text-lg font-semibold text-gray-900">Ventas por Vendedor</h3>
          <span class="text-sm text-gray-500">{{ vendedores.length }} vendedores</span>
        </div>

        <div v-if="loading" class="flex flex-col gap-3">
          <Skeleton v-for="item in 4" :key="item" height="2.5rem" />
        </div>

        <div v-else-if="vendedores.length === 0" class="py-10 text-center text-gray-500">
          No hay ventas registradas.
        </div>

        <div v-else class="flex flex-col gap-3">
          <div v-for="vendedor in vendedores" :key="vendedor.nombre" class="grid grid-cols-[11rem_1fr_8rem] items-center gap-3">
            <div class="min-w-0">
              <p class="truncate text-sm font-medium text-gray-900">{{ vendedor.nombre }}</p>
              <p class="text-xs text-gray-500">{{ vendedor.cantidad }} ventas</p>
            </div>
            <div class="h-3 overflow-hidden rounded bg-gray-100">
              <div class="h-full rounded bg-yellow-500" :style="{ width: `${vendedor.porcentaje}%` }"></div>
            </div>
            <p class="text-right text-sm font-semibold text-gray-900">$ {{ formatPrecio(vendedor.total) }}</p>
          </div>
        </div>
      </section>

      <section class="rounded-lg border border-gray-200 bg-white p-4">
        <h3 class="mb-4 text-lg font-semibold text-gray-900">Estado General</h3>

        <div class="flex flex-col gap-3">
          <div class="flex items-center justify-between border-b border-gray-100 pb-3">
            <span class="text-sm text-gray-500">Usuarios</span>
            <span class="font-semibold text-gray-900">{{ usuarios.length }}</span>
          </div>
          <div class="flex items-center justify-between border-b border-gray-100 pb-3">
            <span class="text-sm text-gray-500">Clientes</span>
            <span class="font-semibold text-gray-900">{{ clientes.length }}</span>
          </div>
          <div class="flex items-center justify-between border-b border-gray-100 pb-3">
            <span class="text-sm text-gray-500">Vehiculos</span>
            <span class="font-semibold text-gray-900">{{ vehiculos.length }}</span>
          </div>
          <div class="flex items-center justify-between border-b border-gray-100 pb-3">
            <span class="text-sm text-gray-500">Clientes activos</span>
            <span class="font-semibold text-gray-900">{{ clientesActivos }}</span>
          </div>
          <div class="flex items-center justify-between border-b border-gray-100 pb-3">
            <span class="text-sm text-gray-500">Vehiculos activos</span>
            <span class="font-semibold text-gray-900">{{ vehiculosActivos }}</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-500">Reservas completadas</span>
            <span class="font-semibold text-gray-900">{{ reservasCompletadas.length }}</span>
          </div>
        </div>
      </section>
    </div>

    <section class="rounded-lg border border-gray-200 bg-white p-4">
      <div class="mb-4 flex items-center justify-between gap-3">
        <h3 class="text-lg font-semibold text-gray-900">Ventas Recientes</h3>
        <span class="text-sm text-gray-500">Ultimos registros</span>
      </div>

      <DataTable
        :value="ventasRecientes"
        :loading="loading"
        tableStyle="min-width: 64rem"
        size="small"
        stripedRows
      >
        <template #empty>
          <p class="p-4 text-center">No hay ventas recientes.</p>
        </template>

        <Column field="fecha" header="Fecha">
          <template #body="slotProps">
            {{ formatFecha(slotProps.data.fecha) }}
          </template>
        </Column>
        <Column field="vendedor" header="Vendedor">
          <template #body="slotProps">
            {{ slotProps.data.vendedor || 'Sin vendedor' }}
          </template>
        </Column>
        <Column field="cliente" header="Cliente" />
        <Column field="vehiculo" header="Vehiculo" />
        <Column field="estado_venta" header="Estado">
          <template #body="slotProps">
            <Tag :value="slotProps.data.estado_venta" :severity="slotProps.data.estado_venta === 'Completada' ? 'success' : 'warning'" />
          </template>
        </Column>
        <Column field="tipo_cambio_usado" header="TC">
          <template #body="slotProps">
            {{ formatTipoCambio(slotProps.data.tipo_cambio_usado) }}
          </template>
        </Column>
        <Column field="precio_total" header="USD">
          <template #body="slotProps">
            $ {{ formatPrecio(slotProps.data.precio_total) }}
          </template>
        </Column>
        <Column field="monto_bob_calculado" header="BOB">
          <template #body="slotProps">
            Bs {{ formatPrecio(slotProps.data.monto_bob_calculado) }}
          </template>
        </Column>
      </DataTable>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { server } from '~/server/server';
import Button from 'primevue/button';
import Column from 'primevue/column';
import DataTable from 'primevue/datatable';
import Skeleton from 'primevue/skeleton';
import Tag from 'primevue/tag';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';

definePageMeta({ layout: 'menu-admin' });

const toast = useToast();
const loading = ref(true);
const usuarios = ref<any[]>([]);
const clientes = ref<any[]>([]);
const vehiculos = ref<any[]>([]);
const ventas = ref<any[]>([]);

const ventasCompletadas = computed(() => ventas.value.filter(venta => venta.estado_venta === 'Completada' && venta.estado_pago === 'Pagado completo'));
const reservas = computed(() => ventas.value.filter(venta => venta.tipo_venta === 'Reserva'));
const reservasPendientes = computed(() => reservas.value.filter(venta => venta.estado_venta !== 'Completada' && venta.estado_venta !== 'Anulada'));
const reservasCompletadas = computed(() => reservas.value.filter(venta => venta.estado_venta === 'Completada' && venta.estado_pago === 'Pagado completo'));
const ventasMixtas = computed(() => ventas.value.filter(venta => venta.metodo_pago === 'Mixto'));
const ventasConTipoCambio = computed(() => ventas.value.filter(venta => Number(venta.tipo_cambio_usado || 0) > 0));
const totalVendidoUSD = computed(() => ventasCompletadas.value.reduce((total, venta) => total + Number(venta.precio_total || 0), 0));
const totalVendidoBOB = computed(() => ventasCompletadas.value.reduce((total, venta) => total + Number(venta.monto_bob_calculado || 0), 0));
const totalReservadoUSD = computed(() => reservasPendientes.value.reduce((total, venta) => total + Number(venta.cuota_inicial || 0), 0));
const saldoReservasUSD = computed(() => reservasPendientes.value.reduce((total, venta) => total + Number(venta.saldo || 0), 0));
const saldoReservasBOB = computed(() => reservasPendientes.value.reduce((total, venta) => total + Number(venta.saldo_bob || 0), 0));
const totalPagoMixtoUSD = computed(() => ventasMixtas.value.reduce((total, venta) => total + Number(venta.pago_usd || 0), 0));
const totalPagoMixtoBOB = computed(() => ventasMixtas.value.reduce((total, venta) => total + Number(venta.pago_bob || 0), 0));
const tipoCambioPromedio = computed(() => {
  if (ventasConTipoCambio.value.length === 0) return 0;
  return ventasConTipoCambio.value.reduce((total, venta) => total + Number(venta.tipo_cambio_usado || 0), 0) / ventasConTipoCambio.value.length;
});
const margenInventario = computed(() => vehiculos.value.reduce((total, vehiculo) => {
  const margenUnidad = Number(vehiculo.margen_ganancia ?? (Number(vehiculo.precio || 0) - Number(vehiculo.precio_compra || 0)));
  return total + (margenUnidad * Number(vehiculo.cantidad_disponible || 0));
}, 0));
const clientesActivos = computed(() => clientes.value.filter(cliente => cliente.estado === 'Activo').length);
const vehiculosActivos = computed(() => vehiculos.value.filter(vehiculo => vehiculo.estado === 'Activo').length);

const vendedores = computed(() => {
  const acumulado = new Map<string, { nombre: string; cantidad: number; total: number }>();

  for (const venta of ventasCompletadas.value) {
    const nombre = venta.vendedor || 'Sin vendedor';
    const actual = acumulado.get(nombre) || { nombre, cantidad: 0, total: 0 };
    actual.cantidad += 1;
    actual.total += Number(venta.precio_total || 0);
    acumulado.set(nombre, actual);
  }

  const maxTotal = Math.max(...Array.from(acumulado.values()).map(item => item.total), 1);
  return Array.from(acumulado.values())
    .sort((a, b) => b.total - a.total)
    .map(item => ({
      ...item,
      porcentaje: Math.max(6, Math.round((item.total / maxTotal) * 100))
    }));
});

const ventasRecientes = computed(() => ventas.value.slice(0, 8));

onMounted(async () => {
  await cargarDashboard();
});

async function cargarDashboard() {
  loading.value = true;

  try {
    const [resUsuarios, resClientes, resVehiculos, resVentas] = await Promise.all([
      $fetch(server.HOST + '/api/v1/usuarios', { method: 'GET' }),
      $fetch(server.HOST + '/api/v1/clientes', { method: 'GET' }),
      $fetch(server.HOST + '/api/v1/vehiculos', { method: 'GET' }),
      $fetch(server.HOST + '/api/v1/ventas', { method: 'GET' })
    ]);

    usuarios.value = Array.isArray(resUsuarios) ? resUsuarios : [];
    clientes.value = Array.isArray(resClientes) ? resClientes : [];
    vehiculos.value = Array.isArray(resVehiculos) ? resVehiculos : [];
    ventas.value = Array.isArray(resVentas) ? resVentas : [];
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar dashboard', life: 3000 });
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

function formatTipoCambio(value: number) {
  return Number(value || 0).toLocaleString('es-BO', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 4
  });
}
</script>
