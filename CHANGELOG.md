# Changelog
Changelog for NeoFS LOCODE database

## [Unreleased]

## [0.8.1] - 2025-09-08

### Changed
- Minimal supported Go version is 1.24 now (#58)

### Updated
- golang.org/x/text dependency to v0.29.0 (#58)
- LOCODE DB revision, bringing CNFEH and FRBEC (with the same 2024-2 base version, #58)

## [0.8.0] - 2025-02-28

### Removed
- Useless exported types (#54)

### Updated
- UN/LOCODE to `2024-2` version (#51)
- github.com/stretchr/testify dependency to v1.10.0 (#55)
- golang.org/x/text dependency to v0.22.0 (#55)

## [0.7.0] - 2024-08-28

### Changed
- Go 1.22 is required to build now (#40, #43)
- Corrected known bad coordinates via a local override (#38)
- Drop cobra dependency, use standard packages to deal with flags (#42)

### Updated
- UN/LOCODE to `2024-1` version (#44)
- Switch data source to github.com/datasets/un-locode which has more suitable data format (#37)
- Continents file from known source (#45)

### Fixed
- UN/LOCODE files now in UTF-8, without wrong symbols (#37)
- Support float coordinates from UN/LOCODE files (#38)

## [0.6.0] - 2024-02-20

### Changed
- Go 1.20 is required to build now (#33)

### Updated
- UN/LOCODE to `2023-2` version (#33)
- All dependencies to current versions (#33)

## [0.5.0] - 2023-11-28

The DB is now provided as a Go package, import it as
`github.com/nspcc-dev/locode-db/pkg/locodedb`. Old binary DB is no longer
being built and released.

### Changed
- Dropped deb packages (#18)
- License is MIT now (#18)
- DB is stored in a Go package (#22, #23, #28)
- neofs-cli is no longer required to regenerate the DB (#22)

### Updated
- UN/LOCODE to `2023-1` version (#22)
- golang.org/x/text dependency from 0.3.7 to 0.3.8 (#27)

## [0.4.0] - 2023-04-11

### Updated
- UN/LOCODE to `2022-2` version

## [0.3.0] - 2022-10-28

### Changed
- Added Makefile to simplify DB build (#9)
- Added .deb packages support (#11)

### Updated
- UN/LOCODE to `2022-1` version

## [0.2.1] - 2021-11-02

### Changed
- Find the nearest continent for LOCODEs without exact continent match (#3, #6)

## [0.2.0] - 2021-10-21

### Fixed
- Decimal parts of coordinates contains minutes, not degrees (#2)

### Updated
- UN/LOCODE to `2021-1` version

## [0.1.0] - 2021-02-10

Initial release.

[0.1.0]: https://github.com/nspcc-dev/locode-db/releases/tag/v0.1.0
[0.2.0]: https://github.com/nspcc-dev/locode-db/compare/v0.1.0...v0.2.0
[0.2.1]: https://github.com/nspcc-dev/locode-db/compare/v0.2.0...v0.2.1
[0.3.0]: https://github.com/nspcc-dev/locode-db/compare/v0.2.1...v0.3.0
[0.4.0]: https://github.com/nspcc-dev/locode-db/compare/v0.3.0...v0.4.0
[0.5.0]: https://github.com/nspcc-dev/locode-db/compare/v0.4.0...v0.5.0
[0.6.0]: https://github.com/nspcc-dev/locode-db/compare/v0.5.0...v0.6.0
[0.7.0]: https://github.com/nspcc-dev/locode-db/compare/v0.6.0...v0.7.0
[0.8.0]: https://github.com/nspcc-dev/locode-db/compare/v0.7.0...v0.8.0
[Unreleased]: https://github.com/nspcc-dev/locode-db/compare/v0.8.0...master
