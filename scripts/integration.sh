cd scripts/mocklicensing && \
       make build && \
       make start && \
       go test -v -count=1 -tags=integration ../../... ; \
       make stop && make clean