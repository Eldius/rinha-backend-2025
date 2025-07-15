
IMG_NAME := eldius/rinha-backend-2025

env-down:
	cd resources/providers/ ; docker compose \
		-f docker-compose-providers.yaml \
		down -v
	docker compose \
		-f docker-compose.yaml \
		down -v


run-local: env-down
	docker compose \
		-f resources/providers/docker-compose-providers.yaml \
			up -d \
			--force-recreate \
			--build
	docker compose \
		-f docker-compose.yaml \
			up \
				--build \
				--remove-orphans \
				-d

build-image:
	$(eval IMG_VER=$(shell git rev-parse --short HEAD))
	@echo "IMG: $(IMG_NAME):$(IMG_VER)"
	docker buildx build --tag "$(IMG_NAME):$(IMG_VER)" --tag "$(IMG_NAME):dev" ./
