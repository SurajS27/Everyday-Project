FROM node:18-alpine

WORKDIR /app

# Copy package.json and install dependencies
COPY package.json ./
RUN npm install --production

# Copy the application source code
COPY src/ ./src/

# Expose the port the Express app runs on
EXPOSE 3000

# Start the application
CMD ["npm", "start"]