import {defineConfig} from "orval";

export default defineConfig({
    roeApi: {
        output: {
            client: 'react-query',
            target: "./src/client.ts",
            schemas: "./src/model",
            baseUrl: "http://localhost:8000",
        },
        input: {
            target: "./document.yaml",
        },
    },
});
