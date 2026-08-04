package reports

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/security"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type VentaVehiculoReporte struct {
	ID                       uint    `gorm:"column:id"`
	Fecha                    string  `gorm:"column:fecha"`
	FechaVenta               string  `gorm:"column:fecha_venta"`
	TipoVenta                string  `gorm:"column:tipo_venta"`
	Cantidad                 uint    `gorm:"column:cantidad"`
	PrecioUnidad             float64 `gorm:"column:precio_unidad"`
	PrecioTotal              float64 `gorm:"column:precio_total"`
	PrecioUSD                float64 `gorm:"column:precio_usd"`
	TipoCambioUsado          float64 `gorm:"column:tipo_cambio_usado"`
	MontoBOBCalculado        float64 `gorm:"column:monto_bob_calculado"`
	PagoUSD                  float64 `gorm:"column:pago_usd"`
	PagoBOB                  float64 `gorm:"column:pago_bob"`
	DetallePago              string  `gorm:"column:detalle_pago"`
	SaldoBOB                 float64 `gorm:"column:saldo_bob"`
	CuotaInicial             float64 `gorm:"column:cuota_inicial"`
	Saldo                    float64 `gorm:"column:saldo"`
	ValidezProformaDias      uint    `gorm:"column:validez_proforma_dias"`
	FechaVencimientoProforma string  `gorm:"column:fecha_vencimiento_proforma"`
	EstadoVenta              string  `gorm:"column:estado_venta"`
	EstadoPago               string  `gorm:"column:estado_pago"`
	MetodoPago               string  `gorm:"column:metodo_pago"`
	EstadoEntrega            string  `gorm:"column:estado_entrega"`
	FechaEntrega             string  `gorm:"column:fecha_entrega"`
	ReferenciaBancaria       string  `gorm:"column:referencia_bancaria"`
	EstadoDesembolso         string  `gorm:"column:estado_desembolso"`
	BancoEntidadFinanciera   string  `gorm:"column:banco_entidad_financiera"`
	EstadoTramiteBancario    string  `gorm:"column:estado_tramite_bancario"`
	MontoFinanciarBanco      float64 `gorm:"column:monto_financiar_banco"`
	FechaEstimadaDesembolso  string  `gorm:"column:fecha_estimada_desembolso"`
	Observacion              string  `gorm:"column:observacion"`
	Cliente                  string  `gorm:"column:cliente"`
	CICliente                string  `gorm:"column:ci_cliente"`
	Vehiculo                 string  `gorm:"column:vehiculo"`
	Categoria                string  `gorm:"column:categoria"`
	Segmento                 string  `gorm:"column:segmento"`
	Vendedor                 string  `gorm:"column:vendedor"`
}

var QueryVentaVehiculoReporte = `
	select
		vv.id,
		to_char(vv.fecha, 'YYYY-MM-DD') as fecha,
		to_char(coalesce(vv.fecha_venta, vv.fecha), 'YYYY-MM-DD HH24:MI:SS') as fecha_venta,
		vv.tipo_venta,
		vv.cantidad,
		vv.precio_unidad,
		vv.precio_total,
		coalesce(vv.precio_usd, vv.precio_total) as precio_usd,
		coalesce(vv.tipo_cambio_usado, 0) as tipo_cambio_usado,
		coalesce(vv.monto_bob_calculado, 0) as monto_bob_calculado,
		coalesce(pv.total_usd, vv.pago_usd, 0) as pago_usd,
		coalesce(pv.total_bob, vv.pago_bob, 0) as pago_bob,
		coalesce(pv.detalle_pago, '') as detalle_pago,
		coalesce(vv.saldo_bob, 0) as saldo_bob,
		vv.cuota_inicial,
		vv.saldo,
		vv.validez_proforma_dias,
		to_char(vv.fecha_vencimiento_proforma, 'YYYY-MM-DD') as fecha_vencimiento_proforma,
		vv.estado_venta,
		vv.estado_pago,
		case
			when coalesce(pv.cantidad, 0) > 1 then 'Mixto'
			when coalesce(pv.cantidad, 0) = 1 then pv.unico_metodo
			else coalesce(vv.metodo_pago, 'Efectivo')
		end as metodo_pago,
		vv.estado_entrega,
		coalesce(to_char(vv.fecha_entrega, 'YYYY-MM-DD'), '') as fecha_entrega,
		coalesce(vv.referencia_bancaria, '') as referencia_bancaria,
		coalesce(vv.estado_desembolso, '') as estado_desembolso,
		coalesce(vv.banco_entidad_financiera, '') as banco_entidad_financiera,
		coalesce(vv.estado_tramite_bancario, '') as estado_tramite_bancario,
		coalesce(vv.monto_financiar_banco, 0) as monto_financiar_banco,
		coalesce(to_char(vv.fecha_estimada_desembolso, 'YYYY-MM-DD'), '') as fecha_estimada_desembolso,
		coalesce(vv.observacion, '') as observacion,
		concat_ws(' ', c.nombre, c.apellido) as cliente,
		coalesce(c.ci, '') as ci_cliente,
		coalesce(nullif(concat_ws(' ', nullif(trim(v.marca), ''), nullif(trim(v.modelo), ''), nullif(v.anio::text, '0')), ''), v.nombre) as vehiculo,
		coalesce(cat.nombre, '') as categoria,
		coalesce(seg.nombre, '') as segmento,
		coalesce(concat_ws(' ', u.nombre, u.apellido), '') as vendedor
	from ventas_vehiculos vv
	inner join clientes c on c.id = vv.id_cliente
	inner join vehiculos v on v.id = vv.id_vehiculo
	left join categoria_vehiculo cat on cat.id = v.id_categoria
	left join segmento_vehiculo seg on seg.id = v.id_segmento
	left join usuarios u on u.id = vv.id_usuario
	left join lateral (
		select
			string_agg(concat(p.moneda, ' ', p.metodo, ' ', p.monto::text), ' | ' order by p.id) as detalle_pago,
			count(*) as cantidad,
			max(p.metodo) as unico_metodo,
			coalesce(sum(case when p.moneda = 'USD' then p.monto else 0 end), 0) as total_usd,
			coalesce(sum(case when p.moneda = 'BOB' then p.monto else 0 end), 0) as total_bob
		from pagos_venta p
		where p.venta_id = vv.id
	) pv on true
	where vv.id = ?
	limit 1`

