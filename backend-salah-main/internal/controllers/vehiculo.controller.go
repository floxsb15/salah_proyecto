package controllers

import (
	"backend-restaurant-delitto/internal/db"
	"backend-restaurant-delitto/internal/functions"
	"backend-restaurant-delitto/internal/models"
	"backend-restaurant-delitto/internal/querys"
	"backend-restaurant-delitto/internal/security"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var errPrecioCompraNoAutorizado = errors.New("Solo el administrador puede editar el precio de compra")

type VehiculoDAO struct {
	ID                 uint     `json:"id"`
	Nombre             string   `json:"nombre"`
	Descripcion        string   `json:"descripcion"`
	Precio             float64  `json:"precio" format:"%.2f"`
	PrecioUSD          float64  `json:"precio_usd" format:"%.2f"`
	PrecioCompra       *float64 `json:"precio_compra,omitempty" format:"%.2f"`
	MargenGanancia     *float64 `json:"margen_ganancia,omitempty" format:"%.2f"`
	CantidadDisponible uint     `json:"cantidad_disponible"`
	Imagen             string   `json:"imagen"`
	Imagenes           []string `json:"imagenes" gorm:"-"`
	Estado             string   `json:"estado"`
	IDCategoria        uint     `json:"id_categoria"`
	IDSegmento         *uint    `json:"id_segmento"`
	Categoria          string   `json:"categoria"`
	Segmento           string   `json:"segmento"`
	Marca              string   `json:"marca"`
	Modelo             string   `json:"modelo"`
	Anio               uint     `json:"anio"`
	Version            string   `json:"version"`
	TipoTecho          string   `json:"tipo_techo"`
	Combustible        string   `json:"combustible"`
	Traccion           string   `json:"traccion"`
	Transmision        string   `json:"transmision"`
	Asientos           *uint    `json:"asientos"`
	Garantia           string   `json:"garantia"`
	Equipamiento       string   `json:"equipamiento"`
}

type VehiculoMOD struct {
	Nombre             string   `json:"nombre"`
	Descripcion        string   `json:"descripcion"`
	Precio             float64  `json:"precio" format:"%.2f"`
	PrecioUSD          float64  `json:"precio_usd" format:"%.2f"`
	PrecioCompra       *float64 `json:"precio_compra,omitempty" format:"%.2f"`
	MargenGanancia     *float64 `json:"margen_ganancia,omitempty" format:"%.2f"`
	CantidadDisponible uint     `json:"cantidad_disponible"`
	Imagen             string   `json:"imagen"`
	Imagenes           []string `json:"imagenes" gorm:"-"`
	Estado             string   `json:"estado"`
	IDCategoria        uint     `json:"id_categoria"`
	IDSegmento         *uint    `json:"id_segmento"`
	Categoria          string   `json:"categoria"`
	Segmento           string   `json:"segmento"`
	Marca              string   `json:"marca"`
	Modelo             string   `json:"modelo"`
	Anio               uint     `json:"anio"`
	Version            string   `json:"version"`
	TipoTecho          string   `json:"tipo_techo"`
	Combustible        string   `json:"combustible"`
	Traccion           string   `json:"traccion"`
	Transmision        string   `json:"transmision"`
	Asientos           *uint    `json:"asientos"`
	Garantia           string   `json:"garantia"`
	Equipamiento       string   `json:"equipamiento"`
}

func construirNombreVehiculo(marca string, modelo string, anio uint) string {
	partes := make([]string, 0, 3)
	if strings.TrimSpace(marca) != "" {
		partes = append(partes, strings.TrimSpace(marca))
	}
	if strings.TrimSpace(modelo) != "" {
		partes = append(partes, strings.TrimSpace(modelo))
	}
	if anio > 0 {
		partes = append(partes, strconv.FormatUint(uint64(anio), 10))
	}
	if len(partes) == 0 {
		return "Vehiculo"
	}
	return strings.Join(partes, " ")
}

func encodeImageToBase64(path string) (string, error) {
	imgBytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	mimeType := "image/jpeg"
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png":
		mimeType = "image/png"
	case ".jpg":
		mimeType = "image/jpg"
	case ".jpeg":
		mimeType = "image/jpeg"
	default:
		mimeType = "application/octet-stream"
	}

	encoded := base64.StdEncoding.EncodeToString(imgBytes)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}

