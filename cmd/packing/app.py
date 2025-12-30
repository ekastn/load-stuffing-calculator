#!/usr/bin/env python3

from __future__ import annotations

import os
import time
from typing import Any

from flask import Flask, jsonify, request

try:
    from .packing import pack_request
    from .schema import ErrorResponse, HealthResponse, PackSuccessResponse
except ImportError:  # pragma: no cover
    # Allow running as a script: `python cmd/packing/app.py`
    from packing import pack_request
    from schema import ErrorResponse, HealthResponse, PackSuccessResponse


def _json_error(
    status: int,
    code: str,
    message: str,
    *,
    details: dict[str, Any] | None = None,
) -> tuple[Any, int]:
    payload: ErrorResponse = {
        "success": False,
        "error": {"code": code, "message": message, "details": details or {}},  # type: ignore[typeddict-item]
    }
    return jsonify(payload), status


def create_app() -> Flask:
    app = Flask(__name__)


    @app.get("/health")
    def health():
        payload: HealthResponse = {"success": True, "data": {"status": "ok"}}
        return jsonify(payload)

    @app.post("/pack")
    def pack():
        started_at = time.perf_counter()

        body = request.get_json(silent=True)
        if body is None:
            return _json_error(400, "INVALID_JSON", "Request body must be JSON")

        try:
            result: PackSuccessResponse = pack_request(body)
        except ValueError as e:
            return _json_error(400, "INVALID_REQUEST", str(e))
        except Exception as e:  # noqa: BLE001
            return _json_error(500, "PACKING_FAILED", "Packing failed", details={"error": str(e)})

        result["data"]["stats"]["total_time_ms"] = int((time.perf_counter() - started_at) * 1000)
        return jsonify(result)

    return app


def main() -> None:
    app = create_app()

    host = os.environ.get("PACKING_HOST", "0.0.0.0")
    port = int(os.environ.get("PACKING_PORT", "5051"))
    debug = os.environ.get("PACKING_DEBUG", "0") in {"1", "true", "TRUE", "yes", "YES"}

    app.run(host=host, port=port, debug=debug)


if __name__ == "__main__":
    main()
