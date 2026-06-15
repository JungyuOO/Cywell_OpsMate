FROM quay.io/operator-framework/opm:latest

COPY deploy/olm/catalog /configs

ENTRYPOINT ["/bin/opm"]
CMD ["serve", "/configs"]
