<template>
  <Toast />

  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-3 xl:flex-row xl:items-end xl:justify-between">
      <div>
        <h2 class="text-2xl font-bold">Reportes Contables</h2>
        <p class="text-sm text-gray-500">Exportacion de ventas, creditos y reservas por periodo.</p>
      </div>

      <div class="grid grid-cols-1 gap-2 md:grid-cols-[11rem_11rem_13rem_auto_auto] md:items-end">
        <div class="flex flex-col gap-1">
          <label for="fecha_inicio" class="text-sm text-gray-600">Desde</label>
          <InputText id="fecha_inicio" v-model="fechaInicio" type="date" size="small" />
        </div>
        <div class="flex flex-col gap-1">
          <label for="fecha_fin" class="text-sm text-gray-600">Hasta</label>
          <InputText id="fecha_fin" v-model="fechaFin" type="date" size="small" />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm text-gray-600">Reporte</label>
          <Select v-model="tipoReporte" :options="tiposReporte" option-label="label" option-value="value" size="small" fluid />
        </div>
        <Button label="Excel" icon="pi pi-file-excel" size="small" severity="success" :disabled="filasReporte.length === 0" @click="exportarExcel" />
        <Button label="PDF" icon="pi pi-file-pdf" size="small" severity="danger" :disabled="filasReporte.length === 0" @click="exportarPDF" />
      </div>
    </div>

    <div class="grid grid-cols-1 gap-3 md:grid-cols-4">
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Registros</p>
        <p class="text-2xl font-bold text-gray-900">{{ filasReporte.length }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Total USD</p>
        <p class="text-2xl font-bold text-gray-900">$ {{ formatPrecio(totalUSD) }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Total BOB</p>
        <p class="text-2xl font-bold text-gray-900">Bs {{ formatPrecio(totalBOB) }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">TC promedio</p>
        <p class="text-2xl font-bold text-gray-900">{{ formatTipoCambio(tipoCambioPromedio) }}</p>
      </div>
    </div>

    <DataTable
      :value="filasReporte"
      :loading="loading"
      tableStyle="min-width: 86rem"
      size="small"
      stripedRows
      removableSort
      paginator
      :rows="10"
      :rowsPerPageOptions="[10, 20, 50]"
      currentPageReportTemplate="Mostrando {first} a {last} de {totalRecords} registros"
    >
      <template #empty>
        <p class="p-4 text-center">No hay registros para el periodo seleccionado.</p>
      </template>

      <Column field="fecha" header="Fecha" sortable>
        <template #body="slotProps">{{ formatFecha(slotProps.data.fecha) }}</template>
      </Column>
      <Column field="cliente" header="Cliente" sortable />
      <Column field="vehiculo" header="Vehiculo" sortable />
      <Column field="tipo_venta" header="Tipo de venta" sortable>
        <template #body="slotProps">{{ etiquetaTipoVenta(slotProps.data.tipo_venta) }}</template>
      </Column>
      <Column field="metodo_pago" header="Metodo de pago" sortable />
      <Column field="monto_usd" header="USD" sortable>
        <template #body="slotProps">$ {{ formatPrecio(slotProps.data.monto_usd) }}</template>
      </Column>
      <Column field="monto_bob" header="BOB" sortable>
        <template #body="slotProps">Bs {{ formatPrecio(slotProps.data.monto_bob) }}</template>
      </Column>
      <Column field="tipo_cambio" header="TC" sortable>
        <template #body="slotProps">{{ formatTipoCambio(slotProps.data.tipo_cambio) }}</template>
      </Column>
      <Column field="estado_pago" header="Estado de pago" sortable />
      <Column field="saldo_usd" header="Saldo USD" sortable>
        <template #body="slotProps">$ {{ formatPrecio(slotProps.data.saldo_usd) }}</template>
      </Column>
      <Column field="saldo_bob" header="Saldo BOB" sortable>
        <template #body="slotProps">Bs {{ formatPrecio(slotProps.data.saldo_bob) }}</template>
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
import Select from 'primevue/select';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';

definePageMeta({ layout: 'menu-admin' });

const toast = useToast();
const loading = ref(true);
const ventas = ref<any[]>([]);
const fechaInicio = ref(inicioMesActual());
const fechaFin = ref(new Date().toISOString().slice(0, 10));
const tipoReporte = ref('ventas');
const tiposReporte = [
  { label: 'Ventas', value: 'ventas' },
  { label: 'Creditos', value: 'creditos' },
  { label: 'Reservas', value: 'reservas' }
];

const ventasPeriodo = computed(() => ventas.value.filter((venta: any) =>
  fechaEnRango(venta.fecha, fechaInicio.value, fechaFin.value)
));

const filasReporte = computed(() => ventasPeriodo.value
  .filter((venta: any) => perteneceTipoReporte(venta, tipoReporte.value))
  .map((venta: any) => filaContable(venta))
);

const totalUSD = computed(() => filasReporte.value.reduce((total: number, fila: any) => total + Number(fila.monto_usd || 0), 0));
const totalBOB = computed(() => filasReporte.value.reduce((total: number, fila: any) => total + Number(fila.monto_bob || 0), 0));
const tipoCambioPromedio = computed(() => {
  const conTC = filasReporte.value.filter((fila: any) => Number(fila.tipo_cambio || 0) > 0);
  if (conTC.length === 0) return 0;
  return conTC.reduce((total: number, fila: any) => total + Number(fila.tipo_cambio || 0), 0) / conTC.length;
});

onMounted(async () => {
  await cargarVentas();
});

async function cargarVentas() {
  loading.value = true;
  try {
    const res = await $fetch(server.HOST + '/api/v1/ventas', { method: 'GET' });
    ventas.value = Array.isArray(res) ? res : [];
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar reportes', life: 3000 });
  } finally {
    loading.value = false;
  }
}

function perteneceTipoReporte(venta: any, tipo: string) {
  if (tipo === 'creditos') {
    return ['Credito', 'credito_directo', 'credito_bancario'].includes(venta.tipo_venta);
  }
  if (tipo === 'reservas') {
    return venta.tipo_venta === 'Reserva';
  }
  return !['Credito', 'credito_directo', 'credito_bancario', 'Reserva'].includes(venta.tipo_venta);
}

function filaContable(venta: any) {
  const tipoCambio = Number(venta.tipo_cambio_usado || 0);
  const saldoUSD = Number(venta.saldo || 0);
  const saldoBOB = Number(venta.saldo_bob || (tipoCambio > 0 ? saldoUSD * tipoCambio : 0));
  const montoUSD = montoContableUSD(venta);
  const montoBOB = tipoCambio > 0 ? montoUSD * tipoCambio : Number(venta.monto_bob_calculado || 0);

  return {
    id: venta.id,
    fecha: venta.fecha,
    cliente: venta.cliente || 'Sin nombre',
    vehiculo: venta.vehiculo || 'N/A',
    tipo_venta: venta.tipo_venta || 'N/A',
    metodo_pago: venta.metodo_pago || venta.detalle_pago || 'N/A',
    monto_usd: roundMoney(montoUSD),
    monto_bob: roundMoney(montoBOB),
    tipo_cambio: tipoCambio,
    estado_pago: venta.estado_pago || 'N/A',
    saldo_usd: roundMoney(saldoUSD),
    saldo_bob: roundMoney(saldoBOB),
    vendedor: venta.vendedor || 'Sin vendedor'
  };
}

function montoContableUSD(venta: any) {
  if (venta.tipo_venta === 'Reserva') {
    return Number(venta.cuota_inicial || venta.monto_reserva || 0);
  }
  if (['Credito', 'credito_directo', 'credito_bancario'].includes(venta.tipo_venta)) {
    return Number(venta.monto_financiado || venta.monto_financiar_banco || venta.saldo || 0);
  }
  return Number(venta.precio_total || 0);
}

function exportarExcel() {
  const html = tablaExportacionHTML();
  const blob = new Blob([`\ufeff${html}`], { type: 'application/vnd.ms-excel;charset=utf-8;' });
  descargarBlob(blob, `${nombreArchivoReporte()}.xls`);
}

function exportarPDF() {
  const printWindow = window.open('', '_blank', 'width=1100,height=800');
  if (!printWindow) {
    toast.add({ severity: 'warn', summary: 'No se pudo abrir la ventana de impresion', life: 3000 });
    return;
  }

  printWindow.document.write(`
    <!doctype html>
    <html>
      <head>
        <title>${tituloReporte()}</title>
        <style>
          body { font-family: Arial, sans-serif; color: #111827; margin: 24px; }
          h1 { font-size: 20px; margin: 0 0 4px; }
          p { margin: 2px 0; font-size: 12px; }
          table { width: 100%; border-collapse: collapse; margin-top: 14px; font-size: 10px; }
          th, td { border: 1px solid #d1d5db; padding: 6px; text-align: left; }
          th { background: #f3f4f6; }
          .totales { display: grid; grid-template-columns: repeat(4, 1fr); gap: 8px; margin-top: 12px; }
          .box { border: 1px solid #d1d5db; padding: 8px; font-size: 12px; }
        </style>
      </head>
      <body>${tablaExportacionHTML()}</body>
    </html>
  `);
  printWindow.document.close();
  printWindow.focus();
  printWindow.print();
}

function tablaExportacionHTML() {
  const filas = filasReporte.value.map((fila: any) => `
    <tr>
      <td>${escapeHtml(formatFecha(fila.fecha))}</td>
      <td>${escapeHtml(fila.cliente)}</td>
      <td>${escapeHtml(fila.vehiculo)}</td>
      <td>${escapeHtml(etiquetaTipoVenta(fila.tipo_venta))}</td>
      <td>${escapeHtml(fila.metodo_pago)}</td>
      <td>${formatPrecio(fila.monto_usd)}</td>
      <td>${formatPrecio(fila.monto_bob)}</td>
      <td>${formatTipoCambio(fila.tipo_cambio)}</td>
      <td>${escapeHtml(fila.estado_pago)}</td>
      <td>${formatPrecio(fila.saldo_usd)}</td>
      <td>${formatPrecio(fila.saldo_bob)}</td>
    </tr>
  `).join('');

  return `
    <h1>${escapeHtml(tituloReporte())}</h1>
    <p><strong>Periodo:</strong> ${escapeHtml(formatFecha(fechaInicio.value))} al ${escapeHtml(formatFecha(fechaFin.value))}</p>
    <div class="totales">
      <div class="box"><strong>Registros</strong><br>${filasReporte.value.length}</div>
      <div class="box"><strong>Total USD</strong><br>${formatPrecio(totalUSD.value)}</div>
      <div class="box"><strong>Total BOB</strong><br>${formatPrecio(totalBOB.value)}</div>
      <div class="box"><strong>TC promedio</strong><br>${formatTipoCambio(tipoCambioPromedio.value)}</div>
    </div>
    <table>
      <thead>
        <tr>
          <th>Fecha</th>
          <th>Cliente</th>
          <th>Vehiculo</th>
          <th>Tipo de venta</th>
          <th>Metodo de pago</th>
          <th>USD</th>
          <th>BOB</th>
          <th>TC</th>
          <th>Estado de pago</th>
          <th>Saldo USD</th>
          <th>Saldo BOB</th>
        </tr>
      </thead>
      <tbody>${filas}</tbody>
    </table>
  `;
}

function tituloReporte() {
  const actual = tiposReporte.find(item => item.value === tipoReporte.value);
  return `Reporte contable de ${actual?.label || 'registros'}`;
}

function nombreArchivoReporte() {
  return `reporte_${tipoReporte.value}_${fechaInicio.value}_${fechaFin.value}`;
}

function descargarBlob(blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
}

function fechaEnRango(fecha: string, inicio: string, fin: string) {
  if (!fecha) return false;
  const valor = fecha.slice(0, 10);
  return (!inicio || valor >= inicio) && (!fin || valor <= fin);
}

function inicioMesActual() {
  const hoy = new Date();
  return `${hoy.getFullYear()}-${String(hoy.getMonth() + 1).padStart(2, '0')}-01`;
}

function formatPrecio(precio: number) {
  return Number(precio || 0).toLocaleString('es-BO', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

function formatTipoCambio(value: number) {
  return Number(value || 0).toLocaleString('es-BO', { minimumFractionDigits: 2, maximumFractionDigits: 4 });
}

function formatFecha(fecha: string) {
  if (!fecha) return 'N/A';
  return new Date(`${fecha}T00:00:00`).toLocaleDateString('es-BO');
}

function etiquetaTipoVenta(tipo: string) {
  if (tipo === 'credito_directo') return 'Credito directo';
  if (tipo === 'credito_bancario') return 'Credito bancario';
  return tipo || 'N/A';
}

function roundMoney(value: number) {
  return Math.round(Number(value || 0) * 100) / 100;
}

function escapeHtml(value: unknown) {
  return String(value)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;');
}
</script>
