# Microbackends DSIII

Este proyecto contiene un backend en Go que expone un CRUD de palabras conectado a Supabase (Postgres) y desplegado en Vercel.

## Tecnologias
- Go
- Vercel Serverless Functions
- Supabase (Postgres)
- pgx

## Base URL
Produccion:

```text
https://microbackends2-dsiii.vercel.app
```

## Endpoints
- `GET /api/words` -> Lista todas las palabras.
- `POST /api/words` -> Inserta una palabra.
- `PUT /api/words/:id` -> Actualiza la palabra con ese `id`.
- `DELETE /api/words/:id` -> Elimina la palabra con ese `id`.

## Body JSON para POST y PUT
Estos endpoints esperan un cuerpo JSON con esta estructura:

```json
{
  "word": "valor"
}
```

## Ejemplos
Obtener todas las palabras:

```http
GET /api/words
```

Crear una palabra:

```http
POST /api/words
Content-Type: application/json
```

```json
{
  "word": "hola"
}
```

Actualizar una palabra:

```http
PUT /api/words/1
Content-Type: application/json
```

```json
{
  "word": "adios"
}
```

Eliminar una palabra:

```http
DELETE /api/words/1
```

## Respuestas esperadas
- `GET` responde con un arreglo de objetos:

```json
[
  {
    "id": 1,
    "word": "hola"
  }
]
```

- `POST` responde:

```json
{
  "message": "Palabra creada"
}
```

- `PUT` responde:

```json
{
  "message": "Palabra actualizada"
}
```

- `DELETE` responde:

```json
{
  "message": "Palabra eliminada"
}
```

## Variables de entorno
La aplicacion usa `SUPABASE_URL` como cadena de conexion principal a la base de datos.
