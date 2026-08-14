<template>
  <CompletarPedidoForm historial-path="/admin/pedidos" />
</template>

<script setup lang="ts">
import CompletarPedidoForm from '~/components/ventas/CompletarPedidoForm.vue';

definePageMeta({ layout: 'menu-admin' });
</script>
