package reports

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/security"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type ProformaVehicularReporte struct {
	ID               uint    `gorm:"column:id"`
	Fecha            string  `gorm:"column:fecha"`
	ClienteNombre    string  `gorm:"column:cliente_nombre"`
	ClienteDireccion string  `gorm:"column:cliente_direccion"`
	ClienteTelefono  string  `gorm:"column:cliente_telefono"`
	Modalidad        string  `gorm:"column:modalidad"`
	PrecioUnidad     float64 `gorm:"column:precio_unidad"`
	Cantidad         uint    `gorm:"column:cantidad"`
	PrecioTotal      float64 `gorm:"column:precio_total"`
	CuotaInicial     float64 `gorm:"column:cuota_inicial"`
	Saldo            float64 `gorm:"column:saldo"`
	ValidezDias      uint    `gorm:"column:validez_dias"`
	FechaVencimiento string  `gorm:"column:fecha_vencimiento"`
	PrecioCatalogo   float64 `gorm:"column:precio_catalogo_ref"`
	Vehiculo         string  `gorm:"column:vehiculo"`
	Marca            string  `gorm:"column:marca"`
	Modelo           string  `gorm:"column:modelo"`
	Anio             uint    `gorm:"column:anio"`
	Version          string  `gorm:"column:version"`
	Garantia         string  `gorm:"column:garantia"`
	Equipamiento     string  `gorm:"column:equipamiento"`
	Descripcion      string  `gorm:"column:descripcion"`
	TipoTecho        string  `gorm:"column:tipo_techo"`
	Combustible      string  `gorm:"column:combustible"`
	Traccion         string  `gorm:"column:traccion"`
	Transmision      string  `gorm:"column:transmision"`
	Asientos         *uint   `gorm:"column:asientos"`
	Categoria        string  `gorm:"column:categoria"`
	Segmento         string  `gorm:"column:segmento"`
	Imagen           string  `gorm:"column:imagen"`
	Vendedor         string  `gorm:"column:vendedor"`
}

var QueryProformaVehicularReporte = `
	select
		p.id,
		to_char(p.fecha, 'YYYY-MM-DD HH24:MI:SS') as fecha,
		p.cliente_nombre,
		p.cliente_direccion,
		p.cliente_telefono,
		p.modalidad,
		p.precio_unidad,
		p.cantidad,
		p.precio_total,
		p.cuota_inicial,
		p.saldo,
		p.validez_dias,
		to_char(p.fecha_vencimiento, 'YYYY-MM-DD') as fecha_vencimiento,
		p.precio_catalogo_ref,
		coalesce(nullif(concat_ws(' ', nullif(trim(v.marca), ''), nullif(trim(v.modelo), ''), nullif(v.anio::text, '0')), ''), v.nombre) as vehiculo,
		coalesce(v.marca, '') as marca,
		coalesce(v.modelo, '') as modelo,
		coalesce(v.anio, 0) as anio,
		coalesce(v.version, '') as version,
		coalesce(v.garantia, '') as garantia,
		coalesce(v.equipamiento, '') as equipamiento,
		coalesce(v.descripcion, '') as descripcion,
		coalesce(v.tipo_techo, '') as tipo_techo,
		coalesce(v.combustible, '') as combustible,
		coalesce(v.traccion, '') as traccion,
		coalesce(v.transmision, '') as transmision,
		v.asientos,
		coalesce(cat.nombre, '') as categoria,
		coalesce(seg.nombre, '') as segmento,
		coalesce(v.imagen, '') as imagen,
		coalesce(concat_ws(' ', u.nombre, u.apellido), '') as vendedor
	from proformas_vehiculares p
	inner join vehiculos v on v.id = p.id_vehiculo
	left join categoria_vehiculo cat on cat.id = v.id_categoria
	left join segmento_vehiculo seg on seg.id = v.id_segmento
	left join usuarios u on u.id = p.id_usuario
	where p.id = ?
	limit 1`

func ReporteProformaVehicular(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	principal, ok := security.PrincipalFromContext(r.Context())
	if !ok {
		http.Error(w, "Autenticacion requerida", http.StatusUnauthorized)
		return
	}
	if !security.CurrentUserHasRole(r, "admin", "encargado de ventas") {
		var count int64
		if err := db.GDB.Table("proformas_vehiculares").Where("id = ? AND id_usuario = ?", id, principal.ID).Count(&count).Error; err != nil || count != 1 {
			http.Error(w, "Proforma no encontrada", http.StatusNotFound)
			return
		}
	}

	m, proforma, err := makePDFProformaVehicular(id)
	if err != nil {
		respondInternalError(w, "Error al obtener datos de la proforma", err)
		return
	}

	doc, err := m.Generate()
	if err != nil {
		respondInternalError(w, "Error al generar la proforma", err)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"proforma_vehicular_%d.pdf\"", proforma.ID))
	if _, err := w.Write(doc.GetBytes()); err != nil {
		http.Error(w, "Error escribiendo PDF en la respuesta", http.StatusInternalServerError)
		return
	}
}

