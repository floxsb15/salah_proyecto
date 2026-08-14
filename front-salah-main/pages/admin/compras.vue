<template>
  <Toast />

  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <h2 class="text-2xl font-bold">Compras de Autos</h2>
        <p class="text-sm text-gray-500">Registro de ingreso de vehiculos al inventario y costo real.</p>
      </div>

      <div class="flex flex-col gap-2 sm:flex-row sm:items-center">
        <Button label="Nueva compra" icon="pi pi-plus" size="small" @click="abrirCompra" />
        <span class="p-input-icon-left">
          <i class="pi pi-search" />
          <InputText v-model="searchQuery" placeholder="Buscar..." size="small" />
        </span>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-3 md:grid-cols-3">
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Compras</p>
        <p class="text-2xl font-bold text-gray-900">{{ filteredCompras.length }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Egresos USD</p>
        <p class="text-2xl font-bold text-gray-900">$ {{ formatPrecio(totalCostoUSD) }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 bg-white p-4">
        <p class="text-sm text-gray-500">Pendientes</p>
        <p class="text-2xl font-bold text-gray-900">{{ comprasPendientes }}</p>
      </div>
    </div>

    <DataTable
      :value="filteredCompras"
      :loading="loading"
      tableStyle="min-width: 78rem"
      size="small"
      stripedRows
      removableSort
      paginator
      :rows="10"
      :rowsPerPageOptions="[10, 20, 50]"
    >
      <template #empty>
        <p class="p-4 text-center">No hay compras registradas.</p>
      </template>

      <Column field="fecha_compra" header="Fecha" sortable>
        <template #body="slotProps">{{ formatFecha(slotProps.data.fecha_compra) }}</template>
      </Column>
      <Column field="vehiculo" header="Vehiculo" sortable />
      <Column field="proveedor" header="Proveedor" sortable>
        <template #body="slotProps">
          <div>
            <p class="font-medium text-gray-900">{{ slotProps.data.proveedor }}</p>
            <p class="text-xs text-gray-500">CI/NIT: {{ slotProps.data.proveedor_ci_nit || 'N/A' }}</p>
            <p class="text-xs text-gray-500">Tel: {{ slotProps.data.proveedor_telefono || 'N/A' }}</p>
          </div>
        </template>
      </Column>
      <Column field="precio_compra_usd" header="Precio compra" sortable>
        <template #body="slotProps">
          <div>
            <p class="font-medium text-gray-900">$ {{ formatPrecio(slotProps.data.precio_compra_usd) }}</p>
            <p class="text-xs text-gray-500">Bs {{ formatPrecio(valorBOB(slotProps.data.precio_compra_bob, slotProps.data.precio_compra_usd, slotProps.data.tipo_cambio_usado)) }}</p>
            <p class="text-xs text-gray-400">Registrado en {{ slotProps.data.moneda_precio || 'USD' }}</p>
          </div>
        </template>
      </Column>
      <Column field="gastos_adicionales" header="Gastos adicionales" sortable>
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
      <Column field="tipo_cambio_usado" header="TC usado" sortable>
        <template #body="slotProps">{{ formatTipoCambio(slotProps.data.tipo_cambio_usado) }}</template>
      </Column>
      <Column field="metodo_pago" header="Metodo" sortable />
      <Column field="detalle_pago" header="Detalle pago">
        <template #body="slotProps">{{ slotProps.data.detalle_pago || slotProps.data.metodo_pago || 'N/A' }}</template>
      </Column>
      <Column header="Acciones" style="width: 12rem">
        <template #body="slotProps">
          <Button
            v-if="slotProps.data.estado_pago === 'Pendiente'"
            label="Completar pago"
            icon="pi pi-check-circle"
            size="small"
            severity="success"
            variant="text"
            @click="abrirCompletarPago(slotProps.data)"
          />
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="compraVisible" modal header="Nueva compra" :style="{ width: '62rem' }">
      <form class="grid max-h-[78vh] grid-cols-1 gap-4 overflow-y-auto pr-1 lg:grid-cols-2" @submit.prevent="registrarCompra">
        <div class="flex flex-col gap-1 lg:col-span-2">
          <label for="vehiculo">Vehiculo</label>
          <div class="grid grid-cols-1 gap-2 md:grid-cols-[1fr_auto]">
            <Select id="vehiculo" v-model="form.id_vehiculo" :options="vehiculos" option-label="nombreCompleto" option-value="id" placeholder="Seleccione vehiculo" filter fluid size="small" />
            <Button label="Crear vehiculo" icon="pi pi-car" severity="secondary" type="button" size="small" @click="vehiculoVisible = true" />
          </div>
        </div>

        <div class="flex flex-col gap-1 lg:col-span-2">
          <label for="proveedor">Proveedor</label>
          <div class="grid grid-cols-1 gap-2 md:grid-cols-[1fr_auto]">
            <Select id="proveedor" v-model="form.id_proveedor" :options="proveedoresActivos" option-label="nombreCompleto" option-value="id" placeholder="Seleccione proveedor" filter fluid size="small" />
            <Button label="Crear proveedor" icon="pi pi-briefcase" severity="secondary" type="button" size="small" @click="abrirProveedorRapido" />
          </div>
        </div>

        <div class="flex flex-col gap-1">
          <label for="fecha_compra">Fecha de compra</label>
          <InputText id="fecha_compra" v-model="form.fecha_compra" type="date" size="small" />
        </div>
        <div class="flex flex-col gap-1">
          <label for="tipo_cambio">Tipo de cambio</label>
          <InputNumber id="tipo_cambio" v-model="form.tipo_cambio" :min="0" :minFractionDigits="2" :maxFractionDigits="4" suffix=" Bs/USD" fluid size="small" />
        </div>
        <div class="flex flex-col gap-1">
          <label for="moneda_precio">Moneda precio</label>
          <Select id="moneda_precio" v-model="form.moneda_precio" :options="monedasPago" fluid size="small" />
        </div>
        <div class="flex flex-col gap-1">
          <label for="precio_compra">Precio de compra {{ form.moneda_precio }}</label>
          <InputNumber id="precio_compra" v-model="form.precio_compra" mode="currency" :currency="form.moneda_precio" locale="es-BO" :min="0" fluid size="small" />
          <small class="text-gray-500">$ {{ formatPrecio(precioCompraUSD) }} / Bs {{ formatPrecio(precioCompraBOB) }}</small>
        </div>

        <section class="grid grid-cols-1 gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2 md:grid-cols-3">
          <div class="flex flex-col gap-1">
            <label for="gasto_importacion">Importacion</label>
            <InputNumber id="gasto_importacion" v-model="form.gasto_importacion" mode="currency" currency="USD" locale="es-BO" :min="0" fluid size="small" />
            <small class="text-gray-500">Bs {{ formatPrecio(gastoImportacionBOB) }}</small>
          </div>
          <div class="flex flex-col gap-1">
            <label for="gasto_transporte">Transporte</label>
            <InputNumber id="gasto_transporte" v-model="form.gasto_transporte" mode="currency" currency="USD" locale="es-BO" :min="0" fluid size="small" />
            <small class="text-gray-500">Bs {{ formatPrecio(gastoTransporteBOB) }}</small>
          </div>
          <div class="flex flex-col gap-1">
            <label for="gasto_papeleo">Papeleo</label>
            <InputNumber id="gasto_papeleo" v-model="form.gasto_papeleo" mode="currency" currency="USD" locale="es-BO" :min="0" fluid size="small" />
            <small class="text-gray-500">Bs {{ formatPrecio(gastoPapeleoBOB) }}</small>
          </div>
        </section>

        <div class="grid grid-cols-1 gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2 md:grid-cols-4">
          <div>
            <span class="text-xs font-semibold uppercase text-gray-500">Precio USD</span>
            <strong class="block text-gray-900">$ {{ formatPrecio(precioCompraUSD) }}</strong>
            <small class="text-gray-500">Bs {{ formatPrecio(precioCompraBOB) }}</small>
          </div>
          <div>
            <span class="text-xs font-semibold uppercase text-gray-500">Gastos USD</span>
            <strong class="block text-gray-900">$ {{ formatPrecio(gastosAdicionales) }}</strong>
            <small class="text-gray-500">Bs {{ formatPrecio(gastosAdicionalesBOB) }}</small>
          </div>
          <div>
            <span class="text-xs font-semibold uppercase text-gray-500">Costo total USD</span>
            <strong class="block text-gray-900">$ {{ formatPrecio(costoTotalUSD) }}</strong>
            <small class="text-gray-500">Bs {{ formatPrecio(costoTotalBOB) }}</small>
          </div>
          <div>
            <span class="text-xs font-semibold uppercase text-gray-500">Estado sugerido</span>
            <strong class="block text-gray-900">{{ estadoPagoSugerido }}</strong>
          </div>
        </div>

        <div class="flex flex-col gap-1">
          <label for="metodo_pago">Metodo de pago de compra</label>
          <Select id="metodo_pago" v-model="form.metodo_pago" :options="metodosCompra" fluid size="small" />
        </div>

        <section class="flex flex-col gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2">
          <div class="flex items-center justify-between gap-3">
            <h3 class="text-sm font-semibold text-gray-900">Detalle de pago</h3>
            <Button label="Agregar pago" icon="pi pi-plus" size="small" severity="secondary" type="button" @click="agregarPago" />
          </div>

          <div v-for="(pago, index) in form.pagos" :key="pago.key" class="grid grid-cols-1 gap-2 md:grid-cols-[9rem_13rem_1fr_2.5rem]">
            <Select v-model="pago.moneda" :options="monedasPago" placeholder="Moneda" fluid size="small" />
            <Select v-model="pago.metodo" :options="metodosPagoMoneda" placeholder="Metodo" fluid size="small" />
            <InputNumber v-model="pago.monto" mode="decimal" :min="0" :minFractionDigits="2" :maxFractionDigits="2" placeholder="Monto" fluid size="small" />
            <Button icon="pi pi-trash" severity="danger" text rounded type="button" aria-label="Eliminar pago" @click="eliminarPago(index)" />
          </div>

          <p v-if="pagoExcedeTotal" class="text-sm font-semibold text-red-600">El pagado supera el costo total.</p>
        </section>

        <div class="flex flex-col gap-1 lg:col-span-2">
          <label for="observacion">Observacion</label>
          <Textarea id="observacion" v-model="form.observacion" rows="3" auto-resize fluid />
        </div>

        <div class="flex justify-end gap-2 border-t border-gray-100 pt-4 lg:col-span-2">
          <Button label="Cancelar" severity="secondary" type="button" @click="compraVisible = false" />
          <Button label="Registrar compra" icon="pi pi-check" type="submit" :loading="saving" />
        </div>
      </form>
    </Dialog>

    <modalAgregarVehiculo
      v-if="vehiculoVisible"
      :open="vehiculoVisible"
      @close="vehiculoVisible = false"
      @update="obtenerVehiculos"
      @success="toast.add({ severity: 'success', summary: 'Vehiculo agregado', life: 3000 })"
      @error="mostrarError('Error al agregar vehiculo', $event)"
    />

    <Dialog
      v-model:visible="proveedorVisible"
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
            <h2>Nuevo Proveedor</h2>
            <p>Registro para compras de autos</p>
            <span></span>
          </div>
          <button type="button" class="close-button" aria-label="Cerrar" @click="proveedorVisible = false">
            <svg viewBox="0 0 24 24">
              <path d="M18 6 6 18" />
              <path d="m6 6 12 12" />
            </svg>
          </button>
        </header>

        <Form v-slot="$form" :resolver="proveedorResolver" :initialValues="proveedorForm" @submit="guardarProveedorRapido" class="user-form">
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
                <label for="proveedor_nombre">Nombre / Razon social</label>
                <div class="field-control">
                  <svg viewBox="0 0 24 24" class="field-icon">
                    <path d="M3 21h18" />
                    <path d="M5 21V7l8-4v18" />
                    <path d="M19 21V11l-6-4" />
                  </svg>
                  <InputText id="proveedor_nombre" name="nombre" v-model="proveedorForm.nombre" placeholder="Nombre o razon social" class="salah-input" />
                </div>
                <Message v-if="$form.nombre?.invalid" severity="error" size="small" variant="simple">
                  {{ $form.nombre.error?.message }}
                </Message>
              </div>

              <div class="field">
                <label for="proveedor_ci">CI/NIT</label>
                <div class="field-control">
                  <svg viewBox="0 0 24 24" class="field-icon">
                    <rect width="18" height="14" x="3" y="5" rx="2" />
                    <path d="M7 10h4" />
                    <path d="M7 14h7" />
                  </svg>
                  <InputText id="proveedor_ci" name="ci_nit" v-model="proveedorForm.ci_nit" placeholder="Documento fiscal" class="salah-input" :class="{ 'p-invalid': ciNitProveedorExistente }" />
                </div>
                <Message v-if="$form.ci_nit?.invalid" severity="error" size="small" variant="simple">
                  {{ $form.ci_nit.error?.message }}
                </Message>
                <Message v-if="ciNitProveedorExistente" severity="error" size="small" variant="simple">
                  Este CI/NIT ya existe
                </Message>
              </div>

              <div class="field">
                <label for="proveedor_tipo">Tipo</label>
                <div class="field-control">
                  <svg viewBox="0 0 24 24" class="field-icon">
                    <path d="M20 7h-9" />
                    <path d="M14 17H5" />
                    <circle cx="17" cy="17" r="3" />
                    <circle cx="7" cy="7" r="3" />
                  </svg>
                  <Select id="proveedor_tipo" name="tipo" v-model="proveedorForm.tipo" :options="tiposProveedor" placeholder="Seleccione tipo" show-clear class="salah-select" />
                </div>
                <Message v-if="$form.tipo?.invalid" severity="error" size="small" variant="simple">
                  {{ $form.tipo.error?.message }}
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
                <label for="proveedor_tel">Telefono</label>
                <div class="field-control">
                  <svg viewBox="0 0 24 24" class="field-icon">
                    <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.8 19.8 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6A19.8 19.8 0 0 1 2.1 4.18 2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.12.9.32 1.77.59 2.61a2 2 0 0 1-.45 2.11L8 9.69a16 16 0 0 0 6.31 6.31l1.25-1.25a2 2 0 0 1 2.11-.45c.84.27 1.71.47 2.61.59A2 2 0 0 1 22 16.92Z" />
                  </svg>
                  <InputText id="proveedor_tel" name="telefono" v-model="proveedorForm.telefono" placeholder="Numero de telefono" maxlength="12" class="salah-input" />
                </div>
                <Message v-if="$form.telefono?.invalid" severity="error" size="small" variant="simple">
                  {{ $form.telefono.error?.message }}
                </Message>
              </div>

              <div class="field">
                <label for="proveedor_email">Email</label>
                <div class="field-control">
                  <svg viewBox="0 0 24 24" class="field-icon">
                    <rect width="20" height="16" x="2" y="4" rx="2" />
                    <path d="m22 7-10 6L2 7" />
                  </svg>
                  <InputText id="proveedor_email" name="email" v-model="proveedorForm.email" placeholder="correo@dominio.com" class="salah-input" />
                </div>
                <Message v-if="$form.email?.invalid" severity="error" size="small" variant="simple">
                  {{ $form.email.error?.message }}
                </Message>
              </div>

              <div class="field field-wide">
                <label for="proveedor_direccion">Direccion</label>
                <div class="field-control">
                  <svg viewBox="0 0 24 24" class="field-icon">
                    <path d="M20 10c0 5-8 12-8 12S4 15 4 10a8 8 0 1 1 16 0Z" />
                    <circle cx="12" cy="10" r="3" />
                  </svg>
                  <InputText id="proveedor_direccion" name="direccion" v-model="proveedorForm.direccion" placeholder="Direccion del proveedor" class="salah-input" />
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
              <label for="proveedor_observaciones">Notas internas</label>
              <Textarea id="proveedor_observaciones" name="observaciones" v-model="proveedorForm.observaciones" rows="8" auto-resize class="salah-textarea" />
              <Message v-if="$form.observaciones?.invalid" severity="error" size="small" variant="simple">
                {{ $form.observaciones.error?.message }}
              </Message>
            </div>
          </section>

          <footer class="modal-actions">
            <Button type="submit" label="Registrar proveedor" class="salah-submit" :loading="savingProveedor" />
          </footer>
        </Form>
      </div>
    </Dialog>

    <Dialog v-model:visible="completarPagoVisible" modal header="Completar pago de compra" :style="{ width: '58rem', maxWidth: '96vw' }">
      <div v-if="compraPagoSeleccionada" class="flex flex-col gap-4">
        <section class="rounded-lg border border-gray-200 bg-white p-4">
          <h3 class="text-lg font-semibold text-gray-900">{{ compraPagoSeleccionada.vehiculo }}</h3>
          <p class="text-sm text-gray-500">
            {{ compraPagoSeleccionada.proveedor }} - CI/NIT: {{ compraPagoSeleccionada.proveedor_ci_nit || 'N/A' }}
          </p>

          <div class="mt-4 grid grid-cols-1 gap-3 text-sm md:grid-cols-4">
            <div class="rounded-md bg-gray-50 p-3">
              <span class="text-gray-500">Precio compra</span>
              <strong class="block text-gray-900">$ {{ formatPrecio(compraPagoSeleccionada.precio_compra_usd) }}</strong>
              <small class="block text-gray-500">Bs {{ formatPrecio(valorBOB(compraPagoSeleccionada.precio_compra_bob, compraPagoSeleccionada.precio_compra_usd, compraPagoSeleccionada.tipo_cambio_usado)) }}</small>
            </div>
            <div class="rounded-md bg-gray-50 p-3">
              <span class="text-gray-500">Gastos adicionales</span>
              <strong class="block text-gray-900">$ {{ formatPrecio(compraPagoSeleccionada.gastos_adicionales) }}</strong>
              <small class="block text-gray-500">Bs {{ formatPrecio(valorBOB(compraPagoSeleccionada.gastos_adicionales_bob, compraPagoSeleccionada.gastos_adicionales, compraPagoSeleccionada.tipo_cambio_usado)) }}</small>
            </div>
            <div class="rounded-md bg-gray-50 p-3">
              <span class="text-gray-500">Costo total</span>
              <strong class="block text-gray-900">$ {{ formatPrecio(compraPagoSeleccionada.costo_total_usd) }}</strong>
              <small class="block text-gray-500">Bs {{ formatPrecio(valorBOB(compraPagoSeleccionada.costo_total_bob, compraPagoSeleccionada.costo_total_usd, compraPagoSeleccionada.tipo_cambio_usado)) }}</small>
            </div>
            <div class="rounded-md bg-gray-50 p-3">
              <span class="text-gray-500">TC usado</span>
              <strong class="block text-gray-900">{{ formatTipoCambio(compraPagoSeleccionada.tipo_cambio_usado) }}</strong>
            </div>
            <div class="rounded-md bg-gray-50 p-3">
              <span class="text-gray-500">Pagado USD</span>
              <strong class="block text-gray-900">$ {{ formatPrecio(pagadoCompraUSD) }}</strong>
            </div>
            <div class="rounded-md bg-gray-50 p-3">
              <span class="text-gray-500">Pagado BOB</span>
              <strong class="block text-gray-900">Bs {{ formatPrecio(pagadoCompraBOB) }}</strong>
            </div>
            <div class="rounded-md bg-gray-50 p-3">
              <span class="text-gray-500">Saldo USD</span>
              <strong class="block text-gray-900">$ {{ formatPrecio(saldoCompraUSD) }}</strong>
            </div>
            <div class="rounded-md bg-gray-50 p-3">
              <span class="text-gray-500">Saldo BOB</span>
              <strong class="block text-gray-900">Bs {{ formatPrecio(saldoCompraBOB) }}</strong>
              <small class="block text-gray-500">TC {{ formatTipoCambio(compraPagoSeleccionada.tipo_cambio_usado) }}</small>
            </div>
          </div>
        </section>

        <section class="flex flex-col gap-3 rounded-lg border border-gray-200 bg-white p-4">
          <div class="flex items-center justify-between gap-3">
            <div>
              <h3 class="text-sm font-semibold text-gray-900">Detalle de pago restante</h3>
              <p class="text-xs text-gray-500">Nuevo pago equivalente: $ {{ formatPrecio(nuevoPagoCompraUSD) }}</p>
            </div>
            <Button label="Agregar pago" icon="pi pi-plus" size="small" severity="secondary" type="button" @click="agregarPagoCompra" />
          </div>

          <div v-for="(pago, index) in completarPagoForm.pagos" :key="pago.key" class="grid grid-cols-1 gap-2 md:grid-cols-[9rem_13rem_1fr_2.5rem]">
            <Select v-model="pago.moneda" :options="monedasPago" placeholder="Moneda" fluid size="small" />
            <Select v-model="pago.metodo" :options="metodosPagoMoneda" placeholder="Metodo" fluid size="small" />
            <InputNumber v-model="pago.monto" mode="decimal" :min="0" :minFractionDigits="2" :maxFractionDigits="2" placeholder="Monto" fluid size="small" />
            <Button icon="pi pi-trash" severity="danger" text rounded type="button" aria-label="Eliminar pago" @click="eliminarPagoCompra(index)" />
          </div>

          <p v-if="pagoCompraExcedeSaldo" class="text-sm font-semibold text-red-600">El pago supera el saldo pendiente.</p>

          <div class="flex justify-end gap-2 border-t border-gray-100 pt-4">
            <Button label="Cancelar" severity="secondary" type="button" @click="completarPagoVisible = false" />
            <Button label="Completar pago" icon="pi pi-check" type="button" :loading="savingPagoCompra" @click="guardarPagoCompra" />
          </div>
        </section>
      </div>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { zodResolver } from '@primevue/forms/resolvers/zod';
import { computed, onMounted, reactive, ref } from 'vue';
import { z } from 'zod';
import { server } from '~/server/server';
import modalAgregarVehiculo from '~/components/admin/vehiculos/modalAgregarProducto.vue';
import Button from 'primevue/button';
import Column from 'primevue/column';
import DataTable from 'primevue/datatable';
import Dialog from 'primevue/dialog';
import InputNumber from 'primevue/inputnumber';
import InputText from 'primevue/inputtext';
import Message from 'primevue/message';
import Select from 'primevue/select';
import Tag from 'primevue/tag';
import Textarea from 'primevue/textarea';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';

definePageMeta({ layout: 'menu-admin' });

const toast = useToast();
const loading = ref(true);
const saving = ref(false);
const compraVisible = ref(false);
const vehiculoVisible = ref(false);
const proveedorVisible = ref(false);
const completarPagoVisible = ref(false);
const compraPagoSeleccionada = ref<any>(null);
const compras = ref<any[]>([]);
const vehiculos = ref<any[]>([]);
const proveedores = ref<any[]>([]);
const searchQuery = ref('');
const metodosCompra = ref(['Contado', 'Credito', 'Mixto']);
const metodosPagoMoneda = ref(['QR', 'Transferencia', 'Efectivo', 'Tarjeta']);
const monedasPago = ref(['USD', 'BOB']);
const tiposProveedor = ref(['Persona natural', 'Concesionaria', 'Importadora']);
const savingProveedor = ref(false);
const savingPagoCompra = ref(false);
const ciNitProveedorExistente = ref(false);

const form = reactive({
  id_vehiculo: null as number | null,
  id_proveedor: null as number | null,
  fecha_compra: new Date().toISOString().slice(0, 10),
  moneda_precio: 'USD',
  precio_compra: 0,
  tipo_cambio: null as number | null,
  gasto_importacion: 0,
  gasto_transporte: 0,
  gasto_papeleo: 0,
  metodo_pago: 'Contado',
  pagos: [crearPago()],
  observacion: ''
});
const proveedorForm = reactive({
  nombre: '',
  ci_nit: '',
  telefono: '',
  email: '',
  direccion: '',
  tipo: '',
  observaciones: ''
});
const completarPagoForm = reactive({
  pagos: [crearPago()]
});
const proveedorResolver = ref(zodResolver(
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
    observaciones: z.string().trim().max(500, { message: 'Observaciones no debe superar 500 caracteres.' }).optional()
  })
));

