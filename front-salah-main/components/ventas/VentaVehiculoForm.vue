<template>
  <Toast />

  <div class="flex flex-col gap-4">
    <div class="flex items-center justify-between gap-3">
      <div>
        <h2 class="text-2xl font-bold">{{ esReserva ? 'Registrar reserva' : 'Registrar venta' }}</h2>
        <p class="text-sm text-gray-500">{{ esReserva ? 'Reserva con monto anticipado y proforma vigente.' : 'Venta al contado, credito directo o credito bancario.' }}</p>
      </div>
      <Button label="Volver" icon="pi pi-arrow-left" size="small" severity="secondary" @click="router.push(props.catalogoPath)" />
    </div>

    <div v-if="loading" class="grid grid-cols-1 gap-4 lg:grid-cols-[22rem_1fr]">
      <Skeleton height="22rem" />
      <Skeleton height="22rem" />
    </div>

    <div v-else-if="!vehiculo" class="flex flex-col items-center justify-center gap-2 py-16 text-gray-500">
      <i class="pi pi-car text-4xl"></i>
      <p>No se encontro el vehiculo seleccionado.</p>
    </div>

    <div v-else class="grid grid-cols-1 gap-4 xl:grid-cols-[24rem_1fr]">
      <section class="overflow-hidden rounded-lg border border-gray-200 bg-white">
        <div class="flex aspect-[16/10] items-center justify-center bg-gray-100">
          <img v-if="vehiculo.imagen && vehiculo.imagen !== 'N/A'" :src="vehiculo.imagen" :alt="etiquetaVehiculo(vehiculo)" class="h-full w-full object-cover">
          <i v-else class="pi pi-car text-5xl text-gray-400"></i>
        </div>

        <div class="flex flex-col gap-3 p-4">
          <div class="flex items-start justify-between gap-3">
            <div class="min-w-0">
              <h3 class="truncate text-lg font-semibold text-gray-900">{{ etiquetaVehiculo(vehiculo) }}</h3>
              <p class="text-sm text-gray-500">{{ vehiculo.categoria || 'Sin categoria' }}<span v-if="vehiculo.segmento"> / {{ vehiculo.segmento }}</span></p>
              <p v-if="vehiculo.version" class="text-xs text-gray-500">{{ vehiculo.version }}</p>
            </div>
            <Tag :value="vehiculo.estado" :severity="vehiculo.estado === 'Activo' ? 'success' : 'danger'" />
          </div>

          <div class="grid grid-cols-3 gap-2 text-sm">
            <div class="rounded-md bg-gray-50 p-3">
              <span class="text-gray-500">Precio USD</span>
              <strong class="block text-gray-900">$ {{ formatPrecio(vehiculo.precio) }}</strong>
            </div>
            <div class="rounded-md bg-gray-50 p-3">
              <span class="text-gray-500">Stock fisico</span>
              <strong class="block text-gray-900">{{ vehiculo.cantidad_disponible || 0 }}</strong>
            </div>
            <div class="rounded-md bg-gray-50 p-3">
              <span class="text-gray-500">Disponible venta</span>
              <strong class="block text-gray-900">{{ disponibilidadVenta }}</strong>
            </div>
          </div>

          <div class="rounded-md border border-yellow-200 bg-yellow-50 p-3 text-sm text-gray-700">
            La reserva aplica solo mientras la venta este registrada y la proforma no este vencida.
            <span v-if="cantidadReservadaVigente > 0" class="block pt-1 font-semibold">
              Reservado por proformas vigentes: {{ cantidadReservadaVigente }}.
            </span>
          </div>
        </div>
      </section>

      <section class="rounded-lg border border-gray-200 bg-white p-4">
        <form class="grid grid-cols-1 gap-4 lg:grid-cols-2" @submit.prevent @keydown.enter.prevent>
          <div class="flex flex-col gap-1 lg:col-span-2">
            <div class="flex items-center justify-between gap-3">
              <label for="cliente">Cliente</label>
              <Button label="Agregar cliente" icon="pi pi-plus" size="small" severity="secondary" type="button" @click="abrirModalAgregarCliente" />
            </div>
            <Select id="cliente" v-model="form.id_cliente" :options="clientes" option-label="nombreCompleto" option-value="id" placeholder="Seleccione un cliente" filter fluid size="small" />
          </div>

          <div class="flex flex-col gap-1">
            <label for="fecha">Fecha</label>
            <InputText id="fecha" v-model="form.fecha" type="date" fluid size="small" />
          </div>

          <div class="flex flex-col gap-1">
            <label for="tipo_venta">Tipo de venta</label>
            <Select id="tipo_venta" v-model="form.tipo_venta" :options="tiposVenta" fluid size="small" />
          </div>

          <div class="flex flex-col gap-1">
            <label for="cantidad">Cantidad</label>
            <InputNumber id="cantidad" v-model="form.cantidad" :min="1" :max="Math.max(disponibilidadVenta, 1)" :useGrouping="false" fluid size="small" />
          </div>

          <div class="flex flex-col gap-1">
            <label for="validez">Validez de proforma</label>
            <InputNumber id="validez" v-model="form.validez_proforma_dias" :min="1" :useGrouping="false" suffix=" dias" fluid size="small" />
          </div>

          <div class="flex flex-col gap-1">
            <label for="tipo_cambio">Tipo de cambio</label>
            <InputNumber
              id="tipo_cambio"
              v-model="form.tipo_cambio"
              :min="0"
              :minFractionDigits="2"
              :maxFractionDigits="4"
              suffix=" Bs/USD"
              fluid
              size="small"
            />
          </div>

          <div class="grid grid-cols-1 gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2 md:grid-cols-3">
            <div>
              <span class="text-xs font-semibold uppercase text-gray-500">Precio USD</span>
              <strong class="block text-gray-900">$ {{ formatPrecio(precioUnidad) }}</strong>
            </div>
            <div>
              <span class="text-xs font-semibold uppercase text-gray-500">Total USD</span>
              <strong class="block text-gray-900">$ {{ formatPrecio(precioTotal) }}</strong>
            </div>
            <div>
              <span class="text-xs font-semibold uppercase text-gray-500">Total BOB</span>
              <strong class="block text-gray-900">Bs {{ formatPrecio(montoBOB) }}</strong>
            </div>
          </div>

          <div v-if="!esReserva" class="flex flex-col gap-1">
            <label for="estado_venta">Estado de venta</label>
            <Select id="estado_venta" v-model="form.estado_venta" :options="estadosVenta" :disabled="esContado || esCredito" fluid size="small" />
          </div>

          <div class="flex flex-col gap-1">
            <label for="estado_pago">Estado de pago</label>
            <Select id="estado_pago" v-model="form.estado_pago" :options="estadosPago" :disabled="esReserva || esContado" fluid size="small" />
          </div>

          <div v-if="!esCredito" class="flex flex-col gap-1">
            <label for="metodo_pago">Metodo de pago</label>
            <Select id="metodo_pago" v-model="form.metodo_pago" :options="metodosPagoDisponibles" fluid size="small" />
          </div>

          <template v-if="esPagoMixto && !esReserva && !esCredito">
            <section class="flex flex-col gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2">
              <div class="flex items-center justify-between gap-3">
                <div>
                  <h3 class="text-sm font-semibold text-gray-900">Detalle de pago</h3>
                </div>
                <Button label="Agregar pago" icon="pi pi-plus" size="small" severity="secondary" type="button" @click="agregarPago" />
              </div>

              <div v-for="(pago, index) in form.pagos" :key="pago.key" class="grid grid-cols-1 gap-2 md:grid-cols-[9rem_13rem_1fr_2.5rem]">
                <Select v-model="pago.moneda" :options="monedasPago" placeholder="Moneda" fluid size="small" />
                <Select v-model="pago.metodo" :options="metodosPagoMoneda" placeholder="Metodo" fluid size="small" />
                <InputNumber v-model="pago.monto" mode="decimal" :min="0" :minFractionDigits="2" :maxFractionDigits="2" placeholder="Monto" fluid size="small" />
                <Button icon="pi pi-trash" severity="danger" text rounded type="button" aria-label="Eliminar pago" @click="eliminarPago(index)" />
              </div>

              <p v-if="pagoMixtoExcedeTotal" class="text-sm font-semibold text-red-600">
                El pagado supera el total de la venta.
              </p>
            </section>

            <div class="grid grid-cols-1 gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2 md:grid-cols-3">
              <div>
                <span class="text-xs font-semibold uppercase text-gray-500">Pagado equivalente USD</span>
                <strong class="block text-gray-900">$ {{ formatPrecio(pagoDetalleUSD) }}</strong>
              </div>
              <div>
                <span class="text-xs font-semibold uppercase text-gray-500">Saldo USD</span>
                <strong class="block text-gray-900">$ {{ formatPrecio(saldoDetalleUSD) }}</strong>
              </div>
              <div>
                <span class="text-xs font-semibold uppercase text-gray-500">Saldo BOB</span>
                <strong class="block text-gray-900">Bs {{ formatPrecio(saldoDetalleBOB) }}</strong>
              </div>
              <div>
                <span class="text-xs font-semibold uppercase text-gray-500">Estado sugerido</span>
                <strong class="block text-gray-900">{{ estadoPagoSugerido }}</strong>
              </div>
            </div>
          </template>

          <template v-if="esCreditoDirecto">
            <div class="flex flex-col gap-1">
              <label for="numero_cuotas">Numero de cuotas</label>
              <InputNumber id="numero_cuotas" v-model="form.numero_cuotas" :min="1" :useGrouping="false" fluid size="small" />
            </div>

            <div class="flex flex-col gap-1">
              <label for="fecha_inicio_credito">Fecha inicio pagos</label>
              <InputText id="fecha_inicio_credito" v-model="form.fecha_inicio_credito" type="date" fluid size="small" />
            </div>

            <div class="flex flex-col gap-1">
              <label for="frecuencia_pago">Frecuencia</label>
              <Select id="frecuencia_pago" v-model="form.frecuencia_pago" :options="frecuenciasPago" fluid size="small" />
            </div>

            <div class="grid grid-cols-1 gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2 md:grid-cols-4">
              <div>
                <span class="text-xs font-semibold uppercase text-gray-500">Monto financiado</span>
                <strong class="block text-gray-900">$ {{ formatPrecio(montoFinanciado) }}</strong>
              </div>
              <div>
                <span class="text-xs font-semibold uppercase text-gray-500">Monto por cuota</span>
                <strong class="block text-gray-900">$ {{ formatPrecio(montoCuota) }}</strong>
              </div>
              <div>
                <span class="text-xs font-semibold uppercase text-gray-500">Resta</span>
                <strong class="block text-gray-900">$ {{ formatPrecio(saldo) }}</strong>
              </div>
            </div>

            <section class="flex flex-col gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2">
              <div class="flex items-center justify-between gap-3">
                <div>
                  <h3 class="text-sm font-semibold text-gray-900">Detalle de pago de cuota</h3>
                </div>
                <Button label="Agregar pago" icon="pi pi-plus" size="small" severity="secondary" type="button" @click="agregarPago" />
              </div>

              <div v-for="(pago, index) in form.pagos" :key="pago.key" class="grid grid-cols-1 gap-2 md:grid-cols-[9rem_13rem_1fr_2.5rem]">
                <Select v-model="pago.moneda" :options="monedasPago" placeholder="Moneda" fluid size="small" />
                <Select v-model="pago.metodo" :options="metodosPagoMoneda" placeholder="Metodo" fluid size="small" />
                <InputNumber v-model="pago.monto" mode="decimal" :min="0" :minFractionDigits="2" :maxFractionDigits="2" placeholder="Monto" fluid size="small" />
                <Button icon="pi pi-trash" severity="danger" text rounded type="button" aria-label="Eliminar pago" @click="eliminarPago(index)" />
              </div>

              <p v-if="pagoMixtoExcedeTotal" class="text-sm font-semibold text-red-600">
                El pagado supera el total de la venta.
              </p>
            </section>

            <div class="grid grid-cols-1 gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2 md:grid-cols-3">
              <div>
                <span class="text-xs font-semibold uppercase text-gray-500">Pagado equivalente USD</span>
                <strong class="block text-gray-900">$ {{ formatPrecio(pagoDetalleUSD) }}</strong>
              </div>
              <div>
                <span class="text-xs font-semibold uppercase text-gray-500">Saldo USD</span>
                <strong class="block text-gray-900">$ {{ formatPrecio(saldoDetalleUSD) }}</strong>
              </div>
              <div>
                <span class="text-xs font-semibold uppercase text-gray-500">Saldo BOB</span>
                <strong class="block text-gray-900">Bs {{ formatPrecio(saldoDetalleBOB) }}</strong>
              </div>
            </div>
          </template>

          <template v-if="esCreditoBancario">
            <section class="flex flex-col gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2">
              <div class="flex items-center justify-between gap-3">
                <div>
                  <h3 class="text-sm font-semibold text-gray-900">Detalle de pago del enganche</h3>
                </div>
                <Button label="Agregar pago" icon="pi pi-plus" size="small" severity="secondary" type="button" @click="agregarPagoEnganche" />
              </div>

              <div v-for="(pago, index) in form.pagos_enganche" :key="pago.key" class="grid grid-cols-1 gap-2 md:grid-cols-[9rem_13rem_1fr_2.5rem]">
                <Select v-model="pago.moneda" :options="monedasPago" placeholder="Moneda" fluid size="small" />
                <Select v-model="pago.metodo" :options="metodosPagoMoneda" placeholder="Tipo de pago" fluid size="small" />
                <InputNumber v-model="pago.monto" mode="decimal" :min="0" :minFractionDigits="2" :maxFractionDigits="2" placeholder="Monto" fluid size="small" />
                <Button icon="pi pi-trash" severity="danger" text rounded type="button" aria-label="Eliminar pago" @click="eliminarPagoEnganche(index)" />
              </div>

              <p v-if="engancheExcedeTotal" class="text-sm font-semibold text-red-600">
                El enganche supera el total de la venta.
              </p>
            </section>

            <div class="flex flex-col gap-1">
              <label for="banco_entidad_financiera">Banco/entidad financiera</label>
              <InputText id="banco_entidad_financiera" v-model="form.banco_entidad_financiera" placeholder="Nombre del banco" fluid size="small" />
            </div>

            <div class="flex flex-col gap-1">
              <label for="estado_tramite_bancario">Estado del tramite</label>
              <Select id="estado_tramite_bancario" v-model="form.estado_tramite_bancario" :options="estadosTramiteBancario" fluid size="small" />
            </div>

            <div class="flex flex-col gap-1">
              <label for="fecha_estimada_desembolso">Fecha estimada de desembolso</label>
              <InputText id="fecha_estimada_desembolso" v-model="form.fecha_estimada_desembolso" type="date" fluid size="small" />
            </div>

            <div class="grid grid-cols-1 gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2 md:grid-cols-3">
              <div>
                <span class="text-xs font-semibold uppercase text-gray-500">Monto inicial</span>
                <strong class="block text-gray-900">$ {{ formatPrecio(montoInicialUSD) }}</strong>
              </div>
              <div>
                <span class="text-xs font-semibold uppercase text-gray-500">Financia banco</span>
                <strong class="block text-gray-900">$ {{ formatPrecio(montoFinanciadoBanco) }}</strong>
              </div>
              <div>
                <span class="text-xs font-semibold uppercase text-gray-500">Saldo pendiente</span>
                <strong class="block text-gray-900">$ {{ formatPrecio(saldo) }}</strong>
              </div>
            </div>
          </template>

          <div class="flex flex-col gap-1">
            <label for="estado_entrega">Estado de entrega</label>
            <Select id="estado_entrega" v-model="form.estado_entrega" :options="estadosEntrega" :disabled="esReserva" fluid size="small" />
          </div>

          <div class="flex flex-col gap-1">
            <label for="fecha_entrega">Fecha de entrega</label>
            <InputText id="fecha_entrega" v-model="form.fecha_entrega" type="date" fluid size="small" />
          </div>

          <template v-if="esCreditoDirecto">
            <div class="flex items-center gap-2">
              <Checkbox id="tiene_respaldo" v-model="form.tiene_respaldo" binary />
              <label for="tiene_respaldo">Cliente con respaldo/garantia</label>
            </div>

            <div class="flex flex-col gap-1">
              <label for="tipo_garantia">Tipo de garantia</label>
              <InputText id="tipo_garantia" v-model="form.tipo_garantia" placeholder="Ej. garante, inmueble, vehiculo" fluid size="small" />
            </div>

            <div class="flex flex-col gap-1 lg:col-span-2">
              <label for="documento_garantia">Documento o imagen de la garantia</label>
              <div class="flex flex-col gap-2 rounded-md border border-dashed border-gray-300 bg-gray-50 p-3">
                <FileUpload
                  id="documento_garantia"
                  name="documento_garantia"
                  mode="basic"
                  choose-label="Subir garantia"
                  accept="image/*,application/pdf"
                  severity="secondary"
                  class="p-button-outlined"
                  custom-upload
                  auto
                  @select="onDocumentoGarantiaSelect"
                  @clear="limpiarDocumentoGarantia"
                  @remove="limpiarDocumentoGarantia"
                />
                <small class="text-xs text-gray-500">{{ documentoGarantiaNombre || 'Formatos: imagen o PDF, maximo 10MB.' }}</small>
              </div>
            </div>

            <div class="flex flex-col gap-1 lg:col-span-2">
              <label for="datos_garante">Datos del garante/respaldo</label>
              <Textarea id="datos_garante" v-model="form.datos_garante" rows="3" auto-resize fluid />
            </div>
          </template>

          <div v-if="esReserva" class="flex flex-col gap-1 lg:col-span-2">
            <label for="monto_reserva">Monto de reserva</label>
            <InputNumber
              id="monto_reserva"
              v-model="form.monto_reserva"
              mode="currency"
              currency="USD"
              locale="es-BO"
              :min="0"
              :max="precioTotal"
              fluid
              size="small"
            />
          </div>

          <template v-if="esReserva">
            <section class="flex flex-col gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2">
              <div class="flex items-center justify-between gap-3">
                <div>
                  <h3 class="text-sm font-semibold text-gray-900">Detalle de pago</h3>
                </div>
                <Button v-if="esPagoMixto" label="Agregar pago" icon="pi pi-plus" size="small" severity="secondary" type="button" @click="agregarPago" />
              </div>

              <div v-if="!esPagoMixto" class="grid grid-cols-1 gap-2 md:grid-cols-[9rem_13rem_1fr]">
                <Select v-model="form.pagos[0].moneda" :options="monedasPago" placeholder="Moneda" fluid size="small" />
                <Select v-model="form.metodo_pago" :options="metodosPagoMoneda" placeholder="Metodo" fluid size="small" />
                <InputNumber v-model="form.pagos[0].monto" mode="decimal" :min="0" :minFractionDigits="2" :maxFractionDigits="2" placeholder="Monto" fluid size="small" />
              </div>

              <div v-for="(pago, index) in form.pagos" v-else :key="pago.key" class="grid grid-cols-1 gap-2 md:grid-cols-[9rem_13rem_1fr_2.5rem]">
                <Select v-model="pago.moneda" :options="monedasPago" placeholder="Moneda" fluid size="small" />
                <Select v-model="pago.metodo" :options="metodosPagoMoneda" placeholder="Metodo" fluid size="small" />
                <InputNumber v-model="pago.monto" mode="decimal" :min="0" :minFractionDigits="2" :maxFractionDigits="2" placeholder="Monto" fluid size="small" />
                <Button icon="pi pi-trash" severity="danger" text rounded type="button" aria-label="Eliminar pago" @click="eliminarPago(index)" />
              </div>

              <p v-if="pagoMixtoExcedeTotal" class="text-sm font-semibold text-red-600">
                El pagado supera el monto de reserva.
              </p>
            </section>

            <div class="grid grid-cols-1 gap-3 rounded-md border border-gray-200 bg-gray-50 p-3 lg:col-span-2 md:grid-cols-3">
              <div>
                <span class="text-xs font-semibold uppercase text-gray-500">Pagado equivalente USD</span>
                <strong class="block text-gray-900">$ {{ formatPrecio(pagoDetalleUSD) }}</strong>
                <span class="block text-sm font-semibold text-gray-700">Bs {{ formatPrecio(pagoDetalleBOB) }}</span>
              </div>
              <div>
                <span class="text-xs font-semibold uppercase text-gray-500">Saldo reserva USD</span>
                <strong class="block text-gray-900">$ {{ formatPrecio(saldoDetalleUSD) }}</strong>
                <span class="block text-sm font-semibold text-gray-700">Bs {{ formatPrecio(saldoDetalleBOB) }}</span>
              </div>
              <div>
                <span class="text-xs font-semibold uppercase text-gray-500">Estado sugerido</span>
                <strong class="block text-gray-900">{{ estadoPagoSugerido }}</strong>
              </div>
            </div>
          </template>

          <div class="flex flex-col gap-1 lg:col-span-2">
            <label for="observacion">Observacion</label>
            <Textarea id="observacion" v-model="form.observacion" rows="3" auto-resize fluid />
          </div>

          <div class="flex justify-end gap-2 border-t border-gray-100 pt-4 lg:col-span-2">
            <Button label="Cancelar" severity="secondary" type="button" @click="router.push(props.catalogoPath)" />
            <Button :label="esReserva ? 'Registrar reserva' : 'Registrar venta'" icon="pi pi-check" type="button" :loading="saving" @click="registrarVenta" />
          </div>
        </form>
      </section>
    </div>
  </div>

  <modalAgregarCliente
    v-if="modalClienteVisible"
    :open="modalClienteVisible"
    @close="modalClienteVisible = false"
    @update="actualizarClientesDespuesDeCrear"
    @success="toast.add({ severity: 'success', summary: 'Cliente agregado', life: 3000 })"
    @error="toast.add({ severity: 'error', summary: 'Error al agregar cliente', life: 3000 })"
  />
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import modalAgregarCliente from '~/components/admin/clientes/modalAgregarCliente.vue';
import { server } from '~/server/server';
import Button from 'primevue/button';
import Checkbox from 'primevue/checkbox';
import FileUpload from 'primevue/fileupload';
import InputNumber from 'primevue/inputnumber';
import InputText from 'primevue/inputtext';
import Select from 'primevue/select';
import Skeleton from 'primevue/skeleton';
import Tag from 'primevue/tag';
import Textarea from 'primevue/textarea';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';

