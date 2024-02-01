FROM node:20

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app

COPY package.json /app

RUN npm install

COPY /dist /app

# Install app dependencies

# Expose the port the app runs on
EXPOSE 3000
# Run npm start when the container launches
CMD ["npm", "run", "start"]
