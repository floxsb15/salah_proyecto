package querys

var Gastos = `
	select 
		gv.id,
		gv.nombre,
		gv.unidad_medida
	from gastos_varios gv`

var Gasto = `
	select
		id,
		nombre, 
		unidad_medida
	from gastos_varios
	where id = ?`
