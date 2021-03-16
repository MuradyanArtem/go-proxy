.PHONY: cert
cert:
	rm -rf ssl; mkdir ssl
	chmod +x scripts/gen_ca.sh
	./scripts/gen_ca.sh
