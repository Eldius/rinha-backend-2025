
env-down:
	docker compose \
		-f docker-compose.yaml \
		-f resources/providers/docker-compose-providers.yaml \
			down


run-local: env-down
	docker compose \
		-f docker-compose.yaml \
		-f resources/providers/docker-compose-providers.yaml \
			up \
				--build \
				--remove-orphans
