export function ok(res, data, message = 'success') {
  return res.json({ code: 0, message, data })
}

export function fail(res, code, message) {
  return res.status(code >= 400 ? code : 400).json({ code, message, data: null })
}
