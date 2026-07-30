<template>
  <Toast />

  <!-- Title -->
  <h2 class="text-2xl font-bold mb-4">Vehículos</h2>

  <!-- Controls Bar -->
  <div class="flex justify-between items-center mb-4">
    <div class="flex gap-2">
      <Button label="Nuevo" icon="pi pi-plus" size="small" @click="AgregarVehiculo = true"/>
      <Button label="Generar Reporte" icon="pi pi-file-pdf" size="small" severity="secondary" @click="reporteVehiculos" />
    </div>
    <span class="p-input-icon-left">
      <i class="pi pi-search" />
      <InputText v-model="searchQuery" placeholder="Buscar..." />
    </span>
  </div>

  <!-- Data Table -->
  <div>
    <DataTable 
      :value="filteredVehiculos" 
      tableStyle="min-width: 50rem" 
      size="small"
      stripedRows
      removableSort
      paginator :rows="10"
      :rowsPerPageOptions="[10, 20, 50]"
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
      currentPageReportTemplate="Mostrando {first} a {last} de {totalRecords} vehículos"
      :pt="{ thead: { class: 'bg-amber-300' } }"
    >
      <template #empty>
        <p class="text-center p-4">No hay vehículos para mostrar.</p>
      </template>
      
      <Column field="id" header="ID" sortable />
      <Column header="Vehiculo" sortable>
        <template #body="slotProps">
          {{ etiquetaVehiculo(slotProps.data) }}
        </template>
      </Column>
      <Column field="marca" header="Marca" sortable />
      <Column field="modelo" header="Modelo" sortable />
      <Column field="anio" header="Anio" sortable />
      <Column field="precio" header="Precio" sortable >
        <template #body="slotProps">
          $ {{ slotProps.data.precio }}
        </template>
      </Column>
      <Column field="precio_compra" header="Precio compra" sortable>
        <template #body="slotProps">
          $ {{ formatPrecio(slotProps.data.precio_compra) }}
        </template>
      </Column>
      <Column field="margen_ganancia" header="Margen" sortable>
        <template #body="slotProps">
          $ {{ formatPrecio(slotProps.data.margen_ganancia) }}
        </template>
      </Column>
      <Column field="cantidad_disponible" header="Disponible" sortable />
      <Column field="categoria" header="Categoría" sortable />
      <Column field="segmento" header="Segmento" sortable />
      <Column header="Acciones" style="width: 10rem">
        <template #body="slotProps">
          <div class="flex items-center justify-center">
            <Button label="Editar" variant="text" size="small" @click="idVehiculo = slotProps.data.id; ModificarVehiculo = true"/>
          </div>
        </template>
      </Column>
    </DataTable>
  </div>

  <!-- Modals -->
  <modalAgregarVehiculo 
    :open="AgregarVehiculo"
    v-if="AgregarVehiculo"
    @close="AgregarVehiculo = false"
    @update="obtenerVehiculos"
    @success="toast.add({ severity: 'success', summary: 'Agregado Exitoso', life: 3000 })"
    @error="mostrarError('Error al Agregar', $event)" />

  <modalModificarVehiculo 
    :open="ModificarVehiculo"
    :id="idVehiculo"
    v-if="ModificarVehiculo"
    @close="ModificarVehiculo = false"
    @update="obtenerVehiculos"
    @success="toast.add({ severity: 'success', summary: 'Modificado Exitoso', life: 3000 })"
    @error="mostrarError('Error al Modificar', $event)" />
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import modalAgregarVehiculo from '~/components/admin/vehiculos/modalAgregarProducto.vue';
import modalModificarVehiculo from '~/components/admin/vehiculos/modalModificarProducto.vue';
import { server } from '~/server/server';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';

definePageMeta({ layout : 'menu-admin' });

const toast = useToast();
const Vehiculos = ref<any[]>([]);
const AgregarVehiculo = ref(false);
const ModificarVehiculo = ref(false);
const idVehiculo = ref(0);
const searchQuery = ref('');
const topSearch = useTopSearch();
const activeSearchQuery = computed(() => (topSearch.value || searchQuery.value).trim().toLowerCase());

const filteredVehiculos = computed(() => {
  if (!activeSearchQuery.value) {
    return Vehiculos.value;
  }
  const query = activeSearchQuery.value;
  return Vehiculos.value.filter(v => 
    (v.id?.toString() || '').includes(query) || 
    etiquetaVehiculo(v).toLowerCase().includes(query) ||
    (v.categoria?.toLowerCase() || '').includes(query) ||
    (v.segmento?.toLowerCase() || '').includes(query) ||
    (v.marca?.toLowerCase() || '').includes(query) ||
    (v.modelo?.toLowerCase() || '').includes(query) ||
    (v.anio?.toString() || '').includes(query) ||
    (v.cantidad_disponible?.toString() || '').includes(query) ||
    (v.version?.toLowerCase() || '').includes(query) ||
    (v.combustible?.toLowerCase() || '').includes(query) ||
    (v.estado?.toLowerCase() || '').includes(query)
  );
});

onMounted(async () => {
  await obtenerVehiculos();
});

async function obtenerVehiculos() {
  try {
    const res:any[] = await $fetch(server.HOST + '/api/v1/vehiculos', {
      method: 'GET'
    });
    res.sort((a, b) => a.categoria.localeCompare(b.categoria));
    Vehiculos.value = res;
  } catch (err) {
    console.error(err);
  }
}

async function reporteVehiculos() {
  try {
    const reporte = await fetch(server.HOST + '/api/v1/reportes/vehiculos', {
      method: 'GET'
    });
    const blob = await reporte.blob();
    const url = URL.createObjectURL(blob);
    window.open(url, '_blank');
  } catch(err) {
    console.error(err);
  }
}

function etiquetaVehiculo(vehiculo: any) {
  return [vehiculo.marca, vehiculo.modelo, vehiculo.anio].filter(Boolean).join(' ') || vehiculo.nombre || 'Vehiculo';
}

function formatPrecio(precio: number) {
  return Number(precio || 0).toLocaleString('es-BO', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  });
}

function mostrarError(summary: string, err: any) {
  toast.add({
    severity: 'error',
    summary,
    detail: err?.data || err?.message || 'Revise los datos enviados.',
    life: 4000
  });
}
</script>
