<template>
  <Toast />

  <h2 class="text-2xl font-bold mb-4">Clientes</h2>

  <div class="flex justify-between items-center mb-4">
    <Button label="Nuevo" icon="pi pi-plus" size="small" @click="AgregarCliente = true"/>
    <span class="p-input-icon-left">
      <i class="pi pi-search" />
      <InputText v-model="searchQuery" placeholder="Buscar..." />
    </span>
  </div>

  <div>
    <DataTable
      :value="filteredClientes"
      tableStyle="min-width: 50rem"
      size="small"
      stripedRows
      removableSort
      paginator :rows="10"
      :rowsPerPageOptions="[5, 10, 25, 50]"
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
      currentPageReportTemplate="Mostrando {first} a {last} de {totalRecords} clientes"
      :pt="{ thead: { class: 'bg-amber-300' } }"
    >
      <template #empty>
        <p class="text-center p-4">No hay clientes para mostrar.</p>
      </template>

      <Column field="id" header="ID" sortable />
      <Column field="nombreCompleto" header="Nombre Completo" sortable />
      <Column field="ci" header="CI/NIT" sortable />
      <Column field="celular" header="Celular" />
      <Column field="direccion" header="Dirección" />
      <Column field="estado" header="Estado" sortable>
        <template #body="slotProps">
          <Tag :value="slotProps.data.estado" :severity="slotProps.data.estado === 'Activo' ? 'success' : 'danger'" />
        </template>
      </Column>
      <Column header="Acciones" style="width: 18rem">
        <template #body="slotProps">
          <div class="flex items-center justify-center gap-2">
            <Button label="Detalle" variant="text" size="small" @click="idCliente = slotProps.data.id; DetalleCliente = true"/>
            <Button label="Editar" variant="text" size="small" @click="idCliente = slotProps.data.id; ModificarCliente = true"/>
            <Button
              v-if="slotProps.data.estado === 'Activo'"
              label="Desactivar"
              variant="text"
              size="small"
              severity="danger"
              @click="desactivarCliente(slotProps.data)"
            />
          </div>
        </template>
      </Column>
    </DataTable>
  </div>

  <modalAgregarCliente
    :open="AgregarCliente"
    v-if="AgregarCliente"
    @close="AgregarCliente = false"
    @update="obtenerClientes"
    @success="toast.add({ severity: 'success', summary: 'Agregado Exitoso', life: 3000 })"
    @error="toast.add({ severity: 'error', summary: 'Error al Agregar', life: 3000 })" />

  <modalModificarCliente
    :open="ModificarCliente"
    :id="idCliente"
    v-if="ModificarCliente"
    @close="ModificarCliente = false"
    @update="obtenerClientes"
    @success="toast.add({ severity: 'success', summary: 'Modificado Exitoso', life: 3000 })"
    @error="toast.add({ severity: 'error', summary: 'Error al Modificar', life: 3000 })" />

  <modalDetalleCliente
    :open="DetalleCliente"
    :id="idCliente"
    v-if="DetalleCliente"
    @close="DetalleCliente = false"
    @error="toast.add({ severity: 'error', summary: 'Error al cargar detalle', life: 3000 })" />
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import modalAgregarCliente from '~/components/admin/clientes/modalAgregarCliente.vue';
import modalDetalleCliente from '~/components/admin/clientes/modalDetalleCliente.vue';
import modalModificarCliente from '~/components/admin/clientes/modalModificarCliente.vue';
import { server } from '~/server/server';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import Tag from 'primevue/tag';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';

definePageMeta({ layout : 'menu-ventas' });

const toast = useToast();
const Clientes = ref<any[]>([]);
const AgregarCliente = ref(false);
const DetalleCliente = ref(false);
const ModificarCliente = ref(false);
const idCliente = ref(0);
const searchQuery = ref('');
const topSearch = useTopSearch();
const activeSearchQuery = computed(() => (topSearch.value || searchQuery.value).trim().toLowerCase());

const clientesConNombreCompleto = computed(() => {
  return Clientes.value.map(cliente => ({
    ...cliente,
    nombreCompleto: `${cliente.nombre || ''} ${cliente.apellido || ''}`.trim()
  }));
});

const filteredClientes = computed(() => {
  if (!activeSearchQuery.value) {
    return clientesConNombreCompleto.value;
  }

  const query = activeSearchQuery.value;
  return clientesConNombreCompleto.value.filter(cliente =>
    (cliente.id?.toString() || '').includes(query) ||
    (cliente.nombreCompleto?.toLowerCase() || '').includes(query) ||
    (cliente.ci?.toLowerCase() || '').includes(query) ||
    (cliente.celular?.toLowerCase() || '').includes(query) ||
    (cliente.direccion?.toLowerCase() || '').includes(query) ||
    (cliente.estado?.toLowerCase() || '').includes(query)
  );
});

onMounted(async () => {
  await obtenerClientes();
});

async function obtenerClientes() {
  try {
    const res = await $fetch(server.HOST + '/api/v1/clientes', {
      method: 'GET'
    });
    Clientes.value = Array.isArray(res) ? res : [];
  } catch (err) {
    console.error(err);
  }
}

async function desactivarCliente(cliente: any) {
  try {
    await $fetch(server.HOST + '/api/v1/clientes/' + cliente.id, {
      method: 'PUT',
      body: {
        nombre: cliente.nombre,
        apellido: cliente.apellido,
        ci: cliente.ci,
        celular: cliente.celular,
        direccion: cliente.direccion,
        estado: 'Inactivo'
      }
    });

    await obtenerClientes();
    toast.add({ severity: 'success', summary: 'Cliente Desactivado', life: 3000 });
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al Desactivar', life: 3000 });
  }
}
</script>