const filteredCompras = computed(() => {
  const query = searchQuery.value.trim().toLowerCase();
  if (!query) return compras.value;
  return compras.value.filter((compra: any) =>
    (compra.fecha_compra?.toLowerCase() || '').includes(query) ||
    (compra.vehiculo?.toLowerCase() || '').includes(query) ||
    (compra.proveedor?.toLowerCase() || '').includes(query) ||
    (compra.proveedor_ci_nit?.toLowerCase() || '').includes(query) ||
    (compra.estado_pago?.toLowerCase() || '').includes(query) ||
    (compra.metodo_pago?.toLowerCase() || '').includes(query) ||
    (compra.detalle_pago?.toLowerCase() || '').includes(query)
  );
});
const totalCostoUSD = computed(() => filteredCompras.value.reduce((total: number, compra: any) => total + Number(compra.costo_total_usd || 0), 0));
const comprasPendientes = computed(() => filteredCompras.value.filter((compra: any) => compra.estado_pago === 'Pendiente').length);
const proveedoresActivos = computed(() => proveedores.value
  .filter((proveedor: any) => proveedor.estado === 'Activo')
  .map((proveedor: any) => ({
    ...proveedor,
    nombreCompleto: `${proveedor.nombre || ''} - ${proveedor.ci_nit || 'Sin CI/NIT'}`.trim()
  })));
