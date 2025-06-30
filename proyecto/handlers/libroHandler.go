/*
@Autor: Kevin Pérez
@Descipcion: Módulo que maneja las operaciones CRUD para la entidad Libro en la interfaz web.
*/

package handlers

import (
	"fmt"             // Paquete para formatear cadenas.
	"html/template"   // Paquete para trabajar con plantillas HTML.
	"log"             // Paquete para logging.
	"net/http"        // Paquete para manejar solicitudes HTTP.
	"proyecto/models" // Importa el paquete models para interactuar con los datos de libros.
	"strconv"         // Paquete para conversión de tipos.
	"time"            // Paquete para obtener la fecha y hora actual.

	"github.com/gorilla/mux" // Router HTTP para manejar rutas.
)

// RecuperarLibros maneja la solicitud para listar todos los libros en la interfaz web.
func RecuperarLibros(w http.ResponseWriter, r *http.Request) {
	// Obtiene todos los libros de la base de datos.
	libros, err := models.GetAllLibros()
	if err != nil {
		// Si hay un error, se envía una respuesta de error 500.
		http.Error(w, "Error al recuperar los libros: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parsea los archivos de plantilla base.html y libros.html.
	tmpl, err := template.ParseFiles("templates/base.html", "templates/libros.html")
	if err != nil {
		// Si hay un error al cargar las plantillas, se registra el error y se envía una respuesta de error 500.
		log.Printf("Error al cargar el template: %v", err)
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}

	// Imprime los libros en la consola del servidor (útil para depuración).
	fmt.Println(libros)

	// Ejecuta la plantilla "base" pasando los libros como datos.
	err = tmpl.ExecuteTemplate(w, "base", libros)
	if err != nil {
		// Si hay un error al ejecutar la plantilla, se envía una respuesta de error 500.
		http.Error(w, "Error al ejecutar el template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateLibroGetHandler muestra el formulario HTML para crear un nuevo libro.
func CreateLibroGetHandler(w http.ResponseWriter, r *http.Request) {
	// Parsea los archivos de plantilla base.html y crearLibro.html.
	tmpl, err := template.ParseFiles("templates/base.html", "templates/crearLibro.html")
	if err != nil {
		// Si hay un error al cargar las plantillas, se registra el error y se envía una respuesta de error 500.
		log.Printf("Error al cargar el template: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// Pasa el año actual a la plantilla para el valor máximo del campo AnioPublicacion.
	data := struct {
		CurrentYear int
	}{
		CurrentYear: time.Now().Year(),
	}

	// Ejecuta la plantilla "base" sin pasar datos inicialmente.
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		// Si hay un error al ejecutar la plantilla, se registra el error y se envía una respuesta de error 500.
		log.Printf("Error al ejecutar el template: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}
}

// CreateLibroPostHandler procesa los datos del formulario para crear un nuevo libro.
func CreateLibroPostHandler(w http.ResponseWriter, r *http.Request) {
	// Verifica que la solicitud sea de tipo POST.
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parsea el formulario para acceder a los valores enviados.
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear el formulario: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Recupera los valores de los campos del formulario.
	Autor := r.FormValue("Autor")
	Titulo := r.FormValue("Titulo")
	AnioPublicacionStr := r.FormValue("AnioPublicacion")
	Editorial := r.FormValue("Editorial")
	Prestado := r.FormValue("Prestado")

	// Validaciones básicas de los campos del formulario.
	if Titulo == "" || Autor == "" || AnioPublicacionStr == "" || Editorial == "" || Prestado == "" {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	// Convierte el año de publicación de string a int.
	AnioPublicacion, err := strconv.Atoi(AnioPublicacionStr)
	if err != nil {
		http.Error(w, "El año de publicación debe ser un número válido", http.StatusBadRequest)
		return
	}

	// Llama a la función CreateLibro del modelo para guardar el libro en la base de datos.
	err = models.CreateLibro(Autor, Titulo, AnioPublicacion, Editorial, Prestado)
	if err != nil {
		http.Error(w, "Error al crear el libro: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirige al usuario a la lista de libros después de una creación exitosa.
	http.Redirect(w, r, "/libros", http.StatusSeeOther)

}

// UpdateLibroGetHandler muestra el formulario para editar un libro existente.
func UpdateLibroGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["Id"])
	if err != nil {
		http.Error(w, "ID de libro inválido", http.StatusBadRequest)
		return
	}

	libro, err := models.GetLibroByID(id)

	if err != nil {
		http.Error(w, "Libro no encontrado: "+err.Error(), http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/editarLibro.html")

	if err != nil {
		log.Printf("Error al cargar el template: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// Estructura para pasar el libro y el año actual a la plantilla.

	data := struct {
		models.Libro
		CurrentYear int
	}{
		Libro:       libro,
		CurrentYear: time.Now().Year(),
	}

	err = tmpl.ExecuteTemplate(w, "base", data)

	if err != nil {
		log.Printf("Error al ejecutar el template: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

}

// UpdateLibroPostHandler procesa los datos del formulario para actualizar un libro.

func UpdateLibroPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear el formulario: "+err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["Id"])
	if err != nil {
		http.Error(w, "ID de libro inválido", http.StatusBadRequest)
		return
	}

	Autor := r.FormValue("Autor")
	Titulo := r.FormValue("Titulo")
	AnioPublicacionStr := r.FormValue("AnioPublicacion")
	Editorial := r.FormValue("Editorial")
	Prestado := r.FormValue("Prestado")

	if Titulo == "" || Autor == "" || AnioPublicacionStr == "" || Editorial == "" || Prestado == "" {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	AnioPublicacion, err := strconv.Atoi(AnioPublicacionStr)

	if err != nil {
		http.Error(w, "El año de publicación debe ser un número válido", http.StatusBadRequest)
		return
	}

	// Crea una instancia de Libro con los datos actualizados.

	libro := models.Libro{
		Id:              id,
		Autor:           Autor,
		Titulo:          Titulo,
		AnioPublicacion: AnioPublicacion,
		Editorial:       Editorial,
		Prestado:        Prestado,
	}

	err = models.UpdateLibro(libro)

	if err != nil {
		http.Error(w, "Error al actualizar el libro: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/libros", http.StatusSeeOther)
}

// DeleteLibroHandler maneja la solicitud para eliminar un libro.
func DeleteLibroHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["Id"])
	if err != nil {
		http.Error(w, "ID de libro inválido", http.StatusBadRequest)
		return
	}

	err = models.DeleteLibro(id)

	if err != nil {
		http.Error(w, "Error al eliminar el libro: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/libros", http.StatusSeeOther)
}
