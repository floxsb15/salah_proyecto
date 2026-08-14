<template>
  <RecibirPedidoForm historial-path="/ventas/pedidos" />
</template>

<script setup lang="ts">
import RecibirPedidoForm from '~/components/ventas/RecibirPedidoForm.vue';

definePageMeta({ layout: 'menu-ventas' });
</script>