const props = withDefaults(defineProps<{
  catalogoPath?: string;
  historialPath?: string;
}>(), {
  catalogoPath: '/ventas/catalogo-vehiculos',
  historialPath: '/ventas/historial-ventas'
});

const route = useRoute();
const router = useRouter();
const toast = useToast();
const { descargarPDFVenta } = useVentaPdf();

const loading = ref(true);
const saving = ref(false);
const vehiculo = ref<any>(null);
const clientes = ref<any[]>([]);
const ventas = ref<any[]>([]);
const modalClienteVisible = ref(false);
const clienteIdsAntesModal = ref<Set<number>>(new Set());
const documentoGarantia = ref<File | null>(null);
const documentoGarantiaNombre = ref('');
const tiposVenta = ref(['Contado', 'credito_directo', 'credito_bancario', 'Reserva']);
const estadosVenta = ref(['Registrada', 'Completada', 'en_credito']);
const estadosPago = ref(['Pendiente', 'Parcial', 'Pagado completo']);
const metodosPago = ref(['QR', 'Transferencia', 'Efectivo', 'Mixto']);
const metodosPagoMoneda = ref(['QR', 'Transferencia', 'Efectivo']);
const monedasPago = ref(['USD', 'BOB']);
const estadosEntrega = ref(['Pendiente', 'Entregado']);
const estadosDesembolso = ref(['Pendiente', 'Registrado']);
const estadosTramiteBancario = ref(['Pendiente', 'En evaluacion', 'Aprobado', 'Rechazado']);
const frecuenciasPago = ref(['mensual', 'quincenal', 'semanal']);

