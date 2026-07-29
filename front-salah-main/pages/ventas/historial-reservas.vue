<template>
  <HistorialReservas completar-path="/ventas/completar-reserva" />
</template>

<script setup lang="ts">
import HistorialReservas from '~/components/ventas/HistorialReservas.vue';

definePageMeta({ layout: 'menu-ventas' });
</script>
