# nopricey

Personal project (Norway) to find where weekly grocery deals are cheapest so
Patrik knows which store to go to for that week's shopping.

## What it does

- Aggregates weekly discounts ("tilbud") from Norwegian grocery chains.
- Patrik writes a shopping list; the tool matches it against current
  discounts across stores and returns where things are cheapest.
- Local web UI (not just a CLI) — Patrik types the list into a page and gets
  results back.

## Stores, in build order

1. Rema1000
2. Kiwi
3. Bunnpris
4. Extra
5. Coop

Start with the weekly-discount data for the first three; add normal (non-sale)
pricing and the remaining two chains later.

## Data sourcing (per store, not yet finalized for all)

- **Bunnpris / Rema1000 / Kiwi**: all three publish their weekly kundeavis
  through **Tjek** (formerly eTilbudsavis/ShopGun), a third-party catalog
  platform — `https://api.etilbudsavis.dk/v2/offers`, public, no auth
  required. Retailers submit their catalog to Tjek (as a PDF Tjek
  digitizes, or a structured feed), so this is legitimate published data,
  not a reverse-engineered app endpoint — no ToS/legal concern like the
  store apps would carry. Norwegian dealer IDs (found via
  `/v2/dealers/search?query=<name>&r_lat=59.9139&r_lng=10.7522&r_radius=50000`
  — the geolocation params matter, plain query search returns Danish
  dealers first): Bunnpris `5b11sm`, Rema1000 `faa0Ym`, Kiwi `257bxm`.
  Pagination is `offset`/`limit` (not `page`/`p`, which are silently
  ignored), `limit` capped at 100 server-side. Implemented in
  `internal/compute/tjek.go` (shared fetch/pagination) plus one file per
  store (`bunnpris.go`, `rema1000.go`, `kiwi.go`) with just the dealer ID.
  This makes the previous "reverse-engineered store app API, blocked on a
  Norwegian phone number" plan unnecessary for these three chains.
- **Extra / Coop**: not investigated yet — likely also on Tjek (Extra
  `80742m` confirmed live during the Bunnpris/Rema1000/Kiwi check; Coop's
  several sub-banners were spottier, worth re-checking before relying on
  it).

Since sourcing differs per store, keep each store's fetcher isolated (own
file) behind a common interface, rather than assuming one scraping strategy
fits all chains — even though Bunnpris/Rema1000/Kiwi happen to share Tjek as
a backend, a future chain may not.

## Legal/monetization context

This is a personal tool for now, not a monetized product. If monetization
ever comes up, note that price data via reverse-engineered APIs carries real
legal risk (ToS breach, possible anti-circumvention exposure, EU sui generis
database rights) — flag this rather than assuming it's fine, and don't add
paid/reselling features without that conversation happening explicitly.

## Current state

Bare skeleton: `cmd/main.go` is a `net/http` "Testing" stub on `:10000`, no
scrapers/fetchers, no web UI, no persistence yet. Treat architecture as
wide open — nothing here is a constraint yet beyond what's written above.