const form = reactive({
  id_cliente: null as number | null,
  fecha: new Date().toISOString().slice(0, 10),
  tipo_venta: 'Contado',
  cantidad: 1,
  tipo_cambio: null as number | null,
  pago_usd: 0,
  pago_bob: 0,
  pagos: [
    crearPago()
  ],
  pagos_enganche: [
    crearPago()
  ],
  validez_proforma_dias: 15,
  estado_venta: 'Completada',
  estado_pago: 'Pagado completo',
  metodo_pago: 'Efectivo',
  estado_entrega: 'Pendiente',
  fecha_entrega: '',
  referencia_bancaria: '',
  estado_desembolso: 'Pendiente',
  banco_entidad_financiera: '',
  estado_tramite_bancario: 'Pendiente',
  monto_financiar_banco: 0,
  fecha_estimada_desembolso: '',
  fecha_desembolso_banco: '',
  monto_desembolsado_banco: 0,
  monto_reserva: 0,
  monto_inicial: 0,
  numero_cuotas: 12,
  fecha_inicio_credito: new Date().toISOString().slice(0, 10),
  frecuencia_pago: 'mensual',
  tiene_respaldo: false,
  tipo_garantia: '',
  datos_garante: '',
  observacion: ''
});

const precioUnidad = computed(() => Number(vehiculo.value?.precio_usd || vehiculo.value?.precio || 0));
const precioTotal = computed(() => precioUnidad.value * Number(form.cantidad || 0));
const tipoCambio = computed(() => Number(form.tipo_cambio || 0));
const montoBOB = computed(() => precioTotal.value * tipoCambio.value);
const esContado = computed(() => form.tipo_venta === 'Contado');
const esReserva = computed(() => form.tipo_venta === 'Reserva');
const esCreditoDirecto = computed(() => form.tipo_venta === 'credito_directo');
const esCreditoBancario = computed(() => form.tipo_venta === 'credito_bancario');
const esCredito = computed(() => esCreditoDirecto.value || esCreditoBancario.value || form.tipo_venta === 'Credito');
const metodosPagoDisponibles = computed(() => esCredito.value ? metodosPagoMoneda.value : metodosPago.value);
const esPagoMixto = computed(() => form.metodo_pago === 'Mixto');
const muestraDetallePago = computed(() => (esPagoMixto.value && !esReserva.value && !esCreditoBancario.value) || esCreditoDirecto.value || esReserva.value);
const pagosParaCalculo = computed(() => {
  if (esReserva.value && !esPagoMixto.value) {
    return form.pagos.slice(0, 1);
  }
  return form.pagos;
});
const pagoMixtoUSDDirecto = computed(() => pagosParaCalculo.value
  .filter((pago: any) => pago.moneda === 'USD')
  .reduce((total: number, pago: any) => total + Number(pago.monto || 0), 0));