func parseVehiculoImagenPaths(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" || value == "N/A" {
		return nil
	}

	var paths []string
	if strings.HasPrefix(value, "[") && json.Unmarshal([]byte(value), &paths) == nil {
		result := make([]string, 0, len(paths))
		for _, path := range paths {
			path = strings.TrimSpace(path)
			if path != "" && path != "N/A" {
				result = append(result, path)
			}
		}
		return result
	}

	return []string{value}
}

func encodeVehiculoImagenes(value string) ([]string, string) {
	paths := parseVehiculoImagenPaths(value)
	imagenes := make([]string, 0, len(paths))
	for _, path := range paths {
		encoded, err := encodeImageToBase64(path)
		if err == nil {
			imagenes = append(imagenes, encoded)
		}
	}

	if len(imagenes) == 0 {
		return nil, "N/A"
	}
	return imagenes, imagenes[0]
}

func guardarImagenesVehiculo(r *http.Request, required bool) (string, []string, error) {
	if r.MultipartForm == nil {
		return "N/A", nil, nil
	}

	fileHeaders := r.MultipartForm.File["fotos"]
	if len(fileHeaders) == 0 {
		fileHeaders = r.MultipartForm.File["foto"]
	}
	if required && len(fileHeaders) == 0 {
		return "", nil, errors.New("Debe subir entre 1 y 5 fotos")
	}
	if len(fileHeaders) == 0 {
		return "N/A", nil, nil
	}
	if len(fileHeaders) > 5 {
		return "", nil, errors.New("Solo se permiten hasta 5 fotos")
	}

	if err := os.MkdirAll("internal/images/vehiculos", 0750); err != nil {
		return "", nil, errors.New("Error al preparar la carpeta de fotos")
	}

	paths := make([]string, 0, len(fileHeaders))
	for _, fileHeader := range fileHeaders {
		extension, validationErr := validateImageUpload(fileHeader, maxVehiclePhotoSize)
		if validationErr != nil {
			eliminarImagenesVehiculo(paths)
			return "", nil, validationErr
		}
		file, err := fileHeader.Open()
		if err != nil {
			eliminarImagenesVehiculo(paths)
			return "", nil, errors.New("Error al obtener la foto")
		}

		nombreImagen := fmt.Sprintf("vehiculo-%s%s", uuid.New().String(), extension)
		rutaImagen := "internal/images/vehiculos/" + nombreImagen
		outFile, err := createUploadFile(rutaImagen)
		if err != nil {
			file.Close()
			eliminarImagenesVehiculo(paths)
			return "", nil, errors.New("Error al guardar la foto")
		}

		_, copyErr := io.Copy(outFile, file)
		closeErr := outFile.Close()
		file.Close()
		if copyErr != nil || closeErr != nil {
			_ = os.Remove(rutaImagen)
			eliminarImagenesVehiculo(paths)
			return "", nil, errors.New("Error al escribir la foto")
		}
		paths = append(paths, rutaImagen)
	}

	return serializarImagenesVehiculo(paths), paths, nil
}

func serializarImagenesVehiculo(paths []string) string {
	if len(paths) == 0 {
		return "N/A"
	}
	if len(paths) == 1 {
		return paths[0]
	}
	encoded, err := json.Marshal(paths)
	if err != nil {
		return paths[0]
	}
	return string(encoded)
}

func eliminarImagenesVehiculo(paths []string) {
	for _, path := range paths {
		if strings.TrimSpace(path) != "" && path != "N/A" {
			_ = os.Remove(path)
		}
	}
}

func ObtenerVehiculos(w http.ResponseWriter, r *http.Request) {
	vehiculos := make([]VehiculoDAO, 0)
	esAdmin := requestActorEsAdmin(r, "")

	err := db.GDB.Raw(querys.Vehiculos).Scan(&vehiculos).Error
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}

	for i, vehiculo := range vehiculos {
		if !esAdmin {
			vehiculos[i].PrecioCompra = nil
			vehiculos[i].MargenGanancia = nil
		}
		vehiculos[i].Imagenes, vehiculos[i].Imagen = encodeVehiculoImagenes(vehiculo.Imagen)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehiculos)
}

