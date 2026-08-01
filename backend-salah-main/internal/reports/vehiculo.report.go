package reports

import (
	"backend-restaurant-delitto/internal/db"
	"fmt"
	"net/http"
	"time"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

var QueryVehiculos = `
	SELECT
		COALESCE(NULLIF(CONCAT_WS(' ', NULLIF(TRIM(p.marca), ''), NULLIF(TRIM(p.modelo), ''), NULLIF(p.anio::text, '0')), ''), p.nombre) as nombre,
		p.precio,
		COALESCE(cp.nombre, 'Sin Categoria') as categoria,
		COALESCE(sv.nombre, '') as segmento
	FROM vehiculos p
	LEFT JOIN categoria_vehiculo cp ON p.id_categoria = cp.id
	LEFT JOIN segmento_vehiculo sv ON p.id_segmento = sv.id
	WHERE p.estado = true
	ORDER BY cp.nombre, sv.nombre, p.nombre`

type VehiculoReporte struct {
	Nombre    string  `json:"nombre"`
	Precio    float64 `json:"precio"`
	Categoria string  `json:"categoria"`
	Segmento  string  `json:"segmento"`
}

func ReporteProductos(w http.ResponseWriter, r *http.Request) {
	m, err := makePDFVehiculos()
	if err != nil {
		respondInternalError(w, "Error al obtener datos del reporte de vehiculos", err)
		return
	}

	doc, err := m.Generate()
	if err != nil {
		respondInternalError(w, "Error al generar el reporte de vehiculos", err)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"reporte_vehiculos.pdf\"")

	if _, err := w.Write(doc.GetBytes()); err != nil {
		http.Error(w, "Error escribiendo PDF en la respuesta", http.StatusInternalServerError)
		return
	}
}

func makePDFVehiculos() (core.Maroto, error) {
	var vehiculos []VehiculoReporte
	err := db.GDB.Raw(QueryVehiculos).Scan(&vehiculos).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener vehiculos: %v", err)
	}

	if len(vehiculos) == 0 {
		return nil, fmt.Errorf("no se encontraron vehiculos")
	}

	cfg := config.NewBuilder().
		WithLeftMargin(15).
		WithRightMargin(15).
		WithTopMargin(15).
		WithBottomMargin(15).
		Build()

	m := maroto.New(cfg)

	m.AddRows(text.NewRow(12, "SALAH MOTORS", props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Align: align.Center,
		Size:  18,
	}))

	m.AddRow(8,
		text.NewCol(12, "REPORTE DE VEHICULOS", props.Text{
			Top:   1.0,
			Align: align.Center,
			Style: fontstyle.Bold,
			Size:  14,
		}),
	)

	fechaActual := time.Now().Format("02/01/2006 15:04:05")
	m.AddRow(8,
		text.NewCol(12, "Fecha de generacion: "+fechaActual, props.Text{
			Top:   1.0,
			Align: align.Center,
		}),
	)

	m.AddRows(line.NewRow(2))

	m.AddRow(8,
		text.NewCol(4, "VEHICULO", props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
		text.NewCol(2, "PRECIO (Bs)", props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
		text.NewCol(3, "CATEGORIA", props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
		text.NewCol(3, "SEGMENTO", props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
		}),
	)

	m.AddRows(line.NewRow(1))

	categoriaActual := ""
	for _, vehiculo := range vehiculos {
		if categoriaActual != vehiculo.Categoria {
			if categoriaActual != "" {
				m.AddRows(line.NewRow(1))
			}
			categoriaActual = vehiculo.Categoria
		}

		m.AddRow(7,
			text.NewCol(4, vehiculo.Nombre, props.Text{
				Top:   1.0,
				Align: align.Center,
			}),
			text.NewCol(2, fmt.Sprintf("%.2f", vehiculo.Precio), props.Text{
				Top:   1.0,
				Align: align.Center,
			}),
			text.NewCol(3, vehiculo.Categoria, props.Text{
				Top:   1.0,
				Align: align.Center,
			}),
			text.NewCol(3, vehiculo.Segmento, props.Text{
				Top:   1.0,
				Align: align.Center,
			}),
		)
	}

	m.AddRows(line.NewRow(2))

	m.AddRow(8,
		text.NewCol(12, fmt.Sprintf("Total de vehiculos: %d", len(vehiculos)), props.Text{
			Top:   1.0,
			Align: align.Right,
			Style: fontstyle.Bold,
		}),
	)

	return m, nil
}