func makePDFProformaVehicular(id string) (core.Maroto, ProformaVehicularReporte, error) {
	var proformas []ProformaVehicularReporte
	if err := db.GDB.Raw(QueryProformaVehicularReporte, id).Scan(&proformas).Error; err != nil {
		return nil, ProformaVehicularReporte{}, fmt.Errorf("error al obtener proforma: %v", err)
	}
	if len(proformas) == 0 {
		return nil, ProformaVehicularReporte{}, fmt.Errorf("proforma no encontrada")
	}

	proforma := proformas[0]
	cfg := config.NewBuilder().
		WithLeftMargin(14).
		WithRightMargin(14).
		WithTopMargin(12).
		WithBottomMargin(12).
		Build()

	m := maroto.New(cfg)
	m.AddRow(10,
		text.NewCol(5, "SALAH MOTORS", props.Text{Top: 1, Style: fontstyle.Bold, Size: 18}),
		text.NewCol(7, "PROFORMA VEHICULAR", props.Text{Top: 2, Align: align.Right, Style: fontstyle.Bold, Size: 14}),
	)
	m.AddRow(7,
		text.NewCol(7, "Direccion: Av. principal | Celular: 70000000 | Email: ventas@salahmotors.com", props.Text{Top: 1, Size: 8}),
		text.NewCol(5, fmt.Sprintf("Proforma Nro: %d", proforma.ID), props.Text{Top: 1, Align: align.Right, Style: fontstyle.Bold, Size: 9}),
	)
	m.AddRow(6,
		text.NewCol(6, "Fecha: "+formatDateTimePDF(proforma.Fecha), props.Text{Top: 1, Size: 8}),
		text.NewCol(6, "Vendedor: "+emptyPDF(proforma.Vendedor), props.Text{Top: 1, Align: align.Right, Size: 8}),
	)
	m.AddRows(line.NewRow(2))

	addSectionTitle(m, "Informacion del cliente")
	addPairRow(m, "Cliente", proforma.ClienteNombre, "Telefono", proforma.ClienteTelefono)
	addSingleRow(m, "Direccion", proforma.ClienteDireccion)

	addSectionTitle(m, "Informacion del vehiculo")
	imagePath := firstExistingImagePath(proforma.Imagen)
	if imagePath != "" {
		m.AddRow(34,
			image.NewFromFileCol(4, imagePath, props.Rect{Center: true, Percent: 90}),
			text.NewCol(8, vehicleSummaryPDF(proforma), props.Text{Top: 2, Size: 9}),
		)
	} else {
		addSingleRow(m, "Vehiculo", proforma.Vehiculo)
	}
	addPairRow(m, "Marca", proforma.Marca, "Modelo", proforma.Modelo)
	addPairRow(m, "Anio", fmt.Sprintf("%d", proforma.Anio), "Version", proforma.Version)
	addPairRow(m, "Categoria", proforma.Categoria, "Segmento", proforma.Segmento)
	addPairRow(m, "Garantia", proforma.Garantia, "Equipamiento", proforma.Equipamiento)

	addSectionTitle(m, "Especificaciones tecnicas")
	asientos := "N/A"
	if proforma.Asientos != nil {
		asientos = fmt.Sprintf("%d", *proforma.Asientos)
	}
	addPairRow(m, "Combustible", proforma.Combustible, "Traccion", proforma.Traccion)
	addPairRow(m, "Transmision", proforma.Transmision, "Asientos", asientos)
	addPairRow(m, "Tipo techo", proforma.TipoTecho, "Modalidad", proforma.Modalidad)
	addSingleRow(m, "Descripcion", proforma.Descripcion)

	addSectionTitle(m, "Precio y financiamiento")
	addPairRow(m, "Precio unidad", moneyPDF(proforma.PrecioUnidad), "Cantidad", fmt.Sprintf("%d", proforma.Cantidad))
	addPairRow(m, "Precio total", moneyPDF(proforma.PrecioTotal), "Cuota inicial", moneyPDF(proforma.CuotaInicial))
	addPairRow(m, "Saldo", moneyPDF(proforma.Saldo), "Precio catalogo", moneyPDF(proforma.PrecioCatalogo))
	addPairRow(m, "Validez", fmt.Sprintf("%d dias", proforma.ValidezDias), "Vence", formatDatePDF(proforma.FechaVencimiento))

	m.AddRows(line.NewRow(2))
	m.AddRow(8, text.NewCol(12, "Esta proforma no modifica el precio base del vehiculo en catalogo.", props.Text{Top: 1, Align: align.Center, Size: 8}))
	m.AddRow(14,
		text.NewCol(6, "____________________________\nFirma vendedor", props.Text{Top: 3, Align: align.Center, Size: 9}),
		text.NewCol(6, "____________________________\nFirma cliente", props.Text{Top: 3, Align: align.Center, Size: 9}),
	)
	m.AddRows(line.NewRow(1))
	m.AddRow(8, text.NewCol(12, "Salah Motors | Toyota | Nissan | Suzuki | Hyundai | Contacto comercial", props.Text{Top: 2, Align: align.Center, Size: 8}))

	return m, proforma, nil
}

func vehicleSummaryPDF(proforma ProformaVehicularReporte) string {
	parts := []string{
		"Vehiculo: " + emptyPDF(proforma.Vehiculo),
		"Marca: " + emptyPDF(proforma.Marca),
		"Modelo: " + emptyPDF(proforma.Modelo),
		fmt.Sprintf("Anio: %d", proforma.Anio),
		"Version: " + emptyPDF(proforma.Version),
	}
	return strings.Join(parts, "\n")
}

func firstExistingImagePath(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || value == "N/A" {
		return ""
	}

	candidates := []string{value}
	if strings.HasPrefix(value, "[") {
		var parsed []string
		if err := json.Unmarshal([]byte(value), &parsed); err == nil {
			candidates = parsed
		}
	}

	for _, candidate := range candidates {
		candidate = strings.TrimSpace(candidate)
		if candidate == "" || candidate == "N/A" {
			continue
		}
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ""
}
