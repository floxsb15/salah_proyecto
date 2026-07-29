<template>
  <Toast />
  
  <!-- Title -->
  <h2 class="text-2xl font-bold mb-4">Usuarios</h2>

  <!-- Controls Bar -->
  <div class="flex justify-between items-center mb-4">
    <Button label="Nuevo" icon="pi pi-plus" size="small" @click="AgregarUsuario = true"/>
    <span class="p-input-icon-left">
      <i class="pi pi-search" />
      <InputText v-model="searchQuery" placeholder="Buscar..." />
    </span>
  </div>

  <!-- Data Table -->
  <div>
    <DataTable 
      :value="filteredUsuarios" 
      tableStyle="min-width: 50rem" 
      size="small"
      stripedRows
      removableSort
      paginator :rows="10"
      :rowsPerPageOptions="[5, 10, 25, 50]"
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
      currentPageReportTemplate="Mostrando {first} a {last} de {totalRecords} usuarios"
      :pt="{ thead: { class: 'bg-amber-300' } }"
    >
      <template #empty>
        <p class="text-center p-4">No hay usuarios para mostrar.</p>
      </template>
      
      <Column field="id" header="ID" sortable />
      <Column field="nombre" header="Nombre Completo" sortable />
      <Column field="ci" header="CI" />
      <Column field="usuario" header="Usuario" sortable />
      <Column field="rol" header="Rol" sortable />
      <Column header="Acciones" style="width: 10rem">
        <template #body="slotProps">
          <div class="flex items-center justify-center">
            <Button label="Editar" variant="text" size="small" @click="idUsuario = slotProps.data.id; ModificarUsuario = true"/>
          </div>
        </template>
      </Column>
    </DataTable>
  </div>

  <!-- Modals -->
  <modalAgregarUsuario 
    :open="AgregarUsuario"
    v-if="AgregarUsuario"
    @close="AgregarUsuario = false"
    @update="obtenerUsuarios"
    @success="toast.add({ severity: 'success', summary: 'Agregado Exitoso', life: 3000 })"
    @error="toast.add({ severity: 'error', summary: 'Error al Agregar', life: 3000 })" />

  <modalModificarUsuario 
    :open="ModificarUsuario"
    :id="idUsuario"
    v-if="ModificarUsuario"
    @close="ModificarUsuario = false"
    @update="obtenerUsuarios"
    @success="toast.add({ severity: 'success', summary: 'Modificado Exitoso', life: 3000 })"
    @error="toast.add({ severity: 'error', summary: 'Error al Modificar', life: 3000 })" />
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import modalAgregarUsuario from '~/components/admin/usuarios/modalAgregarUsuario.vue';
import modalModificarUsuario from '~/components/admin/usuarios/modalModificarUsuario.vue';
import { server } from '~/server/server';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';

definePageMeta({ layout : 'menu-admin' });

const toast = useToast();
const Usuarios = ref<any[]>([]);
const AgregarUsuario = ref(false), ModificarUsuario = ref(false);
const idUsuario = ref(0);
const searchQuery = ref('');
const topSearch = useTopSearch();
const activeSearchQuery = computed(() => (topSearch.value || searchQuery.value).trim().toLowerCase());

const filteredUsuarios = computed(() => {
  if (!activeSearchQuery.value) {
    return Usuarios.value;
  }
  const query = activeSearchQuery.value;
  return Usuarios.value.filter(user =>
    (user.id?.toString() || '').includes(query) ||
    (user.nombre?.toLowerCase() || '').includes(query) ||
    (user.ci?.toLowerCase() || '').includes(query) ||
    (user.usuario?.toLowerCase() || '').includes(query) ||
    (user.rol?.toLowerCase() || '').includes(query) ||
    (user.direccion?.toLowerCase() || '').includes(query)
  );
});

onMounted(async () => {
  await obtenerUsuarios();
});

async function obtenerUsuarios() {
  try {
    const res:any[] = await $fetch(server.HOST + '/api/v1/usuarios', {
      method: 'GET'
    });
    Usuarios.value = res;
  } catch (err) {
    console.error(err);
  }
}
</script>