const tipoCambio = computed(() => Number(form.tipo_cambio || 0));
const precioCompraUSD = computed(() => {
  const precio = Number(form.precio_compra || 0);
  if (form.moneda_precio === 'BOB') {
    return tipoCambio.value > 0 ? roundMoney(precio / tipoCambio.value) : 0;
  }
  return precio;
});
const precioCompraBOB = computed(() => {
  const precio = Number(form.precio_compra || 0);
  if (form.moneda_precio === 'BOB') {
    return precio;
  }
  return roundMoney(precioCompraUSD.value * tipoCambio.value);
});
const gastoImportacionBOB = computed(() => roundMoney(Number(form.gasto_importacion || 0) * tipoCambio.value));
const gastoTransporteBOB = computed(() => roundMoney(Number(form.gasto_transporte || 0) * tipoCambio.value));
const gastoPapeleoBOB = computed(() => roundMoney(Number(form.gasto_papeleo || 0) * tipoCambio.value));
const gastosAdicionales = computed(() =>
  Number(form.gasto_importacion || 0) +
  Number(form.gasto_transporte || 0) +
  Number(form.gasto_papeleo || 0)
);
const costoTotalUSD = computed(() => precioCompraUSD.value + gastosAdicionales.value);
const gastosAdicionalesBOB = computed(() => roundMoney(gastosAdicionales.value * tipoCambio.value));
const costoTotalBOB = computed(() => roundMoney(costoTotalUSD.value * tipoCambio.value));
const pagoUSDDirecto = computed(() => form.pagos.filter((pago: any) => pago.moneda === 'USD').reduce((total: number, pago: any) => total + Number(pago.monto || 0), 0));
const pagoBOBDirecto = computed(() => form.pagos.filter((pago: any) => pago.moneda === 'BOB').reduce((total: number, pago: any) => total + Number(pago.monto || 0), 0));
const pagoEquivalenteUSD = computed(() => pagoUSDDirecto.value + (tipoCambio.value > 0 ? pagoBOBDirecto.value / tipoCambio.value : 0));
const pagoExcedeTotal = computed(() => pagoEquivalenteUSD.value > costoTotalUSD.value);
const estadoPagoSugerido = computed(() => pagoEquivalenteUSD.value >= costoTotalUSD.value && costoTotalUSD.value > 0 ? 'Pagado completo' : 'Pendiente');
const pagadoCompraUSD = computed(() => totalPagadoCompraUSD(compraPagoSeleccionada.value));
const pagadoCompraBOB = computed(() => roundMoney(pagadoCompraUSD.value * Number(compraPagoSeleccionada.value?.tipo_cambio_usado || 0)));
const saldoCompraUSD = computed(() => Math.max(roundMoney(Number(compraPagoSeleccionada.value?.costo_total_usd || 0) - pagadoCompraUSD.value), 0));
const saldoCompraBOB = computed(() => roundMoney(saldoCompraUSD.value * Number(compraPagoSeleccionada.value?.tipo_cambio_usado || 0)));
const nuevoPagoCompraUSD = computed(() => {
  const tc = Number(compraPagoSeleccionada.value?.tipo_cambio_usado || 0);
  return completarPagoForm.pagos.reduce((total: number, pago: any) => {
    const monto = Number(pago.monto || 0);
    if (pago.moneda === 'BOB') {
      return roundMoney(total + (tc > 0 ? monto / tc : 0));
    }
    return roundMoney(total + monto);
  }, 0);
});
const pagoCompraExcedeSaldo = computed(() => nuevoPagoCompraUSD.value > saldoCompraUSD.value);

