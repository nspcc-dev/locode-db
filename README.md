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

This repository contains package for getting additional information about continent, country code and name, location, geo-position, subdivision based on locode.

[Locodes](https://github.com/nspcc-dev/locode-db/locodedb/locodes.csv.gz) and [countries](https://github.com/nspcc-dev/locode-db/locodedb/countries.csv.gz) databases can be regenerated using locode-db CLI command.

## Prerequisites

- [UN/LOCODE](https://unece.org/trade/cefact/UNLOCODE-Download)
  database in CSV format
- [OpenFlight Airports](https://raw.githubusercontent.com/jpatokal/openflights/master/data/airports.dat)
  database
- [OpenFlight Countries](https://raw.githubusercontent.com/jpatokal/openflights/master/data/countries.dat)
  database

## Quick start

Just run `make` to regenerate CSV files with [locodes](https://github.com/nspcc-dev/locode-db/locodedb/locodes.csv.gz) and [countries](https://github.com/nspcc-dev/locode-db/locodedb/countries.csv.gz).

``` shell
$ make
```
## License

This project is licensed under the MIT license - see the [LICENSE.md](LICENSE.md)
file for details.
