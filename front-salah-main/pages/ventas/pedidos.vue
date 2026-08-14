<template>
  <PedidosVehiculos
    recibir-path="/ventas/recibir-pedido"
    completar-path="/ventas/completar-pedido"
  />
</template>

<script setup lang="ts">
import PedidosVehiculos from '~/components/ventas/PedidosVehiculos.vue';

definePageMeta({ layout: 'menu-ventas' });
</script>
