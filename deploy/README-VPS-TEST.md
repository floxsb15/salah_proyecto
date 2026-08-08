# Prueba en VPS sin dominio y sin Caddy

Este flujo publica el frontend unicamente en `127.0.0.1:3000` del VPS. PostgreSQL y el backend permanecen en la red privada de Docker. Para abrir la aplicacion se utiliza un tunel SSH.

## Antes de comenzar

- Docker Engine y Docker Compose deben estar instalados.
- El repositorio debe estar clonado en el VPS. Puede estar en `/opt/salah` o en una carpeta dentro del usuario; los scripts no dependen de una ruta concreta.
- El usuario del VPS debe poder ejecutar Docker.
- No copie secretos al repositorio ni los envie por chat.

## 1. Preparar el entorno

Desde la raiz del repositorio en el VPS:

```bash
chmod +x deploy/*.sh
./deploy/prepare-vps-test.sh
```

El script crea `deploy/.env` con permisos `600`, genera secretos aleatorios y configura el primer usuario como `administrador`. No reemplaza un archivo existente.

## 2. Levantar la prueba

```bash
./deploy/start-vps-test.sh
```

Se levantan solamente `postgres`, `backend` y `frontend`. El servicio `caddy` no se inicia.

## 3. Acceder desde Windows

Mantenga abierto este comando en PowerShell:

```powershell
ssh -L 3000:127.0.0.1:3000 german@srv989157.hstgr.cloud
```

Luego abra `http://localhost:3000/login`. El puerto 3000 no necesita abrirse en el firewall del VPS.

Para consultar el usuario y la clave inicial, hagalo dentro de su terminal SSH:

```bash
grep -E '^INITIAL_ADMIN_(USER|PASSWORD)=' deploy/.env
```

No publique esa salida ni la envie por chat.

## 4. Cerrar la inicializacion

Despues de comprobar el acceso y guardar la clave inicial:

```bash
./deploy/finish-vps-test-setup.sh --password-saved
```

Esto cambia `RUN_MIGRATIONS` a `false`, retira la clave inicial del entorno y recrea backend y frontend.

## Diagnostico

Estado de los servicios:

```bash
docker compose --env-file deploy/.env -f deploy/compose.yml -f deploy/compose.local-test.yml ps
```

Logs sin exponer las variables de entorno:

```bash
docker compose --env-file deploy/.env -f deploy/compose.yml -f deploy/compose.local-test.yml logs --tail=100 backend frontend postgres
```