onMounted(async () => {
  await cargarDatos();
});

async function cargarDatos() {
  loading.value = true;
  try {
    await Promise.all([obtenerCompras(), obtenerVehiculos(), obtenerProveedores()]);
  } finally {
    loading.value = false;
  }
}

async function obtenerCompras() {
  try {
    const res = await $fetch(server.HOST + '/api/v1/compras-autos', { method: 'GET' });
    compras.value = Array.isArray(res) ? res : [];
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar compras', life: 3000 });
  }
}

async function obtenerVehiculos() {
  try {
    const res = await $fetch(server.HOST + '/api/v1/vehiculos', { method: 'GET' });
    vehiculos.value = Array.isArray(res) ? res.map((vehiculo: any) => ({
      ...vehiculo,
      nombreCompleto: etiquetaVehiculo(vehiculo)
    })) : [];
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar vehiculos', life: 3000 });
  }
}

async function obtenerProveedores() {
  try {
    const res = await $fetch(server.HOST + '/api/v1/proveedores-autos', { method: 'GET' });
    proveedores.value = Array.isArray(res) ? res : [];
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar proveedores', life: 3000 });
  }
}

function abrirCompra() {
  resetForm();
  compraVisible.value = true;
}

async function registrarCompra() {
  if (!form.id_vehiculo) {
    toast.add({ severity: 'warn', summary: 'Seleccione un vehiculo', life: 3000 });
    return;
  }
  if (!form.id_proveedor) {
    toast.add({ severity: 'warn', summary: 'Seleccione proveedor', life: 3000 });
    return;
  }
  if (tipoCambio.value <= 0) {
    toast.add({ severity: 'warn', summary: 'Ingrese tipo de cambio', life: 3000 });
    return;
  }
  if (Number(form.precio_compra || 0) <= 0) {
    toast.add({ severity: 'warn', summary: 'Ingrese precio de compra', life: 3000 });
    return;
  }
  if (!validarPagos()) {
    return;
  }

  saving.value = true;
  try {
    await $fetch(server.HOST + '/api/v1/compras-autos', {
      method: 'POST',
      body: {
        id_vehiculo: form.id_vehiculo,
        id_proveedor: form.id_proveedor,
        fecha_compra: form.fecha_compra,
        moneda_precio: form.moneda_precio,
        precio_compra: Number(form.precio_compra || 0),
        tipo_cambio: tipoCambio.value,
        gasto_importacion: Number(form.gasto_importacion || 0),
        gasto_transporte: Number(form.gasto_transporte || 0),
        gasto_papeleo: Number(form.gasto_papeleo || 0),
        metodo_pago: form.metodo_pago,
        pagos: form.pagos.map((pago: any) => ({
          moneda: pago.moneda,
          metodo: pago.metodo,
          monto: Number(pago.monto || 0)
        })),
        observacion: form.observacion
      }
    });
    toast.add({ severity: 'success', summary: 'Compra registrada', life: 3000 });
    compraVisible.value = false;
    await Promise.all([obtenerCompras(), obtenerVehiculos(), obtenerProveedores()]);
  } catch (err: any) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al registrar compra', detail: err?.data || err?.message, life: 4000 });
  } finally {
    saving.value = false;
  }
}

