package querys

var SegmentosVehiculos = `
	select s.id, s.nombre, s.descripcion, s.id_categoria, c.nombre as categoria,
		case when s.estado
		then 'Activo'
		else 'Inactivo'
		end as estado
	from segmento_vehiculo s
	left join categoria_vehiculo c on c.id = s.id_categoria`

var SegmentoVehiculo = `
	select s.id, s.nombre, s.descripcion, s.id_categoria, c.nombre as categoria,
		case when s.estado
		then 'Activo'
		else 'Inactivo'
		end as estado
	from segmento_vehiculo s
	left join categoria_vehiculo c on c.id = s.id_categoria
	where s.id = ?
	limit 1;`
