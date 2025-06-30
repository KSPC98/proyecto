/*
@Autor: Kevin Pérez
@Descripcion: Módulo que maneja las operaciones CRUD para la entidad Libro en la base de datos.
*/

package models

import (
	"database/sql" // Paquete para interactuar con bases de datos SQL.
	"fmt"          // Paquete para formatear cadenas.
	"log"          // Paquete para logging de errores y mensajes.
	"proyecto/db"  // Importa el paquete db para obtener la conexión a la base de datos.
)

// Libro representa la estructura de un libro en la base de datos.
// Los nombres de los campos deben coincidir con los nombres de las columnas de la tabla.
type Libro struct {
	Id              int    // ID único del libro (clave primaria).
	Titulo          string // Título del libro.
	Autor           string // Autor del libro.
	AnioPublicacion int    // Año de publicación del libro.
	Editorial       string // Editorial del libro.
	Prestado        string // Estado de préstamo del libro (ej. "Si", "No").
}

// GetAllLibros consulta la base de datos y devuelve una lista de todos los libros.
func GetAllLibros() ([]Libro, error) {
	var libros []Libro // Declara una slice para almacenar los libros.
	// Establece una conexión a la base de datos.
	DB, err := db.Connect()
	if err != nil {
		log.Printf("Error al conectar a la base de datos en GetAllLibros: %v", err)
		return nil, fmt.Errorf("error al conectar a la base de datos: %w", err)
	}

	defer DB.Close() // Asegura que la conexión a la base de datos se cierre al finalizar la función.
	// Ejecuta la consulta SQL para seleccionar todos los campos de todos los libros.
	rows, err := DB.Query("SELECT Id, Titulo, Autor, AnioPublicacion, Editorial, Prestado FROM libros")
	if err != nil {
		log.Printf("Error al ejecutar la consulta en GetAllLibros: %v", err)
		return nil, fmt.Errorf("error al ejecutar la consulta: %w", err)
	}

	defer rows.Close() // Asegura que las filas de resultados se cierren al finalizar la función.
	// Itera sobre cada fila de resultados.
	for rows.Next() {
		var libro Libro // Declara una variable Libro para almacenar los datos de la fila actual.
		// Escanea los valores de la fila en los campos de la estructura Libro.
		err := rows.Scan(&libro.Id, &libro.Titulo, &libro.Autor, &libro.AnioPublicacion, &libro.Editorial, &libro.Prestado)
		if err != nil {
			log.Printf("Error al escanear los resultados en GetAllLibros: %v", err)
			return nil, fmt.Errorf("error al escanear los resultados: %w", err)
		}
		libros = append(libros, libro) // Agrega el libro a la slice de libros.
	}

	// Verifica si hubo algún error durante la iteración de las filas.
	if err = rows.Err(); err != nil {
		log.Printf("Error al procesar los resultados de la base de datos en GetAllLibros: %v", err)
		return nil, fmt.Errorf("error al procesar los resultados: %w", err)
	}

	return libros, nil // Devuelve la lista de libros y nil si no hay errores.
}

// CreateLibro inserta un nuevo libro en la base de datos.
func CreateLibro(Autor string, Titulo string, AnioPublicacion int, Editorial string, Prestado string) error {
	// Establece una conexión a la base de datos.
	DB, err := db.Connect()
	if err != nil {
		log.Printf("Error al conectar a la base de datos en CreateLibro: %v", err)
		return fmt.Errorf("error al conectar a la base de datos: %w", err)
	}
	defer DB.Close() // Asegura que la conexión se cierre.

	// Prepara la sentencia SQL para insertar un nuevo libro.
	// Esto ayuda a prevenir inyecciones SQL y mejora el rendimiento.
	stmt, err := DB.Prepare("INSERT INTO libros (Autor, Titulo, AnioPublicacion, Editorial, Prestado) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("Error al preparar la sentencia INSERT en CreateLibro: %v", err)
		return fmt.Errorf("error al preparar la sentencia: %w", err)
	}
	defer stmt.Close() // Asegura que la sentencia preparada se cierre.

	// Ejecuta la sentencia preparada con los valores proporcionados.
	resultado, err := stmt.Exec(Autor, Titulo, AnioPublicacion, Editorial, Prestado)
	if err != nil {
		log.Printf("Error al ejecutar la inserción del libro: %v", err)
		return fmt.Errorf("error al insertar el libro: %w", err)
	}

	// Obtiene el ID del último libro insertado.
	lastInsertId, err := resultado.LastInsertId()
	if err != nil {
		log.Printf("Error al obtener el ID del último libro insertado en CreateLibro: %v", err)
		return fmt.Errorf("error al obtener el ID del último libro insertado: %w", err)
	}
	log.Printf("Libro insertado con éxito. ID: %d", lastInsertId)

	return nil // Devuelve nil si la inserción fue exitosa.
}