func ObtenerVehiculo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var vehiculo VehiculoMOD
	esAdmin := requestActorEsAdmin(r, "")

	err := db.GDB.Raw(querys.Vehiculo, id).Scan(&vehiculo).Error
	if err != nil {
		http.Error(w, "No existe el vehiculo", http.StatusInternalServerError)
		return
	}

	vehiculo.Imagenes, vehiculo.Imagen = encodeVehiculoImagenes(vehiculo.Imagen)
	if !esAdmin {
		vehiculo.PrecioCompra = nil
		vehiculo.MargenGanancia = nil
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehiculo)
}

func AgregarVehiculo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		http.Error(w, "Error al parsear el formulario", http.StatusInternalServerError)
		return
	}

	direccionImagen, imagenesGuardadas, err := guardarImagenesVehiculo(r, true)
	if err != nil {
		http.Error(w, "Imagen no valida", http.StatusBadRequest)
		return
	}

	nuevoPrecio, err := strconv.ParseFloat(r.FormValue("precio"), 64)
	if err != nil {
		http.Error(w, "Precio no valido", http.StatusBadRequest)
		return
	}
	precioCompra, err := parsePrecioCompraAdmin(r)
	if err != nil {
		http.Error(w, messagePrecioCompra(err), statusCodePrecioCompra(err))
		return
	}
	nuevoEstado, err := functions.ActualizarEstado(r.FormValue("estado"))
	if err != nil {
		http.Error(w, "Estado no valido", http.StatusBadRequest)
		return
	}
	cantidadDisponible, err := parseUintFormValueWithDefault(r.FormValue("cantidad_disponible"), 1)
	if err != nil || cantidadDisponible == 0 {
		http.Error(w, "Cantidad disponible no valida", http.StatusBadRequest)
		return
	}

	nuevoIdCategoria, err := strconv.ParseInt(r.FormValue("id_categoria"), 10, 64)
	if err != nil {
		http.Error(w, "ID no valida", http.StatusBadRequest)
		return
	}
	nuevoIdSegmento, err := parseOptionalUintFormValue(r.FormValue("id_segmento"))
	if err != nil {
		http.Error(w, "Segmento no valido", http.StatusBadRequest)
		return
	}
	if err := validarSegmentoCategoria(nuevoIdSegmento, uint(nuevoIdCategoria)); err != nil {
		http.Error(w, "Segmento no pertenece a la categoria", http.StatusBadRequest)
		return
	}
	asientos, err := parseOptionalUintFormValue(r.FormValue("asientos"))
	if err != nil {
		http.Error(w, "Numero de asientos no valido", http.StatusBadRequest)
		return
	}
	modelo := strings.TrimSpace(r.FormValue("modelo"))
	marca := strings.TrimSpace(r.FormValue("marca"))
	version := strings.TrimSpace(r.FormValue("version"))
	anio, err := parseAnioFormValue(r.FormValue("anio"))
	if err != nil {
		http.Error(w, "Anio requerido", http.StatusBadRequest)
		return
	}
	if modelo == "" {
		http.Error(w, "Modelo requerido", http.StatusBadRequest)
		return
	}
	if marca == "" {
		http.Error(w, "Marca requerida", http.StatusBadRequest)
		return
	}
	if err := validarMarcaVehiculo(marca); err != nil {
		http.Error(w, "Marca no valida", http.StatusBadRequest)
		return
	}
	if err := validarAnioVehiculo(anio); err != nil {
		http.Error(w, "Anio no valido", http.StatusBadRequest)
		return
	}

	nuevoVehiculo := models.Vehiculo{
		Nombre:             construirNombreVehiculo(marca, modelo, anio),
		Descripcion:        "",
		Precio:             nuevoPrecio,
		PrecioCompra:       precioCompra,
		CantidadDisponible: cantidadDisponible,
		Imagen:             direccionImagen,
		Estado:             nuevoEstado,
		IDCategoria:        uint(nuevoIdCategoria),
		IDSegmento:         nuevoIdSegmento,
		Marca:              marca,
		Modelo:             modelo,
		Anio:               anio,
		Version:            version,
		TipoTecho:          r.FormValue("tipo_techo"),
		Combustible:        r.FormValue("combustible"),
		Traccion:           r.FormValue("traccion"),
		Transmision:        r.FormValue("transmision"),
		Asientos:           asientos,
		Garantia:           r.FormValue("garantia"),
		Equipamiento:       r.FormValue("equipamiento"),
	}

	tx := db.GDB.Begin()
	if err := tx.Create(&nuevoVehiculo).Error; err != nil {
		tx.Rollback()
		respondInternalError(w, "Error al guardar el vehiculo", err)
		eliminarImagenesVehiculo(imagenesGuardadas)
		return
	}
	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nuevoVehiculo)
}

