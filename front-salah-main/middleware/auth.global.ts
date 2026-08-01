export default defineNuxtRouteMiddleware((to) => {
  if (import.meta.server) {
    return;
  }

  const user = readStoredUser();
  const isPublic = to.path === '/login';
  if (!user && !isPublic) {
    return navigateTo('/login');
  }
  if (!user) {
    return;
  }
  if (user.must_change_password && to.path !== '/cambiar-contrasena') {
    return navigateTo('/cambiar-contrasena');
  }
  if (to.path === '/cambiar-contrasena' && !user.must_change_password) {
    return navigateTo(homeForRole(user.rol));
  }
  if (isPublic) {
    return navigateTo(homeForRole(user.rol));
  }
  if (to.path.startsWith('/admin') && user.rol !== 'admin') {
    return navigateTo(homeForRole(user.rol));
  }
  if (to.path.startsWith('/ventas') && !['admin', 'encargado de ventas', 'vendedor'].includes(user.rol)) {
    localStorage.removeItem('user');
    return navigateTo('/login');
  }
});

function readStoredUser(): { rol: string; must_change_password?: boolean } | null {
  try {
    const value = localStorage.getItem('user');
    return value ? JSON.parse(value) : null;
  } catch {
    localStorage.removeItem('user');
    return null;
  }
}

function homeForRole(role: string): string {
  return role === 'admin' ? '/admin/dashboard' : '/ventas/catalogo-vehiculos';
}
