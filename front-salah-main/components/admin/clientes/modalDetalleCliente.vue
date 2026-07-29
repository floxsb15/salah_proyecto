<template>
  <Dialog
    v-model:visible="visible"
    modal
    :showHeader="false"
    :style="{ width: '78rem', maxWidth: '96vw' }"
    :pt="{ content: { class: 'salah-detail-content' } }"
  >
    <div class="client-detail">
      <header class="detail-header">
        <div>
          <p class="eyebrow">Detalle del cliente</p>
          <h2>{{ nombreCompleto || 'Cliente' }}</h2>
          <p class="muted">CI/NIT: {{ cliente?.ci || 'N/A' }} · Celular: {{ cliente?.celular || 'N/A' }}</p>
        </div>
        <button type="button" class="close-button" aria-label="Cerrar" @click="visible = false">
          <i class="pi pi-times"></i>
        </button>
      </header>

      <div class="tabs" role="tablist">
        <button type="button" :class="{ active: activeTab === 'info' }" @click="activeTab = 'info'">Informacion</button>
        <button type="button" :class="{ active: activeTab === 'historial' }" @click="activeTab = 'historial'">Historial de ventas</button>
      </div>

      <section v-if="activeTab === 'info'" class="panel">
        <div class="summary-grid">
          <div class="summary-item">
            <span>Estado</span>
            <Tag :value="cliente?.estado || 'N/A'" :severity="cliente?.estado === 'Activo' ? 'success' : 'danger'" />
          </div>
          <div class="summary-item">
            <span>Direccion</span>
            <strong>{{ cliente?.direccion || 'N/A' }}</strong>
          </div>
          <div class="summary-item">
            <span>Total gastado</span>
            <strong>$ {{ formatPrecio(historial.total_gastado) }}</strong>
          </div>
          <div class="summary-item">
            <span>Saldo pendiente</span>
            <strong>$ {{ formatPrecio(historial.total_pendiente) }}</strong>
          </div>
        </div>
      </section>

      <section v-else class="panel">
        <div class="history-header">
          <div>
            <h3>Historial de ventas</h3>
            <p>{{ historial.compras.length }} registros encontrados</p>
          </div>
          <div class="history-total">
            <span>Total gastado</span>
            <strong>$ {{ formatPrecio(historial.total_gastado) }}</strong>
          </div>
        </div>

        <DataTable
          :value="historial.compras"
          :loading="loading"
          tableStyle="min-width: 72rem"
          size="small"
          stripedRows
          paginator
          :rows="5"
          :rowsPerPageOptions="[5, 10, 20]"
        >
          <template #empty>
            <p class="p-4 text-center">No hay ventas ni reservas registradas.</p>
          </template>

          <Column field="vehiculo" header="Vehiculo" sortable>
            <template #body="slotProps">
              <div>
                <p class="font-medium text-gray-900">{{ slotProps.data.vehiculo }}</p>
                <p class="text-xs text-gray-500">{{ slotProps.data.categoria || 'Sin categoria' }}</p>
              </div>
            </template>
          </Column>
          <Column field="fecha" header="Fecha" sortable>
            <template #body="slotProps">
              {{ formatFecha(slotProps.data.fecha) }}
            </template>
          </Column>
          <Column field="tipo" header="Tipo" sortable />
          <Column field="monto_total" header="Monto total" sortable>
            <template #body="slotProps">
              $ {{ formatPrecio(slotProps.data.monto_total) }}
            </template>
          </Column>
          <Column field="monto_pagado" header="Monto pagado" sortable>
            <template #body="slotProps">
              $ {{ formatPrecio(slotProps.data.monto_pagado) }}
            </template>
          </Column>
          <Column field="saldo_pendiente" header="Saldo pendiente" sortable>
            <template #body="slotProps">
              $ {{ formatPrecio(slotProps.data.saldo_pendiente) }}
            </template>
          </Column>
          <Column field="estado" header="Estado" sortable>
            <template #body="slotProps">
              <Tag :value="slotProps.data.estado" :severity="slotProps.data.estado === 'Pagado' ? 'success' : 'warning'" />
            </template>
          </Column>
        </DataTable>
      </section>
    </div>
  </Dialog>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { server } from '~/server/server';
import Column from 'primevue/column';
import DataTable from 'primevue/datatable';
import Dialog from 'primevue/dialog';
import Tag from 'primevue/tag';

interface Props { open: boolean, id: number }

const props = defineProps<Props>();
const emit = defineEmits(['close', 'error']);

