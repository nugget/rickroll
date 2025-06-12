FROM golang AS builder

WORKDIR /build
COPY . ./
RUN make rickrolld

FROM scratch
COPY --from=builder /build/dist/rickrolld /bin/rickrolld
COPY --from=builder /build/lyrics.dat /data/lyrics.dat
EXPOSE 23
ENTRYPOINT ["/bin/rickrolld"]
