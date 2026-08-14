<template>
  <PedidosVehiculos
    general
    recibir-path="/admin/recibir-pedido"
    completar-path="/admin/completar-pedido"
  />
</template>

<script setup lang="ts">
import PedidosVehiculos from '~/components/ventas/PedidosVehiculos.vue';

definePageMeta({ layout: 'menu-admin' });
</script>
