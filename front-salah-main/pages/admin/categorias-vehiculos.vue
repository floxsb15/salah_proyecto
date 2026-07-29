<template>
  <Toast />

  <h2 class="mb-4 text-2xl font-bold">Categorias de Vehiculo</h2>

  <div class="mb-4 flex items-center justify-between">
    <div class="flex gap-2">
      <Button label="Nueva Categoria" icon="pi pi-plus" size="small" @click="AgregarCategoria = true" />
      <Button label="Nuevo Segmento" icon="pi pi-sitemap" size="small" severity="secondary" @click="AgregarSegmento = true" />
      <Button label="Nueva Marca" icon="pi pi-tag" size="small" severity="secondary" @click="AgregarMarca = true" />
      <Button label="Nuevo Anio" icon="pi pi-calendar" size="small" severity="secondary" @click="AgregarAnio = true" />
    </div>
    <span class="p-input-icon-left">
      <i class="pi pi-search" />
      <InputText v-model="searchQuery" placeholder="Buscar..." />
    </span>
  </div>

  <div>
    <DataTable
      :value="filteredCategorias"
      tableStyle="min-width: 50rem"
      size="small"
      stripedRows
      removableSort
      paginator
      :rows="10"
      :rowsPerPageOptions="[10, 20, 50]"
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
      currentPageReportTemplate="Mostrando {first} a {last} de {totalRecords} categorias"
      :pt="{ thead: { class: 'bg-amber-300' } }"
    >
      <template #empty>
        <p class="p-4 text-center">No hay categorias para mostrar.</p>
      </template>

      <Column field="id" header="ID" sortable />
      <Column field="nombre" header="Nombre" sortable />
      <Column field="descripcion" header="Descripcion" />
      <Column field="estado" header="Estado" sortable />
      <Column header="Acciones" style="width: 10rem">
        <template #body="slotProps">
          <div class="flex items-center justify-center">
            <Button label="Editar" variant="text" size="small" @click="idCategoria = slotProps.data.id; ModificarCategoria = true" />
          </div>
        </template>
      </Column>
    </DataTable>
  </div>

  <h3 class="mb-4 mt-6 text-xl font-bold">Marcas de Vehiculo</h3>

  <div>
    <DataTable
      :value="filteredMarcas"
      tableStyle="min-width: 50rem"
      size="small"
      stripedRows
      removableSort
      paginator
      :rows="10"
      :rowsPerPageOptions="[10, 20, 50]"
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
      currentPageReportTemplate="Mostrando {first} a {last} de {totalRecords} marcas"
      :pt="{ thead: { class: 'bg-amber-300' } }"
    >
      <template #empty>
        <p class="p-4 text-center">No hay marcas para mostrar.</p>
      </template>

      <Column field="id" header="ID" sortable />
      <Column field="nombre" header="Marca" sortable />
      <Column field="estado" header="Estado" sortable />
    </DataTable>
  </div>

  <h3 class="mb-4 mt-6 text-xl font-bold">Anios de Vehiculo</h3>

  <div>
    <DataTable
      :value="filteredAnios"
      tableStyle="min-width: 50rem"
      size="small"
      stripedRows
      removableSort
      paginator
      :rows="10"
      :rowsPerPageOptions="[10, 20, 50]"
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
      currentPageReportTemplate="Mostrando {first} a {last} de {totalRecords} anios"
      :pt="{ thead: { class: 'bg-amber-300' } }"
    >
      <template #empty>
        <p class="p-4 text-center">No hay anios para mostrar.</p>
      </template>

      <Column field="id" header="ID" sortable />
      <Column field="valor" header="Anio" sortable />
      <Column field="estado" header="Estado" sortable />
    </DataTable>
  </div>

  <h3 class="mb-4 mt-6 text-xl font-bold">Segmentos de Vehiculo</h3>


  <div>
    <DataTable
      :value="filteredSegmentos"
      tableStyle="min-width: 50rem"
      size="small"
      stripedRows
      removableSort
      paginator
      :rows="10"
      :rowsPerPageOptions="[10, 20, 50]"
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
      currentPageReportTemplate="Mostrando {first} a {last} de {totalRecords} segmentos"
      :pt="{ thead: { class: 'bg-amber-300' } }"
    >
      <template #empty>
        <p class="p-4 text-center">No hay segmentos para mostrar.</p>
      </template>

      <Column field="id" header="ID" sortable />
      <Column field="categoria" header="Categoria" sortable />
      <Column field="nombre" header="Segmento" sortable />
      <Column field="descripcion" header="Descripcion" />
      <Column field="estado" header="Estado" sortable />
      <Column header="Acciones" style="width: 10rem">
        <template #body="slotProps">
          <div class="flex items-center justify-center">
            <Button label="Editar" variant="text" size="small" @click="idSegmento = slotProps.data.id; ModificarSegmento = true" />
          </div>
        </template>
      </Column>
    </DataTable>
  </div>

  <modalAgregarCategoria
    :open="AgregarCategoria"
    v-if="AgregarCategoria"
    @close="AgregarCategoria = false"
    @update="obtenerCategorias"
    @success="toast.add({ severity: 'success', summary: 'Agregado Exitoso', life: 3000 })"
    @error="toast.add({ severity: 'error', summary: 'Error al Agregar', life: 3000 })" />

  <modalModificarCategoria
    :open="ModificarCategoria"
    :id="idCategoria"
    v-if="ModificarCategoria"
    @close="ModificarCategoria = false"
    @update="obtenerCategorias"
    @success="toast.add({ severity: 'success', summary: 'Modificado Exitoso', life: 3000 })"
    @error="toast.add({ severity: 'error', summary: 'Error al Modificar', life: 3000 })" />

  <modalAgregarSegmento
    :open="AgregarSegmento"
    v-if="AgregarSegmento"
    @close="AgregarSegmento = false"
    @update="obtenerSegmentos"
    @success="toast.add({ severity: 'success', summary: 'Segmento Agregado', life: 3000 })"
    @error="toast.add({ severity: 'error', summary: 'Error al Agregar Segmento', life: 3000 })" />

  <modalAgregarMarca
    :open="AgregarMarca"
    v-if="AgregarMarca"
    @close="AgregarMarca = false"
    @update="obtenerMarcas"
    @success="toast.add({ severity: 'success', summary: 'Marca Agregada', life: 3000 })"
    @error="toast.add({ severity: 'error', summary: 'Error al Agregar Marca', life: 3000 })" />

  <modalAgregarAnio
    :open="AgregarAnio"
    v-if="AgregarAnio"
    @close="AgregarAnio = false"
    @update="obtenerAnios"
    @success="toast.add({ severity: 'success', summary: 'Anio Agregado', life: 3000 })"
    @error="toast.add({ severity: 'error', summary: 'Error al Agregar Anio', life: 3000 })" />

  <modalModificarSegmento
    :open="ModificarSegmento"
    :id="idSegmento"
    v-if="ModificarSegmento"
    @close="ModificarSegmento = false"
    @update="obtenerSegmentos"
    @success="toast.add({ severity: 'success', summary: 'Segmento Modificado', life: 3000 })"
    @error="toast.add({ severity: 'error', summary: 'Error al Modificar Segmento', life: 3000 })" />
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import modalAgregarCategoria from '~/components/admin/categorias-vehiculos/modalAgregarCategoria.vue';
import modalModificarCategoria from '~/components/admin/categorias-vehiculos/modalModificarCategoria.vue';
import modalAgregarMarca from '~/components/admin/categorias-vehiculos/modalAgregarMarcaVehiculo.vue';
import modalAgregarAnio from '~/components/admin/categorias-vehiculos/modalAgregarAnioVehiculo.vue';
import modalAgregarSegmento from '~/components/admin/segmentos-vehiculos/modalAgregarSegmento.vue';
import modalModificarSegmento from '~/components/admin/segmentos-vehiculos/modalModificarSegmento.vue';
import { server } from '~/server/server';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';