const pagoMixtoBOBDirecto = computed(() => pagosParaCalculo.value
  .filter((pago: any) => pago.moneda === 'BOB')
  .reduce((total: number, pago: any) => total + Number(pago.monto || 0), 0));
const pagoMixtoUSD = computed(() => pagoMixtoUSDDirecto.value + (tipoCambio.value > 0 ? pagoMixtoBOBDirecto.value / tipoCambio.value : 0));
const pagoMixtoBOB = computed(() => pagoMixtoBOBDirecto.value + (pagoMixtoUSDDirecto.value * tipoCambio.value));
const engancheUSDDirecto = computed(() => form.pagos_enganche
  .filter((pago: any) => pago.moneda === 'USD')
  .reduce((total: number, pago: any) => total + Number(pago.monto || 0), 0));
const engancheBOBDirecto = computed(() => form.pagos_enganche
  .filter((pago: any) => pago.moneda === 'BOB')
  .reduce((total: number, pago: any) => total + Number(pago.monto || 0), 0));
const engancheEquivalenteUSD = computed(() => engancheUSDDirecto.value + (tipoCambio.value > 0 ? engancheBOBDirecto.value / tipoCambio.value : 0));
const objetivoDetallePagoUSD = computed(() => esReserva.value ? Number(form.monto_reserva || 0) : precioTotal.value);
const objetivoDetallePagoBOB = computed(() => objetivoDetallePagoUSD.value * tipoCambio.value);
const pagoDetalleUSD = computed(() => pagoMixtoUSD.value);
const pagoDetalleBOB = computed(() => pagoMixtoBOB.value);
const saldoDetalleUSD = computed(() => Math.max(objetivoDetallePagoUSD.value - pagoDetalleUSD.value, 0));
const saldoDetalleBOB = computed(() => Math.max(objetivoDetallePagoBOB.value - pagoDetalleBOB.value, 0));
const pagoMixtoExcedeTotal = computed(() => muestraDetallePago.value && pagoDetalleUSD.value > objetivoDetallePagoUSD.value);
const engancheExcedeTotal = computed(() => esCreditoBancario.value && engancheEquivalenteUSD.value > precioTotal.value);
const estadoPagoSugerido = computed(() => {
  if (!muestraDetallePago.value) return form.estado_pago;
  if (pagoDetalleUSD.value <= 0) return 'Pendiente';
  if (pagoDetalleUSD.value >= objetivoDetallePagoUSD.value) return 'Pagado completo';
  return 'Parcial';
});
const montoInicialUSD = computed(() => {
  if (esCreditoDirecto.value) {
    return pagoDetalleUSD.value;
  }
  if (esCreditoBancario.value) {
    return engancheEquivalenteUSD.value;
  }
  return Number(form.monto_inicial || 0);
});
const montoInicial = computed(() => {
  if (esReserva.value) {
    return Number(form.monto_reserva || 0);
  }
  if (esCredito.value) {
    return montoInicialUSD.value;
  }
  return 0;
});
const saldo = computed(() => Math.max(precioTotal.value - montoInicial.value, 0));
const montoFinanciado = computed(() => esCredito.value ? Math.max(precioTotal.value - montoInicialUSD.value, 0) : 0);
const montoFinanciadoBanco = computed(() => esCreditoBancario.value ? Math.max(precioTotal.value - montoInicialUSD.value, 0) : 0);
const montoCuota = computed(() => Number(form.numero_cuotas || 0) > 0 ? montoFinanciado.value / Number(form.numero_cuotas || 1) : 0);
const cantidadReservadaVigente = computed(() => {
  const idVehiculo = Number(route.query.id || 0);
  return ventas.value
    .filter((venta: any) =>
      Number(venta.id_vehiculo || 0) === idVehiculo &&
      venta.estado_venta === 'Registrada' &&
      !venta.proforma_vencida
    )
    .reduce((total: number, venta: any) => total + Number(venta.cantidad || 0), 0);
});
const disponibilidadVenta = computed(() => {
  return Math.max(Number(vehiculo.value?.cantidad_disponible || 0) - cantidadReservadaVigente.value, 0);
});

