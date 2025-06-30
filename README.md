# Sistema de Gestión de Libros

**Datos del Grupo:**
* Kevin Sebastian Perez Cervantes

**Fecha:** 29 de junio de 2025

---

## 📚 Objetivo del Programa

Este proyecto es un sistema básico de gestión de una biblioteca o librería, desarrollado en Go. Su objetivo principal es permitir el registro, consulta, edición y eliminación de información sobre libros, así como visualizar un resumen del inventario disponible y prestado. Busca demostrar la integración de conocimientos adquiridos en la asignatura, incluyendo el desarrollo web con Go, interacción con bases de datos y la creación de servicios web RESTful.

## ✨ Funcionalidades Principales

El aplicativo ofrece las siguientes funcionalidades clave:

1.  **Dashboard de Resumen:** Visualización de métricas importantes sobre el inventario de libros (Total de Libros, Libros Disponibles).
2.  **Gestión de Libros (CRUD):**
    * **Listar Libros:** Muestra una tabla con todos los libros registrados en el sistema.
    * **Crear Nuevo Libro:** Permite añadir nuevos registros de libros a la base de datos.
    * **Editar Libro:** Posibilita modificar la información de un libro existente, incluyendo su estado de "prestado".
    * **Eliminar Libro:** Permite remover libros de la base de datos.
3.  **Servicios Web RESTful (API JSON):** Exposición de las funcionalidades principales del sistema a través de una API para que otras aplicaciones puedan interactuar programáticamente con los datos de los libros.

## 🚀 Cómo Ejecutar el Proyecto

1.  **Clonar el repositorio:**
    ```bash
    git clone [https://github.com/tu-usuario/proyecto.git](https://github.com/tu-usuario/proyecto.git)
    cd proyecto
    ```
2.  **Configurar la Base de Datos:**
    * El proyecto utiliza MySQLite. Asegúrate de que el archivo de la base de datos (ej. `database.db`) esté en la ruta correcta o se cree si no existe.
    * La estructura de la tabla `libros` es la siguiente:
        ```sql
        CREATE TABLE libros (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            titulo TEXT NOT NULL,
            autor TEXT NOT NULL,
            anio_publicacion INTEGER,
            editorial TEXT,
            prestado BOOLEAN NOT NULL DEFAULT FALSE
        );
        ```
3.  **Instalar dependencias de Go:**
    ```bash
    go mod tidy
    ```
4.  **Ejecutar la aplicación:**
    ```bash
    go run inicio.go
    ```
5.  **Acceder a la aplicación:**
    Abre tu navegador web y visita `http://localhost:8080/` (o el puerto configurado en `inicio.go`).

## 💻 Estructura del Proyecto

* `/`: Archivo principal `inicio.go` y `go.mod`, `go.sum`.
* `/database`: Contiene la lógica para la conexión a la base de datos.
* `/handlers`: Contiene las funciones que manejan las solicitudes HTTP (tanto para vistas HTML como para la API JSON).
* `/models`: Define las estructuras de datos (ej. `Libro`).
* `/static`: Archivos estáticos como CSS (`style.css`).
* `/templates`: Archivos HTML para las vistas de la aplicación.

## 🔧 Tecnologías Utilizadas

* **Go (Golang):** Lenguaje de programación principal.
* **Gorilla Mux:** Enrutador HTTP para Go.
* **HTML/CSS:** Para la interfaz de usuario.
* **SQLite3:** Base de datos.
* **JSON:** Para la serialización de datos en los servicios web.

## 🔗 Enlaces de Interés

* Repositorio de GitHub: `https://github.com/KSPC98/proyecto.git` 
