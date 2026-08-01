package functions

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/models"
	"backend-restaurant-delitto/internal/security"
	"errors"
	"os"
	"strings"
)

func CreacionInicial() error {
	if err := Roles(); err != nil {
		return err
	}
	if err := CategoriasVehiculos(); err != nil {
		return err
	}
	if err := SegmentosVehiculos(); err != nil {
		return err
	}
	if err := PrimerUsuario(); err != nil {
		return err
	}
	return nil
}

func Roles() error {
	if err := db.GDB.Model(&models.Rol{}).
		Where("nombre = ?", "ventas").
		Updates(map[string]interface{}{
			"nombre":   "encargado de ventas",
			"permisos": "ventas",
		}).Error; err != nil {
		return err
	}

	roles := []models.Rol{
		{Nombre: "admin", Permisos: "todos los permisos"},
		{Nombre: "encargado de ventas", Permisos: "ventas"},
		{Nombre: "vendedor", Permisos: "ventas"},
		{Nombre: "contador", Permisos: "dashboard, creditos, ventas, reservas"},
		{Nombre: "cliente", Permisos: "productos"},
	}

	for _, rol := range roles {
		var count int64
		if err := db.GDB.Model(&models.Rol{}).Where("nombre = ?", rol.Nombre).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			if err := db.GDB.Create(&rol).Error; err != nil {
				return err
			}
		} else if err := db.GDB.Model(&models.Rol{}).
			Where("nombre = ?", rol.Nombre).
			Update("permisos", rol.Permisos).Error; err != nil {
			return err
		}
	}

	return nil
}

