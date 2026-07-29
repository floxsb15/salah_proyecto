## Setup

En la raiz del proyecto ejecutar las siguientes lineas
```bash
docker build -t app_delitto .
```
Para montar el docker del proyecto ejecutar:
```bash
docker run -p 3000:3000 -d app_delitto
```

## Development Server

Start the development server on `http://localhost:3000`:
