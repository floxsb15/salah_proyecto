<template>
  <Toast />

  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <h2 class="text-2xl font-bold">Historial de Pedidos</h2>
        <p class="text-sm text-gray-500">Vehiculos solicitados para importacion con adelanto registrado.</p>
      </div>

      <div class="flex flex-col gap-2 sm:flex-row sm:items-center">
        <Button label="Nuevo pedido" icon="pi pi-plus" size="small" @click="abrirNuevoPedido" />
        <span class="p-input-icon-left">
          <i class="pi pi-search" />
          <InputText v-model="searchQuery" placeholder="Buscar..." size="small" />
        </span>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-3 md:grid-cols-4">
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Pedidos</p>
        <p class="text-2xl font-bold text-gray-900">{{ filteredPedidos.length }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">En aduana</p>
        <p class="text-2xl font-bold text-gray-900">{{ pedidosEnTransito }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Adelantos</p>
        <p class="text-2xl font-bold text-gray-900">$ {{ formatPrecio(totalAdelantos) }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Saldo pendiente</p>
        <p class="text-2xl font-bold text-gray-900">$ {{ formatPrecio(totalPendiente) }}</p>
      </div>
    </div>

    <DataTable
      :value="filteredPedidos"
      :loading="loading"
      tableStyle="min-width: 86rem"
      size="small"
      stripedRows
      removableSort
      paginator
      :rows="10"
      :rowsPerPageOptions="[10, 20, 50]"
    >
      <template #empty>
        <p class="p-4 text-center">No hay pedidos registrados.</p>
      </template>

      <Column field="id" header="ID" sortable />
      <Column field="fecha" header="Fecha" sortable>
        <template #body="slotProps">{{ formatFecha(slotProps.data.fecha) }}</template>
      </Column>
      <Column field="cliente" header="Cliente" sortable>
        <template #body="slotProps">
          <div>
            <p class="font-medium text-gray-900">{{ slotProps.data.cliente || 'Sin nombre' }}</p>
            <p class="text-xs text-gray-500">CI/NIT: {{ slotProps.data.ci_cliente || 'N/A' }}</p>
          </div>
        </template>
      </Column>
      <Column field="marca" header="Vehiculo solicitado" sortable>
        <template #body="slotProps">
          <div>
            <p class="font-medium text-gray-900">{{ vehiculoSolicitado(slotProps.data) }}</p>
            <p class="text-xs text-gray-500">{{ slotProps.data.color || 'Sin color' }}<span v-if="slotProps.data.version"> / {{ slotProps.data.version }}</span></p>
          </div>
        </template>
      </Column>
      <Column field="pais_origen" header="Pais de origen" sortable />
      <Column field="precio_estimado_usd" header="Precio estimado" sortable>
        <template #body="slotProps">$ {{ formatPrecio(slotProps.data.precio_estimado_usd) }}</template>
      </Column>
      <Column field="adelanto_pagado_usd" header="Adelanto pagado" sortable>
        <template #body="slotProps">
          <div>
            <p class="font-medium text-gray-900">$ {{ formatPrecio(totalPagadoUSD(slotProps.data)) }}</p>
            <p class="text-xs text-gray-500">TC {{ formatTipoCambio(slotProps.data.tipo_cambio_usado) }}</p>
          </div>
        </template>
      </Column>
      <Column field="saldo_pendiente_usd" header="Saldo pendiente" sortable>
        <template #body="slotProps">$ {{ formatPrecio(slotProps.data.saldo_pendiente_usd) }}</template>
      </Column>
      <Column field="estado" header="Estado" sortable>
        <template #body="slotProps">
          <div class="flex items-center gap-2">
            <Tag :value="estadoPedidoLabel(slotProps.data.estado)" :severity="estadoPedidoSeverity(slotProps.data.estado)" />
            <i
              v-if="vehiculoNoDisponible(slotProps.data)"
              class="pi pi-lock text-amber-600"
              title="Vehiculo no disponible para venta hasta completar el pedido"
            ></i>
          </div>
        </template>
      </Column>
      <Column field="fecha_vencimiento" header="Vence" sortable>
        <template #body="slotProps">{{ formatFecha(slotProps.data.fecha_vencimiento) }}</template>
      </Column>
      <Column header="Acciones" style="width: 15rem">
        <template #body="slotProps">
          <div class="flex items-center gap-1">
            <Button
              v-if="normalizarEstadoUI(slotProps.data.estado) === 'Pedido registrado'"
              label="Aduana"
              icon="pi pi-send"
              size="small"
              severity="info"
              variant="text"
              :loading="changingStateId === slotProps.data.id"
              @click="marcarEnTransito(slotProps.data)"
            />
            <Button
              v-if="esEstadoTransito(slotProps.data.estado)"
              label="Recibir"
              icon="pi pi-box"
              size="small"
              severity="secondary"
              variant="text"
              @click="router.push({ path: recibirPath, query: { id: slotProps.data.id } })"
            />
            <Button
              v-if="esEstadoRecibido(slotProps.data.estado) && Number(slotProps.data.saldo_pendiente_usd || 0) > 0"
              label="Completar"
              icon="pi pi-check"
              size="small"
              severity="success"
              variant="text"
              @click="router.push({ path: completarPath, query: { id: slotProps.data.id } })"
            />
          </div>
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="pedidoVisible" modal header="Nuevo pedido" :style="{ width: '68rem', maxWidth: '96vw' }">
      <form class="grid max-h-[78vh] grid-cols-1 gap-4 overflow-y-auto pr-1 lg:grid-cols-2" @submit.prevent="registrarPedido">
        <div class="flex flex-col gap-1 lg:col-span-2">
          <label for="cliente">Cliente</label>
          <div class="grid grid-cols-1 gap-2 md:grid-cols-[1fr_auto]">
            <Select id="cliente" v-model="form.id_cliente" :options="clientesActivos" option-label="nombreCompleto" option-value="id" placeholder="Seleccione cliente" filter fluid size="small" />
            <Button label="Crear cliente" icon="pi pi-user-plus" severity="secondary" type="button" size="small" @click="clienteVisible = true" />
          </div>
        </div>

        <section class="grid grid-cols-1 gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2 md:grid-cols-3">
          <div class="flex flex-col gap-1">
            <label for="marca">Marca</label>
            <InputText id="marca" v-model="form.marca" placeholder="Ej. Toyota" fluid size="small" />
          </div>
          <div class="flex flex-col gap-1">
            <label for="modelo">Modelo</label>
            <InputText id="modelo" v-model="form.modelo" placeholder="Ej. Land Cruiser" fluid size="small" />
          </div>
          <div class="flex flex-col gap-1">
            <label for="anio">Anio</label>
            <InputNumber id="anio" v-model="form.anio" :min="1900" :useGrouping="false" fluid size="small" />
          </div>
          <div class="flex flex-col gap-1">
            <label for="color">Color</label>
            <InputText id="color" v-model="form.color" placeholder="Ej. Blanco" fluid size="small" />
          </div>
          <div class="flex flex-col gap-1">
            <label for="pais_origen">Pais de origen</label>
            <InputText id="pais_origen" v-model="form.pais_origen" placeholder="Ej. Japon" fluid size="small" />
          </div>
          <div class="flex flex-col gap-1">
            <label for="fecha_llegada_estimada">Llegada estimada</label>
            <InputText id="fecha_llegada_estimada" v-model="form.fecha_llegada_estimada" type="date" fluid size="small" />
          </div>
          <div class="flex flex-col gap-1 md:col-span-3">
            <label for="version">Version / especificaciones</label>
            <Textarea id="version" v-model="form.version" rows="2" auto-resize fluid />
          </div>
        </section>

        <div class="flex flex-col gap-1 lg:col-span-2">
          <label for="proveedor">Proveedor</label>
          <div class="grid grid-cols-1 gap-2 md:grid-cols-[1fr_auto]">
            <Select id="proveedor" v-model="form.id_proveedor" :options="proveedoresActivos" option-label="nombreCompleto" option-value="id" placeholder="Opcional" show-clear filter fluid size="small" />
            <Button label="Crear proveedor" icon="pi pi-briefcase" severity="secondary" type="button" size="small" @click="proveedorVisible = true" />
          </div>
        </div>

        <div class="flex flex-col gap-1">
          <label for="fecha">Fecha</label>
          <InputText id="fecha" v-model="form.fecha" type="date" fluid size="small" />
        </div>
        <div class="flex flex-col gap-1">
          <label for="tipo_cambio">Tipo de cambio del dia</label>
          <InputNumber id="tipo_cambio" v-model="form.tipo_cambio" :min="0" :minFractionDigits="2" :maxFractionDigits="4" suffix=" Bs/USD" fluid size="small" />
        </div>
        <div class="flex flex-col gap-1">
          <label for="precio_estimado_usd">Precio estimado USD</label>
          <InputNumber id="precio_estimado_usd" v-model="form.precio_estimado_usd" mode="currency" currency="USD" locale="es-BO" :min="0" fluid size="small" />
        </div>
        <div class="flex flex-col gap-1">
          <label for="adelanto_modo">Tipo de adelanto</label>
          <Select id="adelanto_modo" v-model="form.adelanto_modo" :options="modosAdelanto" option-label="label" option-value="value" fluid size="small" />
        </div>
        <div class="flex flex-col gap-1">
          <label for="adelanto_porcentaje">Adelanto %</label>
          <InputNumber id="adelanto_porcentaje" v-model="form.adelanto_porcentaje" :min="0" :max="100" suffix=" %" :disabled="form.adelanto_modo !== 'porcentaje'" fluid size="small" />
        </div>
        <div class="flex flex-col gap-1">
          <label for="adelanto_requerido_usd">Adelanto requerido USD</label>
          <InputNumber id="adelanto_requerido_usd" v-model="form.adelanto_requerido_usd" mode="currency" currency="USD" locale="es-BO" :min="0" :disabled="form.adelanto_modo === 'porcentaje'" fluid size="small" />
        </div>

        <section class="flex flex-col gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2">
          <div class="flex items-center justify-between gap-3">
            <div>
              <h3 class="text-sm font-semibold text-gray-900">Detalle de pago del adelanto</h3>
              <p class="text-xs text-gray-500">TC {{ formatTipoCambio(tipoCambio) }} / Pagado: $ {{ formatPrecio(pagoEquivalenteUSD) }} / Bs {{ formatPrecio(pagoEquivalenteBOB) }}</p>
            </div>
            <Button label="Agregar pago" icon="pi pi-plus" size="small" severity="secondary" type="button" @click="agregarPago" />
          </div>

          <div v-for="(pago, index) in form.pagos" :key="pago.key" class="grid grid-cols-1 gap-2 md:grid-cols-[9rem_13rem_1fr_2.5rem]">
            <Select v-model="pago.moneda" :options="monedasPago" placeholder="Moneda" fluid size="small" />
            <Select v-model="pago.metodo" :options="metodosPago" placeholder="Tipo de pago" fluid size="small" />
            <InputNumber v-model="pago.monto" mode="decimal" :min="0" :minFractionDigits="2" :maxFractionDigits="2" placeholder="Monto" fluid size="small" />
            <Button icon="pi pi-trash" severity="danger" text rounded type="button" aria-label="Eliminar pago" @click="eliminarPago(index)" />
          </div>

          <p v-if="pagoMenorAdelanto" class="text-sm font-semibold text-red-600">El pago debe ser igual o mayor al adelanto requerido.</p>
          <p v-else-if="pagoExcedeAdelanto" class="text-sm font-semibold text-amber-600">El pago supera el adelanto requerido.</p>
        </section>

        <div class="grid grid-cols-1 gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2 md:grid-cols-4">
          <div>
            <span class="text-xs font-semibold uppercase text-gray-500">Precio estimado</span>
            <strong class="block text-gray-900">$ {{ formatPrecio(form.precio_estimado_usd) }}</strong>
            <span class="block text-xs text-gray-500">Bs {{ formatPrecio(precioEstimadoBOB) }}</span>
          </div>
          <div>
            <span class="text-xs font-semibold uppercase text-gray-500">Adelanto requerido</span>
            <strong class="block text-gray-900">$ {{ formatPrecio(adelantoRequeridoUSD) }}</strong>
            <span class="block text-xs text-gray-500">Bs {{ formatPrecio(adelantoRequeridoBOB) }}</span>
          </div>
          <div>
            <span class="text-xs font-semibold uppercase text-gray-500">Adelanto Bs</span>
            <strong class="block text-gray-900">Bs {{ formatPrecio(adelantoRequeridoBOB) }}</strong>
            <span class="block text-xs text-gray-500">$ {{ formatPrecio(adelantoRequeridoUSD) }}</span>
          </div>
          <div>
            <span class="text-xs font-semibold uppercase text-gray-500">Saldo estimado</span>
            <strong class="block text-gray-900">$ {{ formatPrecio(saldoEstimadoUSD) }}</strong>
            <span class="block text-xs text-gray-500">Bs {{ formatPrecio(saldoEstimadoBOB) }}</span>
          </div>
        </div>

        <div class="flex flex-col gap-1 lg:col-span-2">
          <label for="observacion">Observacion</label>
          <Textarea id="observacion" v-model="form.observacion" rows="3" auto-resize fluid />
        </div>

        <div class="flex justify-end gap-2 border-t border-gray-100 pt-4 lg:col-span-2">
          <Button label="Cancelar" severity="secondary" type="button" @click="pedidoVisible = false" />
          <Button label="Registrar pedido" icon="pi pi-check" type="submit" :loading="saving" :disabled="pagoMenorAdelanto" />
        </div>
      </form>
    </Dialog>

    <modalAgregarCliente
      v-if="clienteVisible"
      :open="clienteVisible"
      @close="clienteVisible = false"
      @update="obtenerClientes"
      @success="toast.add({ severity: 'success', summary: 'Cliente agregado', life: 3000 })"
      @error="mostrarError('Error al agregar cliente', $event)"
    />

    <Dialog v-model:visible="proveedorVisible" modal header="Nuevo proveedor" :style="{ width: '34rem', maxWidth: '96vw' }">
      <form class="grid grid-cols-1 gap-3" @submit.prevent="guardarProveedor">
        <div class="flex flex-col gap-1">
          <label for="proveedor_nombre">Nombre / Razon social</label>
          <InputText id="proveedor_nombre" v-model="proveedorForm.nombre" fluid size="small" />
        </div>
        <div class="flex flex-col gap-1">
          <label for="proveedor_ci">CI/NIT</label>
          <InputText id="proveedor_ci" v-model="proveedorForm.ci_nit" fluid size="small" />
        </div>
        <div class="flex flex-col gap-1">
          <label for="proveedor_telefono">Telefono</label>
          <InputText id="proveedor_telefono" v-model="proveedorForm.telefono" fluid size="small" />
        </div>
        <div class="flex justify-end gap-2 border-t border-gray-100 pt-3">
          <Button label="Cancelar" severity="secondary" type="button" @click="proveedorVisible = false" />
          <Button label="Crear proveedor" icon="pi pi-check" type="submit" :loading="savingProveedor" />
        </div>
      </form>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import { server } from '~/server/server';
import Button from 'primevue/button';
import Column from 'primevue/column';
import DataTable from 'primevue/datatable';
import Dialog from 'primevue/dialog';
import InputNumber from 'primevue/inputnumber';
import InputText from 'primevue/inputtext';
import Select from 'primevue/select';
import Tag from 'primevue/tag';
import Textarea from 'primevue/textarea';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';
import modalAgregarCliente from '~/components/admin/clientes/modalAgregarCliente.vue';

const props = withDefaults(defineProps<{
  general?: boolean;
  recibirPath?: string;
  completarPath?: string;
}>(), {
  general: false,
  recibirPath: '/ventas/recibir-pedido',
  completarPath: '/ventas/completar-pedido'
});

const router = useRouter();
const toast = useToast();
const pedidos = ref<any[]>([]);
const clientes = ref<any[]>([]);
const proveedores = ref<any[]>([]);
const loading = ref(true);
const saving = ref(false);
const savingProveedor = ref(false);
const pedidoVisible = ref(false);
const clienteVisible = ref(false);
const proveedorVisible = ref(false);
const searchQuery = ref('');
const changingStateId = ref<number | null>(null);
const monedasPago = ref(['USD', 'BOB']);
const metodosPago = ref(['Efectivo', 'QR', 'Transferencia', 'Tarjeta']);
const modosAdelanto = ref([
  { label: 'Porcentaje del precio', value: 'porcentaje' },
  { label: 'Monto manual USD', value: 'manual' }
]);

const form = reactive({
  id_cliente: null as number | null,
  id_proveedor: null as number | null,
  fecha: new Date().toISOString().slice(0, 10),
  marca: '',
  modelo: '',
  anio: new Date().getFullYear(),
  color: '',
  version: '',
  pais_origen: '',
  precio_estimado_usd: 0,
  tipo_cambio: null as number | null,
  fecha_llegada_estimada: '',
  adelanto_modo: 'porcentaje',
  adelanto_requerido_usd: 0,
  adelanto_porcentaje: 10,
  pagos: [crearPago()],
  observacion: ''
});

const proveedorForm = reactive({
  nombre: '',
  ci_nit: '',
  telefono: ''
});

const clientesActivos = computed(() => clientes.value
  .filter((cliente: any) => cliente.estado === 'Activo')
  .map((cliente: any) => ({
    ...cliente,
    nombreCompleto: `${cliente.nombre || ''} ${cliente.apellido || ''} - ${cliente.ci || 'Sin CI'}`.trim()
  })));
const proveedoresActivos = computed(() => proveedores.value
  .filter((proveedor: any) => proveedor.estado === 'Activo')
  .map((proveedor: any) => ({
    ...proveedor,
    nombreCompleto: `${proveedor.nombre || ''} - ${proveedor.ci_nit || 'Sin CI/NIT'}`.trim()
  })));
const filteredPedidos = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  if (!query) return pedidos.value;
  return pedidos.value.filter((pedido: any) =>
    (pedido.id?.toString() || '').includes(query) ||
    (pedido.cliente?.toLowerCase() || '').includes(query) ||
    (pedido.marca?.toLowerCase() || '').includes(query) ||
    (pedido.modelo?.toLowerCase() || '').includes(query) ||
    (pedido.pais_origen?.toLowerCase() || '').includes(query) ||
    (pedido.estado?.toLowerCase() || '').includes(query)
  );
});
const tipoCambio = computed(() => Number(form.tipo_cambio || 0));
const precioEstimadoBOB = computed(() => roundMoney(Number(form.precio_estimado_usd || 0) * tipoCambio.value));
const adelantoRequeridoUSD = computed(() => {
  const porcentaje = Number(form.adelanto_porcentaje || 0);
  if (form.adelanto_modo === 'porcentaje' && porcentaje > 0) {
    return roundMoney(Number(form.precio_estimado_usd || 0) * porcentaje / 100);
  }
  return Number(form.adelanto_requerido_usd || 0);
});
const adelantoRequeridoBOB = computed(() => roundMoney(adelantoRequeridoUSD.value * tipoCambio.value));
const pagoUSDDirecto = computed(() => form.pagos.filter((pago: any) => pago.moneda === 'USD').reduce((total: number, pago: any) => total + Number(pago.monto || 0), 0));
const pagoBOBDirecto = computed(() => form.pagos.filter((pago: any) => pago.moneda === 'BOB').reduce((total: number, pago: any) => total + Number(pago.monto || 0), 0));
const pagoEquivalenteUSD = computed(() => roundMoney(pagoUSDDirecto.value + (tipoCambio.value > 0 ? pagoBOBDirecto.value / tipoCambio.value : 0)));
const pagoEquivalenteBOB = computed(() => roundMoney(pagoEquivalenteUSD.value * tipoCambio.value));
const pagoExcedeAdelanto = computed(() => pagoEquivalenteUSD.value > adelantoRequeridoUSD.value);
const pagoMenorAdelanto = computed(() => roundMoney(pagoEquivalenteUSD.value) < roundMoney(adelantoRequeridoUSD.value));
const saldoEstimadoUSD = computed(() => Math.max(roundMoney(Number(form.precio_estimado_usd || 0) - pagoEquivalenteUSD.value), 0));
const saldoEstimadoBOB = computed(() => roundMoney(saldoEstimadoUSD.value * tipoCambio.value));
const totalAdelantos = computed(() => filteredPedidos.value.reduce((total: number, pedido: any) => total + totalPagadoUSD(pedido), 0));
const totalPendiente = computed(() => filteredPedidos.value.reduce((total: number, pedido: any) => total + Number(pedido.saldo_pendiente_usd || 0), 0));
const pedidosEnTransito = computed(() => filteredPedidos.value.filter((pedido: any) => esEstadoTransito(pedido.estado)).length);

onMounted(async () => {
  await cargarDatos();
});

watch([() => form.adelanto_porcentaje, () => form.precio_estimado_usd, () => form.adelanto_modo], () => {
  if (form.adelanto_modo === 'porcentaje') {
    form.adelanto_requerido_usd = adelantoRequeridoUSD.value;
  }
});

async function cargarDatos() {
  loading.value = true;
  try {
    await Promise.all([obtenerPedidos(), obtenerClientes(), obtenerProveedores()]);
  } finally {
    loading.value = false;
  }
}

async function obtenerPedidos() {
  try {
    const query: Record<string, any> = {};
    if (!props.general) {
      const user = localStorage.getItem('user');
      const userId = user ? JSON.parse(user)?.id : null;
      if (userId) query.id_usuario = userId;
    }
    const res = await $fetch(server.HOST + '/api/v1/pedidos', { method: 'GET', query });
    pedidos.value = Array.isArray(res) ? res : [];
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar pedidos', life: 3000 });
  }
}

async function obtenerClientes() {
  try {
    const res = await $fetch(server.HOST + '/api/v1/clientes', { method: 'GET' });
    clientes.value = Array.isArray(res) ? res : [];
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar clientes', life: 3000 });
  }
}

async function obtenerProveedores() {
  try {
    const res = await $fetch(server.HOST + '/api/v1/proveedores-autos', { method: 'GET' });
    proveedores.value = Array.isArray(res) ? res : [];
  } catch (err) {
    proveedores.value = [];
  }
}

function abrirNuevoPedido() {
  resetForm();
  pedidoVisible.value = true;
}

async function registrarPedido() {
  if (!validarPedido()) return;
  saving.value = true;
  try {
    await $fetch(server.HOST + '/api/v1/pedidos', {
      method: 'POST',
      body: {
        id_cliente: form.id_cliente,
        id_proveedor: form.id_proveedor,
        fecha: form.fecha,
        marca: form.marca,
        modelo: form.modelo,
        anio: Number(form.anio || 0),
        color: form.color,
        version: form.version,
        pais_origen: form.pais_origen,
        precio_estimado_usd: Number(form.precio_estimado_usd || 0),
        tipo_cambio: tipoCambio.value,
        fecha_llegada_estimada: form.fecha_llegada_estimada,
        adelanto_requerido_usd: adelantoRequeridoUSD.value,
        adelanto_porcentaje: form.adelanto_modo === 'porcentaje' ? Number(form.adelanto_porcentaje || 0) : 0,
        pagos: form.pagos.map((pago: any) => ({
          moneda: pago.moneda,
          metodo: pago.metodo,
          monto: Number(pago.monto || 0)
        })),
        observacion: form.observacion
      }
    });
    toast.add({ severity: 'success', summary: 'Pedido registrado', life: 3000 });
    pedidoVisible.value = false;
    await obtenerPedidos();
  } catch (err: any) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al registrar pedido', detail: err?.data || err?.message, life: 4000 });
  } finally {
    saving.value = false;
  }
}