onMounted(async () => {
  await cargarDatos();
});

watch(() => form.tipo_venta, () => {
  if (esReserva.value) {
    form.estado_venta = 'Registrada';
    form.estado_pago = 'Pagado completo';
    form.estado_entrega = 'Pendiente';
    form.fecha_entrega = '';
    return;
  }

  if (esContado.value) {
    form.estado_venta = 'Completada';
    form.estado_pago = 'Pagado completo';
    return;
  }

  if (esCredito.value) {
    form.estado_venta = 'en_credito';
    form.estado_pago = 'Pendiente';
    if (form.metodo_pago === 'Mixto') {
      form.metodo_pago = 'Efectivo';
    }
  }

  if (esCreditoBancario.value) {
    form.numero_cuotas = 0;
    form.fecha_inicio_credito = '';
    form.frecuencia_pago = 'mensual';
    form.estado_desembolso = 'Pendiente';
  }

  if (esCreditoDirecto.value && !form.fecha_inicio_credito) {
    form.fecha_inicio_credito = new Date().toISOString().slice(0, 10);
    form.numero_cuotas = form.numero_cuotas || 12;
  }

  if (!esCreditoDirecto.value) {
    limpiarDocumentoGarantia();
  }
}, { immediate: true });

watch(precioTotal, () => {
  if (Number(form.monto_reserva || 0) > precioTotal.value) {
    form.monto_reserva = precioTotal.value;
  }
  if (Number(form.monto_inicial || 0) > precioTotal.value) {
    form.monto_inicial = precioTotal.value;
  }
});

