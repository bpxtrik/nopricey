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

- **Rema1000 / Kiwi**: reverse-engineered store app APIs. Blocked for now —
  Patrik needs to activate a Norwegian phone number to access the apps and
  capture the API calls. Don't build against guessed/placeholder endpoints;
  wait for real captured requests.
- **Bunnpris**: no API decided yet — check for a public webshop/kundeavis
  source before assuming reverse engineering is needed.
- **Extra / Coop**: not investigated yet.

Since sourcing differs per store, keep each store's fetcher isolated (own
package/file) behind a common interface, rather than assuming one scraping
strategy fits all chains.

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
