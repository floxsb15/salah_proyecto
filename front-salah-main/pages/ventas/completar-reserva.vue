<template>
  <CompletarReservaForm historial-path="/ventas/historial-reservas" />
</template>

<script setup lang="ts">
import CompletarReservaForm from '~/components/ventas/CompletarReservaForm.vue';

definePageMeta({ layout: 'menu-ventas' });
</script>
