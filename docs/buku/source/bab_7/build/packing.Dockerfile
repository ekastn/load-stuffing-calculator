FROM python:3.11-slim

WORKDIR /app

# Set environment variables
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PYTHONPATH=/app

# Create a non-root user 'app'
RUN useradd -m -r app

# Install dependencies
COPY cmd/packing/requirements.txt /app/cmd/packing/requirements.txt
RUN pip install --no-cache-dir -r /app/cmd/packing/requirements.txt

# Copy application code
COPY cmd/ /app/cmd/

# Change ownership to 'app' user
RUN chown -R app:app /app

# Switch to non-root user
USER app

EXPOSE 5051

# Run with Gunicorn
CMD ["gunicorn", "--workers", "2", "--threads", "4", "--timeout", "60", "--worker-class", "gthread", "--bind", "0.0.0.0:5051", "cmd.packing.app:create_app()"]