watch(estadoPagoSugerido, (estado) => {
  if (muestraDetallePago.value) {
    form.estado_pago = estado;
  }
});

async function cargarDatos() {
  loading.value = true;
  try {
    const idVehiculo = route.query.id;
    if (!idVehiculo) {
      vehiculo.value = null;
      return;
    }

    const [resVehiculo, resClientes, resVentas] = await Promise.all([
      $fetch(server.HOST + '/api/v1/vehiculos/' + idVehiculo, { method: 'GET' }),
      $fetch(server.HOST + '/api/v1/clientes', { method: 'GET' }),
      $fetch(server.HOST + '/api/v1/ventas', { method: 'GET' })
    ]);

    vehiculo.value = resVehiculo;
    ventas.value = Array.isArray(resVentas) ? resVentas : [];
    asignarClientes(resClientes);
    if (Number(form.cantidad || 0) > disponibilidadVenta.value) {
      form.cantidad = Math.max(disponibilidadVenta.value, 1);
    }
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al cargar la venta', life: 3000 });
  } finally {
    loading.value = false;
  }
}

function abrirModalAgregarCliente() {
  clienteIdsAntesModal.value = new Set(clientes.value.map((cliente: any) => Number(cliente.id)));
  modalClienteVisible.value = true;
}

async function actualizarClientesDespuesDeCrear() {
  try {
    await obtenerClientes();
    const nuevoCliente = clientes.value.find((cliente: any) => !clienteIdsAntesModal.value.has(Number(cliente.id)));

    if (nuevoCliente?.id) {
      form.id_cliente = Number(nuevoCliente.id);
    }
  } catch (err) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al actualizar clientes', life: 3000 });
  }
}

async function obtenerClientes() {
  const resClientes = await $fetch(server.HOST + '/api/v1/clientes', { method: 'GET' });
  asignarClientes(resClientes);
}

function asignarClientes(resClientes: any) {
  const clientesActivos = Array.isArray(resClientes) ? resClientes.filter((cliente: any) => cliente.estado === 'Activo') : [];
  clientes.value = clientesActivos.map((cliente: any) => ({
    ...cliente,
    nombreCompleto: `${cliente.nombre || ''} ${cliente.apellido || ''} - ${cliente.ci || 'Sin CI'}`.trim()
  }));
}

