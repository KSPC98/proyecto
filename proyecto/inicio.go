/*
@Autor: Kevin Pérez
@Descripcion: Aplicación web para gestionar una biblioteca de libros con una interfaz web y una API
*/

package main

import (
	"log"               // Paquete para logging.
	"net/http"          // Paquete para manejar solicitudes y respuestas HTTP.
	"proyecto/db"       // Importa el paquete db para la conexión a la base de datos.
	"proyecto/handlers" // Importa el paquete handlers que contiene los manejadores de rutas.

	"github.com/gorilla/mux" // Router HTTP para Go.
)

func main() {
	// Establece la conexión a la base de datos al inicio de la aplicación.
	// Si la conexión falla, el programa terminará (panic).
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err) // Usa Fatalf para terminar el programa con un mensaje.
	}
	// `defer database.Close()` asegura que la conexión a la base de datos se cierre cuando la función main termine.
	defer database.Close()
	log.Println("Conexión a la base de datos establecida correctamente.")

	// Crea un nuevo enrutador de Gorilla Mux.
	// Mux es un enrutador HTTP que nos permite definir rutas de URL de manera flexible.
	r := mux.NewRouter()

	// Sirve archivos estáticos desde el directorio "static"
	// Esto permite que el navegador cargue CSS, JavaScript, imágenes, etc.
	// Por ejemplo, una solicitud a /static/style.css buscará el archivo en el directorio "static/style.css".
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Rutas para la interfaz web (HTML).
	// Cada HandleFunc asocia una URL con una función manejadora y un método HTTP (GET, POST, etc.).
	r.HandleFunc("/", handlers.HomeHandler(database)).Methods("GET")                     // Ruta para la página de inicio.
	r.HandleFunc("/libros", handlers.RecuperarLibros).Methods("GET")                     // Ruta para listar todos los libros.
	r.HandleFunc("/libros/crear", handlers.CreateLibroGetHandler).Methods("GET")         // Muestra el formulario para crear un libro.
	r.HandleFunc("/libros/crear", handlers.CreateLibroPostHandler).Methods("POST")       // Procesa el envío del formulario para crear un libro.
	r.HandleFunc("/libros/editar/{Id}", handlers.UpdateLibroGetHandler).Methods("GET")   // Muestra el formulario para editar un libro por su ID.
	r.HandleFunc("/libros/editar/{Id}", handlers.UpdateLibroPostHandler).Methods("POST") // Procesa el envío del formulario para actualizar un libro.
	r.HandleFunc("/libros/eliminar/{Id}", handlers.DeleteLibroHandler).Methods("GET")    // Elimina un libro por su ID.

	// Rutas para la API (para aplicaciones cliente-servidor, ej. JavaScript frontend)
	// Se crea un sub-enrutador para las rutas de la API, todas comenzarán con /api.
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/libros", handlers.ApiListarLibros).Methods("GET")          // API para listar todos los libros.
	apiRouter.HandleFunc("/libros/{Id}", handlers.ApiObtenerLibro).Methods("GET")     // API para obtener un libro por ID.
	apiRouter.HandleFunc("/libros", handlers.ApiCrearLibro).Methods("POST")           // API para crear un nuevo libro.
	apiRouter.HandleFunc("/libros/{Id}", handlers.ApiActualizarLibro).Methods("PUT")  // API para actualizar un libro existente.
	apiRouter.HandleFunc("/libros/{Id}", handlers.ApiEliminarLibro).Methods("DELETE") // API para eliminar un libro.

	// Mensaje de log que indica que el servidor se está iniciando.
	log.Println("Servidor iniciado en http://localhost:8000")
	// Inicia el servidor HTTP en el puerto 8080.
	// `log.Fatal` se usa aquí para que si el servidor no puede arrancar (ej. puerto ya en uso),
	// el error se registre y el programa termine.
	log.Fatal(http.ListenAndServe(":8000", r))
}