definePageMeta({ layout : 'menu-admin' });

const toast = useToast();
const Categorias = ref<any[]>([]);
const Segmentos = ref<any[]>([]);
const Marcas = ref<any[]>([]);
const Anios = ref<any[]>([]);
const AgregarCategoria = ref(false);
const ModificarCategoria = ref(false);
const AgregarSegmento = ref(false);
const ModificarSegmento = ref(false);
const AgregarMarca = ref(false);
const AgregarAnio = ref(false);
const idCategoria = ref(0);
const idSegmento = ref(0);
const searchQuery = ref('');
const topSearch = useTopSearch();
const activeSearchQuery = computed(() => (topSearch.value || searchQuery.value).trim().toLowerCase());

const filteredCategorias = computed(() => {
  if (!activeSearchQuery.value) {
    return Categorias.value;
  }
  const query = activeSearchQuery.value;
  return Categorias.value.filter(cat =>
    (cat.id?.toString() || '').includes(query) ||
    (cat.nombre?.toLowerCase() || '').includes(query) ||
    (cat.descripcion?.toLowerCase() || '').includes(query) ||
    (cat.estado?.toLowerCase() || '').includes(query)
  );
});

const filteredSegmentos = computed(() => {
  if (!activeSearchQuery.value) {
    return Segmentos.value;
  }
  const query = activeSearchQuery.value;
  return Segmentos.value.filter(segmento =>
    (segmento.id?.toString() || '').includes(query) ||
    (segmento.nombre?.toLowerCase() || '').includes(query) ||
    (segmento.descripcion?.toLowerCase() || '').includes(query) ||
    (segmento.categoria?.toLowerCase() || '').includes(query) ||
    (segmento.estado?.toLowerCase() || '').includes(query)
  );
});

