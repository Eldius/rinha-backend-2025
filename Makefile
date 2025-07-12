
env-down:
	docker compose \
		-f docker-compose.yaml \
		-f providers/docker-compose-providers.yaml \
			down


run-local: env-down
	docker compose \
		-f docker-compose.yaml \
		-f providers/docker-compose-providers.yaml \
			up \
				--build \
				--remove-orphans