async function registrarVenta() {
  if (!form.id_cliente || !vehiculo.value) {
    toast.add({ severity: 'warn', summary: 'Seleccione cliente y vehiculo', life: 3000 });
    return;
  }
  if (disponibilidadVenta.value <= 0) {
    toast.add({ severity: 'warn', summary: 'Sin disponibilidad vendible', detail: 'El stock esta reservado por proformas vigentes.', life: 4000 });
    return;
  }
  if (Number(form.cantidad || 0) > disponibilidadVenta.value) {
    toast.add({ severity: 'warn', summary: 'Cantidad mayor a la disponible', detail: `Disponible para venta: ${disponibilidadVenta.value}`, life: 4000 });
    return;
  }
  if (tipoCambio.value <= 0) {
    toast.add({ severity: 'warn', summary: 'Ingrese el tipo de cambio', life: 3000 });
    return;
  }
  if (esReserva.value && Number(form.monto_reserva || 0) <= 0) {
    toast.add({ severity: 'warn', summary: 'Ingrese el monto de reserva', life: 3000 });
    return;
  }
  if (esReserva.value && Number(form.monto_reserva || 0) > precioTotal.value) {
    toast.add({ severity: 'warn', summary: 'Reserva mayor al precio total', life: 3000 });
    return;
  }
  if (muestraDetallePago.value && pagoMixtoExcedeTotal.value) {
    toast.add({ severity: 'warn', summary: esReserva.value ? 'Pago mayor al monto de reserva' : 'Pago mayor al precio total', life: 3000 });
    return;
  }
  if (muestraDetallePago.value && !validarDetallePago()) {
    return;
  }
  if (esCreditoDirecto.value) {
    if (montoInicialUSD.value <= 0 || montoInicialUSD.value >= precioTotal.value) {
      toast.add({ severity: 'warn', summary: 'Detalle de pago no valido', life: 3000 });
      return;
    }
    if (Number(form.numero_cuotas || 0) <= 0) {
      toast.add({ severity: 'warn', summary: 'Ingrese el numero de cuotas', life: 3000 });
      return;
    }
    if (!form.fecha_inicio_credito) {
      toast.add({ severity: 'warn', summary: 'Ingrese la fecha de inicio de pagos', life: 3000 });
      return;
    }
  }
  if (esCreditoBancario.value) {
    if (!validarDetalleEnganche()) {
      return;
    }
    if (montoInicialUSD.value <= 0 || montoInicialUSD.value >= precioTotal.value) {
      toast.add({ severity: 'warn', summary: 'Monto inicial no valido', detail: 'Debe ser mayor a 0 y menor al total.', life: 4000 });
      return;
    }
    if (engancheExcedeTotal.value) {
      toast.add({ severity: 'warn', summary: 'Enganche mayor al precio total', life: 3000 });
      return;
    }
    if (!form.banco_entidad_financiera.trim()) {
      toast.add({ severity: 'warn', summary: 'Ingrese banco o entidad financiera', life: 3000 });
      return;
    }
    if (!form.estado_tramite_bancario) {
      toast.add({ severity: 'warn', summary: 'Seleccione estado del tramite', life: 3000 });
      return;
    }
  }
  if (esCreditoDirecto.value && (!form.tiene_respaldo || !form.tipo_garantia.trim() || !form.datos_garante.trim())) {
    toast.add({ severity: 'warn', summary: 'Credito directo bloqueado', detail: 'Debe registrar respaldo, tipo de garantia y datos del garante.', life: 4000 });
    return;
  }
  if (esCreditoDirecto.value && documentoGarantia.value && documentoGarantia.value.size > 10 * 1024 * 1024) {
    toast.add({ severity: 'warn', summary: 'Documento muy grande', detail: 'El documento de garantia no debe superar los 10MB.', life: 4000 });
    return;
  }

  saving.value = true;
  try {
    const user = localStorage.getItem('user');
    const userId = user ? JSON.parse(user)?.id : null;
    const payload = construirVentaFormData(Number(userId || 0));

    const ventaRegistrada: any = await $fetch(server.HOST + '/api/v1/ventas', {
      method: 'POST',
      body: payload
    });

    toast.add({
      severity: 'success',
      summary: esReserva.value ? 'Reserva exitosa' : 'Venta exitosa',
      detail: esReserva.value ? 'La reserva fue registrada correctamente.' : 'La venta fue registrada correctamente.',
      life: 3000
    });
    if (ventaRegistrada?.id) {
      await descargarPDFVenta(ventaRegistrada.id);
    }
    setTimeout(() => {
      router.push(props.historialPath);
    }, 900);
  } catch (err: any) {
    console.error(err);
    toast.add({ severity: 'error', summary: 'Error al registrar venta', detail: err?.data || err?.message, life: 4000 });
  } finally {
    saving.value = false;
  }
}

