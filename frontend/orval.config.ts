import {defineConfig} from "orval";

export default defineConfig({
    roeApi: {
        output: {
            client: 'react-query',
            target: "./src/client.ts",
            schemas: "./src/model",
        },
        input: {
            target: "./document.yaml",
        },
    },
});