async function marcarEnTransito(pedido: any) {
  changingStateId.value = pedido.id;
  try {
    await $fetch(server.HOST + `/api/v1/pedidos/${pedido.id}/aduana`, { method: 'PATCH' });
    toast.add({ severity: 'success', summary: 'Pedido marcado en aduana', life: 3000 });
    await obtenerPedidos();
  } catch (err: any) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'No se pudo actualizar el pedido', detail: err?.data || err?.message, life: 4000 });
  } finally {
    changingStateId.value = null;
  }
}

function validarPedido() {
  if (!form.id_cliente) {
    toast.add({ severity: 'warn', summary: 'Seleccione un cliente', life: 3000 });
    return false;
  }
  if (!form.marca.trim() || !form.modelo.trim() || !form.pais_origen.trim()) {
    toast.add({ severity: 'warn', summary: 'Complete marca, modelo y pais de origen', life: 3000 });
    return false;
  }
  if (Number(form.anio || 0) < 1900 || Number(form.precio_estimado_usd || 0) <= 0 || tipoCambio.value <= 0 || !form.fecha_llegada_estimada) {
    toast.add({ severity: 'warn', summary: 'Complete anio, precio, tipo de cambio y llegada estimada', life: 4000 });
    return false;
  }
  if (adelantoRequeridoUSD.value <= 0 || adelantoRequeridoUSD.value > Number(form.precio_estimado_usd || 0)) {
    toast.add({ severity: 'warn', summary: 'Adelanto no valido', life: 3000 });
    return false;
  }
  if (form.pagos.some((pago: any) => !pago.moneda || !pago.metodo || Number(pago.monto || 0) <= 0)) {
    toast.add({ severity: 'warn', summary: 'Complete moneda, tipo de pago y monto en cada fila', life: 4000 });
    return false;
  }
  if (pagoMenorAdelanto.value) {
    toast.add({ severity: 'warn', summary: 'Pago insuficiente', detail: 'El pago debe ser igual o mayor al adelanto requerido', life: 4000 });
    return false;
  }
  return true;
}

