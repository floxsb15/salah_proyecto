package querys

var Vehiculos = `
	select p.id,
		coalesce(nullif(concat_ws(' ', nullif(trim(p.marca), ''), nullif(trim(p.modelo), ''), nullif(p.anio::text, '0')), ''), p.nombre) as nombre,
		p.descripcion, p.precio, p.precio as precio_usd, coalesce(p.precio_compra, 0) as precio_compra,
		(p.precio - coalesce(p.precio_compra, 0)) as margen_ganancia,
		coalesce(p.cantidad_disponible, 1) as cantidad_disponible, p.imagen, p.id_categoria, p.id_segmento,
		p.marca, p.modelo, p.anio, p.version, p.tipo_techo, p.combustible, p.traccion, p.transmision,
		p.asientos, p.garantia, p.equipamiento,
		case when p.estado
		then 'Activo'
		else 'Inactivo'
		end as estado, cp.nombre as categoria, sv.nombre as segmento
	from vehiculos p 
	left join categoria_vehiculo as cp on cp.id = p.id_categoria
	left join segmento_vehiculo as sv on sv.id = p.id_segmento
	order by p.id asc;`

var Vehiculo = `
	select p.id,
		coalesce(nullif(concat_ws(' ', nullif(trim(p.marca), ''), nullif(trim(p.modelo), ''), nullif(p.anio::text, '0')), ''), p.nombre) as nombre,
		p.descripcion, p.precio, p.precio as precio_usd, coalesce(p.precio_compra, 0) as precio_compra,
		(p.precio - coalesce(p.precio_compra, 0)) as margen_ganancia,
		coalesce(p.cantidad_disponible, 1) as cantidad_disponible, p.imagen, p.id_categoria, p.id_segmento,
		p.marca, p.modelo, p.anio, p.version, p.tipo_techo, p.combustible, p.traccion, p.transmision,
		p.asientos, p.garantia, p.equipamiento,
		case when p.estado
		then 'Activo'
		else 'Inactivo'
		end as estado, cp.nombre as categoria, sv.nombre as segmento
	from vehiculos p 
	left join categoria_vehiculo as cp on cp.id = p.id_categoria
	left join segmento_vehiculo as sv on sv.id = p.id_segmento
	where p.id = ?
	limit 1;`

var VehiculoPorNombre = `
	select id from vehiculos
	where nombre = ?
	limit 1;
`
