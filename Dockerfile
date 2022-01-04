FROM scratch

COPY dnsstat /

ENTRYPOINT ["/dnsstat"]
