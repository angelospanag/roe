# Roe UI

<!-- TOC -->
* [Roe UI](#roe-ui)
  * [Prerequisites](#prerequisites)
    * [Install runtimes](#install-runtimes)
    * [Create a `.env` file in the root of the project with the following content:](#create-a-env-file-in-the-root-of-the-project-with-the-following-content)
    * [Install dependencies](#install-dependencies)
    * [Run development server](#run-development-server)
<!-- TOC -->

## Prerequisites

### Install runtimes

- [Node.js 22](https://nodejs.org/en/)
- [pnpm](https://pnpm.io/)

```bash
brew install node@22 pnpm
````

### Create a `.env` file in the root of the project with the following content:

```dotenv
VITE_API_BASE_URL=http://localhost:8000
```

### Install dependencies

```bash
pnpm install
```

### Run development server

```bash
pnpm run dev
```
