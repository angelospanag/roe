import { defineConfig } from "orval";

export default defineConfig({
  petstore: {
    output: {
      // httpClient: "fetch",
      // client: "fetch",
      client: "react-query",
      target: "./src/client.ts",
      schemas: "./src/model",
      baseUrl: "http://localhost:8000",
    },
    input: {
      target: "./document.yaml",
    },
  },
});
