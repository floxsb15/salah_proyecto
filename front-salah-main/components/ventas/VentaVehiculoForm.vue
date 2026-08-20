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

    <div v-else-if="!vehiculo && !esReserva" class="flex flex-col items-center justify-center gap-2 py-16 text-gray-500">
      <i class="pi pi-car text-4xl"></i>
      <p>No se encontro el vehiculo seleccionado.</p>
    </div>

    <div v-else class="grid grid-cols-1 gap-4 xl:grid-cols-[24rem_1fr]">
      <section class="overflow-hidden rounded-lg border border-gray-200 bg-white">
        <div class="flex aspect-[16/10] items-center justify-center bg-gray-100">
          <img v-if="vehiculoEnStock && vehiculo.imagen && vehiculo.imagen !== 'N/A'" :src="vehiculo.imagen" :alt="etiquetaVehiculo(vehiculo)" class="h-full w-full object-cover">
          <i v-else class="pi pi-car text-5xl text-gray-400"></i>
        </div>

        <div class="flex flex-col gap-3 p-4">
          <div class="flex items-start justify-between gap-3">
            <div class="min-w-0">
              <h3 class="truncate text-lg font-semibold text-gray-900">{{ vehiculoEnStock ? etiquetaVehiculo(vehiculo) : etiquetaVehiculoPedido }}</h3>
              <p class="text-sm text-gray-500">
                <template v-if="vehiculoEnStock">
                  {{ vehiculo.categoria || 'Sin categoria' }}<span v-if="vehiculo.segmento"> / {{ vehiculo.segmento }}</span>
                </template>
                <template v-else>Reserva a pedido / Importacion</template>
              </p>
              <p v-if="vehiculoEnStock && vehiculo.version" class="text-xs text-gray-500">{{ vehiculo.version }}</p>
              <p v-if="!vehiculoEnStock && form.pedido_pais_origen" class="text-xs text-gray-500">Origen: {{ form.pedido_pais_origen }}</p>
            </div>
            <Tag :value="vehiculoEnStock ? vehiculo.estado : 'A pedido'" :severity="vehiculoEnStock && vehiculo.estado === 'Activo' ? 'success' : 'warning'" />
          </div>

          <div class="grid grid-cols-3 gap-2 text-sm">
            <div class="rounded-md bg-gray-50 p-3">
              <span class="text-gray-500">Precio USD</span>
              <strong class="block text-gray-900">$ {{ formatPrecio(precioUnidad) }}</strong>
            </div>
            <div class="rounded-md bg-gray-50 p-3">
              <span class="text-gray-500">Stock fisico</span>
              <strong class="block text-gray-900">{{ vehiculoEnStock ? (vehiculo.cantidad_disponible || 0) : 'N/A' }}</strong>
            </div>
            <div class="rounded-md bg-gray-50 p-3">
              <span class="text-gray-500">Disponible venta</span>
              <strong class="block text-gray-900">{{ vehiculoEnStock ? disponibilidadVenta : 'Pedido' }}</strong>
            </div>
          </div>

          <div v-if="vehiculoEnStock" class="rounded-md border border-yellow-200 bg-yellow-50 p-3 text-sm text-gray-700">
            La reserva aplica solo mientras la venta este registrada y la proforma no este vencida.
            <span v-if="cantidadReservadaVigente > 0" class="block pt-1 font-semibold">
              Reservado por proformas vigentes: {{ cantidadReservadaVigente }}.
            </span>
          </div>
          <div v-else class="rounded-md border border-orange-200 bg-orange-50 p-3 text-sm text-gray-700">
            Esta reserva queda en estado Importando hasta que el vehiculo llegue y se registre en inventario.
          </div>
        </div>
      </section>

      <section class="rounded-lg border border-gray-200 bg-white p-4">
        <form class="grid grid-cols-1 gap-4 lg:grid-cols-2" @submit.prevent @keydown.enter.prevent>
          <div class="flex flex-col gap-1 lg:col-span-2">
            <label for="cliente">Cliente</label>
            <Select id="cliente" v-model="form.id_cliente" :options="clientes" option-label="nombreCompleto" option-value="id" placeholder="Seleccione un cliente" filter fluid size="small" />
          </div>

          <div v-if="esReserva" class="flex flex-col gap-1 lg:col-span-2">
            <label for="tipo_reserva">Tipo de reserva</label>
            <Select id="tipo_reserva" v-model="form.tipo_reserva" :options="tiposReserva" option-label="label" option-value="value" fluid size="small" />
          </div>

          <div v-if="esReserva && !esReservaPedido" class="flex flex-col gap-1 lg:col-span-2">
            <label for="id_vehiculo">Vehiculo en stock</label>
            <Select id="id_vehiculo" v-model="form.id_vehiculo" :options="vehiculos" option-label="nombreCompleto" option-value="id" placeholder="Seleccione vehiculo del catalogo" filter fluid size="small" />
          </div>

          <section v-if="esReservaPedido" class="grid grid-cols-1 gap-3 rounded-md border border-orange-200 bg-orange-50 p-3 lg:col-span-2 md:grid-cols-3">
            <div class="flex flex-col gap-1">
              <label for="pedido_marca">Marca</label>
              <InputText id="pedido_marca" v-model="form.pedido_marca" placeholder="Ej. Toyota" fluid size="small" />
            </div>
            <div class="flex flex-col gap-1">
              <label for="pedido_modelo">Modelo</label>
              <InputText id="pedido_modelo" v-model="form.pedido_modelo" placeholder="Ej. Land Cruiser" fluid size="small" />
            </div>
            <div class="flex flex-col gap-1">
              <label for="pedido_anio">Anio</label>
              <InputNumber id="pedido_anio" v-model="form.pedido_anio" :min="1900" :useGrouping="false" fluid size="small" />
            </div>
            <div class="flex flex-col gap-1">
              <label for="pedido_color">Color</label>
              <InputText id="pedido_color" v-model="form.pedido_color" placeholder="Ej. Blanco perlado" fluid size="small" />
            </div>
            <div class="flex flex-col gap-1">
              <label for="pedido_pais_origen">Pais de origen</label>
              <InputText id="pedido_pais_origen" v-model="form.pedido_pais_origen" placeholder="Ej. Japon" fluid size="small" />
            </div>
            <div class="flex flex-col gap-1">
              <label for="pedido_proveedor">Proveedor</label>
              <InputText id="pedido_proveedor" v-model="form.pedido_proveedor" placeholder="Opcional" fluid size="small" />
            </div>
            <div class="flex flex-col gap-1">
              <label for="pedido_llegada_estimada">Llegada estimada</label>
              <InputText id="pedido_llegada_estimada" v-model="form.pedido_llegada_estimada" placeholder="Fecha o rango" fluid size="small" />
            </div>
            <div class="flex flex-col gap-1">
              <label for="precio_pedido">Precio pactado USD</label>
              <InputNumber id="precio_pedido" v-model="form.precio_pedido" mode="currency" currency="USD" locale="es-BO" :min="0" fluid size="small" />
            </div>
            <div class="flex flex-col gap-1 md:col-span-3">
              <label for="pedido_version">Version / especificaciones</label>
              <Textarea id="pedido_version" v-model="form.pedido_version" rows="2" auto-resize fluid />
            </div>
          </section>

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

          <CreditoFormFields
            v-if="esCredito"
            :form="form"
            :monto-financiado="montoFinanciado"
            :monto-cuota="montoCuota"
            :saldo="saldo"
            :pago-detalle-u-s-d="pagoDetalleUSD"
            :pago-detalle-b-o-b="pagoMixtoBOBDirecto"
            :saldo-detalle-u-s-d="saldoDetalleUSD"
            :saldo-detalle-b-o-b="saldoDetalleBOB"
            :pago-excede-total="pagoMixtoExcedeTotal"
            :documento-garantia-nombre="documentoGarantiaNombre"
            @documento-select="onDocumentoGarantiaSelect"
            @documento-clear="limpiarDocumentoGarantia"
          />

          <template v-if="false && esCredito">
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

          <div v-if="!esCredito" class="flex flex-col gap-1">
            <label for="estado_entrega">Estado de entrega</label>
            <Select id="estado_entrega" v-model="form.estado_entrega" :options="estadosEntrega" :disabled="esReserva" fluid size="small" />
          </div>

          <div v-if="!esCredito" class="flex flex-col gap-1">
            <label for="fecha_entrega">Fecha de entrega</label>
            <InputText id="fecha_entrega" v-model="form.fecha_entrega" type="date" fluid size="small" />
          </div>

          <template v-if="false && esCreditoBancario">
            <div class="flex flex-col gap-1">
              <label for="referencia_bancaria">Referencia bancaria</label>
              <InputText id="referencia_bancaria" v-model="form.referencia_bancaria" placeholder="Texto libre" fluid size="small" />
            </div>

            <div class="flex flex-col gap-1">
              <label for="estado_desembolso">Estado desembolso</label>
              <Select id="estado_desembolso" v-model="form.estado_desembolso" :options="estadosDesembolso" fluid size="small" />
            </div>
          </template>

          <template v-if="false && esCreditoDirecto">
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
              </div>
              <div>
                <span class="text-xs font-semibold uppercase text-gray-500">Saldo reserva USD</span>
                <strong class="block text-gray-900">$ {{ formatPrecio(saldoDetalleUSD) }}</strong>
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
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
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
const vehiculos = ref<any[]>([]);
const clientes = ref<any[]>([]);
const ventas = ref<any[]>([]);
const documentoGarantia = ref<File | null>(null);
const documentoGarantiaNombre = ref('');
const tiposVenta = ref(['Contado', 'credito_directo', 'credito_bancario', 'Reserva']);
const estadosVenta = ref(['Registrada', 'Completada', 'en_credito']);
const estadosPago = ref(['Pendiente', 'Parcial', 'Pagado completo']);
const metodosPago = ref(['QR', 'Transferencia', 'Efectivo', 'Mixto']);
const metodosPagoMoneda = ref(['QR', 'Transferencia', 'Efectivo']);
const monedasPago = ref(['USD', 'BOB']);
const estadosEntrega = ref(['Pendiente', 'Entregado']);
const estadosDesembolso = ref(['Pendiente', 'Desembolsado']);
const frecuenciasPago = ref(['mensual', 'quincenal', 'semanal']);
const tiposReserva = ref([
  { label: 'Vehiculo en stock', value: 'stock' },
  { label: 'Vehiculo a pedido', value: 'pedido' }
]);

