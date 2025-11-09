# 🪣 Stori API – Lambda de carga y validación de CSV

**Autor:** Joseph Gutiérrez  
**Lenguaje:** Go (Golang)  
**Entorno:** AWS Lambda + API Gateway + S3

---

## 📖 Descripción general

Este proyecto implementa un servicio pequeño y enfocado que expone un **endpoint HTTP (API Gateway + Lambda)** para:

1. Recibir una solicitud de carga de archivo (CSV).
2. Validar que el archivo:
    - Exista.
    - Tenga extensión o tipo de contenido `.csv`.
    - No esté vacío.
    - Tenga la cabecera esperada: `Id,Date,Transaction`.
3. Subir el archivo validado a un **bucket de S3**.
4. Devolver una respuesta JSON clara indicando éxito o errores de validación.

Está diseñado para que **otra Lambda**, ubicada en otro servicio o repositorio, procese posteriormente el archivo subido a S3.

El código sigue un estilo **hexagonal (puertos y adaptadores)**, separando la lógica de dominio del detalle técnico de AWS o HTTP.

---

## 🧱 Estructura del proyecto

```text
.
├── cmd/
│   └── lambda_api/              # Punto de entrada de la Lambda (main.go)
│
├── internal/
│   ├── core/
│   │   ├── application/         # Casos de uso (UploadService)
│   │   └── ports/               # Interfaces (puertos in/out)
│   │
│   ├── infra/
│   │   ├── aws/                 # Cliente AWS SDK (S3)
│   │   ├── bootstrap/           # Inicialización de dependencias
│   │   ├── config/              # Configuración por variables de entorno
│   │   └── logger/              # Logging estructurado
│   │
│   └── interfaces/
│       ├── in/
│       │   └── apigw/           # Adaptador de entrada (API Gateway Handler)
│       └── out/
│           └── s3uploader/      # Adaptador de salida (S3 uploader)
│
├── Dockerfile                   # Build de la imagen Lambda
├── docker-compose.yml           # Entorno local / LocalStack
├── Makefile                     # Automatización de build / test / docker
├── go.mod                       # Definición del módulo Go
└── README.md
```

---

## ⚙️ Configuración

La Lambda lee su configuración desde **variables de entorno**, generalmente definidas en IaC o AWS Console:

- `S3_BUCKET_NAME` – bucket destino para las cargas CSV.
- `S3_REGION` – región AWS del bucket.
- Otras variables opcionales para observabilidad o logging.

---

## 🧪 Ejecución de pruebas

Las pruebas unitarias se encuentran bajo `internal/...` e incluyen dominio, servicios y adaptadores.

```bash
# Pruebas unitarias
make test

# Todas las pruebas (placeholder para integración futura)
make test-all
```

---

## 🐳 Entorno local

Puedes correr servicios de apoyo (como LocalStack para S3 + API Gateway) usando Docker Compose:

```bash
make compose-up     # Levanta el entorno local
make compose-down   # Detiene los contenedores
make rebuild        # Reconstruye las imágenes y recrea contenedores
make reset          # Limpia todo (contenedores + volúmenes locales)
```

> 💡 *El archivo docker-compose.yml puede configurarse para que la Lambda se ejecute localmente y use un S3 simulado en lugar de AWS real.*

---

## 🚀 Build y publicación (imagen Lambda)

Construir la imagen Docker y etiquetarla para ECR:

```bash
make build
```

Iniciar sesión en ECR usando tu perfil AWS (`personal` por defecto):

```bash
make login
```

Publicar la imagen:

```bash
make publish
```

Después de esto, puedes apuntar tu Lambda (tipo imagen) al ECR correspondiente.

---

## 🧹 Comandos del Makefile

Resumen de los principales targets:

```text
make clean           # Limpia y organiza dependencias Go
make build           # Construye y etiqueta la imagen
make publish         # Publica en ECR
make login           # Login en ECR

make compose-up      # Levanta entorno local
make compose-down    # Detiene entorno
make rebuild         # Rebuild completo
make reset           # Limpieza total

make test            # Ejecuta pruebas unitarias
make test-integration# Placeholder
make test-all        # Todas las pruebas
```

---

## 🌐 Flujo del API

El API Gateway invoca la Lambda mediante un evento HTTP (versión 2.0).  
El flujo típico es:

1. Cliente envía `POST /upload` con un archivo CSV.
2. El adaptador (`interfaces/in/apigw`) valida y entrega la solicitud al servicio de aplicación.
3. El servicio valida el CSV y usa el puerto S3 para subirlo.
4. La Lambda responde algo como:

```json
{
  "success": true,
  "bucket": "stori-uploads-dev",
  "key": "uploads/txns-2025-11-08-123456.csv"
}
```

O, en caso de error:

```json
{
  "success": false,
  "error": "cabecera inválida, se esperaba: Id,Date,Transaction"
}
```

---

## 💬 Notas

- La segunda Lambda (procesadora) puede escuchar los eventos del bucket S3 para continuar el pipeline.
- Mantener este servicio en un repositorio separado mejora la claridad durante una evaluación técnica.
- El código es deliberadamente simple y legible para que un revisor entienda la arquitectura rápidamente.

---

**© 2025 — Joseph Gutiérrez**