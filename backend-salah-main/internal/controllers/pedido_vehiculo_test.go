package controllers

import (
	"strings"
	"testing"
)

func TestConstruirPedidoVehiculoAllowsAdvanceOverRequired(t *testing.T) {
	payload := pedidoVehiculoPayloadForTest()
	payload.Pagos = []PagoVentaDAO{
		{Moneda: "USD", Metodo: "Efectivo", Monto: 150},
	}

	pedido, _, err := construirPedidoVehiculo(payload, 1)
	if err != nil {
		t.Fatalf("construirPedidoVehiculo() error = %v", err)
	}
	if pedido.AdelantoRequeridoUSD != 100 {
		t.Fatalf("AdelantoRequeridoUSD = %v", pedido.AdelantoRequeridoUSD)
	}
	if pedido.SaldoPendienteUSD != 850 {
		t.Fatalf("SaldoPendienteUSD = %v", pedido.SaldoPendienteUSD)
	}
}

func TestConstruirPedidoVehiculoRejectsAdvanceBelowRequiredWithBOBConversion(t *testing.T) {
	payload := pedidoVehiculoPayloadForTest()
	payload.Pagos = []PagoVentaDAO{
		{Moneda: "BOB", Metodo: "QR", Monto: 680},
	}

	_, _, err := construirPedidoVehiculo(payload, 1)
	if err == nil {
		t.Fatal("construirPedidoVehiculo() error = nil")
	}
	if !strings.Contains(err.Error(), "igual o mayor") {
		t.Fatalf("construirPedidoVehiculo() error = %v", err)
	}
}

func pedidoVehiculoPayloadForTest() PedidoVehiculoDAO {
	return PedidoVehiculoDAO{
		IDCliente:            1,
		Fecha:                "2026-08-15",
		Marca:                "Toyota",
		Modelo:               "Land Cruiser",
		Anio:                 2026,
		PaisOrigen:           "Japon",
		PrecioEstimadoUSD:    1000,
		TipoCambio:           6.96,
		FechaLlegadaEstimada: "2026-09-15",
		AdelantoPorcentaje:   10,
		ValidezDias:          15,
	}
}