function validarPagos() {
  if (form.pagos.length === 0) {
    toast.add({ severity: 'warn', summary: 'Agregue al menos una linea de pago', life: 3000 });
    return false;
  }
  if (form.pagos.some((pago: any) => !pago.moneda || !pago.metodo || Number(pago.monto || 0) <= 0)) {
    toast.add({ severity: 'warn', summary: 'Complete moneda, metodo y monto mayor a cero en cada linea', life: 4000 });
    return false;
  }
  if (pagoExcedeTotal.value) {
    toast.add({ severity: 'warn', summary: 'Pago mayor al costo total', life: 3000 });
    return false;
  }
  return true;
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

function agregarPagoCompra() {
  completarPagoForm.pagos.push(crearPago());
}

function eliminarPagoCompra(index: number) {
  completarPagoForm.pagos.splice(index, 1);
  if (completarPagoForm.pagos.length === 0) {
    agregarPagoCompra();
  }
}

function resetForm() {
  form.id_vehiculo = null;
  form.id_proveedor = null;
  form.fecha_compra = new Date().toISOString().slice(0, 10);
  form.moneda_precio = 'USD';
  form.precio_compra = 0;
  form.tipo_cambio = null;
  form.gasto_importacion = 0;
  form.gasto_transporte = 0;
  form.gasto_papeleo = 0;
  form.metodo_pago = 'Contado';
  form.pagos.splice(0, form.pagos.length, crearPago());
  form.observacion = '';
}

function abrirProveedorRapido() {
  Object.assign(proveedorForm, {
    nombre: '',
    ci_nit: '',
    telefono: '',
    email: '',
    direccion: '',
    tipo: '',
    observaciones: ''
  });
  ciNitProveedorExistente.value = false;
  proveedorVisible.value = true;
}

function abrirCompletarPago(compra: any) {
  compraPagoSeleccionada.value = compra;
  completarPagoForm.pagos.splice(0, completarPagoForm.pagos.length, {
    ...crearPago(),
    moneda: 'USD',
    metodo: metodoPagoCompraPorDefecto(compra),
    monto: saldoCompraUSD.value
  });
  completarPagoVisible.value = true;
}

async function guardarPagoCompra() {
  if (!compraPagoSeleccionada.value) return;
  if (!validarPagosCompra()) return;

  savingPagoCompra.value = true;
  try {
    await $fetch(server.HOST + `/api/v1/compras-autos/${compraPagoSeleccionada.value.id}/completar-pago`, {
      method: 'PATCH',
      body: {
        pagos: completarPagoForm.pagos.map((pago: any) => ({
          moneda: pago.moneda,
          metodo: pago.metodo,
          monto: Number(pago.monto || 0)
        }))
      }
    });
    toast.add({ severity: 'success', summary: 'Pago de compra actualizado', life: 3000 });
    completarPagoVisible.value = false;
    compraPagoSeleccionada.value = null;
    await obtenerCompras();
  } catch (err: any) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al completar pago', detail: err?.data || err?.message, life: 4000 });
  } finally {
    savingPagoCompra.value = false;
  }
}

