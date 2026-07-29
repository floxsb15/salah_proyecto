package querys

var CategoriasVehiculos = `
	select cp.id, cp.nombre, cp.descripcion,
		case when cp.estado 
		then 'Activo'
		else 'Inactivo'
		end as estado
	from categoria_vehiculo cp
	order by cp.id asc;`

var CategoriaVehiculo = `
	select cp.nombre, cp.descripcion,
		case when cp.estado 
		then 'Activo'
		else 'Inactivo'
		end as estado
	from categoria_vehiculo cp
	where cp.id = ?
	limit 1;`