func ModificarVehiculo(w http.ResponseWriter, r *http.Request) {
	id_vehiculo := mux.Vars(r)["id"]
	var vehiculoExistente models.Vehiculo

	err := db.GDB.Where("id = ?", id_vehiculo).First(&vehiculoExistente).Error
	if err != nil {
		http.Error(w, "Vehiculo no encontrado", http.StatusNotFound)
		return
	}

	err = r.ParseMultipartForm(20 << 20)
	if err != nil {
		http.Error(w, "Solicitud no valida", http.StatusBadRequest)
		return
	}

	nuevaImagen, imagenesGuardadas, err := guardarImagenesVehiculo(r, false)
	if err != nil {
		http.Error(w, "Imagen no valida", http.StatusBadRequest)
		return
	}
	if len(imagenesGuardadas) > 0 {
		eliminarImagenesVehiculo(parseVehiculoImagenPaths(vehiculoExistente.Imagen))
		vehiculoExistente.Imagen = nuevaImagen
	}

	vehiculoExistente.Descripcion = ""
	nuevoPrecio, err := strconv.ParseFloat(r.FormValue("precio"), 64)
	if err != nil {
		http.Error(w, "Precio invalido", http.StatusBadRequest)
		return
	}
	vehiculoExistente.Precio = nuevoPrecio
	if r.FormValue("precio_compra") != "" {
		precioCompra, err := parsePrecioCompraAdmin(r)
		if err != nil {
			http.Error(w, messagePrecioCompra(err), statusCodePrecioCompra(err))
			return
		}
		vehiculoExistente.PrecioCompra = precioCompra
	}
	cantidadDisponible, err := parseUintFormValueWithDefault(r.FormValue("cantidad_disponible"), 1)
	if err != nil {
		http.Error(w, "Cantidad disponible no valida", http.StatusBadRequest)
		return
	}
	vehiculoExistente.CantidadDisponible = cantidadDisponible
	nuevoEstado, err := functions.ActualizarEstado(r.FormValue("estado"))
	if err != nil {
		http.Error(w, "Estado no valido", http.StatusBadRequest)
		return
	}
	vehiculoExistente.Estado = nuevoEstado
	nuevaCategoria, err := strconv.ParseInt(r.FormValue("id_categoria"), 10, 64)
	if err != nil {
		http.Error(w, "Categoria no valida", http.StatusBadRequest)
		return
	}
	vehiculoExistente.IDCategoria = uint(nuevaCategoria)
	modelo := strings.TrimSpace(r.FormValue("modelo"))
	marca := strings.TrimSpace(r.FormValue("marca"))
	version := strings.TrimSpace(r.FormValue("version"))
	anio, err := parseAnioFormValue(r.FormValue("anio"))
	if err != nil {
		http.Error(w, "Anio requerido", http.StatusBadRequest)
		return
	}
	if modelo == "" {
		http.Error(w, "Modelo requerido", http.StatusBadRequest)
		return
	}
	if marca == "" {
		http.Error(w, "Marca requerida", http.StatusBadRequest)
		return
	}
	if err := validarMarcaVehiculo(marca); err != nil {
		http.Error(w, "Marca no valida", http.StatusBadRequest)
		return
	}
	if err := validarAnioVehiculo(anio); err != nil {
		http.Error(w, "Anio no valido", http.StatusBadRequest)
		return
	}
	vehiculoExistente.Nombre = construirNombreVehiculo(marca, modelo, anio)
	nuevoSegmento, err := parseOptionalUintFormValue(r.FormValue("id_segmento"))
	if err != nil {
		http.Error(w, "Segmento no valido", http.StatusBadRequest)
		return
	}
	if err := validarSegmentoCategoria(nuevoSegmento, uint(nuevaCategoria)); err != nil {
		http.Error(w, "Segmento no pertenece a la categoria", http.StatusBadRequest)
		return
	}
	vehiculoExistente.IDSegmento = nuevoSegmento
	asientos, err := parseOptionalUintFormValue(r.FormValue("asientos"))
	if err != nil {
		http.Error(w, "Numero de asientos no valido", http.StatusBadRequest)
		return
	}
	vehiculoExistente.Marca = marca
	vehiculoExistente.Modelo = modelo
	vehiculoExistente.Anio = anio
	vehiculoExistente.Version = version
	vehiculoExistente.TipoTecho = r.FormValue("tipo_techo")
	vehiculoExistente.Combustible = r.FormValue("combustible")
	vehiculoExistente.Traccion = r.FormValue("traccion")
	vehiculoExistente.Transmision = r.FormValue("transmision")
	vehiculoExistente.Asientos = asientos
	vehiculoExistente.Garantia = r.FormValue("garantia")
	vehiculoExistente.Equipamiento = r.FormValue("equipamiento")

	if err := db.GDB.Save(&vehiculoExistente).Error; err != nil {
		http.Error(w, "Error al actualizar vehiculo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehiculoExistente)
}

func parseOptionalUintFormValue(value string) (*uint, error) {
	value = strings.TrimSpace(value)
	if value == "" || value == "0" || value == "<nil>" || value == "null" {
		return nil, nil
	}

	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return nil, err
	}

	result := uint(parsed)
	return &result, nil
}

func parseUintFormValueWithDefault(value string, defaultValue uint) (uint, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(parsed), nil
}

func parseAnioFormValue(value string) (uint, error) {
	value = strings.TrimSpace(value)
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil || parsed < 1900 || parsed > 2100 {
		return 0, fmt.Errorf("anio invalido")
	}

	return uint(parsed), nil
}

func parsePrecioCompraAdmin(r *http.Request) (float64, error) {
	value := strings.TrimSpace(r.FormValue("precio_compra"))
	if value == "" {
		return 0, nil
	}
	if !requestActorEsAdmin(r, "") {
		return 0, errPrecioCompraNoAutorizado
	}
	precioCompra, err := strconv.ParseFloat(value, 64)
	if err != nil || precioCompra < 0 {
		return 0, errors.New("Precio de compra no valido")
	}
	return precioCompra, nil
}

func statusCodePrecioCompra(err error) int {
	if errors.Is(err, errPrecioCompraNoAutorizado) {
		return http.StatusForbidden
	}
	return http.StatusBadRequest
}

func messagePrecioCompra(err error) string {
	if errors.Is(err, errPrecioCompraNoAutorizado) {
		return "Acceso denegado"
	}
	return "Precio de compra no valido"
}

func requestActorEsAdmin(r *http.Request, actorID string) bool {
	_ = actorID
	return security.CurrentUserHasRole(r, "admin")
}

func validarSegmentoCategoria(idSegmento *uint, idCategoria uint) error {
	if idSegmento == nil {
		return nil
	}

	var segmento models.SegmentoVehiculo
	return db.GDB.Where("id = ? and id_categoria = ?", *idSegmento, idCategoria).First(&segmento).Error
}

func validarMarcaVehiculo(nombre string) error {
	var marca models.MarcaVehiculo
	return db.GDB.Where("lower(nombre) = lower(?) and estado = true", strings.TrimSpace(nombre)).First(&marca).Error
}

func validarAnioVehiculo(anio uint) error {
	var anioVehiculo models.AnioVehiculo
	return db.GDB.Where("valor = ? and estado = true", anio).First(&anioVehiculo).Error
}
