<template>
  <div class="flex h-dvh bg-gray-100">
    <!-- SIDEBAR -->
    <aside class="w-60 bg-neutral-950 text-white flex flex-col">
      <!-- Logo -->
      <div class="h-16 flex items-center justify-center gap-3 border-b border-neutral-800 px-3">
        <img src="/logosalah.png" alt="SALAH MOTORS" class="h-10 w-10 object-contain">
        <span class="text-base font-semibold tracking-wide text-yellow-400">SALAH MOTORS</span>
      </div>

      <!-- Menu -->
      <nav class="flex-1 p-2 space-y-2">
        <router-link to="/admin/dashboard" custom v-slot="{ href, navigate, isActive }">
          <a :href="href" @click="navigate" class="flex items-center gap-3 rounded-lg px-3 py-2 transition-colors" 
             :class="isActive ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'">
            <i class="pi pi-chart-bar text-yellow-400" :class="{ '!text-black': isActive }"></i>
            <span>Dashboard</span>
          </a>
        </router-link>
        <router-link to="/admin/usuarios" custom v-slot="{ href, navigate, isActive }">
          <a :href="href" @click="navigate" class="flex items-center gap-3 rounded-lg px-3 py-2 transition-colors" 
             :class="isActive ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'">
            <i class="pi pi-users text-yellow-400" :class="{ '!text-black': isActive }"></i>
            <span>Usuarios</span>
          </a>
        </router-link>
        <router-link to="/admin/clientes" custom v-slot="{ href, navigate, isActive }">
          <a :href="href" @click="navigate" class="flex items-center gap-3 rounded-lg px-3 py-2 transition-colors" 
             :class="isActive ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'">
            <i class="pi pi-id-card text-yellow-400" :class="{ '!text-black': isActive }"></i>
            <span>Clientes</span>
          </a>
        </router-link>
        <router-link to="/admin/categorias-vehiculos" custom v-slot="{ href, navigate, isActive }">
          <a :href="href" @click="navigate" class="flex items-center gap-3 rounded-lg px-3 py-2 transition-colors" 
             :class="isActive ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'">
            <i class="pi pi-tags text-yellow-400" :class="{ '!text-black': isActive }"></i>
            <span>Categorias Vehiculo</span>
          </a>
        </router-link>
        <router-link to="/admin/vehiculos" custom v-slot="{ href, navigate, isActive }">
          <a :href="href" @click="navigate" class="flex items-center gap-3 rounded-lg px-3 py-2 transition-colors" 
             :class="isActive ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'">
            <i class="pi pi-car text-yellow-400" :class="{ '!text-black': isActive }"></i>
            <span>Vehiculos</span>
          </a>
        </router-link>
        <router-link to="/admin/catalogo-vehiculos" custom v-slot="{ href, navigate, isActive }">
          <a :href="href" @click="navigate" class="flex items-center gap-3 rounded-lg px-3 py-2 transition-colors" 
             :class="isActive ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'">
            <i class="pi pi-images text-yellow-400" :class="{ '!text-black': isActive }"></i>
            <span>Catalogo Vehiculo</span>
          </a>
        </router-link>
        <div>
          <button
            type="button"
            class="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-left transition-colors"
            :class="creditosOpen || route.path.includes('/admin/creditos') ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'"
            @click="creditosOpen = !creditosOpen"
          >
            <i class="pi pi-credit-card text-yellow-400" :class="{ '!text-black': creditosOpen || route.path.includes('/admin/creditos') }"></i>
            <span class="flex-1">Creditos</span>
            <i class="pi text-xs" :class="creditosOpen ? 'pi-chevron-up' : 'pi-chevron-down'"></i>
          </button>
          <div v-if="creditosOpen" class="mt-1 space-y-1 pl-6">
            <router-link to="/admin/creditos" custom v-slot="{ href, navigate, isActive }">
              <a :href="href" @click="navigate" class="flex items-center gap-2 rounded-lg px-3 py-2 text-sm transition-colors"
                 :class="isActive ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'">
                <i class="pi pi-list text-yellow-400" :class="{ '!text-black': isActive }"></i>
                <span>Historial General</span>
              </a>
            </router-link>
            <router-link to="/admin/creditos-personal" custom v-slot="{ href, navigate, isActive }">
              <a :href="href" @click="navigate" class="flex items-center gap-2 rounded-lg px-3 py-2 text-sm transition-colors"
                 :class="isActive ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'">
                <i class="pi pi-user text-yellow-400" :class="{ '!text-black': isActive }"></i>
                <span>Historial Personal</span>
              </a>
            </router-link>
          </div>
        </div>
        <div>
          <button
            type="button"
            class="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-left transition-colors"
            :class="ventasOpen || route.path.includes('/admin/historial-ventas') ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'"
            @click="ventasOpen = !ventasOpen"
          >
            <i class="pi pi-history text-yellow-400" :class="{ '!text-black': ventasOpen || route.path.includes('/admin/historial-ventas') }"></i>
            <span class="flex-1">Ventas</span>
            <i class="pi text-xs" :class="ventasOpen ? 'pi-chevron-up' : 'pi-chevron-down'"></i>
          </button>
          <div v-if="ventasOpen" class="mt-1 space-y-1 pl-6">
            <router-link to="/admin/historial-ventas" custom v-slot="{ href, navigate, isActive }">
              <a :href="href" @click="navigate" class="flex items-center gap-2 rounded-lg px-3 py-2 text-sm transition-colors"
                 :class="isActive ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'">
                <i class="pi pi-list text-yellow-400" :class="{ '!text-black': isActive }"></i>
                <span>Historial General</span>
              </a>
            </router-link>
            <router-link to="/admin/historial-ventas-personal" custom v-slot="{ href, navigate, isActive }">
              <a :href="href" @click="navigate" class="flex items-center gap-2 rounded-lg px-3 py-2 text-sm transition-colors"
                 :class="isActive ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'">
                <i class="pi pi-user text-yellow-400" :class="{ '!text-black': isActive }"></i>
                <span>Historial Personal</span>
              </a>
            </router-link>
          </div>
        </div>
        <div>
          <button
            type="button"
            class="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-left transition-colors"
            :class="reservasOpen || route.path.includes('/admin/historial-reservas') ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'"
            @click="reservasOpen = !reservasOpen"
          >
            <i class="pi pi-bookmark text-yellow-400" :class="{ '!text-black': reservasOpen || route.path.includes('/admin/historial-reservas') }"></i>
            <span class="flex-1">Reservas</span>
            <i class="pi text-xs" :class="reservasOpen ? 'pi-chevron-up' : 'pi-chevron-down'"></i>
          </button>
          <div v-if="reservasOpen" class="mt-1 space-y-1 pl-6">
            <router-link to="/admin/historial-reservas" custom v-slot="{ href, navigate, isActive }">
              <a :href="href" @click="navigate" class="flex items-center gap-2 rounded-lg px-3 py-2 text-sm transition-colors"
                 :class="isActive ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'">
                <i class="pi pi-list text-yellow-400" :class="{ '!text-black': isActive }"></i>
                <span>Historial General</span>
              </a>
            </router-link>
            <router-link to="/admin/historial-reservas-personal" custom v-slot="{ href, navigate, isActive }">
              <a :href="href" @click="navigate" class="flex items-center gap-2 rounded-lg px-3 py-2 text-sm transition-colors"
                 :class="isActive ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'">
                <i class="pi pi-user text-yellow-400" :class="{ '!text-black': isActive }"></i>
                <span>Historial Personal</span>
              </a>
            </router-link>
          </div>
        </div>
      </nav>
    </aside>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col">
      <!-- HEADER -->
      <header class="h-16 flex items-center justify-between gap-4 px-4 bg-neutral-950 text-white sticky top-0 z-50 border-b border-neutral-800">
        <!-- Search Bar -->
        <div class="relative w-full max-w-md">
          <i class="pi pi-search absolute left-3 top-1/2 -translate-y-1/2 text-yellow-400"></i>
          <input v-model="topSearch" type="text" placeholder="Buscar clientes, vehiculos o usuarios" class="h-10 w-full rounded-md border border-neutral-700 bg-neutral-900 pl-10 pr-10 text-sm text-white outline-none transition-colors placeholder:text-neutral-400 focus:border-yellow-400">
          <button
            v-if="topSearch"
            type="button"
            class="absolute right-2 top-1/2 inline-flex h-7 w-7 -translate-y-1/2 items-center justify-center rounded-md text-neutral-300 transition-colors hover:bg-neutral-800 hover:text-yellow-400"
            aria-label="Limpiar busqueda"
            @click="topSearch = ''"
          >
            <i class="pi pi-times text-xs"></i>
          </button>
        </div>
        
        <!-- User Info and Logout -->
        <client-only>
          <div class="flex items-center gap-3">
            <div class="hidden items-center gap-2 lg:flex">
              <NuxtLink to="/admin/dashboard" class="inline-flex h-9 items-center gap-2 rounded-md px-3 text-sm text-white transition-colors hover:bg-red-900/60">
                <i class="pi pi-chart-bar text-yellow-400"></i>
                <span>Dashboard</span>
              </NuxtLink>
              <NuxtLink to="/admin/clientes" class="inline-flex h-9 items-center gap-2 rounded-md px-3 text-sm text-white transition-colors hover:bg-red-900/60">
                <i class="pi pi-id-card text-yellow-400"></i>
                <span>Clientes</span>
              </NuxtLink>
              <NuxtLink to="/admin/vehiculos" class="inline-flex h-9 items-center gap-2 rounded-md px-3 text-sm text-white transition-colors hover:bg-red-900/60">
                <i class="pi pi-car text-yellow-400"></i>
                <span>Vehiculos</span>
              </NuxtLink>
              <NuxtLink to="/admin/catalogo-vehiculos" class="inline-flex h-9 items-center gap-2 rounded-md px-3 text-sm text-white transition-colors hover:bg-red-900/60">
                <i class="pi pi-shopping-cart text-yellow-400"></i>
                <span>Ventas</span>
              </NuxtLink>
              <NuxtLink to="/admin/historial-ventas" class="inline-flex h-9 items-center gap-2 rounded-md px-3 text-sm text-white transition-colors hover:bg-red-900/60">
                <i class="pi pi-history text-yellow-400"></i>
                <span>Historial</span>
              </NuxtLink>
              <NuxtLink to="/admin/historial-reservas" class="inline-flex h-9 items-center gap-2 rounded-md px-3 text-sm text-white transition-colors hover:bg-red-900/60">
                <i class="pi pi-bookmark text-yellow-400"></i>
                <span>Reservas</span>
              </NuxtLink>
              <NuxtLink to="/admin/creditos" class="inline-flex h-9 items-center gap-2 rounded-md px-3 text-sm text-white transition-colors hover:bg-red-900/60">
                <i class="pi pi-credit-card text-yellow-400"></i>
                <span>Creditos</span>
              </NuxtLink>
            </div>
            <SalesNotifications :userData="userData" />
            <div class="flex items-center gap-3 rounded-xl border border-neutral-800 bg-neutral-900/80 px-2 py-1 transition-colors hover:border-yellow-400/60 hover:bg-neutral-800 cursor-pointer" @click="isProfileModalVisible = true">
                <div class="h-10 w-10 overflow-hidden rounded-full border-2 border-yellow-400 bg-yellow-400 text-black shadow-[0_0_0_3px_rgba(255,215,0,0.12)]">
                  <img v-if="userData.foto && userData.foto !== 'N/A'" :src="userData.foto" alt="Foto de usuario" class="h-full w-full object-cover">
                  <div v-else class="flex h-full w-full items-center justify-center text-sm font-extrabold">
                    {{ userInitials }}
                  </div>
                </div>
                <div class="hidden text-sm sm:block">
                    <div class="font-semibold leading-5">{{ userData.nombre || 'Usuario' }}</div>
                    <div class="text-xs text-neutral-300">{{ userData.rol }}</div>
                </div>
                <i class="pi pi-chevron-down hidden text-xs text-yellow-400 sm:block"></i>
            </div>
            <button
              type="button"
              @click="cerrarSesion"
              class="inline-flex h-10 items-center gap-2 rounded-md px-3 text-sm transition-colors hover:bg-red-900/60"
            >
              <i class="pi pi-sign-out text-yellow-400"></i>
              <span>Cerrar Sesión</span>
            </button>
          </div>
        </client-only>
      </header>

      <!-- PAGE CONTENT -->
      <main class="flex-1 p-4 overflow-y-auto">
        <div class="bg-white p-4 rounded-lg border border-gray-200 shadow-sm">
          <NuxtPage />
        </div>
      </main>
    </div>

    <!-- User Profile Modal -->
    <UserProfileModal v-model:visible="isProfileModalVisible" :userData="userData" />
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { server } from '~/server/server';
import UserProfileModal from '~/components/shared/UserProfileModal.vue';
import SalesNotifications from '~/components/shared/SalesNotifications.vue';

const userData = ref({
  id: 0,
  nombre: '',
  apellido: '',
  ci: '',
  celular: '',
  usuario: '',
  rol: '',
  estado: '',
  direccion: '',
  foto: '',
});

const isProfileModalVisible = ref(false);
const topSearch = useTopSearch();
const route = useRoute();
const ventasOpen = ref(false);
const reservasOpen = ref(false);
const creditosOpen = ref(false);
const userInitials = computed(() => {
  const nombre = userData.value.nombre?.trim()?.[0] || 'U';
  const apellido = userData.value.apellido?.trim()?.[0] || '';
  return `${nombre}${apellido}`.toUpperCase();
});

onMounted(async () => {
  const user = localStorage.getItem('user');
  if (user) {
    try {
      const res:any = await $fetch(server.HOST + '/api/v1/me', { method: 'GET' });
      Object.assign(userData.value, res);
    } catch (err) {
      console.error('Error fetching user data:', err);
    }
  }
});

async function cerrarSesion() {
  try {
    await $fetch('/api/v1/logout', { method: 'POST' });
  } finally {
    localStorage.removeItem('user');
    await navigateTo('/login');
  }
}
</script>
