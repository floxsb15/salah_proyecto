<template>
  <Toast />

  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <h2 class="text-2xl font-bold">Proveedores</h2>
        <p class="text-sm text-gray-500">Personas, concesionarias e importadoras de quienes se compran vehiculos.</p>
      </div>
      <div class="flex flex-col gap-2 sm:flex-row sm:items-center">
        <Button label="Nuevo" icon="pi pi-plus" size="small" @click="abrirNuevo" />
        <span class="p-input-icon-left">
          <i class="pi pi-search" />
          <InputText v-model="searchQuery" placeholder="Buscar..." size="small" />
        </span>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-3 md:grid-cols-4">
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Proveedores</p>
        <p class="text-2xl font-bold text-gray-900">{{ filteredProveedores.length }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Activos</p>
        <p class="text-2xl font-bold text-gray-900">{{ proveedoresActivos }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Total comprado USD</p>
        <p class="text-2xl font-bold text-gray-900">$ {{ formatPrecio(totalCompradoUSD) }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Total comprado BOB</p>
        <p class="text-2xl font-bold text-gray-900">Bs {{ formatPrecio(totalCompradoBOB) }}</p>
      </div>
    </div>

    <DataTable
      :value="filteredProveedores"
      :loading="loading"
      tableStyle="min-width: 72rem"
      size="small"
      stripedRows
      removableSort
      paginator
      :rows="10"
      :rowsPerPageOptions="[10, 20, 50]"
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
      currentPageReportTemplate="Mostrando {first} a {last} de {totalRecords} proveedores"
      :pt="{ thead: { class: 'bg-amber-300' } }"
    >
      <template #empty>
        <p class="p-4 text-center">No hay proveedores para mostrar.</p>
      </template>

      <Column field="nombre" header="Nombre" sortable />
      <Column field="ci_nit" header="CI/NIT" sortable />
      <Column field="telefono" header="Telefono" />
      <Column field="tipo" header="Tipo" sortable />
      <Column field="cantidad_compras" header="Compras" sortable />
      <Column field="total_comprado_usd" header="Total USD" sortable>
        <template #body="slotProps">$ {{ formatPrecio(slotProps.data.total_comprado_usd) }}</template>
      </Column>
      <Column field="total_comprado_bob" header="Total BOB" sortable>
        <template #body="slotProps">Bs {{ formatPrecio(slotProps.data.total_comprado_bob) }}</template>
      </Column>
      <Column field="estado" header="Estado" sortable>
        <template #body="slotProps">
          <Tag :value="slotProps.data.estado" :severity="slotProps.data.estado === 'Activo' ? 'success' : 'danger'" />
        </template>
      </Column>
      <Column header="Acciones" style="width: 16rem">
        <template #body="slotProps">
          <div class="flex items-center justify-center gap-2">
            <Button label="Detalle" variant="text" size="small" @click="abrirDetalle(slotProps.data)" />
            <Button label="Editar" variant="text" size="small" @click="abrirEditar(slotProps.data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <Dialog
      v-model:visible="formVisible"
      modal
      :showHeader="false"
      :style="{ width: '62rem', maxWidth: '96vw' }"
      :pt="{ content: { class: 'salah-dialog-content' } }"
    >
      <div class="salah-user-modal">
        <header class="modal-header">
          <div class="header-icon" aria-hidden="true">
            <svg viewBox="0 0 24 24">
              <path d="M3 21h18" />
              <path d="M5 21V7l8-4v18" />
              <path d="M19 21V11l-6-4" />
              <path d="M9 9h1" />
              <path d="M9 13h1" />
              <path d="M9 17h1" />
            </svg>
          </div>
          <div class="header-copy">
            <h2>{{ form.id ? 'Editar Proveedor' : 'Nuevo Proveedor' }}</h2>
            <p>Registro de compras de Salah Motors</p>
            <span></span>
          </div>
          <button type="button" class="close-button" aria-label="Cerrar" @click="formVisible = false">
            <svg viewBox="0 0 24 24">
              <path d="M18 6 6 18" />
              <path d="m6 6 12 12" />
            </svg>
          </button>
        </header>

        <Form v-slot="$form" :resolver="resolver" :initialValues="form" @submit="guardarProveedor" class="user-form">
          <section class="form-card identity-card">
            <div class="section-title">
              <span aria-hidden="true">
                <svg viewBox="0 0 24 24">
                  <rect width="18" height="14" x="3" y="5" rx="2" />
                  <path d="M7 10h4" />
                  <path d="M7 14h7" />
                  <circle cx="17" cy="12" r="1" />
                </svg>
              </span>
              <h3>Identificacion</h3>
            </div>

            <div class="fields-grid">
              <div class="field field-wide">
                <label for="nombre">Nombre / Razon social</label>
                <div class="field-control">
                  <svg viewBox="0 0 24 24" class="field-icon">
                    <path d="M3 21h18" />
                    <path d="M5 21V7l8-4v18" />
                    <path d="M19 21V11l-6-4" />
                  </svg>
                  <InputText id="nombre" name="nombre" v-model="form.nombre" placeholder="Nombre o razon social" class="salah-input" />
                </div>
                <Message v-if="$form.nombre?.invalid" severity="error" size="small" variant="simple">
                  {{ $form.nombre.error?.message }}
                </Message>
              </div>

              <div class="field">
                <label for="ci_nit">CI/NIT</label>
                <div class="field-control">
                  <svg viewBox="0 0 24 24" class="field-icon">
                    <rect width="18" height="14" x="3" y="5" rx="2" />
                    <path d="M7 10h4" />
                    <path d="M7 14h7" />
                  </svg>
                  <InputText id="ci_nit" name="ci_nit" v-model="form.ci_nit" placeholder="Documento fiscal" class="salah-input" :class="{ 'p-invalid': ciNitExistente }" />
                </div>
                <Message v-if="$form.ci_nit?.invalid" severity="error" size="small" variant="simple">
                  {{ $form.ci_nit.error?.message }}
                </Message>
                <Message v-if="ciNitExistente" severity="error" size="small" variant="simple">
                  Este CI/NIT ya existe
                </Message>
              </div>

              <div class="field">
                <label for="tipo">Tipo</label>
                <div class="field-control">
                  <svg viewBox="0 0 24 24" class="field-icon">
                    <path d="M20 7h-9" />
                    <path d="M14 17H5" />
                    <circle cx="17" cy="17" r="3" />
                    <circle cx="7" cy="7" r="3" />
                  </svg>
                  <Select id="tipo" name="tipo" v-model="form.tipo" :options="tiposProveedor" placeholder="Seleccione tipo" show-clear class="salah-select" />
                </div>
                <Message v-if="$form.tipo?.invalid" severity="error" size="small" variant="simple">
                  {{ $form.tipo.error?.message }}
                </Message>
              </div>

              <div class="field field-wide">
                <label for="estado">Estado</label>
                <div class="field-control">
                  <svg viewBox="0 0 24 24" class="field-icon">
                    <path d="M12 2v20" />
                    <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7H14a3.5 3.5 0 0 1 0 7H6" />
                  </svg>
                  <Select id="estado" name="estado" v-model="form.estado" :options="estados" placeholder="Seleccione estado" class="salah-select" />
                </div>
                <Message v-if="$form.estado?.invalid" severity="error" size="small" variant="simple">
                  {{ $form.estado.error?.message }}
                </Message>
              </div>
            </div>
          </section>

          <section class="form-card contact-card">
            <div class="section-title">
              <span aria-hidden="true">
                <svg viewBox="0 0 24 24">
                  <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.8 19.8 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6A19.8 19.8 0 0 1 2.1 4.18 2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.12.9.32 1.77.59 2.61a2 2 0 0 1-.45 2.11L8 9.69a16 16 0 0 0 6.31 6.31l1.25-1.25a2 2 0 0 1 2.11-.45c.84.27 1.71.47 2.61.59A2 2 0 0 1 22 16.92Z" />
                </svg>
              </span>
              <h3>Contacto</h3>
            </div>

            <div class="fields-grid">
              <div class="field">
                <label for="telefono">Telefono</label>
                <div class="field-control">
                  <svg viewBox="0 0 24 24" class="field-icon">
                    <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.8 19.8 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6A19.8 19.8 0 0 1 2.1 4.18 2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.12.9.32 1.77.59 2.61a2 2 0 0 1-.45 2.11L8 9.69a16 16 0 0 0 6.31 6.31l1.25-1.25a2 2 0 0 1 2.11-.45c.84.27 1.71.47 2.61.59A2 2 0 0 1 22 16.92Z" />
                  </svg>
                  <InputText id="telefono" name="telefono" v-model="form.telefono" placeholder="Numero de telefono" maxlength="12" class="salah-input" />
                </div>
                <Message v-if="$form.telefono?.invalid" severity="error" size="small" variant="simple">
                  {{ $form.telefono.error?.message }}
                </Message>
              </div>

              <div class="field">
                <label for="email">Email</label>
                <div class="field-control">
                  <svg viewBox="0 0 24 24" class="field-icon">
                    <rect width="20" height="16" x="2" y="4" rx="2" />
                    <path d="m22 7-10 6L2 7" />
                  </svg>
                  <InputText id="email" name="email" v-model="form.email" placeholder="correo@dominio.com" class="salah-input" />
                </div>
                <Message v-if="$form.email?.invalid" severity="error" size="small" variant="simple">
                  {{ $form.email.error?.message }}
                </Message>
              </div>

              <div class="field field-wide">
                <label for="direccion">Direccion</label>
                <div class="field-control">
                  <svg viewBox="0 0 24 24" class="field-icon">
                    <path d="M20 10c0 5-8 12-8 12S4 15 4 10a8 8 0 1 1 16 0Z" />
                    <circle cx="12" cy="10" r="3" />
                  </svg>
                  <InputText id="direccion" name="direccion" v-model="form.direccion" placeholder="Direccion del proveedor" class="salah-input" />
                </div>
                <Message v-if="$form.direccion?.invalid" severity="error" size="small" variant="simple">
                  {{ $form.direccion.error?.message }}
                </Message>
              </div>
            </div>
          </section>

          <section class="form-card notes-card">
            <div class="section-title">
              <span aria-hidden="true">
                <svg viewBox="0 0 24 24">
                  <path d="M21 15a4 4 0 0 1-4 4H7l-4 4V7a4 4 0 0 1 4-4h10a4 4 0 0 1 4 4z" />
                </svg>
              </span>
              <h3>Observaciones</h3>
            </div>
            <div class="field">
              <label for="observaciones">Notas internas</label>
              <Textarea id="observaciones" name="observaciones" v-model="form.observaciones" rows="8" auto-resize class="salah-textarea" />
              <Message v-if="$form.observaciones?.invalid" severity="error" size="small" variant="simple">
                {{ $form.observaciones.error?.message }}
              </Message>
            </div>
          </section>

          <footer class="modal-actions">
            <Button type="submit" :label="form.id ? 'Guardar cambios' : 'Registrar proveedor'" class="salah-submit" :loading="saving" />
          </footer>
        </Form>
      </div>
    </Dialog>

    <Dialog v-model:visible="detalleVisible" modal header="Detalle del proveedor" :style="{ width: '78rem', maxWidth: '96vw' }">
      <div class="flex flex-col gap-4">
        <div class="grid grid-cols-1 gap-3 md:grid-cols-4">
          <div class="rounded-lg border border-gray-200 bg-white p-3">
            <p class="text-sm text-gray-500">Proveedor</p>
            <p class="font-semibold text-gray-900">{{ proveedorDetalle?.nombre || 'N/A' }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 bg-white p-3">
            <p class="text-sm text-gray-500">CI/NIT</p>
            <p class="font-semibold text-gray-900">{{ proveedorDetalle?.ci_nit || 'N/A' }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 bg-white p-3">
            <p class="text-sm text-gray-500">Compras</p>
            <p class="font-semibold text-gray-900">{{ historial.length }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 bg-white p-3">
            <p class="text-sm text-gray-500">Total comprado</p>
            <p class="font-semibold text-gray-900">$ {{ formatPrecio(proveedorDetalle?.total_comprado_usd) }}</p>
          </div>
        </div>

        <DataTable :value="historial" :loading="historialLoading" tableStyle="min-width: 62rem" size="small" stripedRows paginator :rows="5" :rowsPerPageOptions="[5, 10, 20]">
          <template #empty>
            <p class="p-4 text-center">No hay compras registradas para este proveedor.</p>
          </template>
          <Column field="fecha_compra" header="Fecha" sortable>
            <template #body="slotProps">{{ formatFecha(slotProps.data.fecha_compra) }}</template>
          </Column>
          <Column field="vehiculo" header="Vehiculo" sortable />
          <Column field="precio_compra_usd" header="Precio compra" sortable>
            <template #body="slotProps">
              <div>
                <p class="font-medium text-gray-900">$ {{ formatPrecio(slotProps.data.precio_compra_usd) }}</p>
                <p class="text-xs text-gray-500">Bs {{ formatPrecio(valorBOB(slotProps.data.precio_compra_bob, slotProps.data.precio_compra_usd, slotProps.data.tipo_cambio_usado)) }}</p>
              </div>
            </template>
          </Column>
          <Column field="gastos_adicionales" header="Gastos" sortable>
            <template #body="slotProps">
              <div>
                <p class="font-medium text-gray-900">$ {{ formatPrecio(slotProps.data.gastos_adicionales) }}</p>
                <p class="text-xs text-gray-500">Bs {{ formatPrecio(valorBOB(slotProps.data.gastos_adicionales_bob, slotProps.data.gastos_adicionales, slotProps.data.tipo_cambio_usado)) }}</p>
              </div>
            </template>
          </Column>
          <Column field="costo_total_usd" header="Costo total" sortable>
            <template #body="slotProps">
              <div>
                <p class="font-medium text-gray-900">$ {{ formatPrecio(slotProps.data.costo_total_usd) }}</p>
                <p class="text-xs text-gray-500">Bs {{ formatPrecio(valorBOB(slotProps.data.costo_total_bob, slotProps.data.costo_total_usd, slotProps.data.tipo_cambio_usado)) }}</p>
              </div>
            </template>
          </Column>
          <Column field="estado_pago" header="Estado de pago" sortable>
            <template #body="slotProps">
              <Tag :value="slotProps.data.estado_pago" :severity="slotProps.data.estado_pago === 'Pagado completo' ? 'success' : 'warning'" />
            </template>
          </Column>
        </DataTable>
      </div>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { zodResolver } from '@primevue/forms/resolvers/zod';
import { computed, onMounted, reactive, ref } from 'vue';
import { z } from 'zod';
import { server } from '~/server/server';
import Button from 'primevue/button';
import Column from 'primevue/column';
import DataTable from 'primevue/datatable';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Select from 'primevue/select';
import Tag from 'primevue/tag';
import Textarea from 'primevue/textarea';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';

definePageMeta({ layout: 'menu-admin' });

const toast = useToast();
const proveedores = ref<any[]>([]);
const historial = ref<any[]>([]);
const loading = ref(true);
const saving = ref(false);
const historialLoading = ref(false);
const formVisible = ref(false);
const detalleVisible = ref(false);
const proveedorDetalle = ref<any>(null);
const ciNitExistente = ref(false);
const searchQuery = ref('');
const topSearch = useTopSearch();
const tiposProveedor = ref(['Persona natural', 'Concesionaria', 'Importadora']);
const estados = ref(['Activo', 'Inactivo']);

const form = reactive({
  id: 0,
  nombre: '',
  ci_nit: '',
  telefono: '',
  email: '',
  direccion: '',
  tipo: '',
  observaciones: '',
  estado: 'Activo'
});

const resolver = ref(zodResolver(
  z.object({
    nombre: z.string().trim().min(1, { message: 'Nombre o razon social requerido.' }),
    ci_nit: z.string().trim().min(5, { message: 'CI/NIT debe tener al menos 5 caracteres.' }).max(30, { message: 'CI/NIT no debe superar 30 caracteres.' }),
    telefono: z.string().trim()
      .min(7, { message: 'Telefono debe tener al menos 7 digitos.' })
      .max(12, { message: 'Telefono no debe superar 12 caracteres.' })
      .regex(/^[0-9+\-\s]+$/, { message: 'Telefono solo puede contener numeros, +, - o espacios.' }),
    email: z.string().trim().email({ message: 'Email no valido.' }).optional().or(z.literal('')),
    direccion: z.string().trim().max(255, { message: 'Direccion no debe superar 255 caracteres.' }).optional(),
    tipo: z.enum(['Persona natural', 'Concesionaria', 'Importadora']).optional().or(z.literal('')),
    observaciones: z.string().trim().max(500, { message: 'Observaciones no debe superar 500 caracteres.' }).optional(),
    estado: z.enum(['Activo', 'Inactivo'], { message: 'Estado requerido.' })
  })
));

const activeSearchQuery = computed(() => (topSearch.value || searchQuery.value).trim().toLowerCase());
const filteredProveedores = computed(() => {
  if (!activeSearchQuery.value) return proveedores.value;
  const query = activeSearchQuery.value;
  return proveedores.value.filter((proveedor: any) =>
    (proveedor.nombre?.toLowerCase() || '').includes(query) ||
    (proveedor.ci_nit?.toLowerCase() || '').includes(query) ||
    (proveedor.telefono?.toLowerCase() || '').includes(query) ||
    (proveedor.tipo?.toLowerCase() || '').includes(query) ||
    (proveedor.estado?.toLowerCase() || '').includes(query)
  );
});
const proveedoresActivos = computed(() => filteredProveedores.value.filter((proveedor: any) => proveedor.estado === 'Activo').length);
const totalCompradoUSD = computed(() => filteredProveedores.value.reduce((total: number, proveedor: any) => total + Number(proveedor.total_comprado_usd || 0), 0));
const totalCompradoBOB = computed(() => filteredProveedores.value.reduce((total: number, proveedor: any) => total + Number(proveedor.total_comprado_bob || 0), 0));

onMounted(async () => {
  await obtenerProveedores();
});

async function obtenerProveedores() {
  loading.value = true;
  try {
    const res = await $fetch(server.HOST + '/api/v1/proveedores-autos', { method: 'GET' });
    proveedores.value = Array.isArray(res) ? res : [];
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar proveedores', life: 3000 });
  } finally {
    loading.value = false;
  }
}

function abrirNuevo() {
  resetForm();
  ciNitExistente.value = false;
  formVisible.value = true;
}

function abrirEditar(proveedor: any) {
  Object.assign(form, {
    id: proveedor.id,
    nombre: proveedor.nombre || '',
    ci_nit: proveedor.ci_nit || '',
    telefono: proveedor.telefono || '',
    email: proveedor.email || '',
    direccion: proveedor.direccion || '',
    tipo: proveedor.tipo || '',
    observaciones: proveedor.observaciones || '',
    estado: proveedor.estado || 'Activo'
  });
  ciNitExistente.value = false;
  formVisible.value = true;
}

async function abrirDetalle(proveedor: any) {
  proveedorDetalle.value = proveedor;
  detalleVisible.value = true;
  historialLoading.value = true;
  try {
    const res = await $fetch(server.HOST + '/api/v1/proveedores-autos/' + proveedor.id + '/historial-compras', { method: 'GET' });
    historial.value = Array.isArray(res) ? res : [];
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar historial', life: 3000 });
  } finally {
    historialLoading.value = false;
  }
}

async function guardarProveedor({ valid }: any) {
  if (!valid) return;
  const ciNit = form.ci_nit.trim().toLowerCase();
  const duplicado = proveedores.value.find((proveedor: any) =>
    String(proveedor.ci_nit || '').trim().toLowerCase() === ciNit && Number(proveedor.id) !== Number(form.id)
  );
  ciNitExistente.value = !!duplicado;
  if (duplicado) return;

  saving.value = true;
  try {
    const method = form.id ? 'PUT' : 'POST';
    const url = form.id ? `${server.HOST}/api/v1/proveedores-autos/${form.id}` : `${server.HOST}/api/v1/proveedores-autos`;
    await $fetch(url, {
      method,
      body: {
        nombre: form.nombre,
        ci_nit: form.ci_nit,
        telefono: form.telefono,
        email: form.email,
        direccion: form.direccion,
        tipo: form.tipo,
        observaciones: form.observaciones,
        estado: form.estado
      }
    });
    toast.add({ severity: 'success', summary: form.id ? 'Proveedor actualizado' : 'Proveedor registrado', life: 3000 });
    formVisible.value = false;
    await obtenerProveedores();
  } catch (err: any) {
    console.error(err);
    if (err?.response?.status === 409 || err?.statusCode === 409) {
      ciNitExistente.value = true;
      return;
    }
    toast.add({ severity: 'error', summary: 'Error al guardar proveedor', detail: err?.data || err?.message, life: 4000 });
  } finally {
    saving.value = false;
  }
}

function resetForm() {
  Object.assign(form, {
    id: 0,
    nombre: '',
    ci_nit: '',
    telefono: '',
    email: '',
    direccion: '',
    tipo: '',
    observaciones: '',
    estado: 'Activo'
  });
}

function formatPrecio(precio: number) {
  return Number(precio || 0).toLocaleString('es-BO', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

function formatFecha(fecha: string) {
  if (!fecha) return 'N/A';
  return new Date(`${fecha}T00:00:00`).toLocaleDateString('es-BO');
}

function valorBOB(valorGuardado: number, valorUSD: number, tc: number) {
  const guardado = Number(valorGuardado || 0);
  if (guardado > 0) return guardado;
  return Math.round(Number(valorUSD || 0) * Number(tc || 0) * 100) / 100;
}
</script>

<style scoped>
:deep(.salah-dialog-content) {
  padding: 0;
  border-radius: 20px;
  overflow: hidden;
  background: #ffffff;
}

:deep(.p-dialog) {
  border-radius: 20px;
  box-shadow: 0 24px 80px rgba(13, 13, 13, 0.32);
}

.salah-user-modal {
  background: #f7f7f7;
  color: #0d0d0d;
}

.modal-header {
  position: relative;
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px 28px;
  background: linear-gradient(135deg, #0d0d0d 0%, #202020 100%);
  color: #ffffff;
}

.header-icon,
.section-title span {
  display: grid;
  place-items: center;
  flex: 0 0 auto;
}

.header-icon {
  width: 52px;
  height: 52px;
  border-radius: 16px;
  background: #ffd700;
  color: #0d0d0d;
  box-shadow: 0 12px 30px rgba(255, 215, 0, 0.24);
}

svg {
  width: 22px;
  height: 22px;
  fill: none;
  stroke: currentColor;
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.header-copy h2 {
  margin: 0;
  font-size: 1.55rem;
  font-weight: 800;
  letter-spacing: 0;
}

.header-copy p {
  margin: 2px 0 10px;
  color: rgba(255, 255, 255, 0.72);
  font-size: 0.92rem;
}

.header-copy span {
  display: block;
  width: 96px;
  height: 3px;
  border-radius: 999px;
  background: #ffd700;
}

.close-button {
  margin-left: auto;
  width: 42px;
  height: 42px;
  border: 1px solid rgba(255, 255, 255, 0.14);
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.08);
  color: #ffffff;
  cursor: pointer;
  transition: background-color 0.2s ease, transform 0.2s ease;
}

.close-button:hover {
  background: rgba(255, 255, 255, 0.16);
  transform: translateY(-1px);
}

.close-button svg {
  margin: auto;
}

.user-form {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 320px;
  gap: 18px;
  padding: 20px;
}

.form-card {
  border: 1px solid rgba(13, 13, 13, 0.08);
  border-radius: 18px;
  background: #ffffff;
  box-shadow: 0 18px 45px rgba(13, 13, 13, 0.08);
}

.identity-card,
.contact-card {
  grid-column: 1;
  padding: 18px;
}

.notes-card {
  grid-column: 2;
  grid-row: 1 / span 2;
  padding: 18px;
  align-self: stretch;
}

.modal-actions {
  grid-column: 1;
  display: flex;
  justify-content: flex-end;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
}

.section-title span {
  width: 36px;
  height: 36px;
  border-radius: 12px;
  background: rgba(255, 215, 0, 0.18);
  color: #0d0d0d;
}

.section-title h3 {
  margin: 0;
  font-size: 1rem;
  font-weight: 800;
  letter-spacing: 0;
}

.fields-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
}

.field {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 6px;
}

.field-wide {
  grid-column: 1 / -1;
}

.field label {
  color: #2a2a2a;
  font-size: 0.79rem;
  font-weight: 700;
  text-transform: uppercase;
}

.field-control {
  position: relative;
  display: flex;
  align-items: center;
}

.field-icon {
  position: absolute;
  left: 14px;
  z-index: 1;
  width: 19px;
  height: 19px;
  color: #6c6c6c;
  pointer-events: none;
}

.salah-input,
:deep(.salah-select) {
  width: 100%;
}

:deep(.salah-input),
:deep(.salah-select .p-select-label) {
  min-height: 50px;
}

:deep(.salah-input),
:deep(.salah-select),
:deep(.salah-textarea) {
  border: 1px solid #d8d8d8;
  border-radius: 12px;
  background: #ffffff;
  color: #0d0d0d;
  box-shadow: none;
  transition: border-color 0.2s ease, box-shadow 0.2s ease, background-color 0.2s ease;
}

:deep(.salah-input) {
  padding-left: 44px;
}

:deep(.salah-select) {
  min-height: 50px;
  padding-left: 34px;
}

:deep(.salah-select .p-select-label) {
  display: flex;
  align-items: center;
}

:deep(.salah-textarea) {
  width: 100%;
  min-height: 236px;
  resize: vertical;
}

:deep(.salah-input:enabled:focus),
:deep(.salah-select:not(.p-disabled).p-focus),
:deep(.salah-textarea:enabled:focus) {
  border-color: #ffd700;
  box-shadow: 0 0 0 4px rgba(255, 215, 0, 0.18);
}

:deep(.p-invalid) {
  border-color: #e30613;
}

:deep(.salah-submit) {
  min-height: 50px;
  min-width: 210px;
  border: 0;
  border-radius: 12px;
  background: #ffd700;
  color: #0d0d0d;
  font-weight: 800;
  box-shadow: 0 14px 28px rgba(255, 215, 0, 0.26);
  transition: background-color 0.2s ease, box-shadow 0.2s ease, transform 0.2s ease;
}

:deep(.salah-submit:hover) {
  background: #e6c200;
  color: #0d0d0d;
  box-shadow: 0 18px 34px rgba(255, 215, 0, 0.32);
  transform: translateY(-1px);
}

:deep(.p-message-text) {
  font-size: 0.78rem;
}

@media (max-width: 860px) {
  .modal-header {
    padding: 20px;
  }

  .user-form {
    grid-template-columns: 1fr;
    padding: 14px;
  }

  .identity-card,
  .contact-card,
  .notes-card,
  .modal-actions {
    grid-column: 1;
  }

  .notes-card {
    grid-row: auto;
  }

  .fields-grid {
    grid-template-columns: 1fr;
  }

  .modal-actions {
    justify-content: stretch;
  }

  :deep(.salah-submit) {
    width: 100%;
  }
}

@media (max-width: 520px) {
  .modal-header {
    align-items: flex-start;
  }

  .header-icon {
    width: 46px;
    height: 46px;
    border-radius: 14px;
  }

  .header-copy h2 {
    font-size: 1.3rem;
  }

  .close-button {
    position: absolute;
    top: 14px;
    right: 14px;
  }
}
</style>
