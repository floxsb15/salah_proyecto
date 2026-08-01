import { getRequestURL, proxyRequest } from 'h3'

export default defineEventHandler((event) => {
  const config = useRuntimeConfig(event)
  const backendUrl = new URL(config.backendUrl)
  if (!['http:', 'https:'].includes(backendUrl.protocol)) {
    throw createError({ statusCode: 500, statusMessage: 'Backend URL invalida' })
  }

  const requestUrl = getRequestURL(event)
  const target = new URL(requestUrl.pathname + requestUrl.search, backendUrl)
  return proxyRequest(event, target.toString())
})
