<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-100 p-4">
    <div class="w-full max-w-md rounded-xl border border-gray-200 bg-white p-8 shadow-lg">
      <h1 class="mb-2 text-2xl font-bold text-gray-900">Cambiar contraseña</h1>
      <p class="mb-6 text-sm text-gray-600">Debes reemplazar la contraseña heredada antes de continuar.</p>

      <form class="flex flex-col gap-4" @submit.prevent="cambiarContrasena">
        <label class="flex flex-col gap-1 text-sm text-gray-700">
          Contraseña actual
          <Password v-model="actual" :feedback="false" toggle-mask fluid />
        </label>
        <label class="flex flex-col gap-1 text-sm text-gray-700">
          Nueva contraseña
          <Password v-model="nueva" toggle-mask fluid />
        </label>
        <label class="flex flex-col gap-1 text-sm text-gray-700">
          Confirmar contraseña
          <Password v-model="confirmacion" :feedback="false" toggle-mask fluid />
        </label>
        <Message v-if="error" severity="error" size="small">{{ error }}</Message>
        <Button type="submit" label="Guardar contraseña" :loading="guardando" />
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import Button from 'primevue/button';
import Message from 'primevue/message';
import Password from 'primevue/password';

const actual = ref('');
const nueva = ref('');
const confirmacion = ref('');
const error = ref('');
const guardando = ref(false);

async function cambiarContrasena() {
  error.value = '';
  if (nueva.value.length < 12) {
    error.value = 'La nueva contraseña debe contener al menos 12 caracteres.';
    return;
  }
  if (nueva.value !== confirmacion.value) {
    error.value = 'Las contraseñas no coinciden.';
    return;
  }

  guardando.value = true;
  try {
    await $fetch('/api/v1/me/password', {
      method: 'PATCH',
      body: { actual: actual.value, nueva: nueva.value },
    });
    localStorage.removeItem('user');
    await navigateTo('/login');
  } catch (err: any) {
    error.value = typeof err?.data === 'string' ? err.data : 'No se pudo cambiar la contraseña.';
  } finally {
    guardando.value = false;
  }
}
</script>
