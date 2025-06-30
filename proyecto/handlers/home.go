/*
@Author: kevin Perez
@Descripcion: Este es el manejador de la pagina de inicio del proyecto, donde se carga la plantilla base y la plantilla home.
*/

package handlers

import (
	"database/sql"
	"fmt" // ¡Importa fmt para usar Printf en la depuración!
	"html/template"
	"log"
	"net/http"
	// Asegúrate de que esta ruta sea correcta para tu paquete models
)

// Estructura para pasar datos al template del dashboard
type DashboardData struct {
	TotalBooks     int
	AvailableBooks int
	BorrowedBooks  int
}

// HomeHandler ahora recibe una conexión a la base de datos
func HomeHandler(db *sql.DB) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/home.html",
	))

	return func(w http.ResponseWriter, r *http.Request) {
		var totalBooks int
		// Contar el total de libros
		err := db.QueryRow("SELECT COUNT(*) FROM libros").Scan(&totalBooks)
		if err != nil {
			log.Printf("ERROR BD: Error al contar libros totales: %v", err) // Mensaje de error más claro
			http.Error(w, "Error interno del servidor al obtener datos del dashboard", http.StatusInternalServerError)
			return
		}
		// --- LÍNEA DE DEPURACIÓN CLAVE ---
		fmt.Printf("DEBUG HOME: Total de Libros obtenidos de DB: %d\n", totalBooks)

		var availableBooks int
		// Contar libros no prestados (disponibles)
		err = db.QueryRow("SELECT COUNT(*) FROM libros WHERE prestado = FALSE").Scan(&availableBooks)
		if err != nil {
			log.Printf("ERROR BD: Error al contar libros disponibles: %v", err) // Mensaje de error más claro
			http.Error(w, "Error interno del servidor al obtener datos del dashboard", http.StatusInternalServerError)
			return
		}
		// --- LÍNEA DE DEPURACIÓN CLAVE ---
		fmt.Printf("DEBUG HOME: Libros Disponibles obtenidos de DB: %d\n", availableBooks)

		var borrowedBooks int
		// Contar libros prestados
		err = db.QueryRow("SELECT COUNT(*) FROM libros WHERE prestado = TRUE").Scan(&borrowedBooks)
		if err != nil {
			log.Printf("ERROR BD: Error al contar libros prestados: %v", err) // Mensaje de error más claro
			http.Error(w, "Error interno del servidor al obtener datos del dashboard", http.StatusInternalServerError)
			return
		}
		// --- LÍNEA DE DEPURACIÓN CLAVE ---
		fmt.Printf("DEBUG HOME: Libros Prestados obtenidos de DB: %d\n", borrowedBooks)

		// Crear la estructura de datos para el template
		data := DashboardData{
			TotalBooks:     totalBooks,
			AvailableBooks: availableBooks,
			BorrowedBooks:  borrowedBooks,
		}

		// --- LÍNEA DE DEPURACIÓN CLAVE ---
		fmt.Printf("DEBUG HOME: Datos enviados al template: %+v\n", data)

		// Ejecutar el template con los datos obtenidos
		err = tmpl.ExecuteTemplate(w, "base", data)
		if err != nil {
			log.Printf("ERROR TEMPLATE: Error al ejecutar el template home: %v", err) // Mensaje de error más claro
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
	}
}
