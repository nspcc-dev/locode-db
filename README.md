<p align="center">
<img src="./.github/logo.svg" width="500px" alt="NeoFS">
</p>
<p align="center">
  UN/LOCODE database for <a href="https://fs.neo.org">NeoFS</a>
</p>

---
![GitHub release](https://img.shields.io/github/release/nspcc-dev/neofs-locode-db.svg)
![GitHub license](https://img.shields.io/github/license/nspcc-dev/neofs-locode-db.svg?style=popout)

# Overview

This repository contains instructions to generate UN/LOCODE database for NeoFS 
and raw representation of it. NeoFS uses UN/LOCODE in storage node attributes
and storage policies. Inner ring nodes converts UN/LOCODE into human-readable 
set of attributes such as continent, country name, etc. You can find out 
more in [NeoFS Specification](https://github.com/nspcc-dev/neofs-spec).


# Build

## Prerequisites

- Latest [neofs-cli](https://github.com/nspcc-dev/neofs-node)
- [UN/LOCODE](https://unece.org/trade/cefact/UNLOCODE-Download) 
  database in CSV format
- [OpenFlight Airports](https://raw.githubusercontent.com/jpatokal/openflights/master/data/airports.dat)
  database
- [OpenFlight Countries](https://raw.githubusercontent.com/jpatokal/openflights/master/data/countries.dat)
  database

## Quick start

Just run `make` to generate `locode_db` file for use with NeoFS InnerRing nodes.

``` shell
$ make
...
--out locode_db
```

## Building

First unzip file with GeoJSON continents from this repository.
```
$ gunzip continents.geojson
```

Then run neofs-cli command to generate boltDB file.
```
$ neofs-cli util locode generate --help
generate UN/LOCODE database for NeoFS

Usage:
  neofs-cli util locode generate [flags]

Flags:
      --airports string     Path to OpenFlights airport database (csv)
      --continents string   Path to continent polygons (GeoJSON)
      --countries string    Path to OpenFlights country database (csv)
  -h, --help                help for generate
      --in strings          List of paths to UN/LOCODE tables (csv)
      --out string          Target path for generated database
      --subdiv string       Path to UN/LOCODE subdivision database (csv)
      
$ ./neofs-cli util locode generate \
  --airports airports.dat \
  --continents continents.geojson \
  --countries countries.dat \
  --in 2020-2\ UNLOCODE\ CodeListPart1.csv,2020-2\ UNLOCODE\ CodeListPart2.csv,2020-2\ UNLOCODE\ CodeListPart3.csv \
  --subdiv 2020-2\ SubdivisionCodes.csv \
  --out locode_db
```

**Database generation might take some time!**

You can test generated database with neofs-cli.
```
$ neofs-cli util locode info --db locode_db --locode 'RU LED'
Country: Russia
Location: Saint Petersburg (ex Leningrad)
Continent: Europe
Subdivision: [SPE] Sankt-Peterburg
Coordinates: 59.53, 30.15
```


## License

This project is licensed under the CC Attribution-ShareAlike 4.0 International -
see the [LICENSE](LICENSE) file for details
