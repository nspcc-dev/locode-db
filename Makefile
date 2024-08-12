#!/usr/bin/make -f

VERSION ?= "$(shell git describe --tags --match "v*" --dirty --always --abbrev=8 2>/dev/null || cat VERSION 2>/dev/null || echo "develop")"
LOCODECLI ?= locode-db
LOCODEDB ?= pkg/locodedb/data
UNLOCODEREVISION = 3648bfa776701c329d27136bef29fb3e21853f20

.PHONY: all clean version help unlocode generate $(LOCODECLI)

DIRS = in tmp ${LOCODEDB}

space := $(subst ,, )

all: $(DIRS) compress_locodedb

$(DIRS):
	@echo "⇒ Ensure dir: $@"
	@mkdir -p $@

in/airports.dat:
	wget -c https://raw.githubusercontent.com/jpatokal/openflights/master/data/airports.dat -O in/airports.dat

in/countries.dat:
	wget -c https://raw.githubusercontent.com/jpatokal/openflights/master/data/countries.dat -O in/countries.dat

in/continents.geojson: continents.geojson.gz
	gunzip -c $< > in/continents.geojson

in/SubdivisionCodes.csv: | in
	wget -c https://raw.githubusercontent.com/datasets/un-locode/${UNLOCODEREVISION}/data/subdivision-codes.csv -O $@
	awk 'NR>1' $@ > temp && mv temp $@

in/CodeList.csv: | in
	wget -c https://raw.githubusercontent.com/datasets/un-locode/${UNLOCODEREVISION}/data/code-list.csv -O $@
	awk 'NR>1' $@ > temp && mv temp $@

$(LOCODECLI):
	go build -o $(LOCODECLI) ./internal

generate: in/airports.dat in/countries.dat in/continents.geojson in/SubdivisionCodes.csv in/CodeList.csv $(LOCODECLI) | $(LOCODEDB)
	./$(LOCODECLI) generate \
	--airports in/airports.dat \
	--continents in/continents.geojson \
	--countries in/countries.dat \
	--in in/CodeList.csv \
	--subdiv in/SubdivisionCodes.csv \
	--out $(LOCODEDB)

compress_locodedb: generate
	@echo "⇒ Compressing files inside $(LOCODEDB)"
	@for file in $(LOCODEDB)/*.csv; do \
	    if [ -f "$$file" ]; then \
	        bzip2 -cf "$$file" > "$$file.bz2"; \
	        rm "$$file"; \
	    fi \
	done

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
	rm -f tmp/*
	rm -f $(LOCODECLI)

