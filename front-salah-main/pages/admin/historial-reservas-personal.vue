<template>
  <HistorialReservas completar-path="/admin/completar-reserva" />
</template>

<script setup lang="ts">
import HistorialReservas from '~/components/ventas/HistorialReservas.vue';

definePageMeta({ layout: 'menu-admin' });
</script>
