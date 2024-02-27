# Use a suitable base image that has both Golang and Python with CUDA installed
FROM nvidia/cuda:12.3.1-runtime-ubuntu20.04 as cuda

# Install Python and required libraries with CUDA support
FROM cuda as python

# Install Python and PyTorch with CUDA support
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip

# Install PyTorch with CUDA support
RUN pip install torch==1.10.0+cu111 torchvision==0.11.1+cu111 torchaudio==0.10.0+cu111 -f https://download.pytorch.org/whl/torch_stable.html

# Now, combine both Golang and Python images
FROM golang:1.21.6

# Copy Python and PyTorch from the Python image with CUDA
COPY --from=python /usr/local/cuda /usr/local/cuda
COPY --from=python /usr/local/bin /usr/local/bin

# Set up environment variables for Python with CUDA
ENV LD_LIBRARY_PATH="/usr/local/cuda/lib64:${LD_LIBRARY_PATH}"
# Now you can add your Golang code
WORKDIR /go/src/app
COPY . .

# Build your Golang application
RUN go build -o main .

# Define the command to run the Golang application
CMD ["./main"]
