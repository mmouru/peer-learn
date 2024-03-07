FROM golang:1.21.6

# Install Python and pip dependencies
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    python3-dev \
    build-essential \
    unzip \
    zip

# Install PyTorch and torchvision
RUN pip install torch==2.2.1+cpu torchvision==0.17.1+cpu -f https://download.pytorch.org/whl/torch_stable.html --break-system-packages

# Set up environment variables for the Python library path
ENV LD_LIBRARY_PATH="/usr/local/lib:${LD_LIBRARY_PATH}"

# Set the working directory and copy the Go application code
WORKDIR /go/src/app
COPY . .

# Build the Go application
RUN go build -o main .

# Define the command to run the Go application
CMD ["./main", "-sp", "3000"]