const form = reactive({
  id_vehiculo: null as number | null,
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
  validez_proforma_dias: 15,
  estado_venta: 'Completada',
  estado_pago: 'Pagado completo',
  metodo_pago: 'Efectivo',
  estado_entrega: 'Pendiente',
  fecha_entrega: '',
  referencia_bancaria: '',
  estado_desembolso: 'Pendiente',
  monto_reserva: 0,
  monto_inicial: 0,
  numero_cuotas: 12,
  fecha_inicio_credito: new Date().toISOString().slice(0, 10),
  frecuencia_pago: 'mensual',
  tiene_respaldo: false,
  tipo_garantia: '',
  datos_garante: '',
  tipo_reserva: 'stock',
  pedido_marca: '',
  pedido_modelo: '',
  pedido_anio: new Date().getFullYear(),
  pedido_color: '',
  pedido_version: '',
  pedido_pais_origen: '',
  pedido_proveedor: '',
  pedido_llegada_estimada: '',
  precio_pedido: 0,
  observacion: ''
});

const vehiculoEnStock = computed(() => !esReservaPedido.value && !!vehiculo.value);
const esReservaPedido = computed(() => esReserva.value && form.tipo_reserva === 'pedido');
const precioUnidad = computed(() => esReservaPedido.value ? Number(form.precio_pedido || 0) : Number(vehiculo.value?.precio_usd || vehiculo.value?.precio || 0));
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
const muestraDetallePago = computed(() => (esPagoMixto.value && !esReserva.value) || esCredito.value || esReserva.value);
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
const objetivoDetallePagoUSD = computed(() => esReserva.value ? Number(form.monto_reserva || 0) : precioTotal.value);
const pagoDetalleUSD = computed(() => pagoMixtoUSD.value);
const saldoDetalleUSD = computed(() => Math.max(objetivoDetallePagoUSD.value - pagoDetalleUSD.value, 0));
const saldoDetalleBOB = computed(() => saldoDetalleUSD.value * tipoCambio.value);
const pagoMixtoExcedeTotal = computed(() => muestraDetallePago.value && pagoDetalleUSD.value > objetivoDetallePagoUSD.value);
const estadoPagoSugerido = computed(() => {
  if (!muestraDetallePago.value) return form.estado_pago;
  if (pagoDetalleUSD.value <= 0) return 'Pendiente';
  if (pagoDetalleUSD.value >= objetivoDetallePagoUSD.value) return 'Pagado completo';
  return 'Parcial';
});
const montoInicialUSD = computed(() => {
  if (esCredito.value) {
    return pagoDetalleUSD.value;
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
const montoCuota = computed(() => Number(form.numero_cuotas || 0) > 0 ? montoFinanciado.value / Number(form.numero_cuotas || 1) : 0);
const cantidadReservadaVigente = computed(() => {
  const idVehiculo = Number(vehiculo.value?.id || form.id_vehiculo || route.query.id || 0);
  return ventas.value
    .filter((venta: any) =>
      Number(venta.id_vehiculo || 0) === idVehiculo &&
      venta.estado_venta === 'Registrada' &&
      !venta.proforma_vencida
    )
    .reduce((total: number, venta: any) => total + Number(venta.cantidad || 0), 0);
});
const disponibilidadVenta = computed(() => {
  if (esReservaPedido.value) return 1;
  return Math.max(Number(vehiculo.value?.cantidad_disponible || 0) - cantidadReservadaVigente.value, 0);
});
const etiquetaVehiculoPedido = computed(() => [form.pedido_marca, form.pedido_modelo, form.pedido_anio].filter(Boolean).join(' ') || 'Vehiculo a pedido');

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

  if (!esCreditoDirecto.value) {
    limpiarDocumentoGarantia();
  }
}, { immediate: true });

watch(() => form.tipo_reserva, (tipo) => {
  if (tipo === 'pedido') {
    vehiculo.value = null;
    form.id_vehiculo = null;
    form.cantidad = 1;
    form.estado_venta = 'Importando';
    form.estado_entrega = 'Pendiente';
  } else if (esReserva.value) {
    form.estado_venta = 'Registrada';
  }
});

watch(() => form.id_vehiculo, (id) => {
  if (esReservaPedido.value) return;
  vehiculo.value = vehiculos.value.find((item: any) => Number(item.id) === Number(id)) || vehiculo.value;
});

watch(precioTotal, () => {
  if (Number(form.monto_reserva || 0) > precioTotal.value) {
    form.monto_reserva = precioTotal.value;
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
      const [resClientes, resVehiculos, resVentas] = await Promise.all([
        $fetch(server.HOST + '/api/v1/clientes', { method: 'GET' }),
        $fetch(server.HOST + '/api/v1/vehiculos', { method: 'GET' }),
        $fetch(server.HOST + '/api/v1/ventas', { method: 'GET' })
      ]);
      vehiculos.value = Array.isArray(resVehiculos) ? resVehiculos.map((item: any) => ({
        ...item,
        nombreCompleto: etiquetaVehiculo(item)
      })) : [];
      ventas.value = Array.isArray(resVentas) ? resVentas : [];
      const clientesActivos = Array.isArray(resClientes) ? resClientes.filter((cliente: any) => cliente.estado === 'Activo') : [];
      clientes.value = clientesActivos.map((cliente: any) => ({
        ...cliente,
        nombreCompleto: `${cliente.nombre || ''} ${cliente.apellido || ''} - ${cliente.ci || 'Sin CI'}`.trim()
      }));
      form.tipo_venta = 'Reserva';
      form.tipo_reserva = 'pedido';
      return;
    }

    const [resVehiculo, resClientes, resVentas] = await Promise.all([
      $fetch(server.HOST + '/api/v1/vehiculos/' + idVehiculo, { method: 'GET' }),
      $fetch(server.HOST + '/api/v1/clientes', { method: 'GET' }),
      $fetch(server.HOST + '/api/v1/ventas', { method: 'GET' })
    ]);

    vehiculo.value = resVehiculo;
    form.id_vehiculo = Number(idVehiculo);
    vehiculos.value = [{
      ...resVehiculo,
      nombreCompleto: etiquetaVehiculo(resVehiculo)
    }];
    ventas.value = Array.isArray(resVentas) ? resVentas : [];
    const clientesActivos = Array.isArray(resClientes) ? resClientes.filter((cliente: any) => cliente.estado === 'Activo') : [];
    clientes.value = clientesActivos.map((cliente: any) => ({
      ...cliente,
      nombreCompleto: `${cliente.nombre || ''} ${cliente.apellido || ''} - ${cliente.ci || 'Sin CI'}`.trim()
    }));
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

async function registrarVenta() {
  if (!form.id_cliente || (!vehiculo.value && !esReservaPedido.value)) {
    toast.add({ severity: 'warn', summary: 'Seleccione cliente y vehiculo', life: 3000 });
    return;
  }
  if (esReservaPedido.value && !validarReservaPedido()) {
    return;
  }
  if (!esReservaPedido.value && disponibilidadVenta.value <= 0) {
    toast.add({ severity: 'warn', summary: 'Sin disponibilidad vendible', detail: 'El stock esta reservado por proformas vigentes.', life: 4000 });
    return;
  }
  if (!esReservaPedido.value && Number(form.cantidad || 0) > disponibilidadVenta.value) {
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
  if (esCredito.value) {
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
  formData.append('id_vehiculo', esReservaPedido.value ? '' : String(Number(form.id_vehiculo || route.query.id)));
  formData.append('id_usuario', String(userId || 0));
  formData.append('fecha', form.fecha);
  formData.append('tipo_venta', form.tipo_venta);
  formData.append('cantidad', String(Number(form.cantidad)));
  formData.append('tipo_cambio', String(tipoCambio.value));
  formData.append('pago_usd', String(muestraDetallePago.value ? Number(pagoMixtoUSDDirecto.value || 0) : 0));
  formData.append('pago_bob', String(muestraDetallePago.value ? Number(pagoMixtoBOBDirecto.value || 0) : 0));
  formData.append('pagos', JSON.stringify(muestraDetallePago.value ? detallePagoPayload() : []));
  formData.append('validez_proforma_dias', String(Number(form.validez_proforma_dias || 15)));
  formData.append('estado_venta', esReservaPedido.value ? 'Importando' : (esReserva.value ? 'Registrada' : form.estado_venta));
  formData.append('estado_pago', form.estado_pago);
  formData.append('metodo_pago', esCredito.value || esReserva.value ? resumenDetallePago() : form.metodo_pago);
  formData.append('estado_entrega', esReserva.value ? 'Pendiente' : form.estado_entrega);
  formData.append('fecha_entrega', esReserva.value ? '' : form.fecha_entrega);
  formData.append('monto_reserva', String(esReserva.value ? Number(form.monto_reserva || 0) : 0));
  formData.append('monto_inicial', String(esCredito.value ? Number(montoInicialUSD.value || 0) : 0));
  formData.append('numero_cuotas', String(esCredito.value ? Number(form.numero_cuotas || 0) : 0));
  formData.append('fecha_inicio_credito', esCredito.value ? form.fecha_inicio_credito : '');
  formData.append('frecuencia_pago', esCredito.value ? form.frecuencia_pago : '');
  formData.append('tiene_respaldo', String(esCreditoDirecto.value ? form.tiene_respaldo : false));
  formData.append('tipo_garantia', esCreditoDirecto.value ? form.tipo_garantia : '');
  formData.append('datos_garante', esCreditoDirecto.value ? form.datos_garante : '');
  formData.append('referencia_bancaria', esCreditoBancario.value ? form.referencia_bancaria : '');
  formData.append('estado_desembolso', esCreditoBancario.value ? form.estado_desembolso : '');
  formData.append('tipo_reserva', esReserva.value ? form.tipo_reserva : 'stock');
  formData.append('pedido_marca', esReservaPedido.value ? form.pedido_marca : '');
  formData.append('pedido_modelo', esReservaPedido.value ? form.pedido_modelo : '');
  formData.append('pedido_anio', esReservaPedido.value ? String(Number(form.pedido_anio || 0)) : '0');
  formData.append('pedido_color', esReservaPedido.value ? form.pedido_color : '');
  formData.append('pedido_version', esReservaPedido.value ? form.pedido_version : '');
  formData.append('pedido_pais_origen', esReservaPedido.value ? form.pedido_pais_origen : '');
  formData.append('pedido_proveedor', esReservaPedido.value ? form.pedido_proveedor : '');
  formData.append('pedido_llegada_estimada', esReservaPedido.value ? form.pedido_llegada_estimada : '');
  formData.append('precio_pedido', esReservaPedido.value ? String(Number(form.precio_pedido || 0)) : '0');
  formData.append('observacion', form.observacion);
  if (esCreditoDirecto.value && documentoGarantia.value) {
    formData.append('documento_garantia', documentoGarantia.value);
  }
  return formData;
}

function validarReservaPedido() {
  if (!form.pedido_marca.trim() || !form.pedido_modelo.trim() || !form.pedido_pais_origen.trim()) {
    toast.add({ severity: 'warn', summary: 'Complete marca, modelo y pais de origen', life: 3000 });
    return false;
  }
  if (Number(form.pedido_anio || 0) < 1900) {
    toast.add({ severity: 'warn', summary: 'Ingrese un anio valido', life: 3000 });
    return false;
  }
  if (Number(form.precio_pedido || 0) <= 0) {
    toast.add({ severity: 'warn', summary: 'Ingrese el precio pactado', life: 3000 });
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