const visible = ref(props.open);
const activeTab = ref('historial');
const loading = ref(false);
const cliente = ref<any>(null);
const historial = reactive({
  compras: [] as any[],
  cuotas_pendientes: [] as any[],
  total_gastado: 0,
  total_pendiente: 0,
  tiene_credito_activo: false
});

const nombreCompleto = computed(() => {
  return `${cliente.value?.nombre || ''} ${cliente.value?.apellido || ''}`.trim();
});

onMounted(async () => {
  await cargarDetalle();
});

watch(() => props.open, (newValue) => {
  visible.value = newValue;
});

watch(visible, (newValue) => {
  if (!newValue) {
    emit('close');
  }
});

async function cargarDetalle() {
  loading.value = true;

  try {
    const [resCliente, resHistorial]: any[] = await Promise.all([
      $fetch(server.HOST + '/api/v1/clientes/' + props.id, { method: 'GET' }),
      $fetch(server.HOST + '/api/v1/clientes/' + props.id + '/historial-compras', { method: 'GET' })
    ]);

    cliente.value = resCliente || null;
    historial.compras = Array.isArray(resHistorial?.compras) ? resHistorial.compras : [];
    historial.cuotas_pendientes = Array.isArray(resHistorial?.cuotas_pendientes) ? resHistorial.cuotas_pendientes : [];
    historial.total_gastado = Number(resHistorial?.total_gastado || 0);
    historial.total_pendiente = Number(resHistorial?.total_pendiente || 0);
    historial.tiene_credito_activo = !!resHistorial?.tiene_credito_activo;
  } catch (err) {
    console.error(err);
    emit('error');
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
</script>

<style scoped>
:deep(.salah-detail-content) { padding: 0; overflow: hidden; background: #f7f7f7; }
.client-detail { color: #111827; }
.detail-header { display: flex; align-items: center; justify-content: space-between; gap: 16px; padding: 22px 26px; background: #111111; color: #ffffff; }
.eyebrow { margin: 0 0 4px; color: #ffd700; font-size: .76rem; font-weight: 800; text-transform: uppercase; }
.detail-header h2 { margin: 0; font-size: 1.45rem; font-weight: 800; letter-spacing: 0; }
.muted { margin: 4px 0 0; color: rgba(255,255,255,.72); font-size: .9rem; }
.close-button { display: grid; width: 40px; height: 40px; place-items: center; border: 1px solid rgba(255,255,255,.16); border-radius: 8px; background: rgba(255,255,255,.08); color: #ffffff; cursor: pointer; }
.tabs { display: flex; gap: 8px; padding: 14px 18px 0; background: #ffffff; border-bottom: 1px solid #e5e7eb; }
.tabs button { min-height: 42px; border: 0; border-bottom: 3px solid transparent; background: transparent; color: #4b5563; font-weight: 800; cursor: pointer; }
.tabs button.active { border-bottom-color: #ffd700; color: #111111; }
.panel { padding: 18px; }
.summary-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 12px; }
.summary-item, .history-total { min-width: 0; padding: 14px; border: 1px solid #e5e7eb; border-radius: 8px; background: #ffffff; }
.summary-item span, .history-total span { display: block; color: #6b7280; font-size: .8rem; font-weight: 700; }
.summary-item strong, .history-total strong { display: block; margin-top: 6px; color: #111827; font-size: 1rem; }
.history-header { display: flex; align-items: center; justify-content: space-between; gap: 12px; margin-bottom: 14px; }
.history-header h3 { margin: 0; font-size: 1.1rem; font-weight: 800; }
.history-header p, .credit-box p { margin: 3px 0 0; color: #6b7280; font-size: .86rem; }
.credit-box { display: grid; grid-template-columns: minmax(0, 1fr) minmax(260px, .8fr); gap: 12px; margin-bottom: 14px; padding: 14px; border: 1px solid #facc15; border-radius: 8px; background: #fffbeb; }
.credit-box h4 { margin: 0; font-size: 1rem; font-weight: 800; }
.credit-box ul { display: grid; gap: 8px; margin: 0; padding: 0; list-style: none; }
.credit-box li { display: flex; align-items: center; justify-content: space-between; gap: 10px; color: #111827; }
@media (max-width: 860px) {
  .summary-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .credit-box { grid-template-columns: 1fr; }
}
@media (max-width: 560px) {
  .detail-header { align-items: flex-start; padding: 18px; }
  .summary-grid { grid-template-columns: 1fr; }
  .history-header { align-items: stretch; flex-direction: column; }
}
</style>
