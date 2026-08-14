<template>
  <RecibirPedidoForm historial-path="/admin/pedidos" />
</template>

<script setup lang="ts">
import RecibirPedidoForm from '~/components/ventas/RecibirPedidoForm.vue';

definePageMeta({ layout: 'menu-admin' });
</script>
