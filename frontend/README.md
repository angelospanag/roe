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

- [Bun](https://bun.com/)

```bash
brew install oven-sh/bun/bun
````

### Create a `.env` file in the root of the project with the following content:

```dotenv
VITE_API_BASE_URL=http://localhost:8000
```

### Install dependencies

```bash
bun install
```

### Run development server

```bash
bun dev
```
