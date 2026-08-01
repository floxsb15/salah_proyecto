package routers

import (
	c "backend-restaurant-delitto/internal/controllers"
	r "backend-restaurant-delitto/internal/reports"
	"backend-restaurant-delitto/internal/security"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	roleAdmin    = "admin"
	roleManager  = "encargado de ventas"
	roleSeller   = "vendedor"
	maxBodyBytes = 30 << 20
)

var staffRoles = []string{roleAdmin, roleManager, roleSeller}

func InitEndPoints(router *mux.Router) {
	v1 := router.PathPrefix("/api/v1").Subrouter()

	v1.Handle("/login", security.LimitBody(16<<10)(http.HandlerFunc(c.Auth.AuthLoginWeb))).Methods(http.MethodPost)
	protected(v1, "/logout", c.CerrarSesion, http.MethodPost)
	protected(v1, "/me", c.ObtenerPerfilActual, http.MethodGet)
	protected(v1, "/me/password", c.CambiarContrasenaActual, http.MethodPatch)

	protected(v1, "/usuarios/{id}", c.ObtenerUsuario, http.MethodGet, roleAdmin)
	protected(v1, "/usuarios/{id}", c.ModificarUsuario, http.MethodPut, roleAdmin)
	protected(v1, "/usuarios", c.ObtenerUsuarios, http.MethodGet, roleAdmin)
	protected(v1, "/usuarios", c.AgregarUsuario, http.MethodPost, roleAdmin)

	protected(v1, "/clientes/{id}/historial-compras", c.ObtenerHistorialComprasCliente, http.MethodGet, staffRoles...)
	protected(v1, "/clientes/{id}", c.ObtenerCliente, http.MethodGet, staffRoles...)
	protected(v1, "/clientes/{id}", c.ModificarCliente, http.MethodPut, staffRoles...)
	protected(v1, "/clientes", c.ObtenerClientes, http.MethodGet, staffRoles...)
	protected(v1, "/clientes", c.AgregarCliente, http.MethodPost, staffRoles...)

	protected(v1, "/ventas/{id}/cuotas", c.ObtenerCuotasCreditoVenta, http.MethodGet, staffRoles...)
	protected(v1, "/ventas/{id}/completar-reserva", c.CompletarReservaVehiculo, http.MethodPatch, roleAdmin, roleManager)
	protected(v1, "/ventas/{id}/estado", c.ActualizarEstadoVentaVehiculo, http.MethodPatch, roleAdmin, roleManager)
	protected(v1, "/ventas/{id}/anular", c.AnularVentaVehiculo, http.MethodPatch, roleAdmin, roleManager)
	protected(v1, "/ventas/{id}", c.ObtenerVentaVehiculo, http.MethodGet, staffRoles...)
	protected(v1, "/ventas", c.ObtenerVentasVehiculos, http.MethodGet, staffRoles...)
	protected(v1, "/ventas", c.AgregarVentaVehiculo, http.MethodPost, staffRoles...)
	protected(v1, "/cuotas-credito/{id}/pagar", c.PagarCuotaCredito, http.MethodPatch, roleAdmin, roleManager)

	protected(v1, "/vehiculos/{id}", c.ObtenerVehiculo, http.MethodGet, staffRoles...)
	protected(v1, "/vehiculos/{id}", c.ModificarVehiculo, http.MethodPut, roleAdmin, roleManager)
	protected(v1, "/vehiculos", c.ObtenerVehiculos, http.MethodGet, staffRoles...)
	protected(v1, "/vehiculos", c.AgregarVehiculo, http.MethodPost, roleAdmin, roleManager)

	protected(v1, "/categorias-vehiculos/{id}", c.ObtenerCategoria, http.MethodGet, staffRoles...)
	protected(v1, "/categorias-vehiculos/{id}", c.ModificarCategoria, http.MethodPut, roleAdmin, roleManager)
	protected(v1, "/categorias-vehiculos", c.ObtenerCategorias, http.MethodGet, staffRoles...)
	protected(v1, "/categorias-vehiculos", c.AgregarCategoria, http.MethodPost, roleAdmin, roleManager)

	protected(v1, "/segmentos-vehiculos/{id}", c.ObtenerSegmentoVehiculo, http.MethodGet, staffRoles...)
	protected(v1, "/segmentos-vehiculos/{id}", c.ModificarSegmentoVehiculo, http.MethodPut, roleAdmin, roleManager)
	protected(v1, "/segmentos-vehiculos", c.ObtenerSegmentosVehiculos, http.MethodGet, staffRoles...)
	protected(v1, "/segmentos-vehiculos", c.AgregarSegmentoVehiculo, http.MethodPost, roleAdmin, roleManager)

	protected(v1, "/marcas-vehiculos", c.ObtenerMarcasVehiculos, http.MethodGet, staffRoles...)
	protected(v1, "/marcas-vehiculos", c.AgregarMarcaVehiculo, http.MethodPost, roleAdmin, roleManager)
	protected(v1, "/anios-vehiculos", c.ObtenerAniosVehiculos, http.MethodGet, staffRoles...)
	protected(v1, "/anios-vehiculos", c.AgregarAnioVehiculo, http.MethodPost, roleAdmin, roleManager)

	protected(v1, "/gastos/{id}", c.ObtenerGasto, http.MethodGet, roleAdmin)
	protected(v1, "/gastos/{id}", c.ModificarGasto, http.MethodPut, roleAdmin)
	protected(v1, "/gastos", c.ObtenerGastos, http.MethodGet, roleAdmin)
	protected(v1, "/gastos", c.AgregarGasto, http.MethodPost, roleAdmin)
	protected(v1, "/movimientos/{id}", c.AgregarMovimiento, http.MethodPost, roleAdmin)
	protected(v1, "/movimientos/{id}", c.ObtenerMovimientos, http.MethodGet, roleAdmin)

	protected(v1, "/reportes/vehiculos", r.ReporteProductos, http.MethodGet, staffRoles...)
	protected(v1, "/reportes/ventas/{id}", r.ReporteVentaVehiculo, http.MethodGet, staffRoles...)
}

func protected(router *mux.Router, path string, handler http.HandlerFunc, method string, roles ...string) {
	var secured http.Handler = handler
	if len(roles) > 0 {
		secured = security.RequireRoles(roles...)(secured)
	}
	secured = security.Authenticate(secured)
	secured = security.LimitBody(maxBodyBytes)(secured)
	router.Handle(path, secured).Methods(method)
}
