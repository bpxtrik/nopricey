import { useState } from 'react'

function pkuUnit(quantity) {
  const unit = quantity.trim().split(' ').pop()
  return unit === 'l' || unit === 'ml' ? 'kr/l' : 'kr/kg'
}

function groupByStore(offers) {
  const map = new Map()
  for (const o of offers) {
    if (!map.has(o.store)) map.set(o.store, [])
    map.get(o.store).push(o)
  }
  return [...map.entries()].sort(([a], [b]) => a.localeCompare(b))
}

function App() {
  const [query, setQuery] = useState('')
  const [results, setResults] = useState(null)
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  function handleClear() {
    setQuery('')
    setResults(null)
    setMessage('')
    setError('')
  }

  async function handleSubmit(e) {
    e.preventDefault()
    if (!query.trim()) return

    setLoading(true)
    setError('')
    setMessage('')
    setResults(null)

    try {
      const res = await fetch(`/item?name=${encodeURIComponent(query)}`)
      const data = await res.json()

      if (!res.ok) {
        setError(data.detail ?? 'Something went wrong.')
      } else if (Array.isArray(data)) {
        setResults(data)
      } else {
        setMessage(data.detail ?? 'No results.')
      }
    } catch {
      setError('Could not reach the server.')
    } finally {
      setLoading(false)
    }
  }

  const stores = results ? groupByStore(results) : null

  return (
    <div className="min-h-screen bg-slate-50 flex flex-col items-center px-4 py-10 sm:px-6 sm:py-16">
      <div className="w-full max-w-5xl">
        <div className="max-w-xl">
          <h1 className="text-2xl font-semibold text-slate-900 mb-1 sm:text-3xl">
            nopricey
          </h1>
          <p className="text-slate-500 mb-6 sm:mb-8">
            Find where this week's deals are cheapest.
          </p>

          <form onSubmit={handleSubmit} className="flex gap-2 mb-8">
            <input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="e.g. melk"
              autoFocus
              className="min-w-0 flex-1 rounded-lg border border-slate-300 bg-white px-4 py-2 text-slate-900 shadow-sm outline-none focus:ring-2 focus:ring-indigo-400"
            />
            <button
              type="submit"
              disabled={loading}
              className="shrink-0 whitespace-nowrap rounded-lg bg-indigo-600 px-4 py-2 font-medium text-white shadow-sm transition hover:bg-indigo-500 disabled:opacity-50 sm:px-5"
            >
              {loading ? 'Searching…' : 'Search'}
            </button>
            <button
              type="button"
              onClick={handleClear}
              disabled={loading}
              className="shrink-0 whitespace-nowrap rounded-lg border border-slate-300 bg-white px-4 py-2 font-medium text-slate-600 shadow-sm transition hover:bg-slate-50 disabled:opacity-50 sm:px-5"
            >
              Clear
            </button>
          </form>

          {error && <p className="mb-4 text-red-600">{error}</p>}
          {message && <p className="mb-4 text-slate-500">{message}</p>}
        </div>

        {stores && stores.length > 0 && (
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {stores.map(([store, items]) => (
              <div
                key={store}
                className="flex min-w-0 flex-col gap-3 rounded-lg border border-slate-200 bg-white p-4 shadow-sm"
              >
                <h2 className="text-xs font-semibold tracking-wide text-indigo-700 uppercase">
                  {store}
                </h2>
                <ul className="flex flex-col gap-3">
                  {items.map((o, i) => (
                    <li
                      key={i}
                      className="flex items-start justify-between gap-3 border-t border-slate-100 pt-3 first:border-t-0 first:pt-0"
                    >
                      <div className="min-w-0">
                        <p className="break-words text-slate-900">
                          {o.product_name}
                        </p>
                        {o.quantity && (
                          <p className="text-sm text-slate-400">{o.quantity}</p>
                        )}
                      </div>
                      <span className="shrink-0 whitespace-nowrap text-right">
                        <span className="block font-semibold text-slate-900">
                          {o.price.toFixed(2)} kr
                        </span>
                        {o.price_per_kilo_unit > 0 && (
                          <span className="block text-xs text-slate-400">
                            {o.price_per_kilo_unit.toFixed(2)} {pkuUnit(o.quantity)}
                          </span>
                        )}
                      </span>
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

export default App
