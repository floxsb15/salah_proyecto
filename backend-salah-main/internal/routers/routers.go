package routers

import (
	c "backend-restaurant-delitto/internal/controllers"
	r "backend-restaurant-delitto/internal/reports"
	"net/http"

	"github.com/gorilla/mux"
)

func InitEndPoints(r *mux.Router) {
	api := r.PathPrefix("/api").Subrouter()
	endPointsAPI(api)
	reports(api)
}

func endPointsAPI(api *mux.Router) {
	v1 := api.PathPrefix("/v1").Subrouter()

	v1.HandleFunc("/login", c.Auth.AuthLoginWeb).Methods(http.MethodPost)

	v1Usuarios := v1.PathPrefix("/usuarios").Subrouter()
	v1Usuarios.HandleFunc("/{id}", c.ObtenerUsuario).Methods(http.MethodGet)
	v1Usuarios.HandleFunc("/{id}", c.ModificarUsuario).Methods(http.MethodPut)
	v1Usuarios.HandleFunc("", c.ObtenerUsuarios).Methods(http.MethodGet)
	v1Usuarios.HandleFunc("", c.AgregarUsuario).Methods(http.MethodPost)

	v1Clientes := v1.PathPrefix("/clientes").Subrouter()
	v1Clientes.HandleFunc("/{id}/historial-compras", c.ObtenerHistorialComprasCliente).Methods(http.MethodGet)
	v1Clientes.HandleFunc("/{id}", c.ObtenerCliente).Methods(http.MethodGet)
	v1Clientes.HandleFunc("/{id}", c.ModificarCliente).Methods(http.MethodPut)
	v1Clientes.HandleFunc("", c.ObtenerClientes).Methods(http.MethodGet)
	v1Clientes.HandleFunc("", c.AgregarCliente).Methods(http.MethodPost)

	v1Ventas := v1.PathPrefix("/ventas").Subrouter()
	v1Ventas.HandleFunc("/{id}/cuotas", c.ObtenerCuotasCreditoVenta).Methods(http.MethodGet)
	v1Ventas.HandleFunc("/{id}/completar-reserva", c.CompletarReservaVehiculo).Methods(http.MethodPatch)
	v1Ventas.HandleFunc("/{id}/estado", c.ActualizarEstadoVentaVehiculo).Methods(http.MethodPatch)
	v1Ventas.HandleFunc("/{id}/anular", c.AnularVentaVehiculo).Methods(http.MethodPatch)
	v1Ventas.HandleFunc("/{id}", c.ObtenerVentaVehiculo).Methods(http.MethodGet)
	v1Ventas.HandleFunc("", c.ObtenerVentasVehiculos).Methods(http.MethodGet)
	v1Ventas.HandleFunc("", c.AgregarVentaVehiculo).Methods(http.MethodPost)

	v1CuotasCredito := v1.PathPrefix("/cuotas-credito").Subrouter()
	v1CuotasCredito.HandleFunc("/{id}/pagar", c.PagarCuotaCredito).Methods(http.MethodPatch)

	v1Vehiculo := v1.PathPrefix("/vehiculos").Subrouter()
	v1Vehiculo.HandleFunc("/{id}", c.ObtenerVehiculo).Methods(http.MethodGet)
	v1Vehiculo.HandleFunc("/{id}", c.ModificarVehiculo).Methods(http.MethodPut)
	v1Vehiculo.HandleFunc("", c.ObtenerVehiculos).Methods(http.MethodGet)
	v1Vehiculo.HandleFunc("", c.AgregarVehiculo).Methods(http.MethodPost)

	v1CategoriaVehiculo := v1.PathPrefix("/categorias-vehiculos").Subrouter()
	v1CategoriaVehiculo.HandleFunc("/{id}", c.ObtenerCategoria).Methods(http.MethodGet)
	v1CategoriaVehiculo.HandleFunc("/{id}", c.ModificarCategoria).Methods(http.MethodPut)
	v1CategoriaVehiculo.HandleFunc("", c.ObtenerCategorias).Methods(http.MethodGet)
	v1CategoriaVehiculo.HandleFunc("", c.AgregarCategoria).Methods(http.MethodPost)

	v1SegmentoVehiculo := v1.PathPrefix("/segmentos-vehiculos").Subrouter()
	v1SegmentoVehiculo.HandleFunc("/{id}", c.ObtenerSegmentoVehiculo).Methods(http.MethodGet)
	v1SegmentoVehiculo.HandleFunc("/{id}", c.ModificarSegmentoVehiculo).Methods(http.MethodPut)
	v1SegmentoVehiculo.HandleFunc("", c.ObtenerSegmentosVehiculos).Methods(http.MethodGet)
	v1SegmentoVehiculo.HandleFunc("", c.AgregarSegmentoVehiculo).Methods(http.MethodPost)

	v1MarcaVehiculo := v1.PathPrefix("/marcas-vehiculos").Subrouter()
	v1MarcaVehiculo.HandleFunc("", c.ObtenerMarcasVehiculos).Methods(http.MethodGet)
	v1MarcaVehiculo.HandleFunc("", c.AgregarMarcaVehiculo).Methods(http.MethodPost)

	v1AnioVehiculo := v1.PathPrefix("/anios-vehiculos").Subrouter()
	v1AnioVehiculo.HandleFunc("", c.ObtenerAniosVehiculos).Methods(http.MethodGet)
	v1AnioVehiculo.HandleFunc("", c.AgregarAnioVehiculo).Methods(http.MethodPost)

	v1Gastos := v1.PathPrefix("/gastos").Subrouter()
	v1Gastos.HandleFunc("/{id}", c.ObtenerGasto).Methods(http.MethodGet)
	v1Gastos.HandleFunc("/{id}", c.ModificarGasto).Methods(http.MethodPut)
	v1Gastos.HandleFunc("", c.ObtenerGastos).Methods(http.MethodGet)
	v1Gastos.HandleFunc("", c.AgregarGasto).Methods(http.MethodPost)

	v1Movimientos := v1.PathPrefix("/movimientos").Subrouter()
	v1Movimientos.HandleFunc("/{id}", c.AgregarMovimiento).Methods(http.MethodPost)
	v1Movimientos.HandleFunc("/{id}", c.ObtenerMovimientos).Methods(http.MethodGet)
}

func reports(api *mux.Router) {
	v1 := api.PathPrefix("/v1").Subrouter()

	v1Reporte := v1.PathPrefix("/reportes").Subrouter()
	v1Reporte.HandleFunc("/vehiculos", r.ReporteProductos).Methods(http.MethodGet)
	v1Reporte.HandleFunc("/ventas/{id}", r.ReporteVentaVehiculo).Methods(http.MethodGet)
}
