<template>
  <Toast />
  <div class="flex items-center justify-center min-h-screen bg-gray-100 p-4">
    <div class="flex flex-col md:flex-row bg-white rounded-xl border border-gray-200 shadow-lg overflow-hidden max-w-4xl w-full">
      <!-- Image/Illustration Section (Hidden on mobile, shown on desktop) -->
      <div class="hidden md:flex md:w-1/2 bg-black items-center justify-center p-8">
        <!-- Placeholder for an image or illustration -->
        <img src="/logosalah.png" alt="Logo Salah" class="max-w-full h-auto">
      </div>

      <!-- Login Form Section -->
      <div class="w-full md:w-1/2 p-8 flex flex-col justify-center">
        <h2 class="text-3xl font-bold text-center text-gray-800 mb-6">Iniciar Sesión</h2>

        <form method="post" action="/api/v1/login" @submit.prevent="iniciarSesion" class="flex flex-col gap-4">

          <!-- Username Input -->
          <div class="flex flex-col gap-1">
            <label for="usuario" class="text-sm text-gray-700"> Usuario </label>
            <InputGroup>
              <InputGroupAddon>
                <i class="pi pi-user"></i>
              </InputGroupAddon>
              <InputText 
                id="usuario" name="usuario" autocomplete="username" placeholder="Ingrese su nombre de usuario"
                v-model="initialValues.usuario"
                :class="{ 'p-invalid': errors.usuario }" 
                fluid size="small"/>
            </InputGroup>
            <Message
              v-if="errors.usuario"
              severity="error" size="small" variant="simple">
              {{ errors.usuario }}
            </Message>
          </div>

          <!-- Password Input -->
          <div class="flex flex-col gap-1">
            <label for="contra" class="text-sm text-gray-700"> Contraseña </label>
            <InputGroup>
              <InputGroupAddon>
                <i class="pi pi-lock"></i>
              </InputGroupAddon>
              <Password
                id="contra" name="contra" autocomplete="current-password" placeholder="Ingrese su contraseña"
                v-model="initialValues.contra"
                :class="{ 'p-invalid': errors.contra }" 
                :feedback="false" toggleMask
                fluid size="small"/>
            </InputGroup>
            <Message
              v-if="errors.contra"
              severity="error" size="small" variant="simple">
              {{ errors.contra }}
            </Message>
          </div>

          <!-- Login Button -->
          <Button type="submit" label="Iniciar Sesión" class="w-full mt-4"
            :pt="{
              root: { class: 'bg-emerald-500 hover:bg-emerald-600 border-emerald-500 hover:border-emerald-600 text-white' }
            }"
          />
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import { z } from 'zod';
import { server } from '~/server/server';

import InputText from 'primevue/inputtext';
import Password from 'primevue/password';
import Button from 'primevue/button';
import Toast from 'primevue/toast';
import Message from 'primevue/message';
import InputGroup from 'primevue/inputgroup';
import InputGroupAddon from 'primevue/inputgroupaddon';
import { useToast } from 'primevue/usetoast';

const initialValues = reactive({ usuario: '', contra: '' });
const router = useRouter();
const toast = useToast();

const errors = reactive({ usuario: '', contra: '', });

const loginSchema = z.object({
  usuario: z.string().min(1, { message: 'Usuario requerido.' }),
  contra: z.string().min(1, { message: 'Contraseña requerida.' }),
});

async function iniciarSesion() {
  // Clear previous errors
  errors.usuario = '';
  errors.contra = '';

  const result = loginSchema.safeParse(initialValues);

  if (!result.success) {
    result.error.errors.forEach(err => {
      if (err.path[0] === 'usuario') {
        errors.usuario = err.message;
      } else if (err.path[0] === 'contra') {
        errors.contra = err.message;
      }
    });
    return;
  }

  try {
    const response:any = await $fetch(server.HOST + '/api/v1/login', {
      method: 'POST',
      body: JSON.stringify(initialValues),
      headers: { 'Content-Type': 'application/json' }
    });


    localStorage.setItem('user', JSON.stringify(response));
    if (response.must_change_password) {
      router.push('/cambiar-contrasena');
    } else if(response.rol === 'admin'){
      router.push('/admin/dashboard');
    } else if(['encargado de ventas', 'vendedor'].includes(response.rol)){
      router.push('/ventas/catalogo-vehiculos');
    } else {
      toast.add({ severity: 'error', summary: 'Rol no reconocido', detail: 'El rol del usuario no es válido.', life: 3000 });
    }
  } catch (err: any) {
    console.error(err);
    const errorMessage = err.response?._data || 'Credenciales incorrectas';
    toast.add({ severity: 'error', summary: 'Error de autenticación', detail: errorMessage, life: 3000 });
  }
}
</script>