func CategoriasVehiculos() error {
	categorias := []models.CategoriaVehiculo{
		{Nombre: "Automoviles", Descripcion: "Autos urbanos, familiares y particulares", Estado: true},
		{Nombre: "Camionetas", Descripcion: "Pickups y vehiculos utilitarios livianos", Estado: true},
		{Nombre: "Vagonetas y SUVs", Descripcion: "Vehiculos familiares altos, SUV y crossover", Estado: true},
		{Nombre: "Deportivos", Descripcion: "Vehiculos de alto rendimiento o estilo deportivo", Estado: true},
		{Nombre: "Electricos", Descripcion: "Vehiculos con motorizacion electrica", Estado: true},
		{Nombre: "Hibridos", Descripcion: "Vehiculos con motorizacion hibrida", Estado: true},
		{Nombre: "Transporte de pasajeros", Descripcion: "Vans, minibuses y buses", Estado: true},
		{Nombre: "Carga y camiones", Descripcion: "Camiones livianos, medianos y pesados", Estado: true},
		{Nombre: "Maquinaria", Descripcion: "Maquinaria pesada, agricola o de trabajo", Estado: true},
		{Nombre: "Motos", Descripcion: "Motocicletas, scooters y similares", Estado: true},
		{Nombre: "Vehiculos a pedido", Descripcion: "Unidades importadas o gestionadas bajo solicitud", Estado: true},
	}

	for _, categoria := range categorias {
		var count int64
		if err := db.GDB.Model(&models.CategoriaVehiculo{}).Where("nombre = ?", categoria.Nombre).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			if err := db.GDB.Create(&categoria).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func SegmentosVehiculos() error {
	segmentosPorCategoria := map[string][]models.SegmentoVehiculo{
		"Automoviles": {
			{Nombre: "Sedan", Descripcion: "Auto de cuatro puertas con maletero separado", Estado: true},
			{Nombre: "Hatchback", Descripcion: "Auto compacto con compuerta trasera", Estado: true},
			{Nombre: "Coupe", Descripcion: "Auto de dos puertas o perfil deportivo", Estado: true},
			{Nombre: "Convertible", Descripcion: "Auto descapotable", Estado: true},
		},
		"Camionetas": {
			{Nombre: "Pickup cabina simple", Descripcion: "Camioneta pickup de una fila", Estado: true},
			{Nombre: "Pickup cabina doble", Descripcion: "Camioneta pickup de dos filas", Estado: true},
			{Nombre: "Utilitaria", Descripcion: "Camioneta para trabajo o carga liviana", Estado: true},
		},
		"Vagonetas y SUVs": {
			{Nombre: "SUV compacta", Descripcion: "SUV de tamano compacto", Estado: true},
			{Nombre: "SUV mediana", Descripcion: "SUV familiar o de uso mixto", Estado: true},
			{Nombre: "SUV grande", Descripcion: "SUV de gran capacidad", Estado: true},
			{Nombre: "Crossover", Descripcion: "Vehiculo mixto entre automovil y SUV", Estado: true},
			{Nombre: "Vagoneta", Descripcion: "Vehiculo familiar tipo station wagon", Estado: true},
		},
		"Deportivos": {
			{Nombre: "Deportivo coupe", Descripcion: "Vehiculo deportivo cerrado", Estado: true},
			{Nombre: "Deportivo convertible", Descripcion: "Vehiculo deportivo descapotable", Estado: true},
			{Nombre: "Alto rendimiento", Descripcion: "Vehiculo orientado a potencia y prestaciones", Estado: true},
		},
		"Electricos": {
			{Nombre: "Auto electrico", Descripcion: "Automovil con propulsion electrica", Estado: true},
			{Nombre: "SUV electrica", Descripcion: "SUV con propulsion electrica", Estado: true},
			{Nombre: "Camioneta electrica", Descripcion: "Pickup o utilitario electrico", Estado: true},
			{Nombre: "Comercial electrico", Descripcion: "Vehiculo electrico para trabajo o transporte", Estado: true},
		},
		"Hibridos": {
			{Nombre: "Auto hibrido", Descripcion: "Automovil con sistema hibrido", Estado: true},
			{Nombre: "SUV hibrida", Descripcion: "SUV con sistema hibrido", Estado: true},
			{Nombre: "Camioneta hibrida", Descripcion: "Pickup o utilitario hibrido", Estado: true},
		},
		"Transporte de pasajeros": {
			{Nombre: "Van", Descripcion: "Vehiculo multipasajero liviano", Estado: true},
			{Nombre: "Minibus", Descripcion: "Transporte de pasajeros de capacidad media", Estado: true},
			{Nombre: "Bus", Descripcion: "Transporte de pasajeros de alta capacidad", Estado: true},
			{Nombre: "Coaster", Descripcion: "Bus pequeno o mediano tipo coaster", Estado: true},
		},
		"Carga y camiones": {
			{Nombre: "Camion liviano", Descripcion: "Camion de carga liviana", Estado: true},
			{Nombre: "Camion mediano", Descripcion: "Camion de carga mediana", Estado: true},
			{Nombre: "Camion pesado", Descripcion: "Camion de carga pesada", Estado: true},
			{Nombre: "Tracto camion", Descripcion: "Unidad para arrastre de semirremolque", Estado: true},
		},
		"Maquinaria": {
			{Nombre: "Construccion", Descripcion: "Maquinaria para obra y construccion", Estado: true},
			{Nombre: "Agricola", Descripcion: "Maquinaria para trabajo agricola", Estado: true},
			{Nombre: "Industrial", Descripcion: "Maquinaria para uso industrial", Estado: true},
		},
		"Motos": {
			{Nombre: "Moto urbana", Descripcion: "Motocicleta de uso urbano", Estado: true},
			{Nombre: "Scooter", Descripcion: "Moto scooter", Estado: true},
			{Nombre: "Moto deportiva", Descripcion: "Motocicleta deportiva", Estado: true},
		},
		"Vehiculos a pedido": {
			{Nombre: "Importado a pedido", Descripcion: "Vehiculo solicitado por importacion", Estado: true},
			{Nombre: "Especial", Descripcion: "Vehiculo con configuracion o uso especial", Estado: true},
		},
	}

	for nombreCategoria, segmentos := range segmentosPorCategoria {
		var categoria models.CategoriaVehiculo
		if err := db.GDB.Where("nombre = ?", nombreCategoria).First(&categoria).Error; err != nil {
			return err
		}

		for _, segmento := range segmentos {
			segmento.IDCategoria = categoria.ID
			var count int64
			if err := db.GDB.Model(&models.SegmentoVehiculo{}).
				Where("nombre = ? and id_categoria = ?", segmento.Nombre, categoria.ID).
				Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				if err := db.GDB.Create(&segmento).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func PrimerUsuario() error {

	var count int64
	if err := db.GDB.Model(&models.Usuario{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		initialUser := strings.TrimSpace(os.Getenv("INITIAL_ADMIN_USER"))
		initialPassword := os.Getenv("INITIAL_ADMIN_PASSWORD")
		if initialUser == "" || initialPassword == "" {
			return errors.New("INITIAL_ADMIN_USER e INITIAL_ADMIN_PASSWORD son obligatorios para crear el primer administrador")
		}
		hashedPassword, err := security.HashPassword(initialPassword)
		if err != nil {
			return err
		}
		var id_admin uint
		query := `select id from roles where nombre = 'admin' limit 1`
		if err := db.GDB.Raw(query).Scan(&id_admin).Error; err != nil {
			return err
		}
		usuario := models.Usuario{
			Usuario:        initialUser,
			Contra:         hashedPassword,
			IDRol:          id_admin,
			Estado:         true,
			SessionVersion: 1,
		}
		if err := db.GDB.Create(&usuario).Error; err != nil {
			return err
		}
	}
	return nil
}
