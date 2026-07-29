<template>
  <div>
    <Select
      id="tipo" name="tipo"
      v-model="tipo"
      placeholder="Seleccione un tipo" 
      :options="Tipos"
      fluid size="small" />

    <div v-if="tipo === 'Persona particular'">
      <p>CI: <input v-model="ci" /></p>
    </div>

    <div v-if="tipo === 'Empresa'">
      <p>NIT: <input v-model="nit" /></p>
      <p>Nombre del encargado: <input v-model="nombre_encargado" /></p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const tipo = ref('Persona particular');
const ci = ref('');
const nit = ref('');
const nombre_encargado = ref('');

const Tipos = ref(['Persona particular', 'Empresa']);
</script>
