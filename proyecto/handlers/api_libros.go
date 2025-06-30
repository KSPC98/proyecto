/*
@Autor: Kevin Pérez
@Descripcion: Módulo que maneja las operaciones CRUD para la entidad Libro en la API.
*/

package handlers

import (
	"net/http"        // Paquete para manejar solicitudes y respuestas HTTP.
	"proyecto/models" // Importa el paquete models donde se define la estructura Libro y funciones CRUD.
	"strconv"         // Paquete para la conversión de cadenas a tipos numéricos.

	"github.com/goccy/go-json" // Paquete para codificar/decodificar JSON de forma eficiente.
	"github.com/gorilla/mux"   // Router HTTP para manejar las rutas de la aplicación.
)

// LibroSimple es una estructura para representar una versión simplificada de un libro para la API.
// Solo incluye los campos que se desean exponer públicamente en ciertas respuestas de la API.
type LibroSimple struct {
	Autor    string `json:"autor"`    // El autor del libro.
	Titulo   string `json:"titulo"`   // El título del libro.
	Prestado string `json:"prestado"` // El estado de préstamo del libro.
}

// ApiListarLibros maneja la solicitud para obtener una lista simplificada de todos los libros.
func ApiListarLibros(w http.ResponseWriter, r *http.Request) {
	// Obtiene todos los libros de la base de datos a través del modelo.
	libros, err := models.GetAllLibros()
	if err != nil {
		// Si ocurre un error al recuperar los libros, se envía una respuesta de error 500.
		http.Error(w, "Error al recuperar los libros: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Se crea una slice (arreglo dinámico) para almacenar la versión simplificada de los libros.
	var datosSimples []LibroSimple
	// Itera sobre cada libro obtenido y crea un objeto LibroSimple con los campos deseados.
	for _, libro := range libros {
		datosSimples = append(datosSimples, LibroSimple{
			Autor:    libro.Autor,
			Titulo:   libro.Titulo,
			Prestado: libro.Prestado,
		})

	}

	// Establece el encabezado Content-Type de la respuesta a "application/json".
	w.Header().Set("Content-Type", "application/json")

	// Codifica la slice de LibroSimple a formato JSON y la escribe en la respuesta.
	if err := json.NewEncoder(w).Encode(datosSimples); err != nil {
		http.Error(w, "Error al codificar la respuesta JSON: "+err.Error(), http.StatusInternalServerError)
	}
}

// ApiObtenerLibro maneja la solicitud para obtener un libro específico por su ID.
func ApiObtenerLibro(w http.ResponseWriter, r *http.Request) {
	// Extrae las variables de la URL (en este caso, el ID del libro).
	vars := mux.Vars(r)
	// Convierte el ID de la URL (que es una cadena) a un entero.
	id, err := strconv.Atoi(vars["Id"])
	if err != nil {
		// Si el ID no es un número válido, se envía una respuesta de error 400.
		http.Error(w, "ID de libro inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Obtiene el libro de la base de datos por su ID.
	libro, err := models.GetLibroByID(id)
	if err != nil {
		// Si el libro no se encuentra o hay un error en la base de datos, se envía una respuesta de error.
		http.Error(w, "Error al recuperar el libro: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Establece el encabezado Content-Type de la respuesta a "application/json".
	w.Header().Set("Content-Type", "application/json")
	// Codifica el objeto Libro a formato JSON y lo escribe en la respuesta.
	if err := json.NewEncoder(w).Encode(libro); err != nil {
		http.Error(w, "Error al codificar la respuesta JSON: "+err.Error(), http.StatusInternalServerError)
	}
}

// ApiCrearLibro maneja la solicitud para crear un nuevo libro.
func ApiCrearLibro(w http.ResponseWriter, r *http.Request) {
	var libro models.Libro // Declara una variable de tipo Libro para decodificar el JSON del cuerpo de la solicitud.
	// Decodifica el cuerpo de la solicitud JSON en la estructura Libro.
	err := json.NewDecoder(r.Body).Decode(&libro)
	if err != nil {
		// Si el JSON es inválido o incompleto, se envía una respuesta de error 400.
		http.Error(w, "Error al decodificar el JSON del libro: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Llama a la función CreateLibro del modelo para insertar el nuevo libro en la base de datos.
	err = models.CreateLibro(libro.Autor, libro.Titulo, libro.AnioPublicacion, libro.Editorial, libro.Prestado)
	if err != nil {
		// Si hay un error al crear el libro en la base de datos, se envía una respuesta de error 500.
		http.Error(w, "Error al crear el libro en la base de datos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Si la creación es exitosa, se establece el código de estado HTTP 201 (Created).
	w.WriteHeader(http.StatusCreated)

	// Se codifica el libro creado (con su posible ID asignado por la DB si la estructura Libro lo incluyera)
	// y se envía como respuesta JSON.
	if err := json.NewEncoder(w).Encode(libro); err != nil {
		http.Error(w, "Error al codificar la respuesta JSON: "+err.Error(), http.StatusInternalServerError)
	}
}

// ApiActualizarLibro maneja la solicitud para actualizar un libro existente.
func ApiActualizarLibro(w http.ResponseWriter, r *http.Request) {
	// Extrae el ID del libro de la URL.
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["Id"])
	if err != nil {
		http.Error(w, "ID de libro inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	var libro models.Libro // Estructura para decodificar el JSON de la solicitud.
	// Decodifica el cuerpo de la solicitud JSON en la estructura Libro.
	err = json.NewDecoder(r.Body).Decode(&libro)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON del libro: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Asigna el ID de la URL al objeto libro, asegurando que se actualice el libro correcto.
	libro.Id = id

	// Llama a la función UpdateLibro del modelo para actualizar el libro en la base de datos.
	err = models.UpdateLibro(libro)
	if err != nil {
		http.Error(w, "Error al actualizar el libro: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Si la actualización es exitosa, se envía un estado HTTP 200 (OK) y el libro actualizado.
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(libro); err != nil {
		http.Error(w, "Error al codificar la respuesta JSON: "+err.Error(), http.StatusInternalServerError)
	}
}

// ApiEliminarLibro maneja la solicitud para eliminar un libro por su ID.
func ApiEliminarLibro(w http.ResponseWriter, r *http.Request) {

	// Extrae el ID del libro de la URL.
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["Id"])
	if err != nil {
		http.Error(w, "ID de libro inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Llama a la función DeleteLibro del modelo para eliminar el libro de la base de datos.
	err = models.DeleteLibro(id)
	if err != nil {
		http.Error(w, "Error al eliminar el libro: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Si la eliminación es exitosa, se envía un estado HTTP 204 (No Content) para indicar que la acción fue exitosa
	// pero no hay contenido que devolver.
	w.WriteHeader(http.StatusNoContent)
}
