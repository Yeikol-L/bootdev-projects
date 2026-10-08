FROM debian:stable-slim
COPY main.py main.py
RUN apt update
RUN apt install -y python3
COPY books/ books/
CMD ["python3", "main.py"]