function construirVentaFormData(userId: number) {
  const formData = new FormData();
  formData.append('id_cliente', String(form.id_cliente || ''));
  formData.append('id_vehiculo', String(Number(route.query.id)));
  formData.append('id_usuario', String(userId || 0));
  formData.append('fecha', form.fecha);
  formData.append('tipo_venta', form.tipo_venta);
  formData.append('cantidad', String(Number(form.cantidad)));
  formData.append('tipo_cambio', String(tipoCambio.value));
  formData.append('pago_usd', String(esCreditoBancario.value ? Number(engancheUSDDirecto.value || 0) : (muestraDetallePago.value ? Number(pagoMixtoUSDDirecto.value || 0) : 0)));
  formData.append('pago_bob', String(esCreditoBancario.value ? Number(engancheBOBDirecto.value || 0) : (muestraDetallePago.value ? Number(pagoMixtoBOBDirecto.value || 0) : 0)));
  formData.append('pagos', JSON.stringify(muestraDetallePago.value ? detallePagoPayload() : []));
  formData.append('pagos_enganche', JSON.stringify(esCreditoBancario.value ? detalleEnganchePayload() : []));
  formData.append('validez_proforma_dias', String(Number(form.validez_proforma_dias || 15)));
  formData.append('estado_venta', esReserva.value ? 'Registrada' : form.estado_venta);
  formData.append('estado_pago', form.estado_pago);
  formData.append('metodo_pago', esCreditoBancario.value ? resumenDetalleEnganche() : (esCredito.value || esReserva.value ? resumenDetallePago() : form.metodo_pago));
  formData.append('estado_entrega', esReserva.value ? 'Pendiente' : form.estado_entrega);
  formData.append('fecha_entrega', esReserva.value ? '' : form.fecha_entrega);
  formData.append('monto_reserva', String(esReserva.value ? Number(form.monto_reserva || 0) : 0));
  formData.append('monto_inicial', String(esCredito.value ? Number(montoInicialUSD.value || 0) : 0));
  formData.append('numero_cuotas', String(esCreditoDirecto.value ? Number(form.numero_cuotas || 0) : 0));
  formData.append('fecha_inicio_credito', esCreditoDirecto.value ? form.fecha_inicio_credito : '');
  formData.append('frecuencia_pago', esCreditoDirecto.value ? form.frecuencia_pago : '');
  formData.append('tiene_respaldo', String(esCreditoDirecto.value ? form.tiene_respaldo : false));
  formData.append('tipo_garantia', esCreditoDirecto.value ? form.tipo_garantia : '');
  formData.append('datos_garante', esCreditoDirecto.value ? form.datos_garante : '');
  formData.append('referencia_bancaria', '');
  formData.append('estado_desembolso', esCreditoBancario.value ? form.estado_desembolso : '');
  formData.append('banco_entidad_financiera', esCreditoBancario.value ? form.banco_entidad_financiera : '');
  formData.append('estado_tramite_bancario', esCreditoBancario.value ? form.estado_tramite_bancario : '');
  formData.append('monto_financiar_banco', String(esCreditoBancario.value ? Number(montoFinanciadoBanco.value || 0) : 0));
  formData.append('fecha_estimada_desembolso', esCreditoBancario.value ? form.fecha_estimada_desembolso : '');
  formData.append('numero_operacion_banco', '');
  formData.append('fecha_desembolso_banco', '');
  formData.append('monto_desembolsado_banco', '0');
  formData.append('id_usuario_desembolso_banco', '0');
  formData.append('observacion', form.observacion);
  if (esCreditoDirecto.value && documentoGarantia.value) {
    formData.append('documento_garantia', documentoGarantia.value);
  }
  return formData;
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

function agregarPagoEnganche() {
  form.pagos_enganche.push(crearPago());
}

function eliminarPagoEnganche(index: number) {
  form.pagos_enganche.splice(index, 1);
  if (form.pagos_enganche.length === 0) {
    agregarPagoEnganche();
  }
}

function detallePagoPayload() {
  if (esReserva.value && !esPagoMixto.value) {
    return [{
      moneda: form.pagos[0]?.moneda || 'USD',
      metodo: form.metodo_pago,
      monto: Number(form.pagos[0]?.monto || 0)
    }];
  }

  return form.pagos.map((pago: any) => ({
    moneda: pago.moneda,
    metodo: pago.metodo,
    monto: Number(pago.monto || 0)
  }));
}

function detalleEnganchePayload() {
  return form.pagos_enganche.map((pago: any) => ({
    moneda: pago.moneda,
    metodo: pago.metodo,
    monto: Number(pago.monto || 0)
  }));
}

function resumenDetallePago() {
  const pagosValidos = detallePagoPayload().filter((pago: any) => pago.moneda && pago.metodo && Number(pago.monto || 0) > 0);
  if (pagosValidos.length === 0) {
    return 'Efectivo';
  }

  const primerPago = pagosValidos[0];
  const esUnicoTipo = pagosValidos.every((pago: any) =>
    pago.moneda === primerPago.moneda && pago.metodo === primerPago.metodo
  );
  return esUnicoTipo ? primerPago.metodo : 'Mixto';
}

function resumenDetalleEnganche() {
  const pagosValidos = detalleEnganchePayload().filter((pago: any) => pago.moneda && pago.metodo && Number(pago.monto || 0) > 0);
  if (pagosValidos.length === 0) {
    return 'Efectivo';
  }

  const primerPago = pagosValidos[0];
  const esUnicoTipo = pagosValidos.every((pago: any) =>
    pago.moneda === primerPago.moneda && pago.metodo === primerPago.metodo
  );
  return esUnicoTipo ? primerPago.metodo : 'Mixto';
}

function validarDetallePago() {
  if (form.pagos.length === 0) {
    toast.add({ severity: 'warn', summary: 'Agregue al menos una linea de pago', life: 3000 });
    return false;
  }

  const lineaInvalida = form.pagos.find((pago: any) =>
    !pago.moneda || !pago.metodo || Number(pago.monto || 0) <= 0
  );
  if (lineaInvalida) {
    toast.add({ severity: 'warn', summary: 'Complete moneda, metodo y monto mayor a cero en cada linea', life: 4000 });
    return false;
  }

  if (pagoMixtoExcedeTotal.value) {
    toast.add({ severity: 'warn', summary: esReserva.value ? 'Pago mayor al monto de reserva' : 'Pago mayor al precio total', life: 3000 });
    return false;
  }

  return true;
}

function validarDetalleEnganche() {
  if (form.pagos_enganche.length === 0) {
    toast.add({ severity: 'warn', summary: 'Agregue al menos una linea de pago del enganche', life: 3000 });
    return false;
  }

  const lineaInvalida = form.pagos_enganche.find((pago: any) =>
    !pago.moneda || !pago.metodo || Number(pago.monto || 0) <= 0
  );
  if (lineaInvalida) {
    toast.add({ severity: 'warn', summary: 'Complete moneda, tipo de pago y monto mayor a cero en cada linea del enganche', life: 4000 });
    return false;
  }

  if (engancheExcedeTotal.value) {
    toast.add({ severity: 'warn', summary: 'Enganche mayor al precio total', life: 3000 });
    return false;
  }

  return true;
}

function onDocumentoGarantiaSelect(event: any) {
  const file = event.files?.[0] || null;
  if (!file) {
    limpiarDocumentoGarantia();
    return;
  }
  documentoGarantia.value = file;
  documentoGarantiaNombre.value = file.name;
}

function limpiarDocumentoGarantia() {
  documentoGarantia.value = null;
  documentoGarantiaNombre.value = '';
}

function formatPrecio(precio: number) {
  return Number(precio || 0).toLocaleString('es-BO', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  });
}

function etiquetaVehiculo(item: any) {
  return [item?.marca, item?.modelo, item?.anio].filter(Boolean).join(' ') || item?.nombre || 'Vehiculo';
}
</script>
