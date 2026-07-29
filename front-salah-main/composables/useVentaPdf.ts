import { server } from '~/server/server';

export function useVentaPdf() {
  async function descargarPDFVenta(idVenta: number | string) {
    const response = await fetch(`${server.HOST}/api/v1/reportes/ventas/${idVenta}`, {
      method: 'GET'
    });

    if (!response.ok) {
      const message = await response.text();
      throw new Error(message || 'No se pudo generar el PDF de la venta');
    }

    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `venta_${idVenta}.pdf`;
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(url);
  }

  return {
    descargarPDFVenta
  };
}
