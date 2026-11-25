# API_MID_AGROCAMPO

## Arranque manual de todos los módulos

Angular envía `/api` a la MID mediante el proxy de desarrollo. La MID resuelve las rutas de cuenta, veterinarias, contenido, noticias, estadísticas, soporte y reportes. Cada CRUD trabaja con su esquema de PostgreSQL en la base `AGROCAMPO`.

1. Inicia PostgreSQL. Prepara la base con [las instrucciones de `database_v2`](../database_v2/README.md). Si necesitas cuentas, categorías y datos de ejemplo, carga `02_datos_prueba_completa.sql` en la base de pruebas. Los usuarios creados solo en el almacenamiento del navegador deben registrarse en el nuevo formulario.
2. Agrega `AGROCAMPO_CLAVE_INTERNA` con una clave aleatoria de al menos 32 caracteres a `AgroCampo_Apis_Crud_Beego/.env`. Usuarios y la MID leen esa misma clave al arrancar. Si se despliegan sin el archivo compartido, define la misma variable de entorno en ambos servicios. Angular no necesita esa clave.
3. Desde cada carpeta de la tabla, ejecuta `go run .` en una terminal independiente. Los CRUD leen su configuración de PostgreSQL de `AgroCampo_Apis_Crud_Beego/.env`.

| Carpeta dentro de `AgroCampo_Apis_Crud_Beego` | Puerto | Pantallas que dependen del servicio |
|---|---:|---|
| `api_Crud_agropecuario` | 8081 | Precios de ganado, mediciones del sector y datos de exportación |
| `api_Crud_Contenido` | 8082 | Foro, temas guardados y videos |
| `api_Crud_Noticias` | 8083 | Noticias, detalle y categorías |
| `api_Crud_Reportes` | 8084 | Historial y descarga de reportes |
| `api_crud_Soporte` | 8085 | Preguntas frecuentes y registro de incidencias |
| `api_Crud_usuarios` | 8086 | Registro, sesión, perfil y recuperación |
| `api_Crud_Veterinarias` | 8087 | Directorio, registro, edición, eliminación y productos |

4. Ejecuta `go run .` desde `API_MID_AGROCAMPO` y `npm start` desde `frontend_AgroCampo`. Los puertos predeterminados son `8080` y `4200`. Abre `http://localhost:4200/#/login`.

Para recuperar contraseñas, configura en la MID las credenciales y el token OAuth de Gmail. `AGROCAMPO_FRONTEND_URL` puede cambiar la dirección que aparece en los correos. Los enlaces duran 15 minutos. La MID guarda sesiones en memoria y solicita un nuevo inicio de sesión al reiniciarse.

Para tus pruebas, registra una cuenta y comprueba que puedes volver a entrar y modificar su perfil. Abre noticias, foro, contenido, soporte y estadísticas con datos de ejemplo. Con el rol correspondiente, publica contenido o actualiza una cotización. Desde Estadísticas, solicita un PDF o Excel y descárgalo en Mis reportes. La descarga conserva las cotizaciones de las últimas cuatro subastas al momento de la solicitud.

El directorio de veterinarias mantiene su API existente. La autorización del servidor para sus escrituras y la persistencia del campo de especialidad son tareas pendientes.

## Conexión de veterinarias

El cliente principal está en `../frontend_AgroCampo`. Angular usa el proxy de desarrollo para enviar `/api/**` a la MID; la MID consulta el CRUD por HTTP y este persiste en PostgreSQL.

```text
Angular :4200 → MID :8080 → CRUD Veterinarias :8087 → PostgreSQL :5432 / AGROCAMPO
```

Inicia PostgreSQL y, en tres terminales distintas desde la raíz de AgroCampo:

```powershell
cd AgroCampo_Apis_Crud_Beego/api_Crud_Veterinarias
go run .
```

```powershell
cd API_MID_AGROCAMPO
go run .
```

```powershell
cd frontend_AgroCampo
npm start
```

Abre `http://localhost:4200/#/veterinarias`. El CRUD carga las credenciales desde `AgroCampo_Apis_Crud_Beego/.env`; la MID no accede a PostgreSQL para este módulo. Se requiere la base `AGROCAMPO` con el esquema `veterinarias` de `../database_v2`.

