<template>
  <Toast />

  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-3 xl:flex-row xl:items-center xl:justify-between">
      <div>
        <h2 class="text-2xl font-bold">Modulo de Creditos</h2>
        <p class="text-sm text-gray-500">{{ descripcion }}</p>
      </div>

      <div class="flex flex-col gap-2 md:flex-row md:items-center">
        <SelectButton v-model="tipoFiltro" :options="tipoFiltros" option-label="label" option-value="value" size="small" />
        <SelectButton v-model="estadoFiltro" :options="estadoFiltros" option-label="label" option-value="value" size="small" />
        <span class="p-input-icon-left">
          <i class="pi pi-search" />
          <InputText v-model="searchQuery" placeholder="Buscar credito..." size="small" />
        </span>
        <Button icon="pi pi-refresh" size="small" severity="secondary" :loading="loading" aria-label="Actualizar" @click="obtenerCreditos" />
      </div>
    </div>

    <div class="grid grid-cols-1 gap-3 md:grid-cols-4">
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Creditos</p>
        <p class="text-2xl font-bold text-gray-900">{{ filteredCreditos.length }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Financiado</p>
        <p class="text-2xl font-bold text-gray-900">$ {{ formatPrecio(totalFinanciado) }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Saldo pendiente</p>
        <p class="text-2xl font-bold text-gray-900">$ {{ formatPrecio(totalSaldoPendiente) }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Cuotas atrasadas</p>
        <p class="text-2xl font-bold text-red-600">{{ totalCuotasAtrasadas }}</p>
      </div>
    </div>

    <DataTable
      :value="filteredCreditos"
      :loading="loading"
      tableStyle="min-width: 78rem"
      size="small"
      stripedRows
      removableSort
      paginator
      :rows="10"
      :rowsPerPageOptions="[10, 20, 50]"
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
      currentPageReportTemplate="Mostrando {first} a {last} de {totalRecords} creditos"
    >
      <template #empty>
        <p class="p-4 text-center">No hay creditos registrados.</p>
      </template>

      <Column field="id" header="ID" sortable />
      <Column field="fecha" header="Fecha" sortable>
        <template #body="slotProps">{{ formatFecha(slotProps.data.fecha) }}</template>
      </Column>
      <Column field="fecha_venta" header="Hora" sortable>
        <template #body="slotProps">{{ formatHora(slotProps.data.fecha_venta) }}</template>
      </Column>
      <Column v-if="mostrarVendedor" field="vendedor" header="Vendedor" sortable />
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
      <Column field="tipo_venta" header="Tipo" sortable>
        <template #body="slotProps">
          <Tag :value="labelTipoCredito(slotProps.data.tipo_venta)" severity="info" />
        </template>
      </Column>
      <Column field="estado_venta" header="Estado" sortable>
        <template #body="slotProps">
          <Tag :value="labelEstadoCredito(slotProps.data.estado_venta)" :severity="severityEstado(slotProps.data.estado_venta)" />
        </template>
      </Column>
      <Column field="precio_total" header="Total" sortable>
        <template #body="slotProps">$ {{ formatPrecio(slotProps.data.precio_total) }}</template>
      </Column>
      <Column field="cuota_inicial" header="Inicial" sortable>
        <template #body="slotProps">$ {{ formatPrecio(slotProps.data.cuota_inicial) }}</template>
      </Column>
      <Column field="monto_financiado" header="Financiado" sortable>
        <template #body="slotProps">$ {{ formatPrecio(slotProps.data.monto_financiado || slotProps.data.saldo) }}</template>
      </Column>
      <Column field="saldo" header="Saldo" sortable>
        <template #body="slotProps">$ {{ formatPrecio(slotProps.data.saldo) }}</template>
      </Column>
      <Column field="numero_cuotas" header="Cuotas" sortable>
        <template #body="slotProps">
          {{ slotProps.data.cuotas_pagadas || 0 }}/{{ slotProps.data.numero_cuotas || 0 }}
        </template>
      </Column>
      <Column field="monto_cuota" header="Monto cuota" sortable>
        <template #body="slotProps">$ {{ formatPrecio(slotProps.data.monto_cuota) }}</template>
      </Column>
      <Column field="frecuencia_pago" header="Frecuencia" sortable>
        <template #body="slotProps">{{ slotProps.data.frecuencia_pago || 'mensual' }}</template>
      </Column>
      <Column field="proxima_cuota_fecha" header="Proxima cuota" sortable>
        <template #body="slotProps">
          <div>
            <p :class="slotProps.data.tiene_atraso ? 'font-semibold text-red-600' : 'text-gray-900'">
              {{ slotProps.data.proxima_cuota_fecha ? formatFecha(slotProps.data.proxima_cuota_fecha) : 'Sin pendientes' }}
            </p>
            <p v-if="slotProps.data.proxima_cuota_monto" class="text-xs text-gray-500">
              $ {{ formatPrecio(slotProps.data.proxima_cuota_monto) }}
            </p>
          </div>
        </template>
      </Column>
      <Column header="Seguimiento">
        <template #body="slotProps">
          <Button label="Ver cuotas" icon="pi pi-list-check" size="small" text @click="abrirDetalle(slotProps.data)" />
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="detalleVisible" modal header="Seguimiento del credito" :style="{ width: '92rem', maxWidth: '98vw' }">
      <div v-if="creditoSeleccionado" class="flex flex-col gap-4">
        <div class="grid grid-cols-1 gap-3 md:grid-cols-4">
          <div class="rounded-lg border border-gray-200 p-3">
            <p class="text-xs uppercase text-gray-500">Cliente</p>
            <p class="font-semibold text-gray-900">{{ creditoSeleccionado.cliente || 'Sin nombre' }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 p-3">
            <p class="text-xs uppercase text-gray-500">Vehiculo</p>
            <p class="font-semibold text-gray-900">{{ creditoSeleccionado.vehiculo }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 p-3">
            <p class="text-xs uppercase text-gray-500">Saldo</p>
            <p class="font-semibold text-gray-900">$ {{ formatPrecio(creditoSeleccionado.saldo) }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 p-3">
            <p class="text-xs uppercase text-gray-500">Progreso</p>
            <p class="font-semibold text-gray-900">{{ creditoSeleccionado.cuotas_pagadas || 0 }}/{{ creditoSeleccionado.numero_cuotas || 0 }}</p>
          </div>
        </div>

        <div class="flex justify-end">
          <Button label="Imprimir cuotas" icon="pi pi-print" size="small" severity="secondary" @click="imprimirCuotas" />
        </div>

        <DataTable :value="cuotasDetalle" :loading="loadingCuotas" size="small" stripedRows tableStyle="min-width: 86rem">
          <template #empty>
            <p class="p-4 text-center">No hay cuotas para este credito.</p>
          </template>
          <Column field="numero" header="Nro." sortable />
          <Column field="fecha_vencimiento" header="Vencimiento" sortable>
            <template #body="slotProps">{{ formatFecha(slotProps.data.fecha_vencimiento) }}</template>
          </Column>
          <Column field="monto" header="Monto" sortable>
            <template #body="slotProps">$ {{ formatPrecio(slotProps.data.monto) }}</template>
          </Column>
          <Column field="monto_pagado" header="Pagado/Falta" sortable>
            <template #body="slotProps">
              <div>
                <p>$ {{ formatPrecio(slotProps.data.monto_pagado) }}</p>
                <p class="text-xs text-gray-500">Falta: $ {{ formatPrecio(saldoCuota(slotProps.data)) }}</p>
              </div>
            </template>
          </Column>
          <Column field="monto_bob_pagado" header="Pagado Bs" sortable>
            <template #body="slotProps">
              {{ slotProps.data.monto_bob_pagado ? 'Bs ' + formatPrecio(slotProps.data.monto_bob_pagado) : '-' }}
            </template>
          </Column>
          <Column field="estado" header="Estado" sortable>
            <template #body="slotProps">
              <Tag :value="labelEstadoCuota(slotProps.data.estado)" :severity="severityCuota(slotProps.data.estado)" />
            </template>
          </Column>
          <Column field="fecha_pago" header="Fecha pago" sortable>
            <template #body="slotProps">{{ slotProps.data.fecha_pago ? formatFecha(slotProps.data.fecha_pago) : '-' }}</template>
          </Column>
          <Column field="usuario_pago" header="Pago aceptado por" sortable>
            <template #body="slotProps">{{ slotProps.data.usuario_pago || '-' }}</template>
          </Column>
          <Column header="Historial">
            <template #body="slotProps">
              <div v-if="historialCuota(slotProps.data).length" class="space-y-1 text-xs text-gray-600">
                <p v-for="pago in historialCuota(slotProps.data)" :key="pago.id">
                  $ {{ formatPrecio(pago.monto_usd) }} / Bs {{ formatPrecio(pago.monto_bob) }}
                </p>
              </div>
              <span v-else class="text-sm text-gray-500">-</span>
            </template>
          </Column>
          <Column header="Accion">
            <template #body="slotProps">
              <Button
                v-if="slotProps.data.estado !== 'pagada' && !soloLectura"
                label="Pagar"
                icon="pi pi-check"
                size="small"
                severity="success"
                :loading="pagandoCuotaId === slotProps.data.id"
                @click="abrirPagoCuota(slotProps.data)"
              />
              <span v-else class="text-sm text-gray-500">{{ slotProps.data.estado === 'pagada' ? 'Registrada' : 'Solo lectura' }}</span>
            </template>
          </Column>
        </DataTable>
      </div>
    </Dialog>

    <Dialog v-model:visible="pagoVisible" modal header="Registrar pago de cuota" :style="{ width: '30rem', maxWidth: '95vw' }">
      <div v-if="cuotaPago" class="flex flex-col gap-4">
        <div class="grid grid-cols-2 gap-3 text-sm">
          <div class="rounded-lg border border-gray-200 p-3">
            <p class="text-xs uppercase text-gray-500">Cuota</p>
            <p class="font-semibold text-gray-900">Nro. {{ cuotaPago.numero }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 p-3">
            <p class="text-xs uppercase text-gray-500">Monto USD</p>
            <p class="font-semibold text-gray-900">$ {{ formatPrecio(cuotaPago.monto) }}</p>
            <p class="text-xs text-gray-500">Falta: $ {{ formatPrecio(saldoCuota(cuotaPago)) }}</p>
          </div>
        </div>

        <div class="flex flex-col gap-1">
          <label for="monto_pago">Monto a pagar</label>
          <InputNumber
            id="monto_pago"
            v-model="montoPago"
            :min="0"
            :max="saldoAplicablePago"
            :minFractionDigits="2"
            :maxFractionDigits="2"
            fluid
            size="small"
            placeholder="Monto en USD"
          />
          <p class="text-xs text-gray-500">Saldo disponible desde esta cuota: $ {{ formatPrecio(saldoAplicablePago) }}</p>
        </div>

        <div class="flex flex-col gap-1">
          <label for="tipo_cambio_pago">Precio del dolar actual</label>
          <InputNumber
            id="tipo_cambio_pago"
            v-model="tipoCambioPago"
            :min="0"
            :minFractionDigits="2"
            :maxFractionDigits="4"
            fluid
            size="small"
            placeholder="Ej. 6.96"
          />
        </div>

        <div class="rounded-lg border border-gray-200 bg-gray-50 p-3">
          <p class="text-sm text-gray-500">Total a cobrar en bolivianos</p>
          <p class="text-2xl font-bold text-gray-900">Bs {{ formatPrecio(montoPagoBOB) }}</p>
        </div>

        <div class="flex justify-end gap-2 border-t border-gray-100 pt-4">
          <Button label="Cancelar" severity="secondary" @click="cerrarPagoCuota" />
          <Button label="Registrar pago" icon="pi pi-check" severity="success" :loading="pagandoCuotaId === cuotaPago.id" @click="pagarCuota" />
        </div>
      </div>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { server } from '~/server/server';
import Button from 'primevue/button';
import Column from 'primevue/column';
import DataTable from 'primevue/datatable';
import Dialog from 'primevue/dialog';
import InputNumber from 'primevue/inputnumber';
import InputText from 'primevue/inputtext';
import SelectButton from 'primevue/selectbutton';
import Tag from 'primevue/tag';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';

const props = defineProps({
  scope: {
    type: String,
    default: 'personal'
  }
});

const toast = useToast();
const creditos = ref<any[]>([]);
const cuotasPorVenta = ref<Record<number, any[]>>({});
const cuotasDetalle = ref<any[]>([]);
const creditoSeleccionado = ref<any | null>(null);
const detalleVisible = ref(false);
const pagoVisible = ref(false);
const loading = ref(true);
const loadingCuotas = ref(false);
const pagandoCuotaId = ref<number | null>(null);
const cuotaPago = ref<any | null>(null);
const tipoCambioPago = ref<number | null>(null);
const montoPago = ref<number | null>(null);
const searchQuery = ref('');
const tipoFiltro = ref('todos');
const estadoFiltro = ref('todos');
const rolUsuario = ref('');

const mostrarVendedor = computed(() => props.scope === 'general');
const soloLectura = computed(() => rolUsuario.value === 'contador');
const descripcion = computed(() => mostrarVendedor.value
  ? 'Historial general de credito directo y credito bancario.'
  : 'Seguimiento de credito directo y credito bancario del vendedor actual.'
);

const tipoFiltros = [
  { label: 'Todos', value: 'todos' },
  { label: 'Directo', value: 'credito_directo' },
  { label: 'Bancario', value: 'credito_bancario' }
];

const estadoFiltros = [
  { label: 'Todos', value: 'todos' },
  { label: 'Activos', value: 'activos' },
  { label: 'Atrasados', value: 'atrasados' },
  { label: 'Pagados', value: 'pagados' }
];

const filteredCreditos = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();

  return creditos.value.filter((credito: any) => {
    const coincideTipo = tipoFiltro.value === 'todos' || credito.tipo_venta === tipoFiltro.value;
    const coincideEstado =
      estadoFiltro.value === 'todos' ||
      (estadoFiltro.value === 'activos' && credito.estado_venta === 'en_credito') ||
      (estadoFiltro.value === 'atrasados' && credito.tiene_atraso) ||
      (estadoFiltro.value === 'pagados' && credito.estado_venta === 'pagado_completo');
    const coincideBusqueda = !query ||
      (credito.id?.toString() || '').includes(query) ||
      (credito.fecha?.toLowerCase() || '').includes(query) ||
      (credito.fecha_venta?.toLowerCase() || '').includes(query) ||
      (credito.vendedor?.toLowerCase() || '').includes(query) ||
      (credito.cliente?.toLowerCase() || '').includes(query) ||
      (credito.ci_cliente?.toLowerCase() || '').includes(query) ||
      (credito.vehiculo?.toLowerCase() || '').includes(query) ||
      (credito.categoria?.toLowerCase() || '').includes(query) ||
      (credito.tipo_venta?.toLowerCase() || '').includes(query) ||
      (credito.estado_venta?.toLowerCase() || '').includes(query);

    return coincideTipo && coincideEstado && coincideBusqueda;
  });
});

const totalFinanciado = computed(() => {
  return filteredCreditos.value.reduce((total: number, credito: any) => total + Number(credito.monto_financiado || credito.saldo || 0), 0);
});

const totalSaldoPendiente = computed(() => {
  return filteredCreditos.value.reduce((total: number, credito: any) => total + Number(credito.saldo || 0), 0);
});

const totalCuotasAtrasadas = computed(() => {
  return filteredCreditos.value.reduce((total: number, credito: any) => total + Number(credito.cuotas_atrasadas || 0), 0);
});

const montoPagoBOB = computed(() => {
  return Number(montoPago.value || 0) * Number(tipoCambioPago.value || 0);
});

const saldoAplicablePago = computed(() => {
  if (!cuotaPago.value) return 0;
  return cuotasDetalle.value
    .filter((cuota: any) => cuota.numero >= cuotaPago.value.numero && cuota.estado !== 'pagada')
    .reduce((total: number, cuota: any) => total + saldoCuota(cuota), 0);
});

onMounted(async () => {
  rolUsuario.value = obtenerRolUsuarioActual();
  await obtenerCreditos();
});

async function obtenerCreditos() {
  loading.value = true;

  try {
    const query: Record<string, any> = {};
    if (props.scope !== 'general') {
      const user = localStorage.getItem('user');
      const parsedUser = user ? JSON.parse(user) : null;
      if (parsedUser?.id) {
        query.id_usuario = parsedUser.id;
      }
    }

    const res = await $fetch(server.HOST + '/api/v1/ventas', {
      method: 'GET',
      query
    });

    const ventasCredito = (Array.isArray(res) ? res : []).filter((venta: any) =>
      ['credito_directo', 'credito_bancario'].includes(venta.tipo_venta)
    );

    creditos.value = await Promise.all(ventasCredito.map(async (venta: any) => {
      const cuotas = await obtenerCuotasVenta(venta.id);
      return enriquecerCredito(venta, cuotas);
    }));
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar creditos', life: 3000 });
  } finally {
    loading.value = false;
  }
}

async function obtenerCuotasVenta(idVenta: number) {
  if (cuotasPorVenta.value[idVenta]) {
    return cuotasPorVenta.value[idVenta];
  }

  const res = await $fetch(server.HOST + `/api/v1/ventas/${idVenta}/cuotas`, {
    method: 'GET'
  });
  const cuotas = Array.isArray(res) ? res : [];
  cuotasPorVenta.value[idVenta] = cuotas;
  return cuotas;
}

function enriquecerCredito(venta: any, cuotas: any[]) {
  const pendientes = cuotas.filter((cuota: any) => cuota.estado !== 'pagada');
  const atrasadas = cuotas.filter((cuota: any) => cuota.estado === 'atrasada');
  const proximaCuota = pendientes[0];
  const saldoPendiente = cuotas.reduce((total: number, cuota: any) => total + saldoCuota(cuota), 0);

  return {
    ...venta,
    saldo: saldoPendiente,
    cuotas_pagadas: cuotas.filter((cuota: any) => cuota.estado === 'pagada').length,
    cuotas_pendientes: pendientes.length,
    cuotas_atrasadas: atrasadas.length,
    tiene_atraso: atrasadas.length > 0,
    proxima_cuota_fecha: proximaCuota?.fecha_vencimiento || '',
    proxima_cuota_monto: proximaCuota ? saldoCuota(proximaCuota) : 0
  };
}

async function abrirDetalle(credito: any) {
  creditoSeleccionado.value = credito;
  detalleVisible.value = true;
  loadingCuotas.value = true;

  try {
    cuotasPorVenta.value[credito.id] = [];
    const res = await $fetch(server.HOST + `/api/v1/ventas/${credito.id}/cuotas`, {
      method: 'GET'
    });
    cuotasDetalle.value = Array.isArray(res) ? res : [];
    cuotasPorVenta.value[credito.id] = cuotasDetalle.value;
    actualizarCreditoEnLista(credito.id, cuotasDetalle.value);
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar cuotas', life: 3000 });
  } finally {
    loadingCuotas.value = false;
  }
}

function abrirPagoCuota(cuota: any) {
  cuotaPago.value = cuota;
  montoPago.value = saldoCuota(cuota);
  tipoCambioPago.value = null;
  pagoVisible.value = true;
}

function cerrarPagoCuota() {
  pagoVisible.value = false;
  cuotaPago.value = null;
  tipoCambioPago.value = null;
  montoPago.value = null;
}

async function pagarCuota() {
  if (!cuotaPago.value) return;
  const monto = Number(montoPago.value || 0);
  if (monto <= 0) {
    toast.add({ severity: 'warn', summary: 'Ingrese el monto a pagar', life: 3000 });
    return;
  }
  if (monto > saldoAplicablePago.value) {
    toast.add({ severity: 'warn', summary: 'Monto mayor al saldo pendiente del credito', life: 3000 });
    return;
  }
  if (!tipoCambioPago.value || Number(tipoCambioPago.value) <= 0) {
    toast.add({ severity: 'warn', summary: 'Ingrese el precio del dolar actual', life: 3000 });
    return;
  }

  pagandoCuotaId.value = cuotaPago.value.id;
  try {
    const userId = obtenerUsuarioActualId();
    if (!userId) {
      toast.add({ severity: 'warn', summary: 'No se encontro usuario en sesion', life: 3000 });
      return;
    }
    await $fetch(server.HOST + `/api/v1/cuotas-credito/${cuotaPago.value.id}/pagar`, {
      method: 'PATCH',
      body: {
        monto_pago: monto,
        tipo_cambio_pago: Number(tipoCambioPago.value || 0),
        id_usuario_pago: userId
      }
    });
    toast.add({ severity: 'success', summary: 'Pago registrado', detail: `Cobrado: Bs ${formatPrecio(montoPagoBOB.value)}`, life: 3000 });
    cerrarPagoCuota();
    if (creditoSeleccionado.value) {
      await abrirDetalle(creditoSeleccionado.value);
    }
  } catch (err: any) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al pagar cuota', detail: err?.data || err?.message, life: 4000 });
  } finally {
    pagandoCuotaId.value = null;
  }
}

function actualizarCreditoEnLista(idVenta: number, cuotas: any[]) {
  const index = creditos.value.findIndex((credito: any) => credito.id === idVenta);
  if (index === -1) {
    return;
  }

  const actualizado = enriquecerCredito(creditos.value[index], cuotas);
  creditos.value.splice(index, 1, actualizado);
  creditoSeleccionado.value = actualizado;
}

function obtenerUsuarioActualId() {
  const user = localStorage.getItem('user');
  return user ? Number(JSON.parse(user)?.id || 0) : 0;
}

function obtenerRolUsuarioActual() {
  const user = localStorage.getItem('user');
  return user ? String(JSON.parse(user)?.rol || '') : '';
}

function saldoCuota(cuota: any) {
  if (!cuota) return 0;
  const saldo = Number(cuota.saldo_pendiente ?? 0);
  if (saldo > 0 || cuota.estado === 'pagada') {
    return saldo;
  }
  return Math.max(Number(cuota.monto || 0) - Number(cuota.monto_pagado || 0), 0);
}

function historialCuota(cuota: any) {
  if (!cuota?.historial_pagos) return [];
  if (Array.isArray(cuota.historial_pagos)) return cuota.historial_pagos;
  try {
    return JSON.parse(cuota.historial_pagos);
  } catch {
    return [];
  }
}

function imprimirCuotas() {
  if (!creditoSeleccionado.value) {
    return;
  }

  const filas = cuotasDetalle.value.map((cuota: any) => `
    <tr>
      <td>${cuota.numero}</td>
      <td>${formatFecha(cuota.fecha_vencimiento)}</td>
      <td>$ ${formatPrecio(cuota.monto)}</td>
      <td>$ ${formatPrecio(cuota.monto_pagado)}</td>
      <td>$ ${formatPrecio(saldoCuota(cuota))}</td>
      <td>${cuota.monto_bob_pagado ? 'Bs ' + formatPrecio(cuota.monto_bob_pagado) : '-'}</td>
      <td>${labelEstadoCuota(cuota.estado)}</td>
      <td>${cuota.fecha_pago ? formatFecha(cuota.fecha_pago) : '-'}</td>
      <td>${escapeHtml(cuota.usuario_pago || '-')}</td>
    </tr>
  `).join('');

  const html = `
    <!doctype html>
    <html>
      <head>
        <title>Cuotas credito ${creditoSeleccionado.value.id}</title>
        <style>
          body { font-family: Arial, sans-serif; color: #111827; margin: 24px; }
          h1 { font-size: 20px; margin: 0 0 8px; }
          p { margin: 2px 0; font-size: 13px; }
          table { width: 100%; border-collapse: collapse; margin-top: 18px; font-size: 12px; }
          th, td { border: 1px solid #d1d5db; padding: 8px; text-align: left; }
          th { background: #f3f4f6; }
          .summary { display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; margin-top: 12px; }
          .box { border: 1px solid #d1d5db; padding: 8px; }
        </style>
      </head>
      <body>
        <h1>Plan de cuotas de credito</h1>
        <p><strong>Credito:</strong> ${creditoSeleccionado.value.id}</p>
        <p><strong>Cliente:</strong> ${escapeHtml(creditoSeleccionado.value.cliente || 'Sin nombre')}</p>
        <p><strong>Vehiculo:</strong> ${escapeHtml(creditoSeleccionado.value.vehiculo || '-')}</p>
        <div class="summary">
          <div class="box"><strong>Total:</strong><br>$ ${formatPrecio(creditoSeleccionado.value.precio_total)}</div>
          <div class="box"><strong>Saldo:</strong><br>$ ${formatPrecio(creditoSeleccionado.value.saldo)}</div>
          <div class="box"><strong>Progreso:</strong><br>${creditoSeleccionado.value.cuotas_pagadas || 0}/${creditoSeleccionado.value.numero_cuotas || 0}</div>
        </div>
        <table>
          <thead>
            <tr>
              <th>Nro.</th>
              <th>Vencimiento</th>
              <th>Monto USD</th>
              <th>Pagado USD</th>
              <th>Falta USD</th>
              <th>Pagado Bs</th>
              <th>Estado</th>
              <th>Fecha pago</th>
              <th>Pago aceptado por</th>
            </tr>
          </thead>
          <tbody>${filas}</tbody>
        </table>
      </body>
    </html>
  `;

  const printWindow = window.open('', '_blank', 'width=900,height=700');
  if (!printWindow) {
    toast.add({ severity: 'warn', summary: 'No se pudo abrir la ventana de impresion', life: 3000 });
    return;
  }
  printWindow.document.write(html);
  printWindow.document.close();
  printWindow.focus();
  printWindow.print();
}

function formatPrecio(precio: number) {
  return Number(precio || 0).toLocaleString('es-BO', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  });
}

function escapeHtml(value: unknown) {
  return String(value)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;');
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

function labelTipoCredito(tipo: string) {
  if (tipo === 'credito_directo') return 'Credito directo';
  if (tipo === 'credito_bancario') return 'Credito bancario';
  return tipo || 'Credito';
}

function labelEstadoCredito(estado: string) {
  if (estado === 'en_credito') return 'En credito';
  if (estado === 'pagado_completo') return 'Pagado completo';
  return estado || 'Pendiente';
}

function labelEstadoCuota(estado: string) {
  if (estado === 'pagada') return 'Pagada';
  if (estado === 'abonada') return 'Abonada';
  if (estado === 'atrasada') return 'Atrasada';
  return 'Pendiente';
}

function severityEstado(estado: string) {
  if (estado === 'pagado_completo') return 'success';
  if (estado === 'Anulada') return 'danger';
  return 'warning';
}

function severityCuota(estado: string) {
  if (estado === 'pagada') return 'success';
  if (estado === 'abonada') return 'info';
  if (estado === 'atrasada') return 'danger';
  return 'warning';
}
</script>
