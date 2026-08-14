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
      <nav class="flex-1 overflow-y-auto p-2 space-y-2">
        <router-link to="/admin/dashboard" custom v-slot="{ href, navigate, isActive }">
          <a :href="href" @click="navigate" class="flex items-center gap-3 rounded-lg px-3 py-2 transition-colors" 
             :class="isActive ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'">
            <i class="pi pi-chart-bar text-yellow-400" :class="{ '!text-black': isActive }"></i>
            <span>Dashboard</span>
          </a>
        </router-link>
        <router-link v-if="!isContador" to="/admin/usuarios" custom v-slot="{ href, navigate, isActive }">
          <a :href="href" @click="navigate" class="flex items-center gap-3 rounded-lg px-3 py-2 transition-colors" 
             :class="isActive ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'">
            <i class="pi pi-users text-yellow-400" :class="{ '!text-black': isActive }"></i>
            <span>Usuarios</span>
          </a>
        </router-link>
        <router-link v-if="!isContador" to="/admin/clientes" custom v-slot="{ href, navigate, isActive }">
          <a :href="href" @click="navigate" class="flex items-center gap-3 rounded-lg px-3 py-2 transition-colors" 
             :class="isActive ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'">
            <i class="pi pi-id-card text-yellow-400" :class="{ '!text-black': isActive }"></i>
            <span>Clientes</span>
          </a>
        </router-link>

        <div v-if="inventarioItems.length">
          <button
            type="button"
            class="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-left transition-colors"
            :class="inventarioOpen || isSectionActive(inventarioItems) ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'"
            @click="inventarioOpen = !inventarioOpen"
          >
            <i class="pi pi-box text-yellow-400" :class="{ '!text-black': inventarioOpen || isSectionActive(inventarioItems) }"></i>
            <span class="flex-1">Inventario</span>
            <i class="pi text-xs" :class="inventarioOpen ? 'pi-chevron-up' : 'pi-chevron-down'"></i>
          </button>
          <div v-if="inventarioOpen" class="mt-1 space-y-1 pl-6">
            <router-link v-for="item in inventarioItems" :key="item.to" :to="item.to" custom v-slot="{ href, navigate, isActive }">
              <a :href="href" @click="navigate" class="flex items-center gap-2 rounded-lg px-3 py-2 text-sm transition-colors"
                 :class="isActive ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'">
                <i :class="[item.icon, isActive ? '!text-black' : 'text-yellow-400']"></i>
                <span>{{ item.label }}</span>
              </a>
            </router-link>
          </div>
        </div>

        <div v-if="ventasItems.length">
          <button
            type="button"
            class="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-left transition-colors"
            :class="ventasOpen || isSectionActive(ventasItems) ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'"
            @click="ventasOpen = !ventasOpen"
          >
            <i class="pi pi-history text-yellow-400" :class="{ '!text-black': ventasOpen || isSectionActive(ventasItems) }"></i>
            <span class="flex-1">Ventas</span>
            <i class="pi text-xs" :class="ventasOpen ? 'pi-chevron-up' : 'pi-chevron-down'"></i>
          </button>
          <div v-if="ventasOpen" class="mt-1 space-y-1 pl-6">
            <router-link v-for="item in ventasItems" :key="item.to" :to="item.to" custom v-slot="{ href, navigate, isActive }">
              <a :href="href" @click="navigate" class="flex items-center gap-2 rounded-lg px-3 py-2 text-sm transition-colors"
                 :class="isActive ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'">
                <i :class="[item.icon, isActive ? '!text-black' : 'text-yellow-400']"></i>
                <span>{{ item.label }}</span>
              </a>
            </router-link>
          </div>
        </div>

        <div v-if="reservasItems.length">
          <button
            type="button"
            class="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-left transition-colors"
            :class="reservasOpen || isSectionActive(reservasItems) ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'"
            @click="reservasOpen = !reservasOpen"
          >
            <i class="pi pi-bookmark text-yellow-400" :class="{ '!text-black': reservasOpen || isSectionActive(reservasItems) }"></i>
            <span class="flex-1">Reservas</span>
            <i class="pi text-xs" :class="reservasOpen ? 'pi-chevron-up' : 'pi-chevron-down'"></i>
          </button>
          <div v-if="reservasOpen" class="mt-1 space-y-1 pl-6">
            <router-link v-for="item in reservasItems" :key="item.to" :to="item.to" custom v-slot="{ href, navigate, isActive }">
              <a :href="href" @click="navigate" class="flex items-center gap-2 rounded-lg px-3 py-2 text-sm transition-colors"
                 :class="isActive ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'">
                <i :class="[item.icon, isActive ? '!text-black' : 'text-yellow-400']"></i>
                <span>{{ item.label }}</span>
              </a>
            </router-link>
          </div>
        </div>

        <div v-if="creditosItems.length">
          <button
            type="button"
            class="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-left transition-colors"
            :class="creditosOpen || isSectionActive(creditosItems) ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'"
            @click="creditosOpen = !creditosOpen"
          >
            <i class="pi pi-credit-card text-yellow-400" :class="{ '!text-black': creditosOpen || isSectionActive(creditosItems) }"></i>
            <span class="flex-1">Creditos</span>
            <i class="pi text-xs" :class="creditosOpen ? 'pi-chevron-up' : 'pi-chevron-down'"></i>
          </button>
          <div v-if="creditosOpen" class="mt-1 space-y-1 pl-6">
            <router-link v-for="item in creditosItems" :key="item.to" :to="item.to" custom v-slot="{ href, navigate, isActive }">
              <a :href="href" @click="navigate" class="flex items-center gap-2 rounded-lg px-3 py-2 text-sm transition-colors"
                 :class="isActive ? 'bg-yellow-400 text-black' : 'text-white hover:bg-red-900/60 hover:text-white'">
                <i :class="[item.icon, isActive ? '!text-black' : 'text-yellow-400']"></i>
                <span>{{ item.label }}</span>
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
const inventarioOpen = ref(false);
const ventasOpen = ref(false);
const reservasOpen = ref(false);
const creditosOpen = ref(false);
const isContador = computed(() => userData.value.rol === 'contador');
const isAdmin = computed(() => userData.value.rol === 'admin');
const inventarioItems = computed(() => {
  const items: Array<{ to: string; label: string; icon: string }> = [];

  if (!isContador.value) {
    items.push(
      { to: '/admin/vehiculos', label: 'Vehiculos', icon: 'pi pi-car' },
      { to: '/admin/categorias-vehiculos', label: 'Categorias Vehiculo', icon: 'pi pi-tags' },
      { to: '/admin/catalogo-vehiculos', label: 'Catalogo Vehiculo', icon: 'pi pi-images' }
    );
  }

  if (isAdmin.value) {
    items.push(
      { to: '/admin/compras', label: 'Compras', icon: 'pi pi-dollar' },
      { to: '/admin/proveedores', label: 'Proveedores', icon: 'pi pi-briefcase' }
    );
  }

  return items;
});
const ventasItems = computed(() => {
  const items: Array<{ to: string; label: string; icon: string }> = [];

  if (!isContador.value) {
    items.push(
      { to: '/admin/historial-ventas', label: 'Historial General', icon: 'pi pi-list' },
      { to: '/admin/historial-ventas-personal', label: 'Historial Personal', icon: 'pi pi-user' },
      { to: '/admin/proformas', label: 'Proformas', icon: 'pi pi-file-pdf' }
    );
  } else {
    items.push({ to: '/admin/historial-ventas', label: 'Historial General', icon: 'pi pi-list' });
  }

  return items;
});
const reservasItems = computed(() => {
  const items = [{ to: '/admin/historial-reservas', label: 'Historial General', icon: 'pi pi-list' }];

  items.push({ to: '/admin/pedidos', label: 'Pedidos', icon: 'pi pi-send' });

  if (!isContador.value) {
    items.push({ to: '/admin/historial-reservas-personal', label: 'Historial Personal', icon: 'pi pi-user' });
  }

  return items;
});
const creditosItems = computed(() => {
  const items = [{ to: '/admin/creditos', label: 'Historial General', icon: 'pi pi-list' }];

  if (!isContador.value) {
    items.push({ to: '/admin/creditos-personal', label: 'Historial Personal', icon: 'pi pi-user' });
  }

  return items;
});
const userInitials = computed(() => {
  const nombre = userData.value.nombre?.trim()?.[0] || 'U';
  const apellido = userData.value.apellido?.trim()?.[0] || '';
  return `${nombre}${apellido}`.toUpperCase();
});

function isSectionActive(items: Array<{ to: string }>) {
  return items.some(item => route.path === item.to);
}

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