func ReporteVentaVehiculo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	principal, ok := security.PrincipalFromContext(r.Context())
	if !ok {
		http.Error(w, "Autenticacion requerida", http.StatusUnauthorized)
		return
	}
	if !security.CurrentUserHasRole(r, "admin", "encargado de ventas", "contador") {
		var count int64
		if err := db.GDB.Table("ventas_vehiculos").Where("id = ? AND id_usuario = ?", id, principal.ID).Count(&count).Error; err != nil || count != 1 {
			http.Error(w, "Venta no encontrada", http.StatusNotFound)
			return
		}
	}
	m, venta, err := makePDFVentaVehiculo(id)
	if err != nil {
		respondInternalError(w, "Error al obtener datos del reporte de venta", err)
		return
	}

	doc, err := m.Generate()
	if err != nil {
		respondInternalError(w, "Error al generar el reporte de venta", err)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s_%d.pdf\"", nombreArchivoVentaPDF(venta), venta.ID))

	if _, err := w.Write(doc.GetBytes()); err != nil {
		http.Error(w, "Error escribiendo PDF en la respuesta", http.StatusInternalServerError)
		return
	}
}

func makePDFVentaVehiculo(id string) (core.Maroto, VentaVehiculoReporte, error) {
	var ventas []VentaVehiculoReporte
	if err := db.GDB.Raw(QueryVentaVehiculoReporte, id).Scan(&ventas).Error; err != nil {
		return nil, VentaVehiculoReporte{}, fmt.Errorf("error al obtener venta: %v", err)
	}
	if len(ventas) == 0 {
		return nil, VentaVehiculoReporte{}, fmt.Errorf("venta no encontrada")
	}

	venta := ventas[0]
	cfg := config.NewBuilder().
		WithLeftMargin(15).
		WithRightMargin(15).
		WithTopMargin(15).
		WithBottomMargin(15).
		Build()

	m := maroto.New(cfg)
	tituloDocumento := "DOCUMENTO DE VENTA"
	numeroDocumento := "Venta Nro"
	montoInicialLabel := "Cuota inicial"
	if venta.TipoVenta == "Reserva" {
		tituloDocumento = "PROFORMA DE RESERVA"
		numeroDocumento = "Proforma Nro"
		montoInicialLabel = "Monto reservado"
	}
	if venta.TipoVenta == "Reserva" && venta.EstadoVenta == "Completada" {
		tituloDocumento = "DOCUMENTO DE VENTA POR RESERVA"
		numeroDocumento = "Venta Nro"
	}
	m.AddRows(text.NewRow(12, "SALAH MOTORS", props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Align: align.Center,
		Size:  18,
	}))
	m.AddRow(8,
		text.NewCol(12, tituloDocumento, props.Text{
			Top:   1,
			Align: align.Center,
			Style: fontstyle.Bold,
			Size:  14,
		}),
	)
	m.AddRow(7,
		text.NewCol(6, fmt.Sprintf("%s: %d", numeroDocumento, venta.ID), props.Text{Top: 1, Style: fontstyle.Bold}),
		text.NewCol(6, "Generado: "+time.Now().Format("02/01/2006 15:04"), props.Text{Top: 1, Align: align.Right}),
	)
	m.AddRows(line.NewRow(2))

	addSectionTitle(m, "Datos de la venta")
	addPairRow(m, "Fecha", formatDatePDF(venta.Fecha), "Tipo", venta.TipoVenta)
	addPairRow(m, "Hora registro", formatDateTimePDF(venta.FechaVenta), "Vendedor", emptyPDF(venta.Vendedor))
	addPairRow(m, "Estado venta", venta.EstadoVenta, "Estado pago", venta.EstadoPago)
	addPairRow(m, "Metodo pago", venta.MetodoPago, "Estado entrega", venta.EstadoEntrega)
	addSingleRow(m, "Fecha entrega", emptyPDF(formatDatePDF(venta.FechaEntrega)))

	addSectionTitle(m, "Cliente y vendedor")
	addPairRow(m, "Cliente", venta.Cliente, "CI/NIT", emptyPDF(venta.CICliente))
	addPairRow(m, "Vehiculo", venta.Vehiculo, "Categoria", emptyPDF(venta.Categoria))
	addSingleRow(m, "Segmento", emptyPDF(venta.Segmento))

	addSectionTitle(m, "Detalle economico")
	addPairRow(m, "Cantidad", fmt.Sprintf("%d", venta.Cantidad), "Precio unidad", moneyPDF(venta.PrecioUnidad))
	addPairRow(m, "Total", moneyPDF(venta.PrecioTotal), montoInicialLabel, moneyPDF(venta.CuotaInicial))
	addPairRow(m, "Tipo cambio", fmt.Sprintf("%.4f Bs/USD", venta.TipoCambioUsado), "Total BOB", moneyBOBPDF(venta.MontoBOBCalculado))
	if venta.MetodoPago == "Mixto" {
		addPairRow(m, "Pago USD", moneyPDF(venta.PagoUSD), "Pago BOB", moneyBOBPDF(venta.PagoBOB))
		addSingleRow(m, "Detalle pago", emptyPDF(venta.DetallePago))
	}
	if venta.TipoVenta == "Reserva" && venta.EstadoVenta == "Completada" {
		addPairRow(m, "Pago final", moneyPDF(venta.PrecioTotal-venta.CuotaInicial), "Saldo", moneyPDF(venta.Saldo))
	} else {
		addPairRow(m, "Saldo", moneyPDF(venta.Saldo), "Validez proforma", fmt.Sprintf("%d dias", venta.ValidezProformaDias))
	}
	addPairRow(m, "Vence proforma", formatDatePDF(venta.FechaVencimientoProforma), "Desembolso", emptyPDF(venta.EstadoDesembolso))

	if venta.TipoVenta == "credito_bancario" {
		addSectionTitle(m, "Credito bancario")
		addPairRow(m, "Banco", emptyPDF(venta.BancoEntidadFinanciera), "Tramite", emptyPDF(venta.EstadoTramiteBancario))
		addPairRow(m, "Financia banco", moneyPDF(venta.MontoFinanciarBanco), "Fecha estimada", emptyPDF(formatDatePDF(venta.FechaEstimadaDesembolso)))
	}

	if venta.ReferenciaBancaria != "" || venta.Observacion != "" {
		addSectionTitle(m, "Notas")
		addSingleRow(m, "Referencia bancaria", emptyPDF(venta.ReferenciaBancaria))
		addSingleRow(m, "Observacion", emptyPDF(venta.Observacion))
	}

	m.AddRows(line.NewRow(2))
	m.AddRow(10,
		text.NewCol(6, "Firma vendedor", props.Text{Top: 4, Align: align.Center}),
		text.NewCol(6, "Firma cliente", props.Text{Top: 4, Align: align.Center}),
	)

	return m, venta, nil
}

