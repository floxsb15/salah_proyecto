package querys

var Clientes = `
	select c.id, c.nombre, c.apellido, c.ci, c.celular, c.direccion,
		case when c.estado
		then 'Activo'
		else 'Inactivo'
		end as estado
	from clientes c
	order by c.id asc;`

var Cliente = `
	select c.id, c.nombre, c.apellido, c.ci, c.celular, c.direccion,
		case when c.estado
		then 'Activo'
		else 'Inactivo'
		end as estado
	from clientes c
	where c.id = ?
	limit 1;`
