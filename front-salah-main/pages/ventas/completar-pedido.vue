<template>
  <CompletarPedidoForm historial-path="/ventas/pedidos" />
</template>

<script setup lang="ts">
import CompletarPedidoForm from '~/components/ventas/CompletarPedidoForm.vue';

definePageMeta({ layout: 'menu-ventas' });
</script>
