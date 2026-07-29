<template>
  <VentaVehiculoForm
    catalogo-path="/admin/catalogo-vehiculos"
    historial-path="/admin/historial-ventas"
  />
</template>

<script setup lang="ts">
import VentaVehiculoForm from '~/components/ventas/VentaVehiculoForm.vue';

definePageMeta({ layout: 'menu-admin' });
</script>
