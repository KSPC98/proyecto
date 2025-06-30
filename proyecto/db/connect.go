/*
@Autor: Kevin Pérez
@Descripcion: Módulo para establecer la conexión a la base de datos MySQL.
*/

package db

import (
	"database/sql" // Paquete para trabajar con bases de datos SQL.
	"fmt"          // Paquete para formatear cadenas.
	"log"          // Paquete para logging de errores y mensajes.
	"os"           // Paquete para interactuar con el sistema operativo (ej. variables de entorno).

	_ "github.com/go-sql-driver/mysql" // Driver de MySQL para Go. El guion bajo indica que se importa solo para sus efectos secundarios (inicializar el driver).
	"github.com/joho/godotenv"         // Paquete para cargar variables de entorno desde un archivo .env.
)

// Connect establece una conexión a la base de datos MySQL.
func Connect() (*sql.DB, error) {

	// Carga las variables de entorno desde el archivo .env.
	// Esto permite mantener la configuración de la base de datos fuera del código fuente.
	err := godotenv.Load()
	if err != nil {

		// Si hay un error al cargar el archivo .env, se devuelve el error.

		return nil, fmt.Errorf("error al cargar el archivo .env: %w", err)

	}

	// Construye la cadena de conexión DSN (Data Source Name) usando las variables de entorno.
	// Esto incluye el usuario, contraseña, host, puerto y nombre de la base de datos.
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),     // Usuario de la base de datos.
		os.Getenv("DB_PASSWORD"), // Contraseña de la base de datos.
		os.Getenv("DB_HOST"),     // Host de la base de datos (ej. localhost).
		os.Getenv("DB_PORT"),     // Puerto de la base de datos (ej. 3306).
		os.Getenv("DB_NAME"),     // Nombre de la base de datos a la que conectarse.
	)

	// Abre una conexión a la base de datos MySQL.
	// sql.Open no establece la conexión inmediatamente, solo valida los parámetros.
	db, err := sql.Open("mysql", dns)
	if err != nil {
		// Si hay un error al abrir la conexión, se devuelve el error.
		return nil, fmt.Errorf("error al abrir la conexión a la base de datos: %w", err)
	}

	// Intenta hacer ping a la base de datos para verificar que la conexión sea válida y esté activa.
	if err := db.Ping(); err != nil {
		// Si el ping falla, significa que la base de datos no está accesible con los datos proporcionados.
		db.Close() // Asegúrate de cerrar la conexión si falla el ping.
		return nil, fmt.Errorf("error al conectar con la base de datos: %w", err)
	}

	// Si todo es exitoso, se imprime un mensaje de éxito y se devuelve la conexión.
	log.Println("Conexión abierta a la base de datos exitosamente.")
	return db, nil
}
