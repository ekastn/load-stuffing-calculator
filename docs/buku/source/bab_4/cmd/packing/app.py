"""Flask application for the Packing Service."""

from __future__ import annotations

import os
import time
from typing import Any

from flask import Flask, jsonify, request

from .packing import pack_request
from .schema import ErrorResponse, HealthResponse, PackSuccessResponse


def _json_error(
    status: int,
    code: str,
    message: str,
    *,
    details: dict[str, Any] | None = None,
) -> tuple[Any, int]:
    """Return a JSON error response."""
    payload: ErrorResponse = {
        "success": False,
        "error": {
            "code": code,
            "message": message,
            "details": details or {},
        },
    }
    return jsonify(payload), status


def create_app() -> Flask:
    """Create and configure the Flask application."""
    app = Flask(__name__)

    @app.get("/health")
    def health():
        """Health check endpoint."""
        payload: HealthResponse = {"success": True, "data": {"status": "ok"}}
        return jsonify(payload)

    @app.post("/pack")
    def pack():
        """Pack items into a container."""
        started_at = time.perf_counter()

        # Parse JSON body
        body = request.get_json(silent=True)
        if body is None:
            return _json_error(400, "INVALID_JSON", "Request body must be JSON")

        try:
            result: PackSuccessResponse = pack_request(body)
        except ValueError as e:
            return _json_error(400, "INVALID_REQUEST", str(e))
        except Exception as e:
            return _json_error(
                500, "PACKING_FAILED", "Packing failed", details={"error": str(e)}
            )

        # Add total time to stats
        result["data"]["stats"]["total_time_ms"] = int(
            (time.perf_counter() - started_at) * 1000
        )
        return jsonify(result)

    return app


def main() -> None:
    """Run the Flask development server."""
    app = create_app()

    host = os.environ.get("PACKING_HOST", "0.0.0.0")
    port = int(os.environ.get("PACKING_PORT", "5000"))
    debug = os.environ.get("PACKING_DEBUG", "0") in {"1", "true", "yes"}

    print(f"Starting Packing Service on {host}:{port}")
    app.run(host=host, port=port, debug=debug)


if __name__ == "__main__":
    main()