// GetLibroByID consulta la base de datos y devuelve un libro específico por su ID.
func GetLibroByID(Id int) (Libro, error) {
	var libro Libro // Declara una variable Libro para almacenar el resultado.
	// Establece una conexión a la base de datos.
	DB, err := db.Connect()
	if err != nil {
		log.Printf("Error al conectar a la base de datos en GetLibroByID: %v", err)
		return libro, fmt.Errorf("error al conectar a la base de datos: %w", err)
	}
	defer DB.Close() // Asegura que la conexión se cierre.

	// Prepara la sentencia SQL para seleccionar un libro por su ID.
	stmt, err := DB.Prepare("SELECT Id, Titulo, Autor, AnioPublicacion, Editorial, Prestado FROM libros WHERE Id = ?")
	if err != nil {
		log.Printf("Error al preparar la consulta en GetLibroByID: %v", err)
		return libro, fmt.Errorf("error al preparar la consulta: %w", err)
	}
	defer stmt.Close() // Asegura que la sentencia preparada se cierre.

	// Ejecuta la consulta y escanea el resultado en la estructura Libro.
	fila := stmt.QueryRow(Id)
	err = fila.Scan(&libro.Id, &libro.Titulo, &libro.Autor, &libro.AnioPublicacion, &libro.Editorial, &libro.Prestado)
	if err != nil {
		if err == sql.ErrNoRows {
			// Si no se encuentra ninguna fila, devuelve un error específico.
			return libro, fmt.Errorf("libro con ID %d no encontrado", Id)
		}
		log.Printf("Error al escanear el libro con ID %d: %v", Id, err)
		return libro, fmt.Errorf("error al obtener el libro: %w", err)
	}
	log.Printf("Libro obtenido con éxito: %+v", libro)
	return libro, nil // Devuelve el libro y nil si no hay errores.
}

// UpdateLibro actualiza un libro existente en la base de datos.
func UpdateLibro(libro Libro) error {
	DB, err := db.Connect()
	if err != nil {
		log.Printf("Error al conectar a la base de datos en UpdateLibro: %v", err)
		return fmt.Errorf("error al conectar a la base de datos: %w", err)
	}
	defer DB.Close()

	// Prepara la sentencia SQL para actualizar un libro.
	stmt, err := DB.Prepare("UPDATE libros SET Titulo = ?, Autor = ?, AnioPublicacion = ?, Editorial = ?, Prestado = ? WHERE Id = ?")
	if err != nil {
		log.Printf("Error al preparar la sentencia UPDATE en UpdateLibro: %v", err)
		return fmt.Errorf("error al preparar la sentencia: %w", err)
	}
	defer stmt.Close()

	// Ejecuta la sentencia preparada con los datos actualizados del libro.
	_, err = stmt.Exec(libro.Titulo, libro.Autor, libro.AnioPublicacion, libro.Editorial, libro.Prestado, libro.Id)
	if err != nil {
		log.Printf("Error al ejecutar la actualización del libro con ID %d: %v", libro.Id, err)
		return fmt.Errorf("error al actualizar el libro: %w", err)
	}
	log.Printf("Libro con ID %d actualizado con éxito.", libro.Id)
	return nil
}

// DeleteLibro elimina un libro de la base de datos por su ID.
func DeleteLibro(Id int) error {
	DB, err := db.Connect()
	if err != nil {
		log.Printf("Error al conectar a la base de datos en DeleteLibro: %v", err)
		return fmt.Errorf("error al conectar a la base de datos: %w", err)
	}
	defer DB.Close()

	// Prepara la sentencia SQL para eliminar un libro.
	stmt, err := DB.Prepare("DELETE FROM libros WHERE Id = ?")
	if err != nil {
		log.Printf("Error al preparar la sentencia DELETE en DeleteLibro: %v", err)
		return fmt.Errorf("error al preparar la sentencia: %w", err)
	}
	defer stmt.Close()

	// Ejecuta la sentencia preparada con el ID del libro a eliminar.
	resultado, err := stmt.Exec(Id)
	if err != nil {
		log.Printf("Error al ejecutar la eliminación del libro con ID %d: %v", Id, err)
		return fmt.Errorf("error al eliminar el libro: %w", err)
	}

	filasAfectadas, err := resultado.RowsAffected()
	if err != nil {
		log.Printf("Error al obtener las filas afectadas en DeleteLibro: %v", err)
		return fmt.Errorf("error al obtener filas afectadas: %w", err)
	}

	if filasAfectadas == 0 {
		return fmt.Errorf("no se encontró ningún libro con ID %d para eliminar", Id)
	}
	log.Printf("Libro con ID %d eliminado con éxito.", Id)
	return nil
}