| Método | Ruta MID | Operación |
|---|---|---|
| GET | `/api/v1/veterinarias` | Listado de clínicas activas con servicios |
| POST | `/api/v1/veterinarias` | Registro de clínica, productos y servicios |
| GET | `/api/v1/veterinarias/{id}` | Detalle con productos y servicios |
| PUT | `/api/v1/veterinarias/{id}` | Editar datos de contacto y presentación |
| DELETE | `/api/v1/veterinarias/{id}` | Eliminar clínica y relaciones según las cascadas del esquema SQL |
| POST | `/api/v1/veterinarias/{id}/productos` | Crear producto |
| PUT | `/api/v1/veterinarias/{id}/productos/{productoId}` | Editar producto |
| DELETE | `/api/v1/veterinarias/{id}/productos/{productoId}` | Eliminar producto |

Los IDs se obtienen del CRUD. Las escrituras actualizan la pantalla después de recibir confirmación; si fallan, se muestra el error y se conserva el formulario. La edición de clínica preserva propietario, aprobación, ubicación y metadatos. La edición/eliminación de productos verifica su asociación con la clínica de la ruta.

## Configuración opcional

### Imágenes de productos

El registro y la edición de productos permiten seleccionar PNG, JPG o GIF de hasta 2 MB. `POST /api/v1/imagenes/productos` recibe un formulario multipart con el campo `imagen`, valida el contenido y devuelve una URL que se guarda en `imagen_url` a través del CRUD. `GET /api/v1/imagenes/productos/{archivo}` sirve la imagen.

Los archivos se guardan en `static/uploads/productos` dentro de la MID, con nombres aleatorios. Conserva esta carpeta al reiniciar o desplegar y respáldala junto con la base. Una imagen subida que se quite del formulario o cuyo producto se elimine permanece en disco; no se realiza limpieza automática.

La opción «Eliminar veterinaria» está en su detalle, para los roles que pueden gestionar veterinarias, y solicita confirmación. El borrado es definitivo y usa las relaciones `ON DELETE CASCADE` del esquema `database_v2/01_esquema.sql`.

La MID lee estas variables del entorno del proceso. Cada URL debe incluir la ruta completa del recurso, sin barra final. Del `.env` compartido solo lee `AGROCAMPO_CLAVE_INTERNA`.

| Variable | Valor predeterminado |
|---|---|
| `API_MID_AGROCAMPO_HTTP_PORT` | `8080` |
| `URL_VETERINARIA` | `http://localhost:8087/v1/veterinaria` |
| `URL_PRODUCTOS_VETERINARIA` | `http://localhost:8087/v1/producto_veterinaria` |
| `URL_SERVICIOS_VETERINARIA` | `http://localhost:8087/v1/servicioveterinaria` |
| `URL_RELACIONES_SERVICIOS_VETERINARIA` | `http://localhost:8087/v1/tr_veterinariaservicio` |

Si cambias el puerto de la MID, ajusta también `frontend_AgroCampo/proxy.conf.json` y reinicia Angular. El proxy de `ng serve` solo existe en desarrollo; un despliegue necesita enrutar `/api` hacia la MID desde su servidor.

## Verificación

Desde esta carpeta:

```powershell
go test ./... -timeout 60s
./scripts/Test-Veterinarias.ps1
./scripts/Test-Veterinarias.ps1 -Escritura
```

El script consulta a través de Angular en `:4200`. Con `-Escritura` crea una clínica temporal con producto y servicio, prueba ediciones y eliminación de productos, vuelve a consultar para comprobar persistencia y limpia sus propios registros. Para probar solo MID → CRUD usa `-BaseUrl http://localhost:8080`. Requiere los servidores iniciados.

Verificado el 19 de noviembre de 2026 contra PostgreSQL local y a través del proxy de Angular, con limpieza de los registros temporales.

## Alcance pendiente

El cliente obtiene la sesión y el rol desde la MID. En el módulo de veterinarias sigue pendiente verificar los permisos del propietario en la API y guardar el campo de especialidad. El mapa requiere una clave de Google Maps para mostrar coordenadas.

El registro de clínica y sus relaciones se realiza con varias solicitudes HTTP. Si falla un producto o servicio después de crear la clínica, la MID informa del guardado parcial; todavía no existe una transacción distribuida ni compensación automática.
