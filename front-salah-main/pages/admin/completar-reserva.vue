<template>
  <CompletarReservaForm historial-path="/admin/historial-reservas" />
</template>

<script setup lang="ts">
import CompletarReservaForm from '~/components/ventas/CompletarReservaForm.vue';

definePageMeta({ layout: 'menu-admin' });
</script>
