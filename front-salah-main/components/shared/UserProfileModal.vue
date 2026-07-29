<template>
  <Dialog
    :visible="visible"
    modal
    :showHeader="false"
    :draggable="false"
    :style="{ width: '34rem', maxWidth: '94vw' }"
    :pt="{ content: { class: 'profile-dialog-content' } }"
    @update:visible="closeModal"
  >
    <div class="profile-modal">
      <header class="profile-header">
        <button type="button" class="close-button" aria-label="Cerrar" @click="closeModal">
          <i class="pi pi-times"></i>
        </button>

        <div class="avatar-frame">
          <img v-if="hasPhoto" :src="userData.foto" alt="Foto de usuario" />
          <span v-else>{{ initials }}</span>
        </div>

        <div class="profile-copy">
          <h2>{{ fullName }}</h2>
          <p>Personal de Salah Motors</p>
          <span class="accent-line"></span>
        </div>

        <div class="role-pill">
          <i class="pi pi-shield"></i>
          <span>{{ userData.rol || 'Sin rol' }}</span>
        </div>
      </header>

      <section class="profile-body">
        <div class="info-card">
          <span class="info-icon"><i class="pi pi-user"></i></span>
          <div>
            <small>Usuario</small>
            <strong>{{ userData.usuario || 'N/A' }}</strong>
          </div>
        </div>

        <div class="info-card">
          <span class="info-icon"><i class="pi pi-id-card"></i></span>
          <div>
            <small>CI</small>
            <strong>{{ userData.ci || 'N/A' }}</strong>
          </div>
        </div>

        <div class="info-card">
          <span class="info-icon"><i class="pi pi-phone"></i></span>
          <div>
            <small>Celular</small>
            <strong>{{ userData.celular || 'N/A' }}</strong>
          </div>
        </div>

        <div class="info-card">
          <span class="info-icon"><i class="pi pi-map-marker"></i></span>
          <div>
            <small>Direccion</small>
            <strong>{{ userData.direccion || 'N/A' }}</strong>
          </div>
        </div>

        <div class="info-card info-wide">
          <span class="info-icon"><i class="pi pi-check-circle"></i></span>
          <div>
            <small>Estado</small>
            <strong>{{ userData.estado || 'N/A' }}</strong>
          </div>
        </div>
      </section>
    </div>
  </Dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import Dialog from 'primevue/dialog';

interface UserData {
  nombre: string;
  apellido: string;
  ci: string;
  rol: string;
  usuario: string;
  celular: string;
  direccion?: string;
  estado?: string;
  foto?: string;
}

const props = defineProps<{
  visible: boolean;
  userData: UserData;
}>();

const emit = defineEmits(['update:visible']);

const fullName = computed(() => {
  const name = `${props.userData.nombre || ''} ${props.userData.apellido || ''}`.trim();
  return name || 'Usuario';
});

const initials = computed(() => {
  const nombre = props.userData.nombre?.trim()?.[0] || 'U';
  const apellido = props.userData.apellido?.trim()?.[0] || '';
  return `${nombre}${apellido}`.toUpperCase();
});

const hasPhoto = computed(() => Boolean(props.userData.foto && props.userData.foto !== 'N/A'));

const closeModal = () => {
  emit('update:visible', false);
};
</script>

<style scoped>
:deep(.profile-dialog-content) {
  padding: 0;
  overflow: hidden;
  border-radius: 22px;
  background: #ffffff;
}

:deep(.p-dialog) {
  border-radius: 22px;
  box-shadow: 0 28px 90px rgba(13, 13, 13, 0.36);
}

.profile-modal {
  background: #f7f7f7;
  color: #0d0d0d;
}

.profile-header {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 30px 24px 24px;
  background: linear-gradient(135deg, #0d0d0d 0%, #202020 100%);
  color: #ffffff;
  text-align: center;
}

.close-button {
  position: absolute;
  top: 16px;
  right: 16px;
  display: grid;
  place-items: center;
  width: 40px;
  height: 40px;
  border: 1px solid rgba(255, 255, 255, 0.14);
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.08);
  color: #ffffff;
  cursor: pointer;
  transition: background-color 0.2s ease, transform 0.2s ease;
}

.close-button:hover {
  background: rgba(255, 255, 255, 0.16);
  transform: translateY(-1px);
}

.avatar-frame {
  display: grid;
  place-items: center;
  width: 116px;
  height: 116px;
  overflow: hidden;
  border: 4px solid #ffd700;
  border-radius: 30px;
  background: #ffd700;
  color: #0d0d0d;
  font-size: 2.1rem;
  font-weight: 900;
  box-shadow: 0 20px 46px rgba(255, 215, 0, 0.24);
}

.avatar-frame img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.profile-copy h2 {
  margin: 16px 0 4px;
  font-size: 1.55rem;
  font-weight: 900;
  letter-spacing: 0;
}

.profile-copy p {
  margin: 0 0 12px;
  color: rgba(255, 255, 255, 0.72);
  font-size: 0.92rem;
}

.accent-line {
  display: block;
  width: 96px;
  height: 3px;
  margin: 0 auto;
  border-radius: 999px;
  background: #ffd700;
}

.role-pill {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin-top: 16px;
  border: 1px solid rgba(255, 215, 0, 0.35);
  border-radius: 999px;
  background: rgba(255, 215, 0, 0.12);
  padding: 8px 14px;
  color: #ffd700;
  font-size: 0.84rem;
  font-weight: 800;
  text-transform: uppercase;
}

.profile-body {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  padding: 18px;
}

.info-card {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
  border: 1px solid rgba(13, 13, 13, 0.08);
  border-radius: 16px;
  background: #ffffff;
  padding: 14px;
  box-shadow: 0 14px 34px rgba(13, 13, 13, 0.07);
}

.info-wide {
  grid-column: 1 / -1;
}

.info-icon {
  display: grid;
  place-items: center;
  flex: 0 0 auto;
  width: 40px;
  height: 40px;
  border-radius: 13px;
  background: rgba(255, 215, 0, 0.18);
  color: #0d0d0d;
}

.info-card small {
  display: block;
  color: #6a6a6a;
  font-size: 0.72rem;
  font-weight: 800;
  text-transform: uppercase;
}

.info-card strong {
  display: block;
  overflow-wrap: anywhere;
  color: #0d0d0d;
  font-size: 0.95rem;
}

@media (max-width: 520px) {
  .profile-body {
    grid-template-columns: 1fr;
    padding: 14px;
  }

  .info-wide {
    grid-column: auto;
  }
}
</style>
