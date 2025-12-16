# Mangocatnotes

Aplicación web para la creación y gestión de notas personales. Permite registrarse, iniciar sesión y acceder a las notas desde cualquier dispositivo con navegador web.

Este proyecto representa una filosofía de simplificación: volver a los fundamentos y demostrar que las mejores soluciones no siempre requieren tecnologías elaboradas.

## Stack Tecnológico

- **Backend:** Go (Chi router)
- **Base de datos:** PostgreSQL
- **Cache/Sesiones:** Redis
- **Frontend:** Templates HTML, HTMX, Alpine.js
- **Estilos:** TailwindCSS

## Requisitos

- Go 1.25+
- PostgreSQL 16+
- Redis 7+
- Node.js 22+ (para compilar TailwindCSS)

## Configuración

1. Clonar el repositorio:

```bash
git clone https://github.com/manuelmtzv/mangocatnotes-api.git
cd mangocatnotes-api
```

2. Copiar el archivo de configuración y ajustar las variables:

```bash
cp .env.example .env
```

3. Instalar dependencias de Node.js:

```bash
npm install
```

4. Ejecutar las migraciones de base de datos:

```bash
make migrate-up
```

## Desarrollo

Iniciar el servidor en modo desarrollo con hot-reload:

```bash
make dev
```

Compilar TailwindCSS en modo watch:

```bash
npm run dev
```

## Producción

Construir la imagen de Docker:

```bash
docker build -t mangocatnotes .
```

Ejecutar el contenedor:

```bash
docker run -p 8080:8080 --env-file .env mangocatnotes
```

## Estructura del Proyecto

```
.
├── cmd/
│   ├── migrate/        # Migraciones de base de datos
│   ├── seed/           # Datos de prueba
│   └── server/         # Punto de entrada de la aplicación
├── internal/
│   ├── config/         # Configuración
│   ├── db/             # Conexión a base de datos
│   ├── handlers/       # Manejadores HTTP
│   ├── kvstore/        # Abstracción de Redis
│   ├── server/         # Servidor HTTP
│   ├── session/        # Gestión de sesiones
│   └── store/          # Acceso a datos
└── web/
    ├── locales/        # Archivos de internacionalización
    ├── static/         # Archivos estáticos (CSS, JS, imágenes)
    └── templates/      # Templates HTML
```