async function guardarProveedor() {
  if (!proveedorForm.nombre.trim()) {
    toast.add({ severity: 'warn', summary: 'Ingrese nombre del proveedor', life: 3000 });
    return;
  }
  savingProveedor.value = true;
  try {
    const nuevo: any = await $fetch(server.HOST + '/api/v1/proveedores-autos', {
      method: 'POST',
      body: {
        ...proveedorForm,
        tipo: 'Importacion',
        estado: 'Activo'
      }
    });
    await obtenerProveedores();
    form.id_proveedor = Number(nuevo?.id || 0) || null;
    proveedorVisible.value = false;
    toast.add({ severity: 'success', summary: 'Proveedor creado', life: 3000 });
  } catch (err: any) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al crear proveedor', detail: err?.data || err?.message, life: 4000 });
  } finally {
    savingProveedor.value = false;
  }
}

function crearPago() {
  return {
    key: `${Date.now()}-${Math.random().toString(36).slice(2)}`,
    moneda: 'USD',
    metodo: 'Efectivo',
    monto: 0
  };
}

function agregarPago() {
  form.pagos.push(crearPago());
}

function eliminarPago(index: number) {
  form.pagos.splice(index, 1);
  if (form.pagos.length === 0) {
    agregarPago();
  }
}

function resetForm() {
  Object.assign(form, {
    id_cliente: null,
    id_proveedor: null,
    fecha: new Date().toISOString().slice(0, 10),
    marca: '',
    modelo: '',
    anio: new Date().getFullYear(),
    color: '',
    version: '',
    pais_origen: '',
    precio_estimado_usd: 0,
    tipo_cambio: null,
    fecha_llegada_estimada: '',
    adelanto_modo: 'porcentaje',
    adelanto_requerido_usd: 0,
    adelanto_porcentaje: 10,
    observacion: ''
  });
  form.pagos.splice(0, form.pagos.length, crearPago());
}

function vehiculoSolicitado(pedido: any) {
  return [pedido.marca, pedido.modelo, pedido.anio].filter(Boolean).join(' ') || 'Vehiculo solicitado';
}

function totalPagadoUSD(pedido: any) {
  const tc = Number(pedido.tipo_cambio_usado || 0);
  return roundMoney(Number(pedido.adelanto_pagado_usd || 0) + (tc > 0 ? Number(pedido.adelanto_pagado_bob || 0) / tc : 0));
}

function normalizarEstadoUI(estado: string) {
  const plano = estadoPlano(estado);
  if (plano.includes('aduana') || plano.includes('transito')) return 'En aduana';
  return String(estado || '').replace('En transito', 'En tránsito').replace('trÃ¡nsito', 'tránsito');
}

function estadoPedidoLabel(estado: string) {
  if (normalizarEstadoUI(estado) === 'En aduana') return 'En aduana';
  return normalizarEstadoUI(estado) === 'En tránsito' ? 'En tránsito' : estado;
}

function estadoPedidoSeverity(estado: string) {
  switch (normalizarEstadoUI(estado)) {
    case 'Pedido registrado':
      return 'secondary';
    case 'En aduana':
      return 'info';
    case 'En tránsito':
      return 'info';
    case 'Recibido':
      return 'warning';
    case 'Completado':
      return 'success';
    default:
      return 'secondary';
  }
}

function estadoPlano(estado: string) {
  return String(estado || '')
    .replace('trÃƒÂ¡nsito', 'transito')
    .replace('trÃ¡nsito', 'transito')
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .toLowerCase();
}

function esEstadoTransito(estado: string) {
  const normalizado = estadoPlano(estado);
  return normalizado.includes('aduana') || normalizado.includes('transito');
}

function esEstadoRecibido(estado: string) {
  return estadoPlano(estado) === 'recibido';
}

function vehiculoNoDisponible(pedido: any) {
  return estadoPlano(pedido.estado) !== 'completado';
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

function roundMoney(value: number) {
  return Math.round(Number(value || 0) * 100) / 100;
}

function mostrarError(summary: string, err: any) {
  toast.add({ severity: 'error', summary, detail: err?.data || err?.message, life: 4000 });
}
</script>
