# Despliegue mínimo seguro en VPS

## Requisitos

- Un dominio con registros DNS apuntando al VPS.
- Docker Engine y Docker Compose.
- PostgreSQL accesible mediante TLS y con un usuario exclusivo para esta aplicación.
- Firewall con sólo `22/tcp`, `80/tcp` y `443/tcp` expuestos. La API y PostgreSQL no deben publicarse a Internet.

## 1. Preparar secretos

Copie `deploy/.env.example` a `deploy/.env`. Genere `AUTH_TOKEN_SECRET` con al menos 32 bytes aleatorios; en PowerShell:

```powershell
$bytes = New-Object byte[] 48
[Security.Cryptography.RandomNumberGenerator]::Create().GetBytes($bytes)
[Convert]::ToBase64String($bytes)
```

No reutilice contraseñas personales, no suba `deploy/.env` a Git y no envíe su contenido por chat o correo. `INITIAL_ADMIN_PASSWORD` debe tener al menos 12 caracteres; se recomienda una frase aleatoria más larga.

## 2. Primer despliegue o migración controlada

Configure temporalmente en `deploy/.env`:

```dotenv
RUN_MIGRATIONS=true
INITIAL_ADMIN_USER=administrador
INITIAL_ADMIN_PASSWORD=una-frase-aleatoria-muy-larga
```

Luego ejecute desde la raíz del repositorio:

```powershell
docker compose --env-file deploy/.env -f deploy/compose.yml up -d --build
docker compose --env-file deploy/.env -f deploy/compose.yml logs backend
```

Las contraseñas antiguas en texto plano se convierten a bcrypt y sus propietarios quedan obligados a cambiarlas al iniciar sesión.

Cuando la migración termine correctamente, cambie `RUN_MIGRATIONS=false`, elimine `INITIAL_ADMIN_PASSWORD` del archivo y recree el backend:

```powershell
docker compose --env-file deploy/.env -f deploy/compose.yml up -d --force-recreate backend
```

## 3. Operación

- Use `DB_SSLMODE=verify-full` y un certificado válido para PostgreSQL.
- Mantenga copias cifradas de la base de datos y del volumen `uploads`; pruebe su restauración.
- Aplique actualizaciones del VPS y reconstruya las imágenes con regularidad.
- Use llaves SSH, deshabilite el acceso SSH por contraseña y el inicio de sesión directo de `root`.
- Revise intentos de acceso, errores `401`, `403` y `429`, y configure alertas de espacio en disco.
- Rotar `AUTH_TOKEN_SECRET` cierra todas las sesiones existentes. Hágalo si el secreto pudo filtrarse.

## 4. Verificación posterior

```powershell
docker compose --env-file deploy/.env -f deploy/compose.yml ps
curl.exe -I https://su-dominio.example/login
```

Compruebe que HTTP redirige a HTTPS, que la respuesta HTTPS contiene HSTS/CSP y que los puertos `3000`, `5000` y `5432` no son accesibles desde Internet.
