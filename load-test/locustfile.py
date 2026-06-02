"""
Locust load test for MyAnimeAPI.

Serves two purposes:
  1. Load test — measures throughput, latency, and error rate under stress.
  2. Data population — generates realistic metric traffic so Grafana/Prometheus
     dashboards light up during a demo.

Usage (via Docker Compose):
  docker compose --profile load-test up --build
  Open http://localhost:8089 -> set users / spawn rate -> Start swarming.

Usage (standalone, api running locally):
  pip install locust
  locust -f load-test/locustfile.py --host http://localhost:8080
"""

import random
import uuid

from locust import HttpUser, between, task

# Anime IDs to probe. Many will 404 — that is intentional, it exercises the
# 4xx path in the metrics middleware and matches realistic dashboard browsing.
ANIME_IDS = list(range(1, 51))

# Common search terms — mix of partial matches and unlikely-to-hit terms.
SEARCH_TERMS = [
    "naruto",
    "one",
    "attack",
    "dragon",
    "demon",
    "hero",
    "tokyo",
    "samurai",
    "magic",
    "ninja",
]

# A small pool of genre slugs the API knows about (adjust to your seed data).
GENRES = ["action", "adventure", "comedy", "drama", "fantasy", "romance", "sci-fi"]


class AnimeApiUser(HttpUser):
    """Simulates a mix of catalogue browsers and detail viewers."""

    wait_time = between(0.05, 0.3)

    # ---------- High-frequency reads (catalogue/dashboard traffic) ----------

    @task(10)
    def list_animes(self) -> None:
        """Catalogue page — highest weight, mirrors landing-page traffic."""
        self.client.get("/v1/animes")

    @task(6)
    def get_anime_by_id(self) -> None:
        """Detail page — random ID lookup. 404s are expected and useful."""
        anime_id = random.choice(ANIME_IDS)
        self.client.get(f"/v1/animes/{anime_id}", name="/v1/animes/[id]")

    @task(4)
    def search_animes(self) -> None:
        """Search bar traffic."""
        term = random.choice(SEARCH_TERMS)
        self.client.get(f"/v1/animes/search?title={term}", name="/v1/animes/search")

    @task(3)
    def list_by_genre(self) -> None:
        """Genre filter clicks."""
        genre = random.choice(GENRES)
        self.client.get(f"/v1/animes/genre/{genre}", name="/v1/animes/genre/[slug]")

    @task(3)
    def list_genres(self) -> None:
        """Sidebar / filter menu fetch."""
        self.client.get("/v1/genres")

    @task(2)
    def list_tags(self) -> None:
        """Tag cloud fetch."""
        self.client.get("/v1/tags")

    # ---------- Probes — keep the readiness/version metrics populated ----------

    @task(1)
    def get_health(self) -> None:
        self.client.get("/v1/health")

    @task(1)
    def get_version(self) -> None:
        self.client.get("/v1/version")

    # ---------- Low-weight write traffic ----------
    # Auth endpoints are rate-limited to 5/min/IP, so we keep weight=1 and
    # accept the 429s — they still populate the POST + status="429" buckets
    # in the dashboard, which is informative during a demo.

    @task(1)
    def register_user(self) -> None:
        """One-shot registration with a random identity — mostly 201 or 429."""
        suffix = uuid.uuid4().hex[:8]
        payload = {
            "username": f"loaduser_{suffix}",
            "email": f"loaduser_{suffix}@example.com",
            "password": "LoadTest!Pass123",
        }
        self.client.post("/v1/users/register", json=payload, name="/v1/users/register")
