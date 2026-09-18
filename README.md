# nopricey

nopricey or (pricey.no) aggregates weekly grocery discounts ("tilbud") from Norwegian
supermarket chains and matches them against a shopping list, so you know
which store is cheapest for that week's shopping.

## How it works

1. Fetchers pull each store's current weekly catalog from
   [Tjek](https://tjek.com/) (formerly eTilbudsavis), the third-party
   platform several Norwegian chains publish their catalogs through.
2. A search endpoint matches a product name against the fetched offers
   across all stores.
3. (Planned) A shopping list of multiple items is matched at once and
   compared per store, to recommend where to shop.

## Stores

| Store    | Status         |
|----------|----------------|
| Rema1000 | Implemented    |
| Kiwi     | Implemented    |
| Bunnpris | Implemented    |
| Extra    | Implemented    |
| Coop     | Not started    |

Weekly-discount data only for now. Normal (non-sale) pricing is a later
addition (if possible).

## Project layout

```
cmd/main.go          entry point, HTTP server setup
internal/api          routes and handlers
internal/compute      matching/search logic over fetched offers
internal/fetch        per-store fetchers (Tjek-backed) + shared client
internal/structs.go   API-facing response types
```

## Running locally

```
go run ./cmd
```

Serves on `:10000`. Requires Go 1.27+.

## API

```
GET /item?name=<query>
```

Returns matching offers for `name` across all stores, as JSON.

## Status

Early stage: backend-only, single-item search working, no frontend or
persistence yet. See [ROADMAP.md](ROADMAP.md) for the full plan.

## Roadmap

[ ] Cache weekly offers instead of fetching live per request.
[ ] Support matching a full shopping list, not just one item, and rank
  stores by total cost.
[ ] Minimal server-rendered frontend.
[ ] Dockerize and deploy publicly on a VPS behind a domain.