func nombreArchivoVentaPDF(venta VentaVehiculoReporte) string {
	if venta.TipoVenta == "Reserva" && venta.EstadoVenta != "Completada" {
		return "proforma_reserva"
	}
	return "venta"
}

func addSectionTitle(m core.Maroto, title string) {
	m.AddRows(line.NewRow(1))
	m.AddRow(8, text.NewCol(12, title, props.Text{
		Top:   1,
		Style: fontstyle.Bold,
		Size:  11,
	}))
}

func addPairRow(m core.Maroto, label1 string, value1 string, label2 string, value2 string) {
	m.AddRow(7,
		text.NewCol(2, label1+":", props.Text{Top: 1, Style: fontstyle.Bold, Size: 9}),
		text.NewCol(4, emptyPDF(value1), props.Text{Top: 1, Size: 9}),
		text.NewCol(2, label2+":", props.Text{Top: 1, Style: fontstyle.Bold, Size: 9}),
		text.NewCol(4, emptyPDF(value2), props.Text{Top: 1, Size: 9}),
	)
}

func addSingleRow(m core.Maroto, label string, value string) {
	m.AddRow(8,
		text.NewCol(3, label+":", props.Text{Top: 1, Style: fontstyle.Bold, Size: 9}),
		text.NewCol(9, emptyPDF(value), props.Text{Top: 1, Size: 9}),
	)
}

func moneyPDF(value float64) string {
	return fmt.Sprintf("$ %.2f", value)
}

func moneyBOBPDF(value float64) string {
	return fmt.Sprintf("Bs %.2f", value)
}

func emptyPDF(value string) string {
	if value == "" {
		return "N/A"
	}
	return value
}

func formatDatePDF(value string) string {
	if value == "" {
		return ""
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return value
	}
	return parsed.Format("02/01/2006")
}

func formatDateTimePDF(value string) string {
	if value == "" {
		return ""
	}
	parsed, err := time.Parse("2006-01-02 15:04:05", value)
	if err != nil {
		return value
	}
	return parsed.Format("02/01/2006 15:04")
}