const filteredMarcas = computed(() => {
  if (!activeSearchQuery.value) {
    return Marcas.value;
  }
  const query = activeSearchQuery.value;
  return Marcas.value.filter(marca =>
    (marca.id?.toString() || '').includes(query) ||
    (marca.nombre?.toLowerCase() || '').includes(query) ||
    (marca.estado?.toLowerCase() || '').includes(query)
  );
});

const filteredAnios = computed(() => {
  if (!activeSearchQuery.value) {
    return Anios.value;
  }
  const query = activeSearchQuery.value;
  return Anios.value.filter(anio =>
    (anio.id?.toString() || '').includes(query) ||
    (anio.valor?.toString() || '').includes(query) ||
    (anio.estado?.toLowerCase() || '').includes(query)
  );
});

onMounted(async () => {
  await Promise.all([obtenerCategorias(), obtenerSegmentos(), obtenerMarcas(), obtenerAnios()]);
});

async function obtenerCategorias() {
  try {
    const res:any[] = await $fetch(server.HOST + '/api/v1/categorias-vehiculos', {
      method: 'GET'
    });
    Categorias.value = res;
  } catch (err) {
    console.error(err);
  }
}

async function obtenerSegmentos() {
  try {
    const res:any[] = await $fetch(server.HOST + '/api/v1/segmentos-vehiculos', {
      method: 'GET'
    });
    Segmentos.value = res;
  } catch (err) {
    console.error(err);
  }
}

async function obtenerMarcas() {
  try {
    const res:any[] = await $fetch(server.HOST + '/api/v1/marcas-vehiculos', {
      method: 'GET'
    });
    Marcas.value = res;
  } catch (err) {
    console.error(err);
  }
}

async function obtenerAnios() {
  try {
    const res:any[] = await $fetch(server.HOST + '/api/v1/anios-vehiculos', {
      method: 'GET'
    });
    Anios.value = res;
  } catch (err) {
    console.error(err);
  }
}
</script>
