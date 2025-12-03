# Repositorio de la API de Mangocatnotes

## Descripción

¡Bienvenido al repositorio de la API de [Mangocatnotes](https://github.com/manuelmtzv/mangocatnotes-web)! Este repositorio contiene el código backend para Mangocatnotes, una aplicación web gratuita de gestión de notas diseñada para ayudar a los usuarios a organizar sus notas sin esfuerzo. Con Mangocatnotes, puedes crear y gestionar tus notas desde cualquier dispositivo con un navegador web, asegurando un acceso universal y una experiencia de toma de notas sin problemas.

## Características

- Gestión de notas: Crea y gestiona tus notas.
- Acceso universal: Accede a tus notas desde cualquier dispositivo con un navegador web, garantizando conveniencia y flexibilidad.
- Etiquetas personalizadas: Organiza tus notas eficientemente añadiendo etiquetas personalizadas, permitiendo una mejor categorización y recuperación de información.

## Instalación

```bash
$ npm install
```

## Configuración

```bash
# Configurar .env
PORT="api port"
DATABASE_URL="mongodb url"
JWT_SECRET="secret string"
```

## Ejecutar API

```bash
# development
$ npm run start

# watch mode
$ npm run start:dev

# production mode
$ npm run start:prod
```

## Test

```bash
# unit tests
$ npm run test

# e2e tests
$ npm run test:e2e

# test coverage
$ npm run test:cov
```

## Documentación del API

- Web: [Swagger Docs](https://mangocatnotes-api.up.railway.app/api)
- Local: http://localhost:3000/api (o en el puerto seleccionado).
