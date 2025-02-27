#!/usr/bin/make -f

VERSION ?= "$(shell git describe --tags --match "v*" --dirty --always --abbrev=8 2>/dev/null || echo "develop")"
LOCODEDB ?= pkg/locodedb/data
UNLOCODEREVISION = 340a08558c84ae43122b86e97606bd2f5a771a06

.PHONY: all clean version help generate lint

DIRS = in ${LOCODEDB}

space := $(subst ,, )

all: $(DIRS) compress_locodedb

$(DIRS):
	@echo "⇒ Ensure dir: $@"
	@mkdir -p $@

in/airports.dat: | in
	wget -c https://raw.githubusercontent.com/jpatokal/openflights/master/data/airports.dat -O $@

in/countries.dat: | in
	wget -c https://raw.githubusercontent.com/jpatokal/openflights/master/data/countries.dat -O $@

# https://gist.githubusercontent.com/hrbrmstr/91ea5cc9474286c72838/raw/59421ff9b268ff0929b051ddafafbeb94a4c1910/continents.json
in/continents.geojson: continents.geojson.gz | in
	gunzip -c $< > $@

in/SubdivisionCodes.csv: | in
	wget -c https://raw.githubusercontent.com/datasets/un-locode/${UNLOCODEREVISION}/data/subdivision-codes.csv -O $@
	awk 'NR>1' $@ > temp && mv temp $@

in/CodeList.csv: | in
	wget -c https://raw.githubusercontent.com/datasets/un-locode/${UNLOCODEREVISION}/data/code-list.csv -O $@
	awk 'NR>1' $@ > temp && mv temp $@

generate: in/airports.dat in/countries.dat in/continents.geojson in/SubdivisionCodes.csv in/CodeList.csv | $(LOCODEDB)
	go run ./internal/generate/ \
	--airports in/airports.dat \
	--continents in/continents.geojson \
	--countries in/countries.dat \
	--in in/CodeList.csv \
	--in override.csv \
	--subdiv in/SubdivisionCodes.csv \
	--out $(LOCODEDB);

compress_locodedb: generate
	@echo "⇒ Compressing files inside $(LOCODEDB)"
	@for file in $(LOCODEDB)/*.csv; do \
	    if [ -f "$$file" ]; then \
	        bzip2 -cf "$$file" > "$$file.bz2"; \
	        rm "$$file"; \
	    fi \
	done

.golangci.yml:
	wget -O $@ https://github.com/nspcc-dev/.github/raw/master/.golangci.yml

# Lint Go code
lint: .golangci.yml
	golangci-lint run

# Print version
version:
	@echo $(VERSION)

# Show this help prompt
help:
	@echo '  Usage:'
	@echo ''
	@echo '    make <target>'
	@echo ''
	@echo '  Targets:'
	@echo ''
	@awk '/^#/{ comment = substr($$0,3) } comment && /^[a-zA-Z][a-zA-Z0-9_-]+ ?:/{ print "   ", $$1, comment }' $(MAKEFILE_LIST) | column -t -s ':' | grep -v 'IGNORE' | sort -u

# Clean up
clean:
	rm -f in/*