function validarPagosCompra() {
  if (completarPagoForm.pagos.length === 0) {
    toast.add({ severity: 'warn', summary: 'Agregue al menos una linea de pago', life: 3000 });
    return false;
  }
  if (completarPagoForm.pagos.some((pago: any) => !pago.moneda || !pago.metodo || Number(pago.monto || 0) <= 0)) {
    toast.add({ severity: 'warn', summary: 'Complete moneda, metodo y monto mayor a cero en cada linea', life: 4000 });
    return false;
  }
  if (pagoCompraExcedeSaldo.value) {
    toast.add({ severity: 'warn', summary: 'Pago mayor al saldo pendiente', life: 3000 });
    return false;
  }
  return true;
}

async function guardarProveedorRapido({ valid }: any) {
  if (!valid) return;
  const ciNit = proveedorForm.ci_nit.trim().toLowerCase();
  const duplicado = proveedores.value.find((proveedor: any) =>
    String(proveedor.ci_nit || '').trim().toLowerCase() === ciNit
  );
  ciNitProveedorExistente.value = !!duplicado;
  if (duplicado) return;

  savingProveedor.value = true;
  try {
    const nuevo: any = await $fetch(server.HOST + '/api/v1/proveedores-autos', {
      method: 'POST',
      body: {
        ...proveedorForm,
        estado: 'Activo'
      }
    });
    await obtenerProveedores();
    form.id_proveedor = Number(nuevo?.id || 0) || null;
    proveedorVisible.value = false;
    toast.add({ severity: 'success', summary: 'Proveedor registrado', life: 3000 });
  } catch (err: any) {
    console.error(err);
    if (err?.response?.status === 409 || err?.statusCode === 409) {
      ciNitProveedorExistente.value = true;
      return;
    }
    toast.add({ severity: 'error', summary: 'Error al registrar proveedor', detail: err?.data || err?.message, life: 4000 });
  } finally {
    savingProveedor.value = false;
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

function formatTipoCambio(value: number) {
  return Number(value || 0).toLocaleString('es-BO', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 4
  });
}

function pagosCompra(compra: any) {
  if (!compra?.pagos) return [];
  if (Array.isArray(compra.pagos)) return compra.pagos;
  if (typeof compra.pagos === 'string') {
    try {
      const parsed = JSON.parse(compra.pagos);
      return Array.isArray(parsed) ? parsed : [];
    } catch {
      return [];
    }
  }
  return [];
}

function totalPagadoCompraUSD(compra: any) {
  const tc = Number(compra?.tipo_cambio_usado || 0);
  return pagosCompra(compra).reduce((total: number, pago: any) => {
    const monto = Number(pago.monto || 0);
    if (pago.moneda === 'BOB') {
      return roundMoney(total + (tc > 0 ? monto / tc : 0));
    }
    return roundMoney(total + monto);
  }, 0);
}

function metodoPagoCompraPorDefecto(compra: any) {
  const metodo = String(compra?.metodo_pago || '').trim();
  return metodosPagoMoneda.value.includes(metodo) ? metodo : 'Efectivo';
}

function valorBOB(valorGuardado: number, valorUSD: number, tc: number) {
  const guardado = Number(valorGuardado || 0);
  if (guardado > 0) return guardado;
  return roundMoney(Number(valorUSD || 0) * Number(tc || 0));
}

function roundMoney(value: number) {
  return Math.round(Number(value || 0) * 100) / 100;
}

function formatFecha(fecha: string) {
  if (!fecha) return 'N/A';
  return new Date(`${fecha}T00:00:00`).toLocaleDateString('es-BO');
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
